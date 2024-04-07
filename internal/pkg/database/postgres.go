package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/naya0000/Advertisement_Manage.git/internal/pkg/models"
)

func Connect() (*sql.DB, error) {
	connStr := "postgres://postgres:admin@db:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

func GetAd(db *sql.DB, params models.QueryParams) (*sql.Rows, error) {
	query := `
		SELECT title, end_at
		FROM advertisement
		WHERE start_at < NOW() AND end_at > NOW()
	`

	// Initialize an empty slice to hold query parameters
	var queryParams []interface{}

	if params.Age != "" {
		ageInt, err := strconv.Atoi(params.Age)
		if err != nil {
			return nil, fmt.Errorf("invalid age parameter: %v", err)
		}
		query += " AND ($1 BETWEEN conditions->'ageStart' AND conditions->'ageEnd')"
		queryParams = append(queryParams, ageInt)
	}
	if params.Gender != "" {
		query += " AND (NOT (conditions?'gender') OR conditions->>'gender' = $2)"
		queryParams = append(queryParams, params.Gender)
	}
	if params.Country != "" {
		jsonCountry := fmt.Sprintf(`["%s"]`, params.Country)
		query += " AND (NOT (conditions?'country') OR conditions->'country' @> $3::jsonb)"
		queryParams = append(queryParams, jsonCountry)
	}
	if params.Platform != "" {
		jsonPlatform := fmt.Sprintf(`["%s"]`, params.Platform)
		query += " AND (NOT (conditions?'platform') OR conditions->'platform' @> $4::jsonb)"
		queryParams = append(queryParams, jsonPlatform)
	}

	// Add pagination to the query
	query += " ORDER BY end_at ASC OFFSET $5 LIMIT $6"
	queryParams = append(queryParams, params.Offset, params.Limit)

	rows, err := db.Query(query, queryParams...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	return rows, nil
}

func CreateAd(db *sql.DB, ad *models.Advertisement) error {
	// starts a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	insertAdStmt, err := tx.Prepare(`
		INSERT INTO advertisement (title, start_at, end_at, conditions) 
		VALUES($1, $2, $3, $4) 
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
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
		return fmt.Errorf("failed to marshal conditions to JSON: %v", err)
	}

	// Insert conditions
	_, err = insertAdStmt.Exec(ad.Title, ad.StartAt, ad.EndAt, []byte(conditionsJSON))
	if err != nil {
		return fmt.Errorf("failed to execute insert statement: %v", err)
	}

	// commits the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
