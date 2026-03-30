package core

import "time"

type LogType string

const (
	LogCompletion LogType = "completion"
	LogSlip       LogType = "slip"
)

type SlipSeverity int

const (
	SlipMinor SlipSeverity = 1
	SlipFull  SlipSeverity = 2
)

type CompletionLog struct {
	ID        string    `json:"id"`
	HabitID   string    `json:"habit_id"`
	UserID    string    `json:"user_id"`
	LogType   LogType   `json:"log_type"`
	Date      time.Time `json:"date"`
	Value     int       `json:"value"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}
