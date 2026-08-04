package game

// MovementHooks lets the movement simulation invoke combat/AI
// behavior that is ported in later phases (phaser and torpedo firing,
// detonation, and enemy cloak shutdown), without move_ships() needing
// to depend directly on those packages/files.
//
// This mirrors calls the original move_ships() and distribute() made
// directly into firing.c/damage.c/special.c; splitting them out as an
// interface keeps movement testable in isolation before those pieces
// exist, and avoids a circular/monolithic single file for the whole
// simulation.
type MovementHooks interface {
	// PhaserFiring is called for a ship whose phaser-firing delay has
	// elapsed this segment, ported from the phaser_firing(sp) call in
	// move_ships() (firing.c).
	PhaserFiring(sp *Ship)

	// TorpedoFiring is called for a ship whose tube-firing delay has
	// elapsed this segment, ported from the torpedo_firing(sp) call in
	// move_ships() (firing.c).
	TorpedoFiring(sp *Ship)

	// TorpDetonate is called when a torpedo/probe's time or proximity
	// fuse has expired, ported from torp_detonate() calls in
	// move_ships() (firing.c).
	TorpDetonate(obj *SpaceObject)

	// ShipDetonate is called when a ship's self-destruct countdown
	// reaches zero, ported from the ship_detonate() call in
	// move_ships() (moveships.c/special.c).
	ShipDetonate(sp *Ship)

	// ECloakOff is called when an enemy ship runs out of energy to
	// sustain its cloaking device, ported from the e_cloak_off(sp,
	// fed) call in distribute() (dist.c/special.c). Returns true if
	// the cloak was turned off.
	ECloakOff(sp *Ship, fed *Ship) bool
}

// NoOpHooks is a MovementHooks implementation that does nothing,
// useful for testing movement/power-distribution mechanics in
// isolation before combat resolution is wired in.
type NoOpHooks struct{}

func (NoOpHooks) PhaserFiring(sp *Ship)              {}
func (NoOpHooks) TorpedoFiring(sp *Ship)             {}
func (NoOpHooks) TorpDetonate(obj *SpaceObject)      {}
func (NoOpHooks) ShipDetonate(sp *Ship)              {}
func (NoOpHooks) ECloakOff(sp *Ship, fed *Ship) bool { return false }
