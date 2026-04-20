# Pane Cursor Counter in Footer

Date: 2026-04-20
Status: Approved (pending spec review)

## Problem

When moving the cursor up and down in the TUI, the user has no visual indication
of how far through the pane they are. For long result sets or large previews,
it's easy to lose track of position.

## Goal

Display the active pane's cursor position and its total item count in the
bottom-left corner of the footer, in the form `N/M` (e.g. `12/234`). The
counter must update on every cursor move, reflect whichever pane holds focus,
and use 1-based indexing.

## Non-goals

- No new navigation commands (page jump, go-to-line, etc.).
- No display of the file's true line count when previewing; the counter
  reflects pane content only, not the underlying file.
- No changes to cursor movement logic or pane layout.
- No color or styling change to the footer beyond the added text.

## Semantics

The counter is computed from state that already exists on `UI` in
`pkg/tui/ui.go`. No new fields are introduced.

| `u.focus`   | N (current)                              | M (total)            |
|-------------|------------------------------------------|----------------------|
| `results`   | `u.selectedIndex + 1`                    | `len(u.Hits)`        |
| `preview`   | `u.previewLine - u.previewStart + 1`     | `len(u.previewLines)`|

Edge cases:

- **Empty pane** (`M == 0`): display `0/0`. Applies when `u.Hits` is empty on
  the results pane. The preview pane always has at least the intro-tips
  content, so `M > 0` there in practice.
- **Input mode**: `u.mode == "input"` is orthogonal to `u.focus`. The counter
  continues to reflect the active focus — it does not change or disappear
  while the user is typing a search/command/locate buffer.
- **Help screen**: displayed by setting `focus = "preview"` and replacing
  `u.previewLines` with help text. The preview-pane formula applies naturally
  (`1 / len(help_lines)` at entry).

## Footer layout

`renderFooter` in `pkg/tui/ui.go:582` currently produces two zones:

```
[        actions (centered)        ] [version]
```

After this change, the footer has three zones:

```
[counter] [       actions (centered in remaining middle)       ] [version]
```

Width allocation, in order of priority as the available width shrinks:

1. **Normal width** — counter on the left, version on the right, actions
   centered between them with a single space separating each adjacent zone.
2. **Middle collapses** — if `counter + actions + version + separators`
   exceeds the available width, drop the actions entirely. Counter stays
   left-aligned, version stays right-aligned, with one space gap at minimum.
3. **Still too narrow** — if `counter + version + one space` still doesn't
   fit, fall back to the existing behaviour: show the version truncated to
   the full width, nothing else. (The counter is the new feature, the
   version is the pre-existing contract; preserving current behaviour in
   the degenerate case minimises risk.)

## Implementation outline

Changes are confined to `pkg/tui/ui.go` and `pkg/tui/ui_test.go`.

1. **New helper** on `UI`:

   ```go
   func (u *UI) cursorPositionText() string
   ```

   Returns the `N/M` string for the currently focused pane, following the
   semantics table above.

2. **`renderFooter` signature change**: accept a new leading `leftText`
   parameter:

   ```go
   func (u *UI) renderFooter(leftText, actions, versionText string, width int) string
   ```

   The function computes widths for `leftText` and `versionText`, centers
   `actions` in the remaining middle, and drops `actions` if the middle
   becomes non-positive. `leftText` and `versionText` are always preserved
   (truncated only if individually wider than the total width).

3. **`View` call-site** updated to pass `u.cursorPositionText()` as the new
   left zone:

   ```go
   footer = u.renderFooter(u.cursorPositionText(), footer, "go-jmp v"+version.VERSION, innerWidth)
   ```

## Tests

All additions go in `pkg/tui/ui_test.go`.

1. **`cursorPositionText` — results focus**
   - 0 hits → `0/0`
   - 5 hits, `selectedIndex = 0` → `1/5`
   - 5 hits, `selectedIndex = 2` → `3/5`
   - 5 hits, `selectedIndex = 4` → `5/5`

2. **`cursorPositionText` — preview focus**
   - `previewStart = 1`, `previewLine = 1`, 10 loaded lines → `1/10`
   - `previewStart = 1`, `previewLine = 7`, 10 loaded lines → `7/10`
   - `previewStart = 100`, `previewLine = 123`, 50 loaded lines → `24/50`
     (tests that the counter reflects pane offset, not file line number)

3. **`renderFooter` layout**
   - Wide enough for all three zones: counter on far left, version on far
     right, actions centered in middle, single-space separators.
   - Narrow enough that the middle span collapses: actions dropped, counter
     and version remain, total output fits `width` exactly.
   - Width so small that even `counter + " " + version` doesn't fit:
     output is the version alone, truncated to the full width — matches
     the existing fallback.

4. **Integration via `View`**: render with known focus/state and assert the
   footer line contains the expected counter substring at the expected
   position (column 1 of the footer's inner content).

## Risks and mitigations

- **Width math regression** — the existing footer has a carefully balanced
  centering calculation. Mitigation: the new tests cover the width-degradation
  path explicitly.
- **Counter drift from pane state** — the helper reads directly from the
  authoritative fields (`focus`, `selectedIndex`, `Hits`, `previewLine`,
  `previewStart`, `previewLines`); there is no cache or copy, so drift is
  not possible.
- **Empty `Hits` producing `0/0`** — cosmetic only; consistent with the
  general "N/M" rule and avoids a special-case format.
