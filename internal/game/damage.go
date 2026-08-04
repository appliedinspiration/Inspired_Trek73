package game

import (
	"fmt"
	"math"

	"github.com/appliedinspiration/inspired_trek73/internal/data"
)

// Damage applies the effects of a single hit to a ship's shield facing
// and internal systems, ported from damage() in damage.c.
//
// facing uses the original's 1-based shield numbering (1-4); it is
// converted to a 0-based index internally.
func Damage(hit int, ep *Ship, facing int, dam *data.DamageProfile, flag int, fed *Ship, r *Rand) []string {
	var messages []string
	s := facing - 1

	// If the shield is at 100% efficiency, no damage at all is taken
	// (except to the shield itself).
	f1 := float64(hit) * (1.0 - ep.Shields[s].Eff*ep.Shields[s].Drain)
	if f1 < 0 {
		return nil
	}

	// Calculate shield damage.
	f2 := 1.0
	switch flag {
	case DamageAntimatter:
		f2 = ep.TorpedoShieldDivisor * 100
	case DamagePhaser:
		f2 = ep.PhaserShieldDivisor * 100
	}
	if s == ShieldForward {
		f2 *= ShieldForwardBonus
	}
	ep.Shields[s].Eff -= math.Max(float64(hit)/f2, 0)
	if ep.Shields[s].Eff < 0.0 {
		ep.Shields[s].Eff = 0.0
	}

	// Calculate loss of fuel, regeneration, etc.
	ep.Eff += f1 / dam.EffDivisor
	ep.Pods -= f1 / dam.FuelDivisor
	ep.Energy -= f1 / dam.FuelDivisor
	ep.Regen -= f1 / dam.RegenDivisor
	if ep.Regen < 0.0 {
		ep.Regen = 0.0
	}
	if ep.Pods < 0.0 {
		ep.Pods = 0.0
	}
	if ep.Energy < 0.0 {
		ep.Energy = 0.0
	}
	if ep.Pods < ep.Energy {
		ep.Energy = ep.Pods
	}

	// Kill some crew.
	if ep.Complement > 0 {
		j := int(f1 * dam.CrewDivisor)
		if j > 0 {
			ep.Complement -= r.Randm(j)
		}
		if ep.Complement < 0 {
			ep.Complement = 0
		}
	}

	// Damage some weapons.
	numWeapons := len(ep.Phasers) + len(ep.Tubes)
	j := int(f1 / dam.WeaponDivisor)
	for i := 0; i < j; i++ {
		k := r.Randm(numWeapons) - 1
		if k < len(ep.Phasers) {
			p := &ep.Phasers[k]
			if p.Status&PhaserDamaged != 0 {
				continue
			}
			p.Status |= PhaserDamaged
			p.Target = nil
			// Reroute the energy back to the engines.
			ep.Energy = math.Min(ep.Pods, ep.Energy+p.Load)
			p.Load = 0
			p.Drain = 0
			if ep == fed {
				messages = append(messages, fmt.Sprintf("   phaser %d damaged", k+1))
			}
		} else {
			k -= len(ep.Phasers)
			t := &ep.Tubes[k]
			if t.Status&TubeDamaged != 0 {
				continue
			}
			// If tubes are damaged, reroute the pods back to the
			// engines.
			ep.Pods += t.Load
			ep.Energy += t.Load
			t.Load = 0
			t.Status |= TubeDamaged
			t.Target = nil
			if ep == fed {
				messages = append(messages, fmt.Sprintf("   tube %d damaged", k+1))
			}
		}
	}

	// Damage the different systems.
	for i := 0; i < NumDamageSystems; i++ {
		if ep.IsDead(i) { // Don't damage a dead system.
			continue
		}
		if float64(r.Randm(dam.Systems[i].Roll)) >= f1 {
			continue
		}
		percent := r.Randm(int(f1))
		ep.Status[i] += percent
		if ep.Status[i] > 100 {
			ep.Status[i] = 100
		}
		if ep == fed {
			if ep.IsDead(i) {
				messages = append(messages, "   "+dam.Systems[i].Message)
			} else {
				messages = append(messages, "   "+data.SystemNames[i]+" damaged.")
			}
		}

		if ep.IsDead(i) {
			// Effects of a totally destroyed system.
			switch i {
			case SysSensor, SysProbe:
				// No bookkeeping needed.
			case SysWarp:
				ep.MaxSpeed = 1.0
			case SysComputer:
				messages = append(messages, CheckLocks(ep, 100, fed, r)...)
			}
		} else {
			// Effects of a partially damaged system.
			switch i {
			case SysSensor, SysProbe:
				// No bookkeeping needed.
			case SysWarp:
				f2 = float64(percent) * ep.OrigMaxSpeed / 100
				ep.MaxSpeed -= f2
				if ep.MaxSpeed < 1.0 {
					ep.MaxSpeed = 1.0
					ep.Status[SysWarp] = 100
				}
			case SysComputer:
				messages = append(messages, CheckLocks(ep, percent, fed, r)...)
			}
		}
	}

	// The original's HISTORICAL-only "if (f1 > 43) ep->delay = 1."
	// self-destruct-on-heavy-hit rule was already disabled (via
	// #ifdef HISTORICAL) in the FreeBSD source; it is intentionally
	// not ported.
	return messages
}

// CheckLocks clears target locks (phaser, tube, and helm) with a
// probability of percent%, ported from check_locks() in damage.c. It
// is invoked after computer damage, since a damaged computer can
// spontaneously drop target locks.
func CheckLocks(ep *Ship, percent int, fed *Ship, r *Rand) []string {
	var messages []string

	var lostPhasers []int
	for i := range ep.Phasers {
		if ep.Phasers[i].Target != nil && r.Randm(100) <= percent {
			ep.Phasers[i].Target = nil
			if ep == fed {
				lostPhasers = append(lostPhasers, i+1)
			}
		}
	}
	if msg := lockLossMessage("Phaser", lostPhasers); msg != "" {
		messages = append(messages, msg)
	}

	var lostTubes []int
	for i := range ep.Tubes {
		if ep.Tubes[i].Target != nil && r.Randm(100) <= percent {
			ep.Tubes[i].Target = nil
			if ep == fed {
				lostTubes = append(lostTubes, i+1)
			}
		}
	}
	if msg := lockLossMessage("Tube", lostTubes); msg != "" {
		messages = append(messages, msg)
	}

	if ep.Target != nil && r.Randm(100) <= percent {
		ep.Target = nil
		ep.RelativeBear = 0
		if ep == fed {
			messages = append(messages, "Computer: "+ep.Name+" has lost helm lock")
		}
	}
	return messages
}

// lockLossMessage formats the "Computer: Phaser(s) 1, 3 have lost
// their target locks." style message from check_locks(), or returns
// "" if no locks were lost.
func lockLossMessage(kind string, indices []int) string {
	if len(indices) == 0 {
		return ""
	}
	msg := fmt.Sprintf("Computer: %s(s) %d", kind, indices[0])
	for _, idx := range indices[1:] {
		msg += fmt.Sprintf(", %d", idx)
	}
	if len(indices) > 1 {
		msg += " have lost their target locks."
	} else {
		msg += " has lost its target lock."
	}
	return msg
}
