package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDecodeUserRoleRequest(t *testing.T) {
	for _, tc := range []struct {
		body string
		ok   bool
	}{
		{`{"role":"MANAGER"}`, true},
		{`{"role":"USER","email":"other@example.test"}`, false},
		{`{"role":"USER"} {"role":"ADMIN"}`, false},
		{`{"role":123}`, false},
		{`[]`, false},
		{`{`, false},
		{``, false},
	} {
		t.Run(tc.body, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = httptest.NewRequest("PATCH", "/", strings.NewReader(tc.body))
			_, err := decodeUserRoleRequest(ctx)
			if (err == nil) != tc.ok {
				t.Fatalf("decode error = %v, valid = %v", err, tc.ok)
			}
		})
	}
}
