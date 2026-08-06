// Package commands implements Trek73's ~32 player command handlers as
// pure functions, plus the numeric command table and text-command
// parser used to route raw player input to them.
package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/appliedinspiration/inspired_trek73/internal/version"
)

// CommandSpec describes one numeric command slot, ported from struct
// cmd (cmds[] in globals.c): a code number, a short description (used
// by Help), and whether issuing it consumes the player's turn (TURN)
// or not (FREE).
type CommandSpec struct {
	Code        int
	Description string
	Turn        bool // true == TURN (ends the turn); false == FREE (loops back for another command).
}

// CommandTable mirrors globals.c's cmds[] array (codes 1-32). It does
// not hold function pointers, matching this port's pure-function
// design: the CLI-wiring layer is responsible for mapping a resolved
// Code to the correct commands.* function call (each of which has its
// own, command-specific parameter list that can't be captured
// generically here).
var CommandTable = []CommandSpec{
	{1, "Fire phasers", true},
	{2, "Fire photon torpedoes", true},
	{3, "Lock phasers onto target", true},
	{4, "Lock tubes onto target", true},
	{5, "Manually rotate phasers", true},
	{6, "Manually rotate tubes", true},
	{7, "Phaser status", false},
	{8, "Tube status", false},
	{9, "Load/unload torpedo tubes", true},
	{10, "Launch antimatter probe", true},
	{11, "Probe control", true},
	{12, "Position report", false},
	{13, "Position display", false},
	{14, "Pursue an enemy vessel", true},
	{15, "Elude an enemy vessel", true},
	{16, "Change course and speed", true},
	{17, "Damage report", false},
	{18, "Scan enemy", true},
	{19, "Alter power distribution", true},
	{20, "Jettison engineering", true},
	{21, "Detonate engineering", true},
	{22, "Alter firing parameters", true},
	{23, "Attempt defenseless ruse", true},
	{24, "Attempt corbomite bluff(s)", true},
	{25, "Surrender", true},
	{26, "Ask enemy to surrender", true},
	{27, "Initiate self-destruct", true},
	{28, "Abort self-destruct", true},
	{29, "Survivors report", false},
	{30, "Print version number", false},
	{31, "Saves game", false},
	{32, "Reprints above list", false},
}

// LookupCommand returns the CommandSpec for the given numeric code, and
// whether it exists (code must be in [1, len(CommandTable)]), ported
// from the linear cmds[] scan in scancmd() (misc.c), simplified since
// exact-integer lookup is unambiguous (see ParseNumeric doc comment).
func LookupCommand(code int) (CommandSpec, bool) {
	if code < 1 || code > len(CommandTable) {
		return CommandSpec{}, false
	}
	return CommandTable[code-1], true
}

// Version returns the version report shown by the "vers"/code-30
// command, ported from vers() in misc.c (adapted to use this port's
// own internal/version metadata rather than the original's compiled-in
// RCS identifiers).
func Version() []string {
	return strings.Split(version.String(), "\n")
}

// Help lists every command's code, description, and TURN/FREE marker,
// ported from help() in misc.c.
func Help() []string {
	msgs := []string{"", "Trek73 Commands:", ""}
	mid := (len(CommandTable) + 1) / 2
	for i := 0; i < mid; i++ {
		left := formatHelpEntry(CommandTable[i])
		rightIndex := i + mid
		if rightIndex < len(CommandTable) {
			right := formatHelpEntry(CommandTable[rightIndex])
			msgs = append(msgs, fmt.Sprintf("%-40s %s", left, right))
			continue
		}
		msgs = append(msgs, left)
	}
	return msgs
}

func formatHelpEntry(c CommandSpec) string {
	marker := "*" // FREE: does not end the turn.
	if c.Turn {
		marker = " "
	}
	return fmt.Sprintf("%2d: %s %-32s", c.Code, marker, c.Description)
}

var numberRe = regexp.MustCompile(`^[+-]?[0-9.]+$`)

// ParseNumeric parses a raw numeric-form command line (e.g. "1 all
// 20"), ported from the non-PARSER half of scancmd()/parsit() in
// misc.c: the first whitespace-separated token is the command code,
// and everything else is passed through as args for the specific
// command handler (in cli-wiring) to interpret.
//
// The original resolved the code via a prefix (strncmp) match against
// cmds[].code_num, which happens to behave identically to exact
// integer parsing for every currently-valid command number (see
// CommandTable) since ambiguous prefixes like "1" are themselves valid
// codes appearing earlier in the table than any code with "1" as a
// prefix ("10"-"19"). This port therefore simplifies it to a direct,
// unambiguous atoi rather than replicating the fragile prefix scan.
func ParseNumeric(input string) (code int, args []string, ok bool) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return 0, nil, false
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, false
	}
	if _, exists := LookupCommand(n); !exists {
		return 0, nil, false
	}
	rest := fields[1:]
	if len(rest) == 0 {
		rest = nil
	}
	return n, rest, true
}

// ParseCommand routes raw player input to either the numeric parser or
// the text parser, ported from the "if (isalpha(*ch)) yyparse(); else
// strcpy(parsed, Input);" dispatch in playit() (main.c). Leading
// whitespace is skipped before checking the first character, matching
// the original.
func ParseCommand(input string) (code int, args []string, ok bool) {
	trimmed := strings.TrimLeft(input, " \t")
	if trimmed == "" {
		return 0, nil, false
	}
	if isAlpha(trimmed[0]) {
		return ParseText(input)
	}
	return ParseNumeric(input)
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
