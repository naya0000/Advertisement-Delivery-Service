package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lib/pq"
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
	AdID     uint8    `json:"ad_id"`
	AgeStart *uint8   `json:"ageStart,omitempty"`
	AgeEnd   *uint8   `json:"ageEnd,omitempty"`
	Gender   *string  `json:"gender,omitempty"`
	Country  []string `json:"country,omitempty"`
	Platform []string `json:"platform,omitempty"`
}

// CreateAdvertisement creates a new advertisement in the database
func CreateAdvertisement(db *sql.DB, ad *Advertisement) error {
	tx, err := db.Begin()
	checkErr(err)

	insertAdStmt, err := tx.Prepare("INSERT INTO advertisement (title, start_at, end_at) VALUES($1, $2, $3) RETURNING id")
	checkErr(err)
	defer insertAdStmt.Close()

	var adInsertId int
	err = insertAdStmt.QueryRow(ad.Title, ad.StartAt, ad.EndAt).Scan(&adInsertId)
	checkErr(err)
	fmt.Println("Last inserted Ad ID =", adInsertId)

	insertConStmt, err := tx.Prepare("INSERT INTO condition (ad_id, age_start, age_end, gender, country, platform) VALUES($1, $2, $3, $4, $5, $6) RETURNING id")
	checkErr(err)
	defer insertConStmt.Close()

	var conInsertId int
	err = insertConStmt.QueryRow(adInsertId, ad.Conditions.AgeStart, ad.Conditions.AgeEnd, ad.Conditions.Gender, pq.Array(ad.Conditions.Country), pq.Array(ad.Conditions.Platform)).Scan(&conInsertId)
	checkErr(err)
	fmt.Println("Last inserted Condition ID =", conInsertId)

	err = tx.Commit()
	checkErr(err)

	return err
	// Convert Conditions to JSON string
	// conditionsJSON, err := json.Marshal(ad.Conditions)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal conditions: %v", err)
	// }

	// Execute the SQL statement
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// var id int
	// var result
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
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
	db, ok := r.Context().Value("DB").(*sql.DB)
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
func GetAdvertisementHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	params := r.URL.Query()
	offsetStr := params.Get("offset")
	limitStr := params.Get("limit")

	// Default values for offset and limit
	offset := 0
	limit := 5

	// Convert offset and limit to integers
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	// Retrieve conditions from query parameters
	age := params.Get("age")
	gender := params.Get("gender")
	country := params.Get("country")
	platform := params.Get("platform")

	// Query database for advertisements based on conditions
	db, ok := r.Context().Value("DB").(*sql.DB)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	query := `
		SELECT title, end_at
		FROM advertisement
		WHERE start_at < NOW() AND end_at > NOW()
	`

	// Add conditions to the query
	if age != "" {
		query += fmt.Sprintf(" AND (%s BETWEEN condition.age_start AND condition.age_end)", age)
	}
	if gender != "" {
		query += fmt.Sprintf(" AND condition.gender = '%s'", gender)
	}
	if country != "" {
		query += fmt.Sprintf(" AND '%s' = ANY (condition.country)", country)
	}
	if platform != "" {
		query += fmt.Sprintf(" AND '%s' = ANY (condition.platform)", platform)
	}

	// Add pagination to the query
	query += fmt.Sprintf(" ORDER BY end_at ASC OFFSET %d LIMIT %d", offset, limit)

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create struct to hold advertisement data
	type AdvertisementResponse struct {
		Title string    `json:"title"`
		EndAt time.Time `json:"endAt"`
	}

	// Create slice to hold advertisement data
	var advertisements []AdvertisementResponse

	// Iterate over query results and append to advertisements slice
	for rows.Next() {
		var ad AdvertisementResponse
		err := rows.Scan(&ad.Title, &ad.EndAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		advertisements = append(advertisements, ad)
	}

	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response struct
	response := struct {
		Items []AdvertisementResponse `json:"items"`
	}{
		Items: advertisements,
	}

	// Marshal response struct to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response content type and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// start api with the pgdb and return a chi router
func StartAPI(db *sql.DB) *chi.Mux {

	//get the router
	r := chi.NewRouter()

	// add middleware to store our DB to use it later
	r.Use(middleware.Logger, middleware.WithValue("DB", db))

	r.Route("/api/v1/ad", func(r chi.Router) {
		r.Post("/", CreateAdvertisementHandler)
		r.Get("/", GetAdvertisementHandler)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up and running"))
	})

	return r
}
