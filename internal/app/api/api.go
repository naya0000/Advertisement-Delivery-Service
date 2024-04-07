package api

import (
	"github.com/gin-gonic/gin"
	"github.com/naya0000/Advertisement_Manage.git/internal/pkg/database"
)

// StartAPI initializes and returns the Gin router
func StartAPI() (*gin.Engine, error) {
	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	// Initialize Gin router
	r := gin.Default()

	// Middleware to store DB in context
	r.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	})

	r.POST("/api/v1/ad", CreateAdHandler)
	r.GET("/api/v1/ad", GetAdHandler)

	return r, nil
}
