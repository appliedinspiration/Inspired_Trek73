package game

// Shutup tracks which "first time only" status messages have already
// been shown during the current turn, ported from the shutup[] global
// array in globals.c/defines.h. The original reset this array once per
// player command prompt (see playit() in main.c), so message
// suppression only applies within a single turn's worth of simulation
// segments.
//
// Rather than replicate the original's flat, offset-indexed byte array
// (DISENGAGE, SHIELDSF, PHASERS+i, TUBES+i, BURNOUT+id, SURRENDER), this
// uses named fields and small per-weapon/per-ship slices, which is
// equivalent in behavior but far less error-prone to extend.
type Shutup struct {
	Disengage          bool
	ShieldsFluctuating bool
	Surrender          bool

	// Phaser and Tube are indexed by weapon-bank index (0-based).
	Phaser [MaxWeapons]bool
	Tube   [MaxWeapons]bool

	// Burnout is indexed by ship ID, since the original's BURNOUT+id
	// offset applied to every ship, not just the player's.
	Burnout map[int]bool
}

// NewShutup returns a freshly reset Shutup, as if a new turn had just
// begun.
func NewShutup() *Shutup {
	return &Shutup{Burnout: make(map[int]bool)}
}

// Reset clears every suppression flag, mirroring the shutup[] reset
// loop at the top of playit()'s "next:" label in main.c.
func (s *Shutup) Reset() {
	*s = Shutup{Burnout: make(map[int]bool)}
}

// BurnoutShown reports whether the warp-burnout message has already
// been shown for the given ship this turn.
func (s *Shutup) BurnoutShown(shipID int) bool {
	return s.Burnout[shipID]
}

// SetBurnoutShown marks the warp-burnout message as shown for the
// given ship for the remainder of this turn.
func (s *Shutup) SetBurnoutShown(shipID int) {
	s.Burnout[shipID] = true
}
