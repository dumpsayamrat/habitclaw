package core

import (
	"testing"
	"time"
)

func TestCalculateHeatmap_EmptyLogs(t *testing.T) {
	habit := newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule())
	from := makeDate(2026, 3, 24)
	to := makeDate(2026, 3, 30)

	days := CalculateHeatmap(habit, nil, nil, from, to)

	if len(days) != 7 {
		t.Fatalf("expected 7 days, got %d", len(days))
	}
	for _, d := range days {
		if d.Value != 0 || d.GoalMet || d.IsPaused || d.IsSlip {
			t.Errorf("day %s: expected zero-value, got %+v", d.Date, d)
		}
	}
}

func TestCalculateHeatmap_CompletionGoalMet(t *testing.T) {
	habit := newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule())
	// GoalValue is 30 (set by newTestHabit)
	from := makeDate(2026, 3, 30)
	to := makeDate(2026, 3, 30)
	logs := []CompletionLog{
		newTestLog(habit.ID, LogCompletion, makeDate(2026, 3, 30), 45),
	}

	days := CalculateHeatmap(habit, logs, nil, from, to)

	if len(days) != 1 {
		t.Fatalf("expected 1 day, got %d", len(days))
	}
	if days[0].Value != 45 {
		t.Errorf("expected Value=45, got %d", days[0].Value)
	}
	if !days[0].GoalMet {
		t.Errorf("expected GoalMet=true for value 45 >= goal 30")
	}
}

func TestCalculateHeatmap_CompletionGoalNotMet(t *testing.T) {
	habit := newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule())
	// GoalValue is 30; log value 15 < 30
	from := makeDate(2026, 3, 30)
	to := makeDate(2026, 3, 30)
	logs := []CompletionLog{
		newTestLog(habit.ID, LogCompletion, makeDate(2026, 3, 30), 15),
	}

	days := CalculateHeatmap(habit, logs, nil, from, to)

	if len(days) != 1 {
		t.Fatalf("expected 1 day, got %d", len(days))
	}
	if days[0].Value != 15 {
		t.Errorf("expected Value=15, got %d", days[0].Value)
	}
	if days[0].GoalMet {
		t.Errorf("expected GoalMet=false for value 15 < goal 30")
	}
}

func TestCalculateHeatmap_SlipLog(t *testing.T) {
	habit := newTestHabit("NoSocialMedia", DirectionAvoid, GoalTypeBoolean, dailySchedule())
	from := makeDate(2026, 3, 30)
	to := makeDate(2026, 3, 30)
	logs := []CompletionLog{
		newTestLog(habit.ID, LogSlip, makeDate(2026, 3, 30), int(SlipFull)),
	}

	days := CalculateHeatmap(habit, logs, nil, from, to)

	if len(days) != 1 {
		t.Fatalf("expected 1 day, got %d", len(days))
	}
	if !days[0].IsSlip {
		t.Errorf("expected IsSlip=true")
	}
	if days[0].Severity != int(SlipFull) {
		t.Errorf("expected Severity=%d, got %d", int(SlipFull), days[0].Severity)
	}
}

func TestCalculateHeatmap_PausedDays(t *testing.T) {
	habit := newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule())
	habitID := habit.ID
	from := makeDate(2026, 3, 28)
	to := makeDate(2026, 3, 30)

	// Global pause (HabitID==nil) covers March 28
	globalPause := newTestPause(nil, makeDate(2026, 3, 28), makeDate(2026, 3, 28))
	// Habit-specific pause covers March 29
	habitPause := newTestPause(&habitID, makeDate(2026, 3, 29), makeDate(2026, 3, 29))

	days := CalculateHeatmap(habit, nil, []Pause{globalPause, habitPause}, from, to)

	if len(days) != 3 {
		t.Fatalf("expected 3 days, got %d", len(days))
	}
	if !days[0].IsPaused {
		t.Errorf("day 0 (Mar 28): expected IsPaused=true (global pause)")
	}
	if !days[1].IsPaused {
		t.Errorf("day 1 (Mar 29): expected IsPaused=true (habit-specific pause)")
	}
	if days[2].IsPaused {
		t.Errorf("day 2 (Mar 30): expected IsPaused=false")
	}
}

func TestCalculateHeatmap_SingleDay(t *testing.T) {
	habit := newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule())
	from := makeDate(2026, 3, 30)
	to := makeDate(2026, 3, 30)

	days := CalculateHeatmap(habit, nil, nil, from, to)

	if len(days) != 1 {
		t.Errorf("expected exactly 1 day, got %d", len(days))
	}
	if days[0].Date != "2026-03-30" {
		t.Errorf("expected date 2026-03-30, got %s", days[0].Date)
	}
}

func TestCalculateHeatmap_BooleanGoalType(t *testing.T) {
	habit := newTestHabit("Meditate", DirectionBuild, GoalTypeBoolean, dailySchedule())
	// For boolean goal, any completion value sets GoalMet=true
	from := makeDate(2026, 3, 30)
	to := makeDate(2026, 3, 30)
	logs := []CompletionLog{
		newTestLog(habit.ID, LogCompletion, makeDate(2026, 3, 30), 1),
	}

	days := CalculateHeatmap(habit, logs, nil, from, to)

	if len(days) != 1 {
		t.Fatalf("expected 1 day, got %d", len(days))
	}
	if !days[0].GoalMet {
		t.Errorf("expected GoalMet=true for boolean goal type regardless of value")
	}
}

func TestCalculateHeatmap_CancelledPauseWithResumeFrom(t *testing.T) {
	habit := newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule())
	from := makeDate(2026, 3, 26)
	to := makeDate(2026, 3, 30)

	// Pause originally from Mar 26 to Mar 30, but cancelled with ResumeFrom Mar 29
	// So only Mar 26, 27, 28 are paused (up to ResumeFrom-1)
	resumeFrom := makeDate(2026, 3, 29)
	cancelledAt := time.Now()
	pause := Pause{
		ID:          "pause-cancelled",
		UserID:      "local",
		HabitID:     nil,
		FromDate:    makeDate(2026, 3, 26),
		ToDate:      makeDate(2026, 3, 30),
		Reason:      "test",
		CancelledAt: &cancelledAt,
		ResumeFrom:  &resumeFrom,
		CreatedAt:   time.Now(),
	}

	days := CalculateHeatmap(habit, nil, []Pause{pause}, from, to)

	if len(days) != 5 {
		t.Fatalf("expected 5 days, got %d", len(days))
	}
	// Mar 26, 27, 28 should be paused
	for i := 0; i < 3; i++ {
		if !days[i].IsPaused {
			t.Errorf("day %d (%s): expected IsPaused=true", i, days[i].Date)
		}
	}
	// Mar 29, 30 should NOT be paused
	for i := 3; i < 5; i++ {
		if days[i].IsPaused {
			t.Errorf("day %d (%s): expected IsPaused=false", i, days[i].Date)
		}
	}
}
