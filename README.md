Inspired Trek73
-------------------
A remake of the classic/retro command line Trek73 computer game

Project Description
--------------------
Inspired Trek73 is a modern Go reimplementation of the classic 1970s/80s
command-line Star Trek combat simulation, Trek73. It aims to faithfully
preserve the original game's mechanics, balance, and command-driven feel
while replacing the original C implementation's Unix-specific and
obsolete platform dependencies with a clean, portable, idiomatic Go
codebase.

This is a from-scratch Go implementation guided by the original C source
as the authoritative gameplay specification. It is not a line-by-line
translation.

Brief History
-------------
Trek73 was originally written in HP-2000 BASIC by William K. Char, Perry
Lee, and Dan Gee. It was later rewritten in C by Dave Pare and
Christopher Williams, and subsequently corrected, completed, and enhanced
by Jeff Okamoto, Peter Yee, Matt Dillon, Dave Sharnoff, Joel Duisman, and
Roger J. Noe. A FreeBSD-era revision of that C code (maintained by Jeff
Okamoto and Peter Yee, with portability fixes by Matt Dillon) served as
the reference source for this modernization project.

Original Author Information
----------------------------
Original BASIC authors: William K. Char, Perry Lee, Dan Gee
C rewrite: Dave Pare, Christopher Williams
Major C enhancements: Jeff Okamoto, Peter Yee, Matt Dillon, Dave Sharnoff,
Joel Duisman, Roger J. Noe

The licensing terms of the original source were not confirmed from the
available source material. No license claim is made here on behalf of
the original authors; only attribution is given. If you are able to
confirm the original license, please update this section accordingly.

Modernization
-------------
Updated to modern operating systems and hardware by Peter S. Lee,
AppliedInspiration.com.

This implementation Copyright 2026 Peter S. Lee (AppliedInspiration.com)
All Rights Reserved.

This software is made available under the permissive MIT License.

This program is provided "as is" without warranties or guarantees of any
kind.

Project Status
--------------
This repository is in active development. The current state is an early
project scaffold: the Go module, CLI entry point, version/build-info
reporting, and startup/exit banner are implemented. The core game engine
(ship state, movement, combat, damage, enemy strategy, and command
handling) has not yet been ported from the original C source. See
prompt_1.md and prompt_1_answers_to_q for the full project brief and the
modernization decisions guiding this work.

Build Instructions
-------------------
Requirements: Go 1.23 or later (no third-party dependencies).

    go build -o inspired_trek73 ./cmd/inspired_trek73        (Linux/macOS)
    go build -o inspired_trek73.exe ./cmd/inspired_trek73    (Windows)

Cross-compilation uses the standard Go toolchain, e.g.:

    GOOS=linux   GOARCH=amd64 go build -o inspired_trek73        ./cmd/inspired_trek73
    GOOS=darwin  GOARCH=arm64 go build -o inspired_trek73        ./cmd/inspired_trek73
    GOOS=windows GOARCH=amd64 go build -o inspired_trek73.exe    ./cmd/inspired_trek73

Run Instructions
-----------------
    ./inspired_trek73          Run the program
    ./inspired_trek73 -version Print version and build information
    ./inspired_trek73 -h       Show all available command-line options

Supported startup options:
    -enemies N             Number of enemy ships (1-9, 0 = random)
    -player-class XX       Player ship class abbreviation (default CA)
    -enemy-class XX        Enemy ship class abbreviation (default CA)
    -race NAME             Enemy race name prefix (default random)
    -allow-silly-race      Include Monty Python race in random selection
    -ship-name NAME        Player ship name (default random Federation ship)

Gameplay Instructions / Command Reference
------------------------------------------
The game engine is not yet playable. This section will be completed as
gameplay is implemented in upcoming phases, and will document the full
original command set (all of which will remain supported), covering
phasers, torpedoes, probes, shields, helm, scanning, power distribution,
surrender/bluff tactics, and self-destruct.

Intentional Differences From the Original
--------------------------------------------
The following differences from the original C implementation are
intentional and are expected to remain as the project progresses:

- The save-game format and custom ship-class file format will be
  replaced with portable, explicit formats. The logical content of the
  original data (ship stats, race data, etc.) will be preserved.
- The optional natural-language command parser (originally implemented
  with lex/yacc) will be replaced with a simpler hand-written text
  parser. All original numeric command codes remain supported.
- Confirmed defects in the original implementation (portability
  accidents, undefined behavior, and clear logic bugs identified during
  source analysis) will be fixed rather than reproduced.
- Real-time turn countdown driven by Unix signals/alarm() will be
  replaced with a portable timer mechanism.

Explanation of Modernization Changes
--------------------------------------
- Rewritten from C to idiomatic Go using only the standard library.
- Clear separation between CLI presentation, command handling, and core
  game engine logic.
- Deterministic, seed-based golden tests are used to verify that ported
  mechanics match the original game's behavior (except where a fix is
  intentional, per above).
- Cross-platform build support for Windows, Linux, and macOS with no
  platform-specific code required.

Instructions for Future Developers
-------------------------------------
- Read prompt_1.md and prompt_1_answers_to_q first; they capture the
  project's guiding constraints and decisions.
- Read the remaining prompt_*.md files for additional context and design decisions.
- The original source under References/FreeBSD/trek73 is the
  authoritative gameplay specification. Consult it before changing any
  game mechanic.
- New gameplay mechanics should be accompanied by golden tests before
  being considered complete.
- Keep the engine (internal/game), command handling (internal/commands),
  static data (internal/data), and CLI (cmd/inspired_trek73) cleanly
  separated.
- Do not add third-party dependencies without first documenting why one
  is needed and getting approval.

