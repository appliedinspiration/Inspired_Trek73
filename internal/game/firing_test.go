package game

import (
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/data"
)

func TestPhaserHitZeroBeyondRange(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.PhaserSpread = 45
	bank := &Phaser{Load: 10}
	hit := PhaserHit(sp, MaxPhaserRange+1, 0, bank, 0)
	if hit != 0 {
		t.Errorf("hit = %d, want 0 beyond max range", hit)
	}
}

func TestPhaserHitZeroOutsideSpread(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.PhaserSpread = 10
	bank := &Phaser{Load: 10}
	// Target due north (bearing 90), phaser pointed due east (trueBear 0):
	// spread is 90 degrees, exceeding the 10-degree cone.
	hit := PhaserHit(sp, 0, 500, bank, 0)
	if hit != 0 {
		t.Errorf("hit = %d, want 0 outside spread cone", hit)
	}
}

func TestPhaserHitPositiveWithinRangeAndSpread(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.PhaserSpread = 45
	sp.PhaserFirePct = 100
	bank := &Phaser{Load: 10}
	hit := PhaserHit(sp, 500, 0, bank, 0) // due east, close range, aimed dead-on
	if hit <= 0 {
		t.Errorf("hit = %d, want positive damage", hit)
	}
}

func TestTorpedoHitZeroBeyondBlastRadius(t *testing.T) {
	hit := TorpedoHit(1, 0, 0, 100000, 100000)
	if hit != 0 {
		t.Errorf("hit = %d, want 0 far outside blast radius", hit)
	}
}

func TestTorpedoHitPositiveAtCloseRange(t *testing.T) {
	hit := TorpedoHit(10, 0, 0, 10, 10)
	if hit <= 0 {
		t.Errorf("hit = %d, want positive damage at close range", hit)
	}
}

func TestDamageFullShieldBlocksAllInternalDamage(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	st := newTestState(sp)
	sp.Shields[0] = Shield{Eff: 1.0, Drain: 1.0}
	origEnergy := sp.Energy
	r := NewRand(1)

	messages := Damage(st, 100, sp, 1, &data.PhaserDamage, DamagePhaser, sp, r)

	want := "hit 100 on TestShip's shield 1"
	if len(messages) != 1 || messages[0] != want {
		t.Errorf("messages = %v, want [%q]", messages, want)
	}
	if sp.Energy != origEnergy {
		t.Errorf("Energy = %v, want unchanged %v", sp.Energy, origEnergy)
	}
}

func TestDamageReducesShieldEfficiencyAndEnergy(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	st := newTestState(sp)
	sp.Shields[0] = Shield{Eff: 0.0, Drain: 0.0} // shield fully down
	sp.Energy = 1000
	sp.Pods = 2000
	r := NewRand(1)

	Damage(st, 50, sp, 1, &data.PhaserDamage, DamagePhaser, sp, r)

	if sp.Energy >= 1000 {
		t.Errorf("Energy = %v, want reduced from 1000", sp.Energy)
	}
}

func TestDamageMarksWeaponDamagedEventually(t *testing.T) {
	// A large hit with the default weapon divisor should damage at
	// least one weapon across repeated applications.
	sp := newTestShip(0, 0, 0)
	st := newTestState(sp)
	sp.Shields[0] = Shield{Eff: 0.0, Drain: 0.0}
	sp.Energy = 100000
	sp.Pods = 200000
	r := NewRand(42)

	anyDamaged := false
	for i := 0; i < 20; i++ {
		Damage(st, 500, sp, 1, &data.PhaserDamage, DamagePhaser, sp, r)
		for _, p := range sp.Phasers {
			if p.Status&PhaserDamaged != 0 {
				anyDamaged = true
			}
		}
		for _, tb := range sp.Tubes {
			if tb.Status&TubeDamaged != 0 {
				anyDamaged = true
			}
		}
	}
	if !anyDamaged {
		t.Error("expected at least one weapon damaged after repeated heavy hits")
	}
}

func TestDamageDoesNotDetonateAtFullWarpStrength(t *testing.T) {
	fed := newTestShip(0, 0, 0)
	enemy := newTestShip(1, 0, 0)
	enemy.Status[SysWarp] = 0
	enemy.Shields[0] = Shield{Eff: 0.0, Drain: 0.0}
	st := newTestState(fed, enemy)
	r := NewRand(1)

	messages := Damage(st, 400, enemy, 1, &data.PhaserDamage, DamagePhaser, fed, r)
	for _, m := range messages {
		if m == "++"+enemy.Name+"++ destruct." {
			t.Fatalf("unexpected detonation at full warp strength: %v", messages)
		}
	}
	if enemy.Complement == -1 {
		t.Fatal("enemy should not have detonated")
	}
}

func TestDamageDetonationChanceUsesWarpDamageAndHitStrength(t *testing.T) {
	r := NewRand(1)
	detonated := false

	for i := 0; i < 20; i++ {
		fed := newTestShip(0, 0, 0)
		fed.Shields[0] = Shield{Eff: 0.0, Drain: 0.0}
		enemy := newTestShip(1, 0, 0)
		enemy.Status[SysWarp] = 60
		enemy.Shields[0] = Shield{Eff: 0.0, Drain: 0.0}
		enemy.Pods = 1000
		st := newTestState(fed, enemy)
		origEnergy := fed.Energy

		messages := Damage(st, 120, enemy, 1, &data.PhaserDamage, DamagePhaser, fed, r)
		for _, m := range messages {
			if m == "++"+enemy.Name+"++ destruct." {
				detonated = true
				if enemy.Complement != -1 {
					t.Fatalf("Complement = %d, want -1", enemy.Complement)
				}
				if fed.Energy >= origEnergy {
					t.Fatalf("fed.Energy = %v, want reduced from %v by blast damage", fed.Energy, origEnergy)
				}
				break
			}
		}
		if detonated {
			break
		}
	}
	if !detonated {
		t.Fatal("expected at least one damage-triggered detonation with warp damage and hit strength")
	}
}

func TestCheckLocksClearsHelmLockAtFullPercent(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	target := newTestShip(1, 100, 100)
	sp.Target = target
	sp.RelativeBear = 20
	r := NewRand(1)

	messages := CheckLocks(sp, 100, sp, r)

	if sp.Target != nil {
		t.Errorf("Target = %v, want nil after 100%% lock-loss check", sp.Target)
	}
	if sp.RelativeBear != 0 {
		t.Errorf("RelativeBear = %v, want 0", sp.RelativeBear)
	}
	found := false
	for _, m := range messages {
		if m == "Computer: "+sp.Name+" has lost helm lock" {
			found = true
		}
	}
	if !found {
		t.Errorf("messages = %v, want helm lock lost message", messages)
	}
}

func TestCheckLocksNeverClearsAtZeroPercent(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	target := newTestShip(1, 100, 100)
	sp.Target = target
	r := NewRand(1)

	CheckLocks(sp, 0, sp, r)

	if sp.Target != target {
		t.Errorf("Target = %v, want unchanged at 0%% chance", sp.Target)
	}
}

func TestAntimatterHitSkipsSource(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	st := newTestState(sp)
	r := NewRand(1)

	// Source is sp itself, at its own position: since sp is excluded,
	// no messages/damage should occur even at zero range.
	messages := AntimatterHit(st, sp, nil, sp.X, sp.Y, 100, r)
	if len(messages) != 0 {
		t.Errorf("messages = %v, want none (source excluded)", messages)
	}
}

func TestAntimatterHitDamagesNearbyShip(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	victim := newTestShip(1, 10, 10)
	victim.Shields[0] = Shield{Eff: 0, Drain: 0}
	st := newTestState(sp, victim)
	r := NewRand(1)

	AntimatterHit(st, sp, nil, 0, 0, 100, r)

	if victim.Energy >= 1000 {
		t.Errorf("victim.Energy = %v, want reduced from 1000", victim.Energy)
	}
}

func TestPhaserFiringDoesNotSwapFiringShip(t *testing.T) {
	// Regression test for the original's "if ((sp = fed) && ...)"
	// assignment-as-comparison bug: an enemy ship (not fed) firing at
	// an already-dead target must clear *its own* phaser lock, not
	// the player's.
	fed := newTestShip(0, 0, 0)
	enemy := newTestShip(1, 5000, 5000)
	deadTarget := newTestShip(2, 6000, 6000)
	deadTarget.Status[SysDead] = FullyDamaged

	enemy.Phasers[0].Status |= PhaserFiring
	enemy.Phasers[0].Target = deadTarget
	fed.Phasers[0].Target = newTestShip(3, 1, 1) // fed's own unrelated lock

	st := newTestState(fed, enemy, deadTarget)
	hooks := NewCombatHooks(st, NewRand(1))

	hooks.PhaserFiring(enemy)

	if enemy.Phasers[0].Target != nil {
		t.Errorf("enemy.Phasers[0].Target = %v, want nil (enemy's own lock cleared)", enemy.Phasers[0].Target)
	}
	if fed.Phasers[0].Target == nil {
		t.Error("fed.Phasers[0].Target was cleared, want unchanged (bug: firing ship swapped with fed)")
	}
}

func TestTorpedoFiringUsesTubeBearingNotPhaserBearing(t *testing.T) {
	// Regression test for the original's bug of reading
	// sp->phasers[i].bearing instead of sp->tubes[i].bearing when no
	// target is locked.
	sp := newTestShip(0, 0, 0)
	sp.Course = 0
	sp.Phasers[0].Bearing = 10 // deliberately different from the tube
	sp.Tubes[0].Bearing = 90
	sp.Tubes[0].Status |= TubeFiring
	sp.Tubes[0].Load = 5
	sp.Tubes[0].Target = nil

	st := newTestState(sp)
	hooks := NewCombatHooks(st, NewRand(1))

	hooks.TorpedoFiring(sp)

	if len(hooks.LaunchedObjects) != 1 {
		t.Fatalf("LaunchedObjects = %d, want 1", len(hooks.LaunchedObjects))
	}
	obj := hooks.LaunchedObjects[0]
	wantCourse := rectify(sp.Tubes[0].Bearing + sp.Course)
	if obj.Course != wantCourse {
		t.Errorf("torpedo course = %v, want %v (tube bearing, not phaser bearing %v)",
			obj.Course, wantCourse, sp.Phasers[0].Bearing)
	}
}

func TestShipDetonateMarksShipDestroyed(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Name = "Trakka"
	other := newTestShip(1, 5, 5)
	st := newTestState(sp, other)
	hooks := NewCombatHooks(st, NewRand(1))

	hooks.ShipDetonate(sp)

	if sp.Complement != -1 {
		t.Errorf("Complement = %d, want -1", sp.Complement)
	}
	if sp.Cloaking != CloakNone {
		t.Errorf("Cloaking = %d, want CloakNone", sp.Cloaking)
	}
	for i := 0; i < NumDamageSystems; i++ {
		if sp.Status[i] != 100 {
			t.Errorf("Status[%d] = %d, want 100", i, sp.Status[i])
		}
	}
	found := false
	for _, m := range hooks.TakeMessages() {
		if m == "++Trakka++ destruct." {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected exact ship destruct message")
	}
}

func TestShipDetonateOnlyOnce(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Name = "Trakka"
	other := newTestShip(1, 5, 5)
	st := newTestState(sp, other)
	hooks := NewCombatHooks(st, NewRand(1))

	hooks.ShipDetonate(sp)
	hooks.ShipDetonate(sp)

	destructMsgs := 0
	for _, m := range hooks.TakeMessages() {
		if m == "++Trakka++ destruct." {
			destructMsgs++
		}
	}
	if destructMsgs != 1 {
		t.Fatalf("destruct message count = %d, want 1", destructMsgs)
	}
}
func TestTorpDetonateRecordsDetonatedObject(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	obj := &SpaceObject{ID: 5, Type: ObjectTorpedo, From: sp, X: 0, Y: 0, Fuel: 10}
	st := newTestState(sp)
	hooks := NewCombatHooks(st, NewRand(1))

	hooks.TorpDetonate(obj)

	if len(hooks.DetonatedObjects) != 1 || hooks.DetonatedObjects[0] != obj {
		t.Errorf("DetonatedObjects = %v, want [obj]", hooks.DetonatedObjects)
	}
	found := false
	for _, m := range hooks.TakeMessages() {
		if m == ":: torp 5 ::" {
			found = true
		}
	}
	if !found {
		t.Error("expected torpedo detonation message")
	}
}

func TestECloakOffTurnsCloakOff(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Cloaking = CloakOn
	fed := newTestShip(1, 0, 0)
	st := newTestState(fed, sp)
	hooks := NewCombatHooks(st, NewRand(1))

	hooks.ECloakOff(sp, fed)

	if sp.Cloaking != CloakOff {
		t.Errorf("Cloaking = %d, want CloakOff", sp.Cloaking)
	}
	if sp.CloakDelay != 4 {
		t.Errorf("CloakDelay = %d, want 4", sp.CloakDelay)
	}
}

func TestECloakOffNoOpWhenNotCloaked(t *testing.T) {
	sp := newTestShip(0, 0, 0)
	sp.Cloaking = CloakOff
	fed := newTestShip(1, 0, 0)
	st := newTestState(fed, sp)
	hooks := NewCombatHooks(st, NewRand(1))

	hooks.ECloakOff(sp, fed)

	if sp.CloakDelay != 0 {
		t.Errorf("CloakDelay = %d, want unchanged 0", sp.CloakDelay)
	}
}
