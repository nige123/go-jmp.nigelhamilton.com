package tui

import (
    "fmt"
    "os"
    "strings"
    "unicode/utf8"

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
	    "Press [i] to search in files for keywords.",
            "Press [o] to run a command and jump on its output.",
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

    u.Title = "jmp in " + trimmed
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
            "i              Search in files (jmp in ...)",
            "o              Jump on files in command output (jmp on ...)",
            "q / x          Quit jmp",
            "h / ?          Show this help",
        }
        return u, nil
    case "i":
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
            var redraw tea.Cmd
            if u.previewHit != nil && u.Editor != nil {
                hit := *u.previewHit
                hit.LineNumber = u.previewLine
                _ = u.Editor.Edit(u.Title, &hit)
                redraw = tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
            }
            return u, redraw
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
    if width < 50 {
        width = 80
    }
    height := u.height
    if height < 18 {
        height = 24
    }
    innerWidth := width - 2
    if innerWidth < 20 {
        innerWidth = 20
    }

    resultsHeight, previewHeight := u.paneHeights(height)

    title := u.Title
    if u.mode == "input" {
        if u.inputContext == "output" {
            title = "jmp on " + u.outputBuffer + greenCursor()
        } else {
            title = "jmp in " + u.searchBuffer + greenCursor()
        }
    }

    footer := "[↑][↓] [←][→] select"
    if u.Searcher != nil {
        footer += "  [i]n search"
    }
    if u.Outputer != nil {
        footer += "  [o]n cmd"
    }
    footer += "  [h]elp  [q]uit"
    footer = u.renderFooter(footer, "go-jmp v"+version.VERSION, innerWidth)

    title = fitToWidth(title, innerWidth)

    resultLines := make([]string, 0, resultsHeight)
    if len(u.Hits) == 0 {
        resultLines = append(resultLines, fitToWidth("(no output lines)", innerWidth))
        for len(resultLines) < resultsHeight {
            resultLines = append(resultLines, strings.Repeat(" ", innerWidth))
        }
    } else {
        start := u.resultsWindowStart(resultsHeight)
        for row := 0; row < resultsHeight; row++ {
            i := start + row
            if i >= len(u.Hits) {
                resultLines = append(resultLines, strings.Repeat(" ", innerWidth))
                continue
            }

            hit := u.Hits[i]
            prefix := "  "
            selected := u.focus == "results" && i == u.selectedIndex
            if selected {
                prefix = "> "
            }
            line := fitToWidth(prefix+hit.Render(), innerWidth)
            if selected {
                line = greenInverse(line)
            }
            resultLines = append(resultLines, line)
        }
    }

    previewLines := make([]string, 0, previewHeight)
    previewOffset := u.previewWindowOffset(previewHeight)
    for row := 0; row < previewHeight; row++ {
        idx := previewOffset + row
        if idx >= len(u.previewLines) {
            previewLines = append(previewLines, strings.Repeat(" ", innerWidth))
            continue
        }

        line := u.previewLines[idx]
        current := u.previewStart + idx
        prefix := "  "
        selected := u.focus == "preview" && current == u.previewLine
        if selected {
            prefix = "> "
        }
        renderedLine := fitToWidth(prefix+line, innerWidth)
        if selected {
            renderedLine = greenInverse(renderedLine)
        }
        previewLines = append(previewLines, renderedLine)
    }

    rendered := make([]string, 0, 4+len(resultLines)+len(previewLines))
    rendered = append(rendered, "╔"+strings.Repeat("═", innerWidth)+"╗")
    rendered = append(rendered, "║"+title+"║")
    rendered = append(rendered, "╟"+strings.Repeat("─", innerWidth)+"╢")
    for _, line := range resultLines {
        rendered = append(rendered, "║"+line+"║")
    }
    rendered = append(rendered, "╟"+strings.Repeat("─", innerWidth)+"╢")
    for _, line := range previewLines {
        rendered = append(rendered, "║"+line+"║")
    }
    rendered = append(rendered, "╟"+strings.Repeat("─", innerWidth)+"╢")
    rendered = append(rendered, "║"+footer+"║")
    rendered = append(rendered, "╚"+strings.Repeat("═", innerWidth)+"╝")

    return strings.Join(rendered, "\n")
}

func (u *UI) resultsWindowStart(windowSize int) int {
    if len(u.Hits) <= windowSize {
        return 0
    }

    half := windowSize / 2
    start := u.selectedIndex - half
    if start < 0 {
        start = 0
    }
    maxStart := len(u.Hits) - windowSize
    if start > maxStart {
        start = maxStart
    }

    return start
}

func (u *UI) previewWindowOffset(windowSize int) int {
    if len(u.previewLines) <= windowSize {
        return 0
    }

    target := u.previewLine - u.previewStart
    if target < 0 {
        target = 0
    }

    half := windowSize / 2
    start := target - half
    if start < 0 {
        start = 0
    }
    maxStart := len(u.previewLines) - windowSize
    if start > maxStart {
        start = maxStart
    }

    return start
}

func (u *UI) paneHeights(totalHeight int) (int, int) {
    // Full frame uses: top border, title, sep, results, sep, preview, sep, footer, bottom border.
    const nonContentRows = 8

    contentRows := totalHeight - nonContentRows
    if contentRows < 2 {
        return 1, 1
    }

    // Keep results pane at 35% of the total terminal height.
    resultsHeight := (totalHeight * 35) / 100

    if resultsHeight < 1 {
        resultsHeight = 1
    }

    // Leave at least one row for preview.
    if resultsHeight > contentRows-1 {
        resultsHeight = contentRows - 1
    }

    previewHeight := contentRows - resultsHeight
    if previewHeight < 1 {
        previewHeight = 1
    }

    return resultsHeight, previewHeight
}

func (u *UI) renderFooter(actions, versionText string, width int) string {
    versionWidth := utf8.RuneCountInString(versionText)
    actionsWidth := width - versionWidth - 1
    if actionsWidth <= 0 {
        return fitToWidth(versionText, width)
    }

    actions = clipToWidth(actions, actionsWidth)
    actionsLen := utf8.RuneCountInString(actions)
    paddingLeft := 0
    paddingRight := 0
    if actionsWidth > actionsLen {
        totalPadding := actionsWidth - actionsLen
        paddingLeft = totalPadding / 2
        paddingRight = totalPadding - paddingLeft
    }

    left := strings.Repeat(" ", paddingLeft) + actions + strings.Repeat(" ", paddingRight)
    left = fitToWidth(left, actionsWidth)

    return left + " " + versionText
}

func fitToWidth(line string, width int) string {
    runes := []rune(line)
    if len(runes) > width {
        return string(runes[:width])
    }
    if len(runes) < width {
        return line + strings.Repeat(" ", width-len(runes))
    }
    return line
}

func clipToWidth(line string, width int) string {
    runes := []rune(line)
    if len(runes) > width {
        return string(runes[:width])
    }
    return line
}

func greenInverse(line string) string {
    return "\x1b[7;32m" + line + "\x1b[0m"
}

func greenCursor() string {
    return "\x1b[32m▌\x1b[0m"
}

func (u *UI) Display() error {
    model := u
    program := tea.NewProgram(model, tea.WithAltScreen())
    _, err := program.Run()
    fmt.Print("\033[H\033[2J")
    return err
}
