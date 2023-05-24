package main

import "net/http"

type RentalController struct {
	DB *RentalsDB
}

func (rc *RentalController) QueryRentals(w http.ResponseWriter, r *http.Request) {
}
func (rc *RentalController) GetRental(w http.ResponseWriter, r *http.Request) {
}
