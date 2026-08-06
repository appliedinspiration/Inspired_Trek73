package game

import (
	"strings"
	"testing"
)

func TestRunTurn_AdvancesAndChecksDisposition(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	fed.Course, fed.NewCourse = 0, 0
	fed.Warp, fed.NewWarp = 1.0, 1.0
	enemy := newTestShip(1, 500000, 0)
	enemy.Course, enemy.NewCourse = 180, 180
	enemy.Warp, enemy.NewWarp = 1.0, 1.0
	st := newTestState(fed, enemy)
	r := NewRand(1)
	hooks := NewCombatHooks(st, r)

	msgs := RunTurn(st, hooks)
	_ = msgs // messages content not asserted; just verifying it runs cleanly.

	if st.PendingFinal != -1 {
		// Not expected to resolve in one turn from this starting
		// configuration, but if it does, that's still a valid outcome
		// worth surfacing rather than silently ignoring.
		t.Logf("battle resolved after one turn: PendingFinal=%d", st.PendingFinal)
	}
}

func TestRunTurn_RemovesDetonatedAndAddsLaunchedObjects(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	enemy := newTestShip(1, 100, 100)
	st := newTestState(fed, enemy)
	r := NewRand(1)
	hooks := NewCombatHooks(st, r)

	// Pre-seed an object that will be "detonated" this turn by
	// directly manipulating hooks bookkeeping, since driving an actual
	// detonation end-to-end depends on many unrelated preconditions.
	stale := &SpaceObject{ID: 99, Type: ObjectTorpedo, TimeDelay: 0}
	st.Objects = append(st.Objects, stale)
	hooks.DetonatedObjects = append(hooks.DetonatedObjects, stale)
	fresh := &SpaceObject{ID: 100, Type: ObjectTorpedo, TimeDelay: 10}
	hooks.LaunchedObjects = append(hooks.LaunchedObjects, fresh)

	RunTurn(st, hooks)

	for _, obj := range st.Objects {
		if obj == stale {
			t.Error("detonated object should have been removed from st.Objects")
		}
	}
	found := false
	for _, obj := range st.Objects {
		if obj == fresh {
			found = true
		}
	}
	if !found {
		t.Error("launched object should have been added to st.Objects")
	}
}

func TestRunTurn_EmitsTorpedoAndDamageMessagesSameTurn(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	fed.Name = "Potempkin"
	fed.Course, fed.NewCourse = 0, 0
	fed.Warp, fed.NewWarp = 1.0, 1.0
	fed.TorpedoFiringDelay = 1
	fed.TubeLaunchSpd = 12
	fed.Tubes[0].Load = 100
	fed.Tubes[0].Status |= TubeFiring
	fed.PhaserFiringDelay = 1
	fed.Phasers[0].Load = 100
	fed.Phasers[0].Status |= PhaserFiring
	fed.PhaserSpread = 45
	fed.PhaserFirePct = 100

	enemy := newTestShip(1, 100, 0)
	enemy.Name = "Meteor"
	enemy.Course, enemy.NewCourse = 180, 180
	enemy.Warp, enemy.NewWarp = 1.0, 1.0
	enemy.Shields[0] = Shield{Eff: 0.0, Drain: 0.0}
	enemy.Shields[1] = Shield{Eff: 0.0, Drain: 0.0}
	enemy.Shields[2] = Shield{Eff: 0.0, Drain: 0.0}
	enemy.Shields[3] = Shield{Eff: 0.0, Drain: 0.0}

	st := newTestState(fed, enemy)
	r := NewRand(1)
	hooks := NewCombatHooks(st, r)

	msgs := RunTurn(st, hooks)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "<<Potempkin frng torpedo") {
		t.Fatalf("expected torpedo fire message this turn, got: %v", msgs)
	}
	if !strings.Contains(joined, "hit ") || !strings.Contains(joined, "on Meteor's shield") {
		t.Fatalf("expected immediate ship damage message this turn, got: %v", msgs)
	}

	nextMsgs := RunTurn(st, hooks)
	nextJoined := strings.Join(nextMsgs, "\n")
	if strings.Contains(nextJoined, "<<Potempkin frng torpedo") {
		t.Fatalf("torpedo fire message leaked into next turn: %v", nextMsgs)
	}
}

func TestRunTurn_DetonatesEachTorpedoOnlyOnce(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	enemy := newTestShip(1, 0, 0)
	enemy.Name = "Millennium Pelican"
	for i := range enemy.Shields {
		enemy.Shields[i] = Shield{Eff: 0.0, Drain: 0.0}
	}

	torp := &SpaceObject{
		ID:        4,
		Type:      ObjectTorpedo,
		From:      fed,
		X:         enemy.X,
		Y:         enemy.Y,
		Course:    0,
		Speed:     0,
		NewSpeed:  0,
		Fuel:      40,
		TimeDelay: SegmentSeconds / 2,
	}
	st := newTestState(fed, enemy)
	st.Objects = append(st.Objects, torp)
	hooks := NewCombatHooks(st, NewRand(1))

	msgs := RunTurn(st, hooks)
	joined := strings.Join(msgs, "\n")
	if strings.Count(joined, ":: torp 4 ::") != 1 {
		t.Fatalf("torpedo detonation should be reported once, got messages: %v", msgs)
	}
	if strings.Count(joined, "Millennium Pelican's shield") != 1 {
		t.Fatalf("torpedo damage should be applied once, got messages: %v", msgs)
	}
}
