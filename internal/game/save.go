package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const saveFileVersion = 1

// DefaultSavePath returns the platform-appropriate default location for a
// saved Trek73 game, mirroring the original's $HOME/trek73.save behavior.
func DefaultSavePath() string {
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return filepath.Join(home, "trek73.save")
	}
	return filepath.Join(".", "trek73.save")
}

// SaveGame writes the current game state and RNG position to disk as a
// portable, versioned snapshot.
func SaveGame(path string, st *State, r *Rand) error {
	if path == "" {
		path = DefaultSavePath()
	}
	payload := saveFile{
		Version: saveFileVersion,
		State:   snapshotState(st),
		Rand:    r.Snapshot(),
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("game: marshal save state: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("game: create save directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("game: write save file: %w", err)
	}
	return nil
}

// LoadGame reads a saved game snapshot and restores both the state and RNG
// sequence so the game can continue from the saved point.
func LoadGame(path string) (*State, *Rand, error) {
	if path == "" {
		path = DefaultSavePath()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("game: read save file: %w", err)
	}

	var saved saveFile
	if err := json.Unmarshal(data, &saved); err != nil {
		return nil, nil, fmt.Errorf("game: parse save file: %w", err)
	}
	if saved.Version != saveFileVersion {
		return nil, nil, fmt.Errorf("game: unsupported save file version %d", saved.Version)
	}

	st := snapshotToState(saved.State)
	r := NewRand(saved.Rand.Seed)
	r.Restore(saved.Rand)
	return st, r, nil
}

type saveFile struct {
	Version int          `json:"version"`
	Rand    RandSnapshot `json:"rand"`
	State   stateSnapshot `json:"state"`
}

type stateSnapshot struct {
	Ships             []shipSnapshot    `json:"ships"`
	Objects           []objectSnapshot  `json:"objects"`
	PlayStatus        int               `json:"playStatus"`
	Crew              CrewNames         `json:"crew"`
	EnemyRaceName     string            `json:"enemyRaceName"`
	EnemyEmpireName   string            `json:"enemyEmpireName"`
	EnemyShipTypeName string            `json:"enemyShipTypeName"`
	EnemyCommander    string            `json:"enemyCommander"`
	RaceChances       RaceChances       `json:"raceChances"`
	Shutup            shutupSnapshot    `json:"shutup"`
	Defenseless       int               `json:"defenseless"`
	Corbomite         int               `json:"corbomite"`
	Surrender         int               `json:"surrender"`
	SurrenderP        int               `json:"surrenderP"`
	Reengaged         bool              `json:"reengaged"`
	PendingFinal      int               `json:"pendingFinal"`
	PendingWarn       int               `json:"pendingWarn"`
	WarnShown         [6]bool           `json:"warnShown"`
	NextID            int               `json:"nextID"`
}

type shutupSnapshot struct {
	Disengage          bool              `json:"disengage"`
	ShieldsFluctuating bool              `json:"shieldsFluctuating"`
	Surrender          bool              `json:"surrender"`
	Phaser             [MaxWeapons]bool  `json:"phaser"`
	Tube               [MaxWeapons]bool  `json:"tube"`
	Burnout            map[int]bool      `json:"burnout"`
}

type shipSnapshot struct {
	ID                   int               `json:"id"`
	Name                 string            `json:"name"`
	Class                string            `json:"class"`
	X                    int               `json:"x"`
	Y                    int               `json:"y"`
	Warp                 float64           `json:"warp"`
	NewWarp              float64           `json:"newWarp"`
	Course               float64           `json:"course"`
	NewCourse            float64           `json:"newCourse"`
	TargetID             int               `json:"targetID"`
	RelativeBear         float64           `json:"relativeBear"`
	Phasers              []phaserSnapshot  `json:"phasers"`
	PhaserSpread         int               `json:"phaserSpread"`
	PhaserFirePct        int               `json:"phaserFirePct"`
	PhaserBlindLeft      int               `json:"phaserBlindLeft"`
	PhaserBlindRight     int               `json:"phaserBlindRight"`
	Tubes                []tubeSnapshot    `json:"tubes"`
	TubeProximity        int               `json:"tubeProximity"`
	TubeDelay            int               `json:"tubeDelay"`
	TubeLaunchSpd        int               `json:"tubeLaunchSpd"`
	TubeBlindLeft        int               `json:"tubeBlindLeft"`
	TubeBlindRight       int               `json:"tubeBlindRight"`
	Shields              [NumShields]shieldSnapshot `json:"shields"`
	ProbeLauncherStatus  int               `json:"probeLauncherStatus"`
	Eff                  float64           `json:"eff"`
	Regen                float64           `json:"regen"`
	Energy               float64           `json:"energy"`
	Pods                 float64           `json:"pods"`
	Complement           int               `json:"complement"`
	Status               [MaxSystems]int   `json:"status"`
	Delay                float64           `json:"delay"`
	OrigMaxSpeed         float64           `json:"origMaxSpeed"`
	MaxSpeed             float64           `json:"maxSpeed"`
	DegPerTurn           float64           `json:"degPerTurn"`
	PhaserShieldDivisor  float64           `json:"phaserShieldDivisor"`
	TorpedoShieldDivisor float64           `json:"torpedoShieldDivisor"`
	Cloaking             int               `json:"cloaking"`
	CloakEnergy          int               `json:"cloakEnergy"`
	CloakDelay           int               `json:"cloakDelay"`
	Strategy             string            `json:"strategy"`
	LastKnown            LastKnownPosition `json:"lastKnown"`
	PhaserFiringDelay    int               `json:"phaserFiringDelay"`
	TorpedoFiringDelay   int               `json:"torpedoFiringDelay"`
}

type phaserSnapshot struct {
	TargetID int     `json:"targetID"`
	Bearing  float64 `json:"bearing"`
	Drain    int     `json:"drain"`
	Load     float64 `json:"load"`
	Status   int     `json:"status"`
}

type tubeSnapshot struct {
	TargetID int     `json:"targetID"`
	Bearing  float64 `json:"bearing"`
	Load     float64 `json:"load"`
	Status   int     `json:"status"`
}

type shieldSnapshot struct {
	Eff          float64 `json:"eff"`
	Drain        float64 `json:"drain"`
	AttemptDrain float64 `json:"attemptDrain"`
}

type objectSnapshot struct {
	ID        int     `json:"id"`
	Type      int     `json:"type"`
	FromID    int     `json:"fromID"`
	X         int     `json:"x"`
	Y         int     `json:"y"`
	Course    float64 `json:"course"`
	Speed     float64 `json:"speed"`
	NewSpeed  float64 `json:"newSpeed"`
	TargetID  int     `json:"targetID"`
	Fuel      int     `json:"fuel"`
	TimeDelay float64 `json:"timeDelay"`
	Proximity int     `json:"proximity"`
	Detonated bool    `json:"detonated"`
}

func snapshotState(st *State) stateSnapshot {
	out := stateSnapshot{
		PlayStatus:        st.PlayStatus,
		Crew:              st.Crew,
		EnemyRaceName:     st.EnemyRaceName,
		EnemyEmpireName:   st.EnemyEmpireName,
		EnemyShipTypeName: st.EnemyShipTypeName,
		EnemyCommander:    st.EnemyCommander,
		RaceChances:       st.RaceChances,
		Defenseless:       st.Defenseless,
		Corbomite:         st.Corbomite,
		Surrender:         st.Surrender,
		SurrenderP:        st.SurrenderP,
		Reengaged:         st.Reengaged,
		PendingFinal:      st.PendingFinal,
		PendingWarn:       st.PendingWarn,
		WarnShown:         st.WarnShown,
		NextID:            st.nextID,
	}
	if st.Shutup != nil {
		out.Shutup = shutupSnapshot{
			Disengage:          st.Shutup.Disengage,
			ShieldsFluctuating: st.Shutup.ShieldsFluctuating,
			Surrender:          st.Shutup.Surrender,
			Phaser:             st.Shutup.Phaser,
			Tube:               st.Shutup.Tube,
			Burnout:            make(map[int]bool, len(st.Shutup.Burnout)),
		}
		for k, v := range st.Shutup.Burnout {
			out.Shutup.Burnout[k] = v
		}
	}
	for _, ship := range st.Ships {
		out.Ships = append(out.Ships, snapshotShip(ship))
	}
	for _, obj := range st.Objects {
		out.Objects = append(out.Objects, snapshotObject(obj))
	}
	return out
}

func snapshotShip(ship *Ship) shipSnapshot {
	out := shipSnapshot{
		ID:                   ship.ID,
		Name:                 ship.Name,
		Class:                ship.Class,
		X:                    ship.X,
		Y:                    ship.Y,
		Warp:                 ship.Warp,
		NewWarp:              ship.NewWarp,
		Course:               ship.Course,
		NewCourse:            ship.NewCourse,
		RelativeBear:         ship.RelativeBear,
		PhaserSpread:         ship.PhaserSpread,
		PhaserFirePct:        ship.PhaserFirePct,
		PhaserBlindLeft:      ship.PhaserBlindLeft,
		PhaserBlindRight:     ship.PhaserBlindRight,
		TubeProximity:        ship.TubeProximity,
		TubeDelay:            ship.TubeDelay,
		TubeLaunchSpd:        ship.TubeLaunchSpd,
		TubeBlindLeft:        ship.TubeBlindLeft,
		TubeBlindRight:       ship.TubeBlindRight,
		ProbeLauncherStatus:  ship.ProbeLauncherStatus,
		Eff:                  ship.Eff,
		Regen:                ship.Regen,
		Energy:               ship.Energy,
		Pods:                 ship.Pods,
		Complement:           ship.Complement,
		Delay:                ship.Delay,
		OrigMaxSpeed:         ship.OrigMaxSpeed,
		MaxSpeed:             ship.MaxSpeed,
		DegPerTurn:           ship.DegPerTurn,
		PhaserShieldDivisor:  ship.PhaserShieldDivisor,
		TorpedoShieldDivisor: ship.TorpedoShieldDivisor,
		Cloaking:             ship.Cloaking,
		CloakEnergy:          ship.CloakEnergy,
		CloakDelay:           ship.CloakDelay,
		Strategy:             ship.Strategy,
		LastKnown:            ship.LastKnown,
		PhaserFiringDelay:    ship.PhaserFiringDelay,
		TorpedoFiringDelay:   ship.TorpedoFiringDelay,
		Status:               ship.Status,
	}
	for i, shield := range ship.Shields {
	out.Shields[i] = shieldSnapshot{
		Eff:          shield.Eff,
		Drain:        shield.Drain,
		AttemptDrain: shield.AttemptDrain,
	}
	}
	if ship.Target != nil {
		out.TargetID = ship.Target.ID
	}
	for _, bank := range ship.Phasers {
		entry := phaserSnapshot{
			Bearing: bank.Bearing,
			Drain:   bank.Drain,
			Load:    bank.Load,
			Status:  bank.Status,
		}
		if bank.Target != nil {
			entry.TargetID = bank.Target.ID
		}
		out.Phasers = append(out.Phasers, entry)
	}
	for _, tube := range ship.Tubes {
		entry := tubeSnapshot{
			Bearing: tube.Bearing,
			Load:    tube.Load,
			Status:  tube.Status,
		}
		if tube.Target != nil {
			entry.TargetID = tube.Target.ID
		}
		out.Tubes = append(out.Tubes, entry)
	}
	return out
}

func snapshotObject(obj *SpaceObject) objectSnapshot {
	out := objectSnapshot{
		ID:        obj.ID,
		Type:      obj.Type,
		X:         obj.X,
		Y:         obj.Y,
		Course:    obj.Course,
		Speed:     obj.Speed,
		NewSpeed:  obj.NewSpeed,
		Fuel:      obj.Fuel,
		TimeDelay: obj.TimeDelay,
		Proximity: obj.Proximity,
		Detonated: obj.Detonated,
	}
	if obj.From != nil {
		out.FromID = obj.From.ID
	}
	if obj.Target != nil {
		out.TargetID = obj.Target.ID
	}
	return out
}

func snapshotToState(saved stateSnapshot) *State {
	st := NewState()
	st.PlayStatus = saved.PlayStatus
	st.Crew = saved.Crew
	st.EnemyRaceName = saved.EnemyRaceName
	st.EnemyEmpireName = saved.EnemyEmpireName
	st.EnemyShipTypeName = saved.EnemyShipTypeName
	st.EnemyCommander = saved.EnemyCommander
	st.RaceChances = saved.RaceChances
	st.Defenseless = saved.Defenseless
	st.Corbomite = saved.Corbomite
	st.Surrender = saved.Surrender
	st.SurrenderP = saved.SurrenderP
	st.Reengaged = saved.Reengaged
	st.PendingFinal = saved.PendingFinal
	st.PendingWarn = saved.PendingWarn
	st.WarnShown = saved.WarnShown
	st.nextID = saved.NextID
	st.Shutup = NewShutup()
	st.Shutup.Disengage = saved.Shutup.Disengage
	st.Shutup.ShieldsFluctuating = saved.Shutup.ShieldsFluctuating
	st.Shutup.Surrender = saved.Shutup.Surrender
	st.Shutup.Phaser = saved.Shutup.Phaser
	st.Shutup.Tube = saved.Shutup.Tube
	st.Shutup.Burnout = make(map[int]bool, len(saved.Shutup.Burnout))
	for k, v := range saved.Shutup.Burnout {
		st.Shutup.Burnout[k] = v
	}

	shipByID := make(map[int]*Ship, len(saved.Ships))
	for _, savedShip := range saved.Ships {
		ship := &Ship{
			ID:                   savedShip.ID,
			Name:                 savedShip.Name,
			Class:                savedShip.Class,
			X:                    savedShip.X,
			Y:                    savedShip.Y,
			Warp:                 savedShip.Warp,
			NewWarp:              savedShip.NewWarp,
			Course:               savedShip.Course,
			NewCourse:            savedShip.NewCourse,
			RelativeBear:         savedShip.RelativeBear,
			PhaserSpread:         savedShip.PhaserSpread,
			PhaserFirePct:        savedShip.PhaserFirePct,
			PhaserBlindLeft:      savedShip.PhaserBlindLeft,
			PhaserBlindRight:     savedShip.PhaserBlindRight,
			TubeProximity:        savedShip.TubeProximity,
			TubeDelay:            savedShip.TubeDelay,
			TubeLaunchSpd:        savedShip.TubeLaunchSpd,
			TubeBlindLeft:        savedShip.TubeBlindLeft,
			TubeBlindRight:       savedShip.TubeBlindRight,
			ProbeLauncherStatus:  savedShip.ProbeLauncherStatus,
			Eff:                  savedShip.Eff,
			Regen:                savedShip.Regen,
			Energy:               savedShip.Energy,
			Pods:                 savedShip.Pods,
			Complement:           savedShip.Complement,
			Delay:                savedShip.Delay,
			OrigMaxSpeed:         savedShip.OrigMaxSpeed,
			MaxSpeed:             savedShip.MaxSpeed,
			DegPerTurn:           savedShip.DegPerTurn,
			PhaserShieldDivisor:  savedShip.PhaserShieldDivisor,
			TorpedoShieldDivisor: savedShip.TorpedoShieldDivisor,
			Cloaking:             savedShip.Cloaking,
			CloakEnergy:          savedShip.CloakEnergy,
			CloakDelay:           savedShip.CloakDelay,
			Strategy:             savedShip.Strategy,
			LastKnown:            savedShip.LastKnown,
			PhaserFiringDelay:    savedShip.PhaserFiringDelay,
			TorpedoFiringDelay:   savedShip.TorpedoFiringDelay,
			Status:               savedShip.Status,
		}
		for i, shield := range savedShip.Shields {
			ship.Shields[i] = Shield{Eff: shield.Eff, Drain: shield.Drain, AttemptDrain: shield.AttemptDrain}
		}
		ship.Phasers = make([]Phaser, len(savedShip.Phasers))
		for i, bank := range savedShip.Phasers {
			ship.Phasers[i] = Phaser{Bearing: bank.Bearing, Drain: bank.Drain, Load: bank.Load, Status: bank.Status}
		}
		ship.Tubes = make([]Tube, len(savedShip.Tubes))
		for i, tube := range savedShip.Tubes {
			ship.Tubes[i] = Tube{Bearing: tube.Bearing, Load: tube.Load, Status: tube.Status}
		}
		shipByID[ship.ID] = ship
		st.Ships = append(st.Ships, ship)
	}
	for i, savedShip := range saved.Ships {
		ship := st.Ships[i]
		if savedShip.TargetID != 0 {
			if target, ok := shipByID[savedShip.TargetID]; ok {
				ship.Target = target
			}
		}
		for j, bank := range savedShip.Phasers {
			if targetID := bank.TargetID; targetID != 0 {
				if target, ok := shipByID[targetID]; ok {
					ship.Phasers[j].Target = target
				}
			}
		}
		for j, tube := range savedShip.Tubes {
			if targetID := tube.TargetID; targetID != 0 {
				if target, ok := shipByID[targetID]; ok {
					ship.Tubes[j].Target = target
				}
			}
		}
	}
	for _, savedObj := range saved.Objects {
		obj := &SpaceObject{
			ID:        savedObj.ID,
			Type:      savedObj.Type,
			X:         savedObj.X,
			Y:         savedObj.Y,
			Course:    savedObj.Course,
			Speed:     savedObj.Speed,
			NewSpeed:  savedObj.NewSpeed,
			Fuel:      savedObj.Fuel,
			TimeDelay: savedObj.TimeDelay,
			Proximity: savedObj.Proximity,
			Detonated: savedObj.Detonated,
		}
		if savedObj.FromID != 0 {
			if from, ok := shipByID[savedObj.FromID]; ok {
				obj.From = from
			}
		}
		if savedObj.TargetID != 0 {
			if target, ok := shipByID[savedObj.TargetID]; ok {
				obj.Target = target
			}
		}
		st.Objects = append(st.Objects, obj)
	}
	return st
}
