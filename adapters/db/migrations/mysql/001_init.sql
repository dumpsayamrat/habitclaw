CREATE TABLE IF NOT EXISTS users (
    id         VARCHAR(36) PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS habits (
    id              VARCHAR(36) PRIMARY KEY,
    user_id         VARCHAR(36) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT DEFAULT '',
    goal_type       VARCHAR(20) NOT NULL CHECK (goal_type IN ('duration', 'count', 'boolean')),
    goal_value      INTEGER DEFAULT 0,
    habit_direction VARCHAR(10) NOT NULL CHECK (habit_direction IN ('build', 'avoid')),
    color           VARCHAR(50) DEFAULT '',
    icon            VARCHAR(50) DEFAULT '',
    archived_at     DATETIME,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS habit_schedules (
    id             VARCHAR(36) PRIMARY KEY,
    habit_id       VARCHAR(36) NOT NULL,
    user_id        VARCHAR(36) NOT NULL,
    schedule_type  VARCHAR(30) NOT NULL,
    days_of_week   TEXT DEFAULT '',
    times_per_week INTEGER DEFAULT 0,
    every_n_days   INTEGER DEFAULT 0,
    day_of_month   INTEGER DEFAULT 0,
    time_of_day    VARCHAR(10) DEFAULT '',
    window_start   VARCHAR(10) DEFAULT '',
    window_end     VARCHAR(10) DEFAULT '',
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS completion_logs (
    id         VARCHAR(36) PRIMARY KEY,
    habit_id   VARCHAR(36) NOT NULL,
    user_id    VARCHAR(36) NOT NULL,
    log_type   VARCHAR(10) NOT NULL CHECK (log_type IN ('completion', 'slip')),
    date       DATE NOT NULL,
    value      INTEGER DEFAULT 0,
    note       TEXT DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pauses (
    id           VARCHAR(36) PRIMARY KEY,
    user_id      VARCHAR(36) NOT NULL,
    habit_id     VARCHAR(36),
    from_date    DATE NOT NULL,
    to_date      DATE NOT NULL,
    reason       TEXT DEFAULT '',
    cancelled_at DATETIME,
    resume_from  DATE,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_habits_user_id ON habits(user_id);
CREATE INDEX idx_habit_schedules_habit_id ON habit_schedules(habit_id);
CREATE INDEX idx_completion_logs_user_id ON completion_logs(user_id);
CREATE INDEX idx_completion_logs_habit_date ON completion_logs(habit_id, date);
CREATE INDEX idx_pauses_user_id ON pauses(user_id);
