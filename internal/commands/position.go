package commands

import (
	"fmt"
	"strings"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

// PosReport renders a position report of every ship and
// probe/jettisoned-engineering section in the battle, ported from
// pos_report() in cmds2.c. Torpedoes are never included, matching the
// original build (the "#ifdef SHOWTORP" branches are not compiled in
// by default).
//
// Unlike the original, which computed a dynamic column width from the
// longest name, this renders a simpler fixed-width table; the
// information conveyed is identical.
func PosReport(st *game.State, sp *game.Ship) []string {
	nameWidth := posReportNameWidth(st)
	lines := []string{
		fmt.Sprintf("%s                     abs           rel   rev rel", strings.Repeat(" ", nameWidth)),
		fmt.Sprintf("%s class warp course bearing range bearing bearing", strings.Repeat(" ", nameWidth)),
	}
	for _, sp1 := range st.Ships {
		if sp1.IsDead(game.SysDead) {
			continue
		}
		x, y := sp1.X, sp1.Y
		speed, course := sp1.Warp, sp1.Course
		who := sp1.Name
		if game.CanSee(sp1) {
			sp1.LastKnown = game.LastKnownPosition{
				X: sp1.X, Y: sp1.Y, Warp: sp1.Warp, Course: sp1.Course,
				Bear:  game.BearingTo(sp.X, sp1.X, sp.Y, sp1.Y),
				Range: game.RangeFind(sp.X, sp1.X, sp.Y, sp1.Y),
			}
		} else {
			x, y = sp1.LastKnown.X, sp1.LastKnown.Y
			speed, course = sp1.LastKnown.Warp, sp1.LastKnown.Course
			who += " *"
		}
		if sp1 == sp {
			lockLine := ""
			if sp.Target != nil {
				lockLine = fmt.Sprintf("helm locked on %s", sp.Target.Name)
			}
			lines = append(lines, fmt.Sprintf("%-*s%5s%6.1f   %3.0f   %s", nameWidth, who, sp1.Class, speed, course, lockLine))
			continue
		}
		bear := game.BearingTo(sp.X, x, sp.Y, y)
		rng := game.RangeFind(sp.X, x, sp.Y, y)
		relBear := game.Rectify(round(bear - sp.Course))
		revRelBear := game.Rectify(round(bear + 180.0 - course))
		lines = append(lines, fmt.Sprintf("%-*s%5s%6.1f   %3.0f    %3.0f   %5d   %3.0f     %3.0f",
			nameWidth,
			who, sp1.Class, speed, course, bear, rng, relBear, revRelBear))
	}
	for _, obj := range st.Objects {
		if obj.Type == game.ObjectTorpedo {
			continue
		}
		who := objectLabel(obj)
		bear := game.BearingTo(sp.X, obj.X, sp.Y, obj.Y)
		rng := game.RangeFind(sp.X, obj.X, sp.Y, obj.Y)
		relBear := game.Rectify(round(bear - sp.Course))
		revRelBear := game.Rectify(round(bear + 180.0 - obj.Course))
		lines = append(lines, fmt.Sprintf("%-*s%5s%6.1f   %3.0f    %3.0f   %5d   %3.0f     %3.0f",
			nameWidth, who, "", obj.Speed, obj.Course, bear, rng, relBear, revRelBear))
	}
	return lines
}

func posReportNameWidth(st *game.State) int {
	maxLen := 0
	for _, sp := range st.Ships {
		if n := len(sp.Name); n > maxLen {
			maxLen = n
		}
	}
	for _, obj := range st.Objects {
		if obj.Type == game.ObjectTorpedo {
			continue
		}
		if labelLen := len(objectLabel(obj)); labelLen > maxLen {
			maxLen = labelLen
		}
	}
	return maxLen + 2
}

func objectLabel(obj *game.SpaceObject) string {
	from := "unknown"
	if obj.From != nil {
		from = obj.From.Name
	}
	switch obj.Type {
	case game.ObjectProbe:
		return fmt.Sprintf("%s probe %d", from, obj.ID)
	case game.ObjectEngineering:
		return fmt.Sprintf("%s engineering", from)
	default:
		return from
	}
}

// round rounds x to the nearest integer, mirroring round() in subs.c
// (used to snap computed bearings to whole degrees for display).
func round(x float64) float64 {
	if x < 0 {
		return float64(int(x - 0.5))
	}
	return float64(int(x + 0.5))
}

// PosDisplay renders an ASCII scope centered on sp showing every other
// ship/object within range megameters, ported from pos_display() in
// cmds2.c.
func PosDisplay(sp *game.Ship, st *game.State, r *game.Rand, rng int) []string {
	if sp.IsDead(game.SysSensor) {
		return []string{"Sensors are damaged."}
	}
	if !sp.SystemWorks(r, game.SysSensor) {
		return []string{"Sensors are temporarily inoperative."}
	}
	if rng < game.MinSensorRange || rng > game.MaxSensorRange {
		return nil
	}

	const hpitch, vpitch = 10, 6
	xScale := rng / hpitch
	yScale := rng / vpitch
	if xScale == 0 {
		xScale = 1
	}
	if yScale == 0 {
		yScale = 1
	}

	grid := make([][]byte, 2*vpitch+1)
	for i := range grid {
		grid[i] = make([]byte, 2*hpitch+1)
		for j := range grid[i] {
			if i == 0 || i == 2*vpitch {
				grid[i][j] = '-'
			} else if j == 0 || j == 2*hpitch {
				grid[i][j] = '|'
			} else {
				grid[i][j] = ' '
			}
		}
	}
	grid[vpitch][hpitch] = '+'

	plot := func(xf, yf float64, c byte) {
		v := int(yf/float64(yScale)+float64(vpitch)) + 0
		h := int(xf/float64(xScale)+float64(hpitch)) + 0
		if v < 0 || v > 2*vpitch || h < 0 || h > 2*hpitch {
			return
		}
		grid[2*vpitch-v][h] = c
	}

	for _, sp1 := range st.Ships {
		if sp1 == sp {
			continue
		}
		var xf, yf float64
		if game.CanSee(sp1) {
			sp1.LastKnown = game.LastKnownPosition{X: sp1.X, Y: sp1.Y, Warp: sp1.Warp, Course: sp1.Course}
			xf = float64(sp1.X - sp.X)
			yf = float64(sp1.Y - sp.Y)
		} else {
			xf = float64(sp1.LastKnown.X - sp.X)
			yf = float64(sp1.LastKnown.Y - sp.Y)
		}
		c := sp1.Name[0]
		if game.CantSee(sp1) {
			c = toLowerByte(c)
		}
		plot(xf, yf, c)
	}
	for _, obj := range st.Objects {
		xf := float64(obj.X - sp.X)
		yf := float64(obj.Y - sp.Y)
		var c byte
		switch obj.Type {
		case game.ObjectTorpedo:
			c = ':'
		case game.ObjectEngineering:
			c = '#'
		case game.ObjectProbe:
			c = '*'
		default:
			c = '?'
		}
		plot(xf, yf, c)
	}

	var lines []string
	for i := range grid {
		lines = append(lines, string(grid[i]))
	}
	return lines
}

func toLowerByte(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}
