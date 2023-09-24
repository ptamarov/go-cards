package user

import "testing"

func TestDailyGoal(t *testing.T) {
	reached := User{DailyGoal: 50, CorrectToday: 50}
	notReached := User{DailyGoal: 50, CorrectToday: 20}

	if !reached.IsDailyGoalReached() {
		t.Error("expected daily goal reached but got false instead")
	}
	if notReached.IsDailyGoalReached() {
		t.Error("expected daily goal not reached but got true instead")
	}

}
