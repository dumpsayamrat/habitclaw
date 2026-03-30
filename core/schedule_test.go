package core

import (
	"testing"
	"time"
)

func TestIsScheduledDay(t *testing.T) {
	tests := []struct {
		name     string
		schedule HabitSchedule
		date     time.Time
		dateStr  string // for error messages
		want     bool
	}{
		// Daily schedule — every day returns true
		{
			name:     "daily: Monday is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleDaily},
			date:     makeDate(2026, 3, 30), // Monday
			dateStr:  "2026-03-30 Mon",
			want:     true,
		},
		{
			name:     "daily: Saturday is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleDaily},
			date:     makeDate(2026, 3, 28), // Saturday
			dateStr:  "2026-03-28 Sat",
			want:     true,
		},
		{
			name:     "daily: Sunday is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleDaily},
			date:     makeDate(2026, 3, 29), // Sunday
			dateStr:  "2026-03-29 Sun",
			want:     true,
		},

		// Specific days [1,3,5] = Mon, Wed, Fri
		{
			name:     "specific_days [1,3,5]: Monday (1) is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 3, 5}},
			date:     makeDate(2026, 3, 30), // Monday
			dateStr:  "2026-03-30 Mon",
			want:     true,
		},
		{
			name:     "specific_days [1,3,5]: Wednesday (3) is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 3, 5}},
			date:     makeDate(2026, 3, 25), // Wednesday
			dateStr:  "2026-03-25 Wed",
			want:     true,
		},
		{
			name:     "specific_days [1,3,5]: Friday (5) is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 3, 5}},
			date:     makeDate(2026, 3, 27), // Friday
			dateStr:  "2026-03-27 Fri",
			want:     true,
		},
		{
			name:     "specific_days [1,3,5]: Tuesday (2) is NOT scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 3, 5}},
			date:     makeDate(2026, 3, 24), // Tuesday
			dateStr:  "2026-03-24 Tue",
			want:     false,
		},
		{
			name:     "specific_days [1,3,5]: Saturday (6) is NOT scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 3, 5}},
			date:     makeDate(2026, 3, 28), // Saturday
			dateStr:  "2026-03-28 Sat",
			want:     false,
		},
		{
			name:     "specific_days [1,3,5]: Sunday (7) is NOT scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 3, 5}},
			date:     makeDate(2026, 3, 29), // Sunday
			dateStr:  "2026-03-29 Sun",
			want:     false,
		},

		// Weekdays: Mon-Fri true, Sat-Sun false
		{
			name:     "weekdays: Monday is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekdays},
			date:     makeDate(2026, 3, 30), // Monday
			dateStr:  "2026-03-30 Mon",
			want:     true,
		},
		{
			name:     "weekdays: Friday is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekdays},
			date:     makeDate(2026, 3, 27), // Friday
			dateStr:  "2026-03-27 Fri",
			want:     true,
		},
		{
			name:     "weekdays: Saturday is NOT scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekdays},
			date:     makeDate(2026, 3, 28), // Saturday
			dateStr:  "2026-03-28 Sat",
			want:     false,
		},
		{
			name:     "weekdays: Sunday is NOT scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekdays},
			date:     makeDate(2026, 3, 29), // Sunday
			dateStr:  "2026-03-29 Sun",
			want:     false,
		},

		// Weekends: Sat-Sun true, Mon-Fri false
		{
			name:     "weekends: Saturday is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekends},
			date:     makeDate(2026, 3, 28), // Saturday
			dateStr:  "2026-03-28 Sat",
			want:     true,
		},
		{
			name:     "weekends: Sunday is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekends},
			date:     makeDate(2026, 3, 29), // Sunday
			dateStr:  "2026-03-29 Sun",
			want:     true,
		},
		{
			name:     "weekends: Monday is NOT scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekends},
			date:     makeDate(2026, 3, 30), // Monday
			dateStr:  "2026-03-30 Mon",
			want:     false,
		},
		{
			name:     "weekends: Wednesday is NOT scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleWeekends},
			date:     makeDate(2026, 3, 25), // Wednesday
			dateStr:  "2026-03-25 Wed",
			want:     false,
		},

		// Monthly with day_of_month=15
		{
			name:     "monthly day 15: the 15th is scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleMonthly, DayOfMonth: 15},
			date:     makeDate(2026, 3, 15),
			dateStr:  "2026-03-15",
			want:     true,
		},
		{
			name:     "monthly day 15: the 14th is NOT scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleMonthly, DayOfMonth: 15},
			date:     makeDate(2026, 3, 14),
			dateStr:  "2026-03-14",
			want:     false,
		},
		{
			name:     "monthly day 15: the 16th is NOT scheduled",
			schedule: HabitSchedule{ScheduleType: ScheduleMonthly, DayOfMonth: 15},
			date:     makeDate(2026, 3, 16),
			dateStr:  "2026-03-16",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schedule.IsScheduledDay(tt.date)
			if got != tt.want {
				t.Errorf("IsScheduledDay(%s) = %v, want %v (weekday=%s)",
					tt.dateStr, got, tt.want, tt.date.Weekday())
			}
		})
	}
}
