# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

HabitClaw — self-hosted habit tracker with MCP integration.
Single Go binary. Local-first. No cloud dependency. AGPL-3.0.

Related repos (not in this tree):
- `habitclaw-cloud` — API server, OAuth, Stripe (private)
- `habitclaw-web` — Next.js frontend for cloud (private)

## Commands

```bash
go build -o habitclaw .          # build binary
go run .                         # run dev server (localhost:3000)
go test ./...                    # run all tests
go test ./core/...               # test core only
go test ./core/ -run TestStreak  # run a single test
go run . --port 3000             # run on specific port
```

## Architecture

Both the web dashboard (`/`) and MCP server (`/mcp`) are served on the same port, sharing a single SQLite database and business logic layer.

**Dependency direction:** `main.go` wires everything. `core/` contains pure business logic and the `HabitStore` interface. Adapters (`adapters/sqlite/`, `adapters/auth/`) implement interfaces. Transport layers (`mcp/`, `web/`) call `core.HabitService`.

**Key design:** Transport is injected at `main.go` — core never imports `net/http`.

## Critical Rules

- **`core/` must NEVER import `net/http`, database drivers, or any adapter package**
- All database access goes through the `HabitStore` interface only
- Streaks are always COMPUTED from completion logs — never cached in a separate table
- Paused days must be excluded from both streak and consistency calculations
- Minor slips (severity=1) do NOT break avoid-habit streaks — only full slips (severity=2) do
- No scheduler — reminders are the AI assistant's responsibility, not ours
- MCP tools must match the spec in `HABITCLAW_SPEC.md`

## Key Domain Types

```go
HabitDirection: build | avoid
GoalType:       duration | count | boolean
ScheduleType:   daily | specific_days | weekdays | weekends |
                times_per_week | every_n_days | weekly | monthly
LogType:        completion (build habits) | slip (avoid habits)
SlipSeverity:   1 = minor (streak intact) | 2 = full (streak breaks)
```

## Database Tables

```
users
habits              (habit_direction: build | avoid, has user_id)
habit_schedules     (has user_id — denormalized, not just habit_id)
completion_logs     (log_type: completion | slip, value: duration/count or slip severity, has user_id)
pauses              (habit_id nullable = all habits, from_date, to_date, reason, has user_id)
```

## MCP Tools

```
list_habits        habits:read     log_habit          habits:write
add_habit          habits:write    log_slip           habits:write
get_streaks        habits:read     get_summary        habits:read
get_goal_alignment habits:read     schedule_pause     pauses:write
cancel_pause       pauses:write    list_pauses        habits:read
```

## Environment Variables

```
HABITCLAW_PORT       3000              HTTP server port
HABITCLAW_DB_PATH    ./habitclaw.db    SQLite database path
HABITCLAW_PASSWORD                     Optional dashboard password
HABITCLAW_USER_ID    local             Local user ID
HABITCLAW_LOG_LEVEL  info              Log level
```

## Dependencies

- `modernc.org/sqlite` — pure Go SQLite driver (no CGO)
- `github.com/mark3labs/mcp-go` — MCP SDK

## Spec

Full specification: `HABITCLAW_SPEC.md`
Cloud specification: `HABITCLAW_CLOUD_SPEC.md`
