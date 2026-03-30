# HabitClaw Core — Specification

> Self-hosted habit tracker with MCP integration for AI assistants.
> Open source. Single binary. Local-first.

---

## Overview

HabitClaw is an open source habit tracking server that runs locally as a single binary. It exposes a web dashboard for manual tracking and an MCP server for AI assistant integration (OpenClaw, Claude Desktop, Cursor, etc.). Data is stored locally in SQLite. No cloud dependency, no account required.

The core has no AI logic. AI reasoning, reminders, and coaching live in the AI assistant (e.g. OpenClaw). HabitClaw just stores honest data and exposes clean tools.

---

## Philosophy

- **Local-first** — data never leaves your machine
- **Single binary** — one file, no runtime, no dependencies
- **AI-optional** — works standalone, works better with AI
- **MCP-native** — designed to be a first-class MCP server
- **No guilt** — consistency scores, not just streaks. Grace pauses. Slip severity.
- **Open core** — foundation for HabitClaw.ai cloud version

---

## Tech Stack

| Layer | Choice | Reason |
|---|---|---|
| Language | Go | Single binary, fast startup, easy cross-compile |
| Database | SQLite (`modernc.org/sqlite`) | Zero config, single file, no server needed |
| Web UI | Embedded HTML + HTMX (`embed.FS`) | No build step, all in binary |
| MCP server | `github.com/mark3labs/mcp-go` | Official Go MCP SDK |
| HTTP router | `net/http` stdlib | No external dependency for open source |
| Config | `.env` / env vars | Simple, standard |

---

## Features

### Core
- Build habits (do more of this) and avoid habits (do less of this)
- Flexible scheduling — daily, specific days, weekdays, weekends, N times/week, every N days, monthly
- Time-of-day targets with optional windows (e.g. anytime before 09:00)
- Log completions with actual value and notes
- Log slips with severity (minor / full) for avoid habits
- Streak calculation computed from logs — never cached, always accurate
- Consistency score (% over 7 / 30 / 90 days) — doesn't collapse on a missed day
- Grace pauses — schedule vacation or sick days, streaks and scores unaffected
- Daily, weekly, monthly summaries
- Goal alignment scoring

### Web Dashboard
- Today's habits with one-click check-off (build) and slip logging (avoid)
- Streak and consistency score per habit
- Weekly heatmap calendar
- Summary stats (completion rate, streaks, best day, goal alignment)
- Add / edit / archive habits with schedule config
- Pause management (upcoming, active, past)
- Accessible at `localhost:3000`

### MCP Server
- Full habit management via natural language
- AI can log completions, log slips, query streaks, get summaries, manage pauses
- Parallel tool calls supported — AI can pause multiple habits in one turn
- Accessible at `localhost:3000/mcp`
- Compatible with OpenClaw, Claude Desktop, Cursor, any MCP client

### Auth
- Single user, local only
- Optional password protection (env var `HABITCLAW_PASSWORD`)
- No JWT, no OAuth — not needed for local use

---

## Architecture

```
localhost:3000/        →  Web dashboard (embedded HTML/HTMX)
localhost:3000/mcp     →  MCP server (AI assistant integration)
```

Both endpoints share the same SQLite database and business logic via Go interfaces. The core never imports HTTP — transport is injected at `main.go`.

### Interfaces (ports)

```go
type HabitStore interface {
    // habits
    CreateHabit(habit Habit) error
    ListHabits(userID string) ([]Habit, error)
    UpdateHabit(habit Habit) error
    ArchiveHabit(id string) error

    // schedules
    SetSchedule(schedule HabitSchedule) error
    GetSchedule(habitID string) (HabitSchedule, error)
    // logs
    LogCompletion(log CompletionLog) error
    LogSlip(log CompletionLog) error
    GetLogs(userID string, from, to time.Time) ([]CompletionLog, error)
    DeleteLog(id string) error

    // pauses
    CreatePause(pause Pause) error
    ListPauses(userID string, status string) ([]Pause, error)
    CancelPause(id string, resumeFrom time.Time) error

    // computed
    GetStreaks(userID string) ([]Streak, error)
    GetSummary(userID string, period Period) (Summary, error)
    GetGoalAlignment(userID string, period Period) (GoalAlignment, error)
    GetHeatmap(habitID string, from, to time.Time) ([]HeatmapDay, error)
}
```

SQLite adapter implements `HabitStore`. Cloud version swaps in PostgreSQL adapter — core unchanged.

---

## Folder Structure

```
habitclaw/
├── main.go                    # entry point, wires adapters
├── go.mod
├── go.sum
│
├── core/
│   ├── habit.go               # Habit, HabitSchedule types
│   ├── log.go                 # CompletionLog, LogType, SlipSeverity types
│   ├── pause.go               # Pause type
│   ├── streak.go              # streak + consistency score computation
│   ├── summary.go             # daily/weekly/monthly summary logic
│   ├── alignment.go           # goal alignment scoring
│   ├── schedule.go            # scheduled day resolution logic
│   ├── interfaces.go          # HabitStore interface
│   └── service.go             # HabitService — pure business logic
│
├── adapters/
│   ├── sqlite/
│   │   ├── store.go           # SQLiteStore implements HabitStore
│   │   └── migrations/
│   │       ├── 001_init.sql
│   │       └── 002_pauses.sql
│   └── auth/
│       └── single_user.go     # optional local password check
│
├── mcp/
│   ├── server.go              # MCP server setup
│   └── handlers.go            # MCP tool handlers → calls core.HabitService
│
├── web/
│   ├── server.go              # HTTP server setup
│   ├── handlers.go            # HTTP handlers → calls core.HabitService
│   └── static/
│       ├── index.html         # dashboard UI
│       ├── app.js             # HTMX interactions
│       └── style.css
│
└── config/
    └── config.go              # env var loading
```

---

## Database Tables

### `users`
```sql
id          TEXT PRIMARY KEY
name        TEXT NOT NULL
created_at  DATETIME NOT NULL
```
Single row for open source. Cloud version has many rows.

---

### User ID scoping design

Every table that belongs to a user carries `user_id` directly — no table relies on a JOIN to determine ownership. This is a deliberate design decision:

```
open source   →  user_id = "local" (from HABITCLAW_USER_ID env var)
cloud         →  user_id = real user ID injected by JWT middleware
```

`HabitService` never decides what `user_id` is — it always receives it as a parameter from the caller. The caller (open source `main.go` or cloud Gin middleware) is responsible for injecting the correct value. Core logic never reads config or auth context directly.

```
Table               Has user_id
────────────────    ────────────
users               is the user (PK)
habits              yes
habit_schedules     yes (denormalized from habits for consistent scoping)
completion_logs     yes
pauses              yes
```

---

### `habits`
```sql
id                TEXT PRIMARY KEY
user_id           TEXT NOT NULL
name              TEXT NOT NULL
description       TEXT
goal_type         TEXT NOT NULL   -- duration | count | boolean
goal_value        INTEGER         -- minutes or count, 0 for boolean
habit_direction   TEXT NOT NULL   -- build | avoid
color             TEXT
icon              TEXT
archived_at       DATETIME
created_at        DATETIME NOT NULL
updated_at        DATETIME NOT NULL
```

---

### `habit_schedules`
```sql
id               TEXT PRIMARY KEY
habit_id         TEXT NOT NULL REFERENCES habits(id)
user_id          TEXT NOT NULL   -- denormalized for consistent user scoping
                                 -- open source: always "local"
                                 -- cloud: real user_id from JWT
schedule_type    TEXT NOT NULL
                 -- daily | specific_days | weekdays | weekends
                 -- times_per_week | every_n_days | weekly | monthly
days_of_week     TEXT            -- JSON array e.g. [1,2,4,5,6] Mon=1 Sun=7
times_per_week   INTEGER         -- used when schedule_type = times_per_week
every_n_days     INTEGER         -- used when schedule_type = every_n_days
day_of_month     INTEGER         -- used when schedule_type = monthly (1-31)
time_of_day      TEXT            -- target time e.g. "06:30", nullable
window_start     TEXT            -- earliest acceptable time e.g. "06:00"
window_end       TEXT            -- latest acceptable time e.g. "09:00"
created_at       DATETIME NOT NULL
updated_at       DATETIME NOT NULL
```

---

### `completion_logs`
```sql
id            TEXT PRIMARY KEY
habit_id      TEXT NOT NULL REFERENCES habits(id)
user_id       TEXT NOT NULL
log_type      TEXT NOT NULL   -- completion | slip
date          DATE NOT NULL   -- the day this log is for
value         INTEGER         -- duration/count for build; severity for avoid (1=minor 2=full)
note          TEXT
created_at    DATETIME NOT NULL
```

Avoid habits use `log_type = slip`. Minor slip = value 1 (streak intact). Full slip = value 2 (streak breaks).

---

### `pauses`
```sql
id            TEXT PRIMARY KEY
user_id       TEXT NOT NULL
habit_id      TEXT            -- NULL = pause all habits
from_date     DATE NOT NULL
to_date       DATE NOT NULL
reason        TEXT            -- vacation | sick | travel | other
cancelled_at  DATETIME        -- set if cancelled early
resume_from   DATE            -- set if cancelled early
created_at    DATETIME NOT NULL
```

---

## Data Models

### Habit

```go
type Habit struct {
    ID          string         `json:"id"`
    UserID      string         `json:"user_id"`
    Name        string         `json:"name"`
    Description string         `json:"description"`
    GoalType    GoalType       `json:"goal_type"`   // duration | count | boolean
    GoalValue   int            `json:"goal_value"`
    Direction   HabitDirection `json:"direction"`   // build | avoid
    Color       string         `json:"color"`
    Icon        string         `json:"icon"`
    Schedule    *HabitSchedule `json:"schedule,omitempty"`
    ArchivedAt  *time.Time     `json:"archived_at,omitempty"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
}
```

### HabitSchedule

```go
type HabitSchedule struct {
    ID           string       `json:"id"`
    HabitID      string       `json:"habit_id"`
    UserID       string       `json:"user_id"`        // always injected by caller, never set internally
    ScheduleType ScheduleType `json:"schedule_type"`
    DaysOfWeek   []int        `json:"days_of_week,omitempty"` // 1=Mon 7=Sun
    TimesPerWeek int          `json:"times_per_week,omitempty"`
    EveryNDays   int          `json:"every_n_days,omitempty"`
    DayOfMonth   int          `json:"day_of_month,omitempty"`
    TimeOfDay    string       `json:"time_of_day,omitempty"`  // "06:30"
    WindowStart  string       `json:"window_start,omitempty"` // "06:00"
    WindowEnd    string       `json:"window_end,omitempty"`   // "09:00"
}
```

### CompletionLog

```go
type CompletionLog struct {
    ID        string    `json:"id"`
    HabitID   string    `json:"habit_id"`
    UserID    string    `json:"user_id"`
    LogType   LogType   `json:"log_type"`  // completion | slip
    Date      time.Time `json:"date"`
    Value     int       `json:"value"`     // duration/count or slip severity (1=minor 2=full)
    Note      string    `json:"note"`
    CreatedAt time.Time `json:"created_at"`
}

type SlipSeverity int
const (
    SlipMinor SlipSeverity = 1  // streak intact
    SlipFull  SlipSeverity = 2  // streak breaks
)
```

### Pause

```go
type Pause struct {
    ID          string     `json:"id"`
    UserID      string     `json:"user_id"`
    HabitID     *string    `json:"habit_id,omitempty"` // nil = all habits
    FromDate    time.Time  `json:"from_date"`
    ToDate      time.Time  `json:"to_date"`
    Reason      string     `json:"reason"`
    CancelledAt *time.Time `json:"cancelled_at,omitempty"`
    ResumeFrom  *time.Time `json:"resume_from,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
}
```

### Streak (computed)

```go
type Streak struct {
    HabitID            string  `json:"habit_id"`
    HabitName          string  `json:"habit_name"`
    Direction          string  `json:"direction"`          // build | avoid
    Current            int     `json:"current"`
    Longest            int     `json:"longest"`
    ConsistencyRate7d  float64 `json:"consistency_rate_7d"`
    ConsistencyRate30d float64 `json:"consistency_rate_30d"`
    LastActivityDate   string  `json:"last_activity_date"`
}
```

Streaks are computed from `completion_logs` and `pauses` on demand. Paused days are excluded from both streak and consistency calculations. Minor slips do not break avoid habit streaks.

### Summary

```go
type Summary struct {
    Period         string      `json:"period"`
    TotalHabits    int         `json:"total_habits"`
    ScheduledDays  int         `json:"scheduled_days"`
    CompletionRate float64     `json:"completion_rate"`
    HabitSummaries []HabitStat `json:"habit_summaries"`
}

type HabitStat struct {
    HabitID       string  `json:"habit_id"`
    Name          string  `json:"name"`
    Direction     string  `json:"direction"`
    ScheduledDays int     `json:"scheduled_days"`
    CompletedDays int     `json:"completed_days"`
    SlippedDays   int     `json:"slipped_days,omitempty"`
    PausedDays    int     `json:"paused_days"`
    GoalMetRate   float64 `json:"goal_met_rate"`
    TotalValue    int     `json:"total_value"`
    CurrentStreak int     `json:"current_streak"`
}
```

---

## MCP Tools

All tools accessible at `localhost:3000/mcp`.

---

### `list_habits`

List all active habits with schedule and direction.

**Input:** none

**Output:**
```json
{
  "habits": [
    {
      "id": "abc123",
      "name": "Running",
      "direction": "build",
      "goal_type": "duration",
      "goal_value": 30,
      "schedule": {
        "schedule_type": "specific_days",
        "days_of_week": [1,2,4,5,6],
        "time_of_day": "06:30"
      }
    },
    {
      "id": "def456",
      "name": "No social media",
      "direction": "avoid",
      "goal_type": "boolean",
      "schedule": {
        "schedule_type": "daily"
      }
    }
  ]
}
```

---

### `add_habit`

Create a new habit with schedule and direction.

**Input:**
```json
{
  "name": "Running",
  "direction": "build",
  "goal_type": "duration",
  "goal_value": 30,
  "description": "Morning run",
  "schedule": {
    "schedule_type": "specific_days",
    "days_of_week": [1,2,4,5,6],
    "time_of_day": "06:30"
  }
}
```

**Output:**
```json
{
  "success": true,
  "habit_id": "abc123",
  "message": "Created: Running (30 min) — Mon Tue Thu Fri Sat at 06:30"
}
```

---

### `log_habit`

Log a completion for a build habit.

**Input:**
```json
{
  "habit_name": "Running",
  "value": 35,
  "note": "felt great today",
  "date": "2026-03-27"
}
```

**Output:**
```json
{
  "success": true,
  "habit": "Running",
  "value": 35,
  "current_streak": 8,
  "consistency_7d": 0.83,
  "message": "Logged 35 min of Running. 8-day streak!"
}
```

---

### `log_slip`

Log a slip for an avoid habit.

**Input:**
```json
{
  "habit_name": "No social media",
  "severity": "minor",
  "note": "checked Instagram for 10 mins at lunch",
  "date": "2026-03-27"
}
```

**Output:**
```json
{
  "success": true,
  "habit": "No social media",
  "severity": "minor",
  "current_streak": 12,
  "streak_status": "intact",
  "message": "Logged minor slip. Streak intact at 12 days — minor slips don't break it."
}
```

---

### `get_streaks`

Get streaks and consistency scores for all habits.

**Input:** none

**Output:**
```json
{
  "streaks": [
    {
      "habit": "Running",
      "direction": "build",
      "current": 8,
      "longest": 21,
      "consistency_7d": 0.83,
      "consistency_30d": 0.76,
      "last_activity": "2026-03-27"
    },
    {
      "habit": "No social media",
      "direction": "avoid",
      "current": 12,
      "longest": 30,
      "consistency_7d": 1.0,
      "consistency_30d": 0.9,
      "last_activity": "2026-03-27"
    }
  ]
}
```

---

### `get_summary`

Get habit summary for a period. Scheduled days exclude paused days.

**Input:**
```json
{
  "period": "today | week | month"
}
```

**Output:**
```json
{
  "period": "week",
  "completion_rate": 0.85,
  "total_habits": 3,
  "habit_summaries": [
    {
      "name": "Running",
      "direction": "build",
      "scheduled_days": 5,
      "completed_days": 4,
      "paused_days": 0,
      "goal_met_rate": 0.9,
      "total_value": 155,
      "current_streak": 8
    },
    {
      "name": "No social media",
      "direction": "avoid",
      "scheduled_days": 7,
      "slipped_days": 1,
      "paused_days": 0,
      "current_streak": 12
    }
  ]
}
```

---

### `get_goal_alignment`

Check how well the user is meeting their habit goals.

**Input:**
```json
{
  "period": "week | month"
}
```

**Output:**
```json
{
  "period": "week",
  "overall_score": 0.82,
  "habits": [
    {
      "name": "Running",
      "goal_minutes": 30,
      "average_actual": 33,
      "alignment": 1.0,
      "status": "exceeding"
    },
    {
      "name": "Reading",
      "goal_minutes": 60,
      "average_actual": 48,
      "alignment": 0.8,
      "status": "on track"
    }
  ]
}
```

---

### `schedule_pause`

Pause tracking for one or more habits during a period. Paused days are excluded from streaks and consistency scores. Omit `habit_ids` to pause all habits.

**Input:**
```json
{
  "from": "2026-04-01",
  "to": "2026-04-07",
  "reason": "vacation",
  "habit_ids": ["abc123", "def456"]
}
```

**Output:**
```json
{
  "success": true,
  "pause_id": "pause_xyz",
  "paused_days": 7,
  "habits_affected": ["Running", "Reading"],
  "message": "Tracking paused Apr 1–7. Your streaks are safe."
}
```

---

### `cancel_pause`

Cancel a scheduled or active pause early.

**Input:**
```json
{
  "pause_id": "pause_xyz",
  "resume_from": "2026-04-05"
}
```

**Output:**
```json
{
  "success": true,
  "message": "Pause cancelled. Tracking resumes from Apr 5."
}
```

---

### `list_pauses`

List pauses by status.

**Input:**
```json
{
  "status": "upcoming | active | past | all"
}
```

**Output:**
```json
{
  "pauses": [
    {
      "id": "pause_xyz",
      "from": "2026-04-01",
      "to": "2026-04-07",
      "reason": "vacation",
      "status": "upcoming",
      "habits_affected": ["Running", "Reading", "Typing"]
    }
  ]
}
```

---

## REST API

Used by the embedded web dashboard. All endpoints prefixed `/api`.

### Habits

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/habits` | List all active habits |
| `POST` | `/api/habits` | Create a new habit |
| `PUT` | `/api/habits/:id` | Update a habit |
| `DELETE` | `/api/habits/:id` | Archive a habit |
| `GET` | `/api/habits/:id/schedule` | Get habit schedule |
| `PUT` | `/api/habits/:id/schedule` | Update habit schedule |

### Logs

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/habits/:id/log` | Log a completion (build habit) |
| `POST` | `/api/habits/:id/slip` | Log a slip (avoid habit) |
| `GET` | `/api/habits/:id/logs` | Get log history |
| `DELETE` | `/api/logs/:id` | Delete a log entry |

### Pauses

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/pauses` | List pauses |
| `POST` | `/api/pauses` | Create a pause |
| `PATCH` | `/api/pauses/:id/cancel` | Cancel a pause early |

### Stats

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/stats/streaks` | Get all streaks + consistency scores |
| `GET` | `/api/stats/summary?period=week` | Get summary |
| `GET` | `/api/stats/heatmap?habit_id=x` | Get calendar heatmap data |
| `GET` | `/api/stats/goal-alignment?period=week` | Get goal alignment |

### System

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/health` | Health check |
| `GET` | `/api/version` | Version info |

---

## Streak & Consistency Calculation Rules

```
Scheduled days  = days habit was due (per schedule), minus paused days
Completed days  = days with a completion log on a scheduled day
Slipped days    = days with a full slip log (avoid habits only)

Streak (build):
  consecutive scheduled days with a completion
  resets on a missed scheduled day
  paused days do not break the streak
  non-scheduled days are ignored

Streak (avoid):
  consecutive scheduled days without a full slip
  minor slips do NOT break the streak
  full slips break the streak
  paused days do not break the streak

Consistency rate:
  completed_days / scheduled_days over the period
  paused days excluded from both numerator and denominator
  never collapses to zero from a single missed day
```

---

## Configuration

| Variable | Default | Description |
|---|---|---|
| `HABITCLAW_PORT` | `3000` | HTTP server port |
| `HABITCLAW_DB_PATH` | `./habitclaw.db` | SQLite database path |
| `HABITCLAW_PASSWORD` | `` | Optional dashboard password |
| `HABITCLAW_USER_ID` | `local` | Local user ID |
| `HABITCLAW_LOG_LEVEL` | `info` | Log level (debug/info/warn) |

---

## Installation

### One-liner (recommended)
```bash
curl -fsSL https://habitclaw.app/install.sh | bash
```

### Homebrew (macOS)
```bash
brew install habitclaw
```

### Go install
```bash
go install github.com/you/habitclaw@latest
```

### From source
```bash
git clone https://github.com/you/habitclaw.git
cd habitclaw
go build -o habitclaw .
./habitclaw
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

```markdown
- Check HabitClaw for incomplete build habits today
- Check HabitClaw for any logged slips on avoid habits
- If any habits not done and it's past 8pm, remind me
- If I have an active pause scheduled, acknowledge it
- If I am on a streak, mention it for motivation
```

Example cron for a scheduled habit:

```bash
# Remind about running at 06:00 on scheduled days
openclaw cron add \
  --name "Running reminder" \
  --cron "0 6 * * 1,2,4,5,6" \
  --session main \
  --message "Check HabitClaw — running is scheduled for today. Remind me." \
  --wake now
```

---

## Roadmap

- [ ] v0.1 — habits, schedules, SQLite, MCP server, basic dashboard
- [ ] v0.2 — avoid habits, slip logging, streak + consistency computation
- [ ] v0.3 — grace pauses, pause MCP tools
- [ ] v0.4 — heatmap calendar, goal alignment, stats UI
- [ ] v0.5 — habit templates, import/export (JSON/CSV)
- [ ] v1.0 — stable MCP spec, polished UI
- [ ] future — HabitClaw.ai cloud version (separate repo)

---

## License

AGPL-3.0 — free to self-host, modifications must be open source.
Cloud hosting rights reserved for HabitClaw.ai.