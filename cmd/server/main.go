package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/naya0000/Advertisement_Manage.git/pkg/api"
	"github.com/naya0000/Advertisement_Manage.git/pkg/db"
)

func main() {
	log.Print("server has started")
	//start the db
	db, err := db.StartDB()
	if err != nil {
		log.Printf("error starting the database %v", err)
		panic("error starting the database")
	}
	//get the router of the API by passing the db
	router := api.StartAPI(db)
	//get the port from the environment variable
	port := os.Getenv("PORT")
	//pass the router and start listening with the server
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Printf("error from router %v\n", err)
		return
	}
	// http.HandleFunc("/api/v1/ad", CreateAdvertisementHandler)
	// http.HandleFunc("/api/v1/getAd", ListAdvertisementHandler)
	// http.ListenAndServe(":8080", nil)
}
