# Pane Cursor Counter Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Show the active pane's cursor position (`N/M`, 1-indexed) in the bottom-left of the TUI footer, updating every time the cursor moves.

**Architecture:** Two surgical changes confined to `pkg/tui/ui.go` — a new pure helper `cursorPositionText()` that reads existing `UI` state, and a refactor of `renderFooter` from a two-zone (actions + version) to a three-zone (counter + actions + version) layout. No new fields, no new dependencies, no changes to cursor movement logic.

**Tech Stack:** Go 1.24, Bubbletea (already a dependency).

**Spec:** `docs/superpowers/specs/2026-04-20-pane-cursor-counter-design.md`

---

## File Structure

| File | Change type | Responsibility |
|------|-------------|----------------|
| `pkg/tui/ui.go` | Modify | Add `cursorPositionText()` method (~15 LOC); change `renderFooter` signature and body; update the single call site in `View()`. |
| `pkg/tui/ui_test.go` | Modify | Add tests for `cursorPositionText` (both panes + empty); add tests for `renderFooter` three-zone layout and width degradation; add one smoke test that confirms the counter shows up in `View()` output. |

---

## Task 1: `cursorPositionText` helper

Pure method on `UI` that formats the counter string from the currently-focused pane's state. No side effects, no dependencies on other new code — can be fully TDD'd in isolation.

**Files:**
- Modify: `pkg/tui/ui.go` (add one method, add `"fmt"` import if not already imported — it is, at line 4)
- Test: `pkg/tui/ui_test.go`

- [ ] **Step 1: Write the failing tests**

Add to `pkg/tui/ui_test.go`:

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

```
cd /home/s3/go-jmp.nigelhamilton.com
go test ./pkg/tui/ -run 'TestCursorPositionText' -v
```

Expected: compile error `ui.cursorPositionText undefined`.

- [ ] **Step 3: Implement `cursorPositionText`**

Add this method to `pkg/tui/ui.go` immediately after `PreviewWindow` (after line 110):

```go
func (u *UI) cursorPositionText() string {
    if u.focus == "preview" {
        total := len(u.previewLines)
        if total < 1 {
            return "0/0"
        }
        current := u.previewLine - u.previewStart + 1
        if current < 1 {
            current = 1
        }
        if current > total {
            current = total
        }
        return fmt.Sprintf("%d/%d", current, total)
    }

    total := len(u.Hits)
    if total < 1 {
        return "0/0"
    }
    current := u.selectedIndex + 1
    if current < 1 {
        current = 1
    }
    if current > total {
        current = total
    }
    return fmt.Sprintf("%d/%d", current, total)
}
```

The clamps are defensive — current cursor-movement code already keeps `selectedIndex` and `previewLine` in range, but the helper stays safe if that ever regresses.

- [ ] **Step 4: Run the new tests and verify they pass**

```
go test ./pkg/tui/ -run 'TestCursorPositionText' -v
```

Expected: both tests PASS.

- [ ] **Step 5: Run the full package test suite to confirm no regressions**

```
go test ./pkg/tui/ -v
```

Expected: all tests PASS, no failures.

- [ ] **Step 6: Commit**

```
git add pkg/tui/ui.go pkg/tui/ui_test.go
git commit -m "feat(tui): add cursorPositionText helper

Returns the active pane's cursor position as N/M (1-indexed). Reads
existing UI state; no new fields. Defensive clamps keep output sane
if upstream movement code ever lets indices drift out of range."
```

---

## Task 2: Three-zone footer + wire into `View()`

Change `renderFooter` from `(actions, version, width)` to `(leftText, actions, version, width)` and update the single call site in `View()` to pass `u.cursorPositionText()`. The caller change and the signature change must land together (they're in the same file and the code won't compile otherwise).

**Files:**
- Modify: `pkg/tui/ui.go:438` (call site in `View`)
- Modify: `pkg/tui/ui.go:582-603` (`renderFooter` body and signature)
- Test: `pkg/tui/ui_test.go`

- [ ] **Step 1: Write failing tests for the new three-zone layout**

Add to `pkg/tui/ui_test.go`:

```go
func TestRenderFooterThreeZones(t *testing.T) {
    ui := NewUI("jmp", nil, nil, nil, nil, nil)

    // Wide: all three zones present, space-separated.
    got := ui.renderFooter("1/5", "actions", "v1.0", 40)
    if len(got) != 40 {
        t.Fatalf("wide: expected length 40, got %d (%q)", len(got), got)
    }
    if got[:3] != "1/5" {
        t.Fatalf("wide: expected counter at left, got %q", got[:3])
    }
    if got[len(got)-4:] != "v1.0" {
        t.Fatalf("wide: expected version at right, got %q", got[len(got)-4:])
    }
    if !strings.Contains(got, "actions") {
        t.Fatalf("wide: expected actions in middle, got %q", got)
    }
}

func TestRenderFooterMiddleCollapses(t *testing.T) {
    ui := NewUI("jmp", nil, nil, nil, nil, nil)

    // Just enough width for counter + space + version, no room for actions.
    // "12/345" (6) + " " + "v69" (3) = 10.
    got := ui.renderFooter("12/345", "actions", "v69", 10)
    if len(got) != 10 {
        t.Fatalf("narrow: expected length 10, got %d (%q)", len(got), got)
    }
    if got[:6] != "12/345" {
        t.Fatalf("narrow: expected counter at left, got %q", got[:6])
    }
    if got[len(got)-3:] != "v69" {
        t.Fatalf("narrow: expected version at right, got %q", got[len(got)-3:])
    }
    if strings.Contains(got, "actions") {
        t.Fatalf("narrow: expected actions to be dropped, got %q", got)
    }
}

func TestRenderFooterTooNarrow(t *testing.T) {
    ui := NewUI("jmp", nil, nil, nil, nil, nil)

    // Width too small for even counter + space + version — fall back to
    // version only, truncated to width. Matches existing behaviour.
    got := ui.renderFooter("1/5", "actions", "v69", 3)
    if len(got) != 3 {
        t.Fatalf("tiny: expected length 3, got %d (%q)", len(got), got)
    }
    if got != "v69" {
        t.Fatalf("tiny: expected %q, got %q", "v69", got)
    }
}
```

Add `"strings"` to the test file's imports if it isn't already present (check first — current test file doesn't import it).

- [ ] **Step 2: Run tests to verify they fail**

```
go test ./pkg/tui/ -run 'TestRenderFooter' -v
```

Expected: compile errors (wrong number of args to `renderFooter`, and `strings` unused if not imported) and/or test failures.

- [ ] **Step 3: Replace the `renderFooter` body**

Replace the entire function at `pkg/tui/ui.go:582-603` with:

```go
func (u *UI) renderFooter(leftText, actions, versionText string, width int) string {
    leftWidth := utf8.RuneCountInString(leftText)
    versionWidth := utf8.RuneCountInString(versionText)

    // Degenerate: not enough room for even left + " " + version.
    // Fall back to version-only (matches pre-existing behaviour).
    if leftWidth+1+versionWidth > width {
        return fitToWidth(versionText, width)
    }

    // Layout: [left][" "][middle][" "][version]
    // Middle is where centered actions live (if any room remains).
    middleWidth := width - leftWidth - versionWidth - 2

    if middleWidth <= 0 {
        // No room for actions; absorb the slack between left and version.
        gap := width - leftWidth - versionWidth
        return leftText + strings.Repeat(" ", gap) + versionText
    }

    actions = clipToWidth(actions, middleWidth)
    actionsLen := utf8.RuneCountInString(actions)
    totalPadding := middleWidth - actionsLen
    paddingLeft := totalPadding / 2
    paddingRight := totalPadding - paddingLeft

    middle := strings.Repeat(" ", paddingLeft) + actions + strings.Repeat(" ", paddingRight)

    return leftText + " " + middle + " " + versionText
}
```

- [ ] **Step 4: Update the call site in `View()`**

At `pkg/tui/ui.go:438`, change:

```go
    footer = u.renderFooter(footer, "go-jmp v"+version.VERSION, innerWidth)
```

to:

```go
    footer = u.renderFooter(u.cursorPositionText(), footer, "go-jmp v"+version.VERSION, innerWidth)
```

- [ ] **Step 5: Run the new renderFooter tests**

```
go test ./pkg/tui/ -run 'TestRenderFooter' -v
```

Expected: all three PASS.

- [ ] **Step 6: Run the full package test suite**

```
go test ./pkg/tui/ -v
```

Expected: all tests PASS (including pre-existing `TestClampPreviewLine`, `TestPreviewWindow`, `TestPaneHeights`, `TestLoadCommandOutput`, and the Task-1 cursor-position tests).

- [ ] **Step 7: Add an integration smoke test for `View()`**

Append to `pkg/tui/ui_test.go`:

```go
func TestViewFooterContainsCursorCounter(t *testing.T) {
    hits := []model.Renderable{
        file.NewHit("a", "a", 1, ""),
        file.NewHit("b", "b", 1, ""),
        file.NewHit("c", "c", 1, ""),
    }
    ui := NewUI("jmp", nil, hits, nil, nil, nil)
    ui.width = 80
    ui.height = 24
    ui.selectedIndex = 1

    out := ui.View()
    if !strings.Contains(out, "2/3") {
        t.Fatalf("expected rendered output to contain counter \"2/3\", got:\n%s", out)
    }
}
```

- [ ] **Step 8: Run the integration test**

```
go test ./pkg/tui/ -run 'TestViewFooterContainsCursorCounter' -v
```

Expected: PASS.

- [ ] **Step 9: Build the binary to confirm no downstream compile breakage**

```
go build ./cmd/jmp
```

Expected: no output, exit 0, `jmp` binary produced at repo root.

- [ ] **Step 10: Run the whole test suite one more time**

```
go test ./...
```

Expected: all tests PASS across every package.

- [ ] **Step 11: Commit**

```
git add pkg/tui/ui.go pkg/tui/ui_test.go
git commit -m "feat(tui): show pane cursor position in footer

Footer now has three zones: pane cursor counter (N/M) on the left,
action hints centered, version on the right. Counter reflects the
active pane — results pane shows selectedIndex+1 of len(Hits),
preview pane shows position within the loaded window.

When width shrinks, actions collapse first; counter and version are
always preserved. A degenerate case (too narrow for even counter +
version) keeps the pre-existing version-only fallback."
```

---

## Self-Review

Checked against the spec at `docs/superpowers/specs/2026-04-20-pane-cursor-counter-design.md`:

1. **Spec coverage**
   - Semantics table (results + preview formulas, empty-pane `0/0`) → covered by `TestCursorPositionTextResults` + `TestCursorPositionTextPreview`, implemented in Task 1 Step 3.
   - Input-mode orthogonality → no code change needed (focus drives the helper regardless of mode); implicit but correct by construction.
   - Three-zone footer layout + degradation rules → covered by `TestRenderFooterThreeZones` + `TestRenderFooterMiddleCollapses` + `TestRenderFooterTooNarrow`, implemented in Task 2 Step 3.
   - `View()` wiring → Task 2 Step 4 + `TestViewFooterContainsCursorCounter`.
   - Single helper on `UI` (`cursorPositionText`) → Task 1.
   - `renderFooter` signature gains leading `leftText` → Task 2 Step 3.

2. **Placeholder scan** — no TBDs, no "similar to above", every code step has complete code and every command step has the exact command.

3. **Type / name consistency** — `cursorPositionText`, `renderFooter(leftText, actions, versionText, width)`, `u.focus == "preview"`, `u.previewLines`, `u.previewStart`, `u.previewLine`, `u.selectedIndex`, `u.Hits` — all consistent between tasks and match the live code in `pkg/tui/ui.go`.

No gaps found.
