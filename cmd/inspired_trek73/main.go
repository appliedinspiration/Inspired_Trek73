// Package main is the CLI entry point for Inspired Trek73.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/appliedinspiration/inspired_trek73/internal/banner"
	"github.com/appliedinspiration/inspired_trek73/internal/commands"
	"github.com/appliedinspiration/inspired_trek73/internal/data"
	"github.com/appliedinspiration/inspired_trek73/internal/game"
	"github.com/appliedinspiration/inspired_trek73/internal/version"
)

type cliOptions struct {
	ShowVersion        bool
	ShowHelp           bool
	EnemyCountProvided bool
	Init               game.InitOptions
}

func main() {
	opts := parseFlags()

	if opts.ShowHelp {
		flag.Usage()
		return
	}

	// Go's flag package treats "-version" and "--version" identically,
	// satisfying both spellings requested for the CLI.
	banner.Print(os.Stdout)

	if opts.ShowVersion {
		fmt.Println()
		fmt.Println(version.String())
		fmt.Println()
		banner.Print(os.Stdout)
		return
	}

	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func parseFlags() cliOptions {
	opts := cliOptions{}

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
	}

	flag.BoolVar(&opts.ShowVersion, "version", false, "print version information and exit")
	flag.BoolVar(&opts.ShowHelp, "help", false, "show available options and exit")
	flag.BoolVar(&opts.ShowHelp, "h", false, "show available options and exit")
	flag.IntVar(&opts.Init.EnemyCount, "enemies", 0, "number of enemy ships (1-9, 0=random)")
	flag.StringVar(&opts.Init.PlayerClassAbbr, "player-class", "", "player ship class abbreviation (default: CA)")
	flag.StringVar(&opts.Init.EnemyClassAbbr, "enemy-class", "", "enemy ship class abbreviation (default: CA)")
	flag.StringVar(&opts.Init.RaceName, "race", "", "enemy race name prefix (default: random)")
	flag.BoolVar(&opts.Init.AllowSillyRace, "allow-silly-race", false, "allow Monty Python race in random selection")
	flag.StringVar(&opts.Init.PlayerShipName, "ship-name", "", "player ship name (default: random Federation ship)")
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "enemies" {
			opts.EnemyCountProvided = true
		}
	})

	return opts
}

func run(opts cliOptions) error {
	reader := bufio.NewReader(os.Stdin)
	r := game.NewRandFromTime()

	crew, err := promptCrewNames(reader, os.Stdout, r)
	if err != nil {
		return err
	}

	initOpts := opts.Init
	if !opts.EnemyCountProvided {
		enemyCount, err := promptEnemyVesselCount(reader, os.Stdout)
		if err != nil {
			return err
		}
		initOpts.EnemyCount = enemyCount
	}

	res, err := game.NewGame(initOpts, r)
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

func promptEnemyVesselCount(reader *bufio.Reader, w io.Writer) (int, error) {
	for {
		line, err := promptLine(reader, w, "I'm expecting [1-9] enemy vessels: ")
		if err != nil {
			return 0, err
		}
		n, err := strconv.Atoi(strings.TrimSpace(line))
		if err == nil && n >= 1 && n <= game.MaxEnemyShips {
			return n, nil
		}
		fmt.Fprintln(w, "Please enter a number from 1 to 9.")
	}
}
