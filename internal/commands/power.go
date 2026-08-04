package commands

import (
	"fmt"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

// AlterPower applies new shield and phaser drain settings, ported from
// alter_power() in cmds3.c. shieldDrains and phaserDrains must each
// have exactly len(sp.Shields)/len(sp.Phasers) entries (the caller is
// responsible for expanding the original's "*" "same value for the
// rest" shorthand into a fully-specified slice before calling this
// function, since parsing/prompting is not this function's job). If
// any value is out of range, nothing is applied and a single
// bad-parameters message is returned, matching the original's
// "goto badparam" behavior.
func AlterPower(st *game.State, sp *game.Ship, shieldDrains []float64, phaserDrains []float64) []string {
	if len(shieldDrains) != len(sp.Shields) {
		return []string{fmt.Sprintf("%s:  Bad parameters, %s.", st.Crew.Engineer, st.Crew.Title)}
	}
	for _, d := range shieldDrains {
		if d < 0.0 || d > 1.0 {
			return []string{fmt.Sprintf("%s:  Bad parameters, %s.", st.Crew.Engineer, st.Crew.Title)}
		}
	}
	if len(phaserDrains) != len(sp.Phasers) {
		return []string{fmt.Sprintf("%s:  Bad parameters, %s.", st.Crew.Engineer, st.Crew.Title)}
	}
	for _, d := range phaserDrains {
		if d < game.MinPhaserDrain || d > game.MaxPhaserDrain {
			return []string{fmt.Sprintf("%s:  Bad parameters, %s.", st.Crew.Engineer, st.Crew.Title)}
		}
	}
	for i, d := range shieldDrains {
		sp.Shields[i].AttemptDrain = d
	}
	for i, d := range phaserDrains {
		sp.Phasers[i].Drain = int(d)
	}
	return nil
}

// DoJettison detaches sp's engineering section into a new space
// object drifting at sp's current course/speed, disabling warp drive
// and firing arcs, ported from do_jettison() in cmds3.c.
func DoJettison(st *game.State, sp *game.Ship) *game.SpaceObject {
	obj := &game.SpaceObject{
		ID:        st.NextObjectID(),
		Type:      game.ObjectEngineering,
		From:      sp,
		X:         sp.X,
		Y:         sp.Y,
		Course:    sp.Course,
		Speed:     sp.Warp,
		NewSpeed:  0.0,
		Fuel:      int(sp.Energy),
		TimeDelay: 10.0,
	}
	sp.Energy = 0
	sp.Pods = 0
	if sp.Warp < 0.0 {
		sp.NewWarp = -0.99
	} else {
		sp.NewWarp = 0.99
	}
	sp.Regen = 0.0
	sp.Status[game.SysEngineering] = game.FullyDamaged
	sp.Status[game.SysWarp] = game.FullyDamaged
	sp.MaxSpeed = 1.0
	sp.Cloaking = game.CloakNone
	sp.TubeBlindLeft, sp.TubeBlindRight = 180, 180
	sp.PhaserBlindLeft, sp.PhaserBlindRight = 180, 180
	return obj
}

// JettisonEngineering jettisons sp's engineering section, ported from
// jettison_engineering() in cmds3.c. Returns nil if it was already
// jettisoned.
func JettisonEngineering(st *game.State, sp *game.Ship) (*game.SpaceObject, []string) {
	if sp.IsDead(game.SysEngineering) {
		return nil, []string{fmt.Sprintf("%s:  But Captain, it's already jettisonned.", st.Crew.Engineer)}
	}
	obj := DoJettison(st, sp)
	return obj, []string{fmt.Sprintf("%s:  Aye, %s.  Jettisoning engineering.", st.Crew.Engineer, st.Crew.Title)}
}

// DetonateEngineering detonates sp's already-jettisoned engineering
// section (or, if confirmDetonateAnyway is true, jettisons it first),
// ported from detonate_engineering() in cmds3.c. jettisoned is
// non-nil only if this call newly jettisoned the section.
func DetonateEngineering(st *game.State, sp *game.Ship, confirmDetonateAnyway bool) (jettisoned *game.SpaceObject, messages []string) {
	if !sp.IsDead(game.SysEngineering) {
		if !confirmDetonateAnyway {
			return nil, nil
		}
		jettisoned = DoJettison(st, sp)
	}
	for _, obj := range st.Objects {
		if obj.Type == game.ObjectEngineering && obj.From == sp {
			obj.TimeDelay = 0.0
			return jettisoned, []string{fmt.Sprintf("%s:  Aye, %s.", st.Crew.Engineer, st.Crew.Title)}
		}
	}
	if jettisoned != nil {
		jettisoned.TimeDelay = 0.0
		return jettisoned, []string{fmt.Sprintf("%s:  Aye, %s.", st.Crew.Engineer, st.Crew.Title)}
	}
	return nil, []string{fmt.Sprintf("%s:  Ours has already detonated.", st.Crew.Engineer)}
}

// AlterFiringParams resets tube and/or phaser firing parameters,
// ported from alterpntparams() in cmds4.c. Values outside their valid
// range are silently ignored (left unchanged), matching the original.
func AlterFiringParams(sp *game.Ship, resetTubes bool, launchSpeed, timeDelay, proxDelay int, resetPhasers bool, firePercent int) {
	if resetTubes {
		if launchSpeed >= 0 && launchSpeed < game.MaxTubeSpeed+1 {
			sp.TubeLaunchSpd = launchSpeed
		}
		if timeDelay >= 0 && float64(timeDelay) < game.MaxTubeTime+1 {
			sp.TubeDelay = timeDelay
		}
		if proxDelay >= 0 && proxDelay < game.MaxTubeProx+1 {
			sp.TubeProximity = proxDelay
		}
	}
	if resetPhasers {
		if firePercent >= 0 && firePercent <= 100 {
			sp.PhaserFirePct = firePercent
		}
	}
}
