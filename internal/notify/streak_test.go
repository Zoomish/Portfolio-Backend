package notify

import "testing"

func TestStreakAlertsEveryLimit(t *testing.T) {
	streak := NewStreak(3)
	if streak.Fail() || streak.Fail() {
		t.Fatal("alerted too early")
	}
	if !streak.Fail() {
		t.Fatal("expected alert on third failure")
	}
	streak.OK()
	if streak.Fail() {
		t.Fatal("counter was not reset")
	}
}
