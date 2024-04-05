package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func StartDB() (*sql.DB, error) {
	fmt.Println("HI!!!")

	db, err := sql.Open("postgres", "postgres://postgres:admin@db:5432/postgres?sslmode=disable")
	checkErr(err)
	// defer db.Close()

	//return the db connection
	return db, err
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
