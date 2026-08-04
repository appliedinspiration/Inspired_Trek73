package game

import (
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/data"
)

func TestNewGameRejectsOutOfRangeEnemyCount(t *testing.T) {
	r := NewRand(1)
	if _, err := NewGame(InitOptions{EnemyCount: 0}, r); err != nil {
		t.Errorf("EnemyCount=0 (random) should be valid, got error: %v", err)
	}
	for _, n := range []int{-1, 10, 99} {
		if _, err := NewGame(InitOptions{EnemyCount: n}, r); err == nil {
			t.Errorf("EnemyCount=%d should be rejected, got no error", n)
		}
	}
}

func TestNewGameDefaults(t *testing.T) {
	r := NewRand(7)
	res, err := NewGame(InitOptions{EnemyCount: 3}, r)
	if err != nil {
		t.Fatalf("NewGame: %v", err)
	}
	st := res.State

	if got := len(st.Ships); got != 4 {
		t.Fatalf("len(Ships) = %d, want 4 (1 player + 3 enemies)", got)
	}

	player := st.Player()
	if player.Class != "CA" {
		t.Errorf("player.Class = %q, want CA (default)", player.Class)
	}
	if player.X != 0 || player.Y != 0 {
		t.Errorf("player position = (%d, %d), want (0, 0)", player.X, player.Y)
	}
	if player.Course != 0 || player.NewCourse != 0 {
		t.Errorf("player course = %v/%v, want 0/0", player.Course, player.NewCourse)
	}
	if player.Cloaking != CloakNone {
		t.Errorf("player.Cloaking = %d, want CloakNone", player.Cloaking)
	}
	playerClass, _ := data.ShipClassByAbbr("CA")
	if len(player.Phasers) != playerClass.NumPhaser {
		t.Errorf("len(player.Phasers) = %d, want %d", len(player.Phasers), playerClass.NumPhaser)
	}
	if len(player.Tubes) != playerClass.NumTorp {
		t.Errorf("len(player.Tubes) = %d, want %d", len(player.Tubes), playerClass.NumTorp)
	}
	if player.Complement != playerClass.OwnCrew {
		t.Errorf("player.Complement = %d, want %d", player.Complement, playerClass.OwnCrew)
	}
	if player.Name == "" {
		t.Error("player.Name is empty, want a random Federation name")
	}

	enemies := st.Enemies()
	if len(enemies) != 3 {
		t.Fatalf("len(Enemies()) = %d, want 3", len(enemies))
	}
	seenNames := map[string]bool{}
	for i, e := range enemies {
		if e.Name == "" {
			t.Errorf("enemy %d has empty name", i)
		}
		if seenNames[e.Name] {
			t.Errorf("enemy name %q assigned more than once", e.Name)
		}
		seenNames[e.Name] = true
		if e.Class != "CA" {
			t.Errorf("enemy %d class = %q, want CA (default)", i, e.Class)
		}
		if e.X == 0 && e.Y == 0 {
			t.Errorf("enemy %d sits at the origin, expected to be positioned away from the player", i)
		}
		for _, ph := range e.Phasers {
			if ph.Load != InitPhaserLoad {
				t.Errorf("enemy %d phaser load = %v, want %v", i, ph.Load, InitPhaserLoad)
			}
		}
	}

	if res.EnemyRaceName == "" {
		t.Error("EnemyRaceName is empty")
	}
	if res.EnemyCommander == "" {
		t.Error("EnemyCommander is empty")
	}
}

func TestNewGameExplicitRace(t *testing.T) {
	r := NewRand(3)
	res, err := NewGame(InitOptions{EnemyCount: 1, RaceName: "Klingon"}, r)
	if err != nil {
		t.Fatalf("NewGame: %v", err)
	}
	if res.EnemyRaceName != "Klingon" {
		t.Errorf("EnemyRaceName = %q, want Klingon", res.EnemyRaceName)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("unexpected warnings for a valid race name: %v", res.Warnings)
	}
}

func TestNewGameUnknownRaceFallsBackWithWarning(t *testing.T) {
	r := NewRand(3)
	res, err := NewGame(InitOptions{EnemyCount: 1, RaceName: "Ferengi"}, r)
	if err != nil {
		t.Fatalf("NewGame: %v", err)
	}
	if res.EnemyRaceName == "" {
		t.Error("expected a fallback race to be chosen")
	}
	if len(res.Warnings) != 1 {
		t.Fatalf("Warnings = %v, want exactly one warning", res.Warnings)
	}
}

// TestNewGameExcludesMontyPythonByDefault verifies that, over many
// random draws without AllowSillyRace, the joke race is never chosen,
// matching the original "silly" flag behavior.
func TestNewGameExcludesMontyPythonByDefault(t *testing.T) {
	r := NewRand(99)
	for i := 0; i < 500; i++ {
		res, err := NewGame(InitOptions{EnemyCount: 1}, r)
		if err != nil {
			t.Fatalf("NewGame: %v", err)
		}
		if res.EnemyRaceName == "Monty Python" {
			t.Fatalf("draw %d: Monty Python chosen despite AllowSillyRace=false", i)
		}
	}
}

func TestNewGameUnknownClassFallsBackWithWarning(t *testing.T) {
	r := NewRand(5)
	res, err := NewGame(InitOptions{EnemyCount: 1, PlayerClassAbbr: "ZZ"}, r)
	if err != nil {
		t.Fatalf("NewGame: %v", err)
	}
	if res.State.Player().Class != "CA" {
		t.Errorf("player.Class = %q, want fallback CA", res.State.Player().Class)
	}
	if len(res.Warnings) != 1 {
		t.Fatalf("Warnings = %v, want exactly one warning", res.Warnings)
	}
}

// TestNewGameGoldenScenario locks in the exact ship placement produced
// by a fixed seed and fixed options. This is the golden baseline for
// Inspired Trek73's own initialization behavior (there is no equivalent
// C output to compare against, since the original random()/srandom()
// algorithm is platform-specific and differs from Go's). If this test
// fails after a deliberate change to NewGame or Rand, the recorded
// values below must be regenerated and the change called out
// explicitly.
func TestNewGameGoldenScenario(t *testing.T) {
	r := NewRand(20260803)
	res, err := NewGame(InitOptions{EnemyCount: 2}, r)
	if err != nil {
		t.Fatalf("NewGame: %v", err)
	}
	st := res.State

	if got, want := len(st.Ships), 3; got != want {
		t.Fatalf("len(Ships) = %d, want %d", got, want)
	}

	wantRace := "Lyran"
	if res.EnemyRaceName != wantRace {
		t.Errorf("EnemyRaceName = %q, want %q", res.EnemyRaceName, wantRace)
	}

	wantEnemyPositions := []struct{ x, y int }{
		{-3764, -1370},
		{66, 3835},
	}
	for i, e := range st.Enemies() {
		if e.X != wantEnemyPositions[i].x || e.Y != wantEnemyPositions[i].y {
			t.Errorf("enemy %d position = (%d, %d), want (%d, %d)",
				i, e.X, e.Y, wantEnemyPositions[i].x, wantEnemyPositions[i].y)
		}
	}
}
