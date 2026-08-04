package game

// Leftovers reports whether any non-ship objects (unexploded
// torpedoes, antimatter probes, or jettisoned engineering sections)
// remain in flight, ported from leftovers() in endgame.c. While any
// remain, an otherwise-resolved battle must show a Warn message
// instead of immediately ending with a Final one, since the original
// let those objects keep threatening ships even after the battle's
// outcome was decided.
func Leftovers(st *State) bool {
	return len(st.Objects) > 0
}

// resetDefeatedShip marks sp destroyed and clears its combat state,
// ported from the "complement <= 0" cleanup block at the top of each
// disposition() loop iteration (moveships.c). A ship whose entire crew
// has died is instantly and fully disabled, even if its equipment is
// nominally still intact.
func resetDefeatedShip(sp *Ship) {
	if sp.Complement > 0 || sp.IsDead(SysDead) {
		return
	}
	sp.Status[SysDead] = FullyDamaged
	sp.NewWarp = 0.0
	sp.NewCourse = sp.Course
	sp.Target = nil
	sp.RelativeBear = 0.0
	for i := range sp.Phasers {
		sp.Phasers[i].Target = nil
		sp.Phasers[i].Drain = MinPhaserDrain
		sp.Phasers[i].Status &^= PhaserFiring
	}
	for i := range sp.Tubes {
		sp.Tubes[i].Target = nil
		sp.Tubes[i].Status &^= TubeFiring
	}
	for i := range sp.Shields {
		sp.Shields[i].AttemptDrain = 0.0
	}
	sp.Regen = 0.0
	sp.Cloaking = CloakNone
}

// shipDisarmed reports whether every phaser bank and torpedo tube on
// sp is damaged, ported from the "all phasers damaged, then all tubes
// damaged" checks used by disposition() to treat a ship as
// effectively defeated even while its crew and hull technically
// survive.
func shipDisarmed(sp *Ship) bool {
	for i := range sp.Phasers {
		if sp.Phasers[i].Status&PhaserDamaged == 0 {
			return false
		}
	}
	for i := range sp.Tubes {
		if sp.Tubes[i].Status&TubeDamaged == 0 {
			return false
		}
	}
	return true
}

// Disposition determines whether the battle has been decided (one
// side defeated, disengaged, or a surrender accepted), ported from
// disposition() in moveships.c. Any resolution is recorded on
// st.PendingWarn (if Leftovers(st) means the game must first wait for
// remaining ordnance to clear) or st.PendingFinal (if the game is over
// immediately); the CLI layer is responsible for presenting the
// corresponding Warn/Final message and, for PendingFinal, ending the
// session.
//
// The original called warn()/final() directly, with final() calling
// exit(1) to end the process immediately; this port instead just
// records the resolved condition and returns, mirroring exit()'s
// "nothing after this runs" effect without an actual process exit.
func Disposition(st *State) []string {
	var messages []string
	fed := st.Player()

	resetDefeatedShip(fed)
	fedstatus := fed.IsDead(SysDead) || shipDisarmed(fed)

	enemies := st.Enemies()
	shipnum := len(enemies)
	kills, others := 0, 0
	for _, ep := range enemies {
		resetDefeatedShip(ep)
		if ep.IsDead(SysDead) {
			kills++
			continue
		}
		if fedstatus {
			continue
		}
		j := RangeFind(fed.X, ep.X, fed.Y, ep.Y)
		if (j > 3500 && ep.IsDead(SysWarp)) || (j > 4500 && ep.Delay < 10.0) {
			others++
			continue
		}
		if j <= 3500 && st.Reengaged {
			st.Reengaged = false
		}
		if ep.Energy > 10 {
			continue
		}
		if !shipDisarmed(ep) {
			continue
		}
		kills++
	}

	// settle records mesg as a Warn (if ordnance is still in flight)
	// or a Final (otherwise) when cond holds, and reports whether a
	// Final was recorded (meaning, like the original's exit(1), no
	// further disposition checks should run).
	settle := func(cond bool, mesg int) bool {
		if !cond {
			return false
		}
		if Leftovers(st) {
			st.PendingWarn = mesg
			return false
		}
		st.PendingFinal = mesg
		return true
	}

	if settle(!fedstatus && st.PlayStatus&StatusEnemySurrender != 0, FinEnemySurrender) {
		return messages
	}
	if settle(fed.IsDead(SysSurrender) && kills+others < shipnum, FinFedSurrender) {
		return messages
	}
	if !fedstatus && kills+others < shipnum {
		return messages // Play continues.
	}
	if fedstatus && kills < shipnum {
		settle(true, FinFedLose)
		return messages
	}
	if !fedstatus && kills == shipnum {
		settle(true, FinEnemyLose)
		return messages
	}
	if !fedstatus && kills+others == shipnum {
		settle(true, FinTactical)
		return messages
	}
	if fedstatus && kills == shipnum {
		// Both sides destroyed: always immediately final, unlike the
		// other conditions above, matching the original (no leftovers
		// check for FIN_COMPLETE).
		st.PendingFinal = FinComplete
	}
	return messages
}
