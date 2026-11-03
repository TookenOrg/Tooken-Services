package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/assets_managements/services"
	"github.com/gin-gonic/gin"
)

func TestRealEstateConflictMissingFields(t *testing.T) {
	fields := []string{"description", "configuration.total_shares"}
	for _, tc := range []struct {
		name string
		err  error
		want []string
	}{
		{"incomplete", fmt.Errorf("publish: %w", &services.MissingFieldsConflict{Message: "incomplete", Fields: fields}), fields},
		{"other conflict", services.ErrRealEstateConflict, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			respondRealEstateError(ctx, tc.err, "Unable to publish.")
			if recorder.Code != http.StatusConflict {
				t.Fatalf("status = %d", recorder.Code)
			}
			var body map[string]json.RawMessage
			if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			var got []string
			if raw, exists := body["missing_fields"]; exists {
				if err := json.Unmarshal(raw, &got); err != nil {
					t.Fatal(err)
				}
				if tc.want == nil {
					t.Fatal("unrelated conflict includes missing_fields")
				}
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("missing_fields = %v, want %v", got, tc.want)
			}
			if _, exists := body["message"]; !exists {
				t.Fatal("existing message field was lost")
			}
		})
	}
}
