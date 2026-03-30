# HabitClaw 🦞

> Self-hosted habit tracker with MCP integration for AI assistants.
> Open source. Single Go binary. Local-first.

HabitClaw runs on your machine as a single binary. It exposes a web dashboard for manual tracking and an MCP server so AI assistants like [OpenClaw](https://openclaw.ai), Claude Desktop, and Cursor can log habits, check streaks, and manage pauses through natural conversation.

Your data stays on your machine. No account required. No cloud dependency.

---

## Features

- **Build habits** — track duration, count, or boolean completion
- **Avoid habits** — track slips with severity (minor keeps streak, full breaks it)
- **Flexible schedules** — daily, specific days, weekdays, N times/week, monthly
- **Consistency score** — percentage over 7/30/90 days, doesn't collapse on one missed day
- **Grace pauses** — schedule vacation or sick days, streaks unaffected
- **Goal alignment** — see how well you're meeting your targets
- **MCP server** — AI assistants can log habits and query stats via natural language
- **Web dashboard** — simple UI at `localhost:3000`

---

## Quick Start

```bash
# Install (one-liner)
curl -fsSL https://habitclaw.app/install.sh | bash

# Or with Go
go install github.com/you/habitclaw@latest

# Run
habitclaw
# → opens localhost:3000
# → MCP ready at localhost:3000/mcp
```

---

## OpenClaw Integration

Add to your OpenClaw config:

```json
{
  "mcpServers": {
    "habitclaw": {
      "url": "http://localhost:3000/mcp"
    }
  }
}
```

Add to your `HEARTBEAT.md`:

```
- Check HabitClaw for incomplete habits today
- If any habits not done and it's past 8pm, remind me
- If I'm on a streak, mention it for motivation
```

Then just talk to your AI:

```
You: "Done reading for 1 hour"
AI:  "Logged! You're on a 5-day reading streak 🔥"

You: "I'm going on vacation April 1-7"
AI:  "Paused all habits Apr 1–7. Your streaks are safe."

You: "How's my week looking?"
AI:  "Reading: 6/7 days (86%). Problem solving: 4/5 days (80%).
      No social media: 7/7 clean days. Strong week overall."
```

---

## MCP Tools

| Tool | Description |
|---|---|
| `list_habits` | List all active habits with schedules |
| `add_habit` | Create a new habit |
| `log_habit` | Log a completion (build habits) |
| `log_slip` | Log a slip with severity (avoid habits) |
| `get_streaks` | Get streaks and consistency scores |
| `get_summary` | Daily, weekly, or monthly summary |
| `get_goal_alignment` | How well you're meeting your goals |
| `schedule_pause` | Pause tracking for vacation/sick days |
| `cancel_pause` | Cancel a pause early |
| `list_pauses` | List upcoming and past pauses |

---

## Configuration

| Variable | Default | Description |
|---|---|---|
| `HABITCLAW_PORT` | `3000` | HTTP server port |
| `HABITCLAW_DB_PATH` | `./habitclaw.db` | SQLite database path |
| `HABITCLAW_PASSWORD` | `` | Optional dashboard password |
| `HABITCLAW_LOG_LEVEL` | `info` | Log level |

---

## Build from Source

```bash
git clone https://github.com/you/habitclaw.git
cd habitclaw
go build -o habitclaw .
./habitclaw
```

Requirements: Go 1.21+

---

## How It Works

HabitClaw runs two servers on the same port:

```
localhost:3000/      →  Web dashboard
localhost:3000/mcp   →  MCP server for AI assistants
```

Both share the same SQLite database. Data never leaves your machine.

Reminders and scheduling are handled by your AI assistant (OpenClaw heartbeat/cron) — HabitClaw just stores your data and exposes clean tools.

---

## Cloud Version

[HabitClaw.app](https://habitclaw.app) is the hosted version — no setup required, OAuth login, Stripe billing. Same core, managed for you.

---

## License

GNU Affero General Public License v3.0 — free to self-host.

See [LICENSE](LICENSE) for details.
