package game

import "math"

// rectify normalizes an angle in degrees into the range [0, 360),
// ported from rectify() in subs.c.
func rectify(deg float64) float64 {
	deg = math.Mod(deg, 360.0)
	if deg < 0 {
		deg += 360.0
	}
	return deg
}

// toDegrees converts radians to degrees, matching the todegrees(x)
// macro in defines.h. Like toRadians, this uses a precise conversion
// rather than the original's truncated constant (57.29578); a
// deliberate accuracy improvement with no meaningful gameplay impact.
func toDegrees(rad float64) float64 {
	return rad * (180 / math.Pi)
}

// rangeFind returns the integer distance between two points, ported
// from rangefind() in subs.c.
func rangeFind(xFrom, xTo, yFrom, yTo int) int {
	x := float64(xTo - xFrom)
	y := float64(yTo - yFrom)
	if x == 0 && y == 0 {
		return 0
	}
	return int(math.Hypot(x, y))
}

// bearingTo returns the bearing, in degrees [0, 360), of (xTo,yTo) as
// seen from (xFrom,yFrom), ported from bearing() in subs.c.
func bearingTo(xFrom, xTo, yFrom, yTo int) float64 {
	x := float64(xTo - xFrom)
	y := float64(yTo - yFrom)
	if x == 0 && y == 0 {
		return 0
	}
	return rectify(toDegrees(math.Atan2(y, x)))
}

// canSee reports whether a ship is visible to enemy sensors (i.e. is
// not currently cloaked), ported from the cansee(x) macro in
// defines.h.
func canSee(sp *Ship) bool {
	return sp.Cloaking != CloakOn
}

// cantSee reports whether a ship is cloaked and thus invisible to
// enemy sensors, ported from the cantsee(x) macro in defines.h.
func cantSee(sp *Ship) bool {
	return sp.Cloaking == CloakOn
}

// RangeFind, BearingTo, Rectify, CanSee, and CantSee are exported
// wrappers around the like-named unexported helpers above, for use by
// packages outside game (e.g. internal/commands) that need to display
// range/bearing information the same way the engine computes it.
func RangeFind(xFrom, xTo, yFrom, yTo int) int     { return rangeFind(xFrom, xTo, yFrom, yTo) }
func BearingTo(xFrom, xTo, yFrom, yTo int) float64 { return bearingTo(xFrom, xTo, yFrom, yTo) }
func Rectify(deg float64) float64                  { return rectify(deg) }
func CanSee(sp *Ship) bool                         { return canSee(sp) }
func CantSee(sp *Ship) bool                        { return cantSee(sp) }
