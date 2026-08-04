package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

var fallbackTitles = []string{"Fag", "Fairy", "Fruit", "Weirdo", "Gumby", "Freak"}

// promptCrewNames ports the interactive portion of name_crew() in
// init.c. The original's "*captain" terse-mode shortcut is intentionally
// omitted here because this CLI has no separate terse presentation mode.
func promptCrewNames(reader *bufio.Reader, w io.Writer, r *game.Rand) (game.CrewNames, error) {
	crew := game.CrewNames{
		Science:  "Spock",
		Engineer: "Scott",
		Com:      "Uhura",
		Nav:      "Chekov",
		Helmsman: "Sulu",
	}

	captain, err := promptRequired(reader, w, "\n\nCaptain: my last name is ", "Captain name is required. Exiting.")
	if err != nil {
		return game.CrewNames{}, err
	}
	crew.Captain = captain

	sex, err := promptRequired(reader, w, fmt.Sprintf("%s: My sex is: ", captain), "Sex is required. Exiting.")
	if err != nil {
		return game.CrewNames{}, err
	}

	switch sex[0] {
	case 'M', 'm':
		crew.Title = "Sir"
	case 'F', 'f':
		crew.Title = "Ma'am"
	default:
		crew.Title = fallbackTitles[r.Randm(len(fallbackTitles))-1]
	}

	return crew, nil
}

func promptRequired(reader *bufio.Reader, w io.Writer, prompt, emptyMessage string) (string, error) {
	value, err := promptLine(reader, w, prompt)
	if err != nil {
		return "", err
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%s", emptyMessage)
	}
	return value, nil
}
