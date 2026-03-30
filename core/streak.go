package core

import "time"

type Streak struct {
	HabitID            string  `json:"habit_id"`
	HabitName          string  `json:"habit_name"`
	Direction          string  `json:"direction"`
	Current            int     `json:"current"`
	Longest            int     `json:"longest"`
	ConsistencyRate7d  float64 `json:"consistency_rate_7d"`
	ConsistencyRate30d float64 `json:"consistency_rate_30d"`
	LastActivityDate   string  `json:"last_activity_date"`
}

// CalculateStreak computes the current and longest streak for a habit
// given its logs, schedule, and pauses. Streaks are always computed, never cached.
func CalculateStreak(habit Habit, logs []CompletionLog, pauses []Pause, today time.Time) Streak {
	schedule := habit.Schedule
	if schedule == nil {
		return Streak{
			HabitID:   habit.ID,
			HabitName: habit.Name,
			Direction: string(habit.Direction),
		}
	}

	pausedDates := make(map[string]bool)
	for _, p := range pauses {
		if p.CancelledAt != nil && p.ResumeFrom != nil {
			// Cancelled early — only paused up to resume date
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

	completionDates := make(map[string]bool)
	slipDates := make(map[string]int) // date -> max severity
	lastActivity := ""
	for _, l := range logs {
		dateStr := l.Date.Format("2006-01-02")
		if l.LogType == LogCompletion {
			completionDates[dateStr] = true
		} else if l.LogType == LogSlip {
			if sev, ok := slipDates[dateStr]; !ok || l.Value > sev {
				slipDates[dateStr] = l.Value
			}
		}
		if dateStr > lastActivity {
			lastActivity = dateStr
		}
	}

	current := 0
	longest := 0
	streakBroken := false

	for d := today; ; d = d.AddDate(0, 0, -1) {
		dateStr := d.Format("2006-01-02")

		if pausedDates[dateStr] {
			continue
		}

		if !schedule.IsScheduledDay(d, habit.CreatedAt) {
			// Don't go back more than 365 days
			if today.Sub(d).Hours() > 365*24 {
				break
			}
			continue
		}

		if habit.Direction == DirectionBuild {
			if completionDates[dateStr] {
				if !streakBroken {
					current++
				}
			} else {
				if !streakBroken {
					streakBroken = true
				}
			}
		} else { // avoid
			severity := slipDates[dateStr]
			if severity >= int(SlipFull) {
				if !streakBroken {
					streakBroken = true
				}
			} else {
				// Clean day or minor slip — streak continues
				if !streakBroken {
					current++
				}
			}
		}

		if today.Sub(d).Hours() > 365*24 {
			break
		}
	}

	tempStreak := 0
	startDate := today.AddDate(-1, 0, 0)
	for d := startDate; !d.After(today); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		if pausedDates[dateStr] {
			continue
		}
		if !schedule.IsScheduledDay(d, habit.CreatedAt) {
			continue
		}

		if habit.Direction == DirectionBuild {
			if completionDates[dateStr] {
				tempStreak++
			} else {
				if tempStreak > longest {
					longest = tempStreak
				}
				tempStreak = 0
			}
		} else {
			severity := slipDates[dateStr]
			if severity >= int(SlipFull) {
				if tempStreak > longest {
					longest = tempStreak
				}
				tempStreak = 0
			} else {
				tempStreak++
			}
		}
	}
	if tempStreak > longest {
		longest = tempStreak
	}
	if current > longest {
		longest = current
	}

	return Streak{
		HabitID:            habit.ID,
		HabitName:          habit.Name,
		Direction:          string(habit.Direction),
		Current:            current,
		Longest:            longest,
		ConsistencyRate7d:  calculateConsistencyRate(habit, completionDates, slipDates, pausedDates, today, 7),
		ConsistencyRate30d: calculateConsistencyRate(habit, completionDates, slipDates, pausedDates, today, 30),
		LastActivityDate:   lastActivity,
	}
}

// calculateConsistencyRate returns completed_days / scheduled_days over the last `days` calendar days.
// Paused days are excluded from both numerator and denominator.
// For build habits: completed = days with a completion log on a scheduled day.
// For avoid habits: completed = scheduled days without a full slip (severity=2).
// Returns 0.0 if there are no scheduled non-paused days in the period.
func calculateConsistencyRate(habit Habit, completionDates map[string]bool, slipDates map[string]int, pausedDates map[string]bool, today time.Time, days int) float64 {
	if habit.Schedule == nil {
		return 0.0
	}
	denominator := 0
	numerator := 0
	for i := 0; i < days; i++ {
		d := today.AddDate(0, 0, -i)
		dateStr := d.Format("2006-01-02")
		if pausedDates[dateStr] {
			continue
		}
		if !habit.Schedule.IsScheduledDay(d, habit.CreatedAt) {
			continue
		}
		denominator++
		if habit.Direction == DirectionBuild {
			if completionDates[dateStr] {
				numerator++
			}
		} else {
			severity := slipDates[dateStr]
			if severity < int(SlipFull) {
				numerator++
			}
		}
	}
	if denominator == 0 {
		return 0.0
	}
	return float64(numerator) / float64(denominator)
}
