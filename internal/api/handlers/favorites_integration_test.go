//go:build integration

// HTTP coverage of /users/{userId}/favorites/real-estates/{realEstateId}:
// an investor bookmarking a real estate, then removing the bookmark. Only the
// investor themselves may do either.
package handlers

import (
	"fmt"
	"net/http"
	"testing"

	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	_ "github.com/lib/pq"
)

func favoritePath(user, realEstate int) string {
	return fmt.Sprintf("/users/%d/favorites/real-estates/%d", user, realEstate)
}

func TestFavoriteRealEstates(t *testing.T) {
	h := newOrderHarness(t)
	fixtures := newAssetFixtures(t, h)
	issuer := fixtures.activeIssuer(t, "Favorites")
	villa := fixtures.asset(t, issuer, true)
	loft := fixtures.asset(t, issuer, true)

	favorites := func(t *testing.T, user, realEstate int) int {
		t.Helper()
		return h.count(t, `
		    SELECT count(*) FROM rel.user_asset_favorite
		    WHERE user_id = $1 AND real_estate_id = $2`, user, realEstate)
	}

	t.Run("a bookmarked real estate is stored once", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)

		code, body := callAPI(h, http.MethodPost, favoritePath(investor.id, villa), investor.token, "")
		expectCode(t, code, body, http.StatusOK)
		expectAPIMessage(t, body, "Favorite added successfully.")
		if n := favorites(t, investor.id, villa); n != 1 {
			t.Fatalf("%d favorite row(s), want 1", n)
		}
	})

	t.Run("bookmarking the same real estate twice keeps a single row", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)
		code, body := callAPI(h, http.MethodPost, favoritePath(investor.id, villa), investor.token, "")
		expectCode(t, code, body, http.StatusOK)

		code, body = callAPI(h, http.MethodPost, favoritePath(investor.id, villa), investor.token, "")
		expectCode(t, code, body, http.StatusConflict)
		expectAPIMessage(t, body, "already is in your favorites")
		if n := favorites(t, investor.id, villa); n != 1 {
			t.Fatalf("%d favorite row(s), want 1", n)
		}
	})

	t.Run("a real estate that does not exist cannot be bookmarked", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)

		code, body := callAPI(h, http.MethodPost, favoritePath(investor.id, unknownID), investor.token, "")
		expectCode(t, code, body, http.StatusNotFound)
		expectAPIMessage(t, body, "Real estate not found.")
		if n := h.count(t, `SELECT count(*) FROM rel.user_asset_favorite WHERE user_id = $1`, investor.id); n != 0 {
			t.Fatalf("%d favorite row(s) written for a missing real estate", n)
		}
	})

	t.Run("removing a bookmark deletes that one only", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)
		other := newTrackedUser(t, h, authUtils.RoleUser)
		for _, add := range []struct {
			who        orderUser
			realEstate int
		}{
			{investor, villa}, {investor, loft}, {other, villa},
		} {
			code, body := callAPI(h, http.MethodPost, favoritePath(add.who.id, add.realEstate), add.who.token, "")
			expectCode(t, code, body, http.StatusOK)
		}

		code, body := callAPI(h, http.MethodDelete, favoritePath(investor.id, villa), investor.token, "")
		expectCode(t, code, body, http.StatusOK)
		expectAPIMessage(t, body, "Favorite removed successfully.")

		if n := favorites(t, investor.id, villa); n != 0 {
			t.Fatalf("the removed favorite is still stored")
		}
		if favorites(t, investor.id, loft) != 1 || favorites(t, other.id, villa) != 1 {
			t.Fatalf("removing one favorite deleted another one")
		}
	})

	// Answering "removed" for a bookmark that was never there would hide a
	// wrong identifier from the caller.
	t.Run("removing a bookmark that does not exist is refused and touches nothing", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)
		code, body := callAPI(h, http.MethodPost, favoritePath(investor.id, loft), investor.token, "")
		expectCode(t, code, body, http.StatusOK)

		code, body = callAPI(h, http.MethodDelete, favoritePath(investor.id, villa), investor.token, "")
		expectCode(t, code, body, http.StatusNotFound)
		expectAPIMessage(t, body, "not in your favorites")
		if n := favorites(t, investor.id, loft); n != 1 {
			t.Fatalf("the existing favorite was removed")
		}
	})

	t.Run("an anonymous caller changes nothing", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)
		code, body := callAPI(h, http.MethodPost, favoritePath(investor.id, loft), investor.token, "")
		expectCode(t, code, body, http.StatusOK)

		code, body = callAPI(h, http.MethodPost, favoritePath(investor.id, villa), "", "")
		expectCode(t, code, body, http.StatusUnauthorized)
		expectAPIMessage(t, body, "Authorization header required")
		code, body = callAPI(h, http.MethodDelete, favoritePath(investor.id, loft), "", "")
		expectCode(t, code, body, http.StatusUnauthorized)

		if favorites(t, investor.id, villa) != 0 || favorites(t, investor.id, loft) != 1 {
			t.Fatalf("an anonymous request changed the favorites")
		}
	})

	// The path names the user, the token proves who calls: the two must be the
	// same person, whatever the role. An administrator has no business in the
	// favorites of an investor either.
	for _, role := range []string{authUtils.RoleUser, authUtils.RoleAdmin} {
		t.Run("a "+role+" cannot change the favorites of somebody else", func(t *testing.T) {
			victim := newTrackedUser(t, h, authUtils.RoleUser)
			intruder := newTrackedUser(t, h, role)
			code, body := callAPI(h, http.MethodPost, favoritePath(victim.id, loft), victim.token, "")
			expectCode(t, code, body, http.StatusOK)

			code, body = callAPI(h, http.MethodPost, favoritePath(victim.id, villa), intruder.token, "")
			expectCode(t, code, body, http.StatusForbidden)
			expectAPIMessage(t, body, "your own favorites")
			code, body = callAPI(h, http.MethodDelete, favoritePath(victim.id, loft), intruder.token, "")
			expectCode(t, code, body, http.StatusForbidden)

			if favorites(t, victim.id, villa) != 0 || favorites(t, victim.id, loft) != 1 {
				t.Fatalf("another user changed the favorites")
			}
			if n := h.count(t, `SELECT count(*) FROM rel.user_asset_favorite WHERE user_id = $1`, intruder.id); n != 0 {
				t.Fatalf("%d favorite row(s) written for the intruder", n)
			}
		})
	}

	// The token is still valid, but it points at nobody: signing in again is
	// what the client must be told, not that the real estate is missing.
	t.Run("a token whose account no longer exists is sent back to sign in", func(t *testing.T) {
		ghost := h.newUser(t, authUtils.RoleUser, "none", "NULL")
		if _, err := h.db.Exec(`DELETE FROM usr.users WHERE id = $1`, ghost.id); err != nil {
			t.Fatal(err)
		}

		code, body := callAPI(h, http.MethodPost, favoritePath(ghost.id, villa), ghost.token, "")
		expectCode(t, code, body, http.StatusUnauthorized)
		expectAPIMessage(t, body, "Sign in again")
		if n := h.count(t, `SELECT count(*) FROM rel.user_asset_favorite WHERE user_id = $1`, ghost.id); n != 0 {
			t.Fatalf("%d favorite row(s) written for a deleted account", n)
		}

		// Deleting the account removed its favorites, so there is nothing
		// left to remove: the answer is the one of any missing favorite.
		code, body = callAPI(h, http.MethodDelete, favoritePath(ghost.id, villa), ghost.token, "")
		expectCode(t, code, body, http.StatusNotFound)
		expectAPIMessage(t, body, "not in your favorites")
	})

	// A failure the caller cannot act on is a 500 with a generic message: the
	// database error stays in the logs.
	t.Run("a database failure is a 500 that does not leak the error", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)
		withClosedDatabase(t, func() {
			for _, method := range []string{http.MethodPost, http.MethodDelete} {
				code, body := callAPI(h, method, favoritePath(investor.id, villa), investor.token, "")
				expectCode(t, code, body, http.StatusInternalServerError)
				expectGenericFailure(t, body)
			}
		})
	})
}
