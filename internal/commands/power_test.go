package commands

import (
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

func TestAlterPower_Valid(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)

	shieldDrains := make([]float64, len(sp.Shields))
	for i := range shieldDrains {
		shieldDrains[i] = 0.5
	}
	phaserDrains := make([]float64, len(sp.Phasers))
	for i := range phaserDrains {
		phaserDrains[i] = float64(game.MinPhaserDrain)
	}

	msgs := AlterPower(st, sp, shieldDrains, phaserDrains)
	if msgs != nil {
		t.Fatalf("expected no messages on success, got %v", msgs)
	}
	for i, sh := range sp.Shields {
		if sh.AttemptDrain != 0.5 {
			t.Errorf("shield %d AttemptDrain = %v, want 0.5", i, sh.AttemptDrain)
		}
	}
}

func TestAlterPower_BadShieldValue(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)

	shieldDrains := make([]float64, len(sp.Shields))
	shieldDrains[0] = 1.5 // out of range
	phaserDrains := make([]float64, len(sp.Phasers))

	msgs := AlterPower(st, sp, shieldDrains, phaserDrains)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 bad-parameter message, got %v", msgs)
	}
}

func TestAlterPower_WrongLength(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)

	msgs := AlterPower(st, sp, []float64{0.5}, make([]float64, len(sp.Phasers)))
	if len(msgs) != 1 {
		t.Fatalf("expected 1 bad-parameter message for mismatched length, got %v", msgs)
	}
}

func TestJettisonEngineering(t *testing.T) {
	sp := newTestShip(0, "Player", 10, 10)
	sp.Course = 90
	sp.Warp = 2.0
	st := newTestState(sp)

	obj, msgs := JettisonEngineering(st, sp)
	if obj == nil {
		t.Fatal("expected a jettisoned SpaceObject")
	}
	if obj.Type != game.ObjectEngineering {
		t.Errorf("obj.Type = %v, want ObjectEngineering", obj.Type)
	}
	if obj.From != sp {
		t.Errorf("obj.From = %v, want sp", obj.From)
	}
	if !sp.IsDead(game.SysEngineering) {
		t.Error("expected engineering to be fully damaged after jettison")
	}
	if !sp.IsDead(game.SysWarp) {
		t.Error("expected warp to be fully damaged after jettison")
	}
	if sp.Energy != 0 {
		t.Errorf("sp.Energy = %v, want 0", sp.Energy)
	}
	if len(msgs) == 0 {
		t.Error("expected a confirmation message")
	}
}

func TestJettisonEngineering_AlreadyJettisoned(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	sp.Status[game.SysEngineering] = game.FullyDamaged

	obj, msgs := JettisonEngineering(st, sp)
	if obj != nil {
		t.Error("expected nil object when already jettisoned")
	}
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %v", msgs)
	}
}

func TestDetonateEngineering_NotJettisonedNoConfirm(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)

	obj, msgs := DetonateEngineering(st, sp, false)
	if obj != nil || msgs != nil {
		t.Errorf("expected no-op, got obj=%v msgs=%v", obj, msgs)
	}
}

func TestDetonateEngineering_NotJettisonedWithConfirm(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)

	obj, msgs := DetonateEngineering(st, sp, true)
	if obj == nil {
		t.Fatal("expected a newly jettisoned object")
	}
	if obj.TimeDelay != 0.0 {
		t.Errorf("obj.TimeDelay = %v, want 0", obj.TimeDelay)
	}
	if len(msgs) == 0 {
		t.Error("expected a confirmation message")
	}
}

func TestDetonateEngineering_AlreadyJettisonedInFlight(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	sp.Status[game.SysEngineering] = game.FullyDamaged
	inFlight := &game.SpaceObject{ID: 1, Type: game.ObjectEngineering, From: sp, TimeDelay: 10.0}
	st.Objects = append(st.Objects, inFlight)

	obj, msgs := DetonateEngineering(st, sp, false)
	if obj != nil {
		t.Error("expected nil newly-jettisoned object since it was already jettisoned")
	}
	if inFlight.TimeDelay != 0.0 {
		t.Errorf("expected in-flight object TimeDelay reset to 0, got %v", inFlight.TimeDelay)
	}
	if len(msgs) == 0 {
		t.Error("expected a confirmation message")
	}
}

func TestDetonateEngineering_AlreadyDetonated(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	sp.Status[game.SysEngineering] = game.FullyDamaged
	// No matching object in st.Objects and jettisoned stays nil (already damaged, no confirm).

	obj, msgs := DetonateEngineering(st, sp, false)
	if obj != nil {
		t.Error("expected nil object")
	}
	if len(msgs) != 1 {
		t.Fatalf("expected 1 'already detonated' message, got %v", msgs)
	}
}

func TestAlterFiringParams_Tubes(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	AlterFiringParams(sp, true, 5, 3, 2, false, 0)
	if sp.TubeLaunchSpd != 5 {
		t.Errorf("TubeLaunchSpd = %v, want 5", sp.TubeLaunchSpd)
	}
	if sp.TubeDelay != 3 {
		t.Errorf("TubeDelay = %v, want 3", sp.TubeDelay)
	}
	if sp.TubeProximity != 2 {
		t.Errorf("TubeProximity = %v, want 2", sp.TubeProximity)
	}
}

func TestAlterFiringParams_OutOfRangeIgnored(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	sp.TubeLaunchSpd = 1
	AlterFiringParams(sp, true, -5, 0, 0, false, 0)
	if sp.TubeLaunchSpd != 1 {
		t.Errorf("TubeLaunchSpd should be unchanged when out of range, got %v", sp.TubeLaunchSpd)
	}
}

func TestAlterFiringParams_Phasers(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	AlterFiringParams(sp, false, 0, 0, 0, true, 75)
	if sp.PhaserFirePct != 75 {
		t.Errorf("PhaserFirePct = %v, want 75", sp.PhaserFirePct)
	}
}
