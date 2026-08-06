package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

func TestSessionBlankCommandShowsHelp(t *testing.T) {
	r := game.NewRand(1)
	res, err := game.NewGame(game.InitOptions{EnemyCount: 1}, r)
	if err != nil {
		t.Fatalf("NewGame() error = %v", err)
	}
	res.State.Crew = game.CrewNames{
		Captain:  "Kirk",
		Title:    "Sir",
		Science:  "Spock",
		Engineer: "Scott",
		Com:      "Uhura",
		Nav:      "Chekov",
		Helmsman: "Sulu",
	}

	in := bufio.NewReader(strings.NewReader("\n"))
	var out bytes.Buffer
	session := newCLISession(in, &out, res.State, r)
	if err := session.run(); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Trek73 Commands:") {
		t.Fatalf("output missing help list: %q", output)
	}
	if strings.Contains(output, "Illegal input.") {
		t.Fatalf("blank command should not be treated as illegal: %q", output)
	}
}
