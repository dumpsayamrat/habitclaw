package core

import "testing"

func TestCalculateStreak(t *testing.T) {
	dailySchedule := &HabitSchedule{ScheduleType: ScheduleDaily}
	// Mon/Wed/Fri schedule
	mwfSchedule := &HabitSchedule{ScheduleType: ScheduleSpecificDays, DaysOfWeek: []int{1, 3, 5}}

	tests := []struct {
		name        string
		habit       Habit
		logs        []CompletionLog
		pauses      []Pause
		today       func() string // returns date string for readability
		todayTime   func() interface{}
		wantCurrent int
		wantLongest int
	}{
		{
			name:  "build habit: consecutive completions increment streak",
			habit: newTestHabit("Running", DirectionBuild, GoalDuration, dailySchedule),
			logs: func() []CompletionLog {
				// 5 consecutive days of completions ending today (2026-03-30)
				habitID := "test-Running"
				return []CompletionLog{
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 26), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 27), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 28), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 29), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 30), 30),
				}
			}(),
			pauses:      nil,
			wantCurrent: 5,
			wantLongest: 5,
		},
		{
			name:  "build habit: missed day resets streak",
			habit: newTestHabit("Running", DirectionBuild, GoalDuration, dailySchedule),
			logs: func() []CompletionLog {
				// Completed 26-27, missed 28, completed 29-30
				habitID := "test-Running"
				return []CompletionLog{
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 26), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 27), 30),
					// gap on March 28
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 29), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 30), 30),
				}
			}(),
			pauses:      nil,
			wantCurrent: 2, // only 29-30
			wantLongest: 2, // both segments are length 2
		},
		{
			name:  "build habit: paused days do NOT break streak",
			habit: newTestHabit("Running", DirectionBuild, GoalDuration, dailySchedule),
			logs: func() []CompletionLog {
				// Completed 26-27, paused 28, completed 29-30
				habitID := "test-Running"
				return []CompletionLog{
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 26), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 27), 30),
					// paused on March 28
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 29), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 30), 30),
				}
			}(),
			pauses: []Pause{
				newTestPause(nil, makeDate(2026, 3, 28), makeDate(2026, 3, 28)),
			},
			wantCurrent: 4, // 26, 27, 29, 30 — pause skipped
			wantLongest: 4,
		},
		{
			name:  "build habit: non-scheduled days are ignored",
			habit: newTestHabit("Running", DirectionBuild, GoalDuration, mwfSchedule),
			logs: func() []CompletionLog {
				// 2026-03-30 is Monday (1), 2026-03-27 is Friday (5), 2026-03-25 is Wednesday (3)
				habitID := "test-Running"
				return []CompletionLog{
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 25), 30), // Wed
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 27), 30), // Fri
					// Sat 28, Sun 29 are not scheduled — should be skipped
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 30), 30), // Mon
				}
			}(),
			pauses:      nil,
			wantCurrent: 3, // Wed, Fri, Mon — non-scheduled days ignored
			wantLongest: 3,
		},
		{
			name:  "avoid habit: clean days increment streak",
			habit: newTestHabit("NoSocialMedia", DirectionAvoid, GoalBoolean, dailySchedule),
			logs:  nil, // no slips = clean days
			pauses: nil,
			wantCurrent: 367, // walks back 365+ days with no slips on daily schedule (inclusive of today and start)
			wantLongest: 367,
		},
		{
			name:  "avoid habit: minor slip (severity=1) does NOT break streak",
			habit: newTestHabit("NoSocialMedia", DirectionAvoid, GoalBoolean, dailySchedule),
			logs: func() []CompletionLog {
				habitID := "test-NoSocialMedia"
				return []CompletionLog{
					// Minor slip on March 28 — streak should continue through it
					newTestLog(habitID, LogSlip, makeDate(2026, 3, 28), int(SlipMinor)),
				}
			}(),
			pauses:      nil,
			wantCurrent: 367, // minor slip does not break streak
			wantLongest: 367,
		},
		{
			name:  "avoid habit: full slip (severity=2) breaks streak",
			habit: newTestHabit("NoSocialMedia", DirectionAvoid, GoalBoolean, dailySchedule),
			logs: func() []CompletionLog {
				habitID := "test-NoSocialMedia"
				return []CompletionLog{
					// Full slip on March 28
					newTestLog(habitID, LogSlip, makeDate(2026, 3, 28), int(SlipFull)),
				}
			}(),
			pauses:      nil,
			wantCurrent: 2,   // only March 29-30 (after the full slip)
			wantLongest: 363, // long clean period before the full slip on March 28
		},
		{
			name:  "avoid habit: paused days do NOT break streak",
			habit: newTestHabit("NoSocialMedia", DirectionAvoid, GoalBoolean, dailySchedule),
			logs: func() []CompletionLog {
				habitID := "test-NoSocialMedia"
				return []CompletionLog{
					// Full slip on March 25
					newTestLog(habitID, LogSlip, makeDate(2026, 3, 25), int(SlipFull)),
				}
			}(),
			pauses: []Pause{
				// Pause on March 27 — should not break the streak between 26 and 28
				newTestPause(nil, makeDate(2026, 3, 27), makeDate(2026, 3, 27)),
			},
			wantCurrent: 4,   // March 26, 28, 29, 30 (27 paused, skipped)
			wantLongest: 360, // long clean period before the slip on March 25 (year walk minus paused day and days after slip start)
		},
	}

	today := makeDate(2026, 3, 30)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateStreak(tt.habit, tt.logs, tt.pauses, today)

			if result.Current != tt.wantCurrent {
				t.Errorf("current streak: got %d, want %d", result.Current, tt.wantCurrent)
			}
			if result.Longest != tt.wantLongest {
				t.Errorf("longest streak: got %d, want %d", result.Longest, tt.wantLongest)
			}
			if result.HabitID != tt.habit.ID {
				t.Errorf("habit ID: got %q, want %q", result.HabitID, tt.habit.ID)
			}
			if result.Direction != string(tt.habit.Direction) {
				t.Errorf("direction: got %q, want %q", result.Direction, string(tt.habit.Direction))
			}
		})
	}
}

func TestCalculateStreak_NoSchedule(t *testing.T) {
	// A habit with no schedule should return a zero streak
	habit := newTestHabit("NoSchedule", DirectionBuild, GoalDuration, nil)
	result := CalculateStreak(habit, nil, nil, makeDate(2026, 3, 30))

	if result.Current != 0 {
		t.Errorf("expected current streak 0 for habit with no schedule, got %d", result.Current)
	}
	if result.Longest != 0 {
		t.Errorf("expected longest streak 0 for habit with no schedule, got %d", result.Longest)
	}
}
