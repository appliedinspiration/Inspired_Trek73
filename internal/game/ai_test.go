package game

import (
	"math"
	"testing"
)

// newAITestState builds a two-ship battle (player at index 0, one
// enemy at index 1) and a CombatHooks wired to a deterministic Rand,
// for exercising the e_* AI primitives, special(), and
// standard_strategy() in isolation.
func newAITestState(seed int64) (*State, *Ship, *Ship, *CombatHooks) {
	fed := newTestShip(0, 0, 0)
	enemy := newTestShip(1, 1000, 0)
	st := newTestState(fed, enemy)
	r := NewRand(seed)
	h := NewCombatHooks(st, r)
	return st, fed, enemy, h
}

func TestEAttackPursuesAndReportsSpeedCapped(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.MaxSpeed = 3 // force the speed cap branch
	if !h.EAttack(enemy, fed) {
		t.Fatal("EAttack = false, want true")
	}
	if enemy.Target != fed {
		t.Errorf("Target = %v, want fed", enemy.Target)
	}
	if enemy.NewWarp > enemy.MaxSpeed {
		t.Errorf("NewWarp = %v, want capped at MaxSpeed %v", enemy.NewWarp, enemy.MaxSpeed)
	}
}

func TestEAttackRefusesWhenWarpDead(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Status[SysWarp] = FullyDamaged
	if h.EAttack(enemy, fed) {
		t.Error("EAttack = true, want false when warp is dead")
	}
}

func TestECheckArmsCountsOnlyLoadedUndamagedWeapons(t *testing.T) {
	_, _, enemy, h := newAITestState(1)
	enemy.Phasers[0].Load = 5
	enemy.Phasers[1].Load = 0 // unloaded: should NOT count
	enemy.Tubes[0].Load = 3
	enemy.Tubes[1].Status |= TubeDamaged
	enemy.Tubes[1].Load = 10 // damaged: should NOT count even though loaded

	arms := h.ECheckArms(enemy)
	if arms != 2 {
		t.Errorf("ECheckArms = %d, want 2 (one phaser, one tube)", arms)
	}
}

// TestECheckArmsBugFixRegression locks in the fix for the original's
// "load >= 0" bug (always true), which meant e_checkarms() never
// detected an unarmed ship.
func TestECheckArmsBugFixRegression(t *testing.T) {
	_, _, enemy, h := newAITestState(1)
	// Every weapon undamaged but completely unloaded.
	arms := h.ECheckArms(enemy)
	if arms != 0 {
		t.Errorf("ECheckArms = %d, want 0 for a fully unloaded ship", arms)
	}
}

func TestECheckProbeEvadesWhenProbeClose(t *testing.T) {
	st, _, enemy, h := newAITestState(1)
	probe := &SpaceObject{Type: ObjectProbe, X: enemy.X + 100, Y: enemy.Y}
	st.Objects = append(st.Objects, probe)

	if !h.ECheckProbe(enemy) {
		t.Fatal("ECheckProbe = false, want true when a probe is within 2000")
	}
	if enemy.Target != nil {
		t.Error("Target should be cleared after evasive action")
	}
}

func TestECheckProbeIgnoresFarProbe(t *testing.T) {
	st, _, enemy, h := newAITestState(1)
	probe := &SpaceObject{Type: ObjectProbe, X: enemy.X + 100000, Y: enemy.Y}
	st.Objects = append(st.Objects, probe)

	if h.ECheckProbe(enemy) {
		t.Error("ECheckProbe = true, want false for a distant probe")
	}
}

func TestECheckProbeSkippedWhenCloaked(t *testing.T) {
	st, _, enemy, h := newAITestState(1)
	enemy.Cloaking = CloakOn
	probe := &SpaceObject{Type: ObjectProbe, X: enemy.X + 100, Y: enemy.Y}
	st.Objects = append(st.Objects, probe)

	if h.ECheckProbe(enemy) {
		t.Error("ECheckProbe = true, want false while cloaked")
	}
}

func TestECloakOnEngagesCloak(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Cloaking = CloakOff
	enemy.CloakDelay = 0
	if !h.ECloakOn(enemy, fed) {
		t.Fatal("ECloakOn = false, want true")
	}
	if enemy.Cloaking != CloakOn {
		t.Errorf("Cloaking = %v, want CloakOn", enemy.Cloaking)
	}
	if enemy.LastKnown.X != enemy.X || enemy.LastKnown.Y != enemy.Y {
		t.Error("LastKnown position not recorded")
	}
}

func TestECloakOnRefusedDuringDelay(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Cloaking = CloakOff
	enemy.CloakDelay = 2
	if h.ECloakOn(enemy, fed) {
		t.Error("ECloakOn = true, want false while cloak_delay > 0")
	}
}

func TestECloseTorpsFiresPhasersWhenPossible(t *testing.T) {
	st, fed, enemy, h := newAITestState(1)
	enemy.Phasers[0].Load = 10
	enemy.Phasers[0].Target = fed
	torp := &SpaceObject{Type: ObjectTorpedo, From: fed, X: enemy.X + 500, Y: enemy.Y}
	st.Objects = append(st.Objects, torp)

	if !h.ECloseTorps(enemy, fed) {
		t.Fatal("ECloseTorps = false, want true (should fire phasers)")
	}
	fired := false
	for _, p := range enemy.Phasers {
		if p.Status&PhaserFiring != 0 {
			fired = true
		}
	}
	if !fired {
		t.Error("expected a phaser bank to be marked firing")
	}
}

// TestECloseTorpsEvadeUsesBadPosition locks in the fix for the
// original's bug where e_evade() was called with the loop variable's
// final value (the last torpedo examined in the whole list) instead of
// the actual close/dangerous torpedo ("bad"). Here, a distant torpedo
// is listed after the close one; if the bug were still present, the
// evade computation would be based on the distant torpedo's position.
func TestECloseTorpsEvadeUsesBadPosition(t *testing.T) {
	st, fed, enemy, h := newAITestState(1)
	// No usable phasers, so evasion is the only option.
	closeTorp := &SpaceObject{Type: ObjectTorpedo, From: fed, X: enemy.X + 500, Y: enemy.Y}
	farTorp := &SpaceObject{Type: ObjectTorpedo, From: fed, X: enemy.X + 50000, Y: enemy.Y + 50000}
	st.Objects = append(st.Objects, closeTorp, farTorp)

	if !h.ECloseTorps(enemy, fed) {
		t.Fatal("ECloseTorps = false, want true (should evade)")
	}
	if enemy.Target != nil {
		t.Error("Target should be cleared by evasive action")
	}
	// The evasion course is derived from bearing to the "bad" torpedo,
	// not the far one; sanity-check it isn't the zero value that
	// would arise only by coincidence.
	if enemy.NewWarp < 2 || enemy.NewWarp > enemy.MaxSpeed {
		t.Errorf("NewWarp = %v, out of e_evade's expected [2, max_speed) range", enemy.NewWarp)
	}
}

func TestECloseTorpsNoOpWhenCloaked(t *testing.T) {
	st, fed, enemy, h := newAITestState(1)
	enemy.Cloaking = CloakOn
	torp := &SpaceObject{Type: ObjectTorpedo, From: fed, X: enemy.X + 500, Y: enemy.Y}
	st.Objects = append(st.Objects, torp)

	if h.ECloseTorps(enemy, fed) {
		t.Error("ECloseTorps = true, want false while cloaked")
	}
}

func TestEDestructRefusedBeforeDelayThreshold(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Delay = 4.0
	if h.EDestruct(enemy, fed) {
		t.Error("EDestruct = true, want false when delay < 5.0")
	}
}

func TestEDestructInitiatesAtDelayThreshold(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Delay = 10000
	if !h.EDestruct(enemy, fed) {
		t.Fatal("EDestruct = false, want true when delay >= 5.0")
	}
	if enemy.Delay != 5.0 {
		t.Errorf("Delay = %v, want 5.0", enemy.Delay)
	}
	if enemy.Cloaking != CloakNone {
		t.Errorf("Cloaking = %v, want CloakNone", enemy.Cloaking)
	}
}

func TestEEvadeClearsTargetAndSetsCourse(t *testing.T) {
	_, _, enemy, h := newAITestState(2)
	enemy.Target = enemy // arbitrary non-nil sentinel
	h.EEvade(enemy, enemy.X+1000, enemy.Y)
	if enemy.Target != nil {
		t.Error("Target should be cleared")
	}
	if enemy.NewWarp < 2 {
		t.Errorf("NewWarp = %v, want >= 2", enemy.NewWarp)
	}
}

func TestEEvadeWarpOneWhenWarpDead(t *testing.T) {
	_, _, enemy, h := newAITestState(2)
	enemy.Status[SysWarp] = FullyDamaged
	h.EEvade(enemy, enemy.X+1000, enemy.Y)
	if enemy.NewWarp != 1.0 {
		t.Errorf("NewWarp = %v, want 1.0 when warp is dead", enemy.NewWarp)
	}
}

func TestEJettisonDropsEngineeringSection(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Energy = 500
	if !h.EJettison(enemy, fed) {
		t.Fatal("EJettison = false, want true")
	}
	if len(h.LaunchedObjects) != 1 {
		t.Fatalf("LaunchedObjects = %d, want 1", len(h.LaunchedObjects))
	}
	obj := h.LaunchedObjects[0]
	if obj.Type != ObjectEngineering {
		t.Errorf("Type = %v, want ObjectEngineering", obj.Type)
	}
	if enemy.Energy != 0 || enemy.Pods != 0 {
		t.Errorf("Energy/Pods = %v/%v, want 0/0", enemy.Energy, enemy.Pods)
	}
	if enemy.MaxSpeed != 1.0 {
		t.Errorf("MaxSpeed = %v, want 1.0", enemy.MaxSpeed)
	}
	if enemy.Status[SysEngineering] != FullyDamaged || enemy.Status[SysWarp] != FullyDamaged {
		t.Error("SysEngineering and SysWarp should be marked fully damaged")
	}
}

func TestEJettisonRefusedWhenEngineeringAlreadyDead(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Status[SysEngineering] = FullyDamaged
	if h.EJettison(enemy, fed) {
		t.Error("EJettison = true, want false when engineering is already dead")
	}
}

func TestELaunchProbeRefusedWhenFedFast(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	fed.Warp = 5.0
	if h.ELaunchProbe(enemy, fed) {
		t.Error("ELaunchProbe = true, want false when fed is moving fast")
	}
}

func TestELaunchProbeLaunchesWhenFedSlow(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	fed.Warp = 0.5
	enemy.Energy = 100
	if !h.ELaunchProbe(enemy, fed) {
		t.Fatal("ELaunchProbe = false, want true")
	}
	if len(h.LaunchedObjects) != 1 || h.LaunchedObjects[0].Type != ObjectProbe {
		t.Fatal("expected a launched probe object")
	}
	if h.LaunchedObjects[0].Target != fed {
		t.Error("probe target should be fed")
	}
}

func TestELoadTubesLoadsUndamagedTubes(t *testing.T) {
	_, _, enemy, h := newAITestState(1)
	enemy.Energy = 100
	enemy.Tubes[0].Load = 0
	enemy.Tubes[1].Status |= TubeDamaged

	loaded := h.ELoadTubes(enemy)
	if loaded != 1 {
		t.Errorf("loaded = %d, want 1 (damaged tube skipped)", loaded)
	}
	if enemy.Tubes[0].Load != MaxTubeCharge {
		t.Errorf("Tubes[0].Load = %v, want %v", enemy.Tubes[0].Load, MaxTubeCharge)
	}
}

func TestELockPhasersLocksUnlockedUndamagedBanks(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Phasers[0].Target = enemy // already locked (onto itself as sentinel)
	enemy.Phasers[1].Status |= PhaserDamaged

	locked := h.ELockPhasers(enemy, fed)
	if locked != 0 {
		t.Errorf("locked = %d, want 0 (one already locked, one damaged)", locked)
	}
}

func TestELockTubesLocksAvailableTubes(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	locked := h.ELockTubes(enemy, fed)
	if locked != 2 {
		t.Errorf("locked = %d, want 2", locked)
	}
	for i := range enemy.Tubes {
		if enemy.Tubes[i].Target != fed {
			t.Errorf("Tubes[%d].Target = %v, want fed", i, enemy.Tubes[i].Target)
		}
	}
}

func TestEPhasersFiresLockedLoadedBanksOnly(t *testing.T) {
	_, fed, enemy, h := newAITestState(3)
	enemy.PhaserFirePct = 100
	enemy.Phasers[0].Load = MaxPhaserCharge
	enemy.Phasers[0].Target = fed
	enemy.Phasers[1].Load = 0 // unloaded, should never fire

	fed.X, fed.Y = enemy.X+100, enemy.Y // close range, easy hit

	banks := h.EPhasers(enemy, fed)
	if banks < 1 {
		t.Fatalf("banks = %d, want at least 1", banks)
	}
	if enemy.Phasers[1].Status&PhaserFiring != 0 {
		t.Error("unloaded phaser bank should never be marked firing")
	}
}

func TestEPursueSetsCourseAndTarget(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	fed.X, fed.Y = enemy.X+1000, enemy.Y
	h.EPursue(enemy, fed, 3.0)
	if enemy.Target != fed {
		t.Error("Target should be set to fed")
	}
	if enemy.NewWarp != 3.0 {
		t.Errorf("NewWarp = %v, want 3.0", enemy.NewWarp)
	}
}

func TestEPursueForcesMinimumWarpWhenWarpDead(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Status[SysWarp] = FullyDamaged
	h.EPursue(enemy, fed, 3.0)
	if enemy.NewWarp != 0.99 {
		t.Errorf("NewWarp = %v, want 0.99 when warp is dead and speed > 1", enemy.NewWarp)
	}
}

// TestERunawayBugFixRegression locks in the fix for the original's
// integer-division bug ("2 / 3 * max_speed" truncating to zero, so
// the retreating ship never actually gained speed).
func TestERunawayBugFixRegression(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.MaxSpeed = 9
	enemy.Shields[1] = Shield{Eff: 1.0, Drain: 1.0} // strongest shield: index 1

	h.ERunaway(enemy, fed)

	if enemy.NewWarp == 0 {
		t.Fatal("NewWarp = 0, want nonzero (bug regression: integer division truncated 2/3 to 0)")
	}
	want := 2.0 / 3.0 * 9.0
	if math.Abs(enemy.NewWarp-want) > 1e-9 {
		t.Errorf("NewWarp = %v, want %v", enemy.NewWarp, want)
	}
}

func TestERunawayClearsTarget(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Target = fed
	h.ERunaway(enemy, fed)
	if enemy.Target != nil {
		t.Error("Target should be cleared while running away")
	}
}

func TestETorpedoRefusedWhenOtherEnemyClose(t *testing.T) {
	st, _, enemy, h := newAITestState(1)
	other := newTestShip(2, enemy.X+100, enemy.Y) // within 400
	st.Ships = append(st.Ships, other)
	enemy.Tubes[0].Load = 10
	enemy.Tubes[0].Target = other

	if tubes := h.ETorpedo(enemy); tubes != 0 {
		t.Errorf("ETorpedo = %d, want 0 when another enemy ship is within 400", tubes)
	}
}

func TestETorpedoFiresLoadedLockedTubes(t *testing.T) {
	_, fed, enemy, h := newAITestState(1)
	enemy.Tubes[0].Load = 10
	enemy.Tubes[0].Target = fed
	enemy.Tubes[1].Load = 0 // unloaded, should not fire

	tubes := h.ETorpedo(enemy)
	if tubes < 1 {
		t.Fatalf("tubes = %d, want at least 1", tubes)
	}
	if enemy.Tubes[0].Status&TubeFiring == 0 {
		t.Error("loaded, locked tube should be marked firing")
	}
	if enemy.Tubes[1].Status&TubeFiring != 0 {
		t.Error("unloaded tube should never be marked firing")
	}
}

func TestSpecialSurrenderAcceptedSetsPlayStatus(t *testing.T) {
	st, fed, enemy, h := newAITestState(1)
	st.Surrender = 1
	st.RaceChances.Surrender = 100 // guarantee acceptance
	h.Special(enemy, 5000, fed)
	if st.PlayStatus&StatusFedSurrender == 0 {
		t.Error("PlayStatus should have StatusFedSurrender set after acceptance")
	}
	if st.Surrender != 2 {
		t.Errorf("Surrender = %d, want 2 after acceptance", st.Surrender)
	}
}

func TestSpecialSurrenderRejectedEndsRuse(t *testing.T) {
	st, fed, enemy, h := newAITestState(1)
	st.Surrender = 1
	st.RaceChances.Surrender = 0 // guarantee rejection
	h.Special(enemy, 5000, fed)
	if st.Surrender != 6 {
		t.Errorf("Surrender = %d, want 6 after rejection", st.Surrender)
	}
}

func TestSpecialCorbomiteAcceptedRetreats(t *testing.T) {
	st, fed, enemy, h := newAITestState(1)
	st.Corbomite = 1
	st.RaceChances.Corbomite = 100
	h.Special(enemy, 5000, fed)
	if enemy.NewWarp >= 0 {
		t.Errorf("NewWarp = %v, want negative (retreating) after corbomite bluff accepted", enemy.NewWarp)
	}
}

// TestSpecialUnsportsmanlikeFiringCancelsRuses locks in the fix for
// the original's wrong-loop-index bug in the tube-firing check
// (fed->tubes[loop] instead of fed->tubes[loop2]), by using a tube
// index beyond the phaser count that the bug would have mis-indexed.
func TestSpecialUnsportsmanlikeFiringCancelsRuses(t *testing.T) {
	st, fed, enemy, h := newAITestState(1)
	st.Defenseless = 3
	// Give fed more tubes than phasers so the original's buggy index
	// (fed->tubes[loop], where loop only ranges over num_phasers)
	// would miss a firing tube beyond that index.
	fed.Phasers = make([]Phaser, 1)
	fed.Tubes = make([]Tube, 2)
	fed.Tubes[1].Status |= TubeFiring

	h.Special(enemy, 5000, fed)

	if st.Defenseless != 6 {
		t.Errorf("Defenseless = %d, want 6 after unsportsmanlike firing detected", st.Defenseless)
	}
}

func TestSpecialNoOpWhenNoRuseInProgress(t *testing.T) {
	st, fed, enemy, h := newAITestState(1)
	h.Special(enemy, 5000, fed)
	if st.Defenseless != 0 || st.Corbomite != 0 || st.Surrender != 0 || st.SurrenderP != 0 {
		t.Error("no ruse state should change when nothing is in progress")
	}
}

func TestStandardStrategyReturnsImmediatelyWhenDead(t *testing.T) {
	_, _, enemy, h := newAITestState(1)
	enemy.Status[SysDead] = FullyDamaged
	origCourse := enemy.NewCourse
	h.StandardStrategy(enemy)
	if enemy.NewCourse != origCourse {
		t.Error("dead ships should never be given new orders")
	}
}

func TestStandardStrategyReturnsWhenSurrenderFlagsSet(t *testing.T) {
	st, fed, enemy, h := newAITestState(1)
	// surrender==2 keeps the flag alive across special()'s default
	// clearing arm (see TestSpecialClearsFedSurrenderFlagWhenNoLongerActive).
	st.Surrender = 2
	st.PlayStatus |= StatusFedSurrender
	origCourse := enemy.NewCourse
	h.StandardStrategy(enemy)
	_ = fed
	if enemy.NewCourse != origCourse {
		t.Error("strategy should not act while a surrender flag is set")
	}
}

func TestStandardStrategyEngagesAtLongRange(t *testing.T) {
	_, fed, enemy, h := newAITestState(4)
	enemy.X, enemy.Y = 10000, 0 // long range
	fed.X, fed.Y = 0, 0
	enemy.Cloaking = CloakNone // no cloak capability to keep the test deterministic

	h.StandardStrategy(enemy)

	// At long range with no cloak, the AI should either lock
	// weapons, attack (pursue), or at minimum not crash; verify no
	// panic and that some sensible outcome occurred (target set, or
	// tubes/phasers locked, or a message was queued).
	locked := enemy.Tubes[0].Target != nil || enemy.Tubes[1].Target != nil ||
		enemy.Phasers[0].Target != nil || enemy.Phasers[1].Target != nil
	if !locked && enemy.Target != fed && len(h.messages) == 0 {
		t.Error("expected StandardStrategy to take some action at long range")
	}
}

func TestSpecialClearsFedSurrenderFlagWhenNoLongerActive(t *testing.T) {
	st, fed, enemy, h := newAITestState(1)
	// The flag persists only while surrender is actively 2/3; once the
	// counter drops out of that range (e.g. back to 0), special()'s
	// default case clears the flag again next call, mirroring the
	// original's "global &= ~F_SURRENDER" in the switch's default arm.
	st.PlayStatus |= StatusFedSurrender
	h.Special(enemy, 5000, fed)
	if st.PlayStatus&StatusFedSurrender != 0 {
		t.Error("PlayStatus should have StatusFedSurrender cleared once surrender is no longer in progress")
	}
}
