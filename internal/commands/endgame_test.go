package commands

import (
	"strings"
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

func newEndgameState() *game.State {
	fed := newTestShip(1, "Enterprise", 0, 0)
	enemy := newTestShip(2, "Klathis", 1, 1)
	st := newTestState(fed, enemy)
	st.EnemyRaceName = "Klingon"
	st.EnemyEmpireName = "Empire"
	st.EnemyShipTypeName = "Battlecruiser"
	st.EnemyCommander = "Kor"
	return st
}

func TestPlural(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, ""},
		{1, ""},
		{2, "s"},
		{5, "s"},
	}
	for _, tc := range tests {
		if got := plural(tc.n); got != tc.want {
			t.Errorf("plural(%d) = %q, want %q", tc.n, got, tc.want)
		}
	}
}

func TestVowelStr(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"Andorian", "n"},
		{"orion", "n"},
		{"Klingon", ""},
		{"", ""},
	}
	for _, tc := range tests {
		if got := vowelStr(tc.s); got != tc.want {
			t.Errorf("vowelStr(%q) = %q, want %q", tc.s, got, tc.want)
		}
	}
}

func TestStarfleet(t *testing.T) {
	msgs := Starfleet()
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "Starfleet Command") {
			found = true
		}
	}
	if !found {
		t.Error("Starfleet() missing header line")
	}
}

func TestFinal_FedLose(t *testing.T) {
	st := newEndgameState()
	msgs := Final(st, game.FinFedLose, false)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "destroyed by") {
		t.Errorf("Final(FinFedLose) missing destruction line: %v", msgs)
	}
	if !strings.Contains(joined, "Klingon") {
		t.Errorf("Final(FinFedLose) missing race name: %v", msgs)
	}
}

func TestFinal_EnemyLose(t *testing.T) {
	st := newEndgameState()
	msgs := Final(st, game.FinEnemyLose, false)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "fine performance") {
		t.Errorf("Final(FinEnemyLose) missing commendation: %v", msgs)
	}
}

func TestFinal_Tactical_PromptsReengage(t *testing.T) {
	st := newEndgameState()
	msgs := Final(st, game.FinTactical, true)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "re-engage") {
		t.Errorf("Final(FinTactical, true) should prompt to re-engage: %v", msgs)
	}
	if !st.Reengaged {
		t.Error("Final(FinTactical, true) should set st.Reengaged = true")
	}
}

func TestFinal_Tactical_DeclineReengage(t *testing.T) {
	st := newEndgameState()
	msgs := Final(st, game.FinTactical, false)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "out-maneuvered") {
		t.Errorf("Final(FinTactical, false) should render final tactical message: %v", msgs)
	}
	if st.Reengaged {
		t.Error("Final(FinTactical, false) should not set st.Reengaged")
	}
}

func TestFinal_Tactical_AlreadyReengaged(t *testing.T) {
	st := newEndgameState()
	st.Reengaged = true
	msgs := Final(st, game.FinTactical, true)
	joined := strings.Join(msgs, "\n")
	if strings.Contains(joined, "Do you wish to re-engage?") {
		t.Errorf("Final(FinTactical) should not re-prompt once already reengaged: %v", msgs)
	}
	if !strings.Contains(joined, "out-maneuvered") {
		t.Errorf("Final(FinTactical) should fall through to final message: %v", msgs)
	}
}

func TestFinal_FedSurrender(t *testing.T) {
	st := newEndgameState()
	msgs := Final(st, game.FinFedSurrender, false)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "surrendered") {
		t.Errorf("Final(FinFedSurrender) missing surrender text: %v", msgs)
	}
}

func TestFinal_EnemySurrender(t *testing.T) {
	st := newEndgameState()
	msgs := Final(st, game.FinEnemySurrender, false)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "surrendered") {
		t.Errorf("Final(FinEnemySurrender) missing surrender text: %v", msgs)
	}
}

func TestFinal_Complete(t *testing.T) {
	st := newEndgameState()
	msgs := Final(st, game.FinComplete, false)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "wreckage") {
		t.Errorf("Final(FinComplete) missing wreckage line: %v", msgs)
	}
}

func TestFinal_SurvivorsReported(t *testing.T) {
	st := newEndgameState()
	msgs := Final(st, game.FinEnemyLose, false)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "Survivors Reported") {
		t.Errorf("Final() should report survivors when ships remain alive: %v", msgs)
	}
}

func TestFinal_ListsDestroyedShipsInSurvivorsReport(t *testing.T) {
	st := newEndgameState()
	enemy := st.Enemies()[0]
	enemy.Complement = 0
	enemy.Status[game.SysDead] = game.FullyDamaged

	msgs := Final(st, game.FinEnemyLose, false)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "Klathis -- destroyed") {
		t.Errorf("Final() should list destroyed ships in the survivors report: %v", msgs)
	}
}

func TestFinal_NoSurvivors(t *testing.T) {
	st := newEndgameState()
	for _, sp := range st.Ships {
		sp.Complement = 0
		sp.Status[game.SysDead] = game.FullyDamaged
	}
	msgs := Final(st, game.FinEnemyLose, false)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "No survivors reported") {
		t.Errorf("Final() should report no survivors when all ships dead: %v", msgs)
	}
}

func TestFinal_UnknownCode(t *testing.T) {
	st := newEndgameState()
	msgs := Final(st, 999, false)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "How did we get here") {
		t.Errorf("Final() with unknown code should return fallback message: %v", msgs)
	}
}

func TestWarn(t *testing.T) {
	tests := []struct {
		name string
		mesg int
		want string
	}{
		{"FedLose", game.FinFedLose, "Message to the Federation"},
		{"EnemyLose", game.FinEnemyLose, "destroyed or crippled"},
		{"Tactical", game.FinTactical, "falling behind"},
		{"FedSurrender", game.FinFedSurrender, "informing Starfleet"},
		{"EnemySurrender", game.FinEnemySurrender, "surrendered"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st := newEndgameState()
			msgs := Warn(st, tc.mesg)
			joined := strings.Join(msgs, "\n")
			if !strings.Contains(joined, tc.want) {
				t.Errorf("Warn(%d) = %v, want to contain %q", tc.mesg, msgs, tc.want)
			}
		})
	}
}

func TestWarn_TacticalSuppressedAfterReengage(t *testing.T) {
	st := newEndgameState()
	st.Reengaged = true
	msgs := Warn(st, game.FinTactical)
	if msgs != nil {
		t.Errorf("Warn(FinTactical) after reengage should be nil, got %v", msgs)
	}
}

func TestWarn_UnknownCode(t *testing.T) {
	st := newEndgameState()
	msgs := Warn(st, 999)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "How did we get here") {
		t.Errorf("Warn() with unknown code should return fallback message: %v", msgs)
	}
}

func TestMission(t *testing.T) {
	st := newEndgameState()
	msgs := Mission(st, 0, "7301.1")
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "final frontier") {
		t.Errorf("Mission() missing intro line: %v", msgs)
	}
	if !strings.Contains(joined, "Enterprise") {
		t.Errorf("Mission() missing ship name: %v", msgs)
	}
	if !strings.Contains(joined, "7301.1") {
		t.Errorf("Mission() missing stardate: %v", msgs)
	}
	if !strings.Contains(joined, "Klingon") {
		t.Errorf("Mission() missing enemy race name: %v", msgs)
	}
	if !strings.Contains(joined, "Kor") {
		t.Errorf("Mission() missing enemy commander name: %v", msgs)
	}
	// Single enemy ship: should identify "it" as "a Klingon Battlecruiser,".
	if !strings.Contains(joined, "Sensors identify it as a Klingon Battlecruiser,") {
		t.Errorf("Mission() single-ship identify line wrong: %v", msgs)
	}
}

func TestMission_MultipleEnemies(t *testing.T) {
	fed := newTestShip(1, "Enterprise", 0, 0)
	e1 := newTestShip(2, "Klathis", 1, 1)
	e2 := newTestShip(3, "Kalyx", 2, 2)
	st := newTestState(fed, e1, e2)
	st.EnemyRaceName = "Klingon"
	st.EnemyShipTypeName = "Battlecruiser"
	st.EnemyCommander = "Kor"

	msgs := Mission(st, 0, "7301.1")
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "Sensors identify them as Klingon Battlecruisers,") {
		t.Errorf("Mission() multi-ship identify line wrong: %v", msgs)
	}
	if !strings.Contains(joined, "picking up 2 vessels") {
		t.Errorf("Mission() should report vessel count for multiple enemies: %v", msgs)
	}
}

func TestMission_OutOfRangeIndex(t *testing.T) {
	st := newEndgameState()
	// Should not panic and should simply skip the briefing lines.
	msgs := Mission(st, -1, "7301.1")
	if len(msgs) == 0 {
		t.Error("Mission() with out-of-range index should still return base messages")
	}
	msgs = Mission(st, 99999, "7301.1")
	if len(msgs) == 0 {
		t.Error("Mission() with out-of-range index should still return base messages")
	}
}

func TestAlert_SingleEnemy(t *testing.T) {
	st := newEndgameState()
	msgs := Alert(st)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "Klingon") {
		t.Errorf("Alert() missing race name: %v", msgs)
	}
	if !strings.Contains(joined, "Klathis") {
		t.Errorf("Alert() missing enemy ship name: %v", msgs)
	}
	if !strings.Contains(joined, "Enterprise") {
		t.Errorf("Alert() missing player ship name: %v", msgs)
	}
}

func TestAlert_MultipleEnemies(t *testing.T) {
	fed := newTestShip(1, "Enterprise", 0, 0)
	e1 := newTestShip(2, "Klathis", 1, 1)
	e2 := newTestShip(3, "Kalyx", 2, 2)
	st := newTestState(fed, e1, e2)
	st.EnemyRaceName = "Klingon"

	msgs := Alert(st)
	joined := strings.Join(msgs, "\n")
	if !strings.Contains(joined, "Klathis") || !strings.Contains(joined, "Kalyx") {
		t.Errorf("Alert() should mention every enemy ship name: %v", msgs)
	}
}
