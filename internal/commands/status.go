package commands

import (
	"fmt"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

// PhaserStatus renders sp's phaser control panel, ported from
// phaser_status() in cmds3.c.
func PhaserStatus(sp *game.Ship) []string {
	lines := []string{"Phasers", ""}
	lines = append(lines, FormatPhaserControl(sp))
	lines = append(lines, FormatPhaserBearing(sp))
	lines = append(lines, FormatPhaserLevel(sp))

	var drain string
	for i := range sp.Phasers {
		bank := &sp.Phasers[i]
		if bank.Status&game.PhaserDamaged != 0 {
			drain += "\t"
		} else {
			drain += fmt.Sprintf("\t%d", bank.Drain)
		}
	}
	lines = append(lines, "  Drain: "+drain)
	lines = append(lines, "", fmt.Sprintf("Firing percentage: %d", sp.PhaserFirePct))
	lines = append(lines, "", FiringAngles(sp, sp.PhaserBlindLeft, sp.PhaserBlindRight))
	return lines
}

// TubeStatus renders sp's torpedo control panel, ported from
// tube_status() in cmds3.c.
func TubeStatus(sp *game.Ship) []string {
	lines := []string{"Torpedos", ""}
	lines = append(lines, FormatTubeControl(sp))
	lines = append(lines, FormatTubeBearing(sp))
	lines = append(lines, FormatTubeLevel(sp))
	lines = append(lines, "", fmt.Sprintf("Launch speed: %d", sp.TubeLaunchSpd))
	lines = append(lines, fmt.Sprintf("  Time delay: %d", sp.TubeDelay))
	lines = append(lines, fmt.Sprintf("  Prox delay: %d", sp.TubeProximity))
	lines = append(lines, "", FiringAngles(sp, sp.TubeBlindLeft, sp.TubeBlindRight))
	return lines
}

// Survivors renders a survivors report for every ship in the battle,
// ported from survivors() in cmds3.c.
//
// The original had a bug here: the loop iterated ep over every ship,
// but every branch of the if/else chain inside it referenced sp (the
// ship whose command this is, always the player) instead of ep - so a
// destroyed enemy always printed "<player's name> -- destructed"
// rather than the enemy's own name and status. This port checks and
// reports each ep's own Complement/Name, fixing that bug.
func Survivors(st *game.State) []string {
	lines := []string{"", "Survivors reported:"}
	for _, ep := range st.Ships {
		switch {
		case ep.Complement < 0:
			lines = append(lines, fmt.Sprintf("   %s -- destructed", ep.Name))
		case game.CantSee(ep):
			lines = append(lines, fmt.Sprintf("   %s -- ???", ep.Name))
		default:
			lines = append(lines, fmt.Sprintf("   %s -- %d", ep.Name, ep.Complement))
		}
	}
	return lines
}
