package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type RentalController struct {
	DB *RentalsDB
}

func (rc *RentalController) QueryRentals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Query parameters
	priceMin := r.URL.Query().Get("price_min")
	priceMax := r.URL.Query().Get("price_max")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	ids := r.URL.Query().Get("ids")
	near := r.URL.Query().Get("near")
	sort := r.URL.Query().Get("sort")

	field_map := map[string]string{"price": "price_per_day"}
	sort_field, ok := field_map[sort]
	if !ok {
		sort_field = sort
	}

	// Database query

	wheres := []string{}
	if priceMin != "" {
		wheres = append(wheres, "price_per_day >= "+priceMin)
	}
	if priceMax != "" {
		wheres = append(wheres, "price_per_day <= "+priceMax)
	}
	if len(ids) > 0 {
		wheres = append(wheres, "rentals.id = ANY(ARRAY["+ids+"])")
	}
	if near != "" {
		pos := strings.Split(near, ",")
		wheres = append(wheres, "ST_Distance(CONCAT('POINT (', lat, ' ', lng, ')')::geography, 'POINT("+pos[0]+" "+pos[1]+")'::geography) < 100*1609.34")
	}

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
		ON users.id = rentals.user_id`
	if len(wheres) > 0 {
		query += " WHERE " + strings.Join(wheres, " AND ")
	}
	if sort != "" {
		query += " ORDER BY " + sort_field + " DESC"
	}
	if limit != "" {
		query += " LIMIT " + limit
	}
	if offset != "" {
		query += " OFFSET " + offset
	}

	rows, err := rc.DB.Query(query)

	// Error handling
	if err != nil {
		fmt.Fprintf(w, "Error querying database: %s", err)
		return
	}

	// Retrieve rows from DB
	var rentals []Rental
	for rows.Next() {
		var rental Rental
		var s string = ""
		rows.Scan(&s)
		err := rows.Scan(&rental.ID, &rental.Name, &rental.Description,
			&rental.Type, &rental.Make, &rental.Model, &rental.Year, &rental.Length,
			&rental.Sleeps, &rental.PrimaryImageURL, &rental.Price.Day,
			&rental.Location.City, &rental.Location.State, &rental.Location.Zip,
			&rental.Location.Country, &rental.Location.Lat, &rental.Location.Lng,
			&rental.User.ID, &rental.User.FirstName, &rental.User.LastName,
		)
		if err != nil {
			fmt.Fprintf(w, "Error scanning row: %s", err)
			return
		}
		rentals = append(rentals, rental)
	}

	// Return as json
	json.NewEncoder(w).Encode(rentals)
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
