package core

import "testing"

func TestHabitIsValid(t *testing.T) {
	tests := []struct {
		name    string
		habit   Habit
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid build habit with duration goal",
			habit:   Habit{Name: "Running", GoalType: GoalTypeDuration, Direction: DirectionBuild},
			wantErr: false,
		},
		{
			name:    "valid avoid habit with boolean goal",
			habit:   Habit{Name: "No social media", GoalType: GoalTypeBoolean, Direction: DirectionAvoid},
			wantErr: false,
		},
		{
			name:    "valid build habit with count goal",
			habit:   Habit{Name: "Push-ups", GoalType: GoalTypeCount, Direction: DirectionBuild},
			wantErr: false,
		},
		{
			name:    "empty name is rejected",
			habit:   Habit{Name: "", GoalType: GoalTypeDuration, Direction: DirectionBuild},
			wantErr: true,
			errMsg:  "habit name is required",
		},
		{
			name:    "invalid goal type is rejected",
			habit:   Habit{Name: "Test", GoalType: "invalid", Direction: DirectionBuild},
			wantErr: true,
			errMsg:  "invalid goal type",
		},
		{
			name:    "invalid direction is rejected",
			habit:   Habit{Name: "Test", GoalType: GoalTypeDuration, Direction: "invalid"},
			wantErr: true,
			errMsg:  "invalid direction",
		},
		{
			name:    "negative goal value is rejected",
			habit:   Habit{Name: "Test", GoalType: GoalTypeDuration, Direction: DirectionBuild, GoalValue: -1},
			wantErr: true,
			errMsg:  "goal value must be >= 0",
		},
		{
			name:    "boolean with goal_value > 0 is rejected",
			habit:   Habit{Name: "Test", GoalType: GoalTypeBoolean, Direction: DirectionBuild, GoalValue: 5},
			wantErr: true,
			errMsg:  "goal value must be 0 for boolean goal type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.habit.IsValid()
			if tt.wantErr && errStr == "" {
				t.Errorf("expected error containing %q, got empty string", tt.errMsg)
			}
			if !tt.wantErr && errStr != "" {
				t.Errorf("expected no error, got %q", errStr)
			}
			if tt.wantErr && errStr != "" && tt.errMsg != "" {
				if !contains(errStr, tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, errStr)
				}
			}
		})
	}
}

func TestHabitScheduleIsValid(t *testing.T) {
	tests := []struct {
		name     string
		schedule HabitSchedule
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid daily schedule",
			schedule: HabitSchedule{ScheduleType: ScheduleDaily},
			wantErr:  false,
		},
		{
			name:     "valid specific_days schedule",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 3, 5}},
			wantErr:  false,
		},
		{
			name:     "valid weekdays schedule",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekdays},
			wantErr:  false,
		},
		{
			name:     "valid weekends schedule",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekends},
			wantErr:  false,
		},
		{
			name:     "valid times_per_week schedule",
			schedule: HabitSchedule{ScheduleType: ScheduleTimesPerWeek, TimesPerWeek: 3},
			wantErr:  false,
		},
		{
			name:     "valid every_n_days schedule",
			schedule: HabitSchedule{ScheduleType: ScheduleEveryNDays, EveryNDays: 2},
			wantErr:  false,
		},
		{
			name:     "valid weekly schedule",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekly},
			wantErr:  false,
		},
		{
			name:     "valid monthly schedule",
			schedule: HabitSchedule{ScheduleType: ScheduleMonthly, DayOfMonth: 15},
			wantErr:  false,
		},
		{
			name:     "all valid day_of_week values 1-7 pass",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}},
			wantErr:  false,
		},
		{
			name:     "invalid schedule type is rejected",
			schedule: HabitSchedule{ScheduleType: "invalid"},
			wantErr:  true,
			errMsg:   "invalid schedule type",
		},
		{
			name:     "specific_days with empty array is rejected",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{}},
			wantErr:  true,
			errMsg:   "specific_days schedule requires days_of_week",
		},
		{
			name:     "specific_days day value 0 is rejected (must be 1-7)",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{0, 3}},
			wantErr:  true,
			errMsg:   "days_of_week values must be between 1 (Mon) and 7 (Sun)",
		},
		{
			name:     "specific_days day value 8 is rejected (must be 1-7)",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 8}},
			wantErr:  true,
			errMsg:   "days_of_week values must be between 1 (Mon) and 7 (Sun)",
		},
		{
			name:     "times_per_week = 0 is rejected",
			schedule: HabitSchedule{ScheduleType: ScheduleTimesPerWeek, TimesPerWeek: 0},
			wantErr:  true,
			errMsg:   "times_per_week must be >= 1",
		},
		{
			name:     "every_n_days = 0 is rejected",
			schedule: HabitSchedule{ScheduleType: ScheduleEveryNDays, EveryNDays: 0},
			wantErr:  true,
			errMsg:   "every_n_days must be >= 1",
		},
		{
			name:     "monthly day_of_month = 0 is rejected",
			schedule: HabitSchedule{ScheduleType: ScheduleMonthly, DayOfMonth: 0},
			wantErr:  true,
			errMsg:   "day_of_month must be between 1 and 31",
		},
		{
			name:     "monthly day_of_month = 32 is rejected",
			schedule: HabitSchedule{ScheduleType: ScheduleMonthly, DayOfMonth: 32},
			wantErr:  true,
			errMsg:   "day_of_month must be between 1 and 31",
		},
		{
			name:     "invalid time_of_day format is rejected",
			schedule: HabitSchedule{ScheduleType: ScheduleDaily, TimeOfDay: "6:30"},
			wantErr:  true,
			errMsg:   "time_of_day must be in HH:MM format",
		},
		{
			name:     "valid time_of_day 06:30",
			schedule: HabitSchedule{ScheduleType: ScheduleDaily, TimeOfDay: "06:30"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.schedule.IsValid()
			if tt.wantErr && errStr == "" {
				t.Errorf("expected error containing %q, got empty string", tt.errMsg)
			}
			if !tt.wantErr && errStr != "" {
				t.Errorf("expected no error, got %q", errStr)
			}
			if tt.wantErr && errStr != "" && tt.errMsg != "" {
				if !contains(errStr, tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, errStr)
				}
			}
		})
	}
}

// contains checks if s contains substr.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
