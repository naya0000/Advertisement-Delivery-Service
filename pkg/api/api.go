package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

// Define a Redis client
var redisClient *redis.Client

func dbConnect() *sql.DB {
	connStr := "postgres://postgres:admin@db:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

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

// CreateAd creates a new advertisement in the database
func CreateAd(db *sql.DB, ad *Advertisement) error {
	// starts a transaction to insert ad and its conditions
	tx, err := db.Begin()
	checkErr(err)
	defer tx.Rollback()

	insertAdStmt, err := tx.Prepare(`
		INSERT INTO advertisement (title, start_at, end_at, conditions) 
		VALUES($1, $2, $3, $4) 
		RETURNING id
	`)
	checkErr(err)
	defer insertAdStmt.Close()

	// set default value to 1
	if ad.Conditions.AgeStart == nil {
		defaultAgeStart := uint8(1)
		ad.Conditions.AgeStart = &defaultAgeStart
	}

	if ad.Conditions.AgeEnd == nil {
		defaultAgeEnd := uint8(100)
		ad.Conditions.AgeEnd = &defaultAgeEnd
	}

	// Marshal conditions to JSONB
	conditionsJSON, err := json.Marshal(ad.Conditions)
	if err != nil {
		return err
	}

	// Insert conditions
	_, err = insertAdStmt.Exec(ad.Title, ad.StartAt, ad.EndAt, []byte(conditionsJSON))
	if err != nil {
		return err
	}

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
func CreateAdHandler(c *gin.Context) {
	var ad Advertisement
	if err := c.BindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//get the db from context
	db, exists := c.Get("DB")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB not found in context"})
		return
	}

	// Insert the advertisement into the database
	err := CreateAd(db.(*sql.DB), &ad)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Write a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Advertisement created successfully"})
}

func GetAdHandler(c *gin.Context) {
	// Parse query parameters
	age := c.Query("age")
	gender := c.Query("gender")
	country := c.Query("country")
	platform := c.Query("platform")
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")

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

	// Define a key for Redis cache
	cacheKey := fmt.Sprintf("ad:%s:%s:%s:%s:%d:%d", age, gender, country, platform, offset, limit)

	// Check if data is available in Redis cache
	cachedData, err := redisClient.Get(c, cacheKey).Result()
	if err == nil {
		// Data found in cache, return cached data
		log.Print("Return cached data")
		c.Data(http.StatusOK, "application/json", []byte(cachedData))
		return
	} else if err != redis.Nil {
		// Error occurred while retrieving data from Redis
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Query database for advertisements based on conditions
	db, exists := c.Get("DB")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB not found in context"})
		return
	}

	query := `
		SELECT title, end_at
		FROM advertisement
		WHERE start_at < NOW() AND end_at > NOW()
	`

	// Initialize an empty slice to hold query parameters
	var queryParams []interface{}

	if age != "" {
		ageInt, _ := strconv.Atoi(age)
		query += " AND ($1 BETWEEN conditions->>'ageStart' AND conditions->>'ageEnd')"

		queryParams = append(queryParams, ageInt)
	}
	if gender != "" {
		query += " AND (NOT (conditions?'gender') OR conditions->>'gender' = $2)"
		queryParams = append(queryParams, gender)
	}
	if country != "" {
		jsonCountry := fmt.Sprintf(`["%s"]`, country)
		query += " AND (NOT (conditions?'country') OR conditions->'country' @> $3::jsonb)"
		queryParams = append(queryParams, jsonCountry)
	}
	if platform != "" {
		jsonPlatform := fmt.Sprintf(`["%s"]`, platform)
		query += " AND (NOT (conditions?'platform') OR conditions->'platform' @> $4::jsonb)"
		queryParams = append(queryParams, jsonPlatform)
	}

	// Add pagination to the query
	query += " ORDER BY end_at ASC OFFSET $5 LIMIT $6"
	queryParams = append(queryParams, offset, limit)

	// Execute the query with prepared statement and query parameters
	rows, err := db.(*sql.DB).Query(query, queryParams...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		advertisements = append(advertisements, ad)
	}

	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cache the response in Redis for future requests
	err = redisClient.Set(c, cacheKey, jsonResponse, time.Second).Err() // Cache for 10 minutes
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set response content type and write JSON response
	c.Data(http.StatusOK, "application/json", jsonResponse)
}

// start api with the pgdb and return a chi router
func StartAPI() *gin.Engine {

	db := dbConnect()

	InitializeRedis()

	// Initialize Gin router
	r := gin.Default()

	// Middleware to store DB in context
	r.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	})

	r.POST("/api/v1/ad", CreateAdHandler)

	r.GET("/api/v1/ad", GetAdHandler)

	return r
}
