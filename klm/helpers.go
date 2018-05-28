package klm

import "time"

// Returns a slice of dates for the `dayOfWeek` in the next `numberOfweeks`
func GetDayOfWeek(dayOfWeek time.Weekday, numberOfWeeks int) []time.Time {
	// first we need to find the next `dayOfWeek`
	today := time.Now().Weekday()
	var offset time.Weekday
	offset = 7 - today + dayOfWeek

	start := time.Now().AddDate(0, 0, int(offset))

	res := make([]time.Time, 0, numberOfWeeks)
	for i := 0; i < numberOfWeeks; i++ {
		res = append(res, start.AddDate(0, 0, 7*i))
	}

	return res
}

func GetMultipleDaysOfWeek(daysOfWeek []time.Weekday, numberOfWeeks int) []time.Time {
	result := make([]time.Time, 0)
	for _, d := range daysOfWeek {
		for _, r := range GetDayOfWeek(d, numberOfWeeks) {
			result = append(result, r)
		}
	}

	return result
}

func RangeSlice(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
