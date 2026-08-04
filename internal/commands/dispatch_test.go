package commands

import (
	"reflect"
	"strings"
	"testing"
)

func TestCommandTable_Has32Entries(t *testing.T) {
	if len(CommandTable) != 32 {
		t.Fatalf("len(CommandTable) = %d, want 32", len(CommandTable))
	}
	for i, c := range CommandTable {
		if c.Code != i+1 {
			t.Errorf("CommandTable[%d].Code = %d, want %d", i, c.Code, i+1)
		}
	}
}

func TestLookupCommand(t *testing.T) {
	c, ok := LookupCommand(1)
	if !ok || c.Description != "Fire phasers" {
		t.Errorf("LookupCommand(1) = %+v, %v", c, ok)
	}
	if _, ok := LookupCommand(0); ok {
		t.Error("LookupCommand(0) should not exist")
	}
	if _, ok := LookupCommand(33); ok {
		t.Error("LookupCommand(33) should not exist")
	}
}

func TestVersion(t *testing.T) {
	msgs := Version()
	if len(msgs) == 0 {
		t.Error("expected version output")
	}
}

func TestHelp(t *testing.T) {
	msgs := Help()
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "Fire phasers") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected help to list 'Fire phasers', got %v", msgs)
	}
}

func TestParseNumeric(t *testing.T) {
	tests := []struct {
		input    string
		wantCode int
		wantArgs []string
		wantOk   bool
	}{
		{"1 all 20", 1, []string{"all", "20"}, true},
		{"9 l all", 9, []string{"l", "all"}, true},
		{"32", 32, nil, true},
		{"99", 0, nil, false},
		{"", 0, nil, false},
		{"notanumber", 0, nil, false},
	}
	for _, tt := range tests {
		code, args, ok := ParseNumeric(tt.input)
		if code != tt.wantCode || ok != tt.wantOk || !reflect.DeepEqual(args, tt.wantArgs) {
			t.Errorf("ParseNumeric(%q) = %d, %v, %v; want %d, %v, %v",
				tt.input, code, args, ok, tt.wantCode, tt.wantArgs, tt.wantOk)
		}
	}
}

func TestParseCommand_RoutesNumericVsText(t *testing.T) {
	code, _, ok := ParseCommand("1 20")
	if !ok || code != 1 {
		t.Errorf("ParseCommand numeric = %d, %v", code, ok)
	}
	code, _, ok = ParseCommand("fire phasers 1")
	if !ok || code != 1 {
		t.Errorf("ParseCommand text = %d, %v", code, ok)
	}
	if _, _, ok := ParseCommand(""); ok {
		t.Error("ParseCommand(\"\") should fail")
	}
	if _, _, ok := ParseCommand("   "); ok {
		t.Error("ParseCommand(whitespace) should fail")
	}
}

func TestParseText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCode int
		wantArgs []string
		wantOk   bool
	}{
		{"fire phaser spread", "fire phasers 1 spread 20", 1, []string{"1", "20"}, true},
		{"fire phaser only", "fire phaser 1", 1, []string{"1"}, true},
		{"fire all phaser spread", "fire all phasers spread 20", 1, []string{"all", "20"}, true},
		{"fire all phaser", "fire all phaser", 1, []string{"all"}, true},
		{"fire phaser bare", "fire phaser", 1, nil, true},
		{"fire tube", "fire tube 2", 2, []string{"2"}, true},
		{"fire all tube", "fire all tube", 2, []string{"all"}, true},
		{"fire tube bare", "fire tube", 2, nil, true},

		{"lock phaser with name", "lock phaser 1 enterprise", 3, []string{"1", "enterprise"}, true},
		{"lock phaser no name", "lock phaser 1", 3, []string{"1"}, true},
		{"lock all phaser with name", "lock all phaser enterprise", 3, []string{"all", "enterprise"}, true},
		{"lock all phaser", "lock all phaser", 3, []string{"all"}, true},
		{"lock phaser bare", "lock phaser", 3, nil, true},
		{"lock tube with name", "lock tube 1 enterprise", 4, []string{"1", "enterprise"}, true},
		{"lock all tube", "lock all tube", 4, []string{"all"}, true},

		{"rotate all phaser num", "rotate all phaser 90", 5, []string{"all", "90"}, true},
		{"rotate all phaser", "rotate all phaser", 5, []string{"all"}, true},
		{"rotate phaser two nums", "rotate phaser 1 90", 5, []string{"1", "90"}, true},
		{"rotate phaser one num", "rotate phaser 1", 5, []string{"1"}, true},
		{"rotate all tube num", "rotate all tube 90", 6, []string{"all", "90"}, true},
		{"rotate tube two nums", "rotate tube 1 90", 6, []string{"1", "90"}, true},

		{"phaser status", "phaser status", 7, nil, true},
		{"tube status", "tube status", 8, nil, true},

		{"load all tube", "load all tube", 9, []string{"l", "all"}, true},
		{"unload all tube", "unload all tube", 9, []string{"u", "all"}, true},
		{"load tube num", "load tube 1", 9, []string{"l", "1"}, true},
		{"unload tube num", "unload tube 1", 9, []string{"u", "1"}, true},
		{"load tube bare", "load tube", 9, []string{"l"}, true},
		{"unload tube bare", "unload tube", 9, []string{"u"}, true},

		{"launch probe bare", "launch probe", 10, nil, true},
		{"launch probe num", "launch probe 5", 10, []string{"5"}, true},
		{"launch probe delay", "launch probe 5 delay 10", 10, []string{"5", "10"}, true},
		{"launch probe prox", "launch probe 5 delay 10 proximity 2", 10, []string{"5", "10", "2"}, true},
		{"launch probe toward", "launch probe 5 delay 10 proximity 2 toward enterprise", 10, []string{"5", "10", "2", "enterprise"}, true},
		{"launch probe course", "launch probe 5 delay 10 proximity 2 course 90", 10, []string{"5", "10", "2", "90"}, true},

		{"control probe num", "control probe 1", 11, []string{"1"}, true},
		{"control probe bare", "control probe", 11, nil, true},
		{"probe control num", "probe control 1", 11, []string{"1"}, true},
		{"probe control bare", "probe control", 11, nil, true},

		{"tactical", "tactical", 12, nil, true},

		{"display num", "display 3", 13, []string{"3"}, true},
		{"display bare", "display", 13, nil, true},

		{"pursue name warp", "pursue enterprise warp 5", 14, []string{"enterprise", "5"}, true},
		{"pursue name num", "pursue enterprise 5", 14, []string{"enterprise", "5"}, true},
		{"pursue name", "pursue enterprise", 14, []string{"enterprise"}, true},
		{"pursue bare", "pursue", 14, nil, true},

		{"elude name warp", "elude enterprise warp 5", 15, []string{"enterprise", "5"}, true},
		{"elude bare", "elude", 15, nil, true},

		{"course warp", "course 90 warp 5", 16, []string{"90", "5"}, true},
		{"course two nums", "course 90 5", 16, []string{"90", "5"}, true},
		{"course one num", "course 90", 16, []string{"90"}, true},
		{"course course warp", "course course 90 warp 5", 16, []string{"90", "5"}, true},
		{"turn synonym", "turn 90 warp 5", 16, []string{"90", "5"}, true},

		{"damage", "damage", 17, nil, true},
		{"damage report phrase", "damage report", 17, nil, true},

		{"scan name", "scan enterprise", 18, []string{"enterprise"}, true},
		{"scan num", "scan 3", 18, []string{"3"}, true},
		{"scan bare", "scan", 18, nil, true},

		{"power", "power", 19, nil, true},
		{"alter power phrase", "alter power", 19, nil, true},

		{"jettison eng", "jettison eng", 20, nil, true},
		{"jett eng synonym", "jett eng", 20, nil, true},

		{"detonate eng name", "detonate eng enterprise", 21, []string{"enterprise"}, true},
		{"detonate eng bare", "detonate eng", 21, nil, true},

		{"param", "param", 22, nil, true},
		{"parameters synonym", "parameters", 22, nil, true},

		{"dead name", "dead enterprise", 23, []string{"enterprise"}, true},
		{"dead bare", "dead", 23, nil, true},
		{"play dead phrase", "play dead", 23, nil, true},

		{"corbomite", "corbomite", 24, nil, true},
		{"corbomite bluff phrase", "corbomite bluff", 24, nil, true},

		{"surrender", "surrender", 25, nil, true},

		{"demand surrender", "demand surrender", 26, nil, true},

		{"destruct", "destruct", 27, nil, true},
		{"self destruct phrase", "self destruct", 27, nil, true},

		{"abort destruct", "abort destruct", 28, nil, true},

		{"survivors", "survivors", 29, nil, true},
		{"survivors report phrase", "survivors report", 29, nil, true},

		{"version", "version", 30, nil, true},
		{"save", "save", 31, nil, true},
		{"help", "help", 32, nil, true},

		{"garbage input", "xyzzy", 0, nil, false},
		{"empty input", "", 0, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, args, ok := ParseText(tt.input)
			if code != tt.wantCode || ok != tt.wantOk || !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("ParseText(%q) = %d, %v, %v; want %d, %v, %v",
					tt.input, code, args, ok, tt.wantCode, tt.wantArgs, tt.wantOk)
			}
		})
	}
}

func TestParseText_CaseInsensitiveKeywords(t *testing.T) {
	code, args, ok := ParseText("FIRE PHASERS 1 SPREAD 20")
	if !ok || code != 1 || !reflect.DeepEqual(args, []string{"1", "20"}) {
		t.Errorf("ParseText uppercase = %d, %v, %v", code, args, ok)
	}
}
