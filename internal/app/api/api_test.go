package api

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestStartAPI(t *testing.T) {
	tests := []struct {
		name    string
		want    *gin.Engine
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StartAPI()
			if (err != nil) != tt.wantErr {
				t.Errorf("StartAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}
