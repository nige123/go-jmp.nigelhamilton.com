# jmp (Go)

`jmp` helps you jump from search results, command output, and recent history straight into your editor.

This repository is the Go implementation of the Raku tool at `/home/s3/jmp.nigelhamilton.com`, with command contract parity and matching workflows.

## Build

```bash
go build -o jmp ./cmd/jmp
```

## Install locally

```bash
go install ./cmd/jmp
```

## Commands

```text
jmp                                         show most recent hits
jmp to '[<search-terms> ...]'               search files and jump to matching lines
jmp on '<command ...>'                      parse files from command output (stdout + stderr)
jmp config                                  edit ~/.jmp config to set editor and search commands
jmp help                                    show command help
jmp version                                 show version
jmp edit <filename> [<line-number>]         start editing at a line number
jmp edit <filename> '[<search-terms> ...]'  start editing at a matching line
```

`jmp on` with no command defaults to `jmp back` behavior.
Unknown top-level commands fall through to command-output mode for backward compatibility.

## TUI keys

- Up/Down: move through hits
- Right/Enter on results: open preview
- Right/Enter on preview: open editor at selected line
- Left/Esc in preview: return to results
- `t`: prompt for `jmp to` search
- `o`: prompt for `jmp on` command
- `h` or `?`: help in preview
- `q`/`x`/Esc: quit

## Tests

```bash
go test ./...
```

## Deployable binaries

- `Makefile` includes local and cross-platform targets
- CI workflow builds Linux/macOS/Windows artifacts via GoReleaser

## Migration

See `MIGRATION_PLAN.md` for the incremental parity and hardening plan.
