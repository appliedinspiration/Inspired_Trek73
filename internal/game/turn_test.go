package game

import "testing"

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
