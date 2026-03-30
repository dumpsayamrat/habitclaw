package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/dumpsayamrat/habitclaw/core"
)

func setupTestDB(t *testing.T) (*sql.DB, Dialect, *Store) {
	t.Helper()
	database, dialect, err := OpenDatabase("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := Migrate(database, dialect); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	store := NewStore(database, dialect)

	// Seed local user
	_, err = database.Exec("INSERT INTO users (id, name, created_at) VALUES ('local', 'Local User', ?)", time.Now().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	t.Cleanup(func() { database.Close() })
	return database, dialect, store
}

func TestCreateAndListHabits(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{
		UserID:    "local",
		Name:      "Running",
		GoalType:  core.GoalTypeDuration,
		GoalValue: 30,
		Direction: core.DirectionBuild,
	}
	if err := store.CreateHabit(h); err != nil {
		t.Fatalf("create habit: %v", err)
	}

	habits, err := store.ListHabits("local")
	if err != nil {
		t.Fatalf("list habits: %v", err)
	}
	if len(habits) != 1 {
		t.Fatalf("expected 1 habit, got %d", len(habits))
	}
	if habits[0].Name != "Running" {
		t.Errorf("expected name Running, got %s", habits[0].Name)
	}
	if habits[0].GoalType != core.GoalTypeDuration {
		t.Errorf("expected goal_type duration, got %s", habits[0].GoalType)
	}
	if habits[0].Direction != core.DirectionBuild {
		t.Errorf("expected direction build, got %s", habits[0].Direction)
	}
	if habits[0].ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestUpdateHabit(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{
		ID:        "habit-update",
		UserID:    "local",
		Name:      "Running",
		GoalType:  core.GoalTypeDuration,
		GoalValue: 30,
		Direction: core.DirectionBuild,
	}
	store.CreateHabit(h)

	h.Name = "Jogging"
	h.GoalValue = 45
	if err := store.UpdateHabit(h); err != nil {
		t.Fatalf("update habit: %v", err)
	}

	habits, _ := store.ListHabits("local")
	if habits[0].Name != "Jogging" {
		t.Errorf("expected Jogging, got %s", habits[0].Name)
	}
}

func TestArchiveHabit(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{
		ID:        "habit-archive",
		UserID:    "local",
		Name:      "Running",
		GoalType:  core.GoalTypeDuration,
		GoalValue: 30,
		Direction: core.DirectionBuild,
	}
	store.CreateHabit(h)
	store.ArchiveHabit(h.ID)

	habits, _ := store.ListHabits("local")
	if len(habits) != 0 {
		t.Errorf("archived habit should not appear in list, got %d", len(habits))
	}
}

func TestSetAndGetSchedule(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{
		ID:        "habit-sched",
		UserID:    "local",
		Name:      "Running",
		GoalType:  core.GoalTypeDuration,
		GoalValue: 30,
		Direction: core.DirectionBuild,
	}
	store.CreateHabit(h)

	sched := core.HabitSchedule{
		HabitID:      h.ID,
		UserID:       "local",
		ScheduleType: core.ScheduleSpecificDays,
		DaysOfWeek:   []int{1, 3, 5},
		TimeOfDay:    "06:30",
	}
	if err := store.SetSchedule(sched); err != nil {
		t.Fatalf("set schedule: %v", err)
	}

	got, err := store.GetSchedule(h.ID)
	if err != nil {
		t.Fatalf("get schedule: %v", err)
	}
	if got.ScheduleType != core.ScheduleSpecificDays {
		t.Errorf("expected specific_days, got %s", got.ScheduleType)
	}
	if len(got.DaysOfWeek) != 3 || got.DaysOfWeek[0] != 1 || got.DaysOfWeek[1] != 3 || got.DaysOfWeek[2] != 5 {
		t.Errorf("expected days [1,3,5], got %v", got.DaysOfWeek)
	}
	if got.TimeOfDay != "06:30" {
		t.Errorf("expected time_of_day 06:30, got %s", got.TimeOfDay)
	}

	// Verify SetSchedule replaces (DELETE+INSERT)
	sched2 := core.HabitSchedule{
		HabitID:      h.ID,
		UserID:       "local",
		ScheduleType: core.ScheduleDaily,
	}
	store.SetSchedule(sched2)
	got2, _ := store.GetSchedule(h.ID)
	if got2.ScheduleType != core.ScheduleDaily {
		t.Errorf("expected daily after replace, got %s", got2.ScheduleType)
	}
}

func TestSetScheduleJoinedInListHabits(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{
		ID:        "habit-sched-join",
		UserID:    "local",
		Name:      "Reading",
		GoalType:  core.GoalTypeDuration,
		GoalValue: 60,
		Direction: core.DirectionBuild,
	}
	store.CreateHabit(h)
	store.SetSchedule(core.HabitSchedule{
		HabitID:      h.ID,
		UserID:       "local",
		ScheduleType: core.ScheduleWeekdays,
	})

	habits, _ := store.ListHabits("local")
	if habits[0].Schedule == nil {
		t.Fatal("expected schedule to be joined in ListHabits")
	}
	if habits[0].Schedule.ScheduleType != core.ScheduleWeekdays {
		t.Errorf("expected weekdays, got %s", habits[0].Schedule.ScheduleType)
	}
}

func TestLogCompletionAndSlip(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{
		ID:        "habit-log",
		UserID:    "local",
		Name:      "Running",
		GoalType:  core.GoalTypeDuration,
		GoalValue: 30,
		Direction: core.DirectionBuild,
	}
	store.CreateHabit(h)

	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Log completion
	if err := store.LogCompletion(core.CompletionLog{
		HabitID: h.ID,
		UserID:  "local",
		Date:    today,
		Value:   35,
		Note:    "felt great",
	}); err != nil {
		t.Fatalf("log completion: %v", err)
	}

	// Log slip
	if err := store.LogSlip(core.CompletionLog{
		HabitID: h.ID,
		UserID:  "local",
		Date:    today.AddDate(0, 0, -1),
		Value:   int(core.SlipMinor),
		Note:    "minor slip",
	}); err != nil {
		t.Fatalf("log slip: %v", err)
	}

	logs, err := store.GetLogs("local", today.AddDate(0, 0, -7), today)
	if err != nil {
		t.Fatalf("get logs: %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}

	// Verify completion
	var completion, slip *core.CompletionLog
	for i := range logs {
		if logs[i].LogType == core.LogCompletion {
			completion = &logs[i]
		} else {
			slip = &logs[i]
		}
	}

	if completion == nil {
		t.Fatal("expected a completion log")
	}
	if completion.Value != 35 {
		t.Errorf("expected completion value 35, got %d", completion.Value)
	}

	if slip == nil {
		t.Fatal("expected a slip log")
	}
	if slip.LogType != core.LogSlip {
		t.Errorf("expected slip log type, got %s", slip.LogType)
	}
}

func TestDeleteLog(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{ID: "habit-dellog", UserID: "local", Name: "Test", GoalType: core.GoalTypeBoolean, Direction: core.DirectionBuild}
	store.CreateHabit(h)

	today := time.Now().UTC().Truncate(24 * time.Hour)
	log := core.CompletionLog{ID: "log-del", HabitID: h.ID, UserID: "local", Date: today, Value: 1}
	store.LogCompletion(log)

	if err := store.DeleteLog("log-del"); err != nil {
		t.Fatalf("delete log: %v", err)
	}

	logs, _ := store.GetLogs("local", today.AddDate(0, 0, -1), today)
	if len(logs) != 0 {
		t.Errorf("expected 0 logs after delete, got %d", len(logs))
	}
}

func TestCreateAndListPauses(t *testing.T) {
	_, _, store := setupTestDB(t)

	from := time.Now().UTC().Truncate(24 * time.Hour)
	to := from.AddDate(0, 0, 7)

	p := core.Pause{
		UserID:   "local",
		FromDate: from,
		ToDate:   to,
		Reason:   "vacation",
	}
	if err := store.CreatePause(p); err != nil {
		t.Fatalf("create pause: %v", err)
	}

	pauses, err := store.ListPauses("local", "all")
	if err != nil {
		t.Fatalf("list pauses: %v", err)
	}
	if len(pauses) != 1 {
		t.Fatalf("expected 1 pause, got %d", len(pauses))
	}
	if pauses[0].Reason != "vacation" {
		t.Errorf("expected reason vacation, got %s", pauses[0].Reason)
	}
}

func TestCancelPause(t *testing.T) {
	_, _, store := setupTestDB(t)

	from := time.Now().UTC().Truncate(24 * time.Hour)
	to := from.AddDate(0, 0, 7)
	resumeFrom := from.AddDate(0, 0, 3)

	p := core.Pause{ID: "pause-cancel", UserID: "local", FromDate: from, ToDate: to, Reason: "vacation"}
	store.CreatePause(p)

	if err := store.CancelPause("pause-cancel", resumeFrom); err != nil {
		t.Fatalf("cancel pause: %v", err)
	}

	pauses, _ := store.ListPauses("local", "all")
	if pauses[0].CancelledAt == nil {
		t.Error("expected cancelled_at to be set")
	}
	if pauses[0].ResumeFrom == nil {
		t.Error("expected resume_from to be set")
	}
}

func TestGetStreaks(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{
		ID:        "habit-streak",
		UserID:    "local",
		Name:      "Running",
		GoalType:  core.GoalTypeDuration,
		GoalValue: 30,
		Direction: core.DirectionBuild,
	}
	store.CreateHabit(h)
	store.SetSchedule(core.HabitSchedule{
		HabitID:      h.ID,
		UserID:       "local",
		ScheduleType: core.ScheduleDaily,
	})

	today := time.Now().UTC().Truncate(24 * time.Hour)
	for i := 0; i < 3; i++ {
		store.LogCompletion(core.CompletionLog{
			HabitID: h.ID,
			UserID:  "local",
			Date:    today.AddDate(0, 0, -i),
			Value:   35,
		})
	}

	streaks, err := store.GetStreaks("local")
	if err != nil {
		t.Fatalf("get streaks: %v", err)
	}
	if len(streaks) != 1 {
		t.Fatalf("expected 1 streak, got %d", len(streaks))
	}
	if streaks[0].Current < 3 {
		t.Errorf("expected current streak >= 3, got %d", streaks[0].Current)
	}
	if streaks[0].HabitName != "Running" {
		t.Errorf("expected habit name Running, got %s", streaks[0].HabitName)
	}
}

func TestGetSummary(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{
		ID:        "habit-summary",
		UserID:    "local",
		Name:      "Running",
		GoalType:  core.GoalTypeDuration,
		GoalValue: 30,
		Direction: core.DirectionBuild,
	}
	store.CreateHabit(h)
	store.SetSchedule(core.HabitSchedule{
		HabitID:      h.ID,
		UserID:       "local",
		ScheduleType: core.ScheduleDaily,
	})

	today := time.Now().UTC().Truncate(24 * time.Hour)
	store.LogCompletion(core.CompletionLog{
		HabitID: h.ID,
		UserID:  "local",
		Date:    today,
		Value:   35,
	})

	summary, err := store.GetSummary("local", core.PeriodWeek)
	if err != nil {
		t.Fatalf("get summary: %v", err)
	}
	if summary.TotalHabits != 1 {
		t.Errorf("expected 1 total habit, got %d", summary.TotalHabits)
	}
	if summary.Period != "week" {
		t.Errorf("expected period week, got %s", summary.Period)
	}
}

func TestGetGoalAlignment(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{
		ID:        "habit-align",
		UserID:    "local",
		Name:      "Running",
		GoalType:  core.GoalTypeDuration,
		GoalValue: 30,
		Direction: core.DirectionBuild,
	}
	store.CreateHabit(h)
	store.SetSchedule(core.HabitSchedule{
		HabitID:      h.ID,
		UserID:       "local",
		ScheduleType: core.ScheduleDaily,
	})

	today := time.Now().UTC().Truncate(24 * time.Hour)
	store.LogCompletion(core.CompletionLog{
		HabitID: h.ID,
		UserID:  "local",
		Date:    today,
		Value:   35,
	})

	alignment, err := store.GetGoalAlignment("local", core.PeriodWeek)
	if err != nil {
		t.Fatalf("get goal alignment: %v", err)
	}
	if alignment.Period != "week" {
		t.Errorf("expected period week, got %s", alignment.Period)
	}
}

func TestGetHeatmap(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{
		ID:        "habit-heatmap",
		UserID:    "local",
		Name:      "Running",
		GoalType:  core.GoalTypeDuration,
		GoalValue: 30,
		Direction: core.DirectionBuild,
	}
	store.CreateHabit(h)

	today := time.Now().UTC().Truncate(24 * time.Hour)
	store.LogCompletion(core.CompletionLog{
		HabitID: h.ID,
		UserID:  "local",
		Date:    today,
		Value:   35,
	})

	days, err := store.GetHeatmap(h.ID, today, today)
	if err != nil {
		t.Fatalf("get heatmap: %v", err)
	}
	if len(days) != 1 {
		t.Fatalf("expected 1 day, got %d", len(days))
	}
	if days[0].Value != 35 {
		t.Errorf("expected value 35, got %d", days[0].Value)
	}
	if !days[0].GoalMet {
		t.Error("expected GoalMet=true for 35 >= 30")
	}
}

func TestDialectRebind(t *testing.T) {
	sqlite := SQLiteDialect{}
	if got := sqlite.Rebind("SELECT * FROM t WHERE id = ?"); got != "SELECT * FROM t WHERE id = ?" {
		t.Errorf("sqlite rebind should be identity, got %s", got)
	}

	pg := PostgresDialect{}
	got := pg.Rebind("SELECT * FROM t WHERE id = ? AND name = ?")
	expected := "SELECT * FROM t WHERE id = $1 AND name = $2"
	if got != expected {
		t.Errorf("postgres rebind: expected %s, got %s", expected, got)
	}

	mysql := MySQLDialect{}
	if got := mysql.Rebind("SELECT * FROM t WHERE id = ?"); got != "SELECT * FROM t WHERE id = ?" {
		t.Errorf("mysql rebind should be identity, got %s", got)
	}
}

func TestMigrateIdempotent(t *testing.T) {
	database, dialect, _ := setupTestDB(t)

	// Run migrate again — should be a no-op
	if err := Migrate(database, dialect); err != nil {
		t.Fatalf("second migrate should succeed: %v", err)
	}

	// Verify tables still exist
	var count int
	err := database.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 migration record, got %d", count)
	}
}

func TestForeignKeysEnabled(t *testing.T) {
	database, _, _ := setupTestDB(t)

	var fkEnabled int
	err := database.QueryRow("PRAGMA foreign_keys").Scan(&fkEnabled)
	if err != nil {
		t.Fatalf("query foreign_keys pragma: %v", err)
	}
	if fkEnabled != 1 {
		t.Errorf("expected foreign_keys=1, got %d", fkEnabled)
	}
}

func TestNewIDVersionVariant(t *testing.T) {
	id := newID()
	if len(id) != 36 {
		t.Fatalf("expected 36-char UUID, got %d: %s", len(id), id)
	}

	// Version nibble (position 14) should be '4'
	if id[14] != '4' {
		t.Errorf("expected version nibble '4' at position 14, got '%c'", id[14])
	}

	// Variant nibble (position 19) should be 8, 9, a, or b
	v := id[19]
	if v != '8' && v != '9' && v != 'a' && v != 'b' {
		t.Errorf("expected variant nibble in [8,9,a,b] at position 19, got '%c'", v)
	}
}

func TestPauseWithHabitID(t *testing.T) {
	_, _, store := setupTestDB(t)

	h := core.Habit{ID: "habit-pause-specific", UserID: "local", Name: "Reading", GoalType: core.GoalTypeBoolean, Direction: core.DirectionBuild}
	store.CreateHabit(h)

	habitID := h.ID
	from := time.Now().UTC().Truncate(24 * time.Hour)
	to := from.AddDate(0, 0, 3)

	store.CreatePause(core.Pause{
		UserID:   "local",
		HabitID:  &habitID,
		FromDate: from,
		ToDate:   to,
		Reason:   "sick",
	})

	pauses, _ := store.ListPauses("local", "all")
	if len(pauses) != 1 {
		t.Fatalf("expected 1 pause, got %d", len(pauses))
	}
	if pauses[0].HabitID == nil || *pauses[0].HabitID != habitID {
		t.Error("expected habit_id to round-trip")
	}
}
