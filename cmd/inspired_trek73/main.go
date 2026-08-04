// Package main is the CLI entry point for Inspired Trek73.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/appliedinspiration/inspired_trek73/internal/banner"
	"github.com/appliedinspiration/inspired_trek73/internal/commands"
	"github.com/appliedinspiration/inspired_trek73/internal/data"
	"github.com/appliedinspiration/inspired_trek73/internal/game"
	"github.com/appliedinspiration/inspired_trek73/internal/version"
)

func main() {
	// Go's flag package treats "-version" and "--version" identically,
	// satisfying both spellings requested for the CLI.
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "print version information and exit")
	flag.Parse()

	banner.Print(os.Stdout)

	if showVersion {
		fmt.Println()
		fmt.Println(version.String())
		fmt.Println()
		banner.Print(os.Stdout)
		return
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func run() error {
	reader := bufio.NewReader(os.Stdin)
	r := game.NewRandFromTime()

	crew, err := promptCrewNames(reader, os.Stdout, r)
	if err != nil {
		return err
	}

	res, err := game.NewGame(game.InitOptions{}, r)
	if err != nil {
		return err
	}

	st := res.State
	st.Crew = crew

	printLines(os.Stdout, res.Warnings)

	stardate := formatStardate(time.Now())
	missionIndex := r.Randm(len(data.MissionBriefings)) - 1
	printLines(os.Stdout, commands.Mission(st, missionIndex, stardate))
	printLines(os.Stdout, commands.Alert(st))

	session := newCLISession(reader, os.Stdout, st, r)
	return session.run()
}

func formatStardate(now time.Time) string {
	return fmt.Sprintf("%02d%02d.%02d", now.Year()%100, int(now.Month()), now.Day())
}

func printLines(w io.Writer, lines []string) {
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
}
