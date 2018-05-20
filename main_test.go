package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetDayOfWeekLen(t *testing.T) {
	dates := getDayOfWeek(time.Monday, 10)
	if len(dates) != 10 {
		t.Errorf("Required 10 dates")
	}
}

func TestGetDayOfWeekNextDatte(t *testing.T) {
	today := time.Now().Weekday()
	next := getDayOfWeek(today, 1)
	fmt.Printf("%v", next.Weekday())
}
