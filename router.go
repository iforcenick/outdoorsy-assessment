package main

import "github.com/gorilla/mux"

func CreateRouter(controller RentalController) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/rentals/{id}", controller.GetRental).Methods("GET")
	router.HandleFunc("/rentals", controller.QueryRentals).Methods("GET")
	return router
}
