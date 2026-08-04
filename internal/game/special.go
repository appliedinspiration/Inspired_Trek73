package game

import (
	"fmt"
	"math"
)

// Special runs one turn's worth of special-action bookkeeping for the
// enemy ship sp against the player's ship fed at the given range,
// ported from special() in special.c. It resolves in-progress
// "playing defenseless", "corbomite bluff", "we surrender", and
// "you surrender" ruses one state-transition at a time (mirroring the
// original's switch-fallthrough state machine, translated using Go's
// explicit fallthrough keyword), and checks for "unsportsmanlike
// firing" - the player firing on the enemy while one of those ruses
// is in progress, which immediately cancels every ruse in progress.
//
// The original called warn()/final() (not yet ported - they belong to
// the end-of-game presentation layer) to announce and act on a
// surrender being accepted; this port instead records the pending
// condition on st.PendingWarn/st.PendingFinal for a future phase to
// notice and act on.
func (h *CombatHooks) Special(sp *Ship, rng int, fed *Ship) {
	st := h.State
	r := h.Rand

	// Play-dead ("defenseless") ruse effects.
	switch st.Defenseless {
	case 1:
		defenselessChance := st.DefenselessChance(r)
		if r.Randm(100) > defenselessChance {
			// Didn't work. Too bad.
			if canSee(sp) {
				h.emit(fmt.Sprintf("%s:   No apparent change in the enemy's actions.", st.Crew.Helmsman))
			}
			st.Defenseless = 6
			break
		}
		st.Defenseless = 2
		fallthrough
	case 2, 3:
		// Okay, he's fallen for it. Choose his action.
		if r.Randm(2) == 1 {
			sp.Target = nil
			sp.NewWarp = 0.0
		} else {
			sp.NewWarp = 1.0
		}
		if canSee(sp) {
			if sp.Target != nil {
				h.emit(fmt.Sprintf("%s:   The %s is cautiously advancing.", st.Crew.Helmsman, sp.Name))
			} else {
				h.emit(fmt.Sprintf("%s:   The %s is turning away.", st.Crew.Helmsman, sp.Name))
			}
		}
		fallthrough
	case 4, 5:
		// Now he might get suspicious. If he's moving too fast, or if
		// we're close enough, or if his shields are up, we'll be
		// spotted.
		if sp.Target != nil && (math.Abs(sp.Target.Warp) > 1.0 || rng < 200) {
			st.Defenseless = 6
		} else if sp.Target != nil {
			for i := 0; i < NumShields; i++ {
				if sp.Target.Shields[i].Drain != 0 {
					st.Defenseless = 6
				}
			}
		}
	}

	// Corbomite bluff effects.
	switch st.Corbomite {
	case 1:
		corbomiteChance := st.CorbomiteChance(r)
		if r.Randm(100) > corbomiteChance {
			// He didn't fall for it.
			h.emit(fmt.Sprintf("%s:  Message coming in from the %ss.", st.Crew.Com, st.EnemyRaceName))
			h.emit(fmt.Sprintf("%s:  Put it on audio.", st.Crew.Captain))
			if r.Randm(2) == 1 {
				h.emit(fmt.Sprintf("%s:  Ha, ha, ha, %s.  You lose.", st.EnemyCommander, st.Crew.Captain))
			} else {
				h.emit(fmt.Sprintf("%s:  I fell for that the last time we met, idiot!", st.EnemyCommander))
			}
			st.Corbomite = 6
			break
		}
		if canSee(sp) {
			h.emit(fmt.Sprintf("%s:   %ss giving ground, Captain.  Obviously they", st.Crew.Science, st.EnemyRaceName))
			h.emit("   tapped in as you expected them to.")
			h.emit(fmt.Sprintf("%s:  A logical assumption, Mr. %s.  Are they still", st.Crew.Captain, st.Crew.Science))
			h.emit("   retreating?")
			h.emit(fmt.Sprintf("%s:  Yes, %s", st.Crew.Science, st.Crew.Title))
			h.emit(fmt.Sprintf("%s:  Good.  All hands, stand by.", st.Crew.Captain))
		}
		st.Corbomite = 2
		fallthrough
	case 2:
		// He fell for it; retrograde out of here!
		sp.Target = nil
		sp.NewWarp = -(3.0 + float64(r.Randm(7)))
	case 3, 4, 5:
		// Begin to get suspicious.
		if sp.Target != nil && math.Abs(sp.Target.Warp) > 2.0 {
			st.Corbomite = 6
		}
	}

	// Will the enemy accept your surrender?
	switch st.Surrender {
	case 1:
		surrenderChance := st.SurrenderChance(r)
		// Just a little reminder.
		if surrenderChance <= 10 {
			h.emit(fmt.Sprintf("%s:  The %ss do not take prisoners.", st.Crew.Nav, st.EnemyRaceName))
		}
		if r.Randm(100) > surrenderChance {
			// Tough luck.
			if r.Randm(2) == 1 {
				h.emit(fmt.Sprintf("%s:  Message coming in from the %ss.", st.Crew.Com, st.EnemyRaceName))
				h.emit(fmt.Sprintf("%s:  Put it on audio.", st.Crew.Captain))
				h.emit(fmt.Sprintf("%s:  Prepare to die, Chicken %s!", st.EnemyCommander, st.Crew.Captain))
			} else {
				h.emit(fmt.Sprintf("%s:  No reply from the %ss", st.Crew.Com, st.EnemyRaceName))
			}
			st.Surrender = 6
			break
		}
		// He took it!
		h.emit(fmt.Sprintf("%s:  Message coming in from the %ss.", st.Crew.Com, st.EnemyRaceName))
		h.emit(fmt.Sprintf("%s:  Put it on audio.", st.Crew.Captain))
		h.emit(fmt.Sprintf("%s:  On behalf of the %s %s, I accept your surrender.", st.EnemyCommander, st.EnemyRaceName, st.EnemyEmpireName))
		h.emit("   You have five seconds to drop your shields, cut")
		h.emit("   warp, and prepare to be boarded.")
		st.PlayStatus |= StatusFedSurrender
		fallthrough
	case 2, 3:
		if st.Surrender == 1 {
			st.Surrender = 2
		} else {
			st.PendingWarn = FinFedSurrender
		}
		sp.Target = fed
		sp.NewWarp = sp.MaxSpeed
		h.ECloakOff(sp, fed)
	case 4, 5:
		// Begin checking surrender conditions.
		shieldsUp := false
		for i := 0; i < NumShields; i++ {
			if sp.Target != nil && sp.Target.Shields[i].Drain != 0 {
				shieldsUp = true
				break
			}
		}
		if shieldsUp {
			break
		}
		if rng <= 1400 {
			sp.NewWarp = 1.0
		}
		if rng <= 1000 && sp.Target != nil && math.Abs(sp.Target.Warp) <= 1.0 {
			fed.Status[SysSurrender] = 100
			st.PendingFinal = FinFedSurrender
		}
		if st.Surrender == 4 {
			break
		}
		// The original had a bug here: due to a missing pair of
		// braces, only the first printf below was actually
		// conditional on !shutup[SURRENDER] - the "resuming our
		// attack" message and the surrender=6 reset ran
		// unconditionally every time this code was reached, ending
		// the ruse on its very first re-check instead of only once.
		// This port wraps all three statements in the guard, as the
		// original indentation clearly intended.
		if !st.Shutup.Surrender {
			h.emit(fmt.Sprintf("%s:  Captain %s, you have not fulfilled our terms.", st.EnemyCommander, st.Crew.Captain))
			h.emit("  We are resuming our attack.")
			st.Surrender = 6
		}
		st.Shutup.Surrender = true
		fallthrough
	default:
		st.PlayStatus &^= StatusFedSurrender
	}

	// Enemy surrenders?
	switch st.SurrenderP {
	case 1:
		outmatched := false
		for _, other := range st.Enemies() {
			if !other.IsDead(SysEngineering) && sp.Complement > 100 {
				h.emit(fmt.Sprintf("%s:  Message coming in from the %ss.", st.Crew.Com, st.EnemyRaceName))
				h.emit(fmt.Sprintf("%s:  Put it on audio.", st.Crew.Captain))
				h.emit(fmt.Sprintf("%s:  You must be joking, Captain %s.", st.EnemyCommander, st.Crew.Captain))
				h.emit("  Why don't you surrender?")
				st.SurrenderP = 6
				outmatched = true
				break
			}
		}
		if outmatched {
			break
		}
		surrenderPChance := st.SurrenderPChance(r)
		if r.Randm(100) > surrenderPChance {
			h.emit(fmt.Sprintf("%s:  I'll never surrender to you, %s", st.EnemyCommander, st.Crew.Captain))
			st.SurrenderP = 6
			break
		}
		h.emit(fmt.Sprintf("%s:  As much as I hate to, Captain %s, we will surrender.", st.EnemyCommander, st.Crew.Captain))
		h.emit("   We are dropping shields.  You may board us.")
		fallthrough
	case 2, 3:
		if st.SurrenderP == 1 {
			st.SurrenderP = 2
		} else {
			st.PendingWarn = FinEnemySurrender
		}
		for i := range sp.Shields {
			sp.Shields[i].AttemptDrain = 0.0
		}
		sp.NewWarp = 0.0
		for _, s := range st.Enemies() {
			s.Status[SysSurrender] = 100
		}
		st.PlayStatus |= StatusEnemySurrender
	}

	// Unsportsmanlike firing: if the player fires while a ruse is in
	// progress, every ruse in progress is immediately cancelled.
	if betw(float64(st.Defenseless), 0, 6) || betw(float64(st.Corbomite), 0, 6) ||
		betw(float64(st.Surrender), 0, 6) || betw(float64(st.SurrenderP), 0, 6) {
		phaserFired := false
		for i := range fed.Phasers {
			if fed.Phasers[i].Status&PhaserFiring != 0 {
				phaserFired = true
				break
			}
		}
		tubeFired := false
		for i := range fed.Tubes {
			if fed.Tubes[i].Status&TubeFiring != 0 {
				tubeFired = true
				break
			}
		}
		probeFired := fed.ProbeLauncherStatus != ProbeNormal
		if phaserFired || tubeFired || probeFired {
			h.emit(fmt.Sprintf("%s: How dare you fire on us!  We are resuming our attack!", st.EnemyCommander))
			st.PlayStatus = StatusNormal
			if betw(float64(st.Defenseless), 0, 6) {
				st.Defenseless = 6
			}
			if betw(float64(st.Corbomite), 0, 6) {
				st.Corbomite = 6
			}
			if betw(float64(st.Surrender), 0, 6) {
				st.Surrender = 6
			}
			if betw(float64(st.SurrenderP), 0, 6) {
				st.SurrenderP = 6
			}
			for _, s := range st.Ships {
				s.Status[SysSurrender] = 0
			}
		}
	}
}
