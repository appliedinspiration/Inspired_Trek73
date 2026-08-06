package commands

import (
	"strings"
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

func TestPosReportListsShipsAndSelf(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 1000, 0))
	sp := st.Player()
	lines := PosReport(st, sp)
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "Enterprise") {
		t.Errorf("missing own ship, got: %s", joined)
	}
	if !strings.Contains(joined, "Klingon") {
		t.Errorf("missing enemy ship, got: %s", joined)
	}
}

func TestPosReportShowsHelmLock(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 1000, 0))
	sp := st.Player()
	sp.Target = st.Enemies()[0]
	lines := PosReport(st, sp)
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "helm locked on Klingon") {
		t.Errorf("missing helm lock line, got: %s", joined)
	}
}

func TestPosReportSkipsDeadShips(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 1000, 0))
	dead := st.Enemies()[0]
	dead.Status[game.SysDead] = game.FullyDamaged
	lines := PosReport(st, st.Player())
	joined := strings.Join(lines, "\n")
	if strings.Contains(joined, "Klingon") {
		t.Errorf("should not list dead ship, got: %s", joined)
	}
}

func TestPosReportMarksCloakedWithLastKnownPosition(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 1000, 0))
	target := st.Enemies()[0]
	target.LastKnown = game.LastKnownPosition{X: 1000, Y: 0, Warp: 1, Course: 0}
	target.Cloaking = game.CloakOn
	lines := PosReport(st, st.Player())
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "Klingon *") {
		t.Errorf("expected cloaked marker, got: %s", joined)
	}
}

func TestPosDisplayRejectsOutOfRangeSensorValue(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	st := newTestState(sp)
	r := game.NewRand(1)
	lines := PosDisplay(sp, st, r, 10) // below MinSensorRange
	if lines != nil {
		t.Errorf("expected nil, got %v", lines)
	}
}

func TestPosDisplayRejectsDamagedSensor(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	sp.Status[game.SysSensor] = game.FullyDamaged
	st := newTestState(sp)
	r := game.NewRand(1)
	lines := PosDisplay(sp, st, r, 500)
	if len(lines) != 1 || !strings.Contains(lines[0], "damaged") {
		t.Errorf("expected damage message, got %v", lines)
	}
}

func TestPosDisplayPlotsNearbyShip(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	enemy := newTestShip(1, "Klingon", 100, 0)
	st := newTestState(sp, enemy)
	r := game.NewRand(1)
	lines := PosDisplay(sp, st, r, 500)
	found := false
	for _, l := range lines {
		if strings.ContainsRune(l, 'K') {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'K' plotted somewhere, got: %v", lines)
	}
}

func TestPosDisplayKeepsFullGridProportions(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	st := newTestState(sp)
	r := game.NewRand(1)
	lines := PosDisplay(sp, st, r, 500)
	if len(lines) != 13 {
		t.Fatalf("len(lines) = %d, want 13 full-height rows", len(lines))
	}
	if lines[0][0] != '-' || lines[0][len(lines[0])-1] != '-' {
		t.Fatalf("top border malformed: %q", lines[0])
	}
	if lines[1][0] != '|' || lines[1][len(lines[1])-1] != '|' {
		t.Fatalf("side borders malformed: %q", lines[1])
	}
}

func TestPosDisplayLowercasesCloakedShip(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	enemy := newTestShip(1, "Klingon", 100, 0)
	enemy.Cloaking = game.CloakOn
	enemy.LastKnown = game.LastKnownPosition{X: 100, Y: 0}
	st := newTestState(sp, enemy)
	r := game.NewRand(1)
	lines := PosDisplay(sp, st, r, 500)
	found := false
	for _, l := range lines {
		if strings.ContainsRune(l, 'k') {
			found = true
		}
	}
	if !found {
		t.Errorf("expected lowercase 'k' plotted for cloaked ship, got: %v", lines)
	}
}
