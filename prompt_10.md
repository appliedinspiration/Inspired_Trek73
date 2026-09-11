The save/restore feature is not currently implemented. Suggest an implementation similar to the original/reference mechanism (as seen from the user's point of view). The implementation/details need not follow the c-code.

The plan:
# Save/restore implementation plan

## Problem
The reference Trek73 game saves an in-memory process image to disk, but the current Go port has no save/restore path at all. The project needs a portable, explicit save format that matches the original UX of saving the current game and later resuming it from a default or custom path.

## Proposed approach
Implement a small, explicit save subsystem that snapshots the active game state to a versioned file, then restores it on startup or on demand. The user experience should feel similar to the C game: a default save path is chosen automatically, a save command can reuse or overwrite it, and restore can happen either from the default file or a user-supplied path.

This will not reproduce the original opaque memory-image save; instead it will use a portable format (likely JSON or a small binary schema) that captures enough game state to resume safely and deterministically. The save layer will live outside the gameplay logic, with a validation/version guard to reject incompatible files.

## Planned work
- Review the current game state model and CLI session workflow to identify exactly which fields must be serialized.
- Add a save package (or equivalent helper) that converts the live `game.State` and RNG/session metadata into a versioned snapshot.
- Add save and restore commands to the CLI flow using the current command table and session loop, with default-file prompts similar to the original UX.
- Add startup support for `-restore` or equivalent flag to resume from a saved game before entering the main command loop.
- Add deterministic tests for serialization, version mismatch rejection, and successful restore semantics.

## Key files likely involved
- `cmd/inspired_trek73/main.go` — CLI flags and startup flow
- `cmd/inspired_trek73/session.go` — user interaction and command dispatch
- `internal/game/state.go` — canonical game state to persist
- `internal/game/init.go` and related gameplay files — validate which state fields are required for a valid resume
- New save helper under `internal/` or `cmd/` for format + versioning logic

## Decisions and constraints
- Keep the save file explicitly versioned and portable; avoid reusing the original binary-image trick.
- Persist only the minimal game state needed to resume: ship state, object state, turn counters, status flags, crew names, and RNG seed/derived state.
- Default save location should be user-friendly and platform-aware, similar to `$HOME/trek73.save` in the original.
- If a saved file is from an incompatible version, fail cleanly with a clear message instead of continuing.
- The save/restore logic should be isolated from gameplay rules so it can be tested independently of the full command loop.

## Todo breakdown
1. Design the save format and default-path UX.
2. Implement state serialization/deserialization.
3. Add save/restore commands to the CLI and startup flow.
4. Add tests for round-trip saves and invalid-version rejects.
