package game

import (
	"fmt"
	"math"
)

// This file ports the sixteen e_* enemy-strategy action primitives
// from enemycom.c. Each one performs (or evaluates whether to perform)
// a single tactical action for an enemy ship sp, given the player's
// ship fed, and reports whether it did anything - exactly mirroring
// the original's int-returning "did this succeed" convention so
// standard_strategy (strategy.go) can chain them together the same
// way strat1.c does.
//
// They are implemented as methods on CombatHooks since they need the
// same State/Rand/message-accumulation plumbing as the rest of
// combat, and since ECloakOff (e_cloak_off) already lives there.

// EAttack turns sp towards fed at a closing speed, ported from
// e_attack() in enemycom.c.
func (h *CombatHooks) EAttack(sp, fed *Ship) bool {
	tmpf := math.Abs(fed.Warp)
	if math.Abs(sp.Warp) >= tmpf+2.0 || sp.IsDead(SysWarp) {
		return false
	}
	speed := tmpf + float64(h.Rand.Randm(2)) + 2.0
	if speed > sp.MaxSpeed {
		speed = sp.MaxSpeed
	}
	h.EPursue(sp, fed, speed)
	if canSee(sp) && fed.SystemWorks(h.Rand, SysSensor) {
		h.emit(fmt.Sprintf("%s:  %s attacking.", h.State.Crew.Helmsman, sp.Name))
	}
	return true
}

// ECheckArms returns the number of currently loaded, undamaged
// weapons sp has, ported from e_checkarms() in enemycom.c.
//
// The original had a bug here: it counted a weapon as armed via
// "load >= 0", which is always true (even a fully unloaded weapon has
// load == 0) - so e_checkarms() never actually detected an unarmed
// ship. This port uses "load > 0", fixing that bug.
func (h *CombatHooks) ECheckArms(sp *Ship) int {
	arms := 0
	for i := range sp.Phasers {
		if sp.Phasers[i].Status&PhaserDamaged == 0 && sp.Phasers[i].Load > 0 {
			arms++
		}
	}
	for i := range sp.Tubes {
		if sp.Tubes[i].Status&TubeDamaged == 0 && sp.Tubes[i].Load > 0 {
			arms++
		}
	}
	return arms
}

// ECheckProbe reports whether sp took evasive action to avoid a
// nearby enemy probe, ported from e_checkprobe() in enemycom.c.
func (h *CombatHooks) ECheckProbe(sp *Ship) bool {
	// Cloaked ships cannot be detected by probes, so there's nothing
	// to check.
	if cantSee(sp) {
		return false
	}
	for _, obj := range h.State.Objects {
		if obj.Type != ObjectProbe {
			continue
		}
		rng := rangeFind(sp.X, obj.X, sp.Y, obj.Y)
		if rng < 2000 {
			h.EEvade(sp, obj.X, obj.Y)
			return true
		}
	}
	return false
}

// ECloakOn engages sp's cloaking device if possible, ported from
// e_cloak_on() in enemycom.c.
func (h *CombatHooks) ECloakOn(sp, fed *Ship) bool {
	if sp.CloakDelay > 0 || sp.Cloaking != CloakOff {
		return false
	}
	sp.Cloaking = CloakOn
	sp.LastKnown.X = sp.X
	sp.LastKnown.Y = sp.Y
	sp.LastKnown.Warp = sp.Warp
	sp.LastKnown.Course = sp.Course
	if fed.SystemWorks(h.Rand, SysSensor) {
		h.emit(fmt.Sprintf("%s:  The %s has disappeared from our sensors.", h.State.Crew.Science, sp.Name))
	}
	return true
}

// ECloseTorps reports whether sp fired phasers at, or evaded, an
// incoming torpedo from fed, ported from e_closetorps() in
// enemycom.c.
//
// The original had a bug here: when no phaser shot could be gotten
// off, it called e_evade() using "tp" - the loop variable, which by
// the end of the loop pointed at the LAST torpedo examined in the
// list (regardless of whether it was actually close), not the
// dangerous one that had been flagged in "bad". This port uses "bad"
// itself, fixing that bug.
func (h *CombatHooks) ECloseTorps(sp, fed *Ship) bool {
	// Proximity fuses cannot affect a cloaked ship, so there's
	// nothing to check.
	if cantSee(sp) {
		return false
	}
	var bad *SpaceObject
	for _, obj := range h.State.Objects {
		if obj.Type != ObjectTorpedo || obj.From != fed {
			continue
		}
		rng := rangeFind(sp.X, obj.X, sp.Y, obj.Y)
		if rng < 1200 {
			bad = obj
			// Fire phasers - hope they're pointing in the right
			// direction!
			if h.EPhasers(sp, nil) > 0 {
				return true
			}
		}
	}
	// We can't get a phaser shot off; try to evade (although
	// hopeless).
	if bad != nil {
		h.EEvade(sp, bad.X, bad.Y)
		return true
	}
	return false
}

// EDestruct initiates sp's self-destruct sequence, ported from
// e_destruct() in enemycom.c.
func (h *CombatHooks) EDestruct(sp, fed *Ship) bool {
	if sp.Delay < 5.0 {
		return false
	}
	sp.Delay = 5.0
	h.ECloakOff(sp, fed)
	sp.Cloaking = CloakNone
	if fed.SystemWorks(h.Rand, SysSensor) {
		h.emit(fmt.Sprintf("%s: The %s is overloading what remains of it's", h.State.Crew.Science, sp.Name))
		h.emit("   antimatter pods -- obviously a suicidal gesture.")
		h.emit("   Estimate detonation in five seconds.")
	}
	return true
}

// EEvade turns sp away from the threat at (x,y), ported from
// e_evade() in enemycom.c. The original's "type" parameter was never
// actually used by the function body (marked "/* LINT */"), so it is
// omitted here.
func (h *CombatHooks) EEvade(sp *Ship, x, y int) bool {
	bear := bearingTo(sp.X, x, sp.Y, y)
	// The original checks the player's own sensors here (shiplist[0]),
	// not fed's - e_evade doesn't even take a fed parameter.
	if canSee(sp) && h.State.Player().SystemWorks(h.Rand, SysSensor) {
		h.emit(fmt.Sprintf("%s taking evasive action!", sp.Name))
	}
	var newCourse float64
	switch h.Rand.Randm(3) {
	case 1:
		newCourse = rectify(bear - 90.0)
	case 2:
		newCourse = rectify(bear + 90.0)
	case 3:
		newCourse = rectify(bear + 180.0)
	}
	sp.Target = nil
	sp.NewCourse = newCourse
	sp.NewWarp = 2 + float64(h.Rand.Randm(int(sp.MaxSpeed-3)))
	if sp.IsDead(SysWarp) {
		sp.NewWarp = 1.0
	}
	return true
}

// EJettison jettisons sp's engineering section, ported from
// e_jettison() in enemycom.c.
func (h *CombatHooks) EJettison(sp, fed *Ship) bool {
	if sp.IsDead(SysEngineering) {
		return false
	}
	h.ECloakOff(sp, fed)
	if h.State.Player().SystemWorks(h.Rand, SysSensor) {
		h.emit(fmt.Sprintf("%s: Sensors indicate debris being left by", h.State.Crew.Science))
		h.emit(fmt.Sprintf("   the %s.  Insufficient mass . . .", sp.Name))
	}
	obj := &SpaceObject{
		ID:        h.State.nextObjectID(),
		Type:      ObjectEngineering,
		From:      sp,
		X:         sp.X,
		Y:         sp.Y,
		Course:    sp.Course,
		Speed:     sp.Warp,
		NewSpeed:  0.0,
		Target:    nil,
		Proximity: 0,
		TimeDelay: 10.0,
		Fuel:      int(sp.Energy),
	}
	h.LaunchedObjects = append(h.LaunchedObjects, obj)

	// Ship slows to warp 1.0 when jettisoning engineering.
	sp.Energy = 0
	sp.Pods = 0
	sp.Regen = 0.0
	if sp.NewWarp < -1.0 {
		sp.NewWarp = -0.99
	}
	if sp.NewWarp > 1.0 {
		sp.NewWarp = 0.99
	}
	sp.MaxSpeed = 1.0
	sp.Status[SysEngineering] = FullyDamaged // Listed as destroyed.
	sp.Status[SysWarp] = FullyDamaged
	sp.Cloaking = CloakNone
	sp.TubeBlindLeft, sp.TubeBlindRight = 180, 180
	sp.PhaserBlindLeft, sp.PhaserBlindRight = 180, 180
	return true
}

// ELaunchProbe launches an antimatter probe at fed, ported from
// e_launchprobe() in enemycom.c.
func (h *CombatHooks) ELaunchProbe(sp, fed *Ship) bool {
	if !sp.SystemWorks(h.Rand, SysProbe) || sp.Energy <= 10 || cantSee(sp) {
		return false
	}
	// The fed ship has to be going slow before we'll launch a probe at
	// it.
	if math.Abs(fed.Warp) > 1.0 {
		return false
	}
	obj := &SpaceObject{
		ID:        h.State.nextObjectID(),
		Type:      ObjectProbe,
		From:      sp,
		X:         sp.X,
		Y:         sp.Y,
		Course:    bearingTo(sp.X, fed.X, sp.Y, fed.Y),
		Speed:     sp.Warp,
		NewSpeed:  3.0,
		Target:    fed,
		Proximity: 200 + h.Rand.Randm(200),
		TimeDelay: 15.0,
	}
	i := h.Rand.Randm(15) + 10
	if float64(i) > sp.Energy {
		i = int(sp.Energy)
	}
	obj.Fuel = i
	sp.Energy -= float64(i)
	sp.Pods -= float64(i)
	h.LaunchedObjects = append(h.LaunchedObjects, obj)
	h.emit(fmt.Sprintf("%s launching probe #%d", sp.Name, obj.ID))
	return true
}

// ELoadTubes loads energy into sp's undamaged torpedo tubes, returning
// the number of tubes loaded, ported from e_loadtubes() in enemycom.c.
func (h *CombatHooks) ELoadTubes(sp *Ship) int {
	const below = 10.0
	loaded := 0
	for i := range sp.Tubes {
		if sp.Energy <= below {
			break
		}
		t := &sp.Tubes[i]
		if t.Status&TubeDamaged != 0 {
			continue
		}
		j := math.Min(sp.Energy, MaxTubeCharge-t.Load)
		if j <= 0 {
			continue
		}
		sp.Energy -= j
		sp.Pods -= j
		t.Load += j
		loaded++
	}
	return loaded
}

// ELockPhasers locks every currently-unlocked, undamaged phaser bank
// onto fed, returning the number of banks locked, ported from
// e_lockphasers() in enemycom.c.
func (h *CombatHooks) ELockPhasers(sp, fed *Ship) int {
	banks := 0
	for i := range sp.Phasers {
		p := &sp.Phasers[i]
		if p.Status&PhaserDamaged != 0 || p.Target != nil {
			continue
		}
		p.Target = fed
		banks++
	}
	return banks
}

// ELockTubes locks every currently-unlocked, undamaged torpedo tube
// onto fed, returning the number of tubes locked, ported from
// e_locktubes() in enemycom.c.
func (h *CombatHooks) ELockTubes(sp, fed *Ship) int {
	tubes := 0
	for i := range sp.Tubes {
		t := &sp.Tubes[i]
		if t.Status&TubeDamaged != 0 || t.Target != nil {
			continue
		}
		t.Target = fed
		tubes++
	}
	return tubes
}

// EPhasers selects and marks a random subset of sp's loaded, undamaged
// phaser banks to fire, returning how many were selected. If fed is
// non-nil, a bank only counts if it is locked onto fed and would
// actually score a hit; if fed is nil (e.g. a snap shot at a nearby
// torpedo), every loaded bank is eligible regardless of lock/hit.
// Ported from e_phasers() in enemycom.c.
func (h *CombatHooks) EPhasers(sp, fed *Ship) int {
	banks := 0
	howMany := h.Rand.Randm(len(sp.Phasers)/2) + len(sp.Phasers)/2
	sp.PhaserSpread = 10 + h.Rand.Randm(12)
	for i := range sp.Phasers {
		p := &sp.Phasers[i]
		if p.Status&PhaserDamaged != 0 || p.Load == 0 {
			continue
		}
		if fed != nil {
			if p.Target == nil {
				continue
			}
			bear := bearingTo(sp.X, fed.X, sp.Y, fed.Y)
			if PhaserHit(sp, fed.X, fed.Y, p, bear) <= 0 {
				continue
			}
		}
		banks++
		p.Status |= PhaserFiring
		if banks >= howMany {
			break
		}
	}
	return banks
}

// EPursue turns sp towards fed at the given speed, slowing down for
// the turn if necessary, ported from e_pursue() in enemycom.c.
func (h *CombatHooks) EPursue(sp, fed *Ship, speed float64) bool {
	bear := bearingTo(sp.X, fed.X, sp.Y, fed.Y)
	// Do a quick turn if our speed is > max_warp - 2 and (thus) we are
	// never going to bear on the fed ship. speed = max_warp / 2 is a
	// magic cookie; feel free to change.
	coursediff := rectify(sp.Course - bear)
	if coursediff > 180.0 {
		coursediff -= 360.0
	}
	if speed >= sp.MaxSpeed-2 && math.Abs(coursediff) > 10 {
		speed = math.Trunc(sp.MaxSpeed / 2)
	}
	sp.Target = fed
	sp.NewCourse = bear
	sp.NewWarp = speed
	if speed > 1 && sp.IsDead(SysWarp) {
		sp.NewWarp = 0.99
	}
	return true
}

// ERunaway turns sp's strongest shield towards fed and accelerates to
// 2/3 maximum speed in reverse (or forward, from behind its strongest
// shield), ported from e_runaway() in enemycom.c.
//
// The original had a bug here: "sp->newwarp = 2 / 3 * sp->max_speed *
// sign;" used integer division, so 2/3 truncated to 0 and every
// retreat speed was silently zeroed - e_runaway() never actually made
// the ship move. This port uses 2.0/3.0, fixing that bug.
func (h *CombatHooks) ERunaway(sp, fed *Ship) bool {
	bear := bearingTo(sp.X, fed.X, sp.Y, fed.Y)

	// Find the strongest shield.
	strong := 0
	strength := 0.0
	for i := 0; i < NumShields; i++ {
		mult := 1.0
		if i == ShieldForward {
			mult = ShieldForwardBonus
		}
		temp := sp.Shields[i].Eff * sp.Shields[i].Drain * mult
		if temp > strength {
			strong = i
			strength = temp
		}
	}
	sign := 1.0
	var course float64
	switch strong {
	case 0:
		course = bear
		sign = -1
	case 1:
		course = rectify(bear - 90)
	case 2:
		course = rectify(bear + 180)
	case 3:
		course = rectify(bear + 90)
	}
	sp.Target = nil
	sp.NewCourse = course
	sp.NewWarp = 2.0 / 3.0 * sp.MaxSpeed * sign
	if sp.NewWarp > 1.0 && sp.IsDead(SysWarp) {
		sp.NewWarp = 0.99
	}
	if canSee(sp) && fed.SystemWorks(h.Rand, SysSensor) {
		h.emit(fmt.Sprintf("%s: The %s is retreating.", h.State.Crew.Helmsman, sp.Name))
	}
	return true
}

// ETorpedo selects and marks a random subset of sp's loaded, locked
// torpedo tubes to fire, returning how many were selected, ported
// from e_torpedo() in enemycom.c.
//
// Note the original (and this port) only checks sp's range to other
// ENEMY ships, not the player's - it never refuses to fire because
// the player's ship itself is in the way, only because another enemy
// ship is.
func (h *CombatHooks) ETorpedo(sp *Ship) int {
	// Don't shoot if someone might be in the way (i.e. the proximity
	// fuse will go off right as the torps leave the tubes!).
	for _, other := range h.State.Enemies() {
		if other == sp {
			continue
		}
		if rangeFind(sp.X, other.X, sp.Y, other.Y) <= 400 {
			return 0
		}
	}
	tubes := 0
	// This is not, and should not be, dependent on the number of
	// tubes one has.
	howMany := h.Rand.Randm(2) + 1
	for i := range sp.Tubes {
		t := &sp.Tubes[i]
		if t.Status&TubeDamaged != 0 || t.Load == 0 || t.Target == nil {
			continue
		}
		tubes++
		t.Status |= TubeFiring
		if tubes >= howMany {
			break
		}
	}
	return tubes
}
