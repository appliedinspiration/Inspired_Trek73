package game

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveGameRoundTrip(t *testing.T) {
	seed := int64(12345)
	r := NewRand(seed)
	res, err := NewGame(InitOptions{EnemyCount: 2}, r)
	if err != nil {
		t.Fatalf("NewGame() error = %v", err)
	}
	res.State.Crew = CrewNames{
		Captain:  "Kirk",
		Science:  "Spock",
		Engineer: "Scott",
		Com:      "Uhura",
		Nav:      "Chekov",
		Helmsman: "Sulu",
		Title:    "Captain",
	}
	for i := 0; i < 12; i++ {
		r.Randm(100)
	}

	path := filepath.Join(t.TempDir(), "save.json")
	if err := SaveGame(path, res.State, r); err != nil {
		t.Fatalf("SaveGame() error = %v", err)
	}

	loadedState, loadedRand, err := LoadGame(path)
	if err != nil {
		t.Fatalf("LoadGame() error = %v", err)
	}

	if loadedState.Crew != res.State.Crew {
		t.Fatalf("loaded crew mismatch: got %#v, want %#v", loadedState.Crew, res.State.Crew)
	}
	if loadedState.Player().Name != res.State.Player().Name {
		t.Fatalf("player name mismatch: got %q, want %q", loadedState.Player().Name, res.State.Player().Name)
	}
	if len(loadedState.Ships) != len(res.State.Ships) {
		t.Fatalf("ship count mismatch: got %d, want %d", len(loadedState.Ships), len(res.State.Ships))
	}

	for i := 0; i < 25; i++ {
		wantVal := r.Randm(100)
		got := loadedRand.Randm(100)
		if got != wantVal {
			t.Fatalf("RNG mismatch at draw %d: got %d, want %d", i, got, wantVal)
		}
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("saved file not created: %v", err)
	}
}
