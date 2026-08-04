package game

import "math"

// Distribute performs one turn's worth of power distribution and
// energy accounting for a single ship, ported from distribute() in
// dist.c. It must be called once per ship at the start of a turn,
// before the segment-by-segment movement simulation in MoveShips
// begins.
//
// It returns any player-visible status messages produced (shields
// down, shields fluctuating, cloak deactivated due to insufficient
// energy), honoring st.Shutup so each message is shown at most once
// per turn, matching the original's shutup[] array.
func Distribute(sp *Ship, st *State, hooks MovementHooks) []string {
	var messages []string
	fed := st.Player()
	isPlayer := sp == fed

	// Granularity of 1 second as far as this loop is concerned.
	for loop := 0; loop < int(SecondsPerTurn); loop++ {
		fuel := sp.Energy + sp.Regen // Slightly unrealistic, per the original comment.

		// Calculate negative phaser drains: a bank with a negative
		// drain setting feeds energy back into the ship.
		for i := range sp.Phasers {
			p := &sp.Phasers[i]
			load := p.Load
			drain := float64(p.Drain)
			if p.Status&PhaserDamaged != 0 || drain >= 0 || load <= 0 {
				continue
			}
			// Drain the lesser of either the current load (if less
			// than the drain) or the drain value.
			effLoad := math.Max(load+drain, 0)
			fuel += load - effLoad
			p.Load = effLoad
		}

		// Calculate shield drains.
		shield := 0.0
		for i := range sp.Shields {
			shield += sp.Shields[i].AttemptDrain
		}
		drain := math.Ceil(shield)

		// If all attempted drains are zero, or we have no fuel, our
		// shields are down.
		if shield*fuel == 0 && !st.Shutup.ShieldsFluctuating && isPlayer {
			messages = append(messages, st.Crew.Engineer+": "+st.Crew.Title+", our shields are down!")
			st.Shutup.ShieldsFluctuating = true
		}

		// If there's not enough fuel to sustain the drains, ration it
		// out in proportion to the attempted drains and say shields
		// are fluctuating.
		if drain <= fuel {
			fuel -= drain
			for i := range sp.Shields {
				sp.Shields[i].Drain = sp.Shields[i].AttemptDrain
			}
		} else {
			if !st.Shutup.ShieldsFluctuating && isPlayer {
				messages = append(messages, st.Crew.Engineer+": "+st.Crew.Title+", our shields are fluctuating!")
				st.Shutup.ShieldsFluctuating = true
			}
			for i := range sp.Shields {
				sp.Shields[i].Drain = sp.Shields[i].AttemptDrain * fuel / drain
			}
			fuel = 0
		}

		// Calculate cloaking device drains. If there's insufficient
		// energy to run the device, it's turned off completely.
		if cantSee(sp) {
			if fuel < float64(sp.CloakEnergy) {
				if isPlayer {
					sp.Cloaking = CloakOff
					messages = append(messages, st.Crew.Engineer+": "+st.Crew.Title+
						", there's not enough energy to keep our cloaking device activated.")
				} else {
					hooks.ECloakOff(sp, fed)
				}
			} else {
				fuel -= float64(sp.CloakEnergy)
			}
		}

		// Calculate positive phaser drains: load phasers either
		// enough to top them off, or by the full drain.
		for i := range sp.Phasers {
			if fuel <= 0 {
				break
			}
			p := &sp.Phasers[i]
			load := p.Load
			drain := float64(p.Drain)
			if p.Status&PhaserDamaged != 0 || load >= MaxPhaserCharge || drain <= 0 {
				continue
			}
			effLoad := math.Min(MaxPhaserCharge, load+math.Min(drain, fuel))
			fuel -= effLoad - load
			p.Load = effLoad
		}

		// Balance the level of energy with the number of pods.
		sp.Energy = math.Min(fuel, sp.Pods)
	}
	return messages
}
