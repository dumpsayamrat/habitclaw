package core

import (
	"math"
	"testing"
)

func approxEqualAlignment(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestCalculateGoalAlignment_Exceeding(t *testing.T) {
	today := makeDate(2026, 3, 30)

	h := newTestHabit("Running", DirectionBuild, GoalTypeCount, dailySchedule())
	h.GoalValue = 30
	h.CreatedAt = makeDate(2026, 1, 1)

	// 7 days, avg = 35 (goal 30) → exceeding
	logs := []CompletionLog{
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 24), 35),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 25), 35),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 26), 35),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 27), 35),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 28), 35),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 29), 35),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 30), 35),
	}

	result := CalculateGoalAlignment([]Habit{h}, logs, nil, PeriodWeek, today)

	if len(result.Habits) != 1 {
		t.Fatalf("expected 1 habit alignment, got %d", len(result.Habits))
	}
	ha := result.Habits[0]

	if !approxEqualAlignment(ha.AverageActual, 35.0, 0.01) {
		t.Errorf("AverageActual = %.4f, want 35.0", ha.AverageActual)
	}
	if !approxEqualAlignment(ha.Alignment, 1.0, 0.001) {
		t.Errorf("Alignment = %.4f, want 1.0 (capped)", ha.Alignment)
	}
	if ha.Status != "exceeding" {
		t.Errorf("Status = %q, want %q", ha.Status, "exceeding")
	}
}

func TestCalculateGoalAlignment_OnTrack(t *testing.T) {
	today := makeDate(2026, 3, 30)

	h := newTestHabit("Running", DirectionBuild, GoalTypeCount, dailySchedule())
	h.GoalValue = 30
	h.CreatedAt = makeDate(2026, 1, 1)

	// avg = 25, goal = 30, ratio = 0.8333 → on track
	logs := []CompletionLog{
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 24), 25),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 25), 25),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 26), 25),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 27), 25),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 28), 25),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 29), 25),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 30), 25),
	}

	result := CalculateGoalAlignment([]Habit{h}, logs, nil, PeriodWeek, today)

	if len(result.Habits) != 1 {
		t.Fatalf("expected 1 habit alignment, got %d", len(result.Habits))
	}
	ha := result.Habits[0]

	wantAlignment := 25.0 / 30.0
	if !approxEqualAlignment(ha.AverageActual, 25.0, 0.01) {
		t.Errorf("AverageActual = %.4f, want 25.0", ha.AverageActual)
	}
	if !approxEqualAlignment(ha.Alignment, wantAlignment, 0.001) {
		t.Errorf("Alignment = %.4f, want %.4f", ha.Alignment, wantAlignment)
	}
	if ha.Status != "on track" {
		t.Errorf("Status = %q, want %q", ha.Status, "on track")
	}
}

func TestCalculateGoalAlignment_NeedsAttention(t *testing.T) {
	today := makeDate(2026, 3, 30)

	h := newTestHabit("Running", DirectionBuild, GoalTypeCount, dailySchedule())
	h.GoalValue = 30
	h.CreatedAt = makeDate(2026, 1, 1)

	// avg = 16, goal = 30, ratio ≈ 0.533 → needs attention
	logs := []CompletionLog{
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 24), 16),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 25), 16),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 26), 16),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 27), 16),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 28), 16),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 29), 16),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 30), 16),
	}

	result := CalculateGoalAlignment([]Habit{h}, logs, nil, PeriodWeek, today)

	if len(result.Habits) != 1 {
		t.Fatalf("expected 1 habit alignment, got %d", len(result.Habits))
	}
	ha := result.Habits[0]

	wantAlignment := 16.0 / 30.0
	if !approxEqualAlignment(ha.AverageActual, 16.0, 0.01) {
		t.Errorf("AverageActual = %.4f, want 16.0", ha.AverageActual)
	}
	if !approxEqualAlignment(ha.Alignment, wantAlignment, 0.001) {
		t.Errorf("Alignment = %.4f, want %.4f", ha.Alignment, wantAlignment)
	}
	if ha.Status != "needs attention" {
		t.Errorf("Status = %q, want %q", ha.Status, "needs attention")
	}
}

func TestCalculateGoalAlignment_OffTrack(t *testing.T) {
	today := makeDate(2026, 3, 30)

	h := newTestHabit("Running", DirectionBuild, GoalTypeCount, dailySchedule())
	h.GoalValue = 30
	h.CreatedAt = makeDate(2026, 1, 1)

	// avg = 10, goal = 30, ratio ≈ 0.333 → off track
	logs := []CompletionLog{
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 24), 10),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 25), 10),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 26), 10),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 27), 10),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 28), 10),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 29), 10),
		newTestLog(h.ID, LogCompletion, makeDate(2026, 3, 30), 10),
	}

	result := CalculateGoalAlignment([]Habit{h}, logs, nil, PeriodWeek, today)

	if len(result.Habits) != 1 {
		t.Fatalf("expected 1 habit alignment, got %d", len(result.Habits))
	}
	ha := result.Habits[0]

	wantAlignment := 10.0 / 30.0
	if !approxEqualAlignment(ha.AverageActual, 10.0, 0.01) {
		t.Errorf("AverageActual = %.4f, want 10.0", ha.AverageActual)
	}
	if !approxEqualAlignment(ha.Alignment, wantAlignment, 0.001) {
		t.Errorf("Alignment = %.4f, want %.4f", ha.Alignment, wantAlignment)
	}
	if ha.Status != "off track" {
		t.Errorf("Status = %q, want %q", ha.Status, "off track")
	}
}

func TestCalculateGoalAlignment_OverallScoreAverages(t *testing.T) {
	today := makeDate(2026, 3, 30)

	h1 := newTestHabit("Running", DirectionBuild, GoalTypeCount, dailySchedule())
	h1.GoalValue = 30
	h1.CreatedAt = makeDate(2026, 1, 1)

	h2 := newTestHabit("Reading", DirectionBuild, GoalTypeCount, dailySchedule())
	h2.GoalValue = 30
	h2.CreatedAt = makeDate(2026, 1, 1)

	// h1: avg=35 → alignment=1.0 (capped)
	// h2: avg=15 → alignment=0.5
	// overall = (1.0 + 0.5) / 2 = 0.75
	logs := []CompletionLog{
		newTestLog(h1.ID, LogCompletion, makeDate(2026, 3, 30), 35),
		newTestLog(h2.ID, LogCompletion, makeDate(2026, 3, 30), 15),
	}

	result := CalculateGoalAlignment([]Habit{h1, h2}, logs, nil, PeriodToday, today)

	if len(result.Habits) != 2 {
		t.Fatalf("expected 2 habit alignments, got %d", len(result.Habits))
	}

	wantOverall := (1.0 + 0.5) / 2.0
	if !approxEqualAlignment(result.OverallScore, wantOverall, 0.001) {
		t.Errorf("OverallScore = %.4f, want %.4f", result.OverallScore, wantOverall)
	}
}

func TestCalculateGoalAlignment_EmptyHabits(t *testing.T) {
	today := makeDate(2026, 3, 30)

	result := CalculateGoalAlignment(nil, nil, nil, PeriodWeek, today)

	if result.OverallScore != 0.0 {
		t.Errorf("OverallScore = %f, want 0.0", result.OverallScore)
	}
	if len(result.Habits) != 0 {
		t.Errorf("Habits length = %d, want 0", len(result.Habits))
	}
	if result.Period != "week" {
		t.Errorf("Period = %q, want %q", result.Period, "week")
	}
}

func TestCalculateGoalAlignment_AvoidHabitsSkipped(t *testing.T) {
	today := makeDate(2026, 3, 30)

	avoid := newTestHabit("No sugar", DirectionAvoid, GoalTypeBoolean, dailySchedule())
	avoid.CreatedAt = makeDate(2026, 1, 1)

	result := CalculateGoalAlignment([]Habit{avoid}, nil, nil, PeriodWeek, today)

	if len(result.Habits) != 0 {
		t.Errorf("expected avoid habits to be skipped, got %d habits", len(result.Habits))
	}
	if result.OverallScore != 0.0 {
		t.Errorf("OverallScore = %f, want 0.0", result.OverallScore)
	}
}
