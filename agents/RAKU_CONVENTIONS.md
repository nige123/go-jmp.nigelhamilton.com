# GO_CONVENTIONS.md

Coding conventions for this Go repo.

## General

- Use 4-space indentation.
- Match the surrounding style where it is consistent.
- Preserve the repo's existing Go module and package layout unless a change requires updates.
- Prefer clear names over terse names.
- Keep routines focused and easy to scan.
- Avoid unnecessary metaprogramming.

## Typical layout

Use standard Go layout unless the repo already uses a different established pattern:

- `cmd/` for executables
- `internal/` for non-exported app logic
- `pkg/` for reusable package logic
- `testdata/` for fixtures
- `*_test.go` for tests

## Functions and types

- Prefer explicit signatures where they improve clarity.
- Use type constraints where they reduce ambiguity and real bugs.
- Do not add interfaces just for decoration.

## Go-native design

Prefer:

- small packages with clear ownership
- structs for domain models
- constructor functions (`NewX`) for initialization and dependency injection
- table-driven tests for behavior matrices
- explicit error returns with context

Avoid forcing patterns from other ecosystems onto Go.

## Parsing

If input has real structure, prefer explicit parsing with clear regexes or tokenization rather than chained ad hoc substitutions.

## Error handling

- Use clear domain outcomes for expected cases.
- Return `error` with context for exceptional conditions.
- Wrap errors when crossing package boundaries.
- Do not swallow failures silently.

## State and side effects

- Avoid hidden mutable global state.
- Inject dependencies where practical.
- Keep filesystem, environment, shell, and time access at the edge.

## Tooling

- Run `go test ./...` after behavior changes.
- Run `go build ./cmd/jmp` before release changes.
- Keep `go.mod` and `go.sum` tidy.

## Output

- Keep human-readable output concise and stable.
- Keep machine-readable output strict and predictable.
- Do not mix debug chatter into normal output paths.

## Comments and docs

- Comment the why, not the obvious what.
- Update examples when command syntax or behaviour changes.
