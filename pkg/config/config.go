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
# uncomment or add your favourite text editor
#--------------------------------------------------------------------

editor.command.template   = nano +[-line-number-] "[-filename-]"

# atom has Raku syntax highlighting and other plugins for Raku
# editor.command.template   = atom [-filename-]:[-line-number-] &

# editor.command.template   = code -g [-filename-]:[-line-number-] &

# editor.command.template   = subl [-filename-]:[-line-number-] &
# editor.command.template   = emacs +[-line-number-]
# editor.command.template   = vim +[-line-number-] [-filename-]

#--------------------------------------------------------------------
# uncomment or add your preferred code searching tool (below)
#--------------------------------------------------------------------

# ripgrep - fast recursive search with stable file:line:text output
find.command.template       = rg --line-number --with-filename --no-heading --color never '[-search-terms-]'

# ag - the silver searcher for generic fast file searching
# find.command.template     = ag --nogroup '[-search-terms-]'

# git grep - for fast search of git repositories
# find.command.template     = git grep --full-name --untracked --text --line-number -e '[-search-terms-]'

# App::Ack - Perl-powered improvement to grep
# find.command.template     = ack --nogroup '[-search-terms-]'

#--------------------------------------------------------------------
# uncomment or add your preferred browser launch command
#--------------------------------------------------------------------

# open a browser at a URL
browser.command.template       = elinks '[-url-]'

# open the default browser on Mac OS
browser.command.template       = open '[-url-]'
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
