package game

// Constants ported from References/FreeBSD/trek73/src/defines.h and
// structs.h. Names are translated to idiomatic Go while preserving the
// original numeric values, since these values directly affect gameplay
// balance.

// System indices, mirroring the S_* defines used to index a ship's
// damage-status array.
const (
	SysComputer    = iota // S_COMP
	SysSensor             // S_SENSOR
	SysProbe              // S_PROBE
	SysWarp               // S_WARP
	SysEngineering        // S_ENG
	SysDead               // S_DEAD
	SysSurrender          // S_SURRENDER
)

// NumDamageSystems is the number of systems that have damage
// descriptions (S_NUMSYSTEMS): computer, sensors, probe launcher, warp.
const NumDamageSystems = 4

// MaxSystems is the size of a ship's status array (MAXSYSTEMS).
const MaxSystems = 7

// Fully damaged systems are represented as 100% in the original; a
// named constant makes the "is dead" check self-documenting.
const FullyDamaged = 100

// Weapon and shield array sizing (MAXWEAPONS, SHIELDS).
const (
	MaxWeapons = 11
	NumShields = 4
)

// Shield array indices: shield 0 is the forward shield and receives a
// bonus multiplier (SHIELD1).
const (
	ShieldForward = 0
	// ShieldForwardBonus is the extra effective strength of the forward
	// shield relative to the others.
	ShieldForwardBonus = 1.5
)

// Play/global status flags (NORMAL, F_SURRENDER, E_SURRENDER).
const (
	StatusNormal         = 0
	StatusFedSurrender   = 1 << 0
	StatusEnemySurrender = 1 << 1
)

// Phaser bank status (P_NORMAL, P_DAMAGED, P_FIRING).
const (
	PhaserNormal  = 0
	PhaserDamaged = 1 << 0
	PhaserFiring  = 1 << 1
)

// Torpedo tube status (T_NORMAL, T_DAMAGED, T_FIRING).
const (
	TubeNormal  = 0
	TubeDamaged = 1 << 0
	TubeFiring  = 1 << 1
)

// Probe launcher status (PR_NORMAL, PR_LAUNCHING, PR_DETONATE, PR_LOCK).
const (
	ProbeNormal    = 0
	ProbeLaunching = 1 << 0
	ProbeDetonate  = 1 << 1
	ProbeLocked    = 1 << 2
)

// Cloaking device status/capability (C_NONE, C_OFF, C_ON).
const (
	CloakNone = 0 // Ship cannot cloak at all.
	CloakOff  = 1 // Ship can cloak, but is not currently cloaked.
	CloakOn   = 2 // Ship is currently cloaked.
)

// CloakDelay is the number of turns after decloaking before a ship may
// cloak again.
const CloakDelay = 2

// Damage source types (D_PHASER, D_ANTIMATTER).
const (
	DamagePhaser     = 0
	DamageAntimatter = 1
)

// HitPerPod and ProxPerPod are conversion constants used when computing
// antimatter blast radii from fuel pod counts.
const (
	HitPerPod  = 5
	ProxPerPod = 50
)

// Roster sizing for enemy ships, Federation ship names, enemy captains,
// and alien races (MAXESHIPS, MAXFEDS, MAXECAPS, MAXFOERACES).
const (
	MaxEnemyShips    = 9
	MaxFedShips      = 9
	MaxEnemyCaptains = 8
	MaxFoeRaces      = 9
)

// MaxShipClass is the number of standard ship classes (DN, CA, CL, DD).
const MaxShipClass = 4

// Turn cost flags for commands (TURN, FREE).
const (
	TurnCosting = 1
	TurnFree    = 0
)

// Space object types (TP_TORPEDO, TP_PROBE, TP_ENGINEERING).
const (
	ObjectTorpedo     = 0
	ObjectProbe       = 1
	ObjectEngineering = 2
)

// Simulation timing (segment, timeperturn in globals.c).
const (
	// SegmentSeconds is the granularity of the movement/combat
	// simulation within a single player turn.
	SegmentSeconds = 0.05
	// SecondsPerTurn is the length of a single player turn.
	SecondsPerTurn = 2.0
)

// DefaultTimeDelay is the default number of real-time seconds a player
// has to enter a command before the turn resolves automatically
// (DEFAULT_TIME).
const DefaultTimeDelay = 30

// Weapon/system tuning constants (MIN_/MAX_ defines).
const (
	MinPhaserSpread = 10
	MaxPhaserSpread = 45
	MaxPhaserRange  = 1000
	MaxPhaserCharge = 10.0
	MinPhaserDrain  = -MaxPhaserCharge
	MaxPhaserDrain  = MaxPhaserCharge
	MaxTubeCharge   = 10
	MaxTubeProx     = 500
	MaxTubeSpeed    = 12
	MaxTubeTime     = 10.0
	MaxProbeDelay   = 15
	MinProbeCharge  = 10
	MinProbeProx    = 50
	MinSensorRange  = 100
	MaxSensorRange  = 50000
)

// Endgame condition codes (FIN_* defines). The actual game-ending
// presentation (leftovers()/warn()/final() in main.c) is not yet
// ported; State.PendingFinal/PendingWarn record which condition was
// triggered so a future phase can act on it.
const (
	FinFedLose        = 0
	FinEnemyLose      = 1
	FinTactical       = 2
	FinFedSurrender   = 3
	FinEnemySurrender = 4
	FinComplete       = 5
)

// Initial weapon/system settings for newly created ships (INIT_* defines).
const (
	InitPhaserSpread  = MinPhaserSpread
	InitPhaserLoad    = MaxPhaserCharge
	InitPhaserDrain   = MaxPhaserDrain
	InitPhaserPercent = 100
	InitTubeLoad      = 0
	InitTubeProx      = 200
	InitTubeTime      = MaxTubeTime
	InitTubeSpeed     = MaxTubeSpeed
)
