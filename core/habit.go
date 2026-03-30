package core

import "time"

// GoalType defines how a habit completion is measured
type GoalType string

const (
	GoalTypeDuration GoalType = "duration"
	GoalTypeCount    GoalType = "count"
	GoalTypeBoolean  GoalType = "boolean"
)

// HabitDirection defines whether this is a build or avoid habit
type HabitDirection string

const (
	DirectionBuild HabitDirection = "build"
	DirectionAvoid HabitDirection = "avoid"
)

// Habit represents a single tracked habit
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

// IsValid returns an error string if the habit is invalid, empty string if valid
func (h Habit) IsValid() string {
	if h.Name == "" {
		return "habit name is required"
	}
	switch h.GoalType {
	case GoalTypeDuration, GoalTypeCount, GoalTypeBoolean:
	default:
		return "invalid goal type: must be duration, count, or boolean"
	}
	switch h.Direction {
	case DirectionBuild, DirectionAvoid:
	default:
		return "invalid direction: must be build or avoid"
	}
	if h.GoalValue < 0 {
		return "goal value must be >= 0"
	}
	if h.GoalType == GoalTypeBoolean && h.GoalValue != 0 {
		return "goal value must be 0 for boolean goal type"
	}
	return ""
}
