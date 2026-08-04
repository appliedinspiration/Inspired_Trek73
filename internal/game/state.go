package game

// CrewNames holds the player crew's names and the captain's title,
// ported from the captain/science/engineer/title (and similar) global
// strings in globals.c. These are normally randomized per game by
// name_crew(), which has not been ported yet; until it is, State
// initializes these to generic placeholders so messages referencing
// them remain sensible.
type CrewNames struct {
	Captain  string
	Science  string
	Engineer string
	Com      string
	Nav      string
	Helmsman string
	Title    string // Captain's rank/title, e.g. "Captain".
}

// defaultCrewNames returns the generic placeholder crew roster used
// until name_crew() is ported.
func defaultCrewNames() CrewNames {
	return CrewNames{
		Captain:  "Captain",
		Science:  "Science officer",
		Engineer: "Engineer",
		Com:      "Communications officer",
		Nav:      "Navigator",
		Helmsman: "Helmsman",
		Title:    "Captain",
	}
}

// RaceChances holds the enemy race's tactical-response probabilities
// for the current game, ported from the aliens[enemynum] fields
// referenced in special.c. Every normal race has these fixed for the
// whole game; the Monty Python joke race instead starts every field
// at -1 ("not yet rolled") and randomly rolls each the first time it
// is needed, caching the result for the rest of the game - matching
// the original's "if (aliens[enemynum].X == -1) aliens[enemynum].X =
// randm(100)" pattern.
type RaceChances struct {
	Defenseless int
	Corbomite   int
	Surrender   int
	SurrenderP  int
}

// resolve returns *field, rolling and caching a fresh 1-100 value
// first if it is still -1 (unrolled).
func resolveChance(field *int, r *Rand) int {
	if *field == -1 {
		*field = r.Randm(100)
	}
	return *field
}

// State holds the complete state of a single Trek73 battle: all ships,
// all objects currently in flight, and the shared status flags that
// affect the whole engagement.
//
// This corresponds to the collection of global variables in the
// original C implementation (globals.c): shiplist[], the head/tail
// object list, and flags such as global, surrender, corbomite, etc.
// Consolidating them into a single struct removes the original's
// reliance on process-wide globals, without changing what state the
// game tracks.
type State struct {
	// Ships holds every ship in the battle. By convention, and
	// matching the original shiplist[0], index 0 is always the
	// player's ship.
	Ships []*Ship

	// Objects holds every torpedo, probe, and jettisoned engineering
	// section currently in flight.
	Objects []*SpaceObject

	// PlayStatus is a bitmask of Status* flags describing the overall
	// engagement state (e.g. StatusFedSurrender, StatusEnemySurrender).
	PlayStatus int

	// Crew holds the player crew's names and title, used in status
	// messages.
	Crew CrewNames

	// EnemyRaceName, EnemyEmpireName, EnemyShipTypeName, and
	// EnemyCommander mirror foerace, empire, foestype, and foename in
	// globals.c: they describe the enemy race/commander for AI
	// messages and mission-briefing text. Set by NewGame.
	EnemyRaceName     string
	EnemyEmpireName   string
	EnemyShipTypeName string
	EnemyCommander    string

	// RaceChances holds the enemy race's cached tactical-response
	// probabilities for this game (see RaceChances doc comment). Set
	// by NewGame from the chosen race's data.RaceInfo.
	RaceChances RaceChances

	// Shutup tracks which "first time only" status messages have
	// already been shown during the current turn. It must be reset
	// (via Shutup.Reset) at the start of each new turn.
	Shutup *Shutup

	// Defenseless, Corbomite, Surrender, and SurrenderP are simple
	// aging counters for in-progress ruses/bluffs, ported from the
	// like-named globals in globals.c (defenseless, corbomite,
	// surrender, surrenderp). They are zero when no such action is in
	// progress, and count up by one turn at a time once set by the
	// (not yet ported) special-actions commands in special.c.
	Defenseless int
	Corbomite   int
	Surrender   int
	SurrenderP  int

	// Reengaged mirrors the "reengaged" flag in globals.c: whether the
	// player is re-engaging enemy ships that had drifted out of range.
	Reengaged bool

	// PendingFinal and PendingWarn record an endgame condition
	// (Fin* constant, or -1 for none) that special-actions/disposition
	// logic determined should end the game, ported from the
	// final()/warn() calls in special.c/moveships.c. The actual
	// end-of-game presentation (leftovers(), scoring, etc.) is a later
	// phase; for now these are just recorded for the caller to notice.
	PendingFinal int
	PendingWarn  int

	// WarnShown tracks which Warn (Fin* code) messages have already
	// been shown once this game, indexed by Fin* constant, ported
	// from the "static int beenhere[5]" array in warn() (endgame.c).
	// commands.Warn() itself is a pure function with no memory of its
	// own; the CLI main loop is expected to check/set this before
	// calling Warn() so each message is only printed the first time
	// its condition arises.
	WarnShown [6]bool

	// nextID is used to hand out unique identifiers to ships and
	// objects, replacing the original's fixed-size slots[] array and
	// new_slot()/return_slot() bookkeeping.
	nextID int
}

// DefenselessChance, CorbomiteChance, SurrenderChance, and
// SurrenderPChance return (rolling and caching if necessary) the enemy
// race's chance of falling for each ruse/bluff/surrender action,
// ported from the aliens[enemynum].{defenseless,corbomite,surrender,
// surrenderp} lazy-roll pattern in special.c.
func (s *State) DefenselessChance(r *Rand) int { return resolveChance(&s.RaceChances.Defenseless, r) }
func (s *State) CorbomiteChance(r *Rand) int   { return resolveChance(&s.RaceChances.Corbomite, r) }
func (s *State) SurrenderChance(r *Rand) int   { return resolveChance(&s.RaceChances.Surrender, r) }
func (s *State) SurrenderPChance(r *Rand) int  { return resolveChance(&s.RaceChances.SurrenderP, r) }

// AgeSpecialActionTimers increments every in-progress ruse/bluff
// counter that is currently active (non-zero), ported from the
// "Ruses, bluffs, surrenders" block at the end of misc_timers() in
// moveships.c.
func (s *State) AgeSpecialActionTimers() {
	if s.Defenseless != 0 {
		s.Defenseless++
	}
	if s.Corbomite != 0 {
		s.Corbomite++
	}
	if s.Surrender != 0 {
		s.Surrender++
	}
	if s.SurrenderP != 0 {
		s.SurrenderP++
	}
}

// NewState returns an empty game state ready to have ships and objects
// added to it.
func NewState() *State {
	return &State{
		PlayStatus:   StatusNormal,
		Crew:         defaultCrewNames(),
		Shutup:       NewShutup(),
		PendingFinal: -1,
		PendingWarn:  -1,
	}
}

// Player returns the player's ship. It panics if no ships have been
// added yet, since a State is never meaningful without one.
func (s *State) Player() *Ship {
	if len(s.Ships) == 0 {
		panic("game: State.Player called before any ships were added")
	}
	return s.Ships[0]
}

// Enemies returns every ship other than the player's.
func (s *State) Enemies() []*Ship {
	if len(s.Ships) <= 1 {
		return nil
	}
	return s.Ships[1:]
}

// nextObjectID returns a fresh, unique identifier for a new ship or
// space object.
func (s *State) nextObjectID() int {
	s.nextID++
	return s.nextID
}

// NextObjectID is the exported form of nextObjectID, for use by
// callers outside the game package (e.g. internal/commands) that need
// to create new space objects such as jettisoned engineering sections
// or launched probes.
func (s *State) NextObjectID() int {
	return s.nextObjectID()
}

// ShipByName finds a living enemy ship whose name has the given string
// as a prefix, ported from ship_name() in subs.c. Only enemies are
// searched, matching the original's "for (i=1; i<=shipnum; i++)" (the
// player's own ship, shiplist[0], is never a valid lock/pursue/elude/
// scan target). Matching is case insensitive (the original only
// uppercased the first character of the typed name before a
// case-sensitive comparison, which meant names typed in any case other
// than "first letter uppercase, rest exactly as stored" would fail to
// match; a fully case-insensitive comparison is used here instead as
// an intentional, low-risk improvement). Ships with negative
// Complement (destroyed) are skipped, matching the original's
// shiplist[i]->complement < 0 check.
func (s *State) ShipByName(name string) *Ship {
	if name == "" {
		return nil
	}
	lower := toLowerASCII(name)
	for _, sp := range s.Enemies() {
		if sp.Complement < 0 {
			continue
		}
		if len(lower) <= len(sp.Name) && toLowerASCII(sp.Name[:len(lower)]) == lower {
			return sp
		}
	}
	return nil
}

func toLowerASCII(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	return string(b)
}
