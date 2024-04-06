package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8" // Import Redis package
	"github.com/lib/pq"
	// "github.com/go-pg/pg"
)

// Define a Redis client
var redisClient *redis.Client

// InitializeRedis initializes the Redis client
func InitializeRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Change this to your Redis server address
		Password: "",           // No password
		DB:       0,            // Use default DB
	})
	// Ping the Redis server to check connection
	ctx := redisClient.Context()
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)
}

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

// CreateAd creates a new advertisement in the database
func CreateAd(db *sql.DB, ad *Advertisement) error {
	// starts a transaction to insert ad and its conditions
	tx, err := db.Begin()
	checkErr(err)

	insertAdStmt, err := tx.Prepare(`
		INSERT INTO advertisement (title, start_at, end_at) 
		VALUES($1, $2, $3) 
		RETURNING id
	`)
	checkErr(err)
	defer insertAdStmt.Close()

	var adInsertId int
	err = insertAdStmt.QueryRow(ad.Title, ad.StartAt, ad.EndAt).Scan(&adInsertId)
	checkErr(err)
	fmt.Println("Last inserted Ad ID =", adInsertId)

	insertConStmt, err := tx.Prepare(`
		INSERT INTO condition (ad_id, age_start, age_end, gender, country, platform) 
		VALUES($1, COALESCE($2, 1), COALESCE($3, 100), $4, $5, $6) 
		RETURNING ad_id
	`)
	checkErr(err)
	defer insertConStmt.Close()

	var conInsertId int
	err = insertConStmt.QueryRow(adInsertId, ad.Conditions.AgeStart, ad.Conditions.AgeEnd, ad.Conditions.Gender, pq.Array(ad.Conditions.Country), pq.Array(ad.Conditions.Platform)).Scan(&conInsertId)
	checkErr(err)
	fmt.Println("Last inserted Condition ID =", conInsertId)

	// commits the transaction
	err = tx.Commit()
	checkErr(err)

	return err
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// CreateAdHandler handles creating advertisements
func CreateAdHandler(w http.ResponseWriter, r *http.Request) {
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
	err = CreateAd(db, &ad)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write a success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Advertisement created successfully")
}
func GetAdHandler(w http.ResponseWriter, r *http.Request) {

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

	// Define a key for Redis cache
	cacheKey := fmt.Sprintf("ad:%s:%s:%s:%s:%s:%s", age, gender, country, platform, offset, limit)

	// Check if data is available in Redis cache
	cachedData, err := redisClient.Get(r.Context(), cacheKey).Result()
	if err == nil {
		// Data found in cache, return cached data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cachedData))
		return
	} else if err != redis.Nil {
		// Error occurred while retrieving data from Redis
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Query database for advertisements based on conditions
	db, ok := r.Context().Value("DB").(*sql.DB)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	query := `
		SELECT a.title, a.end_at
		FROM Advertisement a
		INNER JOIN Condition c ON a.id = c.ad_id
		WHERE a.start_at < NOW() AND a.end_at > NOW()
	`
	// Initialize an empty slice to hold query parameters
	var queryParams []interface{}

	// Add conditions to the query, handling NULL values appropriately
	if age != "" {
		query += " AND (c.age_start IS NULL OR c.age_start <= $1) AND (c.age_end IS NULL OR c.age_end >= $2)"
		ageStart, _ := strconv.Atoi(age)
		ageEnd, _ := strconv.Atoi(age)
		queryParams = append(queryParams, ageStart, ageEnd)
	}
	if gender != "" {
		query += " AND (c.gender IS NULL OR c.gender = $3)"
		queryParams = append(queryParams, gender)
	}
	if country != "" {
		query += " AND (c.country IS NULL OR $4 = ANY (c.country))"
		queryParams = append(queryParams, country)
	}
	if platform != "" {
		query += " AND (c.platform IS NULL OR $5 = ANY (c.platform))"
		queryParams = append(queryParams, platform)
	}

	// Add pagination to the query
	query += fmt.Sprintf(" ORDER BY a.end_at ASC OFFSET $%d LIMIT $%d", len(queryParams)+1, len(queryParams)+2)
	queryParams = append(queryParams, offset, limit)

	// Execute the query with prepared statement and query parameters
	rows, err := db.Query(query, queryParams...)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create struct to hold advertisement data
	type AdResponse struct {
		Title string    `json:"title"`
		EndAt time.Time `json:"endAt"`
	}

	// Create slice to hold advertisement data
	var advertisements []AdResponse

	// Iterate over query results and append to advertisements slice
	for rows.Next() {
		var ad AdResponse
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
		Items []AdResponse `json:"items"`
	}{
		Items: advertisements,
	}

	// Marshal response struct to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cache the response in Redis for future requests
	err = redisClient.Set(r.Context(), cacheKey, jsonResponse, 10*time.Minute).Err() // Cache for 10 minutes
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

	InitializeRedis()

	//get the router
	r := chi.NewRouter()

	// add middleware to store our DB to use it later
	r.Use(middleware.Logger, middleware.WithValue("DB", db))

	// Create a channel to communicate with goroutines
	ch := make(chan http.HandlerFunc)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	// Define the number of concurrent requests
	numRequests := 16

	// Launch goroutines to handle concurrent requests
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- GetAdHandler
		}()
	}

	go func() {
		wg.Wait() // Wait for all goroutines to finish
		close(ch) // Close the channel to signal that all responses have been received
	}()

	// Collect responses from the channel
	for handler := range ch {
		r.Get("/api/v1/ad", handler)
	}

	r.Post("/api/v1/ad", CreateAdHandler)

	return r
}
