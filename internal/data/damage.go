package data

// SystemDamageEntry describes the damage-roll threshold and message for
// a single ship system, ported from struct dam in structs.h.
//
// See damage.c: when a system takes a hit, a die is rolled with Roll
// sides; if the roll is less than the actual hit taken, the system
// accumulates additional damage.
type SystemDamageEntry struct {
	Roll    int    // Number of sides on the damage-check die.
	Message string // Message printed when this system is affected.
}

// DamageProfile describes how a weapon type's damage is distributed
// across a ship's stats and systems, ported from struct damage in
// structs.h and the p_damage/a_damage tables in globals.c.
//
// EffDivisor, FuelDivisor, RegenDivisor, and WeaponDivisor are divisors:
// the portion of a hit that penetrates shields is divided by the
// relevant field, and that amount is subtracted from the corresponding
// ship stat. CrewDivisor works the same way for crew casualties.
type DamageProfile struct {
	EffDivisor    float64
	FuelDivisor   float64
	RegenDivisor  float64
	CrewDivisor   float64
	WeaponDivisor float64

	// Systems holds one entry per system with a damage description, in
	// the same order as the game.Sys* system-index constants
	// (Computer, Sensor, Probe, Warp).
	Systems [4]SystemDamageEntry
}

// PhaserDamage is the damage profile applied by phaser hits, ported
// from p_damage in globals.c.
var PhaserDamage = DamageProfile{
	EffDivisor: 50, FuelDivisor: 2, RegenDivisor: 20, CrewDivisor: 10, WeaponDivisor: 3,
	Systems: [4]SystemDamageEntry{
		{Roll: 1000, Message: "Computer destroyed."},
		{Roll: 500, Message: "Sensors demolished."},
		{Roll: 100, Message: "Probe launcher crushed."},
		{Roll: 50, Message: "Warp drive destroyed."},
	},
}

// AntimatterDamage is the damage profile applied by torpedo/probe
// (antimatter) hits, ported from a_damage in globals.c.
var AntimatterDamage = DamageProfile{
	EffDivisor: 100, FuelDivisor: 3, RegenDivisor: 10, CrewDivisor: 7, WeaponDivisor: 6,
	Systems: [4]SystemDamageEntry{
		{Roll: 1500, Message: "Computer banks pierced."},
		{Roll: 750, Message: "Sensors smashed."},
		{Roll: 150, Message: "Probe launcher shot off."},
		{Roll: 75, Message: "Warp drive disabled."},
	},
}

// SystemNames holds the display name for each of the four systems that
// have damage descriptions (Computer, Sensors, Probe launcher, Warp
// Drive), ported from sysname[] in globals.c.
var SystemNames = []string{
	"Computer",
	"Sensors",
	"Probe launcher",
	"Warp Drive",
}

// SystemDestroyedMessages holds the message printed when a system
// reaches 100% damage, ported from statmsg[] in globals.c. Unlike
// SystemNames, this includes a fifth entry for Engineering.
var SystemDestroyedMessages = []string{
	"computer inoperable",
	"sensors annihilated",
	"probe launcher shot off",
	"warp drive disabled",
	"engineering jettisoned",
}
