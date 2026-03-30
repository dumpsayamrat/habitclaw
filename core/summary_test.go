package core

import (
	"math"
	"testing"
)

func TestPeriodConstants(t *testing.T) {
	tests := []struct {
		name string
		got  Period
		want string
	}{
		{"PeriodToday", PeriodToday, "today"},
		{"PeriodWeek", PeriodWeek, "week"},
		{"PeriodMonth", PeriodMonth, "month"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestSummaryInitialization(t *testing.T) {
	summary := Summary{
		Period:         "week",
		TotalHabits:    3,
		ScheduledDays:  21,
		CompletionRate: 0.85,
		HabitSummaries: []HabitStat{
			{
				HabitID:       "h1",
				Name:          "Running",
				Direction:     "build",
				ScheduledDays: 5,
				CompletedDays: 4,
				PausedDays:    0,
				GoalMetRate:   0.9,
				TotalValue:    155,
				CurrentStreak: 8,
			},
			{
				HabitID:       "h2",
				Name:          "No social media",
				Direction:     "avoid",
				ScheduledDays: 7,
				SlippedDays:   1,
				PausedDays:    0,
				CurrentStreak: 12,
			},
		},
	}

	if summary.Period != "week" {
		t.Errorf("Period = %q, want %q", summary.Period, "week")
	}
	if summary.TotalHabits != 3 {
		t.Errorf("TotalHabits = %d, want %d", summary.TotalHabits, 3)
	}
	if len(summary.HabitSummaries) != 2 {
		t.Errorf("HabitSummaries length = %d, want %d", len(summary.HabitSummaries), 2)
	}
	if summary.HabitSummaries[0].Name != "Running" {
		t.Errorf("first habit name = %q, want %q", summary.HabitSummaries[0].Name, "Running")
	}
	if summary.HabitSummaries[1].SlippedDays != 1 {
		t.Errorf("second habit slipped days = %d, want %d", summary.HabitSummaries[1].SlippedDays, 1)
	}
}

func approxEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func dailySchedule() *HabitSchedule {
	return &HabitSchedule{
		ID:           "sched-1",
		ScheduleType: ScheduleDaily,
	}
}

func TestCalculateSummary_BuildHabit(t *testing.T) {
	// today = 2026-03-30 (Monday)
	today := makeDate(2026, 3, 30)

	// Build habit: goal value = 30, daily schedule
	h := newTestHabit("Running", DirectionBuild, GoalTypeCount, dailySchedule())
	h.GoalValue = 30
	h.CreatedAt = makeDate(2026, 1, 1)

	// 7 days: Mar 24..Mar 30
	// Log completions on 5 days; 3 of them meet goal (value >= 30)
	logs := []CompletionLog{
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 24), 35), // goal met
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 25), 20), // not met
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 26), 30), // goal met
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 28), 10), // not met
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 30), 40), // goal met
	}

	summary := CalculateSummary([]Habit{h}, logs, nil, PeriodWeek, today)

	if summary.TotalHabits != 1 {
		t.Errorf("TotalHabits = %d, want 1", summary.TotalHabits)
	}
	if summary.Period != "week" {
		t.Errorf("Period = %q, want %q", summary.Period, "week")
	}

	stat := summary.HabitSummaries[0]

	if stat.ScheduledDays != 7 {
		t.Errorf("ScheduledDays = %d, want 7", stat.ScheduledDays)
	}
	if stat.CompletedDays != 5 {
		t.Errorf("CompletedDays = %d, want 5", stat.CompletedDays)
	}
	if stat.TotalValue != 135 {
		t.Errorf("TotalValue = %d, want 135", stat.TotalValue)
	}
	// GoalMetRate = 3 goal-met / 7 scheduled ≈ 0.4286
	want := 3.0 / 7.0
	if !approxEqual(stat.GoalMetRate, want, 0.001) {
		t.Errorf("GoalMetRate = %.4f, want %.4f", stat.GoalMetRate, want)
	}
	// CompletionRate = 5 / 7
	wantRate := 5.0 / 7.0
	if !approxEqual(summary.CompletionRate, wantRate, 0.001) {
		t.Errorf("CompletionRate = %.4f, want %.4f", summary.CompletionRate, wantRate)
	}
}

func TestCalculateSummary_AvoidHabit(t *testing.T) {
	today := makeDate(2026, 3, 30)

	h := newTestHabit("No social media", DirectionAvoid, GoalTypeBoolean, dailySchedule())
	h.CreatedAt = makeDate(2026, 1, 1)

	// 7 days: Mar 24..Mar 30
	// 1 full slip on Mar 26
	logs := []CompletionLog{
		newTestLog(h.ID, LogSlip, makeDate(2026, 3, 26), int(SlipFull)),
	}

	summary := CalculateSummary([]Habit{h}, logs, nil, PeriodWeek, today)
	stat := summary.HabitSummaries[0]

	if stat.ScheduledDays != 7 {
		t.Errorf("ScheduledDays = %d, want 7", stat.ScheduledDays)
	}
	if stat.SlippedDays != 1 {
		t.Errorf("SlippedDays = %d, want 1", stat.SlippedDays)
	}
	if stat.CompletedDays != 6 {
		t.Errorf("CompletedDays = %d, want 6", stat.CompletedDays)
	}
	// GoalMetRate = 6/7
	want := 6.0 / 7.0
	if !approxEqual(stat.GoalMetRate, want, 0.001) {
		t.Errorf("GoalMetRate = %.4f, want %.4f", stat.GoalMetRate, want)
	}
}

func TestCalculateSummary_PausedDaysExcluded(t *testing.T) {
	today := makeDate(2026, 3, 30)

	h := newTestHabit("Exercise", DirectionBuild, GoalTypeBoolean, dailySchedule())
	h.CreatedAt = makeDate(2026, 1, 1)

	// Pause Mar 28..Mar 29 (2 days)
	pause := newTestPause(&h.ID, makeDate(2026, 3, 28), makeDate(2026, 3, 29))

	// Complete on all non-paused days (Mar 24, 25, 26, 27, 30 = 5 days)
	logs := []CompletionLog{
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 24), 1),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 25), 1),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 26), 1),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 27), 1),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 30), 1),
	}

	summary := CalculateSummary([]Habit{h}, logs, []Pause{pause}, PeriodWeek, today)
	stat := summary.HabitSummaries[0]

	if stat.PausedDays != 2 {
		t.Errorf("PausedDays = %d, want 2", stat.PausedDays)
	}
	if stat.ScheduledDays != 5 {
		t.Errorf("ScheduledDays = %d, want 5 (7 - 2 paused)", stat.ScheduledDays)
	}
	if stat.CompletedDays != 5 {
		t.Errorf("CompletedDays = %d, want 5", stat.CompletedDays)
	}
}

func TestCalculateSummary_EmptyHabits(t *testing.T) {
	today := makeDate(2026, 3, 30)
	summary := CalculateSummary(nil, nil, nil, PeriodWeek, today)

	if summary.TotalHabits != 0 {
		t.Errorf("TotalHabits = %d, want 0", summary.TotalHabits)
	}
	if summary.ScheduledDays != 0 {
		t.Errorf("ScheduledDays = %d, want 0", summary.ScheduledDays)
	}
	if summary.CompletionRate != 0.0 {
		t.Errorf("CompletionRate = %f, want 0.0", summary.CompletionRate)
	}
	if len(summary.HabitSummaries) != 0 {
		t.Errorf("HabitSummaries length = %d, want 0", len(summary.HabitSummaries))
	}
	if summary.Period != "week" {
		t.Errorf("Period = %q, want %q", summary.Period, "week")
	}
}

