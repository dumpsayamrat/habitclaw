CREATE TABLE IF NOT EXISTS users (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS habits (
    id              TEXT PRIMARY KEY,
    user_id         TEXT NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT DEFAULT '',
    goal_type       TEXT NOT NULL CHECK (goal_type IN ('duration', 'count', 'boolean')),
    goal_value      INTEGER DEFAULT 0,
    habit_direction TEXT NOT NULL CHECK (habit_direction IN ('build', 'avoid')),
    color           TEXT DEFAULT '',
    icon            TEXT DEFAULT '',
    archived_at     DATETIME,
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS habit_schedules (
    id             TEXT PRIMARY KEY,
    habit_id       TEXT NOT NULL,
    user_id        TEXT NOT NULL,
    schedule_type  TEXT NOT NULL,
    days_of_week   TEXT DEFAULT '',
    times_per_week INTEGER DEFAULT 0,
    every_n_days   INTEGER DEFAULT 0,
    day_of_month   INTEGER DEFAULT 0,
    time_of_day    TEXT DEFAULT '',
    window_start   TEXT DEFAULT '',
    window_end     TEXT DEFAULT '',
    created_at     DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at     DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS completion_logs (
    id         TEXT PRIMARY KEY,
    habit_id   TEXT NOT NULL,
    user_id    TEXT NOT NULL,
    log_type   TEXT NOT NULL CHECK (log_type IN ('completion', 'slip')),
    date       DATE NOT NULL,
    value      INTEGER DEFAULT 0,
    note       TEXT DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS pauses (
    id           TEXT PRIMARY KEY,
    user_id      TEXT NOT NULL,
    habit_id     TEXT,
    from_date    DATE NOT NULL,
    to_date      DATE NOT NULL,
    reason       TEXT DEFAULT '',
    cancelled_at DATETIME,
    resume_from  DATE,
    created_at   DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_habits_user_id ON habits(user_id);
CREATE INDEX IF NOT EXISTS idx_habit_schedules_habit_id ON habit_schedules(habit_id);
CREATE INDEX IF NOT EXISTS idx_completion_logs_user_id ON completion_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_completion_logs_habit_date ON completion_logs(habit_id, date);
CREATE INDEX IF NOT EXISTS idx_pauses_user_id ON pauses(user_id);
