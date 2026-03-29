# GO_CONVENTIONS.md

Coding conventions for this Go repo.

## General

- Use 4-space indentation.
- Match the surrounding style where it is consistent.
- Preserve existing package boundaries unless the change requires restructuring.
- Prefer clear names over terse names.
- Keep functions focused and easy to scan.

## Typical layout

- `cmd/` for executables
- `internal/` for non-exported app logic
- `pkg/` for reusable package logic
- `testdata/` for fixtures
- `*_test.go` for tests

## Functions and types

- Prefer explicit function signatures where they improve clarity.
- Use concrete structs by default.
- Introduce interfaces only at package boundaries where they improve testability or decoupling.

## Go-native design

- Prefer constructor functions (`NewX`) for dependency setup.
- Keep side effects at the edge (filesystem, shell, environment, terminal).
- Keep parsing explicit and deterministic.
- Keep CLI dispatch thin and delegate behavior to package logic.

## Error handling

- Return `error` for exceptional conditions.
- Wrap errors with useful context.
- Do not silently ignore failures unless behavior intentionally mirrors upstream tools.

## Output and compatibility

- Keep user-facing output stable and concise.
- Preserve command names, argument meaning, and exit code semantics.
- Keep config and history formats backward compatible.

## Tooling

- Run `go test ./...` after behavior changes.
- Run `go build ./cmd/jmp` after CLI or package changes.
- Keep `go.mod` and `go.sum` tidy.
