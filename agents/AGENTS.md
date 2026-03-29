# AGENTS.md

This repository is a Go codebase for the `jmp` command line utility.

This file is the top-level operating contract for Copilot, VS Code agents, and any other coding agent working in this repo.

## How agents should work here

- Read `agents/START_HERE.md` first.
- Preserve existing command semantics unless the task explicitly changes them.
- Prefer small, reviewable patches over broad rewrites.
- Keep core logic separate from CLI parsing, formatting, filesystem access, and shell integration.
- Protect protocol and behavioural invariants before improving UX.
- Use idiomatic Go design. Do not import Python or JavaScript habits into the codebase.
- Prefer small packages, typed structs, focused interfaces, and explicit boundaries.
- Do not add dependencies without a strong reason.
- Update tests and docs alongside behaviour changes.
- Do not claim success unless you actually ran the relevant checks.

## Read order

1. `agents/START_HERE.md`
2. `agents/PROJECT_INVARIANTS.md`
3. `agents/JMP_SPECIFIC.md`
4. `agents/GO_CONVENTIONS.md`
5. `agents/CHANGE_PROTOCOL.md`
6. `agents/TEST_AND_VERIFY.md`
7. The most relevant role file in `agents/`

## Good changes in this repo

Good changes are:

- small
- explicit
- test-backed
- CLI-safe
- easy to revert
- aligned with Go idioms
- careful about stdout, stderr, exit codes, and help text

## Bad changes in this repo

Avoid:

- broad rewrites without need
- hidden behavioural changes
- regex soup where explicit parsing belongs
- mixing cleanup with behaviour in one patch
- silent changes to output format
- new dependencies for trivial problems
- shelling out when built-in Go facilities are enough
- hidden global mutable state

## Multi-agent order of operations

Use this order when a task needs several lenses:

1. `ARCHITECTURE_ARCHIE.md`
2. `BUILDER_BRIAN.md`
3. `UX_ARIA.md`
4. `RELIABILITY_SEAN.md`
5. `RED_TEAM_NORBERT.md`

## Definition of done

A change is only done when:

- the patch is coherent
- the relevant tests were added or updated
- the touched code compiles
- the CLI contract remains intentional
- docs or examples were updated if needed
- the final summary states what was checked and what was not checked
