# Tokenkitty

**Open-source, self-hosted dashboard for what your team spends on AI coding agents.**

Your team runs Claude Code, Codex, Cursor and Copilot: four bills, four dashboards. Tokenkitty shows who, which repo and which model is spending, and the cost per merged PR, in one place you host yourself.

> **Status: pre-alpha.** Nothing is usable yet. Watch the repo or join the discussions to follow along.

## How it works

- A small **collector** runs on each developer machine (Windows, macOS, Linux). It reads the usage numbers coding agents already write locally: token counts, models and timestamps. It never reads prompts, responses or code.
- The **server** also imports usage from vendor admin APIs (Anthropic, Cursor, GitHub Copilot, OpenAI).
- The **dashboard** shows spend by developer, team, repo, model and tool, with budgets and Slack, email or webhook alerts.

You run it with `docker compose up`. Nothing leaves your infrastructure.

## Planned sources

| Tool | How |
|---|---|
| Claude Code | local logs |
| Codex CLI | local logs |
| Gemini CLI | local logs |
| OpenCode | local database (read-only) |
| Anthropic | Admin API |
| Cursor (Teams) | Admin API |
| GitHub Copilot | GitHub API |
| OpenAI | Usage and Costs API |

## Privacy

Tokenkitty collects usage metadata only: token counts, model, timestamps, the repo name from the git remote, and who ran it. It never collects prompt text, model output, file contents, file paths or shell commands. The collector will have a `--dry-run` mode that prints exactly what it would send.

## Contributing

Contributions are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md) first. Every commit needs a DCO sign-off (`git commit -s`).

To report a security issue, follow [SECURITY.md](SECURITY.md). Please don't open a public issue.

## Acknowledgements

[ccusage](https://github.com/ryoppippi/ccusage) is used as a reference for checking token totals. Tokenkitty has its own parsers and shares no code with it.

## License

[Apache-2.0](LICENSE)
