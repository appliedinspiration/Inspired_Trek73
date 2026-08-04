package commands

import (
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

func TestLaunchProbeRejectsDeadLauncher(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	sp.Status[game.SysProbe] = game.FullyDamaged
	r := game.NewRand(1)
	obj, msgs := LaunchProbe(st, r, sp, 20, 5, 100, nil, 90)
	if obj != nil {
		t.Error("expected nil probe when launcher destroyed")
	}
	if len(msgs) == 0 {
		t.Error("expected a message")
	}
}

func TestLaunchProbeRejectsInsufficientEnergy(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	sp.Energy = 5 // below MinProbeCharge
	r := game.NewRand(1)
	obj, _ := LaunchProbe(st, r, sp, 20, 5, 100, nil, 90)
	if obj != nil {
		t.Error("expected nil probe with insufficient energy")
	}
}

func TestLaunchProbeSucceedsWithCourse(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	sp.Energy = 100
	sp.Pods = 200
	r := game.NewRand(1)
	obj, msgs := LaunchProbe(st, r, sp, 20, 5, 100, nil, 90)
	if obj == nil {
		t.Fatalf("expected a probe, messages: %v", msgs)
	}
	if obj.Course != 90 {
		t.Errorf("Course = %v, want 90", obj.Course)
	}
	if obj.Fuel != 20 {
		t.Errorf("Fuel = %v, want 20", obj.Fuel)
	}
	if sp.Energy != 80 {
		t.Errorf("Energy = %v, want 80", sp.Energy)
	}
	if sp.Pods != 180 {
		t.Errorf("Pods = %v, want 180", sp.Pods)
	}
	if sp.ProbeLauncherStatus != game.ProbeLaunching {
		t.Errorf("ProbeLauncherStatus = %v, want ProbeLaunching", sp.ProbeLauncherStatus)
	}
}

func TestLaunchProbeSucceedsWithTarget(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	sp.Energy = 100
	target := st.Enemies()[0]
	r := game.NewRand(1)
	obj, _ := LaunchProbe(st, r, sp, 20, 5, 100, target, 0)
	if obj == nil || obj.Target != target {
		t.Fatalf("expected probe targeting %v, got %v", target, obj)
	}
}

func TestLaunchProbeRejectsCloakedTarget(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	sp.Energy = 100
	target := st.Enemies()[0]
	target.Cloaking = game.CloakOn
	r := game.NewRand(1)
	obj, msgs := LaunchProbe(st, r, sp, 20, 5, 100, target, 0)
	if obj != nil {
		t.Error("expected nil probe for cloaked target")
	}
	if len(msgs) == 0 {
		t.Error("expected a message")
	}
}

func TestShipProbesFiltersByOwner(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	ep := st.Enemies()[0]
	p1 := &game.SpaceObject{Type: game.ObjectProbe, From: sp}
	p2 := &game.SpaceObject{Type: game.ObjectProbe, From: ep}
	st.Objects = append(st.Objects, p1, p2)
	probes := ShipProbes(st, sp)
	if len(probes) != 1 || probes[0] != p1 {
		t.Errorf("ShipProbes = %v, want [p1]", probes)
	}
}

func TestDetonateAllProbesZeroesTimers(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	p1 := &game.SpaceObject{TimeDelay: 5}
	p2 := &game.SpaceObject{TimeDelay: 10}
	DetonateAllProbes(sp, []*game.SpaceObject{p1, p2})
	if p1.TimeDelay != 0 || p2.TimeDelay != 0 {
		t.Error("expected both probes' timers zeroed")
	}
	if sp.ProbeLauncherStatus != game.ProbeDetonate {
		t.Error("expected ProbeDetonate status")
	}
}

func TestLockProbeSucceedsOnVisibleTarget(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	target := st.Enemies()[0]
	probe := &game.SpaceObject{}
	r := game.NewRand(1)
	LockProbe(st, r, sp, probe, target)
	if probe.Target != target {
		t.Errorf("Target = %v, want %v", probe.Target, target)
	}
}

func TestLockProbeFailsOnCloakedTarget(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	target := st.Enemies()[0]
	target.Cloaking = game.CloakOn
	probe := &game.SpaceObject{}
	r := game.NewRand(1)
	LockProbe(st, r, sp, probe, target)
	if probe.Target != nil {
		t.Error("should not have locked onto cloaked target")
	}
}

func TestCourseProbeSetsCourseAndClearsTarget(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	probe := &game.SpaceObject{Target: st.Player()}
	CourseProbe(st, sp, probe, 270)
	if probe.Course != 270 {
		t.Errorf("Course = %v, want 270", probe.Course)
	}
	if probe.Target != nil {
		t.Error("Target should be cleared")
	}
}
