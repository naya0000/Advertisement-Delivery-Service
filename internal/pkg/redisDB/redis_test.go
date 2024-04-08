package redisDB

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestInitClient(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitClient(); (err != nil) != tt.wantErr {
				t.Errorf("InitClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetClient(t *testing.T) {
	tests := []struct {
		name    string
		want    *redis.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCacheData(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCacheData(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCacheData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCacheData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetCacheData(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		value      interface{}
		expiration time.Duration
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
			if err := SetCacheData(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("SetCacheData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
