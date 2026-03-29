# Migration Plan: Raku jmp -> Go jmp

## Goal

Port all user-visible behavior from the Raku CLI and TUI implementation to Go without functional loss.

## Phase 1: Contract Lock

1. Keep subcommands and defaults stable: `to`, `on`, `edit`, `config`, `help`, `version`, `back`.
2. Preserve fallback behavior: unknown command -> `on` mode; empty `on` -> history.
3. Preserve output contract: `jmp - version <VERSION>` and help text semantics.

## Phase 2: Core Parity (implemented)

1. Config parser for `~/.jmp` with default templates.
2. Template token renderer for `[-key-]` substitutions.
3. Finder support for `file:line:text` search output with malformed-line filtering.
4. Deferred command-output parsing with `HitLater` and language-specific file extraction patterns.
5. Editor handoff and memory persistence to `~/.jmp.hist` with max-entries trim.

## Phase 3: TUI Parity (implemented baseline)

1. Result + preview workflow with directional navigation.
2. In-UI prompts for text search (`t`) and command output (`o`).
3. Preview window clamping and centering logic preserved.
4. Footer shortcuts and version rendering.

## Phase 4: Regression Hardening (next)

1. Add more parser tests for Perl/Raku/Git output examples.
2. Add end-to-end tests for `edit` workflows with fixture files.
3. Add robustness checks for non-UTF8 and control-code heavy outputs.
4. Add snapshot tests for help text and default config generation.

## Phase 5: Distribution

1. Build static-ish binaries with `CGO_ENABLED=0` for portability.
2. Publish multi-arch release artifacts via GoReleaser.
3. Add Homebrew tap and package metadata.
4. Add checksums and signed release artifacts.

## No-Functionality-Loss Checklist

1. Subcommands and defaults identical.
2. Command-output mode supports both stdout/stderr inputs.
3. Preview line math equivalent.
4. Config and memory file formats remain compatible.
5. Unknown command fallback preserved.
6. Interactive key flows preserved.
