package commands

import (
	"strings"
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

func TestPrintDamageReportsSurvivorsAndSystems(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	sp.Status[game.SysComputer] = 40
	lines := PrintDamage(sp)
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "Damages to the Enterprise") {
		t.Errorf("missing header, got: %s", joined)
	}
	if !strings.Contains(joined, "Computer damaged 40%") {
		t.Errorf("missing computer damage line, got: %s", joined)
	}
	if !strings.Contains(joined, "Survivors: 100") {
		t.Errorf("missing survivors line, got: %s", joined)
	}
}

func TestPrintDamageReportsDestroyedSystem(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	sp.Status[game.SysWarp] = game.FullyDamaged
	lines := PrintDamage(sp)
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "warp drive disabled.") {
		t.Errorf("missing destroyed message, got: %s", joined)
	}
}

func TestSelfScanIsPrintDamage(t *testing.T) {
	sp := newTestShip(0, "Enterprise", 0, 0)
	got := SelfScan(sp)
	want := PrintDamage(sp)
	if len(got) != len(want) {
		t.Errorf("len = %d, want %d", len(got), len(want))
	}
}

func TestScanTorpedoRendersProbe(t *testing.T) {
	sp := newTestShip(1, "Klingon", 100, 0)
	obj := &game.SpaceObject{ID: 3, Type: game.ObjectProbe, Target: sp, TimeDelay: 5, Proximity: 100, Fuel: 20}
	lines := ScanTorpedo(obj)
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "probe") || !strings.Contains(joined, "Klingon") {
		t.Errorf("unexpected output: %s", joined)
	}
}

func TestScanShipSucceeds(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	r := game.NewRand(1)
	result, msgs := Scan(st, r, sp, "Klingon")
	if result.Ship == nil {
		t.Fatalf("expected resolved ship, messages: %v", msgs)
	}
}

func TestScanRefusesCloakedShip(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	target := st.Enemies()[0]
	target.Cloaking = game.CloakOn
	r := game.NewRand(1)
	result, msgs := Scan(st, r, sp, "Klingon")
	if result.Ship != nil {
		t.Error("should not resolve a cloaked ship")
	}
	if len(msgs) == 0 {
		t.Error("expected a message")
	}
}

func TestScanRefusesScanningSelf(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	r := game.NewRand(1)
	result, msgs := Scan(st, r, sp, "Enterprise")
	if result.Ship != nil {
		t.Error("should not resolve self-scan as a ship scan")
	}
	if len(msgs) == 0 {
		t.Error("expected a redirect-to-damage-report message")
	}
}

func TestScanEngineeringSection(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0), newTestShip(1, "Klingon", 100, 0))
	sp := st.Player()
	ep := st.Enemies()[0]
	obj := &game.SpaceObject{Type: game.ObjectEngineering, From: ep, Fuel: 50}
	st.Objects = append(st.Objects, obj)
	r := game.NewRand(1)
	result, msgs := Scan(st, r, sp, "#Klingon")
	if result.Object != obj {
		t.Fatalf("expected resolved engineering object, messages: %v", msgs)
	}
}

func TestScanProbeByNumber(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	obj := &game.SpaceObject{ID: 5, Type: game.ObjectProbe}
	st.Objects = append(st.Objects, obj)
	r := game.NewRand(1)
	result, msgs := Scan(st, r, sp, "5")
	if result.Object != obj {
		t.Fatalf("expected resolved probe, messages: %v", msgs)
	}
}

func TestScanRefusesUnknownShip(t *testing.T) {
	st := newTestState(newTestShip(0, "Enterprise", 0, 0))
	sp := st.Player()
	r := game.NewRand(1)
	result, msgs := Scan(st, r, sp, "Romulan")
	if result.Ship != nil || result.Object != nil {
		t.Error("expected no resolution")
	}
	if len(msgs) == 0 {
		t.Error("expected a message")
	}
}
