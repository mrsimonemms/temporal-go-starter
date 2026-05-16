# CLAUDE.md

Guidance for Claude and other coding agents working in this repository.

## What this repository is

`temporal-go-starter` is a customer-facing scaffold for building Temporal
applications in Go. It is shared with Temporal users as a starting point for
their own services. Treat it as starter infrastructure, not a throwaway demo
or a personal project.

The sample workflow and activity are intentionally trivial. Their job is to
demonstrate the project structure, tooling, and idiomatic patterns. The value
of the repository is the setup, not the business logic.

## How to think about changes

* Keep the sample application small and easy to read.
* Prefer small, explicit, idiomatic Go over clever abstractions.
* Do not add features, layers, or dependencies beyond what the scaffold needs
  to demonstrate its purpose. If something looks like it might be useful later,
  do not add it until there is a clear reason.
* Maintain compatibility with the Dev Containers workflow. Dev Containers are
  the canonical development environment for this project. Changes that
  meaningfully degrade that experience should be raised with the maintainer
  before being merged.
* Keep the README accurate. If you change a command, a file path, a port, or
  a behaviour that the README describes, update the README in the same change.

## Temporal-specific rules

These rules apply whenever you touch code in `internal/app/workflows/` or
register new workflow code with the worker.

* **Workflows must be deterministic.** Two replays of the same workflow with
  the same history must produce the same result.
* **Do not use wall-clock time inside a workflow.** Use `workflow.Now`,
  `workflow.Sleep`, and `workflow.NewTimer` instead of `time.Now`, `time.Sleep`,
  and `time.NewTimer`.
* **Do not use randomness inside a workflow.** Use `workflow.SideEffect` or
  generate the value in an activity. Do not call `math/rand`, `crypto/rand`,
  or `uuid.NewString` from workflow code.
* **Do not perform I/O inside a workflow.** No network calls, file system
  access, database access, or external commands. Put those in activities.
* **Do not spawn goroutines from a workflow.** Use `workflow.Go` if you
  genuinely need concurrent execution inside the workflow.
* **Do not read environment variables or global mutable state from a
  workflow.** Pass everything you need as input, or fetch it from an activity.
* **Keep side effects inside activities.** Activities are the right place for
  retries, timeouts, and non-deterministic work.
* **Be careful when changing existing workflow code.** Changing the shape of
  a workflow can break in-flight executions. When in doubt, version the
  workflow or ask the maintainer.

## Code style

* Use British English spelling in documentation and comments where natural
  ("behaviour", "initialise", "organisation"). Code identifiers and existing
  third-party names should remain as they are.
* Do not use em dashes. Prefer commas, full stops, or parentheses.
* Avoid suppressing lint findings. If you need to use `//nolint`, include a
  short comment explaining why.
* Prefer `assert` from `github.com/stretchr/testify/assert` over `require` in
  tests unless the test cannot continue after the assertion fails.
* Use table-driven tests where you have more than a couple of cases.

## Mandatory validation

No task is complete until formatting, tests, and linting all pass. Before
declaring a piece of work done, run:

```sh
pre-commit run -a
go test ./...
golangci-lint run
```

If `pre-commit` is not installed in the environment:

```sh
pip install pre-commit
pre-commit install
```

If any of the three commands fails, fix the underlying issue rather than
suppressing it. Do not commit `//nolint` directives, skipped tests, or
disabled hooks to make a check pass.

## Pull request hygiene

* Use [Conventional Commits](https://www.conventionalcommits.org) for commit
  messages. The `conventional-pre-commit` hook enforces this on commit
  messages.
* Keep commits focused. One logical change per commit.
* If you touch documentation, run `pre-commit run -a markdown-toc` so that
  the table of contents stays in sync with the headings.
