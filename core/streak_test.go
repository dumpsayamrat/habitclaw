package core

import "testing"

func TestCalculateStreak(t *testing.T) {
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
			habit: newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule()),
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
			habit: newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule()),
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
			habit: newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule()),
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
			habit: newTestHabit("Running", DirectionBuild, GoalTypeDuration, mwfSchedule),
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
			habit: newTestHabit("NoSocialMedia", DirectionAvoid, GoalTypeBoolean, dailySchedule()),
			logs:  nil, // no slips = clean days
			pauses: nil,
			wantCurrent: 367, // walks back 365+ days with no slips on daily schedule (inclusive of today and start)
			wantLongest: 367,
		},
		{
			name:  "avoid habit: minor slip (severity=1) does NOT break streak",
			habit: newTestHabit("NoSocialMedia", DirectionAvoid, GoalTypeBoolean, dailySchedule()),
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
			habit: newTestHabit("NoSocialMedia", DirectionAvoid, GoalTypeBoolean, dailySchedule()),
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
			habit: newTestHabit("NoSocialMedia", DirectionAvoid, GoalTypeBoolean, dailySchedule()),
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

func TestCalculateStreak_ConsistencyRate(t *testing.T) {
	today := makeDate(2026, 3, 30)

	tests := []struct {
		name    string
		habit   Habit
		logs    []CompletionLog
		pauses  []Pause
		want7d  float64
		want30d float64
	}{
		{
			name:  "build habit: 5 of 7 completions = ~0.714 consistency_7d",
			habit: newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule()),
			logs: func() []CompletionLog {
				habitID := "test-Running"
				// today=Mar30, window=Mar24..Mar30 (7 days)
				// complete: 24,25,26,28,30 = 5 of 7
				return []CompletionLog{
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 24), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 25), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 26), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 28), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 30), 30),
				}
			}(),
			pauses:  nil,
			want7d:  5.0 / 7.0,
			want30d: 5.0 / 30.0,
		},
		{
			name:  "avoid habit: 1 full slip in 7 days = 6/7 consistency_7d",
			habit: newTestHabit("NoSocialMedia", DirectionAvoid, GoalTypeBoolean, dailySchedule()),
			logs: func() []CompletionLog {
				habitID := "test-NoSocialMedia"
				return []CompletionLog{
					// Full slip on March 27 — only day that breaks consistency
					newTestLog(habitID, LogSlip, makeDate(2026, 3, 27), int(SlipFull)),
				}
			}(),
			pauses:  nil,
			want7d:  6.0 / 7.0,
			want30d: 29.0 / 30.0,
		},
		{
			name:  "avoid habit: minor slip does NOT count against consistency",
			habit: newTestHabit("NoSocialMedia", DirectionAvoid, GoalTypeBoolean, dailySchedule()),
			logs: func() []CompletionLog {
				habitID := "test-NoSocialMedia"
				return []CompletionLog{
					newTestLog(habitID, LogSlip, makeDate(2026, 3, 28), int(SlipMinor)),
				}
			}(),
			pauses:  nil,
			want7d:  1.0,
			want30d: 1.0,
		},
		{
			name:  "paused days excluded from denominator",
			habit: newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule()),
			logs: func() []CompletionLog {
				habitID := "test-Running"
				// complete all non-paused days in 7-day window
				// Mar24..Mar30, pause Mar27 → 6 scheduled days, 6 completions
				return []CompletionLog{
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 24), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 25), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 26), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 28), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 29), 30),
					newTestLog(habitID, LogCompletion, makeDate(2026, 3, 30), 30),
				}
			}(),
			pauses: []Pause{
				newTestPause(nil, makeDate(2026, 3, 27), makeDate(2026, 3, 27)),
			},
			want7d:  1.0, // 6/6 — paused day excluded
			want30d: 6.0 / 29.0,
		},
		{
			name:  "all days completed = 1.0",
			habit: newTestHabit("Running", DirectionBuild, GoalTypeDuration, dailySchedule()),
			logs: func() []CompletionLog {
				habitID := "test-Running"
				var logs []CompletionLog
				for i := 0; i < 30; i++ {
					logs = append(logs, newTestLog(habitID, LogCompletion, today.AddDate(0, 0, -i), 30))
				}
				return logs
			}(),
			pauses:  nil,
			want7d:  1.0,
			want30d: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateStreak(tt.habit, tt.logs, tt.pauses, today)

			const eps = 1e-9
			if diff := result.ConsistencyRate7d - tt.want7d; diff > eps || diff < -eps {
				t.Errorf("ConsistencyRate7d: got %f, want %f", result.ConsistencyRate7d, tt.want7d)
			}
			if diff := result.ConsistencyRate30d - tt.want30d; diff > eps || diff < -eps {
				t.Errorf("ConsistencyRate30d: got %f, want %f", result.ConsistencyRate30d, tt.want30d)
			}
		})
	}
}

func TestCalculateStreak_NoSchedule(t *testing.T) {
	// A habit with no schedule should return a zero streak
	habit := newTestHabit("NoSchedule", DirectionBuild, GoalTypeDuration, nil)
	result := CalculateStreak(habit, nil, nil, makeDate(2026, 3, 30))

	if result.Current != 0 {
		t.Errorf("expected current streak 0 for habit with no schedule, got %d", result.Current)
	}
	if result.Longest != 0 {
		t.Errorf("expected longest streak 0 for habit with no schedule, got %d", result.Longest)
	}
}
