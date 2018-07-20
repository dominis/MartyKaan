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
	dateFormat      = "2006-01-02"
)

var coveredDestDays []string

var coveredBUDDates = []time.Time{
	time.Date(2018, 6, 14, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 6, 21, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 6, 28, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 7, 5, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 7, 12, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 7, 19, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 7, 26, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 8, 30, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 9, 6, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 9, 27, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 11, 8, 1, 3, 4, 4, time.UTC),
}

var coveredAMSdates = []time.Time{
	time.Date(2018, 6, 14, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 6, 21, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 6, 28, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 7, 12, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 7, 19, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 8, 16, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 8, 23, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 8, 30, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 9, 6, 1, 3, 4, 4, time.UTC),
	time.Date(2018, 10, 18, 1, 3, 4, 4, time.UTC),

	time.Date(2018, 11, 22, 1, 3, 4, 4, time.UTC),
}

func main() {
	klm.Init_db()
	thursdays := klm.GetMultipleDaysOfWeek([]time.Weekday{time.Thursday}, 25)
	alldays := klm.GetNextDays(365)
	BUDDates := removeDupDates(thursdays, coveredBUDDates)
	AMSDates := removeDupDates(thursdays, coveredAMSdates)
	EmptyLegsBUD := removeDupDates(alldays, coveredBUDDates)
	EmptyLegsAMS := removeDupDates(alldays, coveredAMSdates)

	avgSlice := make([]float32, 0)

	origin_offers := make([]klm.FlightOffering, 0)
	destination_offers := make([]klm.FlightOffering, 0)

	for _, t := range BUDDates {
		origin_offers = []klm.FlightOffering{}
		destination_offers = []klm.FlightOffering{}

		origin_bud := klm.FlightDetails{
			AirportCode:  "BUD",
			FlightNumber: BUDFlightNumber,
			Date:         t,
		}

		destination_bud := klm.FlightDetails{
			AirportCode:  "BUD",
			FlightNumber: BUDFlightNumber,
			Date:         t,
		}

		for _, d := range AMSDates {
			origin_ams := klm.FlightDetails{
				AirportCode:  "AMS",
				FlightNumber: AMSFlightNumber,
				Date:         d,
			}

			destination_ams := klm.FlightDetails{
				AirportCode:  "AMS",
				FlightNumber: AMSFlightNumber,
				Date:         d,
			}

			orig := klm.GetOffers(origin_bud, destination_ams)
			dest := klm.GetOffers(origin_ams, destination_bud)

			origin_offers = append(origin_offers, orig)
			destination_offers = append(destination_offers, dest)
		}

		o_best := findBest(origin_offers)
		d_best := findBest(destination_offers)
		PrintResult(o_best)
		PrintResult(d_best)
		bestest := []klm.FlightOffering{o_best, d_best}
		fmt.Printf("*** AVG per flight: %.2f EUR ***\n", avgOffers(bestest))
		avgSlice = append(avgSlice, avgOffers(bestest))
	}

	fmt.Println("Empty legs: ")
	for _, t := range BUDDates {
		origin_offers = []klm.FlightOffering{}
		origin_bud := klm.FlightDetails{
			AirportCode:  "BUD",
			FlightNumber: BUDFlightNumber,
			Date:         t,
		}
		for _, d := range EmptyLegsAMS {
			destination_ams := klm.FlightDetails{
				AirportCode:  "AMS",
				FlightNumber: AMSFlightNumber,
				Date:         d,
			}
			orig := klm.GetOffers(origin_bud, destination_ams)
			origin_offers = append(origin_offers, orig)
		}

		o_best := findBest(origin_offers)
		PrintResult(o_best)
	}

	for _, t := range AMSDates {
		origin_offers = []klm.FlightOffering{}
		origin_bud := klm.FlightDetails{
			AirportCode:  "AMS",
			FlightNumber: AMSFlightNumber,
			Date:         t,
		}
		for _, d := range EmptyLegsBUD {
			destination_ams := klm.FlightDetails{
				AirportCode:  "BUD",
				FlightNumber: BUDFlightNumber,
				Date:         d,
			}
			orig := klm.GetOffers(origin_bud, destination_ams)
			origin_offers = append(origin_offers, orig)
		}

		o_best := findBest(origin_offers)
		PrintResult(o_best)
	}

}

func PrintResult(offer klm.FlightOffering) {
	if offer.Price == 0 {
		return
	}
	if offer.Currency == "HUF" {
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
	best.Price = 0

	for _, o := range offers {
		if best.Price == 0 {
			best = o
		}
		if o.Price > 0 &&
			o.Price < best.Price {
			best = o
		}
	}
	return best
}

func avgOffers(o []klm.FlightOffering) float32 {
	len := 0
	sum := float32(0)

	for _, o := range o {
		if o.Currency == "HUF" {
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

func inSlice(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false

}

func removeDupDates(list []time.Time, remove []time.Time) []time.Time {
	ret := make([]time.Time, 0)
	dup := false
	for _, a := range list {
		dup = false
		for _, b := range remove {
			if a.Format(dateFormat) == b.Format(dateFormat) {
				dup = true
			}
		}
		if dup != true {
			ret = append(ret, a)
		}
	}

	return ret
}
