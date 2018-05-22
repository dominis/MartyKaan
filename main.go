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
	travelDays := klm.GetDayOfWeek(time.Thursday, 23)

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

			origin_offers = append(origin_offers, klm.GetOffers(origin_bud, destination_ams))
			destination_offers = append(destination_offers, klm.GetOffers(origin_ams, destination_bud))
			time.Sleep(1000 * time.Millisecond)
		}

		origin_best := findBest(origin_offers)
		destination_best := findBest(destination_offers)

		fmt.Println("Best:")
		fmt.Printf("origin: %s-%s | destination %s-%s | lowestfare: %.2fEUR\n",
			origin_best.Origin.AirportCode,
			origin_best.Origin.DepartureDate(),
			origin_best.Destination.AirportCode,
			origin_best.Destination.DepartureDate(),
			origin_best.Price/EurHuf)
		fmt.Printf("origin: %s-%s | destination %s-%s | lowestfare: %.2fEUR\n",
			destination_best.Origin.AirportCode,
			destination_best.Origin.DepartureDate(),
			destination_best.Destination.AirportCode,
			destination_best.Destination.DepartureDate(),
			destination_best.Price)
		fmt.Printf("SUM: %.2fEUR\n", destination_best.Price+(origin_best.Price/EurHuf))

		fmt.Println("###############")
	}
}

func findBest(offers []klm.FlightOffering) klm.FlightOffering {
	best := klm.FlightOffering{Price: 100000000000}
	for _, o := range offers {
		if o.Price < best.Price &&
			o.Price > 0.0 {
			best = o
		}
	}

	return best
}
