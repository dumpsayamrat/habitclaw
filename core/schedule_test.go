package core

import (
	"testing"
	"time"
)

func TestIsScheduledDay(t *testing.T) {
	// createdAt for every_n_days tests: 2026-03-23 (Monday)
	createdAt := makeDate(2026, 3, 23)

	tests := []struct {
		name      string
		schedule  HabitSchedule
		date      time.Time
		createdAt time.Time
		dateStr   string // for error messages
		want      bool
	}{
		// Daily schedule — every day returns true
		{
			name:      "daily: Monday is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleDaily},
			date:      makeDate(2026, 3, 30), // Monday
			createdAt: time.Time{},
			dateStr:   "2026-03-30 Mon",
			want:      true,
		},
		{
			name:      "daily: Saturday is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleDaily},
			date:      makeDate(2026, 3, 28), // Saturday
			createdAt: time.Time{},
			dateStr:   "2026-03-28 Sat",
			want:      true,
		},
		{
			name:      "daily: Sunday is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleDaily},
			date:      makeDate(2026, 3, 29), // Sunday
			createdAt: time.Time{},
			dateStr:   "2026-03-29 Sun",
			want:      true,
		},

		// Specific days [1,2,4] = Mon, Tue, Thu
		{
			name:      "specific_days [1,2,4]: Monday (1) is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 2, 4}},
			date:      makeDate(2026, 3, 30), // Monday
			createdAt: time.Time{},
			dateStr:   "2026-03-30 Mon",
			want:      true,
		},
		{
			name:      "specific_days [1,2,4]: Tuesday (2) is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 2, 4}},
			date:      makeDate(2026, 3, 31), // Tuesday
			createdAt: time.Time{},
			dateStr:   "2026-03-31 Tue",
			want:      true,
		},
		{
			name:      "specific_days [1,2,4]: Wednesday (3) is NOT scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 2, 4}},
			date:      makeDate(2026, 3, 25), // Wednesday
			createdAt: time.Time{},
			dateStr:   "2026-03-25 Wed",
			want:      false,
		},
		{
			name:      "specific_days [1,2,4]: Thursday (4) is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 2, 4}},
			date:      makeDate(2026, 3, 26), // Thursday
			createdAt: time.Time{},
			dateStr:   "2026-03-26 Thu",
			want:      true,
		},
		{
			name:      "specific_days [1,2,4]: Sunday (7) is NOT scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 2, 4}},
			date:      makeDate(2026, 3, 29), // Sunday
			createdAt: time.Time{},
			dateStr:   "2026-03-29 Sun",
			want:      false,
		},

		// Weekdays: Mon-Fri true, Sat-Sun false
		{
			name:      "weekdays: Monday is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleWeekdays},
			date:      makeDate(2026, 3, 30), // Monday
			createdAt: time.Time{},
			dateStr:   "2026-03-30 Mon",
			want:      true,
		},
		{
			name:      "weekdays: Friday is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleWeekdays},
			date:      makeDate(2026, 3, 27), // Friday
			createdAt: time.Time{},
			dateStr:   "2026-03-27 Fri",
			want:      true,
		},
		{
			name:      "weekdays: Saturday is NOT scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleWeekdays},
			date:      makeDate(2026, 3, 28), // Saturday
			createdAt: time.Time{},
			dateStr:   "2026-03-28 Sat",
			want:      false,
		},
		{
			name:      "weekdays: Sunday is NOT scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleWeekdays},
			date:      makeDate(2026, 3, 29), // Sunday
			createdAt: time.Time{},
			dateStr:   "2026-03-29 Sun",
			want:      false,
		},

		// Weekends: Sat-Sun true, Mon-Fri false
		{
			name:      "weekends: Saturday is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleWeekends},
			date:      makeDate(2026, 3, 28), // Saturday
			createdAt: time.Time{},
			dateStr:   "2026-03-28 Sat",
			want:      true,
		},
		{
			name:      "weekends: Sunday is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleWeekends},
			date:      makeDate(2026, 3, 29), // Sunday
			createdAt: time.Time{},
			dateStr:   "2026-03-29 Sun",
			want:      true,
		},
		{
			name:      "weekends: Monday is NOT scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleWeekends},
			date:      makeDate(2026, 3, 30), // Monday
			createdAt: time.Time{},
			dateStr:   "2026-03-30 Mon",
			want:      false,
		},

		// every_n_days=3: day0=true, day1=false, day2=false, day3=true, day6=true
		{
			name:      "every_n_days=3: day 0 (createdAt) is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleEveryNDays, EveryNDays: 3},
			date:      makeDate(2026, 3, 23), // day 0
			createdAt: createdAt,
			dateStr:   "2026-03-23 day0",
			want:      true,
		},
		{
			name:      "every_n_days=3: day 1 is NOT scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleEveryNDays, EveryNDays: 3},
			date:      makeDate(2026, 3, 24), // day 1
			createdAt: createdAt,
			dateStr:   "2026-03-24 day1",
			want:      false,
		},
		{
			name:      "every_n_days=3: day 2 is NOT scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleEveryNDays, EveryNDays: 3},
			date:      makeDate(2026, 3, 25), // day 2
			createdAt: createdAt,
			dateStr:   "2026-03-25 day2",
			want:      false,
		},
		{
			name:      "every_n_days=3: day 3 is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleEveryNDays, EveryNDays: 3},
			date:      makeDate(2026, 3, 26), // day 3
			createdAt: createdAt,
			dateStr:   "2026-03-26 day3",
			want:      true,
		},
		{
			name:      "every_n_days=3: day 6 is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleEveryNDays, EveryNDays: 3},
			date:      makeDate(2026, 3, 29), // day 6
			createdAt: createdAt,
			dateStr:   "2026-03-29 day6",
			want:      true,
		},

		// Weekly DaysOfWeek=[3] (Wed): Wed=true, Thu=false
		{
			name:      "weekly DaysOfWeek=[3]: Wednesday is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleWeekly, DaysOfWeek: []int{3}},
			date:      makeDate(2026, 3, 25), // Wednesday
			createdAt: time.Time{},
			dateStr:   "2026-03-25 Wed",
			want:      true,
		},
		{
			name:      "weekly DaysOfWeek=[3]: Thursday is NOT scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleWeekly, DaysOfWeek: []int{3}},
			date:      makeDate(2026, 3, 26), // Thursday
			createdAt: time.Time{},
			dateStr:   "2026-03-26 Thu",
			want:      false,
		},

		// Monthly with day_of_month=15
		{
			name:      "monthly day 15: the 15th is scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleMonthly, DayOfMonth: 15},
			date:      makeDate(2026, 3, 15),
			createdAt: time.Time{},
			dateStr:   "2026-03-15",
			want:      true,
		},
		{
			name:      "monthly day 15: the 14th is NOT scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleMonthly, DayOfMonth: 15},
			date:      makeDate(2026, 3, 14),
			createdAt: time.Time{},
			dateStr:   "2026-03-14",
			want:      false,
		},
		{
			name:      "monthly day 15: the 16th is NOT scheduled",
			schedule:  HabitSchedule{ScheduleType: ScheduleMonthly, DayOfMonth: 15},
			date:      makeDate(2026, 3, 16),
			createdAt: time.Time{},
			dateStr:   "2026-03-16",
			want:      false,
		},

		// times_per_week: always returns true
		{
			name:      "times_per_week: Monday returns true",
			schedule:  HabitSchedule{ScheduleType: ScheduleTimesPerWeek, TimesPerWeek: 3},
			date:      makeDate(2026, 3, 30), // Monday
			createdAt: time.Time{},
			dateStr:   "2026-03-30 Mon",
			want:      true,
		},
		{
			name:      "times_per_week: Saturday returns true",
			schedule:  HabitSchedule{ScheduleType: ScheduleTimesPerWeek, TimesPerWeek: 3},
			date:      makeDate(2026, 3, 28), // Saturday
			createdAt: time.Time{},
			dateStr:   "2026-03-28 Sat",
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schedule.IsScheduledDay(tt.date, tt.createdAt)
			if got != tt.want {
				t.Errorf("IsScheduledDay(%s) = %v, want %v (weekday=%s)",
					tt.dateStr, got, tt.want, tt.date.Weekday())
			}
		})
	}
}
