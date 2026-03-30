package core

import "time"

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

// CalculateGoalAlignment computes alignment scores for build habits with duration/count goals.
func CalculateGoalAlignment(habits []Habit, logs []CompletionLog, pauses []Pause, period Period, today time.Time) GoalAlignment {
	var days int
	switch period {
	case PeriodToday:
		days = 1
	case PeriodWeek:
		days = 7
	case PeriodMonth:
		days = 30
	default:
		days = 7
	}

	allPausedDates := make(map[string]bool)
	habitPausedDates := make(map[string]map[string]bool)

	for _, p := range pauses {
		var pausedSet map[string]bool
		if p.HabitID == nil {
			pausedSet = allPausedDates
		} else {
			if habitPausedDates[*p.HabitID] == nil {
				habitPausedDates[*p.HabitID] = make(map[string]bool)
			}
			pausedSet = habitPausedDates[*p.HabitID]
		}

		if p.CancelledAt != nil && p.ResumeFrom != nil {
			for d := p.FromDate; !d.After(p.ResumeFrom.AddDate(0, 0, -1)); d = d.AddDate(0, 0, 1) {
				pausedSet[d.Format("2006-01-02")] = true
			}
		} else if p.CancelledAt != nil {
			continue
		} else {
			for d := p.FromDate; !d.After(p.ToDate); d = d.AddDate(0, 0, 1) {
				pausedSet[d.Format("2006-01-02")] = true
			}
		}
	}

	completionsByHabit := make(map[string]map[string]int)
	for _, l := range logs {
		if l.LogType != LogCompletion {
			continue
		}
		if completionsByHabit[l.HabitID] == nil {
			completionsByHabit[l.HabitID] = make(map[string]int)
		}
		completionsByHabit[l.HabitID][l.Date.Format("2006-01-02")] = l.Value
	}

	habitAlignments := make([]HabitAlignment, 0)
	totalAlignment := 0.0

	for _, h := range habits {
		if h.Direction != DirectionBuild {
			continue
		}
		if h.GoalType == GoalTypeBoolean || h.GoalValue == 0 {
			continue
		}

		isPaused := func(dateStr string) bool {
			return allPausedDates[dateStr] || habitPausedDates[h.ID][dateStr]
		}

		completions := completionsByHabit[h.ID]
		if completions == nil {
			completions = make(map[string]int)
		}

		sumValues := 0
		completedDays := 0

		for i := 0; i < days; i++ {
			d := today.AddDate(0, 0, -(days - 1 - i))
			dateStr := d.Format("2006-01-02")

			if isPaused(dateStr) {
				continue
			}
			if h.Schedule == nil || !h.Schedule.IsScheduledDay(d, h.CreatedAt) {
				continue
			}

			if val, ok := completions[dateStr]; ok {
				sumValues += val
				completedDays++
			}
		}

		var avgActual float64
		if completedDays > 0 {
			avgActual = float64(sumValues) / float64(completedDays)
		}

		ratio := avgActual / float64(h.GoalValue)
		alignment := ratio
		if alignment > 1.0 {
			alignment = 1.0
		}

		var status string
		switch {
		case ratio > 1.0:
			status = "exceeding"
		case ratio >= 0.8:
			status = "on track"
		case ratio >= 0.5:
			status = "needs attention"
		default:
			status = "off track"
		}

		habitAlignments = append(habitAlignments, HabitAlignment{
			Name:          h.Name,
			GoalValue:     h.GoalValue,
			AverageActual: avgActual,
			Alignment:     alignment,
			Status:        status,
		})
		totalAlignment += alignment
	}

	overallScore := 0.0
	if len(habitAlignments) > 0 {
		overallScore = totalAlignment / float64(len(habitAlignments))
	}

	return GoalAlignment{
		Period:       string(period),
		OverallScore: overallScore,
		Habits:       habitAlignments,
	}
}
