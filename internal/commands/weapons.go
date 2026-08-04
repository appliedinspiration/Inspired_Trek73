package commands

import (
	"fmt"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

// FirePhasers fires the selected phaser banks, ported from
// fire_phasers() in cmds1.c. banks is the raw selection string (e.g.
// "13", "all"); spread is the requested phaser spread in degrees. If
// spread is out of range, nothing happens, matching the original
// (which silently returned 0 without printing anything).
//
// Per the project's pure-function command design, the caller (the
// interactive/CLI layer, ported later) is responsible for prompting
// for banks/spread and passing them in already parsed; this function
// never performs I/O itself.
func FirePhasers(sp *game.Ship, banks string, spread int) []string {
	if spread < game.MinPhaserSpread || spread > game.MaxPhaserSpread {
		return nil
	}
	sp.PhaserSpread = spread

	typed := parseSelection(banks, len(sp.Phasers))
	for i, t := range typed {
		if t == 0 {
			continue
		}
		bank := &sp.Phasers[i]
		if bank.Status&(game.PhaserDamaged|game.PhaserFiring) != 0 {
			continue
		}
		bank.Status |= game.PhaserFiring
	}
	var messages []string
	messages = append(messages, checkPhaserDamage(typed, sp, "fire")...)
	messages = append(messages, checkPhaserTurn(typed, sp, true)...)
	return messages
}

// FireTubes fires the selected torpedo tubes, ported from fire_tubes()
// in cmds1.c. banks is the raw selection string.
func FireTubes(sp *game.Ship, banks string) []string {
	typed := parseSelection(banks, len(sp.Tubes))
	for i, t := range typed {
		if t == 0 {
			continue
		}
		tube := &sp.Tubes[i]
		if tube.Status&(game.TubeDamaged|game.TubeFiring) != 0 {
			continue
		}
		tube.Status |= game.TubeFiring
	}
	var messages []string
	messages = append(messages, checkTubeDamage(typed, sp, "fire")...)
	messages = append(messages, checkTubeTurn(typed, sp, true)...)

	var unloaded []int
	for i, t := range typed {
		if t == 0 || sp.Tubes[i].Status&game.TubeFiring == 0 {
			continue
		}
		if sp.Tubes[i].Load == 0 {
			unloaded = append(unloaded, i+1)
		}
	}
	if len(unloaded) > 0 {
		msg := "Computer: Tube(s) "
		for i, n := range unloaded {
			if i > 0 {
				msg += ", "
			}
			msg += fmt.Sprintf("%d", n)
		}
		msg += " " + pluralize(len(unloaded), "not loaded.")
		messages = append(messages, msg)
	}
	return messages
}

// LockPhasers locks the selected phaser banks onto target, ported from
// lock_phasers() in cmds1.c. target must already be resolved (e.g. via
// game.State.ShipByName) by the caller.
func LockPhasers(st *game.State, r *game.Rand, sp *game.Ship, banks string, target *game.Ship) []string {
	if sp.IsDead(game.SysComputer) {
		return []string{fmt.Sprintf("%s:  Impossible %s, our computer is dead.", st.Crew.Science, st.Crew.Title)}
	}
	if !sp.SystemWorks(r, game.SysComputer) {
		return []string{fmt.Sprintf("%s:  Our computer is temporarily buggy", st.Crew.Science)}
	}
	if target == nil {
		return nil
	}
	if game.CantSee(target) {
		return []string{fmt.Sprintf("%s:  %s, unable to lock phasers onto %s.", st.Crew.Nav, st.Crew.Title, target.Name)}
	}
	typed := parseSelection(banks, len(sp.Phasers))
	for i, t := range typed {
		if t == 0 {
			continue
		}
		bank := &sp.Phasers[i]
		if bank.Status&game.PhaserDamaged != 0 {
			continue
		}
		bank.Target = target
	}
	var messages []string
	messages = append(messages, checkPhaserDamage(typed, sp, "lock")...)
	messages = append(messages, checkPhaserTurn(typed, sp, false)...)
	return messages
}

// LockTubes locks the selected tubes onto target, ported from
// lock_tubes() in cmds1.c.
func LockTubes(st *game.State, r *game.Rand, sp *game.Ship, banks string, target *game.Ship) []string {
	if sp.IsDead(game.SysComputer) {
		return []string{fmt.Sprintf("%s:  Impossible %s, our computer is dead.", st.Crew.Science, st.Crew.Title)}
	}
	if !sp.SystemWorks(r, game.SysComputer) {
		return []string{fmt.Sprintf("%s:  Our computer is temporarily buggy", st.Crew.Science)}
	}
	if target == nil {
		return nil
	}
	if game.CantSee(target) {
		return []string{fmt.Sprintf("%s:  %s, unable to lock tubes onto %s.", st.Crew.Nav, st.Crew.Title, target.Name)}
	}
	typed := parseSelection(banks, len(sp.Tubes))
	for i, t := range typed {
		if t == 0 {
			continue
		}
		tube := &sp.Tubes[i]
		if tube.Status&game.TubeDamaged != 0 {
			continue
		}
		tube.Target = target
	}
	var messages []string
	messages = append(messages, checkTubeDamage(typed, sp, "lock")...)
	messages = append(messages, checkTubeTurn(typed, sp, false)...)
	return messages
}

// TurnPhasers manually rotates the selected phaser banks to a relative
// bearing, unlocking any target, ported from turn_phasers() in
// cmds1.c. If bearing is out of [0,360], nothing happens.
func TurnPhasers(sp *game.Ship, banks string, bearing float64) []string {
	if bearing < 0.0 || bearing > 360.0 {
		return nil
	}
	typed := parseSelection(banks, len(sp.Phasers))
	for i, t := range typed {
		if t == 0 {
			continue
		}
		bank := &sp.Phasers[i]
		if bank.Status&game.PhaserDamaged != 0 {
			continue
		}
		bank.Target = nil
		bank.Bearing = bearing
	}
	var messages []string
	messages = append(messages, checkPhaserDamage(typed, sp, "turn")...)
	messages = append(messages, checkPhaserTurn(typed, sp, false)...)
	return messages
}

// TurnTubes manually rotates the selected tubes to a relative bearing,
// unlocking any target, ported from turn_tubes() in cmds1.c.
func TurnTubes(sp *game.Ship, banks string, bearing float64) []string {
	if bearing < 0.0 || bearing > 360.0 {
		return nil
	}
	typed := parseSelection(banks, len(sp.Tubes))
	for i, t := range typed {
		if t == 0 {
			continue
		}
		tube := &sp.Tubes[i]
		if tube.Status&game.TubeDamaged != 0 {
			continue
		}
		tube.Target = nil
		tube.Bearing = bearing
	}
	var messages []string
	messages = append(messages, checkTubeDamage(typed, sp, "turn")...)
	messages = append(messages, checkTubeTurn(typed, sp, false)...)
	return messages
}

// LoadTubes loads or unloads the selected tubes, transferring energy
// between the ship's main store and each tube, ported from
// load_tubes() in cmds1.c. load selects loading (true) vs. unloading
// (false).
func LoadTubes(st *game.State, sp *game.Ship, load bool, banks string) []string {
	typed := parseSelection(banks, len(sp.Tubes))
	for i := range sp.Tubes {
		if typed[i] == 0 {
			continue
		}
		tube := &sp.Tubes[i]
		if tube.Status&game.TubeDamaged != 0 {
			continue
		}
		if load {
			j := minFloat(sp.Energy, game.MaxTubeCharge-tube.Load)
			if j <= 0 {
				continue
			}
			sp.Energy -= j
			sp.Pods -= j
			tube.Load += j
		} else {
			j := tube.Load
			if j <= 0 {
				continue
			}
			sp.Energy += j
			sp.Pods += j
			tube.Load = 0
		}
	}
	msg := fmt.Sprintf("%s: Tubes now ", st.Crew.Engineer)
	for i := range sp.Tubes {
		if sp.Tubes[i].Status&game.TubeDamaged != 0 {
			msg += " -- "
		} else {
			msg += fmt.Sprintf(" %-2d ", int(sp.Tubes[i].Load))
		}
	}
	msg += fmt.Sprintf(" energy at %d/%d", int(sp.Energy), int(sp.Pods))
	return []string{msg}
}
