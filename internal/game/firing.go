package game

import (
	"fmt"

	"github.com/appliedinspiration/inspired_trek73/internal/data"
)

// AntimatterHit applies an antimatter blast at (x,y) with the given
// fuel yield to every ship and object in the battle except the source
// itself, ported from antimatter_hit() in subs.c. Exactly one of
// sourceShip/sourceObj should be non-nil, identifying what not to hit
// (a torpedo/probe/ship-self-destruct never damages itself).
func AntimatterHit(st *State, sourceShip *Ship, sourceObj *SpaceObject, x, y, fuel int, r *Rand) []string {
	var messages []string
	fed := st.Player()

	for _, sp := range st.Ships {
		if sp == sourceShip {
			continue
		}
		hit := TorpedoHit(fuel, x, y, sp.X, sp.Y)
		if hit <= 0 {
			continue
		}
		bear := rectify(bearingTo(sp.X, x, sp.Y, y) - sp.Course)
		facing := shieldFacingFromAntimatterBearing(bear)
		messages = append(messages, Damage(hit, sp, facing, &data.AntimatterDamage, DamageAntimatter, fed, r)...)
	}

	for _, obj := range st.Objects {
		if obj == sourceObj {
			continue
		}
		hit := TorpedoHit(fuel, x, y, obj.X, obj.Y)
		if hit <= 0 {
			continue
		}
		if obj.TimeDelay <= SegmentSeconds {
			continue
		}
		obj.TimeDelay = SegmentSeconds
		messages = append(messages, objectHitMessage(obj))
	}
	return messages
}

func objectHitMessage(obj *SpaceObject) string {
	switch obj.Type {
	case ObjectTorpedo:
		return fmt.Sprintf("hit on torpedo %d", obj.ID)
	case ObjectProbe:
		return fmt.Sprintf("hit on probe %d", obj.ID)
	case ObjectEngineering:
		name := "unknown"
		if obj.From != nil {
			name = obj.From.Name
		}
		return fmt.Sprintf("hit on %s engineering", name)
	default:
		return fmt.Sprintf("hit on unknown item %d", obj.ID)
	}
}

// CombatHooks implements MovementHooks, tying phaser/torpedo firing,
// detonation, and enemy cloak shutdown into a single stateful object
// that MoveShips can call into. It replaces the original's reliance on
// global shiplist[]/head/tail state (firing.c, dist.c/enemycom.c).
type CombatHooks struct {
	State *State
	Rand  *Rand

	// DetonatedObjects collects objects that detonated this call, so
	// the caller can remove them from State.Objects (the original
	// deleted them from the linked list in-place via delitem()).
	DetonatedObjects []*SpaceObject
	// LaunchedObjects collects objects newly created by torpedo fire,
	// so the caller can append them to State.Objects.
	LaunchedObjects []*SpaceObject

	messages []string
}

// NewCombatHooks returns a CombatHooks ready to be passed to
// MoveShips.
func NewCombatHooks(st *State, r *Rand) *CombatHooks {
	return &CombatHooks{State: st, Rand: r}
}

// TakeMessages returns and clears every message accumulated since the
// last call.
func (h *CombatHooks) TakeMessages() []string {
	m := h.messages
	h.messages = nil
	return m
}

func (h *CombatHooks) emit(s string) {
	h.messages = append(h.messages, s)
}

// PhaserFiring fires the first phaser bank on sp with its firing flag
// set, ported from phaser_firing() in firing.c.
//
// The original had a bug here: "if ((sp = fed) && (!shutup[PHASERS+i])
// ...)" used assignment instead of comparison, which silently
// reassigned the local sp variable to always be the player's ship for
// the rest of the function - including the disengage message and the
// sp->phasers[i].target = NULL clear, which should have applied to
// whichever ship was actually firing. This port uses "sp == fed" and
// never reassigns sp, fixing that bug.
func (h *CombatHooks) PhaserFiring(sp *Ship) {
	fed := h.State.Player()

	i := -1
	for idx := range sp.Phasers {
		if sp.Phasers[idx].Status&PhaserFiring != 0 {
			i = idx
			break
		}
	}
	if i < 0 {
		return
	}
	bank := &sp.Phasers[i]
	bank.Status &^= PhaserFiring
	target := bank.Target

	// j is the relative bearing of the phasers relative to the ship;
	// bear is the absolute direction the phasers are pointing.
	var bear, j float64
	if target == nil {
		bear = bank.Bearing + sp.Course
		j = rectify(bank.Bearing)
	} else {
		tx, ty := target.X, target.Y
		if cantSee(target) {
			tx, ty = target.LastKnown.X, target.LastKnown.Y
		}
		bear = bearingTo(sp.X, tx, sp.Y, ty)
		j = rectify(bear - sp.Course)
	}
	if betw(j, float64(sp.PhaserBlindLeft), float64(sp.PhaserBlindRight)) && !sp.IsDead(SysEngineering) {
		return
	}
	if target != nil && target.IsDead(SysDead) {
		if sp == fed && !h.State.Shutup.Phaser[i] && !sp.IsDead(SysDead) {
			h.emit(fmt.Sprintf("%s phaser %d unlocking", sp.Name, i+1))
			h.State.Shutup.Phaser[i] = true
		}
		bank.Target = nil
		return
	}
	if cantSee(sp) {
		h.ECloakOff(sp, fed)
	}
	h.emit(fmt.Sprintf(" <%s frng phasers>", sp.Name))

	for _, ep := range h.State.Ships {
		if ep == sp {
			continue
		}
		hit := PhaserHit(sp, ep.X, ep.Y, bank, bear)
		if hit <= 0 {
			continue
		}
		bearFromTarget := rectify(bearingTo(ep.X, sp.X, ep.Y, sp.Y) - ep.Course)
		facing := shieldFacingFromPhaserBearing(bearFromTarget)
		h.messages = append(h.messages, Damage(hit, ep, facing, &data.PhaserDamage, DamagePhaser, fed, h.Rand)...)
	}
	for _, obj := range h.State.Objects {
		hit := PhaserHit(sp, obj.X, obj.Y, bank, bear)
		if hit <= 0 {
			continue
		}
		if obj.TimeDelay > SegmentSeconds {
			h.emit(objectHitMessage(obj))
			obj.TimeDelay = 0
		}
		obj.Fuel -= hit / 2
		if obj.Fuel < 0 {
			obj.Fuel = 0
		}
	}

	// Reduce the load by the firing percentage.
	bank.Load *= 1.0 - float64(sp.PhaserFirePct)/100
}

// TorpedoFiring fires the first tube on sp with its firing flag set,
// launching a new torpedo object, ported from torpedo_firing() in
// firing.c.
//
// The original had a bug here: when no target was locked, it computed
// the tube's firing bearing from sp->phasers[i].bearing instead of
// sp->tubes[i].bearing - an unguided torpedo fired from an
// unlocked tube would aim wherever phaser bank i happened to be
// pointed, rather than where tube i was aimed. This port uses the
// tube's own bearing, fixing that bug.
func (h *CombatHooks) TorpedoFiring(sp *Ship) {
	fed := h.State.Player()

	i := -1
	for idx := range sp.Tubes {
		if sp.Tubes[idx].Status&TubeFiring != 0 {
			i = idx
			break
		}
	}
	if i < 0 {
		return
	}
	tube := &sp.Tubes[i]
	tube.Status &^= TubeFiring
	th := int(tube.Load)
	if th == 0 {
		return
	}
	target := tube.Target

	var bear, j float64
	if target == nil {
		bear = tube.Bearing + sp.Course
		j = rectify(tube.Bearing)
	} else {
		tx, ty := target.X, target.Y
		if cantSee(target) {
			tx, ty = target.LastKnown.X, target.LastKnown.Y
		}
		bear = bearingTo(sp.X, tx, sp.Y, ty)
		j = rectify(bear - sp.Course)
	}
	if betw(j, float64(sp.TubeBlindLeft), float64(sp.TubeBlindRight)) && !sp.IsDead(SysEngineering) {
		return
	}
	if target != nil && target.IsDead(SysDead) {
		if sp == fed && !h.State.Shutup.Tube[i] && !sp.IsDead(SysDead) {
			h.emit(fmt.Sprintf("   tube %d disengaging", i+1))
			h.State.Shutup.Tube[i] = true
		}
		tube.Target = nil
		return
	}

	tube.Load = 0
	obj := &SpaceObject{
		ID:        h.State.nextObjectID(),
		Type:      ObjectTorpedo,
		From:      sp,
		X:         sp.X,
		Y:         sp.Y,
		Target:    nil,
		Course:    rectify(bear),
		Fuel:      th,
		Speed:     float64(sp.TubeLaunchSpd) + sp.Warp,
		TimeDelay: float64(sp.TubeDelay),
		Proximity: sp.TubeProximity,
	}
	obj.NewSpeed = obj.Speed
	h.LaunchedObjects = append(h.LaunchedObjects, obj)

	if cantSee(sp) {
		h.ECloakOff(sp, fed)
	}
	h.emit(fmt.Sprintf(" <<%s frng torpedo %d>>", sp.Name, obj.ID))
}

// ShipDetonate destroys sp (self-destruct or fatal damage), releasing
// its remaining phaser/tube charge and antimatter pods as a blast,
// ported from ship_detonate() in moveships.c/firing.c.
func (h *CombatHooks) ShipDetonate(sp *Ship) {
	h.emit(fmt.Sprintf("++%s++ destruct.", sp.Name))

	fuel := 0
	for i := range sp.Phasers {
		if sp.Phasers[i].Status&PhaserDamaged == 0 {
			fuel += int(minFloat(sp.Phasers[i].Load, MaxPhaserCharge))
		}
	}
	for i := range sp.Tubes {
		if sp.Tubes[i].Status&TubeDamaged == 0 {
			fuel += int(minFloat(sp.Tubes[i].Load, MaxTubeCharge))
		}
	}
	fuel += int(sp.Pods)

	h.messages = append(h.messages, AntimatterHit(h.State, sp, nil, sp.X, sp.Y, fuel, h.Rand)...)

	for i := 0; i < NumDamageSystems; i++ { // S_NUMSYSTEMS: computer, sensor, probe, warp
		sp.Status[i] = 100
	}
	sp.Cloaking = CloakNone
	sp.Complement = -1
}

// TorpDetonate detonates a torpedo, probe, or jettisoned engineering
// section, ported from torp_detonate() in firing.c. The caller is
// responsible for removing obj from State.Objects afterward (recorded
// in h.DetonatedObjects).
func (h *CombatHooks) TorpDetonate(obj *SpaceObject) {
	switch obj.Type {
	case ObjectTorpedo:
		h.emit(fmt.Sprintf(":: torp %d ::", obj.ID))
	case ObjectProbe:
		h.emit(fmt.Sprintf("** probe %d **", obj.ID))
	case ObjectEngineering:
		name := "unknown"
		if obj.From != nil {
			name = obj.From.Name
		}
		h.emit(fmt.Sprintf("## %s engineering ##", name))
	}
	h.messages = append(h.messages, AntimatterHit(h.State, nil, obj, obj.X, obj.Y, obj.Fuel, h.Rand)...)
	h.DetonatedObjects = append(h.DetonatedObjects, obj)
}

// ECloakOff forcibly decloaks an enemy ship whose cloaking device has
// run out of power, ported from e_cloak_off() in enemycom.c. Returns
// true if the cloaking device was turned off.
func (h *CombatHooks) ECloakOff(sp, fed *Ship) bool {
	if sp.Cloaking != CloakOn {
		return false
	}
	sp.Cloaking = CloakOff
	sp.CloakDelay = 4
	if sp.SystemWorks(h.Rand, SysSensor) {
		rng := rangeFind(sp.X, sp.LastKnown.X, sp.Y, sp.LastKnown.Y)
		h.emit(fmt.Sprintf("%s: The %s has reappeared on our sensors %d", h.State.Crew.Science, sp.Name, rng))
		h.emit("   megameters from its projected position.")
	}
	return true
}

// betw reports whether i is strictly between j and k, ported from the
// betw(i,j,k) macro in defines.h.
func betw(i, j, k float64) bool {
	return j < i && i < k
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
