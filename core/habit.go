package core

import (
	"errors"
	"time"
)

type HabitDirection string

const (
	DirectionBuild HabitDirection = "build"
	DirectionAvoid HabitDirection = "avoid"
)

type GoalType string

const (
	GoalDuration GoalType = "duration"
	GoalCount    GoalType = "count"
	GoalBoolean  GoalType = "boolean"
)

type Habit struct {
	ID          string         `json:"id"`
	UserID      string         `json:"user_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	GoalType    GoalType       `json:"goal_type"`
	GoalValue   int            `json:"goal_value"`
	Direction   HabitDirection `json:"direction"`
	Color       string         `json:"color"`
	Icon        string         `json:"icon"`
	Schedule    *HabitSchedule `json:"schedule,omitempty"`
	ArchivedAt  *time.Time     `json:"archived_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (h Habit) Validate() error {
	if h.Name == "" {
		return errors.New("habit name is required")
	}
	switch h.GoalType {
	case GoalDuration, GoalCount, GoalBoolean:
	default:
		return errors.New("invalid goal type: must be duration, count, or boolean")
	}
	switch h.Direction {
	case DirectionBuild, DirectionAvoid:
	default:
		return errors.New("invalid direction: must be build or avoid")
	}
	return nil
}
