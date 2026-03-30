package core

import "time"

func makeDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func newTestHabit(name string, direction HabitDirection, goalType GoalType, schedule *HabitSchedule) Habit {
	return Habit{
		ID:        "test-" + name,
		UserID:    "local",
		Name:      name,
		GoalType:  goalType,
		GoalValue: 30,
		Direction: direction,
		Schedule:  schedule,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func newTestLog(habitID string, logType LogType, date time.Time, value int) CompletionLog {
	return CompletionLog{
		ID:        "log-" + date.Format("2006-01-02"),
		HabitID:   habitID,
		UserID:    "local",
		LogType:   logType,
		Date:      date,
		Value:     value,
		CreatedAt: time.Now(),
	}
}

func newTestPause(habitID *string, from, to time.Time) Pause {
	return Pause{
		ID:        "pause-test",
		UserID:    "local",
		HabitID:   habitID,
		FromDate:  from,
		ToDate:    to,
		Reason:    "test",
		CreatedAt: time.Now(),
	}
}
