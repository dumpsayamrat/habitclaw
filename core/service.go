package core

import (
	"fmt"
	"time"
)

type HabitService struct {
	store HabitStore
}

func NewHabitService(store HabitStore) *HabitService {
	return &HabitService{store: store}
}

func (s *HabitService) CreateHabit(userID string, habit Habit) error {
	habit.UserID = userID
	if msg := habit.IsValid(); msg != "" {
		return fmt.Errorf("invalid habit: %s", msg)
	}
	return s.store.CreateHabit(habit)
}

func (s *HabitService) ListHabits(userID string) ([]Habit, error) {
	return s.store.ListHabits(userID)
}

func (s *HabitService) UpdateHabit(userID string, habit Habit) error {
	habit.UserID = userID
	if msg := habit.IsValid(); msg != "" {
		return fmt.Errorf("invalid habit: %s", msg)
	}
	return s.store.UpdateHabit(habit)
}

func (s *HabitService) ArchiveHabit(userID string, id string) error {
	return s.store.ArchiveHabit(id)
}

func (s *HabitService) SetSchedule(userID string, schedule HabitSchedule) error {
	schedule.UserID = userID
	if msg := schedule.IsValid(); msg != "" {
		return fmt.Errorf("invalid schedule: %s", msg)
	}
	return s.store.SetSchedule(schedule)
}

func (s *HabitService) GetSchedule(userID string, habitID string) (HabitSchedule, error) {
	return s.store.GetSchedule(habitID)
}

func (s *HabitService) LogCompletion(userID string, log CompletionLog) error {
	log.UserID = userID
	return s.store.LogCompletion(log)
}

func (s *HabitService) LogSlip(userID string, log CompletionLog) error {
	log.UserID = userID
	return s.store.LogSlip(log)
}

func (s *HabitService) GetLogs(userID string, from, to time.Time) ([]CompletionLog, error) {
	return s.store.GetLogs(userID, from, to)
}

func (s *HabitService) DeleteLog(userID string, id string) error {
	return s.store.DeleteLog(id)
}

func (s *HabitService) CreatePause(userID string, pause Pause) error {
	pause.UserID = userID
	return s.store.CreatePause(pause)
}

func (s *HabitService) ListPauses(userID string, status string) ([]Pause, error) {
	return s.store.ListPauses(userID, status)
}

func (s *HabitService) CancelPause(userID string, id string, resumeFrom time.Time) error {
	return s.store.CancelPause(id, resumeFrom)
}

func (s *HabitService) GetStreaks(userID string) ([]Streak, error) {
	return s.store.GetStreaks(userID)
}

func (s *HabitService) GetSummary(userID string, period Period) (Summary, error) {
	return s.store.GetSummary(userID, period)
}

func (s *HabitService) GetGoalAlignment(userID string, period Period) (GoalAlignment, error) {
	return s.store.GetGoalAlignment(userID, period)
}

func (s *HabitService) GetHeatmap(userID string, habitID string, from, to time.Time) ([]HeatmapDay, error) {
	return s.store.GetHeatmap(habitID, from, to)
}
