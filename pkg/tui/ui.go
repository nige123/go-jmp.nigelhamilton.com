package tui

import (
    "fmt"
    "os"
    "strings"

    tea "github.com/charmbracelet/bubbletea"

    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/editor"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/file"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/memory"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/model"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/version"
)

type SearcherFunc func(terms string) ([]model.Renderable, error)
type OutputerFunc func(command string) ([]model.Renderable, error)

type UI struct {
    Title    string
    Editor   *editor.Editor
    Hits     []model.Renderable
    Searcher SearcherFunc
    Outputer OutputerFunc

    selectedIndex int
    focus         string
    mode          string
    inputContext  string
    searchBuffer  string
    outputBuffer  string

    previewLines []string
    previewStart int
    previewLine  int
    previewHit   *file.Hit

    width  int
    height int
}

func NewUI(title string, editor *editor.Editor, hits []model.Renderable, searcher SearcherFunc, outputer OutputerFunc) *UI {
    return &UI{
        Title:         title,
        Editor:        editor,
        Hits:          hits,
        Searcher:      searcher,
        Outputer:      outputer,
        selectedIndex: 0,
        focus:         "results",
        mode:          "command",
        previewLines: []string{
            "Press Right Arrow on a result to preview the file here.",
            "Press Right Arrow again in this pane to open the editor.",
            "Press [o] in results to run a command and jump on its output.",
        },
        previewStart: 1,
        previewLine:  1,
    }
}

func (u *UI) ClampPreviewLine(requestedLine, maxLine int) int {
    if maxLine < 1 {
        return 1
    }
    if requestedLine < 1 {
        return 1
    }
    if requestedLine > maxLine {
        return maxLine
    }
    return requestedLine
}

func (u *UI) PreviewWindow(maxLine, targetLine, windowSize int) (int, int) {
    if maxLine < 1 {
        return 1, 1
    }

    safeTarget := u.ClampPreviewLine(targetLine, maxLine)
    if windowSize < 1 {
        windowSize = 1
    }
    if maxLine <= windowSize {
        return 1, maxLine
    }

    half := windowSize / 2
    start := safeTarget - half
    end := start + windowSize - 1

    if start < 1 {
        start = 1
        end = windowSize
    }
    if end > maxLine {
        end = maxLine
        start = end - windowSize + 1
    }

    return start, end
}

func (u *UI) LoadSearchResults(terms string) bool {
    if u.Searcher == nil {
        return false
    }
    trimmed := strings.TrimSpace(terms)
    if trimmed == "" {
        return false
    }

    hits, err := u.Searcher(trimmed)
    if err != nil {
        return false
    }

    u.Title = "jmp to " + trimmed
    u.Hits = hits
    u.selectedIndex = 0
    u.focus = "results"
    return true
}

func (u *UI) LoadCommandOutput(command string) bool {
    if u.Outputer == nil {
        return false
    }
    trimmed := strings.TrimSpace(command)
    if trimmed == "" {
        return false
    }

    hits, err := u.Outputer(trimmed)
    if err != nil {
        return false
    }

    u.Title = "jmp on " + trimmed
    u.Hits = hits
    u.selectedIndex = 0
    u.focus = "results"
    return true
}

func (u *UI) Init() tea.Cmd {
    return nil
}

func (u *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        u.width = msg.Width
        u.height = msg.Height
    case tea.KeyMsg:
        if u.mode == "input" {
            return u.handleInputKey(msg)
        }
        return u.handleCommandKey(msg)
    }
    return u, nil
}

func (u *UI) handleInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.Type {
    case tea.KeyEnter:
        if u.inputContext == "search" {
            u.LoadSearchResults(u.searchBuffer)
        } else if u.inputContext == "output" {
            u.LoadCommandOutput(u.outputBuffer)
        }
        u.searchBuffer = ""
        u.outputBuffer = ""
        u.inputContext = ""
        u.mode = "command"
        return u, nil
    case tea.KeyEsc:
        u.searchBuffer = ""
        u.outputBuffer = ""
        u.inputContext = ""
        u.mode = "command"
        return u, nil
    case tea.KeyBackspace, tea.KeyDelete:
        if u.inputContext == "output" {
            if len(u.outputBuffer) > 0 {
                u.outputBuffer = u.outputBuffer[:len(u.outputBuffer)-1]
            }
        } else {
            if len(u.searchBuffer) > 0 {
                u.searchBuffer = u.searchBuffer[:len(u.searchBuffer)-1]
            }
        }
        return u, nil
    default:
        if len(msg.String()) == 1 {
            if u.inputContext == "output" {
                u.outputBuffer += msg.String()
            } else {
                u.searchBuffer += msg.String()
            }
        }
    }
    return u, nil
}

func (u *UI) handleCommandKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "ctrl+c", "q", "x":
        return u, tea.Quit
    case "esc":
        if u.focus == "preview" {
            u.focus = "results"
            return u, nil
        }
        return u, tea.Quit
    case "h", "?":
        u.focus = "preview"
        u.previewHit = nil
        u.previewStart = 1
        u.previewLine = 1
        u.previewLines = []string{
            "jmp help",
            "",
            "Right / Enter  Open selected result or preview line",
            "Left           Return focus from preview to results",
            "Up / Down      Move selection",
            "PageUp / l     Page through content",
            "t              Jump to text in files (jmp to ...)",
            "o              Jump on files in command output (jmp on ...)",
            "q / x          Quit jmp",
            "h / ?          Show this help",
        }
        return u, nil
    case "t":
        if u.Searcher != nil {
            u.mode = "input"
            u.inputContext = "search"
            u.searchBuffer = ""
        }
        return u, nil
    case "o":
        if u.Outputer != nil {
            u.mode = "input"
            u.inputContext = "output"
            u.outputBuffer = ""
        }
        return u, nil
    case "up":
        if u.focus == "preview" {
            if u.previewLine > u.previewStart {
                u.previewLine--
            }
            return u, nil
        }
        if u.selectedIndex > 0 {
            u.selectedIndex--
        }
        return u, nil
    case "down":
        if u.focus == "preview" {
            max := u.previewStart + len(u.previewLines) - 1
            if u.previewLine < max {
                u.previewLine++
            }
            return u, nil
        }
        if u.selectedIndex < len(u.Hits)-1 {
            u.selectedIndex++
        }
        return u, nil
    case "right", "enter":
        if u.focus == "preview" {
            if u.previewHit != nil {
                hit := *u.previewHit
                hit.LineNumber = u.previewLine
                _ = u.Editor.Edit(u.Title, &hit)
            }
            return u, nil
        }
        u.previewSelectedHit()
        return u, nil
    case "left":
        if u.focus == "preview" {
            u.focus = "results"
        }
        return u, nil
    }
    return u, nil
}

func (u *UI) previewSelectedHit() {
    if u.selectedIndex < 0 || u.selectedIndex >= len(u.Hits) {
        return
    }

    resolved := u.resolveEditableHit(u.Hits[u.selectedIndex])
    u.focus = "preview"
    u.previewStart = 1
    u.previewLine = 1

    if resolved == nil {
        u.previewHit = nil
        u.previewLines = []string{"No preview available for this entry."}
        return
    }

    u.previewHit = resolved
    payload, err := os.ReadFile(resolved.AbsolutePath)
    if err != nil {
        u.previewLines = []string{"No preview available for this entry."}
        return
    }

    lines := strings.Split(strings.ReplaceAll(string(payload), "\r\n", "\n"), "\n")
    if len(lines) == 0 {
        u.previewLines = []string{"     1:"}
        u.previewLine = 1
        return
    }

    start, end := u.PreviewWindow(len(lines), resolved.LineNumber, 400)
    rendered := make([]string, 0, end-start+1)
    for i := start; i <= end; i++ {
        rendered = append(rendered, fmt.Sprintf("%6d: %s", i, lines[i-1]))
    }

    u.previewLines = rendered
    u.previewStart = start
    u.previewLine = u.ClampPreviewLine(resolved.LineNumber, len(lines))
}

func (u *UI) resolveEditableHit(hit model.Renderable) *file.Hit {
    switch candidate := hit.(type) {
    case *file.Hit:
        if candidate.FileExists() {
            return candidate
        }
    case *file.HitLater:
        candidate.FindFilePath()
        if candidate.FileExists() {
            return &candidate.Hit
        }
    case *memory.Hit:
        if candidate.FileExists() {
            return &candidate.Hit
        }
    }
    return nil
}

func (u *UI) View() string {
    width := u.width
    if width < 40 {
        width = 80
    }

    title := u.Title
    if u.mode == "input" {
        if u.inputContext == "output" {
            title = "jmp on " + u.outputBuffer + "▌"
        } else {
            title = "jmp to " + u.searchBuffer + "▌"
        }
    }

    footer := "[↑][↓] [←][→] select"
    if u.Searcher != nil {
        footer += "  [t]o search"
    }
    if u.Outputer != nil {
        footer += "  [o]n cmd"
    }
    footer += "  [h]elp  [q]uit"

    versionText := "jmp v" + version.VERSION
    if len(footer)+len(versionText)+1 < width {
        spaces := strings.Repeat(" ", width-len(footer)-len(versionText))
        footer = footer + spaces + versionText
    }

    resultLines := make([]string, 0, len(u.Hits))
    if len(u.Hits) == 0 {
        resultLines = append(resultLines, "(no output lines)")
    } else {
        for i, hit := range u.Hits {
            prefix := "  "
            if u.focus == "results" && i == u.selectedIndex {
                prefix = "> "
            }
            resultLines = append(resultLines, prefix+hit.Render())
        }
    }

    previewLines := make([]string, 0, len(u.previewLines))
    for idx, line := range u.previewLines {
        current := u.previewStart + idx
        prefix := "  "
        if u.focus == "preview" && current == u.previewLine {
            prefix = "> "
        }
        previewLines = append(previewLines, prefix+line)
    }

    return strings.Join([]string{
        title,
        strings.Join(resultLines, "\n"),
        strings.Join(previewLines, "\n"),
        footer,
    }, "\n")
}

func (u *UI) Display() error {
    model := u
    program := tea.NewProgram(model, tea.WithAltScreen())
    _, err := program.Run()
    return err
}
