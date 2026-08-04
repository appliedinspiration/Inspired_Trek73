package game

import "math"

// PhaserHit computes the raw damage a phaser bank would inflict on a
// target at (x,y), given the bank's current load and the ship's
// spread/firing-percentage settings, ported from phaser_hit() in
// subs.c.
//
// trueBear is the absolute bearing the phaser bank is actually pointed
// (as opposed to the bearing of the target), used to check whether
// the target falls within the bank's spread cone.
func PhaserHit(sp *Ship, x, y int, bank *Phaser, trueBear float64) int {
	rng := rangeFind(sp.X, x, sp.Y, y)
	if rng >= MaxPhaserRange {
		return 0
	}
	bear := bearingTo(sp.X, x, sp.Y, y)
	spread := math.Min(rectify(bear-trueBear), rectify(trueBear-bear))
	if spread > float64(sp.PhaserSpread) {
		return 0
	}
	d1 := 1.0 - float64(rng)/MaxPhaserRange
	pct := float64(sp.PhaserFirePct) / 100
	d2 := bank.Load * math.Sqrt(d1) * pct
	// This may have to be changed if phaser spread or maximum phaser
	// load is changed.
	d2 = bank.Load * d2 * 45.0 / float64(sp.PhaserSpread) * pct
	return int(d2 / 10.0)
}

// TorpedoHit computes the raw damage an antimatter blast of the given
// fuel yield would inflict at range from (x,y) to (tx,ty), ported from
// torpedo_hit() in subs.c.
func TorpedoHit(fuel, x, y, tx, ty int) int {
	rng := rangeFind(x, tx, y, ty)
	f1 := float64(fuel * HitPerPod)
	f2 := float64(fuel * ProxPerPod)
	if float64(rng) >= f2 {
		return 0
	}
	d1 := 1.0 - float64(rng)/f2
	return int(f1 * math.Sqrt(d1))
}

// shieldFacingFromAntimatterBearing returns the 1-based shield facing
// (matching the original's s=1..4 numbering used by Damage) struck by
// an antimatter blast arriving from the given bearing (relative to the
// target's own course), ported from the facing-selection logic in
// antimatter_hit() (subs.c).
func shieldFacingFromAntimatterBearing(bear float64) int {
	switch {
	case bear <= 45.0 || bear >= 315.0:
		return 1
	case bear <= 135.0:
		return 2
	case bear < 225.0:
		return 3
	default:
		return 4
	}
}

// shieldFacingFromPhaserBearing returns the 1-based shield facing
// struck by a phaser hit arriving from the given bearing (relative to
// the target's own course), ported from the facing-selection logic in
// phaser_firing() (firing.c).
//
// Note the boundary comparisons (> / < rather than antimatter's >= /
// <=) differ slightly from shieldFacingFromAntimatterBearing; this
// mirrors the original, where the two facing calculations were
// written independently and never unified. It is preserved as-is
// rather than "fixed", since the discrepancy only affects the exact
// 45/135/225/315-degree boundary values and was not flagged as an
// unintended defect.
func shieldFacingFromPhaserBearing(bear float64) int {
	switch {
	case bear > 315.0 || bear < 45.0:
		return 1
	case bear < 135.0:
		return 2
	case bear < 225.0:
		return 3
	default:
		return 4
	}
}
