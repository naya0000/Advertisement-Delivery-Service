package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	// "github.com/go-pg/pg"
)

type Advertisement struct {
	Title      string     `json:"title"`
	StartAt    time.Time  `json:"startAt"`
	EndAt      time.Time  `json:"endAt"`
	Conditions Conditions `json:"conditions"`
}

// Conditions represents conditions for displaying an advertisement
type Conditions struct {
	AgeStart *uint8   `json:"ageStart,omitempty"`
	AgeEnd   *uint8   `json:"ageEnd,omitempty"`
	Gender   *string  `json:"gender,omitempty"`
	Country  []string `json:"country,omitempty"`
	Platform []string `json:"platform,omitempty"`
}

// CreateAdvertisement creates a new advertisement in the database
func CreateAdvertisement(db *pg.DB, ad *Advertisement) error {

	sqlStatement := `
		INSERT INTO advertisement (title, start_at, end_at, conditions)
		VALUES (?, ?, ?, ?)
		RETURNING id
	`
	// Convert Conditions to JSON string
	conditionsJSON, err := json.Marshal(ad.Conditions)
	if err != nil {
		return fmt.Errorf("failed to marshal conditions: %v", err)
	}

	// Execute the SQL statement
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// var id int
	_, err = db.ExecContext(ctx, sqlStatement, ad.Title, ad.StartAt, ad.EndAt)
	if err != nil {
		return fmt.Errorf("failed to create advertisement: %v", err)
	}

	fmt.Println("Advertisement created successfully")
	return nil
}

// CreateAdvertisementHandler handles creating advertisements
func CreateAdvertisementHandler(w http.ResponseWriter, r *http.Request) {
	var ad Advertisement
	err := json.NewDecoder(r.Body).Decode(&ad)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//get the db from context
	db, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Call the CreateAdvertisement function to insert the advertisement into the database
	err = CreateAdvertisement(db, &ad)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write a success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Advertisement created successfully")
}

// start api with the pgdb and return a chi router
func StartAPI(db *pg.DB) *chi.Mux {

	//get the router
	r := chi.NewRouter()

	// add middleware
	// store our DB to use it later
	r.Use(middleware.Logger, middleware.WithValue("DB", db))

	r.Route("/api/v1/ad", func(r chi.Router) {
		r.Post("/", CreateAdvertisementHandler)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up and running"))
	})

	return r
}
