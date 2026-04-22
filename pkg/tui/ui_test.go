package tui

import (
    "testing"

    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/file"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/model"
)

func TestClampPreviewLine(t *testing.T) {
    ui := NewUI("jmp in test", nil, []model.Renderable{file.NewHit("testdata/1.txt", "testdata/1.txt", 1, "")}, nil, nil, nil)

    if got := ui.ClampPreviewLine(0, 10); got != 1 {
        t.Fatalf("expected 1, got %d", got)
    }
    if got := ui.ClampPreviewLine(11, 10); got != 10 {
        t.Fatalf("expected 10, got %d", got)
    }
    if got := ui.ClampPreviewLine(4, 10); got != 4 {
        t.Fatalf("expected 4, got %d", got)
    }
    if got := ui.ClampPreviewLine(5, 0); got != 1 {
        t.Fatalf("expected 1, got %d", got)
    }
}

func TestPreviewWindow(t *testing.T) {
    ui := NewUI("jmp", nil, nil, nil, nil, nil)

    assertWindow(t, ui, 0, 1, 400, 1, 1)
    assertWindow(t, ui, 100, 20, 400, 1, 100)
    assertWindow(t, ui, 1000, 25, 200, 1, 200)
    assertWindow(t, ui, 1000, 975, 200, 801, 1000)
    assertWindow(t, ui, 1000, 500, 200, 400, 599)
}

func TestPaneHeights(t *testing.T) {
    ui := NewUI("jmp", nil, nil, nil, nil, nil)

    results, preview := ui.paneHeights(24)
    if results != 8 || preview != 8 {
        t.Fatalf("expected (8,8) for total height 24, got (%d,%d)", results, preview)
    }

    results, preview = ui.paneHeights(40)
    if results != 14 || preview != 18 {
        t.Fatalf("expected (14,18) for total height 40, got (%d,%d)", results, preview)
    }

    results, preview = ui.paneHeights(8)
    if results != 1 || preview != 1 {
        t.Fatalf("expected (1,1) guardrails for total height 8, got (%d,%d)", results, preview)
    }
}

func TestLoadCommandOutput(t *testing.T) {
    sample := []model.Renderable{file.NewHitLater("initial line")}

    uiWithoutOutput := NewUI("jmp test", nil, sample, nil, nil, nil)
    if uiWithoutOutput.LoadCommandOutput("tail -n 10 /tmp/file.log") {
        t.Fatalf("expected false when no output callback is defined")
    }

    seenSearch := ""
    seenCommand := ""
    ui := NewUI("jmp test", nil, sample,
        func(terms string) ([]model.Renderable, error) {
            seenSearch = terms
            return []model.Renderable{file.NewHitLater("search:" + terms)}, nil
        },
        func(command string) ([]model.Renderable, error) {
            seenCommand = command
            return []model.Renderable{
                file.NewHitLater("out:" + command),
                file.NewHitLater("stderr: warning at testdata/1.txt line 1"),
            }, nil
        },
        nil,
    )

    if ui.LoadCommandOutput("   ") {
        t.Fatalf("expected false for blank command")
    }
    if !ui.LoadCommandOutput("  tail -n 1000 /tmp/file.log  ") {
        t.Fatalf("expected true for non-empty command")
    }
    if seenCommand != "tail -n 1000 /tmp/file.log" {
        t.Fatalf("unexpected trimmed command %q", seenCommand)
    }
    if ui.Title != "jmp on tail -n 1000 /tmp/file.log" {
        t.Fatalf("unexpected title %q", ui.Title)
    }
    if len(ui.Hits) != 2 {
        t.Fatalf("expected 2 hits, got %d", len(ui.Hits))
    }
    if !ui.LoadSearchResults("needle") || seenSearch != "needle" {
        t.Fatalf("expected search path to still work")
    }
}

func assertWindow(t *testing.T, ui *UI, maxLine, targetLine, windowSize, expectedStart, expectedEnd int) {
    t.Helper()

    start, end := ui.PreviewWindow(maxLine, targetLine, windowSize)
    if start != expectedStart || end != expectedEnd {
        t.Fatalf("expected (%d,%d), got (%d,%d)", expectedStart, expectedEnd, start, end)
    }
}

func TestCursorPositionTextResults(t *testing.T) {
    // Empty hits.
    ui := NewUI("jmp", nil, nil, nil, nil, nil)
    if got := ui.cursorPositionText(); got != "0/0" {
        t.Fatalf("empty results: expected \"0/0\", got %q", got)
    }

    // Five hits, cursor at start.
    hits := []model.Renderable{
        file.NewHit("a", "a", 1, ""),
        file.NewHit("b", "b", 1, ""),
        file.NewHit("c", "c", 1, ""),
        file.NewHit("d", "d", 1, ""),
        file.NewHit("e", "e", 1, ""),
    }
    ui = NewUI("jmp", nil, hits, nil, nil, nil)
    if got := ui.cursorPositionText(); got != "1/5" {
        t.Fatalf("start of 5: expected \"1/5\", got %q", got)
    }

    // Middle.
    ui.selectedIndex = 2
    if got := ui.cursorPositionText(); got != "3/5" {
        t.Fatalf("middle of 5: expected \"3/5\", got %q", got)
    }

    // End.
    ui.selectedIndex = 4
    if got := ui.cursorPositionText(); got != "5/5" {
        t.Fatalf("end of 5: expected \"5/5\", got %q", got)
    }
}

func TestCursorPositionTextPreview(t *testing.T) {
    ui := NewUI("jmp", nil, nil, nil, nil, nil)
    ui.focus = "preview"

    // Ten loaded lines, cursor at first loaded line.
    ui.previewLines = make([]string, 10)
    ui.previewStart = 1
    ui.previewLine = 1
    if got := ui.cursorPositionText(); got != "1/10" {
        t.Fatalf("preview 1/10: got %q", got)
    }

    // Cursor in middle.
    ui.previewLine = 7
    if got := ui.cursorPositionText(); got != "7/10" {
        t.Fatalf("preview 7/10: got %q", got)
    }

    // Window that starts mid-file — counter reflects pane-relative
    // position, not file line number.
    ui.previewLines = make([]string, 50)
    ui.previewStart = 100
    ui.previewLine = 123
    if got := ui.cursorPositionText(); got != "24/50" {
        t.Fatalf("preview 24/50 (window at 100, line 123): got %q", got)
    }
}
