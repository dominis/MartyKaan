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
	Currency    string
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

	if destination.Date.Sub(origin.Date) < 0 {
		return d
	}

	var has_results bool
	d.Price, has_results = DbGetPrice(d)

	if has_results == false { // no price in the database
		reqJSON := prepareRequestJson(d)
		offers, _ := sendAPIRequest(reqJSON)
		d.Price, d.Currency = FindCheapest(offers)
		// cache result in db
		DbInsertOffer(d)
	}

	if d.Currency == "" { // set currency here
		// ASSUMPTION IS THE MOTHER OF ALL FUCKUPS
		if d.Price < 10000 {
			d.Currency = "EUR"
		} else {
			d.Currency = "HUF"
		}
	}

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

func sendAPIRequest(jsonStr bytes.Buffer) (OfferDetails, error) {
	url := "https://www.klm.com/ams/search-web/api/upsell-products?country=CH&language=en&localeStringDateTime=en&localeStringNumber=de-CH"

	req, err := http.NewRequest("POST", url, &jsonStr)
	req.Header.Set("Content-Type", "application/json")

	var offerDetails OfferDetails
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("%s", err)
		return offerDetails, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&offerDetails)
	if err != nil {
		fmt.Errorf("%s", err)
		return offerDetails, err
	}

	return offerDetails, nil
}

func FindCheapest(offers OfferDetails) (float32, string) {
	lowest := float32(0)
	currency := ""
	for _, p := range offers.UpsellProducts {
		if lowest == 0 {
			lowest = p.Price.DisplayPrice
		}

		if lowest > p.Price.DisplayPrice {
			lowest = p.Price.DisplayPrice
		}

		currency = p.Price.Currency
	}

	return lowest, currency
}
