package game

import "testing"

func TestLeftovers(t *testing.T) {
	st := newTestState(newTestShip(0, 0, 0))
	if Leftovers(st) {
		t.Error("Leftovers should be false with no objects in flight")
	}
	st.Objects = append(st.Objects, &SpaceObject{ID: 1, Type: ObjectTorpedo})
	if !Leftovers(st) {
		t.Error("Leftovers should be true with an object in flight")
	}
}

func TestDisposition_PlayContinues(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	enemy := newTestShip(1, 100, 100)
	st := newTestState(fed, enemy)

	Disposition(st)
	if st.PendingFinal != -1 || st.PendingWarn != -1 {
		t.Errorf("expected no pending condition, got Final=%d Warn=%d", st.PendingFinal, st.PendingWarn)
	}
}

func TestDisposition_FedLose(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	fed.Status[SysDead] = FullyDamaged
	enemy := newTestShip(1, 100, 100)
	st := newTestState(fed, enemy)

	Disposition(st)
	if st.PendingFinal != FinFedLose {
		t.Errorf("PendingFinal = %d, want FinFedLose", st.PendingFinal)
	}
}

func TestDisposition_FedLose_WaitsForLeftovers(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	fed.Status[SysDead] = FullyDamaged
	enemy := newTestShip(1, 100, 100)
	st := newTestState(fed, enemy)
	st.Objects = append(st.Objects, &SpaceObject{ID: 1, Type: ObjectTorpedo})

	Disposition(st)
	if st.PendingFinal != -1 {
		t.Errorf("PendingFinal should stay -1 while ordnance is in flight, got %d", st.PendingFinal)
	}
	if st.PendingWarn != FinFedLose {
		t.Errorf("PendingWarn = %d, want FinFedLose", st.PendingWarn)
	}
}

func TestDisposition_EnemyLose(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	enemy := newTestShip(1, 100, 100)
	enemy.Status[SysDead] = FullyDamaged
	st := newTestState(fed, enemy)

	Disposition(st)
	if st.PendingFinal != FinEnemyLose {
		t.Errorf("PendingFinal = %d, want FinEnemyLose", st.PendingFinal)
	}
}

func TestDisposition_EnemyDisarmedCountsAsKill(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	enemy := newTestShip(1, 100, 100)
	enemy.Energy = 5 // <= 10
	for i := range enemy.Phasers {
		enemy.Phasers[i].Status |= PhaserDamaged
	}
	for i := range enemy.Tubes {
		enemy.Tubes[i].Status |= TubeDamaged
	}
	st := newTestState(fed, enemy)

	Disposition(st)
	if st.PendingFinal != FinEnemyLose {
		t.Errorf("PendingFinal = %d, want FinEnemyLose (enemy disarmed)", st.PendingFinal)
	}
}

func TestDisposition_BothDead(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	fed.Status[SysDead] = FullyDamaged
	enemy := newTestShip(1, 100, 100)
	enemy.Status[SysDead] = FullyDamaged
	st := newTestState(fed, enemy)
	// Even with leftover ordnance, FIN_COMPLETE is always immediately final.
	st.Objects = append(st.Objects, &SpaceObject{ID: 1, Type: ObjectTorpedo})

	Disposition(st)
	if st.PendingFinal != FinComplete {
		t.Errorf("PendingFinal = %d, want FinComplete", st.PendingFinal)
	}
}

func TestDisposition_ComplementZeroForcesShipDead(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	enemy := newTestShip(1, 100, 100)
	enemy.Complement = 0
	st := newTestState(fed, enemy)

	Disposition(st)
	if !enemy.IsDead(SysDead) {
		t.Error("expected enemy with zero complement to be marked dead")
	}
	if st.PendingFinal != FinEnemyLose {
		t.Errorf("PendingFinal = %d, want FinEnemyLose", st.PendingFinal)
	}
}

func TestDisposition_EnemySurrenderAccepted(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	enemy := newTestShip(1, 100, 100)
	st := newTestState(fed, enemy)
	st.PlayStatus |= StatusEnemySurrender

	Disposition(st)
	if st.PendingFinal != FinEnemySurrender {
		t.Errorf("PendingFinal = %d, want FinEnemySurrender", st.PendingFinal)
	}
}

func TestDisposition_FedSurrenderAccepted(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	fed.Status[SysSurrender] = FullyDamaged
	enemy := newTestShip(1, 100, 100)
	st := newTestState(fed, enemy)

	Disposition(st)
	if st.PendingFinal != FinFedSurrender {
		t.Errorf("PendingFinal = %d, want FinFedSurrender", st.PendingFinal)
	}
}
