package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestPromptEnemyVesselCountRepromptsUntilValid(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("0\n10\n5\n"))
	var out bytes.Buffer

	got, err := promptEnemyVesselCount(reader, &out)
	if err != nil {
		t.Fatalf("promptEnemyVesselCount() error = %v", err)
	}
	if got != 5 {
		t.Fatalf("promptEnemyVesselCount() = %d, want 5", got)
	}

	output := out.String()
	if strings.Count(output, "I'm expecting [1-9] enemy vessels: ") != 3 {
		t.Fatalf("prompt count = %d, want 3; output=%q", strings.Count(output, "I'm expecting [1-9] enemy vessels: "), output)
	}
}
