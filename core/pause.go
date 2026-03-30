package core

import "time"

type Pause struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	HabitID     *string    `json:"habit_id,omitempty"`
	FromDate    time.Time  `json:"from_date"`
	ToDate      time.Time  `json:"to_date"`
	Reason      string     `json:"reason"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
	ResumeFrom  *time.Time `json:"resume_from,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
