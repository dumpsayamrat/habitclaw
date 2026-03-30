package core

import (
	"testing"
	"time"
)

func TestPauseStruct(t *testing.T) {
	now := time.Now()
	habitID := "habit-123"

	t.Run("all fields initialized correctly", func(t *testing.T) {
		p := Pause{
			ID:        "pause-1",
			UserID:    "user-1",
			HabitID:   &habitID,
			FromDate:  now,
			ToDate:    now.AddDate(0, 0, 7),
			Reason:    "vacation",
			CreatedAt: now,
		}
		if p.ID != "pause-1" {
			t.Errorf("expected ID %q, got %q", "pause-1", p.ID)
		}
		if p.UserID != "user-1" {
			t.Errorf("expected UserID %q, got %q", "user-1", p.UserID)
		}
		if p.HabitID == nil || *p.HabitID != habitID {
			t.Errorf("expected HabitID %q, got %v", habitID, p.HabitID)
		}
		if p.Reason != "vacation" {
			t.Errorf("expected Reason %q, got %q", "vacation", p.Reason)
		}
	})

	t.Run("nil HabitID means global pause (all habits)", func(t *testing.T) {
		p := Pause{
			ID:       "pause-global",
			UserID:   "user-1",
			HabitID:  nil,
			FromDate: now,
			ToDate:   now.AddDate(0, 0, 3),
			Reason:   "sick",
		}
		if p.HabitID != nil {
			t.Errorf("expected HabitID nil for global pause, got %v", p.HabitID)
		}
	})

	t.Run("non-nil HabitID means specific habit paused", func(t *testing.T) {
		p := Pause{
			ID:      "pause-specific",
			UserID:  "user-1",
			HabitID: &habitID,
			Reason:  "travel",
		}
		if p.HabitID == nil {
			t.Error("expected HabitID to be non-nil for specific habit pause")
		}
		if *p.HabitID != habitID {
			t.Errorf("expected HabitID %q, got %q", habitID, *p.HabitID)
		}
	})

	t.Run("CancelledAt and ResumeFrom are nil by default", func(t *testing.T) {
		p := Pause{
			ID:       "pause-2",
			UserID:   "user-1",
			FromDate: now,
			ToDate:   now.AddDate(0, 0, 5),
			Reason:   "other",
		}
		if p.CancelledAt != nil {
			t.Errorf("expected CancelledAt nil, got %v", p.CancelledAt)
		}
		if p.ResumeFrom != nil {
			t.Errorf("expected ResumeFrom nil, got %v", p.ResumeFrom)
		}
	})

	t.Run("cancelled pause has CancelledAt and ResumeFrom set", func(t *testing.T) {
		cancelledAt := now.Add(24 * time.Hour)
		resumeFrom := now.Add(48 * time.Hour)
		p := Pause{
			ID:          "pause-3",
			UserID:      "user-1",
			FromDate:    now,
			ToDate:      now.AddDate(0, 0, 7),
			Reason:      "vacation",
			CancelledAt: &cancelledAt,
			ResumeFrom:  &resumeFrom,
		}
		if p.CancelledAt == nil {
			t.Error("expected CancelledAt to be set on cancelled pause")
		}
		if p.ResumeFrom == nil {
			t.Error("expected ResumeFrom to be set on cancelled pause")
		}
		if !p.CancelledAt.Equal(cancelledAt) {
			t.Errorf("expected CancelledAt %v, got %v", cancelledAt, *p.CancelledAt)
		}
		if !p.ResumeFrom.Equal(resumeFrom) {
			t.Errorf("expected ResumeFrom %v, got %v", resumeFrom, *p.ResumeFrom)
		}
	})

	t.Run("valid reason values", func(t *testing.T) {
		validReasons := []string{"vacation", "sick", "travel", "other"}
		for _, reason := range validReasons {
			p := Pause{
				ID:       "pause-reason",
				UserID:   "user-1",
				FromDate: now,
				ToDate:   now.AddDate(0, 0, 1),
				Reason:   reason,
			}
			if p.Reason != reason {
				t.Errorf("expected Reason %q, got %q", reason, p.Reason)
			}
		}
	})
}
