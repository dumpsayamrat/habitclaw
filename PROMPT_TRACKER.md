# HabitClaw Prompt Tracker 🦞

> Update status by changing `[ ]` to `[x]` or `[~]` (in progress)

**Legend:** `[x]` done · `[~]` in progress · `[ ]` to do

---

## Progress

```
Phase 1 — Foundation        ██████████  2/2   ✅
Phase 2 — Core types        ██████████  9/9   ✅
Phase 3 — Database          ██████████  3/3   ✅
Phase 4 — MCP server        ░░░░░░░░░░  0/5
Phase 5 — Polish            ░░░░░░░░░░  0/3
─────────────────────────────────────────────
Total                                  14/22
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
| 12 | [x] | SQL migrations + dialect layer | `adapters/db/dialect.go`, `adapters/db/open.go`, `adapters/db/migrate.go`, `adapters/db/migrations/` |
| 13 | [x] | Database adapter (HabitStore impl) | `adapters/db/store.go`, `adapters/db/store_test.go` |
| 14 | [x] | Single user auth + config + wiring | `adapters/auth/single_user.go`, `config/config.go`, `main.go` |

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