package game

// RunTurn advances the battle by one full turn: enemy strategy
// decisions, per-ship power distribution, movement/combat simulation,
// target/lock invalidation, misc timers, and finally disposition
// (win/loss/surrender resolution). It is the Go equivalent of
// alarmtrap() in main.c, with the "for (i=0;i<=shipnum;i++)
// distribute(shiplist[i])" loop that move_ships() used to open with
// folded in explicitly (see MoveShips's doc comment).
//
// Newly-launched objects (e.g. torpedoes fired this turn) are appended
// to st.Objects, and detonated objects are removed from it, so the
// caller does not need to manage hooks.LaunchedObjects/
// DetonatedObjects itself.
func RunTurn(st *State, hooks *CombatHooks) []string {
	var messages []string

	for _, sp := range st.Enemies() {
		hooks.StandardStrategy(sp)
	}
	messages = append(messages, hooks.TakeMessages()...)

	for _, sp := range st.Ships {
		messages = append(messages, Distribute(sp, st, hooks)...)
	}

	messages = append(messages, MoveShips(st, hooks)...)

	if len(hooks.LaunchedObjects) > 0 {
		st.Objects = append(st.Objects, hooks.LaunchedObjects...)
		hooks.LaunchedObjects = nil
	}
	if len(hooks.DetonatedObjects) > 0 {
		detonated := make(map[*SpaceObject]bool, len(hooks.DetonatedObjects))
		for _, obj := range hooks.DetonatedObjects {
			detonated[obj] = true
		}
		remaining := st.Objects[:0]
		for _, obj := range st.Objects {
			if !detonated[obj] {
				remaining = append(remaining, obj)
			}
		}
		st.Objects = remaining
		hooks.DetonatedObjects = nil
	}

	messages = append(messages, CheckTargets(st)...)
	messages = append(messages, MiscTimers(st)...)
	messages = append(messages, Disposition(st)...)
	return messages
}
