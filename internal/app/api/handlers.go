package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/naya0000/Advertisement_Manage.git/internal/pkg/database"
	"github.com/naya0000/Advertisement_Manage.git/internal/pkg/models"
	"github.com/naya0000/Advertisement_Manage.git/internal/pkg/redisDB"
)

// Handles creating advertisements
func CreateAdHandler(c *gin.Context) {
	var ad models.Advertisement
	if err := c.BindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate advertisement start and end time
	if ad.StartAt.After(ad.EndAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}

	//get the db from context
	db, exists := c.Get("DB")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB not found in context"})
		return
	}

	// Insert the advertisement into the database
	err := database.CreateAd(db.(*sql.DB), &ad)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Write a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Advertisement created successfully"})
}

// Handles retrieving advertisements
func GetAdHandler(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "5")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	queryParams := models.QueryParams{
		Age:      c.Query("age"),
		Gender:   c.Query("gender"),
		Country:  c.Query("country"),
		Platform: c.Query("platform"),
		Offset:   offset,
		Limit:    limit,
	}

	// Validate age parameter
	if queryParams.Age != "" {
		_, err := strconv.Atoi(queryParams.Age)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age parameter"})
			return
		}
	}

	// Define a key for Redis cache
	cacheKey := fmt.Sprintf("ad:%s:%s:%s:%s:%d:%d", queryParams.Age, queryParams.Gender, queryParams.Country, queryParams.Platform, queryParams.Offset, queryParams.Limit)

	// Check if data is available in Redis cache
	cachedData, err := redisDB.GetCacheData(c, cacheKey)
	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cached data"})
		return
	}

	if cachedData != "" {
		// Data found in cache
		log.Print("Return cached data")
		c.Data(http.StatusOK, "application/json", []byte(cachedData))
		return
	}

	// Query database for advertisements based on conditions
	db, exists := c.Get("DB")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB not found in context"})
		return
	}
	rows, err := database.GetAd(db.(*sql.DB), queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer rows.Close()

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan advertisement data"})
			return
		}
		advertisements = append(advertisements, ad)
	}

	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate through advertisement data"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal response"})
		return
	}

	// Cache the response in Redis for future requests
	err = redisDB.SetCacheData(c, cacheKey, jsonResponse, 10*time.Second) // Cache for 10 minutes
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache response"})
		return
	}

	// Set response content type and write JSON response
	c.Data(http.StatusOK, "application/json", jsonResponse)
}
