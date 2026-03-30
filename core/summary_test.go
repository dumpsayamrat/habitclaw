package core

import "testing"

func TestPeriodConstants(t *testing.T) {
	// Verify Period constants match the spec values
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
	// Verify Summary struct can be initialized with expected fields
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
