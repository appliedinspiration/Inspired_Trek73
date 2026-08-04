package game

import (
	"fmt"
	"math"
	"strings"

	"github.com/appliedinspiration/inspired_trek73/internal/data"
)

// InitOptions configures a new battle, mirroring the configurable
// portions of init_ships() in init.c. Fields left at their zero value
// take the same defaults the original game used when a setting was not
// supplied via command-line options or the environment.
type InitOptions struct {
	// EnemyCount is the number of enemy ships. Zero means "choose
	// randomly", matching the original prompting for "1-9 enemy
	// vessels" and defaulting to a random count when given a blank
	// answer.
	EnemyCount int

	// PlayerClassAbbr and EnemyClassAbbr select ship classes by their
	// two-letter abbreviation (e.g. "CA"). An empty string defaults to
	// "CA", matching the original.
	//
	// The original game also allowed loading a custom class definition
	// from a binary file ($HOME/.trekXX, written by the shipyard tool).
	// That mechanism is not yet supported here; per project decisions
	// it will be replaced with a portable format rather than ported
	// as-is.
	PlayerClassAbbr string
	EnemyClassAbbr  string

	// RaceName selects the enemy race by (case-sensitive) name prefix,
	// matching the original's strncmp-based lookup. An empty string
	// means "choose randomly".
	RaceName string

	// AllowSillyRace mirrors the original "silly" option: when false,
	// the Monty Python joke race is excluded from random selection (it
	// may still be selected explicitly via RaceName).
	AllowSillyRace bool

	// PlayerShipName names the player's ship. An empty string picks a
	// random name from data.FederationShipNames, matching the original.
	PlayerShipName string
}

// InitResult is the outcome of a successful NewGame call.
type InitResult struct {
	State *State

	// EnemyRaceName, EnemyEmpireName, EnemyShipTypeName, and
	// EnemyCommander mirror foerace, empire, foestype, and foename in
	// the original globals.c: they describe the enemy race/ships for
	// the mission-briefing text, which is not yet implemented here.
	EnemyRaceName     string
	EnemyEmpireName   string
	EnemyShipTypeName string
	EnemyCommander    string

	// Warnings holds non-fatal notices about fallback choices the
	// original game printed directly (e.g. an unrecognized race or
	// class name), left to the caller to display.
	Warnings []string
}

// NewGame creates a new battle: it selects the enemy race and ship
// classes, creates the player's ship and all enemy ships, and positions
// them, mirroring init_ships() in init.c.
//
// Unlike the original, this function is pure with respect to package
// state: it does not mutate the shared data.Races tables, takes its
// randomness from the supplied Rand, and returns validation errors
// instead of exiting the process.
func NewGame(opts InitOptions, r *Rand) (*InitResult, error) {
	if opts.EnemyCount != 0 && (opts.EnemyCount < 1 || opts.EnemyCount > MaxEnemyShips) {
		return nil, fmt.Errorf("game: enemy count must be from 1 to %d", MaxEnemyShips)
	}

	res := &InitResult{}

	shipCount := opts.EnemyCount
	if shipCount == 0 {
		shipCount = r.Randm(MaxEnemyShips)
	}

	raceIndex, warn := chooseRace(opts.RaceName, opts.AllowSillyRace, r)
	if warn != "" {
		res.Warnings = append(res.Warnings, warn)
	}
	race := data.Races[raceIndex]

	playerClassAbbr := opts.PlayerClassAbbr
	if playerClassAbbr == "" {
		playerClassAbbr = "CA"
	}
	playerClass, ok := data.ShipClassByAbbr(playerClassAbbr)
	if !ok {
		res.Warnings = append(res.Warnings, fmt.Sprintf("Unknown player ship class %q; using CA.", playerClassAbbr))
		playerClass, _ = data.ShipClassByAbbr("CA")
	}

	enemyClassAbbr := opts.EnemyClassAbbr
	if enemyClassAbbr == "" {
		enemyClassAbbr = "CA"
	}
	enemyClass, ok := data.ShipClassByAbbr(enemyClassAbbr)
	if !ok {
		res.Warnings = append(res.Warnings, fmt.Sprintf("Unknown enemy ship class %q; using CA.", enemyClassAbbr))
		enemyClass, _ = data.ShipClassByAbbr("CA")
	}

	// The enemy commanding officer's name is drawn from the race's
	// captain pool, re-rolling past any blank entries (several races'
	// captain lists are padded with empty strings), matching the
	// while-loop in init.c.
	foeCaptain := ""
	for foeCaptain == "" {
		foeCaptain = race.Captains[r.Randm(len(race.Captains))-1]
	}

	// Shuffle a private copy of the race's ship-name pool so repeated
	// games don't mutate shared package state, unlike the original,
	// which shuffled the global aliens[] table in place.
	shipNames := append([]string(nil), race.ShipNames...)
	shuffleNames(shipNames, r)

	st := NewState()

	// Create the player's ship first so it always occupies index 0,
	// matching shiplist[0] in the original.
	player := newShipFromClass(st, playerClass, playerClass.OwnEff, playerClass.OwnCrew, playerClass.OwnWarpMax)
	player.Course, player.NewCourse = 0, 0
	player.X, player.Y = 0, 0
	player.Class = playerClass.Abbr
	player.Cloaking = CloakNone
	if opts.PlayerShipName != "" {
		player.Name = opts.PlayerShipName
	} else {
		player.Name = data.FederationShipNames[r.Randm(len(data.FederationShipNames))-1]
	}
	st.Ships = append(st.Ships, player)

	for i := 1; i <= shipCount; i++ {
		enemy := newShipFromClass(st, enemyClass, enemyClass.EnemyEff, enemyClass.EnemyCrew, enemyClass.EnemyWarpMax)
		enemy.Name = shipNames[i-1]
		enemy.Class = enemyClass.Abbr
		enemy.Course = float64(r.Randm(360))
		enemy.NewCourse = enemy.Course

		// Position enemy ships in a rough ring around the player, each
		// successively slightly closer, matching the original's
		// "range = 4100 + randm(300) - i*200" spacing. The player is
		// always at the origin, so its own (never-visible) placement
		// draw is skipped rather than reproduced.
		rng := 4100 + r.Randm(300) - i*200
		bear := toRadians(float64(r.Randm(360)))
		enemy.X = int(float64(rng) * math.Cos(bear))
		enemy.Y = int(float64(rng) * math.Sin(bear))

		if race.Name == "Romulan" {
			enemy.Cloaking = CloakOff
		} else {
			enemy.Cloaking = CloakNone
		}
		st.Ships = append(st.Ships, enemy)
	}

	res.State = st
	res.EnemyRaceName = race.Name
	res.EnemyEmpireName = race.EmpireName
	res.EnemyShipTypeName = race.ShipTypes[enemyClass.ClassNum]
	res.EnemyCommander = foeCaptain

	st.EnemyRaceName = race.Name
	st.EnemyEmpireName = race.EmpireName
	st.EnemyShipTypeName = race.ShipTypes[enemyClass.ClassNum]
	st.EnemyCommander = foeCaptain
	st.RaceChances = RaceChances{
		Defenseless: race.DefenselessChance,
		Corbomite:   race.CorbomiteChance,
		Surrender:   race.SurrenderChance,
		SurrenderP:  race.SurrenderPChance,
	}
	return res, nil
}

// newShipFromClass builds a ship with all the per-class defaults
// shared between the player ship and enemy ships, matching the common
// initialization block in init_ships(). Course, position, name, and
// class-specific overrides not shared between own/enemy ships are set
// by the caller afterward.
func newShipFromClass(st *State, class data.ShipClass, eff float64, crew int, maxSpeed int) *Ship {
	phasers, tubes := initWeapons(class.NumPhaser, class.NumTorp)

	s := &Ship{
		ID:                   st.nextObjectID(),
		Warp:                 1.0,
		NewWarp:              1.0,
		Phasers:              phasers,
		PhaserSpread:         InitPhaserSpread,
		PhaserFirePct:        InitPhaserPercent,
		PhaserBlindLeft:      class.PhaserBlindLeft,
		PhaserBlindRight:     class.PhaserBlindRight,
		Tubes:                tubes,
		TubeProximity:        InitTubeProx,
		TubeDelay:            InitTubeTime,
		TubeLaunchSpd:        InitTubeSpeed,
		TubeBlindLeft:        class.TubeBlindLeft,
		TubeBlindRight:       class.TubeBlindRight,
		Eff:                  eff,
		Regen:                class.Regen,
		Energy:               class.Energy,
		Pods:                 class.Pods,
		Complement:           crew,
		Delay:                10000,
		OrigMaxSpeed:         float64(maxSpeed),
		MaxSpeed:             float64(maxSpeed),
		DegPerTurn:           float64(class.TurnRate),
		PhaserShieldDivisor:  class.PhaserShieldDivisor,
		TorpedoShieldDivisor: class.TorpedoShieldDivisor,
		CloakEnergy:          class.CloakingEnergy,
		CloakDelay:           CloakDelay,
		Strategy:             "standard",
		PhaserFiringDelay:    class.PhaserFiringDelay,
		TorpedoFiringDelay:   class.TorpedoFiringDelay,
	}
	for j := range s.Shields {
		s.Shields[j] = Shield{Eff: 1.0, Drain: 1.0, AttemptDrain: 1.0}
	}
	return s
}

// initWeapons builds freshly initialized phaser and tube slices for a
// ship with the given weapon counts, using the original's per-count
// initial bearing tables.
func initWeapons(numPhasers, numTubes int) ([]Phaser, []Tube) {
	phasers := make([]Phaser, numPhasers)
	for j := range phasers {
		phasers[j] = Phaser{
			Bearing: data.PhaserInitialBearings[numPhasers][j],
			Load:    InitPhaserLoad,
			Drain:   InitPhaserDrain,
			Status:  PhaserNormal,
		}
	}
	tubes := make([]Tube, numTubes)
	for j := range tubes {
		tubes[j] = Tube{
			Bearing: data.TubeInitialBearings[numTubes][j],
			Load:    InitTubeLoad,
			Status:  TubeNormal,
		}
	}
	return phasers, tubes
}

// chooseRace resolves the RaceName option to a race index, falling back
// to a random choice (excluding Monty Python unless allowSilly is set)
// when name is empty or unrecognized. It returns a non-empty warning
// string when a requested name could not be found, matching the
// original's printed "Cannot find race %s." message.
func chooseRace(name string, allowSilly bool, r *Rand) (int, string) {
	if name != "" {
		for i, race := range data.Races {
			if strings.HasPrefix(race.Name, name) {
				return i, ""
			}
		}
		idx := randomRaceIndex(allowSilly, r)
		return idx, fmt.Sprintf("Cannot find race %s.", name)
	}
	return randomRaceIndex(allowSilly, r), ""
}

// randomRaceIndex picks a random race index, excluding Monty Python
// unless allowSilly is true, matching the offset logic in init.c.
func randomRaceIndex(allowSilly bool, r *Rand) int {
	offset := 1
	if allowSilly {
		offset = 0
	}
	return r.Randm(MaxFoeRaces-offset) - 1
}

// shuffleNames performs the same 20-swap shuffle the original game used
// on its enemy ship-name pool.
func shuffleNames(names []string, r *Rand) {
	for i := 0; i < 20; i++ {
		a := r.Randm(len(names)) - 1
		b := r.Randm(len(names)) - 1
		names[a], names[b] = names[b], names[a]
	}
}

// toRadians converts degrees to radians, matching the toradians(x)
// macro in defines.h. The original macro used a truncated constant
// (0.0174533) rather than a precise pi/180; this uses the precise
// value, a deliberate accuracy improvement with no meaningful gameplay
// impact.
func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}
