# Contributing to Tokenkitty

Thanks for helping. This guide covers setup, the rules every change follows, and how to get a pull request merged.

## Before you start

- For anything bigger than a small fix, open an issue first so we can agree on the approach.
- Want to add a new data source? Use the **New data source** issue template. Sources need extra checks (see below).
- Questions go to [Discussions](https://github.com/tokenkitty/tokenkitty/discussions), not issues.

## Development setup

| Tool | Version | Windows install |
|---|---|---|
| Git | latest | `winget install Git.Git` |
| Go | 1.23+ | `winget install GoLang.Go` |
| Node | 20 LTS or newer | `winget install OpenJS.NodeJS.LTS` |
| Task | 3.x | `winget install Task.Task` |
| Docker Desktop | latest | `winget install Docker.DockerDesktop` |
| golangci-lint | see `Taskfile.yml` | `winget install GolangCI.golangci-lint` |
| sqlc | see `Taskfile.yml` | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@<version>` |
| golang-migrate | see `Taskfile.yml` | `go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@<version>` |
| air | see `Taskfile.yml` | `go install github.com/air-verse/air@<version>` |

On macOS and Linux, install the same tools with Homebrew or your package manager.

Then:

```sh
task doctor   # checks every tool above
task dev      # postgres in docker, server with live reload, web dev server
task test     # all Go and web tests
task lint     # golangci-lint, eslint, tsc
```

The same commands work on Windows, macOS and Linux.

## Rules for every change

1. **Privacy first.** Tokenkitty never collects prompt text, model output, file contents, file paths or shell commands. Parsers decode only named usage fields. Every parser has a canary test that fails if content leaks into its output.
2. **No guessing formats.** Agent log formats are undocumented and change between versions. A parser is written against real, verified files, with anonymized fixtures and golden tests.
3. **Idempotent ingest.** Every event has a deterministic `event_id`. Reading the same file twice must never create duplicates.
4. **Cross-platform.** Use `filepath`, `os.UserHomeDir` and `os.UserConfigDir`. Never hard-code `/` paths. CI runs on Windows, macOS and Linux.
5. **Tests are required** for parsers, API endpoints and aggregation queries. Prefer table-driven Go tests. Tests never read your real home directory; use fixtures and `t.TempDir()`.
6. **Money is never a float.** Use integer micro-dollars or a decimal type.
7. **New dependencies need a reason** in the PR description. Prefer the standard library.

## Adding a data source

Before writing a parser:

1. Run 2–3 real sessions with the agent, with different models and at least one cache hit.
2. Find where it stores usage on Windows, macOS and Linux.
3. Identify the exact fields for timestamps, session and turn IDs, model and each token count.
4. List every field that holds **content** (prompts, responses, tool input and output, paths, commands). These go in the source's README as forbidden fields.
5. Commit 2–3 fixtures with all content replaced by `"<redacted>"`.
6. Compare totals with a reference tool where one exists, and state the result in the PR.
7. Record the agent version you tested.

## Commits

- Keep commits small and focused: one logical change each.
- Use [Conventional Commits](https://www.conventionalcommits.org): `feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `chore:`, `ci:`, with an optional scope such as `feat(collector):`.
- Keep the subject short and imperative, for example `fix(codex): skip empty token events`.

### Sign-off (DCO)

Every commit must be signed off under the [Developer Certificate of Origin](https://developercertificate.org). This certifies that you wrote the change or have the right to submit it under the project's license.

```sh
git commit -s -m "fix(codex): skip empty token events"
```

This adds a `Signed-off-by: Your Name <you@example.com>` line that matches your commit author. A check blocks PRs with unsigned-off commits. To fix a branch, run `git rebase --signoff main` and force-push the branch.

## Pull requests

- One topic per PR. Small PRs get reviewed faster.
- Fill in the PR template, including **what you verified manually** (for example "parsed my real Codex sessions on Windows: 41 events, totals match the reference").
- CI must pass on all three operating systems.
- PRs are squash-merged into `main`, so the PR title becomes the commit message. Write it as a Conventional Commit.

## Code of conduct

Everyone taking part follows the [Code of Conduct](CODE_OF_CONDUCT.md).

## License

By contributing, you agree that your contributions are licensed under the [Apache-2.0 License](LICENSE).
