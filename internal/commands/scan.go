package commands

import (
	"fmt"
	"strings"

	"github.com/appliedinspiration/inspired_trek73/internal/data"
	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

// PrintDamage renders a full damage/status report for ep, ported from
// print_damage() in cmds2.c. self_scan()/scan() both call through to
// this (self_scan always targeting the player's own ship, scan for
// other ships).
func PrintDamage(ep *game.Ship) []string {
	var lines []string
	lines = append(lines, "", fmt.Sprintf("Damages to the %s", ep.Name))
	for i := 0; i < game.NumDamageSystems+1; i++ {
		if ep.IsDead(i) {
			lines = append(lines, data.SystemDestroyedMessages[i]+".")
		} else if i < game.NumDamageSystems && ep.Status[i] != 0 {
			lines = append(lines, fmt.Sprintf("%s damaged %d%%", data.SystemNames[i], ep.Status[i]))
		}
	}
	lines = append(lines, fmt.Sprintf("Survivors: %d", ep.Complement))
	if ep.Target == nil {
		lines = append(lines, "Helm lock: none.")
	} else {
		lines = append(lines, fmt.Sprintf("Helm lock: %s", ep.Target.Name))
	}
	lines = append(lines, "", "Phasers Control")
	lines = append(lines, FormatPhaserControl(ep))
	lines = append(lines, FormatPhaserBearing(ep))
	lines = append(lines, FormatPhaserLevel(ep))
	lines = append(lines, "", "Tubes\tcontrol")
	lines = append(lines, FormatTubeControl(ep))
	lines = append(lines, FormatTubeBearing(ep))
	lines = append(lines, FormatTubeLevel(ep))

	var shieldLevels, shieldDrains strings.Builder
	shieldLevels.WriteString("Shields\t levels")
	shieldDrains.WriteString("\t drains")
	for i := range ep.Shields {
		level := int(100 * ep.Shields[i].Eff * ep.Shields[i].Drain)
		fmt.Fprintf(&shieldLevels, "\t%-2d", level)
		fmt.Fprintf(&shieldDrains, "\t%-4.2f", ep.Shields[i].AttemptDrain)
	}
	lines = append(lines, shieldLevels.String(), shieldDrains.String())
	lines = append(lines, fmt.Sprintf("Efficiency: %4.1f\tFuel remaining: %d", ep.Eff, int(ep.Energy)))
	lines = append(lines, fmt.Sprintf("Regeneration: %4.1f\tFuel capacity: %d", ep.Regen, int(ep.Pods)))
	return lines
}

// SelfScan is a damage report on the player's own ship, ported from
// self_scan() in cmds2.c.
func SelfScan(sp *game.Ship) []string {
	return PrintDamage(sp)
}

// ScanTorpedo renders a single-line report on a probe or jettisoned
// engineering section, ported from scan_torpedo() in cmds2.c
// (unguided torpedoes are never scannable this way, matching the
// original: fire_tubes never allocates a scannable ID visible to this
// command since SHOWTORP is not enabled in this build).
func ScanTorpedo(obj *game.SpaceObject) []string {
	kind := "engng"
	if obj.Type == game.ObjectProbe {
		kind = "probe"
	}
	target := "NONE"
	if obj.Target != nil {
		target = obj.Target.Name
	}
	lines := []string{"", "object id time  prox units target"}
	lines = append(lines, fmt.Sprintf("%-7s%2d  %2.1f %4d   %3d   %s", kind, obj.ID, obj.TimeDelay, obj.Proximity, obj.Fuel, target))
	return lines
}

// ScanResult carries what a Scan() call resolved, since a scan can
// target either a ship or a probe/engineering-section object (never
// both), ported from the item/tp/ep locals in scan() (cmds2.c).
type ScanResult struct {
	Ship   *game.Ship
	Object *game.SpaceObject
}

// Scan resolves and validates a "scan [who]" request, ported from
// scan() in cmds2.c. who is exactly the raw string the player typed
// (e.g. "Klingon", "#Klingon" for its jettisoned engineering, or a
// probe number like "3"); the caller is expected to have already
// prompted for it. Returns the resolved target and any messages; if
// Ship and Object are both nil, the scan failed and the messages
// explain why.
func Scan(st *game.State, r *game.Rand, sp *game.Ship, who string) (ScanResult, []string) {
	if who == "" {
		return ScanResult{}, nil
	}
	var result ScanResult

	if who[0] == '#' {
		name := who[1:]
		if name == "" {
			return ScanResult{}, []string{fmt.Sprintf("%s: %s, scan whose engineering?", st.Crew.Science, st.Crew.Title)}
		}
		ep := st.ShipByName(name)
		if ep == nil {
			return ScanResult{}, []string{fmt.Sprintf("%s: %s, no such ship as the %s.", st.Crew.Science, st.Crew.Title, name)}
		}
		var obj *game.SpaceObject
		for _, o := range st.Objects {
			if o.Type == game.ObjectEngineering && o.From == ep {
				obj = o
				break
			}
		}
		if obj == nil {
			return ScanResult{}, []string{fmt.Sprintf("%s:  %s, the %s has not jettisoned it's engineering.", st.Crew.Science, st.Crew.Title, ep.Name)}
		}
		result.Object = obj
	} else if probeNum := atoiOrZero(who); probeNum > 0 {
		var obj *game.SpaceObject
		for _, o := range st.Objects {
			if o.Type == game.ObjectProbe && o.ID == probeNum {
				obj = o
				break
			}
		}
		if obj == nil {
			return ScanResult{}, []string{fmt.Sprintf("%s: %s, there is no probe %d", st.Crew.Science, st.Crew.Title, probeNum)}
		}
		result.Object = obj
	} else {
		ep := st.ShipByName(who)
		if ep == nil {
			return ScanResult{}, []string{fmt.Sprintf("%s: %s, no such ship as the %s.", st.Crew.Science, st.Crew.Title, who)}
		}
		if game.CantSee(ep) {
			return ScanResult{}, []string{fmt.Sprintf("%s:  %s, I am unable to scan the %s.", st.Crew.Science, st.Crew.Title, ep.Name)}
		}
		result.Ship = ep
	}

	scanningSelf := result.Ship == sp
	// The sensor-health checks below always apply when the scan
	// resolved to a probe/engineering-section object (result.Object),
	// matching the original's "sp != ep" test, since ep is NULL for a
	// probe scan and (except for the corner case of a ship scanning
	// its own already-jettisoned engineering) different from sp for an
	// engineering scan.
	if !scanningSelf {
		if sp.IsDead(game.SysSensor) {
			return ScanResult{}, []string{fmt.Sprintf("%s: The sensors are damaged, Captain.", st.Crew.Science)}
		}
		if !sp.SystemWorks(r, game.SysSensor) {
			return ScanResult{}, []string{fmt.Sprintf("%s: %s, sensors are temporarily out.", st.Crew.Science, st.Crew.Title)}
		}
	}
	if scanningSelf {
		return ScanResult{}, []string{fmt.Sprintf("%s: Captain, don't you mean 'Damage Report'?", st.Crew.Science)}
	}

	var messages []string
	if result.Ship != nil {
		messages = PrintDamage(result.Ship)
	} else {
		messages = ScanTorpedo(result.Object)
	}
	return result, messages
}

func atoiOrZero(s string) int {
	n := 0
	neg := false
	i := 0
	if i < len(s) && (s[i] == '-' || s[i] == '+') {
		neg = s[i] == '-'
		i++
	}
	if i == len(s) {
		return 0
	}
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	if neg {
		n = -n
	}
	return n
}
