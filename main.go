package main

import (
	"fmt"
	"time"
)

const upsellJsonTpl = `{"country":"CH","language":"en","cabinClass":"ECONOMY","passengerCount":{"ADULT":1,"CHILD":0,"INFANT":0},"requestedConnections":[{"origin":{"airport":{"code":"{{Origin.AirportCode}}"}},"destination":{"airport":{"code":"{{Destination.AirportCode}}"}},"departureDate":"{{Origin.DepartureDate}}","fareBasis":"ASRHU","segments":[{"marketingCarrier":"KL","marketingFlightNumber":"{{Origin.FlightNumber}}","origin":{"code":"{{Origin.AirportCode}}"},"destination":{"code":"{{Destination.AirportCode}}"},"departureDateTime":"{{Origin.DepartureTime}}","sellingClass":"A"}]},{"origin":{"airport":{"code":"{{Destination.AirportCode}}"}},"destination":{"airport":{"code":"{{Origin.AirportCode}}"}},"departureDate":"{{Destination.DepartureTime}}","fareBasis":"FWKHU","segments":[{"marketingCarrier":"KL","marketingFlightNumber":"{{Destination.FlightNumber}}","origin":{"code":"{{Destination.AirportCode}}"},"destination":{"code":"{{Origin.AirportCode}}"},"departureDateTime":"{{Destination.DepartureTime}}","sellingClass":"F"}]}],"localeStringDateTime":"en","localeStringNumber":"de-CH"}`

type FlightDetails struct {
	AirportCode   string    //AMS
	FlightNumber  int32     //1981
	DepartureTime time.Time //2018-06-21T20:55:00
}

type FlightOffering struct {
	Origin      FlightDetails
	Destination FlightDetails
}

func main() {
	travelDays := getDayOfWeek(time.Thursday, 23)
	for _, d := range travelDays {
		fmt.Println(d)
	}
}

// Returns a slice of dates for the `dayOfWeek` in the next `numberOfweeks`
func getDayOfWeek(dayOfWeek time.Weekday, numberOfWeeks int) []time.Time {
	// first we need to find the next `dayOfWeek`
	today := time.Now().Weekday()
	offset := 7 + today - dayOfWeek
	start := time.Now().AddDate(0, 0, int(offset))

	res := make([]time.Time, 0, numberOfWeeks)
	for i := 0; i < numberOfWeeks; i++ {
		res = append(res, start.AddDate(0, 0, 7*i))
	}

	return res
}

func buildRequestJson() {
}
