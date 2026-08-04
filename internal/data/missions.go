package data

// MissionBriefing is one randomly-selectable mission assignment shown
// at the start of a game, ported from the missiontab[] string pairs in
// mission.c's missionlog(). Line2 is empty for entries that were a
// single line in the original (a bare 0 terminator in missiontab).
//
// The original mission blurbs referenced specific Star Trek (TOS/TMP)
// episodes/films in comments; those comments are omitted here since
// they aren't part of the in-game text, but the blurb text itself
// (including a few of the original's small typos, e.g. "resue" and
// "vincinity") is preserved verbatim as flavor text.
type MissionBriefing struct {
	Line1 string
	Line2 string
}

// MissionBriefings mirrors mission.c's missiontab[] table.
var MissionBriefings = []MissionBriefing{
	{"   We are acting in response to a Priority 1 distress call from", "space station K7."},
	{"   We are orbiting Gamma 2 to make a routine check of automatic", "communications and astrogation stations."},
	{"   We are on course for Epsilon Canares 3 to treat Commissioner", "Headford for Sukaro's disease."},
	{"   We have been assigned to transport ambassadors to a diplomatic", "conference on the planet code named Babel."},
	{"   Our mission is to investigate a find of tritanium on Beta 7.", ""},
	{"   We are orbiting Rigel 4 for therapeutic shore leave.", ""},
	{"   We are orbiting Sigma Iota 2 to study the effects of", "contamination upon a developing culture."},
	{"   We have altered course for a resue mission on the Gamma 7A", "system."},
	{"   We are presently on course for Altair 6 to attend inauguration", "cermonies on the planet."},
	{"   We are on a cartographic mission to Pollux 9.", ""},
	{"   We are headed for Malurian in response to a distress call", "from that system."},
	{"   We are to negotiate a treaty to mine dilithium crystals from", "the Halkans."},
	{"   We are to investigate strange sensor readings reported by a", "scoutship investigating Gamma Triangula 6."},
	{"   We are headed for planets L370 and L374 to investigate the", "disappearance of the starship Constellation in that vincinity."},
	{"   We are ordered, with a skeleton crew, to proceed to Space", "Station K2 to test Dr. Richard Daystrom's computer M5."},
	{"   We have encountered debris from the SS Beagle and are", "proceeding to investigate."},
	{"   We are on course for Ekos to locate John Gill.", ""},
	{"   We are to divert an asteroid from destroying an inhabited", "planet."},
	{"   We are responding to a distresss call form the scientific", "expedition on Triacus."},
	{"   We have been assigned to transport the Medusan Ambassador to", "to his home planet."},
	{"   We are within the Neutral Zone on a mission to rescue the", "Kobayashi Maru."},
}
