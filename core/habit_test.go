package core

import "testing"

func TestHabitValidate(t *testing.T) {
	tests := []struct {
		name    string
		habit   Habit
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid build habit with duration goal",
			habit:   Habit{Name: "Running", GoalType: GoalDuration, Direction: DirectionBuild},
			wantErr: false,
		},
		{
			name:    "valid avoid habit with boolean goal",
			habit:   Habit{Name: "No social media", GoalType: GoalBoolean, Direction: DirectionAvoid},
			wantErr: false,
		},
		{
			name:    "valid build habit with count goal",
			habit:   Habit{Name: "Push-ups", GoalType: GoalCount, Direction: DirectionBuild},
			wantErr: false,
		},
		{
			name:    "empty name is rejected",
			habit:   Habit{Name: "", GoalType: GoalDuration, Direction: DirectionBuild},
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
			habit:   Habit{Name: "Test", GoalType: GoalDuration, Direction: "invalid"},
			wantErr: true,
			errMsg:  "invalid direction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.habit.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("expected error containing %q, got nil", tt.errMsg)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if got := err.Error(); got != tt.errMsg && !contains(got, tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, got)
				}
			}
		})
	}
}

func TestHabitScheduleValidate(t *testing.T) {
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
			name:     "invalid schedule type is rejected",
			schedule: HabitSchedule{ScheduleType: "invalid"},
			wantErr:  true,
			errMsg:   "invalid schedule type",
		},
		{
			name:     "day_of_week value 0 is rejected (must be 1-7)",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{0, 3}},
			wantErr:  true,
			errMsg:   "days_of_week values must be between 1 (Mon) and 7 (Sun)",
		},
		{
			name:     "day_of_week value 8 is rejected (must be 1-7)",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 8}},
			wantErr:  true,
			errMsg:   "days_of_week values must be between 1 (Mon) and 7 (Sun)",
		},
		{
			name:     "all valid day_of_week values 1-7 pass",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7}},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schedule.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("expected error containing %q, got nil", tt.errMsg)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if got := err.Error(); !contains(got, tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, got)
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
