package app

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/config"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/editor"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/finder"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/model"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/memory"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/tui"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/version"
)

type App struct {
    args    []string
    command string

    config *config.Config
    editor *editor.Editor
    finder *finder.Finder
    memory *memory.Memory
}

func New(args []string) (*App, error) {
    cfg, err := config.New()
    if err != nil {
        return nil, err
    }
    mem, err := memory.New()
    if err != nil {
        return nil, err
    }

    ed := editor.New(cfg, mem)
    fd := finder.New(cfg)

    return &App{
        args:    args,
        command: "jmp " + strings.Join(args, " "),
        config:  cfg,
        editor:  ed,
        finder:  fd,
        memory:  mem,
    }, nil
}

func (a *App) Run() error {
    if len(a.args) == 0 {
        return a.recentJmps(100)
    }

    switch a.args[0] {
    case "back":
        count := 100
        if len(a.args) > 1 {
            parsed, err := strconv.Atoi(a.args[1])
            if err != nil {
                return fmt.Errorf("invalid count for back: %w", err)
            }
            count = parsed
        }
        return a.recentJmps(count)
    case "config":
        return a.editor.EditFile(a.config.ConfigFile())
    case "edit":
        return a.runEdit()
    case "version":
        fmt.Printf("jmp - version %s\n", version.VERSION)
        return nil
    case "help":
        fmt.Print(usageText)
        return nil
    case "in":
        terms := strings.Join(a.args[1:], " ")
        return a.searchInFiles(strings.TrimSpace(terms))
    case "to":
        return a.runTo()
    case "on":
        command := strings.TrimSpace(strings.Join(a.args[1:], " "))
        if command == "" {
            return a.recentJmps(100)
        }
        return a.findFilesInCommandOutput(command)
    default:
        command := strings.TrimSpace(strings.Join(a.args, " "))
        if command == "" {
            return a.recentJmps(100)
        }
        return a.findFilesInCommandOutput(command)
    }
}

func (a *App) runEdit() error {
    if len(a.args) < 2 {
        return fmt.Errorf("edit requires at least a filename")
    }

    filename := a.args[1]
    if len(a.args) == 2 {
        return a.editFile(filename, 1)
    }

    if lineNumber, err := strconv.Atoi(a.args[2]); err == nil && len(a.args) == 3 {
        return a.editFile(filename, lineNumber)
    }

    searchTerms := strings.Join(a.args[2:], " ")
    return a.editFileAtMatchingLine(filename, searchTerms)
}

func (a *App) editFile(filename string, lineNumber int) error {
    hit := a.finder.FindLineInFile(filename, lineNumber)
    return a.editor.Edit(a.command, hit)
}

func (a *App) editFileAtMatchingLine(filename, searchTerms string) error {
    hit, err := a.finder.FindMatchingLineInFile(filename, searchTerms)
    if err != nil {
        return err
    }
    return a.editor.Edit(a.command, hit)
}

func (a *App) runTo() error {
    if len(a.args) < 2 {
        return fmt.Errorf("jmp to requires a filename")
    }

    filename := a.args[1]
    lineNumber := 0
    if len(a.args) >= 3 {
        parsed, err := strconv.Atoi(a.args[2])
        if err == nil {
            lineNumber = parsed
        }
    }

    if lineNumber > 0 {
        hit := a.finder.FindLineInFile(filename, lineNumber)
        return a.editor.Edit(a.command, hit)
    }

    return a.locateFiles(filename)
}

func (a *App) locateFiles(searchTerms string) error {
    hits, err := a.finder.FindFilesOnFilesystem(searchTerms)
    if err != nil {
        return err
    }
    return a.displayHits("jmp to "+searchTerms, hits)
}

func (a *App) findFilesInCommandOutput(command string) error {
    hits, err := a.finder.FindFilesInCommandOutput(command)
    if err != nil {
        return err
    }
    return a.displayHits("jmp on "+command, hits)
}

func (a *App) searchInFiles(searchTerms string) error {
    hits, err := a.finder.FindInFiles(searchTerms)
    if err != nil {
        return err
    }
    return a.displayHits("jmp in "+searchTerms, hits)
}

func (a *App) recentJmps(lastNEntries int) error {
    hits := a.memory.GetRecentJmps(lastNEntries)
    return a.displayHits("jmp", hits)
}

func (a *App) displayHits(title string, hits []model.Renderable) error {
    if len(hits) == 0 {
        return nil
    }

    ui := tui.NewUI(
        title,
        a.editor,
        hits,
        func(terms string) ([]model.Renderable, error) {
            return a.finder.FindInFiles(terms)
        },
        func(command string) ([]model.Renderable, error) {
            return a.finder.FindFilesInCommandOutput(command)
        },
        func(filename string) ([]model.Renderable, error) {
            return a.finder.FindFilesOnFilesystem(filename)
        },
    )
    return ui.Display()
}

const usageText = `

jmp - jump to files in your workflow

Usage:

    jmp                                         -- show most recent jmps
    jmp in '[<search-terms> ...]'               -- search in files for matching lines
    jmp to <filename>                           -- locate a file anywhere in the filesystem
    jmp to <filename> <line-number>             -- jump to a specific line in a file
    jmp on '<command ...>'                      -- jump on files in command output

    # jmp on examples:
    jmp on tail /some.log                       -- files mentioned in log files
    jmp on ls                                   -- files in a directory
    jmp on find .                               -- files returned from find
    jmp on git status                           -- files in git
    jmp on perl test.pl                         -- Perl output and errors
    jmp on raku test.raku                       -- Raku output and errors

    jmp config                                  -- edit ~/.jmp config
    jmp help                                    -- show this help

    jmp edit <filename> [<line-number>]         -- start editing at a line number
    jmp edit <filename> '[<search-terms> ...]'  -- start editing at a matching line
`
