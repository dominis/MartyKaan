package klm

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
	Price       float32
}

type OfferDetails struct {
	UpsellProducts []struct {
		Price struct {
			DisplayPrice           float32 `json:"displayPrice"`
			TotalPrice             float32 `json:"totalPrice"`
			Accuracy               float32 `json:"accuracy"`
			PricePerPassengerTypes []struct {
				PassengerType string  `json:"passengerType"`
				Fare          float32 `json:"fare"`
				Taxes         float32 `json:"taxes"`
			} `json:"pricePerPassengerTypes"`
			FlexibilityWaiver bool   `json:"flexibilityWaiver"`
			Currency          string `json:"currency"`
			DisplayType       string `json:"displayType"`
		} `json:"price"`
	}
}

func GetOffers(origin, destination FlightDetails) FlightOffering {
	d := FlightOffering{
		Origin:      origin,
		Destination: destination,
	}

	reqJSON := prepareRequestJson(d)
	offers := sendAPIRequest(reqJSON)
	best := FindCheapest(offers)

	fmt.Printf("origin: %s-%s | destination %s-%s | lowestfare: %.2f\n",
		origin.AirportCode,
		origin.DepartureDate(),
		destination.AirportCode,
		destination.DepartureDate(),
		best)

	d.Price = best

	return d
}

func prepareRequestJson(data FlightOffering) bytes.Buffer {
	t, err := template.ParseFiles("upsell.json")
	if err != nil {
		fmt.Println(err)
	}
	buf := new(bytes.Buffer)
	t.Execute(buf, &data)
	return *buf
}

func sendAPIRequest(jsonStr bytes.Buffer) OfferDetails {
	url := "https://www.klm.com/ams/search-web/api/upsell-products?country=CH&language=en&localeStringDateTime=en&localeStringNumber=de-CH"

	req, err := http.NewRequest("POST", url, &jsonStr)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var offerDetails OfferDetails
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&offerDetails)
	if err != nil {
		panic(err)
	}

	return offerDetails
}

func FindCheapest(offers OfferDetails) float32 {
	lowest := float32(0)
	for _, p := range offers.UpsellProducts {
		if lowest == 0 {
			lowest = p.Price.DisplayPrice
		}

		if lowest > p.Price.DisplayPrice {
			lowest = p.Price.DisplayPrice
		}
	}

	return lowest
}
