package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/router"
	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/gin-gonic/gin"
)

type trexInitSpy struct {
	mu    sync.Mutex
	calls []string
}

func (s *trexInitSpy) record(call string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, call)
}

func (s *trexInitSpy) recorded() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.calls...)
}

func (s *trexInitSpy) DeployAllImplementations(context.Context) ([]server.ContractDetails, error) {
	s.record("implementations")
	return nil, nil
}

func (s *trexInitSpy) DeployIdentityFactory(context.Context) (server.ContractDetails, error) {
	s.record("identity factory")
	return server.ContractDetails{}, nil
}

func (s *trexInitSpy) ConfigureAuthority(context.Context) (server.TxHashName, error) {
	s.record("authority")
	return server.TxHashName{}, nil
}

func (s *trexInitSpy) DeployAndInitTrexFactory(context.Context) error {
	s.record("TREX factory")
	return nil
}

func (s *trexInitSpy) DeployTrexSuite(context.Context) ([]server.ContractDetails, error) {
	s.record("suite")
	return nil, nil
}

func newTrexInitTestEngine(t *testing.T) (*gin.Engine, *Handler, *trexInitSpy) {
	t.Helper()

	spy := &trexInitSpy{}
	h := NewHandler()
	h.trexInitSvc = spy
	t.Cleanup(h.WaitBackground)

	engine, err := router.New(h, "../../../api/openapi.yaml", []string{"http://localhost:3000"})
	if err != nil {
		t.Fatalf("router.New: %v", err)
	}
	return engine, h, spy
}

func trexInitToken(t *testing.T, role string) string {
	t.Helper()
	token, err := authUtils.GenerateJWT(42, "staff@tooken.test", role)
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}
	return token
}

func callTrexInit(engine *gin.Engine, h *Handler, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, globals.BaseURL+path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, req)
	h.WaitBackground()
	return response
}

func TestTrexInitializationRequiresExplicitEnablement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	roles := []struct {
		name  string
		token string
		staff bool
	}{
		{name: "anonymous"},
		{name: "USER", token: trexInitToken(t, authUtils.RoleUser)},
		{name: "MANAGER", token: trexInitToken(t, authUtils.RoleManager), staff: true},
		{name: "ADMIN", token: trexInitToken(t, authUtils.RoleAdmin), staff: true},
	}
	routes := []struct {
		path       string
		operation  string
		wantStatus int
	}{
		{"/trex/init/implementations", "implementations", http.StatusAccepted},
		{"/trex/init/contract/factory/identity", "identity factory", http.StatusOK},
		{"/trex/init/authority/configure", "authority", http.StatusOK},
		{"/trex/init/contract/factory/trexFactory", "TREX factory", http.StatusAccepted},
		{"/trex/init/suite/deploy", "suite", http.StatusAccepted},
	}
	configurations := []struct {
		name    string
		value   string
		unset   bool
		enabled bool
	}{
		{name: "absent", unset: true},
		{name: "empty"},
		{name: "false", value: "false"},
		{name: "invalid", value: "ture"},
		{name: "uppercase", value: "TRUE"},
		{name: "mixed case", value: "True"},
		{name: "leading space", value: " true"},
		{name: "trailing space", value: "true "},
		{name: "trailing newline", value: "true\n"},
		{name: "numeric", value: "1"},
		{name: "yes", value: "yes"},
		{name: "enabled", value: "true", enabled: true},
	}

	for _, cfg := range configurations {
		t.Run(cfg.name, func(t *testing.T) {
			t.Setenv("TREX_INIT_ENABLED", cfg.value)
			if cfg.unset {
				if err := os.Unsetenv("TREX_INIT_ENABLED"); err != nil {
					t.Fatal(err)
				}
			}
			engine, h, spy := newTrexInitTestEngine(t)

			for _, route := range routes {
				for _, role := range roles {
					before := len(spy.recorded())
					response := callTrexInit(engine, h, route.path, role.token)

					wantStatus := http.StatusForbidden
					if role.name == "anonymous" {
						wantStatus = http.StatusUnauthorized
					} else if cfg.enabled && role.staff {
						wantStatus = route.wantStatus
					}
					if response.Code != wantStatus {
						t.Errorf("%s as %s: status = %d, want %d; body = %s",
							route.path, role.name, response.Code, wantStatus, response.Body.String())
					}

					calls := spy.recorded()
					if cfg.enabled && role.staff {
						if len(calls) != before+1 || calls[len(calls)-1] != route.operation {
							t.Errorf("%s as %s: calls = %v, want one %q call",
								route.path, role.name, calls[before:], route.operation)
						}
					} else if len(calls) != before {
						t.Errorf("%s as %s unexpectedly started %v", route.path, role.name, calls[before:])
					}

					if role.name == "anonymous" || role.name == "USER" || (role.staff && !cfg.enabled) {
						var body server.APIResponse
						if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
							t.Fatalf("%s: invalid JSON: %v", route.path, err)
						}
						wantMessage := "TREX initialization is disabled."
						switch role.name {
						case "anonymous":
							wantMessage = "Authorization header required"
						case "USER":
							wantMessage = "Your account is not allowed to perform this operation."
						}
						if body.Message != wantMessage {
							t.Errorf("%s as %s: message = %q, want %q", route.path, role.name, body.Message, wantMessage)
						}
					}
				}
			}

			if !cfg.enabled {
				for _, path := range []string{"/health", "/trex/suite"} {
					req := httptest.NewRequest(http.MethodGet, globals.BaseURL+path, nil)
					response := httptest.NewRecorder()
					engine.ServeHTTP(response, req)
					if response.Code != http.StatusOK {
						t.Errorf("GET %s with initialization disabled = %d, want 200", path, response.Code)
					}
				}
			}
		})
	}
}

func TestTrexInitializationSettingIsCapturedAtStartup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token := trexInitToken(t, authUtils.RoleManager)
	for _, tc := range []struct {
		name       string
		initial    string
		later      string
		unset      bool
		wantStatus int
		wantCalls  int
	}{
		{"closed stays closed", "false", "true", false, http.StatusForbidden, 0},
		{"open stays open", "true", "false", false, http.StatusAccepted, 1},
		{"open stays open after removal", "true", "", true, http.StatusAccepted, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TREX_INIT_ENABLED", tc.initial)
			engine, h, spy := newTrexInitTestEngine(t)
			t.Setenv("TREX_INIT_ENABLED", tc.later)
			if tc.unset {
				if err := os.Unsetenv("TREX_INIT_ENABLED"); err != nil {
					t.Fatal(err)
				}
			}

			response := callTrexInit(engine, h, "/trex/init/implementations", token)
			if response.Code != tc.wantStatus {
				t.Errorf("status after environment change = %d, want %d", response.Code, tc.wantStatus)
			}
			if calls := spy.recorded(); len(calls) != tc.wantCalls {
				t.Errorf("calls after environment change = %v, want %d", calls, tc.wantCalls)
			}
		})
	}
}

func TestNewHandlerSharesBlockchainServiceForTrexInitialization(t *testing.T) {
	t.Setenv("TREX_INIT_ENABLED", "true")
	h := NewHandler()
	if h.blockchainSvc == nil || h.trexInitSvc != h.blockchainSvc {
		t.Fatal("TREX initialization must use the shared blockchain service")
	}
}
