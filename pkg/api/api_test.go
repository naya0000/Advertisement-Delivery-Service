package api

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestCreateAd(t *testing.T) {
	type args struct {
		db *sql.DB
		ad *Advertisement
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateAd(tt.args.db, tt.args.ad); (err != nil) != tt.wantErr {
				t.Errorf("CreateAd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
