package commands

import (
	"fmt"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

// PowerTransferTarget selects where play_dead's power reallocation
// goes, mirroring the buf1[0]=='e'/'p' checks in play_dead().
type PowerTransferTarget int

const (
	// TransferNone represents an unrecognized/empty answer.
	TransferNone PowerTransferTarget = iota
	TransferEngines
	TransferPhasers
)

// PlayDead drops sp's shields and reroutes their power into either the
// engines or the phasers, ported from play_dead() in cmds4.c.
func PlayDead(st *game.State, sp *game.Ship, target PowerTransferTarget) []string {
	msgs := []string{fmt.Sprintf("%s:   %s, drop shields ...", st.Crew.Captain, st.Crew.Nav)}
	if st.Defenseless != 0 {
		msgs = append(msgs, fmt.Sprintf("%s:   %s, the %ss are not that stupid.", st.Crew.Science, st.Crew.Title, st.EnemyRaceName))
		return msgs
	}
	var phaserCharge float64
	switch target {
	case TransferEngines:
		phaserCharge = -game.MaxPhaserCharge
		sp.Energy = minFloat(sp.Energy+phaserCharge, sp.Pods)
	case TransferPhasers:
		phaserCharge = game.MaxPhaserCharge
	default:
		msgs = append(msgs, fmt.Sprintf("%s:   I cannot transfer power there, %s.", st.Crew.Nav, st.Crew.Title))
		return msgs
	}
	for i := range sp.Shields {
		sp.Shields[i].AttemptDrain = 0.0
	}
	for i := range sp.Phasers {
		sp.Phasers[i].Drain = int(phaserCharge)
	}
	st.Defenseless = 1
	return msgs
}

// CorbomiteBluff attempts the classic corbomite-device bluff, ported
// from corbomite_bluff() in cmds4.c. The coin flip (which of the two
// scripted transmissions plays) is passed in via headsFirstBranch so
// callers control determinism (originally randm(2) == 1).
func CorbomiteBluff(st *game.State, sp *game.Ship, headsFirstBranch bool) []string {
	var msgs []string
	if headsFirstBranch {
		msgs = append(msgs,
			fmt.Sprintf("%s:   Open a hailing frequency, ship-to-ship.", st.Crew.Captain),
			fmt.Sprintf("%s:  Hailing frequency open, %s.", st.Crew.Com, st.Crew.Title),
			fmt.Sprintf("%s:  This is the Captain of the %s.  Our respect for", st.Crew.Captain, sp.Name),
			"   other life forms requires that we give you this warning--",
			"   one critical item of information which has never been",
			"   incorporated into the memory banks of any Earth ship.",
			"   Since the early years of space exploration, Earth vessels",
			"   have had incorporated into them a substance know as corbomite.",
		)
		if st.Corbomite == 0 {
			msgs = append(msgs,
				"      It is a material and a device which prevents attack on",
				"   us.  If any destructive energy touchs our vessel, a re-",
				"   verse reaction of equal strength is created, destroying",
				"   the attacker.  It may interest you to know that, since",
				"   the initial use of corbomite for more than two of our",
				"   centuries ago, no attacking vessel has survived the attempt.",
				"   Death has little meaning to us.  If it has none to you,",
				"   then attack us now.  We grow annoyed with your foolishness.",
			)
		}
	} else {
		msgs = append(msgs,
			fmt.Sprintf("%s:   Open a special channel to Starfleet Command.", st.Crew.Captain),
			fmt.Sprintf("%s:   Aye, %s.", st.Crew.Com, st.Crew.Title),
			fmt.Sprintf("%s:   Use Code 2.", st.Crew.Captain),
			fmt.Sprintf("%s:   but, %s, according to our last Starfleet", st.Crew.Com, st.Crew.Title),
			fmt.Sprintf("   Bulletin, the %ss have broken code 2.", st.EnemyRaceName),
			fmt.Sprintf("%s:   That's an order, Lieutenant.  Code 2!", st.Crew.Captain),
			fmt.Sprintf("%s:   Yes, Captain.  Code 2.", st.Crew.Com),
			fmt.Sprintf("%s:   Message from %s to Starfleet Command, this sector.", st.Crew.Captain, sp.Name),
			fmt.Sprintf("   have inadvertantly encroached upon %s neutral zone,", st.EnemyRaceName),
			fmt.Sprintf("   surrounded and under heavy %s attack.  Escape", st.EnemyRaceName),
			"   impossible.  Shields failing.  Will implement destruct",
			"   order using corbomite device recently installed.",
		)
		if st.Corbomite == 0 {
			msgs = append(msgs,
				fmt.Sprintf("   This will result in the destruction of the %s and", sp.Name),
				"   all matter within a 200 megameter diameter and",
				"   establish corresponding dead zone, all Federation",
				"   vessels will aviod this area for the next four solar",
				fmt.Sprintf("   years.  Explosion will take place in one minute.  %s,", st.Crew.Captain),
				fmt.Sprintf("   commanding %s, out.", sp.Name),
			)
		}
	}
	if st.Corbomite == 0 {
		msgs = append(msgs,
			fmt.Sprintf("      Mr. %s.  Stand by.", st.Crew.Helmsman),
			fmt.Sprintf("%s:  Standing by.", st.Crew.Helmsman),
		)
		st.Corbomite = 1
	} else {
		msgs = append(msgs,
			fmt.Sprintf("%s:  I don't believe that they will fall for that maneuver", st.Crew.Science),
			fmt.Sprintf("   again, %s.", st.Crew.Title),
		)
	}
	return msgs
}

// SurrenderShip offers sp's unconditional surrender to the enemy,
// ported from surrender_ship() in cmds4.c.
//
// The original compared foerace == "Romulan" as a pointer comparison
// rather than a string comparison, so that branch could never fire;
// this port fixes that bug and uses a proper string comparison.
func SurrenderShip(st *game.State, sp *game.Ship) []string {
	msgs := []string{
		fmt.Sprintf("%s:   %s, open a channel to the %ss.", st.Crew.Captain, st.Crew.Com, st.EnemyRaceName),
		fmt.Sprintf("%s:   Aye, %s.", st.Crew.Com, st.Crew.Title),
		fmt.Sprintf("%s:   This is Captain %s of the U.S.S. %s.  Will", st.Crew.Captain, st.Crew.Captain, sp.Name),
		"   you accept my unconditional surrender?",
	}
	if st.PlayStatus&game.StatusFedSurrender != 0 {
		msgs = append(msgs, fmt.Sprintf("%s:  %s, we have already surrendered.", st.Crew.Science, st.Crew.Title))
		return msgs
	}
	if st.Surrender != 0 {
		msgs = append(msgs, fmt.Sprintf("%s:  The %ss have already refused.", st.Crew.Science, st.EnemyRaceName))
	} else {
		if st.EnemyRaceName == "Romulan" {
			msgs = append(msgs,
				fmt.Sprintf("%s:  The %ss have not been know to have taken", st.Crew.Science, st.EnemyRaceName),
				"   prisoners.",
			)
		}
		st.Surrender = 1
	}
	return msgs
}

// RequestSurrender demands the enemy's surrender, ported from
// request_surrender() in cmds4.c.
func RequestSurrender(st *game.State, sp *game.Ship) []string {
	msgs := []string{
		fmt.Sprintf("%s:   %s, open a hailing frequency to the %ss.", st.Crew.Com, st.Crew.Captain, st.EnemyRaceName),
		fmt.Sprintf("%s:  Aye, %s.", st.Crew.Com, st.Crew.Title),
		fmt.Sprintf("%s:  This is Captain %s of the U. S. S. %s.  I give you", st.Crew.Captain, st.Crew.Captain, sp.Name),
		"   one last chance to surrender before we resume our attack.",
	}
	if st.PlayStatus&game.StatusEnemySurrender != 0 {
		msgs = append(msgs, fmt.Sprintf("%s:  %s, we are already complying with your previous request!", st.EnemyCommander, st.Crew.Captain))
		return msgs
	}
	if st.SurrenderP != 0 {
		msgs = append(msgs, fmt.Sprintf("%s:   %s, our offer has already been refused.", st.Crew.Science, st.Crew.Title))
	} else {
		st.SurrenderP = 1
	}
	return msgs
}

// SelfDestructDelay is the countdown, in seconds, set on sp.Delay once
// the self-destruct sequence has been fully engaged, ported from the
// "sp->delay = 22." assignment in self_destruct().
const SelfDestructDelay = 22.0

// SelfDestruct engages sp's self-destruct sequence, ported from
// self_destruct() in cmds4.c.
func SelfDestruct(st *game.State, r *game.Rand, sp *game.Ship) []string {
	msgs := []string{
		fmt.Sprintf("%s:   Lieutenant %s, tie in the bridge to the master", st.Crew.Captain, st.Crew.Com),
		"   computer.",
	}
	if sp.IsDead(game.SysComputer) {
		msgs = append(msgs, fmt.Sprintf("%s:  Our computer is down.", st.Crew.Science))
		return msgs
	}
	if !sp.SystemWorks(r, game.SysComputer) {
		msgs = append(msgs, fmt.Sprintf("%s:  That program has been lost.  Restoring from backup.", st.Crew.Science))
		return msgs
	}
	msgs = append(msgs,
		fmt.Sprintf("%s:  Aye, %s.", st.Crew.Com, st.Crew.Title),
		fmt.Sprintf("%s:  Computer.  Destruct sequence.  Are you ready to copy?", st.Crew.Captain),
		"Computer:  Working.",
		fmt.Sprintf("%s:  Computer, this is Captain %s of the U.S.S. %s.", st.Crew.Captain, st.Crew.Captain, sp.Name),
		"   Destruct sequence one, code 1-1A.",
		"Computer:  Voice and code verified and correct.",
		"   Sequence one complete.",
		fmt.Sprintf("%s:  This is Commander %s, Science Officer.  Destruct", st.Crew.Science, st.Crew.Science),
		"   sequence two, code 1-1A-2B.",
		"Computer:  Voice and code verified and correct.  Sequence",
		"   two complete.",
		fmt.Sprintf("%s:  This is Lieutenant Commander %s, Chief Engineering", st.Crew.Engineer, st.Crew.Engineer),
		fmt.Sprintf("   Officer of the U. S. S. %s.  Destruct sequence", sp.Name),
		"   number three, code 1B-2B-3.",
		"Computer:  Voice and code verified and correct.",
		"   Destruct sequence complete and engaged.  Awaiting final",
		"   code for twenty second countdown.",
		fmt.Sprintf("%s:  Computer, this is Captain %s of the U. S. S. %s.", st.Crew.Captain, st.Crew.Captain, sp.Name),
		"   begin countdown, code 0-0-0, destruct 0.",
	)
	sp.Delay = SelfDestructDelay
	return msgs
}

// AbortSelfDestruct attempts to abort sp's in-progress self-destruct
// sequence, ported from abort_self_destruct() in cmds4.c. The original
// slept 4 real seconds mid-command to build suspense before checking
// whether the countdown had already reached zero; that delay is a pure
// presentation concern and is left to the (not yet built) CLI layer to
// simulate if desired, rather than being modeled here.
func AbortSelfDestruct(st *game.State, r *game.Rand, sp *game.Ship) []string {
	msgs := []string{
		fmt.Sprintf("%s:   Computer, this is Captain %s of the U.S.S. %s.", st.Crew.Captain, st.Crew.Captain, sp.Name),
		"   Code 1-2-3 continuity abort destruct order, repeat:",
		"   Code 1-2-3 continuity abort destruct order!",
	}
	if sp.IsDead(game.SysComputer) {
		msgs = append(msgs, fmt.Sprintf("%s:  Our computer is down.", st.Crew.Science))
		return msgs
	}
	if !sp.SystemWorks(r, game.SysComputer) {
		msgs = append(msgs, fmt.Sprintf("%s:  Temporary memory loss.  Unable to find program.", st.Crew.Science))
		return msgs
	}
	if sp.Delay > 1000.0 {
		msgs = append(msgs,
			"Computer:  Self-destruct sequence has not been",
			"   initiated.",
		)
		return msgs
	}
	msgs = append(msgs, "Computer:   Self-destruct order ... ")
	if sp.Delay > 4.0 {
		msgs = append(msgs, "aborted.  Destruct order aborted.")
		sp.Delay = 10000.0
	} else {
		msgs = append(msgs, "cannot be aborted.")
	}
	return msgs
}
