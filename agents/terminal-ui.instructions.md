---
applyTo: "pkg/tui/ui.go"
---

# Bubble Tea Best Practices for JMP

Follow these rules whenever editing the Go TUI implementation.

## Layout contract

- Keep the same conceptual panes: title, results, preview, footer.
- Preserve user-facing key semantics while adapting to Bubble Tea event handling.
- Handle terminal resize (`tea.WindowSizeMsg`) and avoid hard-coding full-screen assumptions.

## Model state

- Keep model state explicit: selected result, focus context, input mode, preview target.
- Avoid hidden global state.
- Keep rendering derived from model state only.

## Key handling

- Keep `t` and `o` input prompts context-sensitive.
- Keep `Right`/`Enter` behavior split by focus: results -> preview, preview -> edit.
- Keep `Left`/`Esc` behavior split by focus: preview -> results, results -> quit.

## Editor handoff

- Launch editor in foreground and restore TUI state after return.
- Do not lose selection context when editor exits.
- Keep memory save behavior attached to successful edit flow.

## Preview math

- Keep line-clamp and preview-window logic deterministic and unit-tested.
- Preserve behavior at file start, file end, mid-file, and empty-file boundaries.

## Output stability

- Keep title and footer wording stable unless a task explicitly changes copy.
- Keep version placement in footer stable.

## Testing

- Keep unit tests for non-interactive methods (`ClampPreviewLine`, `PreviewWindow`, command/search loaders).
- Add integration/smoke tests for keyboard flow when changing navigation behavior.
