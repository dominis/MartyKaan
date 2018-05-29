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

	travelDays := klm.GetNextDays(200)
	//startDate := make([]time.Time, 0)
	//startDate = append(startDate, time.Date(2018, 06, 14, 1, 0, 0, 0, time.UTC))
	startDate := klm.GetMultipleDaysOfWeek([]time.Weekday{time.Thursday}, 25)

	avgSlice := make([]float32, 0)

	for _, t := range startDate {
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

			//PrintResult(o)
			//PrintResult(d)
		}

		o_best := findBest(origin_offers)
		d_best := findBest(destination_offers)
		fmt.Println("#### Bestest ####")
		PrintResult(o_best)
		PrintResult(d_best)
		bestest := []klm.FlightOffering{o_best, d_best}
		fmt.Printf("*** AVG per flight: %.2f EUR ***\n", avgOffers(bestest))
		avgSlice = append(avgSlice, avgOffers(bestest))
	}
	fmt.Println("####")
	fmt.Printf("Summ avg: %.2f EUR\n", calcAvg(avgSlice))

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

func findBest(offers []klm.FlightOffering) klm.FlightOffering {
	best := klm.FlightOffering{}
	best.Price = 1000000
	for _, o := range offers {
		if o.Price > 0 && o.Price < best.Price {
			best = o
		}
	}
	return best
}

func avgOffers(o []klm.FlightOffering) float32 {
	len := 0
	sum := float32(0)

	for _, o := range o {
		if o.Currency != "EUR" {
			o.Price = o.Price / EurHuf
		}
		len++
		sum += o.Price
	}

	return sum / float32(len)
}

func calcAvg(nums []float32) float32 {
	sum := float32(0)
	for _, i := range nums {
		sum += i
	}

	return sum / float32(len(nums))
}
