package commands

import (
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

func TestPursueSetsCourseTowardTarget(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	target := st.Enemies()[0]
	r := game.NewRand(1)
	msgs := Pursue(st, r, sp, target, 5.0)
	if sp.Target != target {
		t.Errorf("Target = %v, want %v", sp.Target, target)
	}
	if sp.NewWarp != 5.0 {
		t.Errorf("NewWarp = %v, want 5", sp.NewWarp)
	}
	if sp.NewCourse != 0 { // target due east
		t.Errorf("NewCourse = %v, want 0", sp.NewCourse)
	}
	if len(msgs) == 0 {
		t.Error("expected a message")
	}
}

func TestPursueRefusesWhenComputerDead(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	sp.Status[game.SysComputer] = game.FullyDamaged
	target := st.Enemies()[0]
	r := game.NewRand(1)
	Pursue(st, r, sp, target, 5.0)
	if sp.Target != nil {
		t.Error("should not have set target with dead computer")
	}
}

func TestPursueCapsWarpAtMaxSpeed(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	sp.MaxSpeed = 3
	target := st.Enemies()[0]
	r := game.NewRand(1)
	Pursue(st, r, sp, target, 10.0)
	if sp.NewWarp != 3 {
		t.Errorf("NewWarp = %v, want capped at 3", sp.NewWarp)
	}
}

func TestPursueForcesWarpDownWhenWarpDriveDead(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	sp.Status[game.SysWarp] = game.FullyDamaged
	target := st.Enemies()[0]
	r := game.NewRand(1)
	Pursue(st, r, sp, target, 5.0)
	if sp.NewWarp != 1.0 {
		t.Errorf("NewWarp = %v, want 1.0", sp.NewWarp)
	}
}

func TestEludeSetsCourseAwayFromTarget(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	target := st.Enemies()[0]
	r := game.NewRand(1)
	Elude(st, r, sp, target, 5.0)
	if sp.RelativeBear != 180.0 {
		t.Errorf("RelativeBear = %v, want 180", sp.RelativeBear)
	}
	if sp.NewCourse != 180 { // target due east, so away is course 180
		t.Errorf("NewCourse = %v, want 180", sp.NewCourse)
	}
}

func TestHelmSetsManualCourse(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	sp.Target = st.Player() // any non-nil sentinel to verify it gets cleared
	msgs := Helm(st, sp, 90.0, 4.0)
	if sp.NewCourse != 90.0 {
		t.Errorf("NewCourse = %v, want 90", sp.NewCourse)
	}
	if sp.NewWarp != 4.0 {
		t.Errorf("NewWarp = %v, want 4", sp.NewWarp)
	}
	if sp.Target != nil {
		t.Error("Target should be cleared")
	}
	if len(msgs) != 2 {
		t.Errorf("messages = %v, want 2 confirmation messages", msgs)
	}
}

func TestHelmRejectsOutOfRangeCourse(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	msgs := Helm(st, sp, 400.0, 4.0)
	if sp.NewCourse == 400.0 {
		t.Error("should not have set out-of-range course")
	}
	if len(msgs) != 1 {
		t.Errorf("messages = %v, want a single rejection message", msgs)
	}
}
