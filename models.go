package main

type Rental struct {
	ID              int           `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Type            string        `json:"type"`
	Make            string        `json:"make"`
	Model           string        `json:"model"`
	Year            int           `json:"year"`
	Length          float32       `json:"length"`
	Sleeps          int           `json:"sleeps"`
	PrimaryImageURL string        `json:"primary_image_url"`
	Price           rentalPrice   `json:"price"`
	Location        rentalAddress `json:"location"`
	User            rentalUser    `json:"user"`
}

type rentalPrice struct {
	Day int `json:"day"`
}

type rentalAddress struct {
	City    string  `json:"city"`
	State   string  `json:"state"`
	Zip     string  `json:"zip"`
	Country string  `json:"country"`
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
}

type rentalUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
