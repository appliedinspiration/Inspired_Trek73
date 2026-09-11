package game

import (
	"math/rand"
	"time"
)

// Rand wraps a seedable random source and provides the die-roll style
// helper the original game relies on throughout its mechanics.
//
// The original C code seeded a single global generator with
// srandom(time(0)) and read values through the randm(x) macro:
//
//	#define randm(x) ((int)random() % (x) + 1)
//
// That macro is slightly non-uniform whenever x does not evenly divide
// the generator's range, and relies on the platform's random()/srandom()
// implementation, which differs across operating systems. Inspired
// Trek73 replaces it with Go's math/rand, seeded explicitly for
// reproducibility, while preserving the same [1, x] inclusive range and
// otherwise-uniform distribution the original game's balance depends on.
//
// math/rand (v1) is used deliberately instead of math/rand/v2: its
// output for a given seed is stable across Go versions, which golden
// tests rely on for long-term reproducibility.
type Rand struct {
	seed int64
	used int64
	r    *rand.Rand
}

// RandSnapshot records the seed and number of draws consumed so the
// generator can be resumed after save/restore without losing the
// exact RNG state.
type RandSnapshot struct {
	Seed int64
	Used int64
}

// Snapshot returns the current RNG state for serialization.
func (rr *Rand) Snapshot() RandSnapshot {
	return RandSnapshot{Seed: rr.seed, Used: rr.used}
}

// Restore rehydrates the RNG from a prior snapshot, advancing the
// generator the same number of draws that had previously been consumed.
func (rr *Rand) Restore(snapshot RandSnapshot) {
	rr.seed = snapshot.Seed
	rr.used = snapshot.Used
	rr.r = rand.New(rand.NewSource(snapshot.Seed))
	for i := int64(0); i < snapshot.Used; i++ {
		rr.r.Int63()
	}
}

// NewRand returns a Rand seeded deterministically, e.g. for tests or a
// reproducible scenario.
func NewRand(seed int64) *Rand {
	return &Rand{seed: seed, r: rand.New(rand.NewSource(seed))}
}

// NewRandFromTime returns a Rand seeded from the current time, matching
// the original game's srandom(time(0)) startup behavior for normal play.
func NewRandFromTime() *Rand {
	return NewRand(time.Now().UnixNano())
}

// Randm returns a pseudo-random integer in the inclusive range [1, n],
// corresponding to the original game's randm(x) macro. It panics if n
// is less than 1, since the original macro's behavior for non-positive
// arguments is undefined and no legitimate game mechanic calls it that
// way.
func (rr *Rand) Randm(n int) int {
	if n < 1 {
		panic("game: Randm called with non-positive n")
	}
	rr.used++
	return rr.r.Intn(n) + 1
}
