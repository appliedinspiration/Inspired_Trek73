package commands

import (
	"fmt"
	"math"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

// clampWarp mirrors the repeated warp-validation block in
// pursue()/elude()/helm() (cmds2.c): a warp magnitude over 1 is
// rejected if the warp drive is dead (forced down to +-1), and any
// magnitude over the ship's current max speed is capped to it.
func clampWarp(st *game.State, sp *game.Ship, warp float64) (float64, []string) {
	var messages []string
	if math.Abs(warp) > 1.0 && sp.IsDead(game.SysWarp) {
		messages = append(messages, fmt.Sprintf("%s: Warp drive is dead, Captain.", st.Crew.Science))
		if warp < 0.0 {
			warp = -1.0
		} else {
			warp = 1.0
		}
	}
	if math.Abs(warp) > sp.MaxSpeed {
		messages = append(messages, fmt.Sprintf("%s: %s, the engines canna go that fast!", st.Crew.Engineer, st.Crew.Title))
		if warp < 0.0 {
			warp = -sp.MaxSpeed
		} else {
			warp = sp.MaxSpeed
		}
	}
	return warp, messages
}

// Pursue sets sp on an intercept course with target at the given warp,
// ported from pursue() in cmds2.c.
func Pursue(st *game.State, r *game.Rand, sp *game.Ship, target *game.Ship, warp float64) []string {
	if sp.IsDead(game.SysComputer) {
		return []string{fmt.Sprintf("%s: Impossible, %s, our computer is dead", st.Crew.Science, st.Crew.Title)}
	}
	if !sp.SystemWorks(r, game.SysComputer) {
		return []string{fmt.Sprintf("%s: Main computer down, %s.  Rebooting.", st.Crew.Science, st.Crew.Title)}
	}
	if target == nil {
		return nil
	}
	if game.CantSee(target) {
		return []string{fmt.Sprintf("%s:  %s, unable to acquire helm lock.", st.Crew.Nav, st.Crew.Title)}
	}
	warp, messages := clampWarp(st, sp, warp)
	sp.NewWarp = warp
	sp.Target = target
	sp.RelativeBear = 0.0
	bear := game.BearingTo(sp.X, target.X, sp.Y, target.Y)
	messages = append(messages, fmt.Sprintf("%s: Aye, %s, coming to course %3.0f.", st.Crew.Nav, st.Crew.Title, bear))
	sp.NewCourse = bear
	return messages
}

// Elude sets sp on a course directly away from target at the given
// warp, ported from elude() in cmds2.c.
func Elude(st *game.State, r *game.Rand, sp *game.Ship, target *game.Ship, warp float64) []string {
	if sp.IsDead(game.SysComputer) {
		return []string{fmt.Sprintf("%s: Impossible, %s, our computer is dead", st.Crew.Science, st.Crew.Title)}
	}
	if !sp.SystemWorks(r, game.SysComputer) {
		return []string{fmt.Sprintf("%s: Main computer down, %s.  Rebooting.", st.Crew.Science, st.Crew.Title)}
	}
	if target == nil {
		return nil
	}
	if game.CantSee(target) {
		return []string{fmt.Sprintf("%s:  %s, unable to acquire helm lock.", st.Crew.Nav, st.Crew.Title)}
	}
	warp, messages := clampWarp(st, sp, warp)
	sp.NewWarp = warp
	sp.Target = target
	sp.RelativeBear = 180.0
	bear := game.Rectify(game.BearingTo(sp.X, target.X, sp.Y, target.Y) + 180.0)
	messages = append(messages, fmt.Sprintf("%s: Aye, %s, coming to course %3.0f.", st.Crew.Nav, st.Crew.Title, bear))
	sp.NewCourse = bear
	return messages
}

// Helm sets sp on a manual course and warp, clearing any pursue/elude
// lock, ported from helm() in cmds2.c.
func Helm(st *game.State, sp *game.Ship, course, warp float64) []string {
	if course < 0.0 || course >= 360.0 {
		return []string{fmt.Sprintf("%s: How's that, %s?", st.Crew.Nav, st.Crew.Title)}
	}
	warp, messages := clampWarp(st, sp, warp)
	sp.NewWarp = warp
	sp.NewCourse = course
	sp.Target = nil
	sp.RelativeBear = 0.0
	messages = append(messages,
		fmt.Sprintf("%s: Aye, %s.", st.Crew.Nav, st.Crew.Title),
		fmt.Sprintf("%s: Aye, %s.", st.Crew.Helmsman, st.Crew.Title),
	)
	return messages
}
