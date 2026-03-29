# go-jmp

go-jmp is a Go implementation of `jmp`: a terminal workflow tool for jumping from search results and command output directly into your editor.

It is designed to preserve the command contract and day-to-day behavior of the Raku version while producing a widely deployable Go binary.

## Quick Start

Build locally:

```bash
go mod tidy
go build -o jmp ./cmd/jmp
./jmp version
```

Install into your Go bin:

```bash
go install ./cmd/jmp
```

If needed:

```bash
export PATH="$HOME/go/bin:$PATH"
```

## First Run

Run:

```bash
jmp config
```

This opens `~/.jmp`, where you configure:

- `editor.command.template`
- `find.command.template`
- `browser.command.template`

Default search uses ripgrep (`rg`) in stable `file:line:text` mode.

## Commands

```text
jmp                                         show most recent hits
jmp back [count]                            show recent history (default 100)
jmp to '[<search-terms> ...]'               search files and jump to matching lines
jmp on '<command ...>'                      parse files from command output (stdout + stderr)
jmp edit <filename> [<line-number>]         start editing at a line number
jmp edit <filename> '[<search-terms> ...]'  start editing at a matching line
jmp config                                  edit ~/.jmp config
jmp help                                    show command help
jmp version                                 show version
```

Compatibility behavior:

- `jmp on` with no command defaults to history view.
- Unknown top-level command falls through to command-output mode.

## TUI Controls

- Up/Down: move selection
- Right/Enter on results: open preview
- Right/Enter on preview: open editor at selected line
- Left/Esc in preview: return to results
- `t`: new `jmp to` query in title bar
- `o`: new `jmp on` command in title bar
- `h` or `?`: help in preview
- `q` or `x`: quit

## Visual Layout

- Double-line outer frame
- Single-line separators between title, results, preview, and footer
- Fixed 15-row results pane
- Static title bar showing current command
- Green inverse highlight for selected row
- Centered footer actions with right-aligned version text

## Common Workflows

Search code and jump:

```bash
jmp to parser token
```

Jump from command output:

```bash
jmp on git status
jmp on tail -n 200 /var/log/syslog
```

Open directly:

```bash
jmp edit pkg/tui/ui.go 120
```

## Build, Test, Release

Run tests:

```bash
go test ./...
```

Build local binary:

```bash
go build -o jmp ./cmd/jmp
```

Cross-platform artifacts:

```bash
make cross
```

Release automation:

- `.github/workflows/release.yml`
- `.goreleaser.yml`

## Documentation

- Migration plan: `MIGRATION_PLAN.md`
- Release history: `CHANGES.md`
- Agent guidance: `AGENTS.md` and `agents/`
