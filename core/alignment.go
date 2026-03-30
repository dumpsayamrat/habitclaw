package core

type GoalAlignment struct {
	Period       string           `json:"period"`
	OverallScore float64          `json:"overall_score"`
	Habits       []HabitAlignment `json:"habits"`
}

type HabitAlignment struct {
	Name          string  `json:"name"`
	GoalValue     int     `json:"goal_value"`
	AverageActual float64 `json:"average_actual"`
	Alignment     float64 `json:"alignment"`
	Status        string  `json:"status"`
}
