package core

import "time"

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

// CalculateSummary computes a Summary for all habits over the given period.
func CalculateSummary(habits []Habit, logs []CompletionLog, pauses []Pause, period Period, today time.Time) Summary {
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
			// Cancelled early — paused up to (but not including) resume date
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

	type logMaps struct {
		completions map[string]int // date -> value
		slips       map[string]int // date -> max severity
	}
	logsByHabit := make(map[string]*logMaps)
	for _, l := range logs {
		lm := logsByHabit[l.HabitID]
		if lm == nil {
			lm = &logMaps{
				completions: make(map[string]int),
				slips:       make(map[string]int),
			}
			logsByHabit[l.HabitID] = lm
		}
		dateStr := l.Date.Format("2006-01-02")
		if l.LogType == LogCompletion {
			lm.completions[dateStr] = l.Value
		} else if l.LogType == LogSlip {
			if sev, ok := lm.slips[dateStr]; !ok || l.Value > sev {
				lm.slips[dateStr] = l.Value
			}
		}
	}

	totalScheduled := 0
	totalCompleted := 0
	habitStats := make([]HabitStat, 0, len(habits))

	for _, h := range habits {
		lm := logsByHabit[h.ID]
		if lm == nil {
			lm = &logMaps{
				completions: make(map[string]int),
				slips:       make(map[string]int),
			}
		}

		isPaused := func(dateStr string) bool {
			return allPausedDates[dateStr] || habitPausedDates[h.ID][dateStr]
		}

		stat := HabitStat{
			HabitID:   h.ID,
			Name:      h.Name,
			Direction: string(h.Direction),
		}

		goalMetDays := 0

		for i := 0; i < days; i++ {
			d := today.AddDate(0, 0, -(days - 1 - i))
			dateStr := d.Format("2006-01-02")

			if isPaused(dateStr) {
				stat.PausedDays++
				continue
			}

			if h.Schedule == nil || !h.Schedule.IsScheduledDay(d, h.CreatedAt) {
				continue
			}

			stat.ScheduledDays++

			if h.Direction == DirectionBuild {
				if val, ok := lm.completions[dateStr]; ok {
					stat.CompletedDays++
					stat.TotalValue += val
					if h.GoalType == GoalTypeBoolean || val >= h.GoalValue {
						goalMetDays++
					}
				}
			} else { // avoid
				severity := lm.slips[dateStr]
				if severity >= int(SlipFull) {
					stat.SlippedDays++
				} else {
					stat.CompletedDays++
					goalMetDays++
				}
			}
		}

		if stat.ScheduledDays > 0 {
			stat.GoalMetRate = float64(goalMetDays) / float64(stat.ScheduledDays)
		}

		streak := CalculateStreak(h, logs, pauses, today)
		stat.CurrentStreak = streak.Current

		totalScheduled += stat.ScheduledDays
		totalCompleted += stat.CompletedDays
		habitStats = append(habitStats, stat)
	}

	completionRate := 0.0
	if totalScheduled > 0 {
		completionRate = float64(totalCompleted) / float64(totalScheduled)
	}

	return Summary{
		Period:         string(period),
		TotalHabits:    len(habits),
		ScheduledDays:  totalScheduled,
		CompletionRate: completionRate,
		HabitSummaries: habitStats,
	}
}
