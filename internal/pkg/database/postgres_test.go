package database

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
	"github.com/naya0000/Advertisement_Manage.git/internal/pkg/models"
)

func TestCreateAd(t *testing.T) {
	type args struct {
		db *sql.DB
		ad *models.Advertisement
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

func TestGetAd(t *testing.T) {
	type args struct {
		db     *sql.DB
		params models.QueryParams
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.Rows
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAd(tt.args.db, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnect(t *testing.T) {
	tests := []struct {
		name    string
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Connect()
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connect() = %v, want %v", got, tt.want)
			}
		})
	}
}
