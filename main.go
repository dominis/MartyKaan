package main

import (
	"fmt"
	"time"

	"github.com/dominis/martykaan/klm"
)

const (
	BUDFlightNumber = 1972
	AMSFlightNumber = 1981
)

func main() {
	travelDays := klm.GetDayOfWeek(time.Thursday, 23)

	for _, t := range travelDays {
		origin := klm.FlightDetails{
			AirportCode:  "BUD",
			FlightNumber: BUDFlightNumber,
			Date:         t,
		}

		for _, d := range travelDays {
			destination := klm.FlightDetails{
				AirportCode:  "AMS",
				FlightNumber: AMSFlightNumber,
				Date:         d,
			}

			klm.GetOffers(origin, destination)
			time.Sleep(1000 * time.Millisecond)
		}

		fmt.Println("##### Next Date #####")
	}
}

func debug() {
	origin := klm.FlightDetails{
		AirportCode:  "BUD",
		FlightNumber: BUDFlightNumber,
		Date:         time.Date(2018, time.June, 7, 1, 0, 0, 0, time.UTC),
	}
	destination := klm.FlightDetails{
		AirportCode:  "AMS",
		FlightNumber: AMSFlightNumber,
		Date:         time.Date(2018, time.June, 21, 1, 0, 0, 0, time.UTC),
	}
	klm.GetOffers(origin, destination)

}
