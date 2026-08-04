package data

import "testing"

func TestShipClassesComplete(t *testing.T) {
	const want = 4
	if len(ShipClasses) != want {
		t.Fatalf("len(ShipClasses) = %d, want %d", len(ShipClasses), want)
	}
	seen := map[string]bool{}
	for _, c := range ShipClasses {
		if c.Abbr == "" {
			t.Errorf("ship class %+v has empty abbreviation", c)
		}
		if seen[c.Abbr] {
			t.Errorf("duplicate ship class abbreviation %q", c.Abbr)
		}
		seen[c.Abbr] = true
		if c.NumPhaser < 1 || c.NumTorp < 1 {
			t.Errorf("ship class %s has invalid weapon counts: phasers=%d torps=%d", c.Abbr, c.NumPhaser, c.NumTorp)
		}
	}
	for _, abbr := range []string{"DN", "CA", "CL", "DD"} {
		if !seen[abbr] {
			t.Errorf("missing expected ship class %q", abbr)
		}
	}
}

func TestShipClassByAbbr(t *testing.T) {
	c, ok := ShipClassByAbbr("CA")
	if !ok {
		t.Fatal("ShipClassByAbbr(\"CA\") not found")
	}
	if c.NumPhaser != 4 || c.NumTorp != 6 {
		t.Errorf("CA class = %+v, want NumPhaser=4 NumTorp=6", c)
	}

	if _, ok := ShipClassByAbbr("ZZ"); ok {
		t.Error("ShipClassByAbbr(\"ZZ\") unexpectedly found")
	}
}

func TestRacesComplete(t *testing.T) {
	const wantRaces = 9
	if len(Races) != wantRaces {
		t.Fatalf("len(Races) = %d, want %d", len(Races), wantRaces)
	}
	for i, r := range Races {
		if r.Name == "" {
			t.Errorf("race %d has empty name", i)
		}
		if len(r.ShipNames) != MaxEnemyShipsPerRace {
			t.Errorf("race %s: len(ShipNames) = %d, want %d", r.Name, len(r.ShipNames), MaxEnemyShipsPerRace)
		}
		if len(r.ShipTypes) != 4 {
			t.Errorf("race %s: len(ShipTypes) = %d, want 4", r.Name, len(r.ShipTypes))
		}
		if len(r.Captains) != MaxCaptainsPerRace {
			t.Errorf("race %s: len(Captains) = %d, want %d", r.Name, len(r.Captains), MaxCaptainsPerRace)
		}
	}
	if Races[MontyPythonRaceIndex].Name != "Monty Python" {
		t.Errorf("Races[MontyPythonRaceIndex] = %q, want \"Monty Python\"", Races[MontyPythonRaceIndex].Name)
	}
}

func TestFederationShipNames(t *testing.T) {
	const want = 9
	if len(FederationShipNames) != want {
		t.Fatalf("len(FederationShipNames) = %d, want %d", len(FederationShipNames), want)
	}
}

func TestDamageProfiles(t *testing.T) {
	for name, dp := range map[string]DamageProfile{"PhaserDamage": PhaserDamage, "AntimatterDamage": AntimatterDamage} {
		for i, sys := range dp.Systems {
			if sys.Roll <= 0 {
				t.Errorf("%s.Systems[%d].Roll = %d, want > 0", name, i, sys.Roll)
			}
			if sys.Message == "" {
				t.Errorf("%s.Systems[%d].Message is empty", name, i)
			}
		}
	}
}

func TestSystemMessageTables(t *testing.T) {
	if len(SystemNames) != 4 {
		t.Errorf("len(SystemNames) = %d, want 4", len(SystemNames))
	}
	if len(SystemDestroyedMessages) != 5 {
		t.Errorf("len(SystemDestroyedMessages) = %d, want 5", len(SystemDestroyedMessages))
	}
}

func TestWeaponBearingTables(t *testing.T) {
	for count := 1; count <= 10; count++ {
		if got := len(PhaserInitialBearings[count]); got != count {
			t.Errorf("PhaserInitialBearings[%d] has %d entries, want %d", count, got, count)
		}
		if got := len(TubeInitialBearings[count]); got != count {
			t.Errorf("TubeInitialBearings[%d] has %d entries, want %d", count, got, count)
		}
	}
}
