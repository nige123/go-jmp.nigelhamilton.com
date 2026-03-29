# CHANGES

## v54 - 2026-03-29

Release focus: Go-port usability, TUI parity, and deployment readiness.

### Added

- New-user README with quick-start, setup, workflows, and troubleshooting notes.
- Release changelog file for ongoing versioned change tracking.
- Cross-platform release flow via GoReleaser + GitHub Actions.

### Changed

- Footer branding shows `go-jmp v<version>` for runtime clarity.
- TUI frame now uses double-line outer border and single-line panel separators.
- Title bar remains static with current command context.
- Input mode title behavior now replaces current command text with prompt text (`jmp to` / `jmp on`) and a green cursor.
- Selected line highlighting uses green inverse style.
- Footer actions are centered with version right-aligned.
- Post-editor return path now forces redraw after exiting external editors (for example `nano`).

### Fixed

- Header/title rendering no longer shifts off-screen from frame-height miscalculation.
- TUI redraw after external editor return no longer leaves a blank screen.

## v53 - 2026-03-29

Initial Go implementation of `jmp` with functional compatibility targets:

- CLI subcommand parity (`to`, `on`, `edit`, `config`, `help`, `version`, `back`).
- Config parsing and default `~/.jmp` generation.
- Finder parsing for `file:line:text` search output.
- Command-output mode with lazy file-path extraction.
- Memory persistence in `~/.jmp.hist` with bounded history.
- Bubble Tea terminal interface with results/preview workflow.
- Test baseline for config, finder, preview math, and version output.
- Initial repository scaffolding, docs, and release tooling.
