You are tasked with recreating a classic computer game called Trek73 as a modern Go application.

==================================================
SOURCE MATERIAL
==================================================

The original source code is located at:

D:\GitHub\appliedinspiration\Inspired_Trek73\References\FreeBSD\trek73

This directory contains the source code for a version of the classic game Trek73. The original program was written in the C programming language and designed to run on FreeBSD/Unix-like systems. It is a completely text-based command-line interface (CLI) game.

Your first task is NOT to write code.

First, carefully examine and understand the existing source code.

Analyze:

- Overall game architecture.
- Game initialization and main loop.
- User interaction model.
- Game rules and mechanics.
- Data structures.
- Combat logic.
- Ship behavior.
- Command processing.
- Random number generation.
- File handling.
- Any assumptions made about the original operating environment.
- Any external dependencies.
- Any platform-specific behavior.
- Any bugs or limitations caused by the original implementation environment.

After completing the analysis, provide:

1. An overview of the original program architecture.
2. A description of the game mechanics.
3. A migration plan explaining how the C implementation will be translated into Go.
4. A list of any questions or ambiguities requiring clarification.

Do not begin implementation until the analysis and migration plan are complete.

==================================================
PROJECT GOAL
==================================================

Create a modern implementation of Trek73 called:

Inspired Trek73

The goal is to create a faithful modernization of the original game.

This means:

- Preserve the original gameplay experience.
- Preserve the original game mechanics and rules.
- Preserve the overall feel and strategy of the original game.
- Modernize the code architecture, maintainability, and compatibility.
- Correct obvious bugs, portability issues, and limitations that existed only because of the original environment.

This is NOT intended to be a simple line-by-line translation of the C source code.

The original game behavior should be treated as the authoritative specification, but the implementation should use modern software engineering practices.

Examples of acceptable modernization:

- Fixing crashes or undefined behavior from the original C implementation.
- Replacing obsolete system-specific code.
- Improving input handling where the original depended on obsolete terminal behavior.
- Improving code organization and maintainability.
- Adding tests to ensure important game behavior remains consistent.

Examples of changes that should NOT be made without approval:

- Changing game balance.
- Adding new gameplay mechanics.
- Changing the command system significantly.
- Adding graphics.
- Adding sound.
- Adding networking or multiplayer features.
- Altering the strategic decisions available to the player.

If you identify a possible gameplay improvement, document it separately and request approval before implementing it.

==================================================
IMPLEMENTATION REQUIREMENTS
==================================================

The new implementation will be written entirely in the Go programming language.

Requirements:

- Use modern idiomatic Go practices.
- Target current Go releases.
- Produce a clean, maintainable codebase.
- Separate concerns clearly.
- Avoid unnecessary complexity.
- Maintain the simplicity appropriate for a classic CLI game.
- Do not create a web application, GUI application, or game engine abstraction. This is intentionally a CLI application

The program must compile and run on:

- Windows
- Linux
- macOS

The code should require no platform-specific modifications.

The project should support normal Go cross-compilation.

==================================================
DEPENDENCY REQUIREMENTS
==================================================

The preferred implementation uses only the Go standard library.

Do NOT add third-party dependencies.

If you believe an external package is necessary:

1. Stop implementation.
2. Explain why the dependency is needed.
3. Explain alternatives considered.
4. Wait for approval before proceeding.

==================================================
PROJECT STRUCTURE
==================================================

Create a maintainable Go project structure with clear separation between:

- Game engine logic.
- User interface / CLI handling.
- Game state management.
- Input parsing.
- Random number generation.
- Configuration/constants.
- Utility functions.

Avoid creating unnecessary abstractions.

The design should remain simple and appropriate for the size and nature of the project.

==================================================
RANDOMNESS AND GAME BALANCE
==================================================

The original game's random number generation affects gameplay balance.

During migration:

- Identify how randomness is used in the original C code.
- Preserve the statistical behavior where practical.
- Avoid replacing random calculations with significantly different algorithms unless necessary.
- Provide a way to use deterministic random seeds during testing.

The goal is to ensure the modernized version maintains the same gameplay balance as the original.

==================================================
GOLDEN TEST REQUIREMENT
==================================================

Before beginning the Go implementation, create a "golden test" strategy.

The purpose of the golden tests is to preserve important behaviors of the original Trek73 game and prevent accidental changes during modernization.

The golden test system should include:

- A documented set of representative game scenarios.
- Expected inputs and outputs.
- Expected game state transitions.
- Important calculations and mechanics that can be verified automatically.

Where practical, create deterministic test cases by controlling random number generation.

The golden tests should cover important areas such as:

- Game initialization.
- Ship creation and configuration.
- Player command processing.
- Movement and navigation.
- Combat calculations.
- Damage handling.
- Victory and defeat conditions.
- Any other mechanics identified during source analysis.

If original behavior cannot be easily tested automatically, create documented manual regression tests.

The golden test suite should become the reference point for determining whether the Go version remains faithful to the original game.

==================================================
DOCUMENTATION REQUIREMENTS
==================================================

The project must include:

README.txt

The README should contain:

- Project description.
- Brief history of the original Trek73 game.
- Original author information, if confirmed from source.
- Build instructions.
- Run instructions.
- Gameplay instructions.
- Command reference.
- Description of intentional differences from the original.
- Explanation of modernization changes.
- Instructions for future developers maintaining the project.

The Go source code should contain concise comments explaining:

- Major design decisions.
- Complex algorithms.
- Areas where the Go implementation differs from the original C implementation.

Avoid excessive comments that merely restate obvious code.

==================================================
BUILD REQUIREMENTS
==================================================

The final executable should be named:

inspired_trek73

with the appropriate operating system extension:

Windows:

inspired_trek73.exe

Linux/macOS:

inspired_trek73

The program should support standard Go build commands, including cross-compilation.

==================================================
PROGRAM INFORMATION DISPLAY
==================================================

When the program starts, before gameplay begins, display:

Inspired Trek73

A modern remake of the classic Trek73 game

Updated to modern operating systems and hardware by

Peter S. Lee
AppliedInspiration.com

The original Trek73 game was created by David A. Smith.

Before adding any licensing statements:

- Inspect the original source code.
- Verify the appropriate attribution and licensing language.
- Do not assume licensing status without confirmation from the source.

If the original source confirms public-domain status, include an appropriate statement acknowledging the original work.

The modern Go implementation should include an appropriate copyright notice:

Copyright 2026 Peter S. Lee (AppliedInspiration.com)

The program should also display the disclaimer:

"This program is provided 'as is' without warranties or guarantees of any kind."

This information should be displayed:

- When the program starts.
- When the program exits.

==================================================
VERSION INFORMATION
==================================================

Provide a command-line option such as:

--version

or

-version

that displays:

- Program name.
- Version number.
- Build date.
- Go compiler version.
- Any relevant build information.

Use Go build metadata where appropriate.

==================================================
DEVELOPMENT PROCESS
==================================================

Work incrementally.

Phase 1:
- Analyze the original C source.
- Identify the original architecture.
- Document game mechanics.
- Identify portability issues.
- Identify bugs or limitations caused by the original environment.
- Produce an architecture summary.
- Produce a modernization plan.

Phase 2:
- Design the golden test strategy.
- Create regression scenarios.
- Document expected behavior from the original implementation.

Phase 3:
- Create the Go project structure.
- Implement core game state and logic.
- Ensure golden tests pass as functionality is migrated.

Phase 4:
- Implement CLI interaction.
- Replace obsolete terminal-specific behavior with portable Go implementations.

Phase 5:
- Compare the Go version against the original behavior.
- Resolve unexpected differences.
- Document intentional modernization changes.

Phase 6:
- Complete README documentation.
- Perform final cleanup.
- Verify cross-platform builds.

At each major phase:

- Summarize what was completed.
- Explain important design decisions.
- Identify differences from the original game.
- Identify questions requiring my input.

Do not make assumptions about unclear behavior.

Ask questions whenever the original source is ambiguous.

==================================================
FINAL OBJECTIVE
==================================================

The final result should be a clean, portable, maintainable Go implementation of Trek73 that feels like the original game to players while taking advantage of modern programming practices.

The finished program should be something that a developer can clone, build, run, understand, and maintain many years into the future.