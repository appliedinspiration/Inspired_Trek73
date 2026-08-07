package commands

import (
	"fmt"

	"github.com/appliedinspiration/inspired_trek73/internal/data"
	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

// plural returns "s" if count is greater than 1, and "" otherwise,
// ported from the plural(n) macro in defines.h.
func plural(count int) string {
	if count > 1 {
		return "s"
	}
	return ""
}

// vowelStr returns "n" if s begins with a vowel, and "" otherwise,
// ported from vowelstr() in misc.c (used to pick "a"/"an").
func vowelStr(s string) string {
	if len(s) == 0 {
		return ""
	}
	switch s[0] {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return "n"
	}
	return ""
}

// Starfleet returns the "Starfleet Command:" header shown before every
// Final message, ported from starfleet() in endgame.c. The original's
// sleep(3) suspense pause is a pure presentation concern left to the
// CLI layer to simulate if desired, rather than modeled here.
func Starfleet() []string {
	return []string{"", "", "Starfleet Command: ", ""}
}

// Final renders the end-of-game message for the given Fin* condition,
// ported from final() in endgame.c. Unlike the original (which called
// exit(1) itself), this only returns the message text; the caller is
// responsible for ending the session.
//
// The FinTactical "offer to re-engage" prompt is interactive (the
// original called Gets() right here); since command functions in this
// package never perform I/O, that decision is instead expected to
// have already been made by the caller and passed in via
// reengageRequested, mirroring the original's "if (*buf==NULL ||
// *buf=='y' || *buf=='Y') { reengaged = 1; return; }" default-to-yes
// behavior.
func Final(st *game.State, mesg int, reengageRequested bool) []string {
	sp := st.Player()
	shipnum := len(st.Enemies())
	var msgs []string

	switch mesg {
	case game.FinFedLose:
		msgs = append(msgs, Starfleet()...)
		msgs = append(msgs,
			fmt.Sprintf("We have recieved confirmation that the U.S.S. %s,", sp.Name),
			fmt.Sprintf("   captained by %s, was destroyed by %s%s", st.Crew.Captain, ifStr(shipnum == 1, "a "), st.EnemyRaceName),
			fmt.Sprintf("   %s%s.  May future Federation officers", st.EnemyShipTypeName, plural(shipnum)),
			"   perform better in their duties.",
			"",
		)
	case game.FinEnemyLose:
		msgs = append(msgs, Starfleet()...)
		msgs = append(msgs,
			fmt.Sprintf("We commend Captain %s and the crew of the %s on their", st.Crew.Captain, sp.Name),
			fmt.Sprintf("   fine performance against the %ss.  They are", st.EnemyRaceName),
			"   an inspiration to all Starfleet personnel.",
			"",
		)
	case game.FinTactical:
		if !st.Reengaged {
			if reengageRequested {
				st.Reengaged = true
				return []string{
					fmt.Sprintf("%s:  %s, we are in a position to either disengage from the", st.Crew.Science, st.Crew.Title),
					fmt.Sprintf("   %ss, or re-engage them in combat.", st.EnemyRaceName),
					"   Do you wish to re-engage?",
				}
			}
		}
		msgs = append(msgs, Starfleet()...)
		msgs = append(msgs,
			fmt.Sprintf("Captain %s of the starship %s has", st.Crew.Captain, sp.Name),
			fmt.Sprintf("   out-maneuvered %s aggressors.  We commend", st.EnemyRaceName),
			"   his tactical ability.",
			"",
		)
	case game.FinFedSurrender:
		msgs = append(msgs, Starfleet()...)
		msgs = append(msgs,
			fmt.Sprintf("Captain %s has surrendered the U.S.S. %s ", st.Crew.Captain, sp.Name),
			fmt.Sprintf("   to the %ss.  May Captain Donsell be remembered.", st.EnemyRaceName),
			"",
		)
	case game.FinEnemySurrender:
		msgs = append(msgs, Starfleet()...)
		msgs = append(msgs,
			fmt.Sprintf("We have recieved word from the %s that the", sp.Name),
			fmt.Sprintf("   %ss have surrendered.", st.EnemyRaceName),
			"",
		)
	case game.FinComplete:
		msgs = append(msgs, Starfleet()...)
		msgs = append(msgs,
			"One of our scout vessels has encountered the wreckage of",
			fmt.Sprintf("   the %s and %d other %s vessel%s.", sp.Name, shipnum, st.EnemyRaceName, plural(shipnum)),
			"",
		)
	default:
		msgs = append(msgs, fmt.Sprintf("How did we get here? final(%d)", mesg))
		return msgs
	}

	msgs = append(msgs, "")
	var survivors []string
	liveShips := 0
	for _, ep := range st.Ships {
		if ep.IsDead(game.SysDead) || ep.Complement <= 0 {
			survivors = append(survivors, fmt.Sprintf("   %s -- destroyed", ep.Name))
		} else {
			liveShips++
			survivors = append(survivors, fmt.Sprintf("   %s -- %d", ep.Name, ep.Complement))
		}
	}
	if liveShips > 0 {
		msgs = append(msgs, "Survivors Reported:", "")
		msgs = append(msgs, survivors...)
	} else {
		msgs = append(msgs, "*** No survivors reported ***", "")
	}
	return msgs
}

// ifStr returns yes if cond is true, and "" otherwise; a small local
// helper for the inline conditional text fragments final() used (e.g.
// "a " only when facing a single enemy ship).
func ifStr(cond bool, yes string) string {
	if cond {
		return yes
	}
	return ""
}

// Warn renders a non-terminal disposition notice for the given Fin*
// condition, ported from warn() in endgame.c. Unlike the original
// (which used a static beenhere[] array to show each message only
// once per game), the caller is expected to track that itself (e.g.
// via st.Shutup) since Warn is a pure function with no memory of its
// own between calls.
func Warn(st *game.State, mesg int) []string {
	if st.Reengaged && mesg == game.FinTactical {
		return nil
	}
	switch mesg {
	case game.FinFedLose:
		return []string{
			"Message to the Federation:  This is Commander",
			fmt.Sprintf("   %s of the %s %s.  We have defeated", st.EnemyCommander, st.EnemyRaceName, st.EnemyEmpireName),
			fmt.Sprintf("   the %s and are departing the quadrant.", st.Player().Name),
		}
	case game.FinEnemyLose:
		return []string{
			fmt.Sprintf("%s: All %s vessels have been either", st.Crew.Science, st.EnemyRaceName),
			"   destroyed or crippled.  We still, however, have",
			"   antimatter devices to avoid.",
		}
	case game.FinTactical:
		return []string{
			fmt.Sprintf("%s: The %ss are falling behind and seem to", st.Crew.Helmsman, st.EnemyRaceName),
			"   be breaking off their attack.",
		}
	case game.FinFedSurrender:
		return []string{
			fmt.Sprintf("%s: I'm informing Starfleet Command of our ", st.Crew.Com),
			"   disposition.",
		}
	case game.FinEnemySurrender:
		return []string{
			fmt.Sprintf("%s: Although the %ss have surrendered,", st.Crew.Science, st.EnemyRaceName),
			"   there are still antimatter devices floating",
			"   around us.",
		}
	default:
		return []string{fmt.Sprintf("How did we get here? final(%d)", mesg)}
	}
}

// Mission renders the game's opening mission briefing, ported from
// mission()/missionlog() in mission.c. missionIndex selects the
// briefing from data.MissionBriefings (the caller supplies it, e.g.
// via r.Randm(len(data.MissionBriefings)), so this function stays a
// pure, deterministic renderer). stardate is a caller-supplied string
// (e.g. derived from the current date) since the original derived it
// from the wall-clock date, which isn't this port's concern to
// reproduce exactly.
func Mission(st *game.State, missionIndex int, stardate string) []string {
	sp := st.Player()
	shipnum := len(st.Enemies())
	onef := shipnum == 1

	msgs := []string{
		"", "", "",
		"Space, the final frontier.",
		fmt.Sprintf("These are the voyages of the starship %s.", sp.Name),
		"Its five year mission: to explore strange new worlds,",
		"to seek out new life and new civilizations,",
		"to boldly go where no man has gone before!",
		"",
		"                    S T A R    T R E K",
		"",
		fmt.Sprintf("%s:  Captain's log, stardate %s", st.Crew.Captain, stardate),
	}

	if missionIndex >= 0 && missionIndex < len(data.MissionBriefings) {
		b := data.MissionBriefings[missionIndex]
		msgs = append(msgs, b.Line1)
		if b.Line2 != "" {
			msgs = append(msgs, "   "+b.Line2)
		}
	}

	countWord := fmt.Sprintf("%d", shipnum)
	if onef {
		countWord = "a"
	}
	identify := "them"
	article := ""
	if onef {
		identify = "it"
		article = "a" + vowelStr(st.EnemyRaceName) + " "
	}
	msgs = append(msgs,
		fmt.Sprintf("%s:  %s, I'm picking up %s vessel%s on an interception", st.Crew.Helmsman, st.Crew.Title, countWord, plural(shipnum)),
		fmt.Sprintf("   course with the %s.", sp.Name),
		fmt.Sprintf("%s:  Sensors identify %s as %s%s %s%s,", st.Crew.Science, identify, article, st.EnemyRaceName, st.EnemyShipTypeName, plural(shipnum)),
		fmt.Sprintf("   probably under the command of Captain %s.", st.EnemyCommander),
		fmt.Sprintf("%s:  Sound general quarters, Lieutenant!", st.Crew.Captain),
		fmt.Sprintf("%s:  Aye, %s!", st.Crew.Com, st.Crew.Title),
	)
	return msgs
}

// Alert renders the "computer is attacking" notice shown just before
// the first turn, ported from alert() in mission.c.
func Alert(st *game.State) []string {
	sp := st.Player()
	enemies := st.Enemies()
	var names string
	if len(enemies) == 1 {
		names = enemies[0].Name
	} else {
		for i, ep := range enemies {
			if i == len(enemies)-1 {
				names += "and the "
			}
			names += ep.Name
			if i == len(enemies)-1 {
				continue
			}
			names += ", "
			if i == 0 || i == 5 {
				names += "\n   "
			}
		}
	}
	return []string{fmt.Sprintf("Computer: The %ss are attacking the %s with the %s.", st.EnemyRaceName, sp.Name, names)}
}
