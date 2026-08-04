package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// tokKind enumerates the lexical token classes from command.l.
type tokKind int

const (
	tNUMBER tokKind = iota
	tNAME
	tABORT
	tALL
	tCONTROL
	tCORB
	tCOURSE
	tDAMAGE
	tDEAD
	tDELAY
	tDEMAND
	tDESTR
	tDET
	tDISPLAY
	tELUDE
	tENG
	tFIRE
	tHELP
	tJETT
	tLAUNCH
	tLOAD
	tLOCK
	tPARAM
	tPHASER
	tPOWER
	tPROBE
	tPROXIMITY
	tPURSUE
	tROTATE
	tSAVE
	tSCAN
	tSPREAD
	tSTATUS
	tSURREND
	tSURV
	tTACTIC
	tTOWARD
	tTUBE
	tUNLOAD
	tVERSION
	tWARP
)

// token is one lexed unit: either a NUMBER (with its float value), a
// NAME (with its literal text, case preserved), or a fixed keyword.
type token struct {
	kind tokKind
	num  float64
	name string
}

// keywordTokens maps every single-word keyword (and simple synonym) in
// command.l to its token kind. Multi-word keywords (e.g. "corbomite
// bluff") are merged into single pseudo-words by mergeMultiWordPhrases
// before this map is consulted.
var keywordTokens = map[string]tokKind{
	"abort":           tABORT,
	"all":             tALL,
	"control":         tCONTROL,
	"corbomitebluff":  tCORB,
	"corbomite":       tCORB,
	"turn":            tCOURSE,
	"course":          tCOURSE,
	"damagereport":    tDAMAGE,
	"damage":          tDAMAGE,
	"playdead":        tDEAD,
	"dead":            tDEAD,
	"delay":           tDELAY,
	"fuse":            tDELAY,
	"demand":          tDEMAND,
	"self-destruct":   tDESTR,
	"selfdestruct":    tDESTR,
	"destruct":        tDESTR,
	"detonate":        tDET,
	"display":         tDISPLAY,
	"elude":           tELUDE,
	"engineering":     tENG,
	"eng":             tENG,
	"fire":            tFIRE,
	"help":            tHELP,
	"jettison":        tJETT,
	"jett":            tJETT,
	"launch":          tLAUNCH,
	"load":            tLOAD,
	"lock":            tLOCK,
	"parameters":      tPARAM,
	"params":          tPARAM,
	"param":           tPARAM,
	"phasers":         tPHASER,
	"phaser":          tPHASER,
	"alterpower":      tPOWER,
	"power":           tPOWER,
	"probe":           tPROBE,
	"prox":            tPROXIMITY,
	"proximity":       tPROXIMITY,
	"pursue":          tPURSUE,
	"rotate":          tROTATE,
	"save":            tSAVE,
	"scan":            tSCAN,
	"spread":          tSPREAD,
	"status":          tSTATUS,
	"surrender":       tSURREND,
	"survivorsreport": tSURV,
	"survivors":       tSURV,
	"surv":            tSURV,
	"tactical":        tTACTIC,
	"toward":          tTOWARD,
	"towards":         tTOWARD,
	"photons":         tTUBE,
	"photon":          tTUBE,
	"torpedos":        tTUBE,
	"torpedo":         tTUBE,
	"torp":            tTUBE,
	"torps":           tTUBE,
	"tubes":           tTUBE,
	"tube":            tTUBE,
	"unload":          tUNLOAD,
	"version":         tVERSION,
	"warpfactor":      tWARP,
	"warp":            tWARP,
}

// fillerWords lists words that command.l recognizes but discards (its
// actions are empty), e.g. "at" in "fire phasers 1 at ship".
var fillerWords = map[string]bool{
	"at": true, "onto": true, "on": true, "pod": true, "pods": true,
	"the": true, "to": true, "with": true,
}

// multiWordPhrases collapses command.l's "word{WS}word" lex rules
// (which match with zero or more spaces between the two halves, i.e.
// effectively as a single lexeme) into single hyphen-free pseudo-words
// before tokenizing, since Go's tokenizer here works word-by-word
// rather than as a single combined regex the way flex did.
var multiWordPhrases = []struct {
	re          *regexp.Regexp
	replacement string
}{
	{regexp.MustCompile(`(?i)corbomite\s+bluff`), "corbomitebluff"},
	{regexp.MustCompile(`(?i)play\s+dead`), "playdead"},
	{regexp.MustCompile(`(?i)self\s+destruct`), "selfdestruct"},
	{regexp.MustCompile(`(?i)damage\s+report`), "damagereport"},
	{regexp.MustCompile(`(?i)survivors\s+report`), "survivorsreport"},
	{regexp.MustCompile(`(?i)alter\s+power`), "alterpower"},
	{regexp.MustCompile(`(?i)warp\s+factor`), "warpfactor"},
}

var numberToken = regexp.MustCompile(`^[+-]?[0-9.]+$`)
var nameToken = regexp.MustCompile(`^[A-Za-z#]+$`)

// tokenize lexes a raw text command line into tokens, ported from
// command.l's flex rules.
func tokenize(input string) []token {
	for _, mw := range multiWordPhrases {
		input = mw.re.ReplaceAllString(input, mw.replacement)
	}
	var toks []token
	for _, word := range strings.Fields(input) {
		lower := strings.ToLower(word)
		if fillerWords[lower] {
			continue
		}
		if kind, isKeyword := keywordTokens[lower]; isKeyword {
			toks = append(toks, token{kind: kind})
			continue
		}
		if numberToken.MatchString(word) {
			if v, err := strconv.ParseFloat(word, 64); err == nil {
				toks = append(toks, token{kind: tNUMBER, num: v})
				continue
			}
		}
		if nameToken.MatchString(word) {
			toks = append(toks, token{kind: tNAME, name: word})
			continue
		}
		// Anything else (stray punctuation, etc.) is silently
		// discarded, matching command.l's catch-all "." rule.
	}
	return toks
}

// parser walks a token slice with a single lookahead position, in
// place of the yacc-generated LALR parser driving grammar.y.
type parser struct {
	toks []token
	pos  int
}

func (p *parser) peek() (token, bool) {
	if p.pos >= len(p.toks) {
		return token{}, false
	}
	return p.toks[p.pos], true
}

func (p *parser) peekIs(k tokKind) bool {
	t, ok := p.peek()
	return ok && t.kind == k
}

func (p *parser) next() (token, bool) {
	t, ok := p.peek()
	if ok {
		p.pos++
	}
	return t, ok
}

func (p *parser) consumeNumber() (float64, bool) {
	if t, ok := p.peek(); ok && t.kind == tNUMBER {
		p.pos++
		return t.num, true
	}
	return 0, false
}

func (p *parser) consumeName() (string, bool) {
	if t, ok := p.peek(); ok && t.kind == tNAME {
		p.pos++
		return t.name, true
	}
	return "", false
}

// formatInt renders a number the way grammar.y's "%.0f" sprintf format
// specifiers do (truncated to the nearest integer, as text).
func formatInt(v float64) string {
	return fmt.Sprintf("%.0f", v)
}

// formatFloat renders a number the way grammar.y's "%f" sprintf format
// specifier does for pursue/elude warp factors, except trimmed of
// trailing zeros for readability; the exact fixed 6-decimal-place text
// the original produced is not otherwise consumed by anything (the
// only reader is atof() in the downstream command handler), so exact
// digit-for-digit replication isn't a behavior worth preserving here.
func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// ParseText parses a natural-language-ish text command line into a
// numeric command code and its argument list, hand-written in place
// of the original's lex/yacc grammar (command.l/grammar.y), per the
// project's "no lex/yacc" constraint. It implements the same keyword
// grammar and produces the same numeric-code-plus-arguments shape as
// the original's sprintf(parsed, ...) productions.
func ParseText(input string) (code int, args []string, ok bool) {
	toks := tokenize(input)
	if len(toks) == 0 {
		return 0, nil, false
	}
	p := &parser{toks: toks}
	first, _ := p.peek()
	switch first.kind {
	case tFIRE:
		return p.parseFire()
	case tLOCK:
		return p.parseLock()
	case tROTATE:
		return p.parseRotate()
	case tPHASER:
		return p.parsePhaserStatus()
	case tTUBE:
		return p.parseTubeStatus()
	case tLOAD, tUNLOAD:
		return p.parseTubeLoad()
	case tLAUNCH:
		return p.parseLaunch()
	case tCONTROL, tPROBE:
		return p.parseProbeControl()
	case tTACTIC:
		return 12, nil, true
	case tDISPLAY:
		return p.parseDisplay()
	case tPURSUE:
		return p.parsePursueElude(14)
	case tELUDE:
		return p.parsePursueElude(15)
	case tCOURSE:
		return p.parseCourse()
	case tDAMAGE:
		return 17, nil, true
	case tSCAN:
		return p.parseScan()
	case tPOWER:
		return 19, nil, true
	case tJETT:
		return p.parseJettison()
	case tDET:
		return p.parseDetonate()
	case tPARAM:
		return 22, nil, true
	case tDEAD:
		return p.parseDead()
	case tCORB:
		return 24, nil, true
	case tSURREND:
		return 25, nil, true
	case tDEMAND:
		return p.parseDemandSurrender()
	case tDESTR:
		return 27, nil, true
	case tABORT:
		return p.parseAbortDestruct()
	case tSURV:
		return 29, nil, true
	case tVERSION:
		return 30, nil, true
	case tSAVE:
		return 31, nil, true
	case tHELP:
		return 32, nil, true
	default:
		return 0, nil, false
	}
}

// parseFire handles phfire (code 1) and tufire (code 2).
func (p *parser) parseFire() (int, []string, bool) {
	p.next() // FIRE
	switch {
	case p.peekIs(tPHASER):
		p.next()
		if num1, ok := p.consumeNumber(); ok {
			if p.peekIs(tSPREAD) {
				p.next()
				if num2, ok2 := p.consumeNumber(); ok2 {
					return 1, []string{formatInt(num1), formatInt(num2)}, true
				}
				return 0, nil, false
			}
			return 1, []string{formatInt(num1)}, true
		}
		return 1, nil, true
	case p.peekIs(tALL):
		p.next()
		switch {
		case p.peekIs(tPHASER):
			p.next()
			if p.peekIs(tSPREAD) {
				p.next()
				if num1, ok := p.consumeNumber(); ok {
					return 1, []string{"all", formatInt(num1)}, true
				}
				return 0, nil, false
			}
			return 1, []string{"all"}, true
		case p.peekIs(tTUBE):
			p.next()
			return 2, []string{"all"}, true
		}
		return 0, nil, false
	case p.peekIs(tTUBE):
		p.next()
		if num1, ok := p.consumeNumber(); ok {
			return 2, []string{formatInt(num1)}, true
		}
		return 2, nil, true
	}
	return 0, nil, false
}

// parseLock handles phlock (code 3) and tulock (code 4).
func (p *parser) parseLock() (int, []string, bool) {
	p.next() // LOCK
	switch {
	case p.peekIs(tPHASER):
		p.next()
		return p.lockRest(3)
	case p.peekIs(tALL):
		p.next()
		switch {
		case p.peekIs(tPHASER):
			p.next()
			if name, ok := p.consumeName(); ok {
				return 3, []string{"all", name}, true
			}
			return 3, []string{"all"}, true
		case p.peekIs(tTUBE):
			p.next()
			if name, ok := p.consumeName(); ok {
				return 4, []string{"all", name}, true
			}
			return 4, []string{"all"}, true
		}
		return 0, nil, false
	case p.peekIs(tTUBE):
		p.next()
		return p.lockRest(4)
	}
	return 0, nil, false
}

// lockRest parses the shared "number1 name? | number1 | <nothing>"
// tail shared by the non-ALL phlock/tulock productions.
func (p *parser) lockRest(code int) (int, []string, bool) {
	if num1, ok := p.consumeNumber(); ok {
		if name, ok2 := p.consumeName(); ok2 {
			return code, []string{formatInt(num1), name}, true
		}
		return code, []string{formatInt(num1)}, true
	}
	return code, nil, true
}

// parseRotate handles phrot (code 5) and turot (code 6).
func (p *parser) parseRotate() (int, []string, bool) {
	p.next() // ROTATE
	switch {
	case p.peekIs(tALL):
		p.next()
		switch {
		case p.peekIs(tPHASER):
			p.next()
			if num1, ok := p.consumeNumber(); ok {
				return 5, []string{"all", formatInt(num1)}, true
			}
			return 5, []string{"all"}, true
		case p.peekIs(tTUBE):
			p.next()
			if num1, ok := p.consumeNumber(); ok {
				return 6, []string{"all", formatInt(num1)}, true
			}
			return 6, []string{"all"}, true
		}
		return 0, nil, false
	case p.peekIs(tPHASER):
		p.next()
		return p.rotateRest(5)
	case p.peekIs(tTUBE):
		p.next()
		return p.rotateRest(6)
	}
	return 0, nil, false
}

// rotateRest parses the shared "number1 number2? " tail of the
// non-ALL phrot/turot productions.
func (p *parser) rotateRest(code int) (int, []string, bool) {
	num1, ok := p.consumeNumber()
	if !ok {
		return 0, nil, false
	}
	if num2, ok2 := p.consumeNumber(); ok2 {
		return code, []string{formatInt(num1), formatInt(num2)}, true
	}
	return code, []string{formatInt(num1)}, true
}

// parsePhaserStatus handles phstat: PHASER STATUS (code 7).
func (p *parser) parsePhaserStatus() (int, []string, bool) {
	p.next() // PHASER
	if p.peekIs(tSTATUS) {
		p.next()
		return 7, nil, true
	}
	return 0, nil, false
}

// parseTubeStatus handles tustat: TUBE STATUS (code 8).
func (p *parser) parseTubeStatus() (int, []string, bool) {
	p.next() // TUBE
	if p.peekIs(tSTATUS) {
		p.next()
		return 8, nil, true
	}
	return 0, nil, false
}

// parseTubeLoad handles tuload (code 9): LOAD/UNLOAD [ALL] TUBE
// [number1].
func (p *parser) parseTubeLoad() (int, []string, bool) {
	verb, _ := p.next() // LOAD or UNLOAD
	dir := "l"
	if verb.kind == tUNLOAD {
		dir = "u"
	}
	if p.peekIs(tALL) {
		p.next()
		if p.peekIs(tTUBE) {
			p.next()
			return 9, []string{dir, "all"}, true
		}
		return 0, nil, false
	}
	if p.peekIs(tTUBE) {
		p.next()
		if num1, ok := p.consumeNumber(); ok {
			return 9, []string{dir, formatInt(num1)}, true
		}
		return 9, []string{dir}, true
	}
	return 0, nil, false
}

// parseLaunch handles probe (code 10): LAUNCH PROBE [number1 [DELAY
// number2 [PROXIMITY number3 [TOWARD name | COURSE number4]]]].
func (p *parser) parseLaunch() (int, []string, bool) {
	p.next() // LAUNCH
	if !p.peekIs(tPROBE) {
		return 0, nil, false
	}
	p.next() // PROBE
	num1, ok := p.consumeNumber()
	if !ok {
		return 10, nil, true
	}
	args := []string{formatInt(num1)}
	if !p.peekIs(tDELAY) {
		return 10, args, true
	}
	p.next() // DELAY
	num2, ok := p.consumeNumber()
	if !ok {
		return 0, nil, false
	}
	args = append(args, formatInt(num2))
	if !p.peekIs(tPROXIMITY) {
		return 10, args, true
	}
	p.next() // PROXIMITY
	num3, ok := p.consumeNumber()
	if !ok {
		return 0, nil, false
	}
	args = append(args, formatInt(num3))
	switch {
	case p.peekIs(tTOWARD):
		p.next()
		name, ok := p.consumeName()
		if !ok {
			return 0, nil, false
		}
		args = append(args, name)
	case p.peekIs(tCOURSE):
		p.next()
		num4, ok := p.consumeNumber()
		if !ok {
			return 0, nil, false
		}
		args = append(args, formatInt(num4))
	}
	return 10, args, true
}

// parseProbeControl handles control (code 11): CONTROL PROBE
// [number1] | PROBE CONTROL [number1].
func (p *parser) parseProbeControl() (int, []string, bool) {
	first, _ := p.next()
	if first.kind == tCONTROL {
		if !p.peekIs(tPROBE) {
			return 0, nil, false
		}
		p.next()
	} else {
		if !p.peekIs(tCONTROL) {
			return 0, nil, false
		}
		p.next()
	}
	if num1, ok := p.consumeNumber(); ok {
		return 11, []string{formatInt(num1)}, true
	}
	return 11, nil, true
}

// parseDisplay handles display (code 13): DISPLAY [number1].
func (p *parser) parseDisplay() (int, []string, bool) {
	p.next() // DISPLAY
	if num1, ok := p.consumeNumber(); ok {
		return 13, []string{formatInt(num1)}, true
	}
	return 13, nil, true
}

// parsePursueElude handles pursue (code 14) and elude (code 15):
// PURSUE/ELUDE [name [WARP? number1]].
func (p *parser) parsePursueElude(code int) (int, []string, bool) {
	p.next() // PURSUE or ELUDE
	name, ok := p.consumeName()
	if !ok {
		return code, nil, true
	}
	if p.peekIs(tWARP) {
		p.next()
	}
	if num1, ok := p.consumeNumber(); ok {
		return code, []string{name, formatFloat(num1)}, true
	}
	return code, []string{name}, true
}

// parseCourse handles course (code 16): COURSE [COURSE] number1 [WARP?
// number2].
func (p *parser) parseCourse() (int, []string, bool) {
	p.next() // COURSE
	if p.peekIs(tCOURSE) {
		p.next() // optional doubled COURSE keyword
	}
	num1, ok := p.consumeNumber()
	if !ok {
		return 0, nil, false
	}
	if p.peekIs(tWARP) {
		p.next()
	}
	if num2, ok2 := p.consumeNumber(); ok2 {
		return 16, []string{formatInt(num1), formatInt(num2)}, true
	}
	return 16, []string{formatInt(num1)}, true
}

// parseScan handles scan (code 18): SCAN [name | number1].
func (p *parser) parseScan() (int, []string, bool) {
	p.next() // SCAN
	if name, ok := p.consumeName(); ok {
		return 18, []string{name}, true
	}
	if num1, ok := p.consumeNumber(); ok {
		return 18, []string{formatInt(num1)}, true
	}
	return 18, nil, true
}

// parseJettison handles jettison (code 20): JETT ENG.
func (p *parser) parseJettison() (int, []string, bool) {
	p.next() // JETT
	if p.peekIs(tENG) {
		p.next()
		return 20, nil, true
	}
	return 0, nil, false
}

// parseDetonate handles detonate (code 21): DET ENG [name].
func (p *parser) parseDetonate() (int, []string, bool) {
	p.next() // DET
	if !p.peekIs(tENG) {
		return 0, nil, false
	}
	p.next() // ENG
	if name, ok := p.consumeName(); ok {
		return 21, []string{name}, true
	}
	return 21, nil, true
}

// parseDead handles dead (code 23): DEAD [name].
func (p *parser) parseDead() (int, []string, bool) {
	p.next() // DEAD
	if name, ok := p.consumeName(); ok {
		return 23, []string{name}, true
	}
	return 23, nil, true
}

// parseDemandSurrender handles esurr (code 26): DEMAND SURREND.
func (p *parser) parseDemandSurrender() (int, []string, bool) {
	p.next() // DEMAND
	if p.peekIs(tSURREND) {
		p.next()
		return 26, nil, true
	}
	return 0, nil, false
}

// parseAbortDestruct handles abort (code 28): ABORT DESTR.
func (p *parser) parseAbortDestruct() (int, []string, bool) {
	p.next() // ABORT
	if p.peekIs(tDESTR) {
		p.next()
		return 28, nil, true
	}
	return 0, nil, false
}
