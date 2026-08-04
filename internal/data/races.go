package data

// RaceInfo holds the per-race data ported from struct race_info in
// structs.h and the aliens[] table in globals.c: naming pools for
// enemy ships and captains, and the tactical-behavior probabilities
// used by the enemy AI.
//
// The exact meaning of SurrenderChance vs SurrenderPChance is preserved
// from the original C comments as written; their precise use in enemy
// decision-making will be confirmed when the surrender/bluff logic
// (special.c, enemycom.c) is ported.
type RaceInfo struct {
	Name       string // Race name (e.g. "Klingon").
	EmpireName string // What the race calls its polity (e.g. "Empire").

	// SurrenderChance is the race's chance (0-100) of accepting a
	// demand from the player to surrender ("Chance they will accept a
	// surrender").
	SurrenderChance int
	// SurrenderPChance is the race's chance (0-100) described in the
	// original source as "Chance they will surrender to you".
	SurrenderPChance int
	// CorbomiteChance is the race's chance (0-100) of falling for a
	// corbomite bluff.
	CorbomiteChance int
	// DefenselessChance is the race's chance (0-100) of falling for a
	// defenseless ruse.
	DefenselessChance int
	// Attitude is a tuning factor consumed by enemy strategy logic.
	Attitude int

	// ShipNames is the pool of enemy ship names for this race (one per
	// possible enemy vessel).
	ShipNames []string
	// ShipTypes is the descriptive ship-type name for each of the four
	// standard ship classes (Dreadnought, Heavy Cruiser, Light Cruiser,
	// Destroyer, in that order).
	ShipTypes []string
	// Captains is the pool of exemplary enemy captain names for this
	// race.
	Captains []string
}

// MontyPythonRaceIndex is the index of the "Monty Python" joke race in
// Races. The original game only offers this race as an enemy when the
// "silly" option is enabled (see init.c); it is preserved here as
// original game content, not removed.
const MontyPythonRaceIndex = 8

// MaxEnemyShipsPerRace is the number of enemy ship names provided per
// race (MAXESHIPS in defines.h).
const MaxEnemyShipsPerRace = 9

// MaxCaptainsPerRace is the number of enemy captain names provided per
// race (MAXECAPS in defines.h). Some races have fewer distinct captains
// than this and pad the remaining entries with empty strings, matching
// the original data.
const MaxCaptainsPerRace = 8

// Races are the alien races available as enemies, ported verbatim from
// aliens[] in globals.c.
var Races = []RaceInfo{
	{
		Name: "Klingon", EmpireName: "Empire",
		SurrenderChance: 25, SurrenderPChance: 25, CorbomiteChance: 75, DefenselessChance: 75, Attitude: 0,
		ShipNames: []string{
			"Annihilation", "Crusher", "Devastator", "Merciless",
			"Nemesis", "Pitiliess", "Ruthless", "Savage", "Vengeance",
		},
		ShipTypes: []string{
			"C-9 Dreadnought", "D-7 Battle Cruiser",
			"D-6 Light Battlecruiser", "F-5L Destroyer",
		},
		Captains: []string{
			"Koloth", "Kang", "Kor", "Krulix", "Korax", "Karg",
			"Kron", "Kumerian",
		},
	},
	{
		Name: "Romulan", EmpireName: "Star Empire",
		SurrenderChance: 5, SurrenderPChance: 5, CorbomiteChance: 80, DefenselessChance: 50, Attitude: 0,
		ShipNames: []string{
			"Avenger", "Defiance", "Fearless", "Harrower", "Intrepid",
			"Relentless", "Seeker", "Torch", "Vigilant",
		},
		ShipTypes: []string{
			"Condor Dreadnought", "Firehawk Heavy Cruiser",
			"Sparrowhawk Light Cruiser", "Skyhawk Destroyer",
		},
		Captains: []string{
			"Tal", "Tiercellus", "Diana", "Tama", "Subeus", "Turm",
			"Strell", "Scipio",
		},
	},
	{
		Name: "Kzinti", EmpireName: "Hegemony",
		SurrenderChance: 50, SurrenderPChance: 50, CorbomiteChance: 50, DefenselessChance: 50, Attitude: 0,
		ShipNames: []string{
			"Black Hole", "Comet", "Ecliptic", "Galaxy", "Meteor",
			"Nova", "Pulsar", "Quasar", "Satellite",
		},
		ShipTypes: []string{
			"Space Control Ship", "Strike Cruiser", "Light Cruiser",
			"Destroyer",
		},
		Captains: []string{
			"Hunter", "\"Cat Who Sleeps With Dogs\"", "Fellus", "Corda",
			"\"Cat Who Fought Fuzzy Bear\"", "", "", "",
		},
	},
	{
		Name: "Gorn", EmpireName: "Confederation",
		SurrenderChance: 80, SurrenderPChance: 50, CorbomiteChance: 50, DefenselessChance: 50, Attitude: 0,
		ShipNames: []string{
			"Chimericon", "Dragonicon", "Ornithocon", "Predatoricon",
			"Reptilicon", "Serpenticon", "Tyranicon", "Vipericon",
			"Wyvericon",
		},
		ShipTypes: []string{
			"Tyrannosaurus Rex Dreadnought", "Allosaurus Heavy Cruiser",
			"Megalosaurus Light Cruiser", "Carnosaurus Destroyer",
		},
		Captains: []string{
			"Sslith", "Dardiss", "Ssor", "Sslitz", "S'Arnath",
			"Zor", "", "",
		},
	},
	{
		Name: "Orion", EmpireName: "Pirates",
		SurrenderChance: 95, SurrenderPChance: 5, CorbomiteChance: 50, DefenselessChance: 60, Attitude: 0,
		ShipNames: []string{
			"Deuce Coupe", "Final Jeopardy", "Long John Dilithium",
			"Millennium Pelican", "Omega Race", "Penzance",
			"Road Warrior", "Scarlet Pimpernel", "Thunderduck",
		},
		ShipTypes: []string{
			"Battle Raider", "Heavy Cruiser", "Raider Cruiser",
			"Light Raider",
		},
		Captains: []string{
			"Daniel \"Deth\" O'Kay", "Neil Ricca", "Delilah Smith",
			"Hamilcar", "Pharoah", "Felna Greymane", "Hacker",
			"Credenza",
		},
	},
	{
		Name: "Hydran", EmpireName: "Monarchy",
		SurrenderChance: 50, SurrenderPChance: 50, CorbomiteChance: 50, DefenselessChance: 50, Attitude: 0,
		ShipNames: []string{
			"Bravery", "Chivalry", "Devotion", "Fortitude", "Loyalty",
			"Modesty", "Purity", "Resolution", "Tenacity",
		},
		ShipTypes: []string{
			"Paladin Dreadnought", "Ranger-class Cruiser",
			"Horseman Light Cruiser", "Lancer Destroyer",
		},
		Captains: []string{
			"Hypantspts", "S'Lenthna", "Hydraxan", "", "", "",
			"", "",
		},
	},
	{
		Name: "Lyran", EmpireName: "Empire",
		SurrenderChance: 50, SurrenderPChance: 50, CorbomiteChance: 50, DefenselessChance: 50, Attitude: 0,
		ShipNames: []string{
			"Bandit", "Claw", "Dangerous", "Fury", "Mysterious",
			"Sleek", "Tiger", "Vicious", "Wildcat",
		},
		ShipTypes: []string{
			"Lion Dreadnought", "Tiger Cruiser",
			"Panther Light Cruiser", "Leopard Destroyer",
		},
		Captains: []string{
			"Kleave", "Leyraf", "Kuhla", "Nashar",
			"Prekor", "Ffarric", "Rippke", "Larkahn",
		},
	},
	{
		Name: "Tholian", EmpireName: "Holdfast",
		SurrenderChance: 75, SurrenderPChance: 25, CorbomiteChance: 50, DefenselessChance: 50, Attitude: 0,
		ShipNames: []string{
			"Bismark", "Centaur", "Draddock", "Forbin", "Kreiger",
			"Shlurg", "Trakka", "Varnor", "Warrior",
		},
		ShipTypes: []string{
			"Tarantula Dreadnought", "Cruiser", "Improved Patrol Cruiser",
			"Patrol Cruiser",
		},
		Captains: []string{
			"Secthane", "Kotheme", "Sectin", "Brezgonne",
			"Loskene", "", "", "",
		},
	},
	{
		// Monty Python (MontyPythonRaceIndex): a joke race only offered
		// as an enemy when the original game's "silly" option is set.
		Name: "Monty Python", EmpireName: "Flying Circus",
		SurrenderChance: -1, SurrenderPChance: -1, CorbomiteChance: -1, DefenselessChance: -1, Attitude: 0,
		ShipNames: []string{
			"Blancmange", "Spam", "R.J. Gumby", "Lumberjack",
			"Dennis Moore", "Ministry of Silly Walks", "Argument Clinic",
			"Piranha Brothers", "Upper Class Twit of the Year",
		},
		ShipTypes: []string{
			"Thingee", "Thingee", "Thingee", "Thingee",
		},
		Captains: []string{
			"Cleese", "Chapman", "Idle", "Jones", "Gilliam", "Bruce",
			"Throatwobblermangrove", "Arthur \"Two Sheds\" Jackson",
		},
	},
}
