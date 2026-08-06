package game

// SpaceObject represents a torpedo, antimatter probe, or jettisoned
// engineering section in flight, ported from struct torpedo in
// structs.h.
//
// The original game distinguished these via a separate linked-list node
// type (I_TORPEDO, I_PROBE, I_ENG in the struct list union); Go's slices
// and a Type field replace that manual linked list, since the list
// mechanics themselves are an implementation detail rather than a
// gameplay behavior worth preserving.
type SpaceObject struct {
	ID   int // Unique identifier (matches original slot numbering).
	Type int // One of the Object* constants (Torpedo, Probe, Engineering).

	From *Ship // Ship that launched this object.

	X, Y int // Current position.

	Course   float64 // Direction of travel, 0-360 degrees.
	Speed    float64 // Current speed.
	NewSpeed float64 // Target speed, for objects still accelerating.

	// Target is only meaningful for probes, which (unlike unguided
	// torpedoes) retain a live target and can be steered/detonated
	// after launch.
	Target *Ship

	Fuel int // Antimatter pods carried; determines blast yield.

	TimeDelay float64 // Seconds remaining until the detonation timer expires.
	Proximity int     // Proximity-fuse trigger distance.
	Detonated bool    // Set once the object has detonated so it cannot fire again this turn before removal.
}
