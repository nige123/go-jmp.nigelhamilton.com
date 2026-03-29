# jmp

**Stop copying filenames from terminal output. Start jumping.**

`jmp` searches your codebase and parses command output to find files and line numbers, then opens them directly in your editor. No more copy-pasting paths.

```bash
jmp to 'parse error'        # search code, jump to matches
jmp on 'git status'          # parse command output, jump to files
jmp edit main.go 42          # open file at line 42
```

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/nige123/go-jmp.nigelhamilton.com/main/scripts/install.sh | bash
```

This detects your OS/architecture, downloads the right binary, and installs to `~/.local/bin/jmp`.

Add to your PATH if needed:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Or download manually from [Releases](https://github.com/nige123/go-jmp.nigelhamilton.com/releases).

## First Run

```bash
jmp config
```

This opens `~/.jmp` where you set your editor, search tool, and browser commands.

## Commands

```text
jmp                          show most recent hits
jmp to '<search-terms>'      search files and jump to matches
jmp on '<command>'           parse command output for files
jmp edit <file> [<line>]     open file at a line number
jmp back [count]             show recent history
jmp config                   edit config
jmp help                     show help
jmp version                  show version
```

## Feedback

This is a Go port of the original [Raku jmp](https://github.com/nige123/jmp.nigelhamilton.com). For feature requests, bug reports, and feedback, please use the issues on the [Raku repo](https://github.com/nige123/jmp.nigelhamilton.com/issues).
