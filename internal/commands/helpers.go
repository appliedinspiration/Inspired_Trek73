package commands

import (
	"fmt"
	"strings"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

// parseSelection parses a weapon-bank selection string, ported from
// the "all"/"ALL"-or-digit-string convention used throughout cmds1.c
// (fire_phasers, fire_tubes, lock_phasers, lock_tubes, turn_phasers,
// turn_tubes, load_tubes). If sel is "all" or "ALL", every bank from 0
// to count-1 is selected. Otherwise each character of sel is treated
// as a 1-based bank index ('1' selects bank 0, etc.); characters that
// don't resolve to a valid bank index are ignored, matching the
// original's silent "continue" on out-of-range indexes.
//
// The returned typed slice has one entry per bank (length count); a
// non-zero entry means that bank's index was named in sel, regardless
// of whether the bank was actually acted on afterward (e.g. because it
// was damaged) - this mirrors the original's typed[] array, which is
// used afterward purely to decide which banks to mention in
// damage/blind-spot warnings.
func parseSelection(sel string, count int) []int {
	typed := make([]int, count)
	if sel == "all" || sel == "ALL" {
		for i := 0; i < count; i++ {
			typed[i]++
		}
		return typed
	}
	for _, c := range sel {
		k := int(c) - '1'
		if k < 0 || k >= count {
			continue
		}
		typed[k]++
	}
	return typed
}

// pluralize appends "are"/"is" and a period-terminated suffix to msg,
// matching the "if (j > 1) ... else if (j == 1) ..." pattern repeated
// throughout misc.c/cmds*.c for building up a list of affected banks
// followed by a verb agreeing in number.
func pluralize(count int, suffix string) string {
	if count > 1 {
		return "are " + suffix
	}
	return "is " + suffix
}

// bankList renders "Computer: <Kind>(s) N, N, N" for every bank index
// in typed for which include(i) is true, followed by a plural-aware
// suffix, or "" if no banks match. Ported from the repeated
// pattern in check_p_damage/check_t_damage/check_p_turn/check_t_turn
// in misc.c.
func bankList(kind string, typed []int, include func(i int) bool, suffix string) string {
	var b strings.Builder
	count := 0
	for i, t := range typed {
		if t == 0 || !include(i) {
			continue
		}
		if count == 0 {
			fmt.Fprintf(&b, "Computer: %s(s) %d", kind, i+1)
		} else {
			fmt.Fprintf(&b, ", %d", i+1)
		}
		count++
	}
	if count == 0 {
		return ""
	}
	b.WriteString(" " + pluralize(count, suffix))
	return b.String()
}

// checkPhaserDamage reports any selected phaser banks that are damaged
// and thus unable to perform action, ported from check_p_damage() in
// misc.c.
func checkPhaserDamage(typed []int, sp *game.Ship, action string) []string {
	msg := bankList("Phaser", typed, func(i int) bool {
		return sp.Phasers[i].Status&game.PhaserDamaged != 0
	}, "damaged and unable to "+action+".")
	if msg == "" {
		return nil
	}
	return []string{msg}
}

// checkTubeDamage reports any selected tubes that are damaged and thus
// unable to perform action, ported from check_t_damage() in misc.c.
//
// The original checked "sp->tubes[i].status & P_DAMAGED" - the phaser
// damage bit instead of T_DAMAGED, the tube damage bit. Since both
// bits happen to have the same numeric value in this codebase, that
// bug was behaviorally inert, but it is fixed here per the project's
// "fix obviously unintended bugs" policy.
func checkTubeDamage(typed []int, sp *game.Ship, action string) []string {
	msg := bankList("Tube", typed, func(i int) bool {
		return sp.Tubes[i].Status&game.TubeDamaged != 0
	}, "damaged and unable to "+action+".")
	if msg == "" {
		return nil
	}
	return []string{msg}
}

// checkPhaserTurn reports any selected phaser banks that are currently
// pointed into the ship's blind arc, ported from check_p_turn() in
// misc.c. If firingOnly is true, only banks that are actually firing
// are considered (matching the "flag" parameter passed as 1 from
// fire_phasers and 0 from lock_phasers/turn_phasers).
func checkPhaserTurn(typed []int, sp *game.Ship, firingOnly bool) []string {
	msg := bankList("Phaser", typed, func(i int) bool {
		if firingOnly && sp.Phasers[i].Status&game.PhaserFiring == 0 {
			return false
		}
		return phaserPointsBlind(sp, i)
	}, "pointing into our blind side.")
	if msg == "" {
		return nil
	}
	return []string{msg}
}

// checkTubeTurn reports any selected tubes that are currently pointed
// into the ship's blind arc, ported from check_t_turn() in misc.c.
func checkTubeTurn(typed []int, sp *game.Ship, firingOnly bool) []string {
	msg := bankList("Tube", typed, func(i int) bool {
		if firingOnly && sp.Tubes[i].Status&game.TubeFiring == 0 {
			return false
		}
		return tubePointsBlind(sp, i)
	}, "pointing into our blind side.")
	if msg == "" {
		return nil
	}
	return []string{msg}
}

func phaserPointsBlind(sp *game.Ship, i int) bool {
	bank := &sp.Phasers[i]
	var k float64
	if bank.Target == nil {
		k = bank.Bearing
	} else {
		target := bank.Target
		tx, ty := target.X, target.Y
		if game.CantSee(target) {
			tx, ty = target.LastKnown.X, target.LastKnown.Y
		}
		bear := game.BearingTo(sp.X, tx, sp.Y, ty)
		k = bear - sp.Course
	}
	k = game.Rectify(k)
	return betw(k, float64(sp.PhaserBlindLeft), float64(sp.PhaserBlindRight)) && !sp.IsDead(game.SysEngineering)
}

func tubePointsBlind(sp *game.Ship, i int) bool {
	tube := &sp.Tubes[i]
	var k float64
	if tube.Target == nil {
		k = tube.Bearing
	} else {
		target := tube.Target
		tx, ty := target.X, target.Y
		if game.CantSee(target) {
			tx, ty = target.LastKnown.X, target.LastKnown.Y
		}
		bear := game.BearingTo(sp.X, tx, sp.Y, ty)
		k = bear - sp.Course
	}
	k = game.Rectify(k)
	return betw(k, float64(sp.TubeBlindLeft), float64(sp.TubeBlindRight)) && !sp.IsDead(game.SysEngineering)
}

// betw reports whether i is strictly between j and k, ported from the
// betw(i,j,k) macro in defines.h (see also game.betw, unexported).
func betw(i, j, k float64) bool {
	return j < i && i < k
}

// FormatPhaserControl, FormatPhaserBearing, and FormatPhaserLevel
// render the "Control"/"Turned"/"Level" rows shared by phaser_status()
// and print_damage() in cmds2.c/cmds3.c: one tab-separated entry per
// phaser bank, "damaged" for a damaged bank, "manual"/bearing for an
// unlocked bank, or the target's (truncated) name for a locked bank.
func FormatPhaserControl(sp *game.Ship) string {
	var b strings.Builder
	b.WriteString("Control: ")
	for i := range sp.Phasers {
		bank := &sp.Phasers[i]
		switch {
		case bank.Status&game.PhaserDamaged != 0:
			b.WriteString("\tdamaged")
		case bank.Target == nil:
			b.WriteString("\tmanual")
		default:
			fmt.Fprintf(&b, "\t%.7s", bank.Target.Name)
		}
	}
	return b.String()
}

func FormatPhaserBearing(sp *game.Ship) string {
	var b strings.Builder
	b.WriteString(" Turned: ")
	for i := range sp.Phasers {
		bank := &sp.Phasers[i]
		switch {
		case bank.Status&game.PhaserDamaged != 0:
			b.WriteString("\t")
		case bank.Target == nil:
			fmt.Fprintf(&b, "\t%.0f", bank.Bearing)
		default:
			b.WriteString("\tLOCKED")
		}
	}
	return b.String()
}

func FormatPhaserLevel(sp *game.Ship) string {
	var b strings.Builder
	b.WriteString("  Level: ")
	for i := range sp.Phasers {
		bank := &sp.Phasers[i]
		if bank.Status&game.PhaserDamaged != 0 {
			b.WriteString("\t")
		} else {
			fmt.Fprintf(&b, "\t%d", int(bank.Load))
		}
	}
	return b.String()
}

// FormatTubeControl, FormatTubeBearing, and FormatTubeLevel are the
// tube equivalents of the phaser formatters above, ported from
// tube_status() and print_damage()'s tube section in cmds2.c/cmds3.c.
func FormatTubeControl(sp *game.Ship) string {
	var b strings.Builder
	b.WriteString("Control: ")
	for i := range sp.Tubes {
		tube := &sp.Tubes[i]
		switch {
		case tube.Status&game.TubeDamaged != 0:
			b.WriteString("\tdamaged")
		case tube.Target == nil:
			b.WriteString("\tmanual")
		default:
			fmt.Fprintf(&b, "\t%.7s", tube.Target.Name)
		}
	}
	return b.String()
}

func FormatTubeBearing(sp *game.Ship) string {
	var b strings.Builder
	b.WriteString(" Turned: ")
	for i := range sp.Tubes {
		tube := &sp.Tubes[i]
		switch {
		case tube.Status&game.TubeDamaged != 0:
			b.WriteString("\t")
		case tube.Target == nil:
			fmt.Fprintf(&b, "\t%.0f", tube.Bearing)
		default:
			b.WriteString("\tLOCKED")
		}
	}
	return b.String()
}

func FormatTubeLevel(sp *game.Ship) string {
	var b strings.Builder
	b.WriteString("  Level: ")
	for i := range sp.Tubes {
		tube := &sp.Tubes[i]
		if tube.Status&game.TubeDamaged != 0 {
			b.WriteString("\t")
		} else {
			fmt.Fprintf(&b, "\t%d", int(tube.Load))
		}
	}
	return b.String()
}

// FiringAngles renders the "Firing angles: ..." line shared by
// phaser_status() and tube_status(), ported from the
// is_dead(sp,S_ENG) check repeated in both.
func FiringAngles(sp *game.Ship, blindLeft, blindRight int) string {
	if sp.IsDead(game.SysEngineering) {
		return "Firing angles: unrestricted."
	}
	return fmt.Sprintf("Firing angles: 0 - %d and %d - 360.", blindLeft, blindRight)
}

// minFloat is a small numeric helper used throughout the command
// layer, mirroring the min() macro in defines.h.
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
