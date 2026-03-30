package core

import (
	"testing"
	"time"
)

func TestLogTypeConstants(t *testing.T) {
	tests := []struct {
		name string
		got  LogType
		want LogType
	}{
		{"LogCompletion equals completion", LogCompletion, "completion"},
		{"LogSlip equals slip", LogSlip, "slip"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestSlipSeverityConstants(t *testing.T) {
	tests := []struct {
		name string
		got  SlipSeverity
		want SlipSeverity
	}{
		{"SlipMinor equals 1", SlipMinor, 1},
		{"SlipFull equals 2", SlipFull, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %d, want %d", tt.got, tt.want)
			}
		})
	}
}

func TestCompletionLogInitialization(t *testing.T) {
	now := time.Now()
	log := CompletionLog{
		ID:        "log-1",
		HabitID:   "habit-1",
		UserID:    "user-1",
		LogType:   LogCompletion,
		Date:      now,
		Value:     30,
		Note:      "felt great",
		CreatedAt: now,
	}

	if log.ID != "log-1" {
		t.Errorf("ID: got %q, want %q", log.ID, "log-1")
	}
	if log.HabitID != "habit-1" {
		t.Errorf("HabitID: got %q, want %q", log.HabitID, "habit-1")
	}
	if log.UserID != "user-1" {
		t.Errorf("UserID: got %q, want %q", log.UserID, "user-1")
	}
	if log.LogType != LogCompletion {
		t.Errorf("LogType: got %q, want %q", log.LogType, LogCompletion)
	}
	if log.Value != 30 {
		t.Errorf("Value: got %d, want %d", log.Value, 30)
	}
	if log.Note != "felt great" {
		t.Errorf("Note: got %q, want %q", log.Note, "felt great")
	}
	if !log.Date.Equal(now) {
		t.Errorf("Date: got %v, want %v", log.Date, now)
	}
	if !log.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt: got %v, want %v", log.CreatedAt, now)
	}
}

func TestLogTypesForHabitDirections(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		log         CompletionLog
		wantLogType LogType
	}{
		{
			name: "completion log for build habit uses LogCompletion",
			log: CompletionLog{
				ID:      "log-build",
				HabitID: "habit-running",
				UserID:  "user-1",
				LogType: LogCompletion,
				Date:    now,
				Value:   45,
			},
			wantLogType: LogCompletion,
		},
		{
			name: "slip log for avoid habit uses LogSlip with severity as Value",
			log: CompletionLog{
				ID:      "log-avoid",
				HabitID: "habit-no-social-media",
				UserID:  "user-1",
				LogType: LogSlip,
				Date:    now,
				Value:   int(SlipFull),
			},
			wantLogType: LogSlip,
		},
		{
			name: "minor slip for avoid habit uses SlipMinor severity",
			log: CompletionLog{
				ID:      "log-minor-slip",
				HabitID: "habit-no-junk-food",
				UserID:  "user-1",
				LogType: LogSlip,
				Date:    now,
				Value:   int(SlipMinor),
			},
			wantLogType: LogSlip,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.log.LogType != tt.wantLogType {
				t.Errorf("LogType: got %q, want %q", tt.log.LogType, tt.wantLogType)
			}
		})
	}

	// Verify slip severity is stored as Value
	slipLog := CompletionLog{
		LogType: LogSlip,
		Value:   int(SlipFull),
	}
	if SlipSeverity(slipLog.Value) != SlipFull {
		t.Errorf("slip Value as SlipSeverity: got %d, want %d", slipLog.Value, SlipFull)
	}

	minorSlipLog := CompletionLog{
		LogType: LogSlip,
		Value:   int(SlipMinor),
	}
	if SlipSeverity(minorSlipLog.Value) != SlipMinor {
		t.Errorf("minor slip Value as SlipSeverity: got %d, want %d", minorSlipLog.Value, SlipMinor)
	}
}
