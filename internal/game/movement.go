package game

import (
	"fmt"
	"math"
)

// MoveShips advances the battle by one full turn (SecondsPerTurn
// seconds), simulating movement, firing, fuses, and course/speed
// changes in segment-sized (SegmentSeconds) steps. Ported from
// move_ships() in moveships.c.
//
// Distribute must already have been called for every ship this turn
// (matching the "for (i=0;i<=shipnum;i++) distribute(shiplist[i])"
// loop at the top of the original move_ships()); MoveShips does not
// call it automatically, so callers can collect Distribute's messages
// separately before movement messages.
//
// Detonated torpedoes/probes/engineering sections and destroyed ships
// are reported to hooks (TorpDetonate/ShipDetonate) rather than
// removed here; the caller is responsible for actually removing
// detonated objects from st.Objects once combat/damage resolution (a
// later phase) has processed them.
func MoveShips(st *State, hooks MovementHooks) []string {
	var messages []string
	fed := st.Player()

	// The value 100 is the number of Megameters per second per
	// warp-factor.
	const mPerSegment = SegmentSeconds * 100.0
	iterations := int(math.Floor(SecondsPerTurn/SegmentSeconds + 0.5)) // What a crock!

	for i := 0; i < iterations; i++ {
		for _, sp := range st.Ships {
			if sp.IsDead(SysDead) {
				continue
			}
			if i%sp.PhaserFiringDelay == 0 {
				hooks.PhaserFiring(sp)
			}
			if i%sp.TorpedoFiringDelay == 0 {
				hooks.TorpedoFiring(sp)
			}

			// Self-destruct countdown.
			sp.Delay -= SegmentSeconds
			if sp.Delay < SegmentSeconds/2 {
				hooks.ShipDetonate(sp)
				continue
			}

			if sp.IsDead(SysDead) {
				continue
			}
			moveOneShip(sp, fed, mPerSegment, st.Shutup, messagesAppender(&messages))
		}

		for _, obj := range st.Objects {
			if obj.Detonated {
				continue
			}
			// Time fuse.
			obj.TimeDelay -= SegmentSeconds
			if obj.TimeDelay <= 0 {
				obj.Detonated = true
				hooks.TorpDetonate(obj)
				continue
			}

			// Proximity fuse.
			if obj.Proximity != 0 {
				if proximityTriggered(st, obj) {
					obj.Detonated = true
					hooks.TorpDetonate(obj)
					continue
				}
			}

			moveOneObject(obj, mPerSegment)
		}
	}
	return messages
}

// messagesAppender returns a func(string) that appends to *msgs,
// letting moveOneShip report messages without importing a logger.
func messagesAppender(msgs *[]string) func(string) {
	return func(s string) { *msgs = append(*msgs, s) }
}

// proximityTriggered reports whether any other visible ship has come
// within a torpedo/probe's proximity-fuse range, ported from the
// proximity-fuse loop in move_ships() (moveships.c). The original
// checked both I_SHIP and I_ENG list entries (both of which store a
// *ship pointer in this port's model: SpaceObject entries of Type ==
// ObjectEngineering don't have a corresponding *Ship, so unlike the
// original - which incorrectly treated an I_ENG entry's struct
// torpedo* as a struct ship* - this only checks live ships, which is
// the intended, bug-fixed behavior.
func proximityTriggered(st *State, obj *SpaceObject) bool {
	for _, sp := range st.Ships {
		if sp == obj.From {
			continue
		}
		if cantSee(sp) {
			continue
		}
		if rangeFind(obj.X, sp.X, obj.Y, sp.Y) < obj.Proximity {
			return true
		}
	}
	return false
}

// moveOneShip performs the per-segment movement simulation for a
// single ship: fuel consumption, warp-drive burnout, destroyed-warp
// speed capping, autopilot course tracking, turn-rate limiting,
// acceleration, and position update. Ported from the "if (sp) { ... }"
// branch of move_ships()'s per-item loop in moveships.c.
func moveOneShip(sp, fed *Ship, mPerSegment float64, shutup *Shutup, emit func(string)) {
	x, y := sp.X, sp.Y
	warp := sp.Warp
	if math.Abs(sp.NewWarp) > sp.MaxSpeed {
		if sp.NewWarp > 0 {
			sp.NewWarp = sp.MaxSpeed
		} else {
			sp.NewWarp = -sp.MaxSpeed
		}
	}
	newWarp := sp.NewWarp
	course := sp.Course
	newCourse := sp.NewCourse
	target := sp.Target
	energy := sp.Energy

	// Fuel consumption.
	var fuelUse float64
	if math.Abs(warp) > 1.0 {
		fuelUse = math.Abs(warp) * sp.Eff * SegmentSeconds
	}
	if fuelUse > energy {
		if !shutup.BurnoutShown(sp.ID) && !sp.IsDead(SysWarp) && canSee(sp) {
			emit(sp.Name + "'s warp drive burning out.")
			shutup.SetBurnoutShown(sp.ID)
		}
		if warp < 0.0 {
			newWarp = -0.99
		} else {
			newWarp = 0.99
		}
		energy = 0
	} else {
		energy -= fuelUse
	}

	// Destroyed warp drive.
	if sp.IsDead(SysWarp) && math.Abs(warp) > 1.0 {
		if warp < 0.0 {
			newWarp = -0.99
		} else {
			newWarp = 0.99
		}
	}

	// Automatic pilot.
	if target != nil {
		if target.IsDead(SysDead) {
			if sp == fed && !shutup.Disengage && !sp.IsDead(SysDead) {
				emit(sp.Name + "'s autopilot disengaging.")
				shutup.Disengage = true
			}
			newCourse = course
			target = nil
			sp.RelativeBear = 0.0
		} else {
			var j float64
			if cantSee(target) {
				j = bearingTo(x, target.LastKnown.X, y, target.LastKnown.Y)
			} else {
				j = bearingTo(x, target.X, y, target.Y)
			}
			j = rectify(j - sp.RelativeBear)
			newCourse = j
		}
	}

	// Turn rate.
	if course != newCourse {
		j := rectify(newCourse - course)
		if j > 180 {
			j -= 360
		}
		// Maximum degrees turned in one turn.
		k := (sp.MaxSpeed + 2.0 - math.Abs(warp)) * sp.DegPerTurn * SegmentSeconds
		// A ship with no warp drive is less maneuverable.
		if sp.IsDead(SysWarp) {
			k /= 2
		}
		sign := 1.0
		if j < 0.0 {
			sign = -1.0
		}
		k = course + sign*math.Min(math.Abs(j), k)
		course = rectify(k)
	}

	// Acceleration.
	tmpf := newWarp - warp
	d0 := math.Abs(tmpf)
	sign := 1.0
	if tmpf < 0.0 {
		sign = -1.0
	}
	warp += sign * math.Sqrt(d0) * SegmentSeconds
	rad := toRadians(course)
	x += int(warp * math.Cos(rad) * mPerSegment)
	y += int(warp * math.Sin(rad) * mPerSegment)

	// Projected position (cloaked), tracked separately from the
	// ship's true position so enemy sensors can only ever see the
	// last-known position and its extrapolated track.
	if cantSee(sp) {
		radLK := toRadians(sp.LastKnown.Course)
		sp.LastKnown.X += int(sp.LastKnown.Warp * math.Cos(radLK) * SegmentSeconds)
		sp.LastKnown.Y += int(sp.LastKnown.Warp * math.Sin(radLK) * SegmentSeconds)
	}

	// Commit.
	sp.X = x
	sp.Y = y
	sp.Warp = warp
	sp.NewWarp = newWarp
	sp.Course = rectify(course)
	sp.NewCourse = rectify(newCourse)
	sp.Energy = energy
	sp.Target = target
}

// moveOneObject performs the per-segment movement simulation for a
// torpedo, probe, or jettisoned engineering section: target homing
// (for probes with a live target), acceleration, and position update.
// Ported from the "else" (tp) branch of move_ships()'s per-item loop
// in moveships.c.
//
// Unlike ships, objects have no relative-bearing offset, no turn-rate
// limit (course snaps directly to the homing bearing every segment,
// matching the original's "if (tp) course = newcourse;"), no fuel
// consumption, and no cloak tracking. Also unlike ships, when an
// object's target dies the object keeps homing on the target's last
// position instead of disengaging - the original's autopilot-abort
// message only ever applied to sp (ship) targets, and this preserves
// that asymmetry rather than treating it as a bug, since flying
// toward a target's last known position/wreckage is plausible
// intended behavior for an unguided munition.
func moveOneObject(obj *SpaceObject, mPerSegment float64) {
	x, y := obj.X, obj.Y
	warp := obj.Speed
	newWarp := obj.NewSpeed
	course := obj.Course
	target := obj.Target

	if target != nil {
		var j float64
		if cantSee(target) {
			j = bearingTo(x, target.LastKnown.X, y, target.LastKnown.Y)
		} else {
			j = bearingTo(x, target.X, y, target.Y)
		}
		course = j
	}

	tmpf := newWarp - warp
	d0 := math.Abs(tmpf)
	sign := 1.0
	if tmpf < 0.0 {
		sign = -1.0
	}
	warp += sign * math.Sqrt(d0) * SegmentSeconds
	rad := toRadians(course)
	x += int(warp * math.Cos(rad) * mPerSegment)
	y += int(warp * math.Sin(rad) * mPerSegment)

	obj.X = x
	obj.Y = y
	obj.Speed = warp
	obj.Course = rectify(course)
	obj.Target = target
}

// CheckTargets validates every ship's target lock, and every phaser
// and tube target lock, clearing any that point at a now-destroyed
// ship (and reporting a player-visible disengage message when it's
// the player's own lock). It also updates the bearing tracking for
// any still-locked phaser/tube target. Ported from check_targets() in
// moveships.c.
//
// The original only decremented cloak_delay for the last ship
// processed in its loop (a bug: the "if (sp->cloak_delay > 0)
// sp->cloak_delay--;" statement sat outside the per-ship loop body,
// so it only ever affected whichever ship the loop variable "sp"
// happened to still reference after its final iteration). This port
// fixes that by decrementing every ship's own CloakDelay inside the
// loop, which is very likely the intended behavior.
func CheckTargets(st *State) []string {
	var messages []string
	fed := st.Player()

	for _, sp := range st.Ships {
		if sp.IsDead(SysDead) {
			continue
		}
		if target := sp.Target; target != nil && target.IsDead(SysDead) {
			if sp == fed && !st.Shutup.Disengage {
				messages = append(messages, "   helm lock disengaging")
				st.Shutup.Disengage = true
			}
			sp.Target = nil
			sp.RelativeBear = 0.0
		}

		for j := range sp.Phasers {
			p := &sp.Phasers[j]
			target := p.Target
			if target != nil && target.IsDead(SysDead) {
				if sp == fed && !st.Shutup.Phaser[j] && !sp.IsDead(SysDead) {
					messages = append(messages, phaserDisengageMessage(j))
					st.Shutup.Phaser[j] = true
				}
				p.Target = nil
			} else if target != nil {
				tx, ty := target.X, target.Y
				if cantSee(target) {
					tx, ty = target.LastKnown.X, target.LastKnown.Y
				}
				p.Bearing = rectify(bearingTo(sp.X, tx, sp.Y, ty) - sp.Course)
			}
		}

		for j := range sp.Tubes {
			t := &sp.Tubes[j]
			target := t.Target
			if target != nil && target.IsDead(SysDead) {
				if sp == fed && !st.Shutup.Tube[j] && !sp.IsDead(SysDead) {
					messages = append(messages, tubeDisengageMessage(j))
					st.Shutup.Tube[j] = true
				}
				t.Target = nil
			} else if target != nil {
				tx, ty := target.X, target.Y
				if cantSee(target) {
					tx, ty = target.LastKnown.X, target.LastKnown.Y
				}
				t.Bearing = rectify(bearingTo(sp.X, tx, sp.Y, ty) - sp.Course)
			}
		}

		if sp.CloakDelay > 0 {
			sp.CloakDelay--
		}
	}
	return messages
}

func phaserDisengageMessage(index int) string {
	return "  phaser " + itoa(index+1) + " disengaging"
}

func tubeDisengageMessage(index int) string {
	return "  tube " + itoa(index+1) + " disengaging"
}

// itoa is a tiny local integer-to-string helper so this file doesn't
// need to import strconv solely for two call sites.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// MiscTimers advances the self-destruct warning, resets the probe
// launch status, and ages ruse/bluff/surrender timers. Ported from
// misc_timers() in moveships.c.
func MiscTimers(st *State) []string {
	var messages []string
	fed := st.Player()

	// Self-destruct warning.
	if fed.Delay < 1000.0 && fed.Delay > SegmentSeconds/2 {
		if fed.IsDead(SysComputer) {
			messages = append(messages, "Science: Self-destruct has been aborted due to computer damage")
			fed.Delay = 10000.0
		} else {
			messages = append(messages, fmt.Sprintf("Computer: %5.2f seconds to self destruct.", fed.Delay))
		}
	}

	fed.ProbeLauncherStatus = ProbeNormal

	// Ruses, bluffs, surrenders: these counters are owned by whatever
	// higher-level "special actions" state tracks them (special.c);
	// MiscTimers only ages them here, matching the original.
	st.AgeSpecialActionTimers()

	return messages
}
