# HabitClaw Prompt Tracker 🦞

> Update status by changing `[ ]` to `[x]` or `[~]` (in progress)

**Legend:** `[x]` done · `[~]` in progress · `[ ]` to do

---

## Progress

```
Phase 1 — Foundation        ██████████  2/2   ✅
Phase 2 — Core types        ██████████  9/9   ✅
Phase 3 — Database          ░░░░░░░░░░  0/3
Phase 4 — MCP server        ░░░░░░░░░░  0/5
Phase 5 — REST API          ░░░░░░░░░░  0/5
Phase 6 — Dashboard         ░░░░░░░░░░  0/6
Phase 7 — Polish            ░░░░░░░░░░  0/3
─────────────────────────────────────────────
Total                                  11/33
```

---

## Phase 1 — Foundation

| # | Status | Prompt | File(s) |
|---|--------|--------|---------|
| 01 | [x] | Project init + health API | `main.go`, `config.go` |
| 02 | [x] | Embedded web + health UI | `web/server.go`, `index.html` |

---

## Phase 2 — Core Types

| # | Status | Prompt | File(s) |
|---|--------|--------|---------|
| 03 | [x] | Habit types | `core/habit.go` |
| 04 | [x] | Log types | `core/log.go` |
| 05 | [x] | Pause type | `core/pause.go` |
| 06 | [x] | HabitStore interface | `core/interfaces.go` |
| 07 | [x] | Schedule resolution | `core/schedule.go` |
| 08 | [x] | Streak calculation | `core/streak.go` |
| 09 | [x] | Summary logic | `core/summary.go` |
| 10 | [x] | Goal alignment | `core/alignment.go` |
| 11 | [x] | HabitService | `core/service.go` |

---

## Phase 3 — Database

| # | Status | Prompt | File(s) |
|---|--------|--------|---------|
| 12 | [ ] | SQL migrations | `adapters/sqlite/migrations/001_init.sql` |
| 13 | [ ] | SQLite adapter | `adapters/sqlite/store.go` |
| 14 | [ ] | Single user auth | `adapters/auth/single_user.go` |

---

## Phase 4 — MCP Server

| # | Status | Prompt | File(s) |
|---|--------|--------|---------|
| 15 | [ ] | MCP server setup | `mcp/server.go` |
| 16 | [ ] | MCP: list + add + log habit | `mcp/handlers.go` |
| 17 | [ ] | MCP: slip + streaks + summary | `mcp/handlers.go` |
| 18 | [ ] | MCP: alignment + pauses | `mcp/handlers.go` |
| 19 | [ ] | Wire MCP into main | `main.go` |


## Phase 5 — Polish

| # | Status | Prompt | File(s) |
|---|--------|--------|---------|
| 31 | [ ] | Install script | `install.sh` |
| 32 | [ ] | Dockerfile | `Dockerfile` |
| 33 | [ ] | Integration tests | `*_integration_test.go` |

---