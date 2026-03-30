package core

import (
	"errors"
	"time"
)

type ScheduleType string

const (
	ScheduleDaily        ScheduleType = "daily"
	ScheduleSpecificDays ScheduleType = "specific_days"
	ScheduleWeekdays     ScheduleType = "weekdays"
	ScheduleWeekends     ScheduleType = "weekends"
	ScheduleTimesPerWeek ScheduleType = "times_per_week"
	ScheduleEveryNDays   ScheduleType = "every_n_days"
	ScheduleWeekly       ScheduleType = "weekly"
	ScheduleMonthly      ScheduleType = "monthly"
)

type HabitSchedule struct {
	ID           string       `json:"id"`
	HabitID      string       `json:"habit_id"`
	ScheduleType ScheduleType `json:"schedule_type"`
	DaysOfWeek   []int        `json:"days_of_week,omitempty"`
	TimesPerWeek int          `json:"times_per_week,omitempty"`
	EveryNDays   int          `json:"every_n_days,omitempty"`
	DayOfMonth   int          `json:"day_of_month,omitempty"`
	TimeOfDay    string       `json:"time_of_day,omitempty"`
	WindowStart  string       `json:"window_start,omitempty"`
	WindowEnd    string       `json:"window_end,omitempty"`
}

func (s HabitSchedule) Validate() error {
	switch s.ScheduleType {
	case ScheduleDaily, ScheduleSpecificDays, ScheduleWeekdays, ScheduleWeekends,
		ScheduleTimesPerWeek, ScheduleEveryNDays, ScheduleWeekly, ScheduleMonthly:
	default:
		return errors.New("invalid schedule type")
	}
	for _, d := range s.DaysOfWeek {
		if d < 1 || d > 7 {
			return errors.New("days_of_week values must be between 1 (Mon) and 7 (Sun)")
		}
	}
	return nil
}

// IsScheduledDay returns true if the given date is a scheduled day for this schedule.
// Uses ISO weekday: Monday=1, Sunday=7.
func (s HabitSchedule) IsScheduledDay(date time.Time) bool {
	isoWeekday := int(date.Weekday())
	if isoWeekday == 0 {
		isoWeekday = 7 // Sunday = 7
	}

	switch s.ScheduleType {
	case ScheduleDaily:
		return true
	case ScheduleSpecificDays:
		for _, d := range s.DaysOfWeek {
			if d == isoWeekday {
				return true
			}
		}
		return false
	case ScheduleWeekdays:
		return isoWeekday >= 1 && isoWeekday <= 5
	case ScheduleWeekends:
		return isoWeekday == 6 || isoWeekday == 7
	case ScheduleTimesPerWeek:
		// For times_per_week, every day is potentially schedulable
		// The constraint is checked at a higher level (weekly count)
		return true
	case ScheduleWeekly:
		// Default to same weekday as habit creation or Monday
		if len(s.DaysOfWeek) > 0 {
			for _, d := range s.DaysOfWeek {
				if d == isoWeekday {
					return true
				}
			}
			return false
		}
		return isoWeekday == 1 // default Monday
	case ScheduleMonthly:
		return date.Day() == s.DayOfMonth
	case ScheduleEveryNDays:
		// Can't determine without a start date reference — return true for now
		return true
	default:
		return false
	}
}
