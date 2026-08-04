package commands

import (
	"strings"
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

func TestPlayDead_TransferEngines(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	sp.Energy = 500
	sp.Pods = 2000
	st := newTestState(sp)
	st.EnemyRaceName = "Klingon"

	msgs := PlayDead(st, sp, TransferEngines)
	if st.Defenseless != 1 {
		t.Errorf("Defenseless = %v, want 1", st.Defenseless)
	}
	if sp.Energy != 500-game.MaxPhaserCharge {
		t.Errorf("Energy = %v, want %v", sp.Energy, 500-game.MaxPhaserCharge)
	}
	for i, ph := range sp.Phasers {
		if ph.Drain != int(-game.MaxPhaserCharge) {
			t.Errorf("phaser %d drain = %v, want %v", i, ph.Drain, -game.MaxPhaserCharge)
		}
	}
	if len(msgs) == 0 {
		t.Error("expected messages")
	}
}

func TestPlayDead_TransferPhasers(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)

	PlayDead(st, sp, TransferPhasers)
	for i, ph := range sp.Phasers {
		if ph.Drain != int(game.MaxPhaserCharge) {
			t.Errorf("phaser %d drain = %v, want %v", i, ph.Drain, game.MaxPhaserCharge)
		}
	}
}

func TestPlayDead_AlreadyDefenseless(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	st.Defenseless = 5

	msgs := PlayDead(st, sp, TransferEngines)
	if st.Defenseless != 5 {
		t.Errorf("Defenseless should be unchanged, got %v", st.Defenseless)
	}
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "not that stupid") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'not that stupid' message, got %v", msgs)
	}
}

func TestPlayDead_InvalidTarget(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)

	msgs := PlayDead(st, sp, TransferNone)
	if st.Defenseless != 0 {
		t.Errorf("Defenseless should remain 0, got %v", st.Defenseless)
	}
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "cannot transfer power") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'cannot transfer power' message, got %v", msgs)
	}
}

func TestCorbomiteBluff_FirstAttempt(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)

	msgs := CorbomiteBluff(st, sp, true)
	if st.Corbomite != 1 {
		t.Errorf("Corbomite = %v, want 1", st.Corbomite)
	}
	if len(msgs) == 0 {
		t.Error("expected messages")
	}
}

func TestCorbomiteBluff_SecondAttempt(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	st.Corbomite = 1

	msgs := CorbomiteBluff(st, sp, false)
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "don't believe") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected skeptical response on repeat bluff, got %v", msgs)
	}
}

func TestSurrenderShip_AlreadySurrendered(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	st.PlayStatus |= game.StatusFedSurrender

	msgs := SurrenderShip(st, sp)
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "already surrendered") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected already-surrendered message, got %v", msgs)
	}
}

func TestSurrenderShip_RomulanBugFixed(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	st.EnemyRaceName = "Romulan"

	msgs := SurrenderShip(st, sp)
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "not been know to have taken") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected Romulan prisoner-taking message (bug fix), got %v", msgs)
	}
	if st.Surrender != 1 {
		t.Errorf("Surrender = %v, want 1", st.Surrender)
	}
}

func TestSurrenderShip_AlreadyRefused(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	st.Surrender = 1
	st.EnemyRaceName = "Klingon"

	msgs := SurrenderShip(st, sp)
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "already refused") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected already-refused message, got %v", msgs)
	}
}

func TestRequestSurrender_AlreadyComplying(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	st.PlayStatus |= game.StatusEnemySurrender

	msgs := RequestSurrender(st, sp)
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "already complying") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected already-complying message, got %v", msgs)
	}
}

func TestRequestSurrender_AlreadyRefused(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	st.SurrenderP = 1

	msgs := RequestSurrender(st, sp)
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "already been refused") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected already-refused message, got %v", msgs)
	}
}

func TestRequestSurrender_FirstRequest(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)

	RequestSurrender(st, sp)
	if st.SurrenderP != 1 {
		t.Errorf("SurrenderP = %v, want 1", st.SurrenderP)
	}
}

func TestSelfDestruct_ComputerDead(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	sp.Status[game.SysComputer] = game.FullyDamaged
	r := game.NewRand(1)

	msgs := SelfDestruct(st, r, sp)
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "computer is down") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected computer-down message, got %v", msgs)
	}
	if sp.Delay == SelfDestructDelay {
		t.Error("Delay should not have been set when computer is dead")
	}
}

func TestSelfDestruct_Success(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	st := newTestState(sp)
	r := game.NewRand(1)

	// Find a seed roll where SystemWorks succeeds (Status[SysComputer]==0
	// means randm(100) > 0 is true unless randm returns exactly 0).
	msgs := SelfDestruct(st, r, sp)
	if sp.Delay != SelfDestructDelay {
		// System might have "failed" this once; retry with a fresh roll.
		msgs = SelfDestruct(st, r, sp)
	}
	if sp.Delay != SelfDestructDelay {
		t.Fatalf("expected Delay = %v after successful self-destruct sequence, got %v (msgs=%v)", SelfDestructDelay, sp.Delay, msgs)
	}
}

func TestAbortSelfDestruct_NotInitiated(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	sp.Delay = 10000.0
	st := newTestState(sp)
	r := game.NewRand(1)

	msgs := AbortSelfDestruct(st, r, sp)
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "has not been") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected not-initiated message, got %v", msgs)
	}
}

func TestAbortSelfDestruct_AbortedSuccessfully(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	sp.Delay = SelfDestructDelay
	st := newTestState(sp)
	r := game.NewRand(1)

	msgs := AbortSelfDestruct(st, r, sp)
	if sp.Delay != 10000.0 {
		t.Errorf("Delay = %v, want 10000 after abort", sp.Delay)
	}
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "aborted.  Destruct order aborted.") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected abort-success message, got %v", msgs)
	}
}

func TestAbortSelfDestruct_TooLateToAbort(t *testing.T) {
	sp := newTestShip(0, "Player", 0, 0)
	sp.Delay = 2.0 // less than 4, past the abort window
	st := newTestState(sp)
	r := game.NewRand(1)

	msgs := AbortSelfDestruct(st, r, sp)
	if sp.Delay != 2.0 {
		t.Errorf("Delay should be unchanged, got %v", sp.Delay)
	}
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "cannot be aborted") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected cannot-be-aborted message, got %v", msgs)
	}
}
