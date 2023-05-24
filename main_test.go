package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	godotenv.Load()
	db, err := InitDB()
	defer db.Close()

	controller := RentalController{db}
	router := CreateRouter(controller)

	// TestGetRentalByID
	req, err := http.NewRequest("GET", "/rentals/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := "{\"id\":1,\"name\":\"'Abaco' VW Bay Window: Westfalia Pop-top\",\"description\":\"ultrices consectetur torquent posuere phasellus urna faucibus convallis fusce sem felis malesuada luctus diam hendrerit fermentum ante nisl potenti nam laoreet netus est erat mi\",\"type\":\"camper-van\",\"make\":\"Volkswagen\",\"model\":\"Bay Window\",\"year\":1978,\"length\":15,\"sleeps\":4,\"primary_image_url\":\"https://res.cloudinary.com/outdoorsy/image/upload/v1528586451/p/rentals/4447/images/yd7txtw4hnkjvklg8edg.jpg\",\"price\":{\"day\":16900},\"location\":{\"city\":\"Costa Mesa\",\"state\":\"CA\",\"zip\":\"92627\",\"country\":\"US\",\"lat\":33.64,\"lng\":-117.93},\"user\":{\"id\":1,\"first_name\":\"John\",\"last_name\":\"Smith\"}}\n"
	assert.Equal(t, expected, rr.Body.String())

	// TestGetAllRentals
	req, err = http.NewRequest("GET", "/rentals", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var rentals []Rental
	err = json.NewDecoder(rr.Body).Decode(&rentals)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, rentals, 30)

	// TestPriceRange
	req, err = http.NewRequest("GET", "/rentals?price_min=9000&price_max=75000", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	err = json.NewDecoder(rr.Body).Decode(&rentals)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, rentals, 24)
	for i := range rentals {
		assert.InDelta(t, rentals[i].Price.Day, 9000, 75000)
	}

	// TestPagination
	req, err = http.NewRequest("GET", "/rentals?limit=3&offset=6", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	err = json.NewDecoder(rr.Body).Decode(&rentals)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, rentals, 3)

	// TestIDs
	req, err = http.NewRequest("GET", "/rentals?ids=3,4,5", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	err = json.NewDecoder(rr.Body).Decode(&rentals)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, rentals, 3)
	for i := range rentals {
		assert.Contains(t, []int{3, 4, 5}, rentals[i].ID)
		assert.NotContains(t, []int{6, 7, 8}, rentals[i].ID)
	}

	// TestNears
	req, err = http.NewRequest("GET", "/rentals?near=33.64,-117.93", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	err = json.NewDecoder(rr.Body).Decode(&rentals)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, rentals, 6)

	// TestSortByPrice
	req, err = http.NewRequest("GET", "/rentals?sort=price", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	err = json.NewDecoder(rr.Body).Decode(&rentals)
	if err != nil {
		t.Fatal(err)
	}
	for i := range rentals {
		if i == 0 {
			continue
		}
		assert.GreaterOrEqual(t, rentals[i-1].Price.Day, rentals[i].Price.Day)
	}

	// TestAll
	req, err = http.NewRequest("GET", "/rentals?near=33.64,-117.93&price_min=9000&price_max=75000&limit=3&offset=6&sort=price", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	err = json.NewDecoder(rr.Body).Decode(&rentals)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, rentals, 0)

}
