package main

import (
	"fmt"
	"time"

	"github.com/dominis/martykaan/klm"
)

const (
	BUDFlightNumber = 1972
	AMSFlightNumber = 1981
	EurHuf          = 312
)

func main() {
	klm.Init_db()

	travelDays := klm.GetMultipleDaysOfWeek([]time.Weekday{time.Tuesday, time.Wednesday, time.Thursday}, 23)

	for _, t := range travelDays {
		origin_offers := make([]klm.FlightOffering, 0)
		destination_offers := make([]klm.FlightOffering, 0)

		origin_bud := klm.FlightDetails{
			AirportCode:  "BUD",
			FlightNumber: BUDFlightNumber,
			Date:         t,
		}
		origin_ams := klm.FlightDetails{
			AirportCode:  "AMS",
			FlightNumber: AMSFlightNumber,
			Date:         t,
		}

		for _, d := range travelDays {
			destination_ams := klm.FlightDetails{
				AirportCode:  "AMS",
				FlightNumber: AMSFlightNumber,
				Date:         d,
			}
			destination_bud := klm.FlightDetails{
				AirportCode:  "BUD",
				FlightNumber: BUDFlightNumber,
				Date:         d,
			}

			o := klm.GetOffers(origin_bud, destination_ams)
			d := klm.GetOffers(origin_ams, destination_bud)

			origin_offers = append(origin_offers, o)
			destination_offers = append(destination_offers, d)

			PrintResult(o)
			PrintResult(d)
		}
	}
}

func PrintResult(offer klm.FlightOffering) {
	if offer.Price == 0 {
		return
	}
	if offer.Currency != "EUR" {
		offer.Price = offer.Price / EurHuf
	}
	fmt.Printf("origin: %s-%s | destination: %s-%s | lowestfare: %.2fEUR\n",
		offer.Origin.AirportCode,
		offer.Origin.DepartureDate(),
		offer.Destination.AirportCode,
		offer.Destination.DepartureDate(),
		offer.Price)
}
