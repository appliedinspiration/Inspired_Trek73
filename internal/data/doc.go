// Package data holds the static game-content tables ported from the
// original Trek73 source: ship class stats, alien race information,
// Federation ship names, damage tables, and weapon initial-bearing
// tables. See References/FreeBSD/trek73/src/ships.c and globals.c for
// the source material.
//
// These tables are treated as authoritative game content: values are
// preserved exactly as in the original, including entries (such as the
// "Monty Python" alien race and its associated names) that later became
// optional via the original game's "silly" flag. Whether a given table
// entry is offered to the player is a game-logic decision made
// elsewhere, not a data-table concern.
package data
