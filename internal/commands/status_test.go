package commands

import (
	"strings"
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

func TestPhaserStatusReportsDamage(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	sp.Phasers[0].Status = game.PhaserDamaged
	lines := PhaserStatus(sp)
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "damaged") {
		t.Errorf("expected damaged bank listed, got: %s", joined)
	}
}

func TestTubeStatusReportsLaunchSettings(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	sp.TubeLaunchSpd = 5
	sp.TubeDelay = 3
	sp.TubeProximity = 200
	lines := TubeStatus(sp)
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "Launch speed: 5") {
		t.Errorf("missing launch speed, got: %s", joined)
	}
	if !strings.Contains(joined, "Time delay: 3") {
		t.Errorf("missing time delay, got: %s", joined)
	}
}

func TestSurvivorsReportsEachShipOwnStatus(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0), newTestShip(2, "Romulan", 200, 0))
	st.Enemies()[0].Complement = -1 // destroyed
	st.Enemies()[1].Cloaking = game.CloakOn
	lines := Survivors(st)
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "Klingon -- destructed") {
		t.Errorf("expected Klingon destructed (not Enterprise), got: %s", joined)
	}
	if !strings.Contains(joined, "Romulan -- ???") {
		t.Errorf("expected Romulan hidden, got: %s", joined)
	}
	if !strings.Contains(joined, "Enterprise -- 100") {
		t.Errorf("expected Enterprise survivor count, got: %s", joined)
	}
	// Regression: the original bug always printed the player's own
	// name for a destroyed ship, regardless of which ship was
	// actually destroyed.
	if strings.Contains(joined, "Enterprise -- destructed") {
		t.Errorf("bug regression: should not attribute destruction to Enterprise, got: %s", joined)
	}
}
