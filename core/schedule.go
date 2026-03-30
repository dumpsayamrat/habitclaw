package core

import "time"

// ScheduleType defines how often a habit recurs
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

// HabitSchedule defines when a habit should be performed
type HabitSchedule struct {
	ID           string       `json:"id"`
	HabitID      string       `json:"habit_id"`
	UserID       string       `json:"user_id"`
	ScheduleType ScheduleType `json:"schedule_type"`
	DaysOfWeek   []int        `json:"days_of_week,omitempty"`
	TimesPerWeek int          `json:"times_per_week,omitempty"`
	EveryNDays   int          `json:"every_n_days,omitempty"`
	DayOfMonth   int          `json:"day_of_month,omitempty"`
	TimeOfDay    string       `json:"time_of_day,omitempty"`
	WindowStart  string       `json:"window_start,omitempty"`
	WindowEnd    string       `json:"window_end,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// IsScheduledDay returns true if the given date is a scheduled day
// for this habit based on its schedule type.
// createdAt is used as the reference start date for every_n_days calculation.
func (s HabitSchedule) IsScheduledDay(date time.Time, createdAt time.Time) bool {
	isoWeekday := int(date.Weekday())
	if isoWeekday == 0 {
		isoWeekday = 7 // Sunday = 7 in our convention
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
		return true
	case ScheduleEveryNDays:
		if s.EveryNDays <= 0 {
			return false
		}
		days := int(date.Sub(createdAt).Hours() / 24)
		return days%s.EveryNDays == 0
	case ScheduleWeekly:
		if len(s.DaysOfWeek) > 0 {
			return s.DaysOfWeek[0] == isoWeekday
		}
		return isoWeekday == 1
	case ScheduleMonthly:
		return date.Day() == s.DayOfMonth
	default:
		return false
	}
}

// isValidTimeFormat checks if a string is in "HH:MM" format
func isValidTimeFormat(t string) bool {
	if len(t) != 5 || t[2] != ':' {
		return false
	}
	h := (t[0]-'0')*10 + (t[1] - '0')
	m := (t[3]-'0')*10 + (t[4] - '0')
	if t[0] < '0' || t[0] > '9' || t[1] < '0' || t[1] > '9' {
		return false
	}
	if t[3] < '0' || t[3] > '9' || t[4] < '0' || t[4] > '9' {
		return false
	}
	return h <= 23 && m <= 59
}

// IsValid returns an error string if the schedule is invalid, empty string if valid
func (s HabitSchedule) IsValid() string {
	switch s.ScheduleType {
	case ScheduleDaily, ScheduleWeekdays, ScheduleWeekends:
		// no extra fields required
	case ScheduleSpecificDays:
		if len(s.DaysOfWeek) == 0 {
			return "specific_days schedule requires days_of_week"
		}
		for _, d := range s.DaysOfWeek {
			if d < 1 || d > 7 {
				return "days_of_week values must be between 1 (Mon) and 7 (Sun)"
			}
		}
	case ScheduleTimesPerWeek:
		if s.TimesPerWeek < 1 {
			return "times_per_week must be >= 1"
		}
	case ScheduleEveryNDays:
		if s.EveryNDays < 1 {
			return "every_n_days must be >= 1"
		}
	case ScheduleWeekly:
		// DaysOfWeek optional, defaults to Monday
	case ScheduleMonthly:
		if s.DayOfMonth < 1 || s.DayOfMonth > 31 {
			return "day_of_month must be between 1 and 31"
		}
	default:
		return "invalid schedule type"
	}

	if s.TimeOfDay != "" && !isValidTimeFormat(s.TimeOfDay) {
		return "time_of_day must be in HH:MM format"
	}
	if s.WindowStart != "" && !isValidTimeFormat(s.WindowStart) {
		return "window_start must be in HH:MM format"
	}
	if s.WindowEnd != "" && !isValidTimeFormat(s.WindowEnd) {
		return "window_end must be in HH:MM format"
	}

	return ""
}
