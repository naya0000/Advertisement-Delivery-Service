package app

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/naya0000/Advertisement_Manage.git/internal/app/api"
)

func Run() {
	log.Print("server has started")

	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	router, err := api.StartAPI()
	if err != nil {
		log.Fatalf("failed to start API: %v", err)
	}

	// Start listening with the server
	log.Printf("server is listening on port %s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Printf("error from router: %v\n", err)
		return
	}
}
