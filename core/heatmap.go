package core

import "time"

// CalculateHeatmap returns a HeatmapDay for each calendar day from `from` to `to` inclusive.
// Paused days (global pauses where HabitID==nil, and habit-specific pauses matching habit.ID)
// are marked IsPaused=true. Completion and slip logs are reflected per day.
func CalculateHeatmap(habit Habit, logs []CompletionLog, pauses []Pause, from, to time.Time) []HeatmapDay {
	// Build paused dates set
	pausedDates := make(map[string]bool)
	for _, p := range pauses {
		// Only include global pauses or pauses for this habit
		if p.HabitID != nil && *p.HabitID != habit.ID {
			continue
		}
		if p.CancelledAt != nil && p.ResumeFrom != nil {
			// Cancelled early — only paused up to resume date (exclusive)
			for d := p.FromDate; !d.After(p.ResumeFrom.AddDate(0, 0, -1)); d = d.AddDate(0, 0, 1) {
				pausedDates[d.Format("2006-01-02")] = true
			}
		} else if p.CancelledAt != nil {
			continue // cancelled without resume = ignore
		} else {
			for d := p.FromDate; !d.After(p.ToDate); d = d.AddDate(0, 0, 1) {
				pausedDates[d.Format("2006-01-02")] = true
			}
		}
	}

	// Build map of date -> log for this habit's logs
	logsByDate := make(map[string]CompletionLog)
	for _, l := range logs {
		if l.HabitID != habit.ID {
			continue
		}
		dateStr := l.Date.Format("2006-01-02")
		// Keep the first log per date (or prefer completion over slip if both exist)
		if existing, ok := logsByDate[dateStr]; !ok {
			logsByDate[dateStr] = l
		} else if existing.LogType == LogSlip && l.LogType == LogCompletion {
			logsByDate[dateStr] = l
		}
	}

	var result []HeatmapDay
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		if pausedDates[dateStr] {
			result = append(result, HeatmapDay{Date: dateStr, IsPaused: true})
			continue
		}

		log, hasLog := logsByDate[dateStr]
		if hasLog && log.LogType == LogCompletion {
			goalMet := habit.GoalType == GoalTypeBoolean || log.Value >= habit.GoalValue
			result = append(result, HeatmapDay{
				Date:    dateStr,
				Value:   log.Value,
				GoalMet: goalMet,
			})
		} else if hasLog && log.LogType == LogSlip {
			result = append(result, HeatmapDay{
				Date:     dateStr,
				IsSlip:   true,
				Severity: log.Value,
			})
		} else {
			result = append(result, HeatmapDay{Date: dateStr})
		}
	}

	return result
}
