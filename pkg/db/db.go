package db

import (
	"fmt"
	"log"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	// "github.com/avast/retry-go"
	// "github.com/go-pg/migrations"
	// "github.com/go-pg/pg"
	// _ "github.com/golang-migrate/migrate/v4/source/file"
	// _ "github.com/lib/pq" // Import PostgreSQL driver
)

// func StartDB() (*sqlx.DB, error) {
// 	var (
// 		db  *sqlx.DB
// 		err error
// 	)

// 	// Retry connecting to the database with exponential backoff
// 	err = retry.Do(
// 		func() error {
// 			// Define the database connection parameters
// 			dbURL := "postgres:admin@172.18.0.1:5435/postgres?sslmode=disable"
// 			// dsn := "host=127.0.0.1 port=5435 user=postgres password=admin dbname=postgres sslmode=disable"
// 			// Connect to the database
// 			db, err = sqlx.Connect("postgres", dbURL)

// 			// db, err = sqlx.Open("postgres", dbURL)
// 			if err != nil {
// 				return fmt.Errorf("failed to connect to the database: %v", err)
// 			}
// 			defer db.Close()

// 			// Check if the connection is successful
// 			ctx := context.Background()
// 			if err := db.PingContext(ctx); err != nil {
// 				return err // Retry will be attempted if there's an error
// 			}

// 			log.Println("Connected to the database")
// 			return nil // Return nil to indicate success
// 		},
// 		retry.Delay(1*time.Second),          // Wait 1 second between retries
// 		retry.Attempts(5),                   // Retry up to 5 times
// 		retry.DelayType(retry.BackOffDelay), // Use exponential backoff
// 	)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to the database after retries: %v", err)
// 	}

// 	// Run migrations after successful connection
// 	err = runMigrations(db)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to run migrations: %v", err)
// 	}
// 	fmt.Println("Migrations executed successfully")

// 	return db, nil
// }

func StartDB() (*pg.DB, error) {
	// Define the database connection parameters
	// dbURL := "postgresql://postgres:admin@localhost:5432/postgres?sslmode=disable"
	// // Connect to the database
	// db, err := pg.Connect("postgres", dbURL)
	// if err != nil {
	// 	log.Fatalf("failed to connect to the database: %v", err)
	// }
	// defer db.Close()

	// // Check if the connection is successful
	// err = db.Ping()
	// if err != nil {
	// 	log.Fatalf("failed to ping the database: %v", err)
	// }
	// log.Println("Connected to the database")

	// // Run migrations
	// err = runMigrations(db)
	// if err != nil {
	// 	log.Fatalf("failed to run migrations: %v", err)
	// }
	// fmt.Println("Migrations executed successfully")
	fmt.Println("HI!!!")

	var (
		opts *pg.Options
		err  error
	)

	//check if we are in prod
	//then use the db url from the env
	if os.Getenv("ENV") == "PROD" {
		// fmt.Print("DATABASE_URL: ", os.Getenv("DATABASE_URL"))
		opts, err = pg.ParseURL(os.Getenv("DATABASE_URL"))
		if err != nil {
			return nil, err
		}
	} else {
		opts = &pg.Options{
			//default port
			//depends on the db service from docker compose
			Addr:     "db:5432",
			User:     "postgres",
			Password: "admin",
		}
	}
	// Retry connecting to the database with exponential backoff

	// err = retry.Do(
	// 	func() error {
	// Your database connection logic
	db := pg.Connect(opts)
	// 		if db == nil {
	// 			return fmt.Errorf("failed to connect to the database")
	// 		}

	// 		// Check if the connection is successful
	// 		ctx := context.Background()
	// 		if err := db.Ping(ctx); err != nil {
	// 			return err // Retry will be attempted if there's an error
	// 		}
	// 		log.Printf("Success to connect to the database")
	// 		return nil // Return nil to indicate success
	// 	},
	// 	retry.Delay(1*time.Second),          // Wait 1 second between retries
	// 	retry.Attempts(5),                   // Retry up to 5 times
	// 	retry.DelayType(retry.BackOffDelay), // Use exponential backoff
	// )

	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to the database after retries: %v", err)
	// }

	// run migrations
	collection := migrations.NewCollection()
	err = collection.DiscoverSQLMigrations("migrations")
	if err != nil {
		return nil, err
	}

	//start the migrations
	_, _, err = collection.Run(db, "init")
	if err != nil {
		return nil, err
	}

	oldVersion, newVersion, err := collection.Run(db, "up")
	if err != nil {
		return nil, err
	}
	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}

	//return the db connection
	return db, err
}

// func runMigrations(db *sqlx.DB) error {
// 	// Initialize the driver for PostgreSQL
// 	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
// 	if err != nil {
// 		return err
// 	}

// 	// Create a new migrate instance
// 	m, err := migrate.NewWithDatabaseInstance(
// 		"Advertisement_Manage/migrations", // Path to migration files
// 		"postgres", driver)
// 	if err != nil {
// 		return err
// 	}

// 	// Migrate the database to the latest version
// 	err = m.Up()
// 	if err != nil && err != migrate.ErrNoChange {
// 		return err
// 	}
// 	log.Printf("Success to migrate the database")
// 	return nil
// }
