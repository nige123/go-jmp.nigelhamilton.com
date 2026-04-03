package config

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    jmptemplate "github.com/nige123/go-jmp.nigelhamilton.com/pkg/template"
)

type Config struct {
    configFile string
    fields     map[string]string
}

const defaultConfig = `
#--------------------------------------------------------------------
# Editor — uncomment or add your favourite text editor
#--------------------------------------------------------------------

# nano — simple terminal editor (Linux, macOS, Windows/WSL)
editor.command.template   = nano +[-line-number-] "[-filename-]"

# VS Code — popular cross-platform editor (Linux, macOS, Windows)
# editor.command.template   = code -g [-filename-]:[-line-number-] &

# Vim — ubiquitous terminal editor (Linux, macOS, Windows)
# editor.command.template   = vim +[-line-number-] [-filename-]

# Neovim — modern Vim fork (Linux, macOS, Windows)
# editor.command.template   = nvim +[-line-number-] [-filename-]

# Emacs — extensible editor (Linux, macOS, Windows)
# editor.command.template   = emacs +[-line-number-] [-filename-]

# Sublime Text — fast GUI editor (Linux, macOS, Windows)
# editor.command.template   = subl [-filename-]:[-line-number-] &

# Helix — post-modern terminal editor (Linux, macOS, Windows)
# editor.command.template   = hx [-filename-]:[-line-number-]

# Micro — intuitive terminal editor (Linux, macOS, Windows)
# editor.command.template   = micro +[-line-number-] [-filename-]

# Atom — hackable editor (Linux, macOS, Windows) [discontinued but still used]
# editor.command.template   = atom [-filename-]:[-line-number-] &

#--------------------------------------------------------------------
# Search — uncomment or add your preferred code searching tool
# Used by: jmp in <search-terms>
#--------------------------------------------------------------------

# ripgrep — fast recursive search (Linux, macOS, Windows)
find.command.template       = rg --line-number --with-filename --no-heading --color never '[-search-terms-]'

# ag — the silver searcher (Linux, macOS, Windows)
# find.command.template     = ag --nogroup '[-search-terms-]'

# git grep — search within git repositories (any platform with git)
# find.command.template     = git grep --full-name --untracked --text --line-number -e '[-search-terms-]'

# ack — Perl-powered grep alternative (Linux, macOS, Windows)
# find.command.template     = ack --nogroup '[-search-terms-]'

# grep — universal fallback (Linux, macOS)
# find.command.template     = grep -rn '[-search-terms-]' .

#--------------------------------------------------------------------
# Locate — uncomment or add your preferred file locating tool
# Used by: jmp to <filename>
#--------------------------------------------------------------------

# locate — fast indexed file search (Linux, macOS with findutils)
locate.command.template     = locate '[-search-terms-]'

# plocate — fast modern replacement for locate (newer Linux distros)
# locate.command.template   = plocate '[-search-terms-]'

# mlocate — common on Ubuntu/Debian
# locate.command.template   = mlocate '[-search-terms-]'

# fd — simple, fast find alternative (Linux, macOS, Windows)
# locate.command.template   = fd --type f '[-search-terms-]'

# find — universal fallback, no index needed (Linux, macOS)
# locate.command.template   = find / -name '*[-search-terms-]*' -type f 2>/dev/null

# where — Windows built-in file search
# locate.command.template   = where /r \\ [-search-terms-]

# mdfind — macOS Spotlight search from terminal
# locate.command.template   = mdfind -name '[-search-terms-]'

# everything — voidtools Everything CLI for Windows (very fast)
# locate.command.template   = es -name '[-search-terms-]'

#--------------------------------------------------------------------
# Browser — uncomment or add your preferred browser launch command
#--------------------------------------------------------------------

# xdg-open — open URLs in default browser (Linux)
browser.command.template       = xdg-open '[-url-]'

# open — open URLs in default browser (macOS)
# browser.command.template   = open '[-url-]'

# start — open URLs in default browser (Windows)
# browser.command.template   = start '[-url-]'

# elinks — terminal-based browser (Linux, macOS)
# browser.command.template   = elinks '[-url-]'

# wslview — open URLs from WSL in Windows browser
# browser.command.template   = wslview '[-url-]'
`

func New() (*Config, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("could not determine HOME directory: %w", err)
    }

    path := filepath.Join(home, ".jmp")
    if _, err := os.Stat(path); os.IsNotExist(err) {
        if err := os.WriteFile(path, []byte(defaultConfig), 0o644); err != nil {
            return nil, fmt.Errorf("could not create default config %s: %w", path, err)
        }
    }

    fields, err := parseConfigFile(path)
    if err != nil {
        return nil, err
    }

    return &Config{configFile: path, fields: fields}, nil
}

func (c *Config) ConfigFile() string {
    return c.configFile
}

func (c *Config) Get(key string) (string, error) {
    value, ok := c.fields[key]
    if !ok {
        return "", fmt.Errorf("key %s does not exist in config file. Please add a value for %s to %s", key, key, c.configFile)
    }
    return value, nil
}

func (c *Config) GetTemplate(key string, params map[string]string) (string, error) {
    value, err := c.Get(key)
    if err != nil {
        return "", err
    }
    return jmptemplate.Render(value, params), nil
}

func parseConfigFile(path string) (map[string]string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("could not open config file %s: %w", path, err)
    }
    defer f.Close()

    fields := map[string]string{}
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        if key != "" {
            fields[key] = value
        }
    }
    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("could not read config file %s: %w", path, err)
    }
    return fields, nil
}
