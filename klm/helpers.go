package klm

import "time"

// Returns a slice of dates for the `dayOfWeek` in the next `numberOfweeks`
func GetDayOfWeek(dayOfWeek time.Weekday, numberOfWeeks int) []time.Time {
	// first we need to find the next `dayOfWeek`
	today := time.Now().Weekday()
	var offset time.Weekday
	if today != dayOfWeek {
		offset = 7 - today + dayOfWeek
	} else {
		offset = 0
	}

	start := time.Now().AddDate(0, 0, int(offset))

	res := make([]time.Time, 0, numberOfWeeks)
	for i := 0; i < numberOfWeeks; i++ {
		res = append(res, start.AddDate(0, 0, 7*i))
	}

	return res
}
