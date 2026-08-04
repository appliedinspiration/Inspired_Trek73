package data

// PhaserInitialBearings gives the initial relative bearing (in degrees)
// assigned to each phaser bank when a ship is created, indexed by the
// ship's number of phasers and then by bank number. Ported verbatim
// from init_p_turn in globals.c.
//
// Row 0 is an unused sentinel from the original table (no ship has zero
// phasers); rows 1-10 have exactly as many entries as the phaser count
// they describe.
var PhaserInitialBearings = [][]float64{
	{666.666},
	{0.0},
	{0.0, 0.0},
	{90.0, 0.0, 270.0},
	{90.0, 0.0, 0.0, 270.0},
	{90.0, 0.0, 0.0, 0.0, 270.0},
	{90.0, 90.0, 0.0, 0.0, 270.0, 270.0},
	{90.0, 90.0, 0.0, 0.0, 0.0, 270.0, 270.0},
	{90.0, 90.0, 0.0, 0.0, 0.0, 0.0, 270.0, 270.0},
	{90.0, 90.0, 90.0, 0.0, 0.0, 0.0, 270.0, 270.0, 270.0},
	{90.0, 90.0, 90.0, 0.0, 0.0, 0.0, 0.0, 270.0, 270.0, 270.0},
}

// TubeInitialBearings gives the initial relative bearing (in degrees)
// assigned to each torpedo tube when a ship is created, indexed by the
// ship's number of tubes and then by tube number. Ported verbatim from
// init_t_turn in globals.c.
//
// Row 0 is an unused sentinel from the original table (no ship has zero
// tubes); rows 1-10 have exactly as many entries as the tube count they
// describe.
var TubeInitialBearings = [][]float64{
	{666.666},
	{0.0},
	{0.0, 0.0},
	{60.0, 0.0, 300.0},
	{60.0, 0.0, 0.0, 300.0},
	{60.0, 0.0, 0.0, 0.0, 300.0},
	{120.0, 60.0, 0.0, 0.0, 300.0, 240.0},
	{120.0, 60.0, 0.0, 0.0, 0.0, 300.0, 240.0},
	{120.0, 60.0, 60.0, 0.0, 0.0, 300.0, 300.0, 240.0},
	{120.0, 60.0, 60.0, 0.0, 0.0, 0.0, 300.0, 300.0, 240.0},
	{120.0, 120.0, 60.0, 60.0, 0.0, 0.0, 300.0, 300.0, 240.0, 240.0},
}
