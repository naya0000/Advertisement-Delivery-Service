package db

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func StartDB() (*bun.DB, error) {
	fmt.Println("HI!!!")
	var err error
	var db *bun.DB
	err = retry.Do(
		func() error {
			config, err := pgx.ParseConfig("postgres://postgres:admin@db:5432/postgres?sslmode=disable")
			if err != nil {
				panic(err)
			}
			config.PreferSimpleProtocol = true

			sqldb := stdlib.OpenDB(*config)
			db = bun.NewDB(sqldb, pgdialect.New())
			// dsn := "postgres://postgres:admin@db:5432/postgres?sslmode=disable"
			// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
			// pgconn := pgdriver.NewConnector(
			// 	// pgdriver.WithNetwork("tcp"),
			// 	pgdriver.WithAddr("db:5432"),
			// 	// pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
			// 	pgdriver.WithUser("postgres"),
			// 	pgdriver.WithPassword("admin"),
			// 	// pgdriver.WithDatabase("postgres"),
			// 	// pgdriver.WithApplicationName("myapp"),
			// 	// pgdriver.WithTimeout(5 * time.Second),
			// 	// pgdriver.WithDialTimeout(5 * time.Second),
			// 	// pgdriver.WithReadTimeout(5 * time.Second),
			// 	// pgdriver.WithWriteTimeout(5 * time.Second),
			// 	// pgdriver.WithConnParams(map[string]interface{}{
			// 	// 	"search_path": "my_search_path",
			// 	// }),
			// )
			// sqldb := sql.OpenDB(pgconn)
			// sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
			// db = bun.NewDB(sqldb, pgdialect.New())

			// Check if the connection is successful
			ctx := context.Background()
			if err := db.PingContext(ctx); err != nil {
				return err
			}
			fmt.Print("Success to connect to DB.")
			return nil
		},
		retry.Delay(1*time.Second),          // Wait 1 second between retries
		retry.Attempts(3),                   // Retry up to 5 times
		retry.DelayType(retry.BackOffDelay), // Use exponential backoff
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database after retries: %v", err)
	}
	// A collection of migrations.
	// var Migrations = migrate.NewMigrations()
	// var sqlMigrations embed.FS
	// if err := Migrations.Discover(sqlMigrations); err != nil {
	// 	panic(err)
	// }

	// run migrations
	// collection := migrations.NewCollection()
	// err = collection.DiscoverSQLMigrations("migrations")
	// if err != nil {
	// 	return nil, err
	// }

	// //start the migrations
	// _, _, err = collection.Run(db, "init")
	// if err != nil {
	// 	return nil, err
	// }

	// oldVersion, newVersion, err := collection.Run(db, "up")
	// if err != nil {
	// 	return nil, err
	// }
	// if newVersion != oldVersion {
	// 	fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	// } else {
	// 	fmt.Printf("version is %d\n", oldVersion)
	// }

	return db, err
}

// run migrations
// collection := migrations.NewCollection()
// err = collection.DiscoverSQLMigrations("migrations")
// if err != nil {
// 	return nil, err
// }

// //start the migrations
// _, _, err = collection.Run(db, "init")
// if err != nil {
// 	return nil, err
// }

// oldVersion, newVersion, err := collection.Run(db, "up")
// if err != nil {
// 	return nil, err
// }
// if newVersion != oldVersion {
// 	log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
// } else {
// 	log.Printf("version is %d\n", oldVersion)
// }

//return the db connection

// }

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
