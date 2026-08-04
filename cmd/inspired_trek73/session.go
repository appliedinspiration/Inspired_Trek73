package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/appliedinspiration/inspired_trek73/internal/commands"
	"github.com/appliedinspiration/inspired_trek73/internal/game"
)

var errInvalidArguments = errors.New("invalid arguments")

type cliSession struct {
	reader *bufio.Reader
	out    io.Writer
	state  *game.State
	rand   *game.Rand
	hooks  *game.CombatHooks
}

func newCLISession(reader *bufio.Reader, out io.Writer, st *game.State, r *game.Rand) *cliSession {
	return &cliSession{
		reader: reader,
		out:    out,
		state:  st,
		rand:   r,
		hooks:  game.NewCombatHooks(st, r),
	}
}

func (c *cliSession) run() error {
	for {
		c.state.Shutup.Reset()

		line, err := promptLine(c.reader, c.out, "\nCommand: ")
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(c.out)
				fmt.Fprintln(c.out, "End of input. Exiting.")
				return nil
			}
			return err
		}

		code, args, ok := commands.ParseCommand(line)
		if !ok {
			fmt.Fprintln(c.out, "Illegal input.")
			continue
		}

		spec, ok := commands.LookupCommand(code)
		if !ok {
			fmt.Fprintln(c.out, "Illegal input.")
			continue
		}

		executed, err := c.execute(code, args)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(c.out)
				fmt.Fprintln(c.out, "End of input. Exiting.")
				return nil
			}
			if errors.Is(err, errInvalidArguments) {
				fmt.Fprintln(c.out, "Invalid arguments for that command.")
				continue
			}
			return err
		}
		if !executed || !spec.Turn {
			continue
		}

		printMessages(c.out, game.RunTurn(c.state, c.hooks))

		ended, err := c.handleEndgame()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(c.out)
				fmt.Fprintln(c.out, "End of input. Exiting.")
				return nil
			}
			return err
		}
		if ended {
			fmt.Fprintln(c.out, "Game over.")
			return nil
		}
	}
}

func (c *cliSession) execute(code int, args []string) (bool, error) {
	sp := c.state.Player()

	switch code {
	case 1:
		banks, spread, err := c.resolveFirePhasersArgs(sp, args)
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.FirePhasers(sp, banks, spread))
		return true, nil
	case 2:
		banks, err := c.resolveBanksArg(args, 0, "   Fire which tube(s) [all or digits]? ")
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.FireTubes(sp, banks))
		return true, nil
	case 3:
		banks, target, err := c.resolveLockArgs(args, "   Lock which phaser bank(s) [all or digits]? ", "   Lock onto which ship? ")
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.LockPhasers(c.state, c.rand, sp, banks, target))
		return true, nil
	case 4:
		banks, target, err := c.resolveLockArgs(args, "   Lock which tube(s) [all or digits]? ", "   Lock onto which ship? ")
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.LockTubes(c.state, c.rand, sp, banks, target))
		return true, nil
	case 5:
		banks, bearing, err := c.resolveTurnArgs(args, "   Rotate which phaser bank(s) [all or digits]? ", "   Bearing [0-360]: ")
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.TurnPhasers(sp, banks, bearing))
		return true, nil
	case 6:
		banks, bearing, err := c.resolveTurnArgs(args, "   Rotate which tube(s) [all or digits]? ", "   Bearing [0-360]: ")
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.TurnTubes(sp, banks, bearing))
		return true, nil
	case 7:
		printMessages(c.out, commands.PhaserStatus(sp))
		return true, nil
	case 8:
		printMessages(c.out, commands.TubeStatus(sp))
		return true, nil
	case 9:
		load, banks, err := c.resolveTubeLoadArgs(args)
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.LoadTubes(c.state, sp, load, banks))
		return true, nil
	case 10:
		obj, messages, err := c.launchProbe(sp, args)
		if err != nil {
			return false, err
		}
		if obj != nil {
			c.state.Objects = append(c.state.Objects, obj)
		}
		printMessages(c.out, messages)
		return true, nil
	case 11:
		return c.controlProbe(sp, args)
	case 12:
		printMessages(c.out, commands.PosReport(c.state, sp))
		return true, nil
	case 13:
		rng, err := c.resolveIntArg(args, 0, "   Display range [100-50000]: ")
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.PosDisplay(sp, c.state, c.rand, rng))
		return true, nil
	case 14:
		target, warp, err := c.resolvePursuitArgs(args, "   Mr. "+c.state.Crew.Nav+", pursue [who] ", "   Mr. "+c.state.Crew.Helmsman+", warp factor: ")
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.Pursue(c.state, c.rand, sp, target, warp))
		return true, nil
	case 15:
		target, warp, err := c.resolvePursuitArgs(args, "   Mr. "+c.state.Crew.Nav+", elude [who] ", "   Mr. "+c.state.Crew.Helmsman+", warp factor: ")
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.Elude(c.state, c.rand, sp, target, warp))
		return true, nil
	case 16:
		course, warp, err := c.resolveHelmArgs(args)
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.Helm(c.state, sp, course, warp))
		return true, nil
	case 17:
		printMessages(c.out, commands.SelfScan(sp))
		return true, nil
	case 18:
		who, err := c.resolveStringArg(args, 0, "   "+c.state.Crew.Science+", scan [who] ")
		if err != nil {
			return false, err
		}
		_, messages := commands.Scan(c.state, c.rand, sp, who)
		printMessages(c.out, messages)
		return true, nil
	case 19:
		shields, phasers, err := c.promptAlterPower(sp)
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.AlterPower(c.state, sp, shields, phasers))
		return true, nil
	case 20:
		obj, messages := commands.JettisonEngineering(c.state, sp)
		if obj != nil {
			c.state.Objects = append(c.state.Objects, obj)
		}
		printMessages(c.out, messages)
		return true, nil
	case 21:
		confirm := false
		if !sp.IsDead(game.SysEngineering) {
			var err error
			confirm, err = c.promptYesNo("Engineering is still attached. Jettison and detonate anyway? [Y/n] ", true)
			if err != nil {
				return false, err
			}
		}
		obj, messages := commands.DetonateEngineering(c.state, sp, confirm)
		if obj != nil {
			c.state.Objects = append(c.state.Objects, obj)
		}
		printMessages(c.out, messages)
		return true, nil
	case 22:
		resetTubes, launchSpeed, timeDelay, proxDelay, resetPhasers, firePercent, err := c.promptAlterFiringParams(sp)
		if err != nil {
			return false, err
		}
		commands.AlterFiringParams(sp, resetTubes, launchSpeed, timeDelay, proxDelay, resetPhasers, firePercent)
		return true, nil
	case 23:
		target, err := c.resolvePowerTransferTarget(args)
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.PlayDead(c.state, sp, target))
		return true, nil
	case 24:
		printMessages(c.out, commands.CorbomiteBluff(c.state, sp, c.rand.Randm(2) == 1))
		return true, nil
	case 25:
		printMessages(c.out, commands.SurrenderShip(c.state, sp))
		return true, nil
	case 26:
		printMessages(c.out, commands.RequestSurrender(c.state, sp))
		return true, nil
	case 27:
		printMessages(c.out, commands.SelfDestruct(c.state, c.rand, sp))
		return true, nil
	case 28:
		printMessages(c.out, commands.AbortSelfDestruct(c.state, c.rand, sp))
		return true, nil
	case 29:
		printMessages(c.out, commands.Survivors(c.state))
		return true, nil
	case 30:
		printMessages(c.out, commands.Version())
		return true, nil
	case 31:
		fmt.Fprintln(c.out, "Save game is not yet implemented in this version.")
		return true, nil
	case 32:
		printMessages(c.out, commands.Help())
		return true, nil
	default:
		return false, errInvalidArguments
	}
}

func (c *cliSession) handleEndgame() (bool, error) {
	if c.state.PendingWarn != -1 {
		code := c.state.PendingWarn
		if code >= 0 && code < len(c.state.WarnShown) && !c.state.WarnShown[code] {
			printMessages(c.out, commands.Warn(c.state, code))
			c.state.WarnShown[code] = true
		}
		c.state.PendingWarn = -1
	}

	if c.state.PendingFinal == -1 {
		return false, nil
	}

	code := c.state.PendingFinal
	c.state.PendingFinal = -1

	if code == game.FinTactical && !c.state.Reengaged {
		answer, err := c.promptYesNo("Do you wish to re-engage? [Y/n] ", true)
		if err != nil {
			return false, err
		}
		printMessages(c.out, commands.Final(c.state, code, answer))
		return !answer, nil
	}

	printMessages(c.out, commands.Final(c.state, code, false))
	return true, nil
}

func (c *cliSession) resolveFirePhasersArgs(sp *game.Ship, args []string) (string, int, error) {
	banks, err := c.resolveBanksArg(args, 0, "   Fire which phaser bank(s) [all or digits]? ")
	if err != nil {
		return "", 0, err
	}
	if len(args) > 1 {
		spread, err := strconv.Atoi(args[1])
		if err != nil {
			return "", 0, errInvalidArguments
		}
		return banks, spread, nil
	}
	spread, err := c.promptIntWithDefault("   Phaser spread [10-45] (blank keeps current): ", sp.PhaserSpread)
	if err != nil {
		return "", 0, err
	}
	return banks, spread, nil
}

func (c *cliSession) resolveLockArgs(args []string, banksPrompt, targetPrompt string) (string, *game.Ship, error) {
	banks, err := c.resolveBanksArg(args, 0, banksPrompt)
	if err != nil {
		return "", nil, err
	}
	targetName := ""
	if len(args) > 1 {
		targetName = args[1]
	} else {
		targetName, err = c.promptNonEmpty(targetPrompt)
		if err != nil {
			return "", nil, err
		}
	}
	target := c.resolveEnemy(targetName)
	if target == nil {
		return "", nil, errInvalidArguments
	}
	return banks, target, nil
}

func (c *cliSession) resolveTurnArgs(args []string, banksPrompt, bearingPrompt string) (string, float64, error) {
	banks, err := c.resolveBanksArg(args, 0, banksPrompt)
	if err != nil {
		return "", 0, err
	}
	bearingText, err := c.resolveStringArg(args, 1, bearingPrompt)
	if err != nil {
		return "", 0, err
	}
	bearing, err := strconv.ParseFloat(bearingText, 64)
	if err != nil {
		return "", 0, errInvalidArguments
	}
	return banks, bearing, nil
}

func (c *cliSession) resolveTubeLoadArgs(args []string) (bool, string, error) {
	action := ""
	if len(args) > 0 {
		action = args[0]
	} else {
		text, err := c.promptNonEmpty("   Load or unload tubes? [l/u] ")
		if err != nil {
			return false, "", err
		}
		action = text
	}

	load := false
	switch strings.ToLower(action[:1]) {
	case "l":
		load = true
	case "u":
		load = false
	default:
		return false, "", errInvalidArguments
	}

	banks, err := c.resolveBanksArg(args, 1, "   Which tube(s) [all or digits]? ")
	if err != nil {
		return false, "", err
	}
	return load, banks, nil
}

func (c *cliSession) launchProbe(sp *game.Ship, args []string) (*game.SpaceObject, []string, error) {
	pods, err := c.resolveIntArg(args, 0, fmt.Sprintf("%s: Number to launch [%d+]: ", c.state.Crew.Captain, game.MinProbeCharge))
	if err != nil {
		return nil, nil, err
	}
	delay, err := c.resolveIntArg(args, 1, fmt.Sprintf("   Set time delay [0-%d]: ", game.MaxProbeDelay))
	if err != nil {
		return nil, nil, err
	}
	prox, err := c.resolveIntArg(args, 2, fmt.Sprintf("   Set proximity delay [%d+]: ", game.MinProbeProx))
	if err != nil {
		return nil, nil, err
	}

	var target *game.Ship
	course := 0.0
	if len(args) > 3 {
		target = c.resolveEnemy(args[3])
		if target == nil {
			course, err = strconv.ParseFloat(args[3], 64)
			if err != nil {
				return nil, nil, errInvalidArguments
			}
		}
	} else {
		targetName, err := promptLine(c.reader, c.out, "   Launch toward [whom, if anyone]? ")
		if err != nil {
			return nil, nil, err
		}
		targetName = strings.TrimSpace(targetName)
		if targetName != "" {
			target = c.resolveEnemy(targetName)
			if target == nil {
				return nil, nil, errInvalidArguments
			}
		} else {
			course, err = c.resolveFloatArg(nil, 0, "   Course [0-360]: ")
			if err != nil {
				return nil, nil, err
			}
		}
	}

	obj, messages := commands.LaunchProbe(c.state, c.rand, sp, pods, delay, prox, target, course)
	return obj, messages, nil
}

func (c *cliSession) controlProbe(sp *game.Ship, args []string) (bool, error) {
	probes := commands.ShipProbes(c.state, sp)
	if len(probes) == 0 {
		fmt.Fprintf(c.out, "%s: What probes?\n", c.state.Crew.Nav)
		return true, nil
	}

	c.printProbeTable(sp, probes)

	var probe *game.SpaceObject
	if len(args) == 0 {
		detonateAll, err := c.promptYesNo(c.state.Crew.Nav+": Detonate all probes? [y/N] ", false)
		if err != nil {
			return false, err
		}
		if detonateAll {
			printMessages(c.out, commands.DetonateAllProbes(sp, probes))
			return true, nil
		}

		id, err := c.resolveIntArg(nil, 0, "   Control probe [#]: ")
		if err != nil {
			return false, err
		}
		probe = probeByID(probes, id)
	} else {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return false, errInvalidArguments
		}
		probe = probeByID(probes, id)
	}
	if probe == nil {
		return false, errInvalidArguments
	}

	detonateIt, err := c.promptYesNo(c.state.Crew.Nav+": Detonate it? [y/N] ", false)
	if err != nil {
		return false, err
	}
	if detonateIt {
		printMessages(c.out, commands.DetonateProbe(sp, probe))
		return true, nil
	}

	targetName, err := promptLine(c.reader, c.out, "   Lock it onto [whom, if anyone]? ")
	if err != nil {
		return false, err
	}
	targetName = strings.TrimSpace(targetName)
	if targetName != "" {
		target := c.resolveEnemy(targetName)
		if target == nil {
			return false, errInvalidArguments
		}
		printMessages(c.out, commands.LockProbe(c.state, c.rand, sp, probe, target))
		return true, nil
	}

	course, err := c.resolveFloatArg(nil, 0, "   Set it to course [0-360]: ")
	if err != nil {
		return false, err
	}
	printMessages(c.out, commands.CourseProbe(c.state, sp, probe, course))
	return true, nil
}

func (c *cliSession) resolvePursuitArgs(args []string, targetPrompt, warpPrompt string) (*game.Ship, float64, error) {
	targetName := ""
	var err error
	if len(args) > 0 {
		targetName = args[0]
	} else {
		targetName, err = c.promptNonEmpty(targetPrompt)
		if err != nil {
			return nil, 0, err
		}
	}
	target := c.resolveEnemy(targetName)
	if target == nil {
		return nil, 0, errInvalidArguments
	}

	warpText, err := c.resolveStringArg(args, 1, warpPrompt)
	if err != nil {
		return nil, 0, err
	}
	warp, err := strconv.ParseFloat(warpText, 64)
	if err != nil {
		return nil, 0, errInvalidArguments
	}
	return target, warp, nil
}

func (c *cliSession) resolveHelmArgs(args []string) (float64, float64, error) {
	course, err := c.resolveFloatArg(args, 0, "   Mr. "+c.state.Crew.Nav+", come to course [0-359]: ")
	if err != nil {
		return 0, 0, err
	}
	warp, err := c.resolveFloatArg(args, 1, "   Mr. "+c.state.Crew.Helmsman+", warp factor: ")
	if err != nil {
		return 0, 0, err
	}
	return course, warp, nil
}

func (c *cliSession) promptAlterPower(sp *game.Ship) ([]float64, []float64, error) {
	fmt.Fprintf(c.out, "\n%s:  Regeneration rate is %5.2f.\n", c.state.Crew.Engineer, sp.Regen)

	shields, err := c.promptFloatSeries(
		len(sp.Shields),
		func(i int) string {
			return fmt.Sprintf("%s:  Shield %d drain is [0.0 to 1.0] ", c.state.Crew.Captain, i+1)
		},
		0.0,
		1.0,
	)
	if err != nil {
		return nil, nil, err
	}

	fmt.Fprintln(c.out)
	phasers, err := c.promptFloatSeries(
		len(sp.Phasers),
		func(i int) string {
			return fmt.Sprintf("%s:  Phaser %d drain is [%.0f to %.0f] ", c.state.Crew.Captain, i+1, game.MinPhaserDrain, game.MaxPhaserDrain)
		},
		game.MinPhaserDrain,
		game.MaxPhaserDrain,
	)
	if err != nil {
		return nil, nil, err
	}

	return shields, phasers, nil
}

func (c *cliSession) promptAlterFiringParams(sp *game.Ship) (bool, int, int, int, bool, int, error) {
	resetTubes, err := c.promptYesNo(fmt.Sprintf("\n%s: Reset tubes, %s? [Y/n] ", c.state.Crew.Nav, c.state.Crew.Title), true)
	if err != nil {
		return false, 0, 0, 0, false, 0, err
	}

	launchSpeed := sp.TubeLaunchSpd
	timeDelay := sp.TubeDelay
	proxDelay := sp.TubeProximity
	if resetTubes {
		launchSpeed, err = c.promptIntWithDefault(fmt.Sprintf("   Set launch speed [0-%d] (blank keeps current): ", game.MaxTubeSpeed), sp.TubeLaunchSpd)
		if err != nil {
			return false, 0, 0, 0, false, 0, err
		}
		timeDelay, err = c.promptIntWithDefault(fmt.Sprintf("   Set time delay [0-%d] (blank keeps current): ", int(game.MaxTubeTime)), sp.TubeDelay)
		if err != nil {
			return false, 0, 0, 0, false, 0, err
		}
		proxDelay, err = c.promptIntWithDefault(fmt.Sprintf("   Set proximity delay [0-%d] (blank keeps current): ", game.MaxTubeProx), sp.TubeProximity)
		if err != nil {
			return false, 0, 0, 0, false, 0, err
		}
	}

	resetPhasers, err := c.promptYesNo(fmt.Sprintf("%s: Reset phasers, %s? [Y/n] ", c.state.Crew.Nav, c.state.Crew.Title), true)
	if err != nil {
		return false, 0, 0, 0, false, 0, err
	}

	firePercent := sp.PhaserFirePct
	if resetPhasers {
		firePercent, err = c.promptIntWithDefault("   Reset firing percentage to [0-100] (blank keeps current): ", sp.PhaserFirePct)
		if err != nil {
			return false, 0, 0, 0, false, 0, err
		}
	}

	return resetTubes, launchSpeed, timeDelay, proxDelay, resetPhasers, firePercent, nil
}

func (c *cliSession) resolvePowerTransferTarget(args []string) (commands.PowerTransferTarget, error) {
	if len(args) > 0 {
		if target := parseTransferTarget(args[0]); target != commands.TransferNone {
			return target, nil
		}
	}

	answer, err := c.promptNonEmpty("   Transfer power to [engines or phasers]: ")
	if err != nil {
		return commands.TransferNone, err
	}
	target := parseTransferTarget(answer)
	if target == commands.TransferNone {
		return commands.TransferNone, errInvalidArguments
	}
	return target, nil
}

func parseTransferTarget(text string) commands.PowerTransferTarget {
	switch strings.ToLower(strings.TrimSpace(text)) {
	case "engines", "engine", "e":
		return commands.TransferEngines
	case "phasers", "phaser", "p":
		return commands.TransferPhasers
	default:
		return commands.TransferNone
	}
}

func (c *cliSession) resolveEnemy(name string) *game.Ship {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}

	if n, err := strconv.Atoi(name); err == nil {
		enemies := c.state.Enemies()
		if n >= 1 && n <= len(enemies) {
			return enemies[n-1]
		}
	}

	return c.state.ShipByName(name)
}

func probeByID(probes []*game.SpaceObject, id int) *game.SpaceObject {
	for _, probe := range probes {
		if probe.ID == id {
			return probe
		}
	}
	return nil
}

func (c *cliSession) printProbeTable(sp *game.Ship, probes []*game.SpaceObject) {
	fmt.Fprintln(c.out)
	fmt.Fprintln(c.out, "probe bearng range course time  prox units target")
	for _, probe := range probes {
		target := "NONE"
		if probe.Target != nil {
			target = probe.Target.Name
		}
		rangeToProbe := game.RangeFind(sp.X, probe.X, sp.Y, probe.Y)
		bearing := game.BearingTo(sp.X, probe.X, sp.Y, probe.Y)
		fmt.Fprintf(c.out, " %2d    %4.1f %5d  %4.0f  %4.1f %5d  %3d  %s\n",
			probe.ID, bearing, rangeToProbe, probe.Course, probe.TimeDelay, probe.Proximity, int(probe.Fuel), target)
	}
}

func (c *cliSession) promptFloatSeries(count int, prompt func(i int) string, min, max float64) ([]float64, error) {
	values := make([]float64, 0, count)
	for i := 0; i < count; i++ {
		text, err := c.promptNonEmpty(prompt(i))
		if err != nil {
			return nil, err
		}

		fillRest := strings.HasSuffix(text, "*")
		text = strings.TrimSpace(strings.TrimSuffix(text, "*"))

		value, err := strconv.ParseFloat(text, 64)
		if err != nil || value < min || value > max {
			return nil, errInvalidArguments
		}

		values = append(values, value)
		if fillRest {
			for len(values) < count {
				values = append(values, value)
			}
			break
		}
	}
	if len(values) != count {
		return nil, errInvalidArguments
	}
	return values, nil
}

func (c *cliSession) resolveBanksArg(args []string, index int, prompt string) (string, error) {
	value, err := c.resolveStringArg(args, index, prompt)
	if err != nil {
		return "", err
	}
	return strings.ToLower(strings.TrimSpace(value)), nil
}

func (c *cliSession) resolveStringArg(args []string, index int, prompt string) (string, error) {
	if index < len(args) {
		value := strings.TrimSpace(args[index])
		if value != "" {
			return value, nil
		}
	}
	return c.promptNonEmpty(prompt)
}

func (c *cliSession) resolveIntArg(args []string, index int, prompt string) (int, error) {
	text, err := c.resolveStringArg(args, index, prompt)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, errInvalidArguments
	}
	return value, nil
}

func (c *cliSession) resolveFloatArg(args []string, index int, prompt string) (float64, error) {
	text, err := c.resolveStringArg(args, index, prompt)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, errInvalidArguments
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, errInvalidArguments
	}
	return value, nil
}

func (c *cliSession) promptIntWithDefault(prompt string, current int) (int, error) {
	text, err := promptLine(c.reader, c.out, prompt)
	if err != nil {
		return 0, err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return current, nil
	}
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, errInvalidArguments
	}
	return value, nil
}

func (c *cliSession) promptNonEmpty(prompt string) (string, error) {
	text, err := promptLine(c.reader, c.out, prompt)
	if err != nil {
		return "", err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return "", errInvalidArguments
	}
	return text, nil
}

func (c *cliSession) promptYesNo(prompt string, defaultYes bool) (bool, error) {
	text, err := promptLine(c.reader, c.out, prompt)
	if err != nil {
		return false, err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultYes, nil
	}
	switch strings.ToLower(text[:1]) {
	case "y":
		return true, nil
	case "n":
		return false, nil
	default:
		return false, errInvalidArguments
	}
}

func promptLine(reader *bufio.Reader, w io.Writer, prompt string) (string, error) {
	if _, err := fmt.Fprint(w, prompt); err != nil {
		return "", err
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && len(line) > 0 {
			return strings.TrimRight(line, "\r\n"), nil
		}
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func printMessages(w io.Writer, messages []string) {
	for _, message := range messages {
		fmt.Fprintln(w, message)
	}
}
