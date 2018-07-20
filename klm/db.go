package klm

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var database, _ = sql.Open("sqlite3", "./klm.db")

func Init_db() {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS offers (id INTEGER PRIMARY KEY, origin TEXT, destination TEXT, created TEXT, price REAL)")
	statement.Exec()
}

func DbInsertOffer(offer FlightOffering) {
	origin := fmt.Sprintf("%s-%s", offer.Origin.AirportCode, offer.Origin.DepartureDate())
	destination := fmt.Sprintf("%s-%s", offer.Destination.AirportCode, offer.Destination.DepartureDate())
	created := time.Now().Format("2006-01-02 15:04:05")

	statement, _ := database.Prepare("INSERT INTO offers (origin, destination, price, created) VALUES (?, ?, ?, ?)")
	statement.Exec(origin, destination, offer.Price, created)
}

func DbGetPrice(offer FlightOffering) (float32, bool) {
	origin := fmt.Sprintf("%s-%s", offer.Origin.AirportCode, offer.Origin.DepartureDate())
	destination := fmt.Sprintf("%s-%s", offer.Destination.AirportCode, offer.Destination.DepartureDate())

	rows, err := database.Query("SELECT price FROM offers WHERE origin=? AND destination=? LIMIT 1", origin, destination)

	if err != nil {
		log.Fatal(err)
	}

	var price float32
	has_results := false
	for rows.Next() {
		has_results = true
		rows.Scan(&price)
	}

	return price, has_results
}
