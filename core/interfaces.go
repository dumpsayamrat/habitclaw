package core

import "time"

// HeatmapDay represents a single day in a heatmap calendar view.
type HeatmapDay struct {
	Date      string `json:"date"`
	Value     int    `json:"value"`
	GoalMet   bool   `json:"goal_met"`
	IsPaused  bool   `json:"is_paused"`
	IsSlip    bool   `json:"is_slip,omitempty"`
	Severity  int    `json:"severity,omitempty"`
}

// HabitStore defines the persistence interface for all habit data.
// Implementations include SQLite (local) and PostgreSQL (cloud).
type HabitStore interface {
	// habits
	CreateHabit(habit Habit) error
	ListHabits(userID string) ([]Habit, error)
	UpdateHabit(habit Habit) error
	ArchiveHabit(id string) error

	// schedules
	SetSchedule(schedule HabitSchedule) error
	GetSchedule(habitID string) (HabitSchedule, error)

	// logs
	LogCompletion(log CompletionLog) error
	LogSlip(log CompletionLog) error
	GetLogs(userID string, from, to time.Time) ([]CompletionLog, error)
	DeleteLog(id string) error

	// pauses
	CreatePause(pause Pause) error
	ListPauses(userID string, status string) ([]Pause, error)
	CancelPause(id string, resumeFrom time.Time) error

	// computed
	GetStreaks(userID string) ([]Streak, error)
	GetSummary(userID string, period Period) (Summary, error)
	GetGoalAlignment(userID string, period Period) (GoalAlignment, error)
	GetHeatmap(habitID string, from, to time.Time) ([]HeatmapDay, error)
}
