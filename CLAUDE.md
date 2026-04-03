# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is jmp

`jmp` is a Go CLI tool that helps developers quickly jump to files and lines. It searches code, parses command output for file:line references, and presents results in an interactive TUI with live preview. Go port of the original Raku implementation.

## Build and Test Commands

```bash
make build          # CGO_ENABLED=0 go build -o jmp ./cmd/jmp (static binary)
make test           # go test ./...
make cross          # Cross-compile for Linux/macOS/Windows (amd64/arm64/386)
make tidy           # go mod tidy
make clean          # Remove build artifacts
go test ./pkg/finder/  # Run tests for a single package
```

## Architecture

**Entry point:** `cmd/jmp/main.go` → `internal/app/app.go` (App struct coordinates everything)

**Command dispatch:** `App.Run()` routes on first arg: `in` (search in files), `to` (locate files, optional line number), `on` (parse command output), `edit` (open file), `back` (history), `config`, `help`, `version`. No args shows recent history. Unknown commands treated as `on`.

**Key packages:**
- `pkg/config/` — Parses `~/.jmp` config, template token rendering (`[-token-]` syntax)
- `pkg/finder/` — Runs search tools (ripgrep default), parses command output for file:line refs
- `pkg/editor/` — Executes editor with line number support
- `pkg/file/` — `Hit` (file:line pair), `HitLater` (deferred parsing of command output using Perl/Raku patterns)
- `pkg/memory/` — History persistence (`~/.jmp.hist`, JSON array, max 100 entries)
- `pkg/tui/` — Bubbletea-based TUI with results/preview panes, command/input modes
- `pkg/model/` — `Renderable` interface for polymorphic display across Hit types

**Design patterns:**
- Side effects (filesystem, shell, terminal) stay at package boundaries; core logic is pure
- Dependency injection via constructors (`NewX`); interfaces at boundaries for testability
- `HitLater` defers parsing command output until needed (lazy evaluation)
- Three renderable types: `file.Hit`, `file.HitLater`, `memory.Hit`

## Conventions

- **4-space indentation** (not tabs)
- Only dependency is `github.com/charmbracelet/bubbletea` — do not add dependencies without strong justification
- Public command behavior (names, args, help text, exit codes, output format) is a contract — preserve unless explicitly changing
- Prefer explicit parsing over regex stacking
- Small, test-backed, easy-to-revert patches — don't mix cleanup with behavior changes
- Run `go test ./...` after behavior changes; run `go build ./cmd/jmp` after CLI/package changes

## Agent Guidance

The `agents/` directory contains role-based guidance files. Read `agents/START_HERE.md` first, then `agents/PROJECT_INVARIANTS.md` for the 8 core invariants (model-first, explicit parsing, side effects at edges, CLI contract preservation, no silent guessing, stable output, tests as contracts, teachable code).
