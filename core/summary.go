package core

type Period string

const (
	PeriodToday Period = "today"
	PeriodWeek  Period = "week"
	PeriodMonth Period = "month"
)

type Summary struct {
	Period         string      `json:"period"`
	TotalHabits    int         `json:"total_habits"`
	ScheduledDays  int         `json:"scheduled_days"`
	CompletionRate float64     `json:"completion_rate"`
	HabitSummaries []HabitStat `json:"habit_summaries"`
}

type HabitStat struct {
	HabitID       string  `json:"habit_id"`
	Name          string  `json:"name"`
	Direction     string  `json:"direction"`
	ScheduledDays int     `json:"scheduled_days"`
	CompletedDays int     `json:"completed_days"`
	SlippedDays   int     `json:"slipped_days,omitempty"`
	PausedDays    int     `json:"paused_days"`
	GoalMetRate   float64 `json:"goal_met_rate"`
	TotalValue    int     `json:"total_value"`
	CurrentStreak int     `json:"current_streak"`
}
