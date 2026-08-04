// Package game contains the core Inspired Trek73 game engine: ship and
// object state, the simulation tick, combat resolution, damage handling,
// enemy strategy, and victory/defeat conditions.
//
// This package is being built incrementally following the project's
// golden-test strategy: each mechanic ported from the original C
// implementation (see References/FreeBSD/trek73) is accompanied by
// deterministic tests before it is considered complete. Until that work
// lands, this package intentionally contains only scaffolding.
package game
