package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type RentalController struct {
	DB *RentalsDB
}

func (rc *RentalController) QueryRentals(w http.ResponseWriter, r *http.Request) {
}

func (rc *RentalController) GetRental(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get parameter from url
	params := mux.Vars(r)
	id := params["id"]

	// Database query
	query := `
		SELECT
			rentals.id, rentals.name, rentals.description,
			rentals.type, rentals.vehicle_make, rentals.vehicle_model, rentals.vehicle_year, rentals.vehicle_length,
			rentals.sleeps, rentals.primary_image_url, rentals.price_per_day,
			rentals.home_city, rentals.home_state, rentals.home_zip,
			rentals.home_country, rentals.lat, rentals.lng,
			rentals.user_id, users.first_name, users.last_name
		FROM rentals
		LEFT JOIN users
		ON users.id = rentals.user_id
		WHERE rentals.id = $1::int
		`
	var rental Rental
	err := rc.DB.QueryRow(query, id).Scan(&rental.ID, &rental.Name, &rental.Description,
		&rental.Type, &rental.Make, &rental.Model, &rental.Year, &rental.Length,
		&rental.Sleeps, &rental.PrimaryImageURL, &rental.Price.Day,
		&rental.Location.City, &rental.Location.State, &rental.Location.Zip,
		&rental.Location.Country, &rental.Location.Lat, &rental.Location.Lng,
		&rental.User.ID, &rental.User.FirstName, &rental.User.LastName)

	// Error handling
	if err != nil {
		fmt.Fprintf(w, "Error querying database: %s", err)
		return
	}

	// Return as json
	json.NewEncoder(w).Encode(rental)
}
