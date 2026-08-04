package commands

import "github.com/appliedinspiration/inspired_trek73/internal/game"

// newTestShip builds a minimal, fully-populated Ship for command unit
// tests, mirroring internal/game's own newTestShip helper (unexported
// there, so duplicated here rather than exported purely for tests).
func newTestShip(id int, name string, x, y int) *game.Ship {
	s := &game.Ship{
		ID:                 id,
		Name:               name,
		Class:              "CA",
		X:                  x,
		Y:                  y,
		Warp:               1.0,
		NewWarp:            1.0,
		Eff:                1.0,
		Energy:             1000,
		Pods:               2000,
		Complement:         100,
		Delay:              10000,
		OrigMaxSpeed:       8,
		MaxSpeed:           8,
		DegPerTurn:         10,
		Phasers:            make([]game.Phaser, 3),
		Tubes:              make([]game.Tube, 3),
		PhaserFiringDelay:  1,
		TorpedoFiringDelay: 1,
		Cloaking:           game.CloakNone,
	}
	for i := range s.Shields {
		s.Shields[i] = game.Shield{Eff: 1.0}
	}
	return s
}

// newTestState builds a State with the given ships (index 0 is always
// the player, matching game.State's convention).
func newTestState(ships ...*game.Ship) *game.State {
	st := game.NewState()
	st.Ships = append(st.Ships, ships...)
	return st
}
