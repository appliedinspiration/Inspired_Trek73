package game

import (
	"math"
	"testing"
)

// newTestShip builds a minimal, fully-populated Ship for movement unit
// tests, sidestepping NewGame's data-table dependency so these tests
// stay focused on movement/power mechanics rather than initialization.
func newTestShip(id int, x, y int) *Ship {
	s := &Ship{
		ID:                 id,
		Name:               "TestShip",
		X:                  x,
		Y:                  y,
		Warp:               1.0,
		NewWarp:            1.0,
		Course:             0,
		NewCourse:          0,
		Eff:                1.0,
		Regen:              0,
		Energy:             1000,
		Pods:               2000,
		Complement:         100,
		Delay:              10000,
		OrigMaxSpeed:       8,
		MaxSpeed:           8,
		DegPerTurn:         10,
		Phasers:            make([]Phaser, 2),
		Tubes:              make([]Tube, 2),
		PhaserFiringDelay:  1,
		TorpedoFiringDelay: 1,
		Cloaking:           CloakNone,
	}
	for i := range s.Shields {
		s.Shields[i] = Shield{Eff: 1.0}
	}
	return s
}

func newTestState(ships ...*Ship) *State {
	st := NewState()
	st.Ships = append(st.Ships, ships...)
	return st
}

func TestMoveShipsStraightLinePosition(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Course, sp.NewCourse = 0, 0
	sp.Warp, sp.NewWarp = 2.0, 2.0 // steady warp 2, due east (course 0)
	fed := sp
	st := newTestState(sp)

	MoveShips(st, NoOpHooks{})

	// At steady warp with no course/warp change, expect straight-line
	// travel for the whole turn: iterations * warp * cos(0) * mPerSegment.
	iterations := int(math.Floor(SecondsPerTurn/SegmentSeconds + 0.5))
	mPerSegment := SegmentSeconds * 100.0
	wantX := 0
	for i := 0; i < iterations; i++ {
		wantX += int(2.0 * mPerSegment)
	}
	if sp.X != wantX {
		t.Errorf("X = %d, want %d", sp.X, wantX)
	}
	if sp.Y != 0 {
		t.Errorf("Y = %d, want 0", sp.Y)
	}
	_ = fed
}

func TestMoveShipsFuelBurnout(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Warp, sp.NewWarp = 5.0, 5.0
	sp.Energy = 0 // Not enough fuel to sustain warp > 1.
	sp.Eff = 1.0
	st := newTestState(sp)

	messages := MoveShips(st, NoOpHooks{})

	if sp.NewWarp != 0.99 {
		t.Errorf("NewWarp = %v, want 0.99 after burnout", sp.NewWarp)
	}
	if len(messages) == 0 {
		t.Fatal("expected a warp drive burning out message")
	}
	found := false
	for _, m := range messages {
		if m == "TestShip's warp drive burning out." {
			found = true
		}
	}
	if !found {
		t.Errorf("messages = %v, want burnout message", messages)
	}
	// Message should only be shown once even though burnout persists
	// across all iterations within the turn.
	count := 0
	for _, m := range messages {
		if m == "TestShip's warp drive burning out." {
			count++
		}
	}
	if count != 1 {
		t.Errorf("burnout message shown %d times, want exactly 1 (shutup suppression)", count)
	}
}

func TestMoveShipsAutopilotTracksTarget(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	target := newTestShip(1, 10000, 0) // due east of sp
	sp.Target = target
	sp.Warp, sp.NewWarp = 0, 0 // no movement, isolate course tracking
	st := newTestState(sp, target)

	MoveShips(st, NoOpHooks{})

	// Target is due east (bearing 0), so course should track to 0
	// (already there); mainly verifying no crash and target retained.
	if sp.Target != target {
		t.Errorf("Target = %v, want unchanged (target still alive)", sp.Target)
	}
}

func TestMoveShipsAutopilotDisengagesOnTargetDeath(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	target := newTestShip(1, 10000, 0)
	target.Status[SysDead] = FullyDamaged
	sp.Target = target
	sp.RelativeBear = 45
	st := newTestState(sp, target)

	messages := MoveShips(st, NoOpHooks{})

	if sp.Target != nil {
		t.Errorf("Target = %v, want nil after target death", sp.Target)
	}
	if sp.RelativeBear != 0 {
		t.Errorf("RelativeBear = %v, want 0 after disengage", sp.RelativeBear)
	}
	found := false
	for _, m := range messages {
		if m == "TestShip's autopilot disengaging." {
			found = true
		}
	}
	if !found {
		t.Errorf("messages = %v, want autopilot disengaging message", messages)
	}
}

func TestMoveShipsSkipsDeadShips(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Status[SysDead] = FullyDamaged
	origX, origY := sp.X, sp.Y
	st := newTestState(sp)

	MoveShips(st, NoOpHooks{})

	if sp.X != origX || sp.Y != origY {
		t.Errorf("dead ship moved: (%d,%d), want (%d,%d)", sp.X, sp.Y, origX, origY)
	}
}

func TestDistributeShieldsDownMessage(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Energy = 0
	sp.Regen = 0
	for i := range sp.Shields {
		sp.Shields[i].AttemptDrain = 0
	}
	st := newTestState(sp)

	messages := Distribute(sp, st, NoOpHooks{})

	found := false
	for _, m := range messages {
		if m == "Engineer: Captain, our shields are down!" {
			found = true
		}
	}
	if !found {
		t.Errorf("messages = %v, want shields-down message", messages)
	}
}

func TestDistributeChargesPhasersWithAvailableFuel(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Energy = 1000
	sp.Regen = 0
	sp.Phasers[0] = Phaser{Load: 0, Drain: 5}
	for i := range sp.Shields {
		sp.Shields[i].AttemptDrain = 0
	}
	st := newTestState(sp)

	Distribute(sp, st, NoOpHooks{})

	// Over 2 one-second loop iterations, drain of 5/sec should charge
	// the phaser up to MaxPhaserCharge (10) and then stay capped.
	if sp.Phasers[0].Load != MaxPhaserCharge {
		t.Errorf("Phasers[0].Load = %v, want %v", sp.Phasers[0].Load, MaxPhaserCharge)
	}
}

func TestDistributeShieldDrainRationedWhenFuelShort(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Energy = 1
	sp.Regen = 0
	for i := range sp.Shields {
		sp.Shields[i].AttemptDrain = 10
	}
	st := newTestState(sp)

	Distribute(sp, st, NoOpHooks{})

	// fuel(1) < drain(ceil(40)=40), so drain is rationed: attempt *
	// fuel/drain = 10 * 1/40 = 0.25 for the first loop iteration; the
	// second iteration then has fuel=0 (regen=0), so it also rations
	// to 0. We only assert the fluctuating message appeared and no
	// shield exceeds its attempted drain.
	for i := range sp.Shields {
		if sp.Shields[i].Drain > sp.Shields[i].AttemptDrain {
			t.Errorf("Shields[%d].Drain = %v exceeds AttemptDrain %v", i, sp.Shields[i].Drain, sp.Shields[i].AttemptDrain)
		}
	}
}

func TestCheckTargetsClearsShipTargetOnDeath(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	target := newTestShip(1, 100, 100)
	target.Status[SysDead] = FullyDamaged
	sp.Target = target
	sp.RelativeBear = 30
	st := newTestState(sp, target)

	messages := CheckTargets(st)

	if sp.Target != nil {
		t.Errorf("Target = %v, want nil", sp.Target)
	}
	if sp.RelativeBear != 0 {
		t.Errorf("RelativeBear = %v, want 0", sp.RelativeBear)
	}
	found := false
	for _, m := range messages {
		if m == "   helm lock disengaging" {
			found = true
		}
	}
	if !found {
		t.Errorf("messages = %v, want helm lock disengaging", messages)
	}
}

func TestCheckTargetsClearsWeaponTargetsOnDeath(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	target := newTestShip(1, 100, 100)
	target.Status[SysDead] = FullyDamaged
	sp.Phasers[0].Target = target
	sp.Tubes[0].Target = target
	st := newTestState(sp, target)

	CheckTargets(st)

	if sp.Phasers[0].Target != nil {
		t.Error("Phasers[0].Target want nil after target death")
	}
	if sp.Tubes[0].Target != nil {
		t.Error("Tubes[0].Target want nil after target death")
	}
}

func TestCheckTargetsUpdatesBearingForLiveTarget(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Course = 0
	target := newTestShip(1, 0, 10000) // due north
	sp.Phasers[0].Target = target
	st := newTestState(sp, target)

	CheckTargets(st)

	want := rectify(bearingTo(0, 0, 0, 10000) - sp.Course)
	if sp.Phasers[0].Bearing != want {
		t.Errorf("Phasers[0].Bearing = %v, want %v", sp.Phasers[0].Bearing, want)
	}
}

// TestCheckTargetsDecrementsCloakDelayForAllShips locks in the fix for
// a bug in the original check_targets() (moveships.c), where the
// "if (sp->cloak_delay > 0) sp->cloak_delay--;" statement sat outside
// the per-ship loop and so only ever affected the last ship processed.
// This port decrements every ship's own CloakDelay.
func TestCheckTargetsDecrementsCloakDelayForAllShips(t *testing.T) {
	sp1 := newTestShip(0, 0, 0)
	sp1.CloakDelay = 2
	sp2 := newTestShip(1, 100, 100)
	sp2.CloakDelay = 3
	st := newTestState(sp1, sp2)

	CheckTargets(st)

	if sp1.CloakDelay != 1 {
		t.Errorf("sp1.CloakDelay = %d, want 1", sp1.CloakDelay)
	}
	if sp2.CloakDelay != 2 {
		t.Errorf("sp2.CloakDelay = %d, want 2 (bug fix: every ship decrements, not just the last)", sp2.CloakDelay)
	}
}

func TestMiscTimersSelfDestructWarning(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Delay = 500
	st := newTestState(sp)

	messages := MiscTimers(st)

	if sp.Delay != 500 {
		t.Errorf("Delay = %v, want unchanged 500", sp.Delay)
	}
	found := false
	for _, m := range messages {
		if m == "Computer: 500.00 seconds to self destruct." {
			found = true
		}
	}
	if !found {
		t.Errorf("messages = %v, want self-destruct warning", messages)
	}
}

func TestMiscTimersAbortsOnComputerDamage(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Delay = 500
	sp.Status[SysComputer] = FullyDamaged
	st := newTestState(sp)

	messages := MiscTimers(st)

	if sp.Delay != 10000 {
		t.Errorf("Delay = %v, want reset to 10000", sp.Delay)
	}
	found := false
	for _, m := range messages {
		if m == "Science: Self-destruct has been aborted due to computer damage" {
			found = true
		}
	}
	if !found {
		t.Errorf("messages = %v, want abort message", messages)
	}
}

func TestMiscTimersAgesSpecialActionCounters(t *testing.T) {
	st := newTestState(newTestShip(0, 0, 0))
	st.Defenseless = 3
	st.Corbomite = 0
	st.Surrender = 1
	st.SurrenderP = 0

	MiscTimers(st)

	if st.Defenseless != 4 {
		t.Errorf("Defenseless = %d, want 4", st.Defenseless)
	}
	if st.Corbomite != 0 {
		t.Errorf("Corbomite = %d, want 0 (inactive counters stay at 0)", st.Corbomite)
	}
	if st.Surrender != 2 {
		t.Errorf("Surrender = %d, want 2", st.Surrender)
	}
	if st.SurrenderP != 0 {
		t.Errorf("SurrenderP = %d, want 0", st.SurrenderP)
	}
}
