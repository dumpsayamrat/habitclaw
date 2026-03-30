package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dumpsayamrat/habitclaw/core"
)

var _ core.HabitStore = (*Store)(nil)

type Store struct {
	db      *sql.DB
	dialect Dialect
}

func NewStore(db *sql.DB, dialect Dialect) *Store {
	return &Store{db: db, dialect: dialect}
}

func newID() string {
	var uuid [16]byte
	rand.Read(uuid[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

func dateStr(t time.Time) string {
	return t.Format("2006-01-02")
}

// parseDate handles date strings from SQLite which may return
// "2006-01-02" or "2006-01-02T15:04:05Z" depending on the driver.
func parseDate(s string) time.Time {
	if len(s) > 10 {
		s = s[:10]
	}
	t, _ := time.Parse("2006-01-02", s)
	return t
}

// parseTimestamp handles timestamp strings that may be in RFC3339 or
// SQLite datetime format ("2006-01-02 15:04:05").
func parseTimestamp(s string) time.Time {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}

func (s *Store) CreateHabit(habit core.Habit) error {
	if habit.ID == "" {
		habit.ID = newID()
	}
	now := time.Now().UTC()
	if habit.CreatedAt.IsZero() {
		habit.CreatedAt = now
	}
	if habit.UpdatedAt.IsZero() {
		habit.UpdatedAt = now
	}

	q := s.dialect.Rebind(`INSERT INTO habits (id, user_id, name, description, goal_type, goal_value, habit_direction, color, icon, archived_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	var archivedAt *string
	if habit.ArchivedAt != nil {
		s := habit.ArchivedAt.UTC().Format(time.RFC3339)
		archivedAt = &s
	}

	_, err := s.db.Exec(q, habit.ID, habit.UserID, habit.Name, habit.Description,
		string(habit.GoalType), habit.GoalValue, string(habit.Direction),
		habit.Color, habit.Icon, archivedAt,
		habit.CreatedAt.UTC().Format(time.RFC3339), habit.UpdatedAt.UTC().Format(time.RFC3339))
	return err
}

func (s *Store) ListHabits(userID string) ([]core.Habit, error) {
	q := s.dialect.Rebind(`SELECT h.id, h.user_id, h.name, h.description, h.goal_type, h.goal_value, h.habit_direction, h.color, h.icon, h.archived_at, h.created_at, h.updated_at,
		hs.id, hs.schedule_type, hs.days_of_week, hs.times_per_week, hs.every_n_days, hs.day_of_month, hs.time_of_day, hs.window_start, hs.window_end
		FROM habits h
		LEFT JOIN habit_schedules hs ON hs.habit_id = h.id
		WHERE h.user_id = ? AND h.archived_at IS NULL
		ORDER BY h.created_at ASC`)

	rows, err := s.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []core.Habit
	for rows.Next() {
		var h core.Habit
		var goalType, direction string
		var archivedAt, createdAt, updatedAt sql.NullString

		var schedID, schedType, schedDaysJSON, schedTimeOfDay, schedWindowStart, schedWindowEnd sql.NullString
		var schedTimesPerWeek, schedEveryNDays, schedDayOfMonth sql.NullInt64

		err := rows.Scan(
			&h.ID, &h.UserID, &h.Name, &h.Description, &goalType, &h.GoalValue, &direction,
			&h.Color, &h.Icon, &archivedAt, &createdAt, &updatedAt,
			&schedID, &schedType, &schedDaysJSON, &schedTimesPerWeek, &schedEveryNDays, &schedDayOfMonth,
			&schedTimeOfDay, &schedWindowStart, &schedWindowEnd,
		)
		if err != nil {
			return nil, err
		}

		h.GoalType = core.GoalType(goalType)
		h.Direction = core.HabitDirection(direction)
		if createdAt.Valid {
			h.CreatedAt = parseTimestamp(createdAt.String)
		}
		if updatedAt.Valid {
			h.UpdatedAt = parseTimestamp(updatedAt.String)
		}
		if archivedAt.Valid {
			t := parseTimestamp(archivedAt.String)
			h.ArchivedAt = &t
		}

		if schedID.Valid {
			sched := &core.HabitSchedule{
				ID:           schedID.String,
				HabitID:      h.ID,
				UserID:       h.UserID,
				ScheduleType: core.ScheduleType(schedType.String),
				TimesPerWeek: int(schedTimesPerWeek.Int64),
				EveryNDays:   int(schedEveryNDays.Int64),
				DayOfMonth:   int(schedDayOfMonth.Int64),
				TimeOfDay:    schedTimeOfDay.String,
				WindowStart:  schedWindowStart.String,
				WindowEnd:    schedWindowEnd.String,
			}
			if schedDaysJSON.Valid && schedDaysJSON.String != "" {
				json.Unmarshal([]byte(schedDaysJSON.String), &sched.DaysOfWeek)
			}
			h.Schedule = sched
		}

		habits = append(habits, h)
	}
	return habits, rows.Err()
}

func (s *Store) UpdateHabit(habit core.Habit) error {
	habit.UpdatedAt = time.Now().UTC()
	q := s.dialect.Rebind(`UPDATE habits SET name = ?, description = ?, goal_type = ?, goal_value = ?, habit_direction = ?, color = ?, icon = ?, updated_at = ?
		WHERE id = ?`)
	_, err := s.db.Exec(q, habit.Name, habit.Description, string(habit.GoalType), habit.GoalValue,
		string(habit.Direction), habit.Color, habit.Icon,
		habit.UpdatedAt.Format(time.RFC3339), habit.ID)
	return err
}

func (s *Store) ArchiveHabit(id string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	q := s.dialect.Rebind(`UPDATE habits SET archived_at = ?, updated_at = ? WHERE id = ?`)
	_, err := s.db.Exec(q, now, now, id)
	return err
}

func (s *Store) SetSchedule(schedule core.HabitSchedule) error {
	if schedule.ID == "" {
		schedule.ID = newID()
	}
	now := time.Now().UTC()

	daysJSON := ""
	if len(schedule.DaysOfWeek) > 0 {
		b, _ := json.Marshal(schedule.DaysOfWeek)
		daysJSON = string(b)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	delQ := s.dialect.Rebind(`DELETE FROM habit_schedules WHERE habit_id = ?`)
	if _, err := tx.Exec(delQ, schedule.HabitID); err != nil {
		tx.Rollback()
		return err
	}

	insQ := s.dialect.Rebind(`INSERT INTO habit_schedules (id, habit_id, user_id, schedule_type, days_of_week, times_per_week, every_n_days, day_of_month, time_of_day, window_start, window_end, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if _, err := tx.Exec(insQ, schedule.ID, schedule.HabitID, schedule.UserID,
		string(schedule.ScheduleType), daysJSON, schedule.TimesPerWeek, schedule.EveryNDays,
		schedule.DayOfMonth, schedule.TimeOfDay, schedule.WindowStart, schedule.WindowEnd,
		now.Format(time.RFC3339), now.Format(time.RFC3339)); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Store) GetSchedule(habitID string) (core.HabitSchedule, error) {
	q := s.dialect.Rebind(`SELECT id, habit_id, user_id, schedule_type, days_of_week, times_per_week, every_n_days, day_of_month, time_of_day, window_start, window_end, created_at, updated_at
		FROM habit_schedules WHERE habit_id = ? LIMIT 1`)

	var sched core.HabitSchedule
	var schedType, daysJSON string
	var createdAt, updatedAt string

	err := s.db.QueryRow(q, habitID).Scan(
		&sched.ID, &sched.HabitID, &sched.UserID, &schedType, &daysJSON,
		&sched.TimesPerWeek, &sched.EveryNDays, &sched.DayOfMonth,
		&sched.TimeOfDay, &sched.WindowStart, &sched.WindowEnd,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return sched, err
	}

	sched.ScheduleType = core.ScheduleType(schedType)
	if daysJSON != "" {
		json.Unmarshal([]byte(daysJSON), &sched.DaysOfWeek)
	}
	sched.CreatedAt = parseTimestamp(createdAt)
	sched.UpdatedAt = parseTimestamp(updatedAt)

	return sched, nil
}

func (s *Store) LogCompletion(log core.CompletionLog) error {
	return s.insertLog(log, core.LogCompletion)
}

func (s *Store) LogSlip(log core.CompletionLog) error {
	return s.insertLog(log, core.LogSlip)
}

func (s *Store) insertLog(log core.CompletionLog, logType core.LogType) error {
	if log.ID == "" {
		log.ID = newID()
	}
	now := time.Now().UTC()
	if log.CreatedAt.IsZero() {
		log.CreatedAt = now
	}

	q := s.dialect.Rebind(`INSERT INTO completion_logs (id, habit_id, user_id, log_type, date, value, note, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := s.db.Exec(q, log.ID, log.HabitID, log.UserID, string(logType),
		dateStr(log.Date), log.Value, log.Note, log.CreatedAt.Format(time.RFC3339))
	return err
}

func (s *Store) GetLogs(userID string, from, to time.Time) ([]core.CompletionLog, error) {
	q := s.dialect.Rebind(`SELECT id, habit_id, user_id, log_type, date, value, note, created_at
		FROM completion_logs WHERE user_id = ? AND date >= ? AND date <= ?
		ORDER BY date ASC, created_at ASC`)

	rows, err := s.db.Query(q, userID, dateStr(from), dateStr(to))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []core.CompletionLog
	for rows.Next() {
		var l core.CompletionLog
		var logType, dateS, createdAt string
		err := rows.Scan(&l.ID, &l.HabitID, &l.UserID, &logType, &dateS, &l.Value, &l.Note, &createdAt)
		if err != nil {
			return nil, err
		}
		l.LogType = core.LogType(logType)
		l.Date = parseDate(dateS)
		l.CreatedAt = parseTimestamp(createdAt)
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func (s *Store) DeleteLog(id string) error {
	q := s.dialect.Rebind(`DELETE FROM completion_logs WHERE id = ?`)
	_, err := s.db.Exec(q, id)
	return err
}

func (s *Store) CreatePause(pause core.Pause) error {
	if pause.ID == "" {
		pause.ID = newID()
	}
	now := time.Now().UTC()
	if pause.CreatedAt.IsZero() {
		pause.CreatedAt = now
	}

	q := s.dialect.Rebind(`INSERT INTO pauses (id, user_id, habit_id, from_date, to_date, reason, cancelled_at, resume_from, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	habitID := pause.HabitID
	var cancelledAt, resumeFrom *string
	if pause.CancelledAt != nil {
		s := pause.CancelledAt.UTC().Format(time.RFC3339)
		cancelledAt = &s
	}
	if pause.ResumeFrom != nil {
		s := dateStr(*pause.ResumeFrom)
		resumeFrom = &s
	}

	_, err := s.db.Exec(q, pause.ID, pause.UserID, habitID,
		dateStr(pause.FromDate), dateStr(pause.ToDate), pause.Reason,
		cancelledAt, resumeFrom, pause.CreatedAt.Format(time.RFC3339))
	return err
}

func (s *Store) ListPauses(userID string, status string) ([]core.Pause, error) {
	today := dateStr(time.Now().UTC())

	var q string
	var args []interface{}

	switch status {
	case "active":
		q = `SELECT id, user_id, habit_id, from_date, to_date, reason, cancelled_at, resume_from, created_at
			FROM pauses WHERE user_id = ? AND from_date <= ? AND to_date >= ? AND cancelled_at IS NULL
			ORDER BY from_date ASC`
		args = []interface{}{userID, today, today}
	case "upcoming":
		q = `SELECT id, user_id, habit_id, from_date, to_date, reason, cancelled_at, resume_from, created_at
			FROM pauses WHERE user_id = ? AND from_date > ? AND cancelled_at IS NULL
			ORDER BY from_date ASC`
		args = []interface{}{userID, today}
	case "past":
		q = `SELECT id, user_id, habit_id, from_date, to_date, reason, cancelled_at, resume_from, created_at
			FROM pauses WHERE user_id = ? AND (to_date < ? OR cancelled_at IS NOT NULL)
			ORDER BY from_date DESC`
		args = []interface{}{userID, today}
	default: // "all"
		q = `SELECT id, user_id, habit_id, from_date, to_date, reason, cancelled_at, resume_from, created_at
			FROM pauses WHERE user_id = ?
			ORDER BY from_date DESC`
		args = []interface{}{userID}
	}

	q = s.dialect.Rebind(q)
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pauses []core.Pause
	for rows.Next() {
		var p core.Pause
		var habitID, cancelledAt, resumeFrom sql.NullString
		var fromDate, toDate, createdAt string

		err := rows.Scan(&p.ID, &p.UserID, &habitID, &fromDate, &toDate, &p.Reason, &cancelledAt, &resumeFrom, &createdAt)
		if err != nil {
			return nil, err
		}

		if habitID.Valid {
			p.HabitID = &habitID.String
		}
		p.FromDate = parseDate(fromDate)
		p.ToDate = parseDate(toDate)
		p.CreatedAt = parseTimestamp(createdAt)
		if cancelledAt.Valid {
			t := parseTimestamp(cancelledAt.String)
			p.CancelledAt = &t
		}
		if resumeFrom.Valid {
			t := parseDate(resumeFrom.String)
			p.ResumeFrom = &t
		}
		pauses = append(pauses, p)
	}
	return pauses, rows.Err()
}

func (s *Store) CancelPause(id string, resumeFrom time.Time) error {
	now := time.Now().UTC().Format(time.RFC3339)
	q := s.dialect.Rebind(`UPDATE pauses SET cancelled_at = ?, resume_from = ? WHERE id = ?`)
	_, err := s.db.Exec(q, now, dateStr(resumeFrom), id)
	return err
}

func periodDays(period core.Period) int {
	switch period {
	case core.PeriodToday:
		return 1
	case core.PeriodWeek:
		return 7
	default:
		return 30
	}
}

// Computed methods — delegate to core pure functions

func (s *Store) GetStreaks(userID string) ([]core.Streak, error) {
	habits, err := s.ListHabits(userID)
	if err != nil {
		return nil, err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	from := today.AddDate(-1, 0, 0)

	logs, err := s.GetLogs(userID, from, today)
	if err != nil {
		return nil, err
	}

	pauses, err := s.ListPauses(userID, "all")
	if err != nil {
		return nil, err
	}

	streaks := make([]core.Streak, 0, len(habits))
	for _, h := range habits {
		streaks = append(streaks, core.CalculateStreak(h, logs, pauses, today))
	}
	return streaks, nil
}

func (s *Store) GetSummary(userID string, period core.Period) (core.Summary, error) {
	habits, err := s.ListHabits(userID)
	if err != nil {
		return core.Summary{}, err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	from := today.AddDate(0, 0, -periodDays(period))

	logs, err := s.GetLogs(userID, from, today)
	if err != nil {
		return core.Summary{}, err
	}

	pauses, err := s.ListPauses(userID, "all")
	if err != nil {
		return core.Summary{}, err
	}

	return core.CalculateSummary(habits, logs, pauses, period, today), nil
}

func (s *Store) GetGoalAlignment(userID string, period core.Period) (core.GoalAlignment, error) {
	habits, err := s.ListHabits(userID)
	if err != nil {
		return core.GoalAlignment{}, err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	from := today.AddDate(0, 0, -periodDays(period))

	logs, err := s.GetLogs(userID, from, today)
	if err != nil {
		return core.GoalAlignment{}, err
	}

	pauses, err := s.ListPauses(userID, "all")
	if err != nil {
		return core.GoalAlignment{}, err
	}

	return core.CalculateGoalAlignment(habits, logs, pauses, period, today), nil
}

func (s *Store) GetHeatmap(habitID string, from, to time.Time) ([]core.HeatmapDay, error) {
	q := s.dialect.Rebind(`SELECT id, user_id, name, description, goal_type, goal_value, habit_direction, color, icon, archived_at, created_at, updated_at
		FROM habits WHERE id = ?`)

	var h core.Habit
	var goalType, direction string
	var archivedAt, createdAt, updatedAt sql.NullString

	err := s.db.QueryRow(q, habitID).Scan(
		&h.ID, &h.UserID, &h.Name, &h.Description, &goalType, &h.GoalValue, &direction,
		&h.Color, &h.Icon, &archivedAt, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	h.GoalType = core.GoalType(goalType)
	h.Direction = core.HabitDirection(direction)
	if createdAt.Valid {
		h.CreatedAt = parseTimestamp(createdAt.String)
	}
	if updatedAt.Valid {
		h.UpdatedAt = parseTimestamp(updatedAt.String)
	}

	logsQ := s.dialect.Rebind(`SELECT id, habit_id, user_id, log_type, date, value, note, created_at
		FROM completion_logs WHERE habit_id = ? AND date >= ? AND date <= ?
		ORDER BY date ASC`)

	rows, err := s.db.Query(logsQ, habitID, dateStr(from), dateStr(to))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []core.CompletionLog
	for rows.Next() {
		var l core.CompletionLog
		var logType, dateS, ca string
		if err := rows.Scan(&l.ID, &l.HabitID, &l.UserID, &logType, &dateS, &l.Value, &l.Note, &ca); err != nil {
			return nil, err
		}
		l.LogType = core.LogType(logType)
		l.Date = parseDate(dateS)
		l.CreatedAt = parseTimestamp(ca)
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	pauses, err := s.ListPauses(h.UserID, "all")
	if err != nil {
		return nil, err
	}

	return core.CalculateHeatmap(h, logs, pauses, from, to), nil
}
