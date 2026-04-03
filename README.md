# jmp

**Stop copying filenames from terminal output. Start jumping.**

`jmp` is a fast CLI tool that finds files and opens them in your editor at the right line. Three verbs, one workflow:

```bash
jmp in 'TODO'               # search in files for matching lines
jmp to main.go              # locate a file anywhere in the filesystem
jmp on 'git status'         # jump on files in command output
```

![jmp demo](images/demo.gif)

## How it works

1. You tell `jmp` what to find — a keyword, a filename, or a command
2. `jmp` shows results in an interactive TUI with a live preview pane
3. You pick a result and `jmp` opens your editor at that exact line

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/nige123/go-jmp.nigelhamilton.com/main/scripts/install.sh | bash
```

Auto-detects your OS and architecture. Installs to `~/.local/bin/jmp`.

Add to your PATH if needed:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**Supported platforms:** Linux (amd64, arm64, 386), macOS (amd64, arm64), Windows (amd64)

Or download manually from [Releases](https://github.com/nige123/go-jmp.nigelhamilton.com/releases).

## Quick start

```bash
jmp config                   # set your editor and search tools
jmp in 'parse error'         # search in files, pick a match, start editing
```

## The three verbs

### `jmp in` — search in files

Search inside files for keywords. Uses [ripgrep](https://github.com/BurntSushi/ripgrep) by default.

```bash
jmp in 'func main'          # find function definitions
jmp in 'TODO'               # find all TODOs in your codebase
jmp in 'import os'          # find imports
```

### `jmp to` — locate files

Find files anywhere in the filesystem. Uses `locate` by default.

```bash
jmp to README.md            # find all README.md files on the system
jmp to main.go              # locate a file by name
jmp to main.go 42           # open main.go at line 42
```

### `jmp on` — jump on command output

Run any command and parse its output for file references. Works with logs, build errors, test output, git — anything that mentions filenames.

```bash
jmp on 'git status'         # files changed in git
jmp on 'git diff --name-only'  # files with diffs
jmp on 'make build'         # jump to build errors
jmp on 'go test ./...'      # jump to test failures
jmp on 'tail -100 app.log'  # files mentioned in logs
jmp on 'find . -name "*.go"'   # files from find
jmp on 'ls'                 # files in current directory
```

## TUI controls

Once results are displayed:

| Key | Action |
|-----|--------|
| Up / Down | Move through results |
| Right / Enter | Preview file or open in editor |
| Left | Return to results from preview |
| `i` | Search in files (jmp in) |
| `o` | Run a command (jmp on) |
| `h` or `?` | Show help |
| `q` or `x` | Quit |

## Configuration

Run `jmp config` to edit `~/.jmp`. The config file lets you set:

- **Editor** — nano, vim, nvim, VS Code, Sublime, Emacs, Helix, Micro
- **Search tool** — ripgrep, ag, git grep, ack, grep
- **Locate tool** — locate, plocate, mlocate, fd, find, mdfind (macOS)
- **Browser** — xdg-open (Linux), open (macOS), start (Windows)

Each option includes commented-out alternatives for different platforms. Uncomment the one that fits your setup.

## All commands

```
jmp                                         show recent jmps
jmp in '<search-terms>'                     search in files for matching lines
jmp to <filename>                           locate a file anywhere
jmp to <filename> <line-number>             jump to a specific line in a file
jmp on '<command>'                          jump on files in command output

jmp config                                  edit ~/.jmp config
jmp edit <filename> [<line-number>]         open file at a line number
jmp edit <filename> '<search-terms>'        open file at a matching line
jmp back [count]                            show recent history
jmp help                                    show help
jmp version                                 show version
```

## Requirements

`jmp` is a single static binary with no runtime dependencies. For full functionality you'll want:

- A search tool for `jmp in` — [ripgrep](https://github.com/BurntSushi/ripgrep) recommended
- A locate tool for `jmp to` — `locate` / `plocate` / `fd` recommended

## Feedback

For feature requests, bug reports, and feedback: [GitHub Issues](https://github.com/nige123/go-jmp.nigelhamilton.com/issues)

This is a Go port of the original [Raku jmp](https://github.com/nige123/jmp.nigelhamilton.com).
