package data

// ShipClass holds the baseline statistics for a ship class, ported from
// struct ship_stat in structs.h and the stats[] table in ships.c.
//
// Own-ship and enemy-ship values differ for several fields (warp speed,
// efficiency, crew) because the original game intentionally gave the
// Federation ship and enemy ships of the "same" class different
// baselines; both are preserved here rather than collapsed into one.
type ShipClass struct {
	Abbr      string // Two-letter class abbreviation (e.g. "CA").
	ClassNum  int    // Index into the class array.
	NumPhaser int    // Number of phaser banks.
	NumTorp   int    // Number of torpedo tubes.

	OwnWarpMax   int // Maximum warp speed when flown by the player.
	EnemyWarpMax int // Maximum warp speed when flown by an enemy.

	OwnEff   float64 // Efficiency (fuel-use factor) when flown by the player.
	EnemyEff float64 // Efficiency when flown by an enemy.

	Regen  float64 // Energy regeneration per turn.
	Energy float64 // Starting energy.
	Pods   float64 // Maximum antimatter pod capacity.

	OwnCrew   int // Crew complement when flown by the player.
	EnemyCrew int // Crew complement when flown by an enemy.

	PhaserShieldDivisor  float64 // Divisor applied to incoming phaser damage.
	TorpedoShieldDivisor float64 // Divisor applied to incoming torpedo damage.

	TurnRate       int // Degrees of course change available per warp-second.
	CloakingEnergy int // Energy required to run the cloaking device.

	TubeBlindLeft    int // Start of the tube blind arc, left side, in degrees.
	TubeBlindRight   int // Start of the tube blind arc, right side, in degrees.
	PhaserBlindLeft  int // Start of the phaser blind arc, left side, in degrees.
	PhaserBlindRight int // Start of the phaser blind arc, right side, in degrees.

	PhaserFiringDelay  int // Delay, in segments, before phasers may fire.
	TorpedoFiringDelay int // Delay, in segments, before tubes may fire.
}

// ShipClasses are the four standard ship classes available in the
// original game: Dreadnought, Heavy Cruiser, Light Cruiser, and
// Destroyer. Values are ported verbatim from stats[] in ships.c.
var ShipClasses = []ShipClass{
	{
		Abbr: "DN", ClassNum: 0, NumPhaser: 6, NumTorp: 8,
		OwnWarpMax: 8, EnemyWarpMax: 10,
		OwnEff: 2.0, EnemyEff: 1.5,
		Regen: 15.0, Energy: 200, Pods: 300,
		OwnCrew: 600, EnemyCrew: 450,
		PhaserShieldDivisor: 5.0, TorpedoShieldDivisor: 4.0,
		TurnRate: 2, CloakingEnergy: 4,
		TubeBlindLeft: 135, TubeBlindRight: 225,
		PhaserBlindLeft: 125, PhaserBlindRight: 235,
		PhaserFiringDelay: 4, TorpedoFiringDelay: 4,
	},
	{
		Abbr: "CA", ClassNum: 1, NumPhaser: 4, NumTorp: 6,
		OwnWarpMax: 9, EnemyWarpMax: 11,
		OwnEff: 1.0, EnemyEff: 0.75,
		Regen: 10.0, Energy: 150, Pods: 200,
		OwnCrew: 450, EnemyCrew: 350,
		PhaserShieldDivisor: 3.0, TorpedoShieldDivisor: 2.0,
		TurnRate: 5, CloakingEnergy: 2,
		TubeBlindLeft: 135, TubeBlindRight: 225,
		PhaserBlindLeft: 125, PhaserBlindRight: 235,
		PhaserFiringDelay: 4, TorpedoFiringDelay: 4,
	},
	{
		Abbr: "CL", ClassNum: 2, NumPhaser: 4, NumTorp: 4,
		OwnWarpMax: 9, EnemyWarpMax: 11,
		OwnEff: 0.75, EnemyEff: 0.5,
		Regen: 10.0, Energy: 125, Pods: 175,
		OwnCrew: 350, EnemyCrew: 250,
		PhaserShieldDivisor: 3.0, TorpedoShieldDivisor: 2.0,
		TurnRate: 6, CloakingEnergy: 2,
		TubeBlindLeft: 150, TubeBlindRight: 210,
		PhaserBlindLeft: 140, PhaserBlindRight: 220,
		PhaserFiringDelay: 4, TorpedoFiringDelay: 4,
	},
	{
		Abbr: "DD", ClassNum: 3, NumPhaser: 2, NumTorp: 4,
		OwnWarpMax: 10, EnemyWarpMax: 12,
		OwnEff: 0.5, EnemyEff: 0.5,
		Regen: 8.0, Energy: 100, Pods: 150,
		OwnCrew: 200, EnemyCrew: 150,
		PhaserShieldDivisor: 2.0, TorpedoShieldDivisor: 1.5,
		TurnRate: 7, CloakingEnergy: 1,
		TubeBlindLeft: 160, TubeBlindRight: 200,
		PhaserBlindLeft: 150, PhaserBlindRight: 210,
		PhaserFiringDelay: 4, TorpedoFiringDelay: 4,
	},
}

// ShipClassByAbbr looks up a standard ship class by its two-letter
// abbreviation (e.g. "CA"), mirroring ship_class() in misc.c. It reports
// false if no such class exists.
func ShipClassByAbbr(abbr string) (ShipClass, bool) {
	for _, c := range ShipClasses {
		if c.Abbr == abbr {
			return c, true
		}
	}
	return ShipClass{}, false
}
