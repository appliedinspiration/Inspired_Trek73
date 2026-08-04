package game

import "testing"

// TestRandmRange verifies every draw stays within the documented
// inclusive [1, n] range across a range of n values.
func TestRandmRange(t *testing.T) {
	r := NewRand(1)
	for _, n := range []int{1, 2, 6, 10, 100, 360, 1000} {
		for i := 0; i < 2000; i++ {
			v := r.Randm(n)
			if v < 1 || v > n {
				t.Fatalf("Randm(%d) = %d, want value in [1, %d]", n, v, n)
			}
		}
	}
}

// TestRandmDeterministic verifies that two Rand instances seeded
// identically produce identical sequences, which the golden-test
// strategy depends on for reproducible scenarios.
func TestRandmDeterministic(t *testing.T) {
	const seed = 42

	a := NewRand(seed)
	b := NewRand(seed)

	for i := 0; i < 500; i++ {
		va := a.Randm(360)
		vb := b.Randm(360)
		if va != vb {
			t.Fatalf("draw %d: seeded sequences diverged: %d != %d", i, va, vb)
		}
	}
}

// TestRandmGoldenSequence locks in the exact sequence produced by a
// fixed seed. This is the golden-test baseline for Inspired Trek73's
// own RNG behavior (there is no equivalent C sequence to compare
// against, since the original random()/srandom() algorithm is
// platform-specific). If this test ever fails after a deliberate change
// to Rand, the recorded sequence below must be regenerated and the
// change called out explicitly, since it will shift every downstream
// mechanic that consumes randomness.
func TestRandmGoldenSequence(t *testing.T) {
	r := NewRand(20260803)

	want := []int{
		5, 7, 3, 2, 9, 1, 3, 8, 3, 9,
		3, 42, 107, 76, 139, 135, 269, 220, 105, 3,
	}

	got := make([]int, 0, len(want))
	for i := 0; i < 10; i++ {
		got = append(got, r.Randm(9))
	}
	for i := 0; i < 10; i++ {
		got = append(got, r.Randm(360))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("draw %d = %d, want %d (full sequence: %v)", i, got[i], want[i], got)
		}
	}
}
