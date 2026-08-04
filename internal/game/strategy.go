package game

import (
	"fmt"
	"math"
)

// StandardStrategy runs one enemy ship's AI decision tree for the
// current turn, ported from standard_strategy() in strat1.c. It
// mutates sp (course/speed/weapon locks/etc.) and returns any
// player-visible messages produced.
func (h *CombatHooks) StandardStrategy(sp *Ship) {
	st := h.State
	fed := st.Player()

	if sp.IsDead(SysDead) {
		return
	}
	rng := rangeFind(sp.X, fed.X, sp.Y, fed.Y)
	bear := bearingTo(sp.X, fed.X, sp.Y, fed.Y)
	bear = rectify(bear - sp.Course)

	// Handle special requests (ruses/bluffs/surrenders).
	h.Special(sp, rng, fed)

	// Now check for surrendering flags.
	if st.PlayStatus&(StatusFedSurrender|StatusEnemySurrender) != 0 {
		return
	}

	// Always turn on the cloaking device if we have it and can afford
	// it.
	if sp.Cloaking == CloakOff && sp.Energy >= 20 && h.ECloakOn(sp, fed) {
		return
	}

	// Check for hostile antimatter devices.
	if sp.Cloaking == CloakOff && h.ECloseTorps(sp, fed) {
		return
	}
	if sp.Cloaking == CloakOff && h.ECheckProbe(sp) {
		return
	}

	// If cloaking is on, and we're running low on energy, drop the
	// cloak.
	if sp.Cloaking == CloakOn && sp.Energy < 30 && h.ECloakOff(sp, fed) {
		return
	}

	// Short range?
	if rng < 1050 {
		if h.ECheckArms(sp) < h.Rand.Randm((len(sp.Phasers)+len(sp.Tubes))/3) {
			if !h.ECloakOn(sp, fed) {
				h.ERunaway(sp, fed)
			}
			return
		}
		if h.ELockPhasers(sp, fed) > 0 {
			return
		}
		if h.EPhasers(sp, fed) > 0 {
			return
		}
		if betw(bear, 90.0, 270.0) && !betw(float64(st.Corbomite), 1, 6) {
			h.EPursue(sp, fed, 1.0)
			return
		}
		if h.ELaunchProbe(sp, fed) {
			return
		}
		if sp.Pods < 20 && sp.Regen < 4.0 && h.EDestruct(sp, fed) {
			return
		}
		// Set course?
		tmpf := math.Abs(fed.Warp)
		if (sp.Target != fed || math.Abs(sp.Warp)+tmpf > 2.0) && !betw(float64(st.Corbomite), 1, 6) {
			h.EPursue(sp, fed, tmpf)
			return
		}
		if h.ECloakOn(sp, fed) {
			return
		}
	}
	if rng < 3800 {
		// Either medium range, or we can't figure out what to do at
		// short range.
		if h.ELockTubes(sp, fed) > 0 {
			return
		}
		if sp.Energy > 30 && sp.Pods > 40 && h.ELoadTubes(sp) > 0 {
			return
		}
		if h.ETorpedo(sp) > 0 {
			return
		}
		// Should we run away; can we?
		if h.ECheckArms(sp) < h.Rand.Randm((len(sp.Phasers)+len(sp.Tubes))/3) {
			if !h.ECloakOn(sp, fed) {
				h.ERunaway(sp, fed)
			}
			return
		}
		// Pursued from behind, low power: jettison engineering!
		if betw(bear, 90.0, 270.0) && sp.Energy < 10 && sp.Regen < 4.0 && h.EJettison(sp, fed) {
			return
		}
		// Put in other junk later.
		if h.ECloakOn(sp, fed) {
			return
		}
	}

	// Either distant range, or we can't figure out what to do at
	// medium range.

	// Warp drive dead and Federation destructing: run away!
	if fed.Delay < 15.0 && sp.IsDead(SysWarp) {
		h.ERunaway(sp, fed)
		return
	}

	// Enemy in our blind area? Make a quick turn.
	if betw(bear, float64(sp.TubeBlindLeft), float64(sp.TubeBlindRight)) && !betw(float64(st.Corbomite), 1, 6) {
		h.EPursue(sp, fed, 1.0)
		return
	}
	if h.ELockTubes(sp, fed) > 0 {
		return
	}
	if h.ELockPhasers(sp, fed) > 0 {
		return
	}

	// Attack?
	if h.EAttack(sp, fed) {
		return
	}
	if sp.Energy > 30 && sp.Pods > 40 && h.ELoadTubes(sp) > 0 {
		return
	}
	if h.ECloakOn(sp, fed) {
		return
	}

	// Gee, there's nothing that we want to do!
	if canSee(sp) {
		h.emit(fmt.Sprintf("%s:  We're being scanned by the %s", st.Crew.Science, sp.Name))
	}
}
