package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	// "github.com/go-pg/migrations/v8"
	// "github.com/go-pg/pg/v10"
	// "github.com/avast/retry-go"
	// "github.com/go-pg/migrations"
	// "github.com/go-pg/pg"
	// _ "github.com/golang-migrate/migrate/v4/source/file"
	// _ "github.com/lib/pq" // Import PostgreSQL driver
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
