package core

import (
	"errors"
	"testing"
	"time"
)

// mockStore implements HabitStore for testing.
type mockStore struct {
	// recorded calls
	createHabitCalled    bool
	createHabitArg       Habit
	listHabitsCalled     bool
	listHabitsUserID     string
	logCompletionCalled  bool
	logCompletionArg     CompletionLog
	getStreaksCalled      bool
	getStreaksUserID      string

	// configurable returns
	createHabitErr   error
	listHabitsReturn []Habit
	listHabitsErr    error
	getStreaksReturn  []Streak
	getStreaksErr     error
}

func (m *mockStore) CreateHabit(habit Habit) error {
	m.createHabitCalled = true
	m.createHabitArg = habit
	return m.createHabitErr
}

func (m *mockStore) ListHabits(userID string) ([]Habit, error) {
	m.listHabitsCalled = true
	m.listHabitsUserID = userID
	return m.listHabitsReturn, m.listHabitsErr
}

func (m *mockStore) UpdateHabit(habit Habit) error        { return nil }
func (m *mockStore) ArchiveHabit(id string) error         { return nil }
func (m *mockStore) SetSchedule(schedule HabitSchedule) error { return nil }
func (m *mockStore) GetSchedule(habitID string) (HabitSchedule, error) {
	return HabitSchedule{}, nil
}

func (m *mockStore) LogCompletion(log CompletionLog) error {
	m.logCompletionCalled = true
	m.logCompletionArg = log
	return nil
}

func (m *mockStore) LogSlip(log CompletionLog) error { return nil }

func (m *mockStore) GetLogs(userID string, from, to time.Time) ([]CompletionLog, error) {
	return nil, nil
}

func (m *mockStore) DeleteLog(id string) error { return nil }

func (m *mockStore) CreatePause(pause Pause) error { return nil }

func (m *mockStore) ListPauses(userID string, status string) ([]Pause, error) {
	return nil, nil
}

func (m *mockStore) CancelPause(id string, resumeFrom time.Time) error { return nil }

func (m *mockStore) GetStreaks(userID string) ([]Streak, error) {
	m.getStreaksCalled = true
	m.getStreaksUserID = userID
	return m.getStreaksReturn, m.getStreaksErr
}

func (m *mockStore) GetSummary(userID string, period Period) (Summary, error) {
	return Summary{}, nil
}

func (m *mockStore) GetGoalAlignment(userID string, period Period) (GoalAlignment, error) {
	return GoalAlignment{}, nil
}

func (m *mockStore) GetHeatmap(habitID string, from, to time.Time) ([]HeatmapDay, error) {
	return nil, nil
}

// mockSetScheduleStore extends mockStore to track SetSchedule calls.
type mockSetScheduleStore struct {
	mockStore
	setScheduleCalled bool
	setScheduleArg    HabitSchedule
	setScheduleErr    error
}

func (m *mockSetScheduleStore) SetSchedule(schedule HabitSchedule) error {
	m.setScheduleCalled = true
	m.setScheduleArg = schedule
	return m.setScheduleErr
}

// Tests

func TestCreateHabit_ValidHabitSetsUserIDAndDelegates(t *testing.T) {
	store := &mockStore{}
	svc := NewHabitService(store)

	habit := Habit{
		Name:      "Exercise",
		GoalType:  GoalTypeDuration,
		GoalValue: 30,
		Direction: DirectionBuild,
	}

	err := svc.CreateHabit("user1", habit)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !store.createHabitCalled {
		t.Fatal("expected store.CreateHabit to be called")
	}
	if store.createHabitArg.UserID != "user1" {
		t.Errorf("expected UserID=user1, got %q", store.createHabitArg.UserID)
	}
}

func TestCreateHabit_InvalidHabitRejectsWithoutCallingStore(t *testing.T) {
	store := &mockStore{}
	svc := NewHabitService(store)

	habit := Habit{
		Name:      "", // invalid: empty name
		GoalType:  GoalTypeDuration,
		GoalValue: 30,
		Direction: DirectionBuild,
	}

	err := svc.CreateHabit("user1", habit)
	if err == nil {
		t.Fatal("expected error for invalid habit, got nil")
	}
	if store.createHabitCalled {
		t.Fatal("store.CreateHabit should not be called for invalid habit")
	}
}

func TestSetSchedule_InvalidScheduleReturnsError(t *testing.T) {
	store := &mockSetScheduleStore{}
	svc := NewHabitService(store)

	schedule := HabitSchedule{
		HabitID:      "h1",
		ScheduleType: ScheduleTimesPerWeek,
		TimesPerWeek: 0, // invalid: must be >= 1
	}

	err := svc.SetSchedule("user1", schedule)
	if err == nil {
		t.Fatal("expected error for invalid schedule, got nil")
	}
	if store.setScheduleCalled {
		t.Fatal("store.SetSchedule should not be called for invalid schedule")
	}
}

func TestSetSchedule_ValidScheduleSetsUserIDAndDelegates(t *testing.T) {
	store := &mockSetScheduleStore{}
	svc := NewHabitService(store)

	schedule := HabitSchedule{
		HabitID:      "h1",
		ScheduleType: ScheduleDaily,
	}

	err := svc.SetSchedule("user1", schedule)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !store.setScheduleCalled {
		t.Fatal("expected store.SetSchedule to be called")
	}
	if store.setScheduleArg.UserID != "user1" {
		t.Errorf("expected UserID=user1, got %q", store.setScheduleArg.UserID)
	}
}

func TestLogCompletion_SetsUserIDAndDelegates(t *testing.T) {
	store := &mockStore{}
	svc := NewHabitService(store)

	log := CompletionLog{
		HabitID: "h1",
		LogType: LogCompletion,
		Date:    makeDate(2026, 3, 30),
		Value:   45,
	}

	err := svc.LogCompletion("user1", log)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !store.logCompletionCalled {
		t.Fatal("expected store.LogCompletion to be called")
	}
	if store.logCompletionArg.UserID != "user1" {
		t.Errorf("expected UserID=user1, got %q", store.logCompletionArg.UserID)
	}
}

func TestListHabits_DelegatesToStore(t *testing.T) {
	expected := []Habit{
		newTestHabit("Exercise", DirectionBuild, GoalTypeDuration, nil),
		newTestHabit("No Sugar", DirectionAvoid, GoalTypeBoolean, nil),
	}
	store := &mockStore{
		listHabitsReturn: expected,
	}
	svc := NewHabitService(store)

	result, err := svc.ListHabits("user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !store.listHabitsCalled {
		t.Fatal("expected store.ListHabits to be called")
	}
	if store.listHabitsUserID != "user1" {
		t.Errorf("expected userID=user1, got %q", store.listHabitsUserID)
	}
	if len(result) != len(expected) {
		t.Errorf("expected %d habits, got %d", len(expected), len(result))
	}
}

func TestListHabits_PropagatesStoreError(t *testing.T) {
	storeErr := errors.New("db error")
	store := &mockStore{listHabitsErr: storeErr}
	svc := NewHabitService(store)

	_, err := svc.ListHabits("user1")
	if !errors.Is(err, storeErr) {
		t.Errorf("expected store error to propagate, got %v", err)
	}
}

func TestGetStreaks_DelegatesToStore(t *testing.T) {
	expected := []Streak{
		{HabitID: "h1", HabitName: "Exercise", Current: 5, Longest: 10},
	}
	store := &mockStore{getStreaksReturn: expected}
	svc := NewHabitService(store)

	result, err := svc.GetStreaks("user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !store.getStreaksCalled {
		t.Fatal("expected store.GetStreaks to be called")
	}
	if store.getStreaksUserID != "user1" {
		t.Errorf("expected userID=user1, got %q", store.getStreaksUserID)
	}
	if len(result) != 1 || result[0].HabitID != "h1" {
		t.Errorf("unexpected streaks result: %v", result)
	}
}
