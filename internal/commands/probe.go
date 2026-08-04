package commands

import (
	"fmt"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

// LaunchProbe validates and launches a new antimatter probe, ported
// from launch_probe() in cmds1.c. Exactly one of target/course should
// be meaningful: if target is non-nil the probe homes in on it,
// otherwise it flies the given course. Returns the newly created probe
// (nil on failure, matching the original's early "goto bad_param"/
// return 0 paths) and any messages.
func LaunchProbe(st *game.State, r *game.Rand, sp *game.Ship, pods, delay, prox int, target *game.Ship, course float64) (*game.SpaceObject, []string) {
	if sp.IsDead(game.SysProbe) {
		return nil, []string{fmt.Sprintf("%s:  Probe launcher destroyed!", st.Crew.Engineer)}
	}
	if !sp.SystemWorks(r, game.SysProbe) {
		return nil, []string{fmt.Sprintf("%s:  Probe launcher temporarily disabled, %s", st.Crew.Engineer, st.Crew.Title)}
	}
	if sp.Energy < game.MinProbeCharge {
		return nil, []string{fmt.Sprintf("%s: We've not enough power, Captain.", st.Crew.Engineer)}
	}
	if pods < game.MinProbeCharge || float64(pods) > sp.Energy {
		return nil, []string{fmt.Sprintf("%s: Bad parameters, %s.", st.Crew.Science, st.Crew.Title)}
	}
	if delay < 0 || delay > game.MaxProbeDelay {
		return nil, []string{fmt.Sprintf("%s: Bad parameters, %s.", st.Crew.Science, st.Crew.Title)}
	}
	if prox < game.MinProbeProx {
		return nil, []string{fmt.Sprintf("%s: Bad parameters, %s.", st.Crew.Science, st.Crew.Title)}
	}
	if target != nil {
		if game.CantSee(target) || !sp.SystemWorks(r, game.SysSensor) {
			return nil, []string{fmt.Sprintf("%s:  %s, unable to lock probe onto the %s.", st.Crew.Helmsman, st.Crew.Title, target.Name)}
		}
	}

	obj := &game.SpaceObject{
		ID:        st.NextObjectID(),
		Type:      game.ObjectProbe,
		From:      sp,
		X:         sp.X,
		Y:         sp.Y,
		Course:    course,
		Speed:     sp.Warp,
		NewSpeed:  3.0,
		Target:    target,
		Fuel:      pods,
		TimeDelay: float64(delay),
		Proximity: prox,
	}
	sp.Pods -= float64(pods)
	sp.Energy -= float64(pods)
	sp.ProbeLauncherStatus = game.ProbeLaunching
	return obj, []string{fmt.Sprintf("%s: Probe %d away", st.Crew.Engineer, obj.ID)}
}

// ShipProbes returns every probe currently in flight that was launched
// by sp, ported from the pnum-collecting loop at the top of
// probe_control() in cmds1.c.
func ShipProbes(st *game.State, sp *game.Ship) []*game.SpaceObject {
	var probes []*game.SpaceObject
	for _, obj := range st.Objects {
		if obj.Type == game.ObjectProbe && obj.From == sp {
			probes = append(probes, obj)
		}
	}
	return probes
}

// DetonateAllProbes immediately detonates every one of sp's probes,
// ported from the "Detonate all probes?" branch of probe_control() in
// cmds1.c.
func DetonateAllProbes(sp *game.Ship, probes []*game.SpaceObject) []string {
	for _, p := range probes {
		p.TimeDelay = 0.0
	}
	sp.ProbeLauncherStatus = game.ProbeDetonate
	return []string{"Aye."}
}

// DetonateProbe immediately detonates a single probe, ported from the
// "Detonate it?" branch of probe_control() in cmds1.c.
func DetonateProbe(sp *game.Ship, probe *game.SpaceObject) []string {
	probe.TimeDelay = 0.0
	sp.ProbeLauncherStatus = game.ProbeDetonate
	return nil
}

// LockProbe re-locks a probe onto a new target, ported from the "lock
// it onto [whom]" branch of probe_control() in cmds1.c.
func LockProbe(st *game.State, r *game.Rand, sp *game.Ship, probe *game.SpaceObject, target *game.Ship) []string {
	sp.ProbeLauncherStatus = game.ProbeLocked
	if target == nil {
		return nil
	}
	if game.CanSee(target) && sp.SystemWorks(r, game.SysSensor) {
		probe.Target = target
		return []string{fmt.Sprintf("%s: locking.", st.Crew.Nav)}
	}
	return []string{fmt.Sprintf("%s:  %s, unable to lock probe on the %s.", st.Crew.Helmsman, st.Crew.Title, target.Name)}
}

// CourseProbe sets a probe to fly an unguided course, ported from the
// "set it to course" branch of probe_control() in cmds1.c.
func CourseProbe(st *game.State, sp *game.Ship, probe *game.SpaceObject, course float64) []string {
	sp.ProbeLauncherStatus = game.ProbeLocked
	if course < 0.0 || course > 360.0 {
		return nil
	}
	probe.Course = course
	probe.Target = nil
	return []string{fmt.Sprintf("%s: setting in new course.", st.Crew.Nav)}
}
