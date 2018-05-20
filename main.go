package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const (
	departureDateFormat = "2006-01-02"
	BUDFlightNumber     = 1972
	AMSFlightNumber     = 1981
)

type FlightDetails struct {
	AirportCode  string
	FlightNumber int32
	Date         time.Time
}

func (f *FlightDetails) DepartureDateTime() string {
	var time string
	if f.AirportCode == "BUD" {
		time = "T06:30:00"
	} else if f.AirportCode == "AMS" {
		time = "T20:55:00"
	}
	return string(f.Date.Format(departureDateFormat) + time)
}

func (f *FlightDetails) DepartureDate() string {
	return f.Date.Format(departureDateFormat)
}

type FlightOffering struct {
	Origin      FlightDetails
	Destination FlightDetails
}

type OfferDetails struct {
	UpsellProducts []struct {
		Price struct {
			DisplayPrice           int `json:"displayPrice"`
			TotalPrice             int `json:"totalPrice"`
			Accuracy               int `json:"accuracy"`
			PricePerPassengerTypes []struct {
				PassengerType string `json:"passengerType"`
				Fare          int    `json:"fare"`
				Taxes         int    `json:"taxes"`
			} `json:"pricePerPassengerTypes"`
			FlexibilityWaiver bool   `json:"flexibilityWaiver"`
			Currency          string `json:"currency"`
			DisplayType       string `json:"displayType"`
		} `json:"price"`
	}
}

func main() {

	travelDays := getDayOfWeek(time.Thursday, 23)

	for _, t := range travelDays {
		origin := FlightDetails{
			AirportCode:  "BUD",
			FlightNumber: BUDFlightNumber,
			Date:         t,
		}

		destination := FlightDetails{
			AirportCode:  "AMS",
			FlightNumber: AMSFlightNumber,
			Date:         t,
		}

		getOffers(origin, destination)
		return
	}
}

// Returns a slice of dates for the `dayOfWeek` in the next `numberOfweeks`
func getDayOfWeek(dayOfWeek time.Weekday, numberOfWeeks int) []time.Time {
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

func getOffers(origin, destination FlightDetails) {
	d := FlightOffering{
		Origin:      origin,
		Destination: destination,
	}
	t, err := template.ParseFiles("upsell.json")
	if err != nil {
		fmt.Println(err)
	}
	buf := new(bytes.Buffer)
	t.Execute(buf, &d)
	sendAPIRequest(*buf)
}

func sendAPIRequest(jsonStr bytes.Buffer) {
	url := "https://www.klm.com/ams/search-web/api/upsell-products?country=CH&language=en&localeStringDateTime=en&localeStringNumber=de-CH"
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("POST", url, &jsonStr)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)

	var offerDetails OfferDetails
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&offerDetails)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", offerDetails)

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// fmt.Println("response Body:", string(body))
}
