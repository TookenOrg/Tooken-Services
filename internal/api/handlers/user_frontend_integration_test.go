//go:build integration

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	blkGlobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/internal/testenv"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

func newUserFrontendHarness(t *testing.T) *orderHarness {
	t.Helper()
	db := testenv.ScratchDatabase(t, "tooken_it_user_frontend")
	oldDB := globals.DB
	globals.DB = db
	t.Cleanup(func() { globals.DB = oldDB })
	logger.Init(false)
	if err := middleware.InitAuth("../../../api/openapi.yaml"); err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group(globals.BaseURL)
	group.Use(middleware.AutoAuthMiddleware())
	handler := NewHandler()
	server.RegisterHandlers(group, handler)
	return &orderHarness{t: t, db: db, r: router, handler: handler, run: time.Now().UnixNano() % 1_000_000_000}
}

func TestUserFrontendEndpoints(t *testing.T) {
	h := newUserFrontendHarness(t)
	user := h.newUser(t, authUtils.RoleUser, "none", "NULL")
	admin := h.newUser(t, authUtils.RoleAdmin, "none", "NULL")
	manager := h.newUser(t, authUtils.RoleManager, "none", "NULL")
	favoritesPath := func(id int) string {
		return fmt.Sprintf("/users/%d/favorites/real-estates", id)
	}
	rolePath := func(id int) string { return fmt.Sprintf("/users/%d/role", id) }

	t.Run("all new routes require authentication", func(t *testing.T) {
		for _, route := range []struct{ method, path, body string }{
			{http.MethodGet, favoritesPath(user.id), ""},
			{http.MethodGet, "/users/me/kyc/onchain", ""},
			{http.MethodPatch, rolePath(user.id), `{"role":"MANAGER"}`},
		} {
			code, body := callAPI(h, route.method, route.path, "", route.body)
			expectCode(t, code, body, http.StatusUnauthorized)
		}
	})

	t.Run("favorites are private and an empty collection is an array", func(t *testing.T) {
		code, body := callAPI(h, http.MethodGet, favoritesPath(user.id), user.token, "")
		expectCode(t, code, body, http.StatusOK)
		if strings.TrimSpace(body) != "[]" {
			t.Fatalf("empty favorites = %s", body)
		}
		code, body = callAPI(h, http.MethodGet, favoritesPath(user.id), admin.token, "")
		expectCode(t, code, body, http.StatusForbidden)
	})

	t.Run("favorites retain catalogue visibility and exclude deleted assets", func(t *testing.T) {
		fixtures := newAssetFixtures(t, h)
		issuer := fixtures.activeIssuer(t, "Favorite listing")
		published := fixtures.asset(t, issuer, true)
		draft := fixtures.asset(t, issuer, false)
		deleted := fixtures.deletedAsset(t, issuer)
		unbookmarked := fixtures.asset(t, issuer, true)
		for _, caller := range []orderUser{user, manager} {
			for _, estateID := range []int{published, draft, deleted} {
				if _, err := h.db.Exec(`INSERT INTO rel.user_asset_favorite (user_id, real_estate_id) VALUES ($1, $2)`,
					caller.id, estateID); err != nil {
					t.Fatal(err)
				}
			}
			code, body := callAPI(h, http.MethodGet, favoritesPath(caller.id), caller.token, "")
			expectCode(t, code, body, http.StatusOK)
			var estates []server.RealEstateSummary
			if err := json.Unmarshal([]byte(body), &estates); err != nil {
				t.Fatal(err)
			}
			seen := map[int]bool{}
			for _, estate := range estates {
				seen[estate.Id] = true
			}
			if !seen[published] || seen[draft] != (caller.id == manager.id) || seen[deleted] || seen[unbookmarked] {
				t.Fatalf("unexpected favorites: %s", body)
			}
			if strings.Contains(body, `"street"`) || strings.Contains(body, `"address"`) {
				t.Fatalf("favorite summary contains an address: %s", body)
			}
		}
	})

	t.Run("on-chain KYC without a wallet is false without provisioning anything", func(t *testing.T) {
		code, body := callAPI(h, http.MethodGet, "/users/me/kyc/onchain", user.token, "")
		expectCode(t, code, body, http.StatusOK)
		var status server.KycOnChainStatus
		if err := json.Unmarshal([]byte(body), &status); err != nil {
			t.Fatal(err)
		}
		if status.UserId != user.id || status.IsVerified || status.WalletAddress != nil {
			t.Fatalf("unexpected status: %s", body)
		}
		if n := h.count(t, `SELECT count(*) FROM blk.user_wallet WHERE user_id = $1`, user.id); n != 0 {
			t.Fatal("reading KYC created a wallet")
		}
	})

	t.Run("an unavailable chain is 503 rather than a false KYC result", func(t *testing.T) {
		if _, err := h.db.Exec(`INSERT INTO blk.user_wallet (user_id, wallet_address, label)
			VALUES ($1, '0x0000000000000000000000000000000000000042', 'Main')`, user.id); err != nil {
			t.Fatal(err)
		}
		oldClient := blkGlobals.EthClient
		blkGlobals.EthClient = nil
		defer func() { blkGlobals.EthClient = oldClient }()
		code, body := callAPI(h, http.MethodGet, "/users/me/kyc/onchain", user.token, "")
		expectCode(t, code, body, http.StatusServiceUnavailable)
	})

	t.Run("only an administrator can manage roles", func(t *testing.T) {
		for _, caller := range []orderUser{user, manager} {
			code, body := callAPI(h, http.MethodPatch, rolePath(user.id), caller.token, `{"role":"ADMIN"}`)
			expectCode(t, code, body, http.StatusForbidden)
		}
	})

	t.Run("invalid role payloads do not change the user", func(t *testing.T) {
		for _, payload := range []string{
			``, `{}`, `null`, `{"role":null}`, `{"role":"owner"}`, `{"role":"user"}`,
			`{"role":1}`, `{"role":"MANAGER","email":"other@example.test"}`,
			`{"role":"MANAGER"} {}`, `[]`,
		} {
			code, body := callAPI(h, http.MethodPatch, rolePath(user.id), admin.token, payload)
			expectCode(t, code, body, http.StatusBadRequest)
		}
		if n := h.count(t, `SELECT count(*) FROM usr.users WHERE id = $1 AND role = 'USER'`, user.id); n != 1 {
			t.Fatal("invalid request changed the user's role")
		}
	})

	t.Run("the last administrator cannot self-demote", func(t *testing.T) {
		code, body := callAPI(h, http.MethodPatch, rolePath(admin.id), admin.token, `{"role":"USER"}`)
		expectCode(t, code, body, http.StatusConflict)
	})

	t.Run("role changes return the updated safe profile and unknown users are 404", func(t *testing.T) {
		for _, role := range []string{"MANAGER", "ADMIN", "USER"} {
			code, body := callAPI(h, http.MethodPatch, rolePath(user.id), admin.token, `{"role":"`+role+`"}`)
			expectCode(t, code, body, http.StatusOK)
			var updated server.User
			if err := json.Unmarshal([]byte(body), &updated); err != nil {
				t.Fatal(err)
			}
			if updated.Id != user.id || updated.Role == nil || string(*updated.Role) != role {
				t.Fatalf("unexpected user: %s", body)
			}
			if strings.Contains(body, "password") {
				t.Fatalf("updated user leaks password: %s", body)
			}
		}
		code, body := callAPI(h, http.MethodPatch, rolePath(999_999_991), admin.token, `{"role":"USER"}`)
		expectCode(t, code, body, http.StatusNotFound)
	})

	t.Run("self-demotion stays forbidden when another administrator exists", func(t *testing.T) {
		otherAdmin := h.newUser(t, authUtils.RoleAdmin, "none", "NULL")
		code, body := callAPI(h, http.MethodPatch, rolePath(admin.id), admin.token, `{"role":"MANAGER"}`)
		expectCode(t, code, body, http.StatusConflict)

		code, body = callAPI(h, http.MethodPatch, rolePath(otherAdmin.id), admin.token, `{"role":"USER"}`)
		expectCode(t, code, body, http.StatusOK)
		code, body = callAPI(h, http.MethodPatch, rolePath(user.id), otherAdmin.token, `{"role":"ADMIN"}`)
		expectCode(t, code, body, http.StatusForbidden)
		if n := h.count(t, `SELECT count(*) FROM usr.users WHERE id = $1 AND role = 'USER'`, user.id); n != 1 {
			t.Fatal("demoted administrator changed a role with a stale JWT")
		}
	})
}
