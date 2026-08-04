package commands

import (
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

func TestFirePhasersSetsFiringFlagAndSpread(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	msgs := FirePhasers(sp, "13", 20)
	if sp.PhaserSpread != 20 {
		t.Errorf("PhaserSpread = %d, want 20", sp.PhaserSpread)
	}
	if sp.Phasers[0].Status&game.PhaserFiring == 0 {
		t.Error("bank 1 not firing")
	}
	if sp.Phasers[1].Status&game.PhaserFiring != 0 {
		t.Error("bank 2 should not be firing")
	}
	if sp.Phasers[2].Status&game.PhaserFiring == 0 {
		t.Error("bank 3 not firing")
	}
	if len(msgs) != 0 {
		t.Errorf("messages = %v, want none", msgs)
	}
}

func TestFirePhasersAllSelector(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	FirePhasers(sp, "all", 10)
	for i := range sp.Phasers {
		if sp.Phasers[i].Status&game.PhaserFiring == 0 {
			t.Errorf("bank %d not firing", i)
		}
	}
}

func TestFirePhasersRejectsOutOfRangeSpread(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	msgs := FirePhasers(sp, "1", 5) // below MinPhaserSpread
	if sp.Phasers[0].Status&game.PhaserFiring != 0 {
		t.Error("should not have fired with bad spread")
	}
	if msgs != nil {
		t.Errorf("messages = %v, want nil", msgs)
	}
}

func TestFirePhasersSkipsDamagedBank(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	sp.Phasers[0].Status = game.PhaserDamaged
	msgs := FirePhasers(sp, "1", 20)
	if sp.Phasers[0].Status&game.PhaserFiring != 0 {
		t.Error("damaged bank should not fire")
	}
	if len(msgs) == 0 {
		t.Error("expected a damage-report message")
	}
}

func TestFireTubesReportsUnloadedTube(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	msgs := FireTubes(sp, "1")
	if sp.Tubes[0].Status&game.TubeFiring == 0 {
		t.Error("tube 1 should be firing")
	}
	found := false
	for _, m := range msgs {
		if m == "Computer: Tube(s) 1 is not loaded." {
			found = true
		}
	}
	if !found {
		t.Errorf("messages = %v, want unloaded-tube message", msgs)
	}
}

func TestFireTubesMultipleUnloaded(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	msgs := FireTubes(sp, "all")
	found := false
	for _, m := range msgs {
		if m == "Computer: Tube(s) 1, 2, 3 are not loaded." {
			found = true
		}
	}
	if !found {
		t.Errorf("messages = %v, want all-tubes-unloaded message", msgs)
	}
}

func TestLockPhasersRefusesWhenComputerDead(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	sp.Status[game.SysComputer] = game.FullyDamaged
	target := newTestShip(1, "Klingon", 100, 0)
	r := game.NewRand(1)
	msgs := LockPhasers(st, r, sp, "1", target)
	if len(msgs) == 0 {
		t.Fatal("expected a computer-dead message")
	}
	if sp.Phasers[0].Target != nil {
		t.Error("should not have locked with dead computer")
	}
}

func TestLockPhasersLocksOntoVisibleTarget(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	target := newTestShip(1, "Klingon", 100, 0)
	r := game.NewRand(1)
	// Force SystemWorks to succeed by ensuring no damage.
	msgs := LockPhasers(st, r, sp, "12", target)
	if sp.Phasers[0].Target != target || sp.Phasers[1].Target != target {
		t.Errorf("banks 1,2 target = %v,%v, want %v", sp.Phasers[0].Target, sp.Phasers[1].Target, target)
	}
	if sp.Phasers[2].Target == target {
		t.Error("bank 3 should not be targeted")
	}
	_ = msgs
}

func TestLockPhasersRefusesCloakedTarget(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	target := newTestShip(1, "Klingon", 100, 0)
	target.Cloaking = game.CloakOn
	r := game.NewRand(1)
	msgs := LockPhasers(st, r, sp, "1", target)
	if len(msgs) == 0 {
		t.Fatal("expected an unable-to-lock message")
	}
	if sp.Phasers[0].Target != nil {
		t.Error("should not have locked onto cloaked target")
	}
}

func TestLockTubesLocksOntoTarget(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	target := newTestShip(1, "Klingon", 100, 0)
	r := game.NewRand(1)
	LockTubes(st, r, sp, "all", target)
	for i := range sp.Tubes {
		if sp.Tubes[i].Target != target {
			t.Errorf("tube %d target = %v, want %v", i, sp.Tubes[i].Target, target)
		}
	}
}

func TestTurnPhasersUnlocksAndSetsBearing(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	sp.Phasers[0].Target = newTestShip(1, "Klingon", 100, 0)
	TurnPhasers(sp, "1", 45.0)
	if sp.Phasers[0].Target != nil {
		t.Error("target should be cleared")
	}
	if sp.Phasers[0].Bearing != 45.0 {
		t.Errorf("Bearing = %v, want 45", sp.Phasers[0].Bearing)
	}
}

func TestTurnPhasersRejectsOutOfRangeBearing(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	msgs := TurnPhasers(sp, "1", 400.0)
	if sp.Phasers[0].Bearing != 0 {
		t.Error("bearing should not have changed")
	}
	if msgs != nil {
		t.Errorf("messages = %v, want nil", msgs)
	}
}

func TestTurnTubesUnlocksAndSetsBearing(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	sp.Tubes[0].Target = newTestShip(1, "Klingon", 100, 0)
	TurnTubes(sp, "1", 90.0)
	if sp.Tubes[0].Target != nil {
		t.Error("target should be cleared")
	}
	if sp.Tubes[0].Bearing != 90.0 {
		t.Errorf("Bearing = %v, want 90", sp.Tubes[0].Bearing)
	}
}

func TestLoadTubesLoadsFromEnergy(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	sp.Energy = 100
	sp.Pods = 200
	LoadTubes(st, sp, true, "1")
	if sp.Tubes[0].Load != game.MaxTubeCharge {
		t.Errorf("Tube[0].Load = %v, want %v", sp.Tubes[0].Load, game.MaxTubeCharge)
	}
	if sp.Energy != 100-game.MaxTubeCharge {
		t.Errorf("Energy = %v, want %v", sp.Energy, 100-game.MaxTubeCharge)
	}
}

func TestLoadTubesUnloadsToEnergy(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	sp.Tubes[0].Load = 5
	sp.Energy = 100
	sp.Pods = 200
	LoadTubes(st, sp, false, "1")
	if sp.Tubes[0].Load != 0 {
		t.Errorf("Tube[0].Load = %v, want 0", sp.Tubes[0].Load)
	}
	if sp.Energy != 105 {
		t.Errorf("Energy = %v, want 105", sp.Energy)
	}
}

func TestLoadTubesSkipsDamagedTube(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	sp.Tubes[0].Status = game.TubeDamaged
	sp.Energy = 100
	LoadTubes(st, sp, true, "all")
	if sp.Tubes[0].Load != 0 {
		t.Error("damaged tube should not have loaded")
	}
	if sp.Tubes[1].Load != game.MaxTubeCharge {
		t.Error("undamaged tube should have loaded")
	}
}
