package game

// Phaser represents a single phaser bank, ported from struct phaser in
// structs.h.
type Phaser struct {
	Target  *Ship   // Locked target, if any.
	Bearing float64 // Direction aimed, relative to ship course, when not locked.
	Drain   int     // Energy drawn from (or fed back to, if negative) the engines.
	// Load is the energy currently stored in the bank. The original C
	// struct declared this as a 16-bit integer (short), which silently
	// truncated fractional charge each turn; that truncation behavior
	// is revisited when the charging/firing mechanics are ported rather
	// than baked into the struct definition here.
	Load   float64
	Status int // Bitmask of Phaser* status flags.
}

// Tube represents a single torpedo tube, ported from struct tube in
// structs.h.
type Tube struct {
	Target  *Ship   // Locked target, if any.
	Bearing float64 // Direction aimed, relative to ship course, when not locked.
	Load    float64 // Energy currently loaded into the tube.
	Status  int     // Bitmask of Tube* status flags.
}

// Shield represents one of a ship's four shield facings, ported from
// struct shield in structs.h.
type Shield struct {
	Eff          float64 // Current efficiency, from 0 to 1.
	Drain        float64 // Actual energy drain allocated by the engines.
	AttemptDrain float64 // Requested energy drain, before allocation limits.
}

// LastKnownPosition records a ship's last observed position, course,
// and speed before it cloaked, ported from struct ship's embedded
// "position" struct in structs.h. Enemy sensors track this instead of a
// cloaked ship's true position.
type LastKnownPosition struct {
	X, Y   int
	Warp   float64
	Range  int
	Bear   float64
	Course float64
}

// Ship represents a single starship, player or enemy, ported from
// struct ship in structs.h.
type Ship struct {
	Name  string
	Class string // Two-letter class abbreviation (e.g. "CA").

	X, Y int // Current position.

	Warp    float64 // Current warp speed.
	NewWarp float64 // Target warp speed, for in-progress speed changes.

	Course    float64 // Current heading, 0-360 degrees.
	NewCourse float64 // Target heading, for in-progress course changes.

	Target       *Ship   // Ship currently being pursued/elueded, if any.
	RelativeBear float64 // Relative bearing to maintain toward Target.

	Phasers       []Phaser
	PhaserSpread  int // Phaser spread setting, in degrees.
	PhaserFirePct int // Phaser firing percentage setting.

	PhaserBlindLeft  int // Start of the phaser blind arc, left side, in degrees.
	PhaserBlindRight int // Start of the phaser blind arc, right side, in degrees.

	Tubes []Tube

	TubeProximity  int // Proximity-fuse setting for launched torpedoes.
	TubeDelay      int // Time-delay-fuse setting for launched torpedoes.
	TubeLaunchSpd  int // Launch speed added to a torpedo's own speed.
	TubeBlindLeft  int // Start of the tube blind arc, left side, in degrees.
	TubeBlindRight int // Start of the tube blind arc, right side, in degrees.

	Shields [NumShields]Shield

	ProbeLauncherStatus int // Bitmask of Probe* status flags.

	Eff    float64 // Efficiency: multiplier controlling fuel consumption.
	Regen  float64 // Energy regenerated per turn.
	Energy float64 // Current effective energy.
	Pods   float64 // Maximum antimatter pod capacity.

	Complement int // Crew members currently alive.

	Status [MaxSystems]int // Damage percentage per system (see Sys* constants).

	Delay float64 // Seconds until self-destruct/detonation, if counting down.

	ID int // Unique identifier (matches original slot numbering).

	OrigMaxSpeed float64 // Maximum warp speed before any damage.
	MaxSpeed     float64 // Current maximum warp speed.
	DegPerTurn   float64 // Degrees of course change available per warp-second.

	PhaserShieldDivisor  float64 // Divisor applied to incoming phaser damage.
	TorpedoShieldDivisor float64 // Divisor applied to incoming torpedo damage.

	Cloaking    int // Cloaking device status/capability (Cloak* constants).
	CloakEnergy int // Energy required to run the cloaking device.
	CloakDelay  int // Turns remaining before the cloak may be reactivated.

	Strategy string // Identifier of the enemy strategy in use, if any.

	LastKnown LastKnownPosition

	PhaserFiringDelay  int // Delay, in segments, before phasers may fire.
	TorpedoFiringDelay int // Delay, in segments, before tubes may fire.
}

// IsDead reports whether the given system has been damaged to 100%,
// mirroring the is_dead(sp, sys) macro in defines.h.
func (s *Ship) IsDead(system int) bool {
	return s.Status[system] == FullyDamaged
}

// SystemWorks reports whether the given system functions this attempt,
// mirroring the syswork(sp, sys) macro in defines.h: a system with X%
// damage has a (100-X)% chance of working on any given check.
//
// The original macro reads: randm(100) > sp->status[sys]. It is
// preserved here as a method taking an explicit Rand so callers control
// determinism.
func (s *Ship) SystemWorks(r *Rand, system int) bool {
	return r.Randm(100) > s.Status[system]
}
