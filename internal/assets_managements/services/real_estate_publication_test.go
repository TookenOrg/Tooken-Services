package services

import (
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/assets_managements/database"
)

// completeRequest is an asset a manager filled in one go: everything an
// investor needs to decide is there.
func completeRequest() server.RealEstateWriteRequest {
	req := validRequest()
	req.Description = strPtrT("A browsable listing")
	req.IssuerId = 1
	req.Media = &[]server.RealEstateMediaInput{{Url: "https://x/1.jpg"}}
	req.Configuration = &server.RealEstateConfigurationInput{
		TotalShares:   "1000",
		PricePerShare: "150",
		CurrencyCode:  "EUR",
		Yield:         strPtrT("4.25"),

		PaymentFrequency:       1,
		PaymentFrequencyTypeId: 1,
	}

	return req
}

func mustWriteDTO(t *testing.T, req server.RealEstateWriteRequest) database.RealEstateWriteDTO {
	t.Helper()

	in, err := toWriteDTO(req, writePatch)
	if err != nil {
		t.Fatalf("unexpected conversion error: %v", err)
	}

	return in
}

// Each requirement is checked on its own: a list that only works as a whole
// would hide the day one of its members stops being tested.
func TestPublicationRequirements(t *testing.T) {
	tests := []struct {
		name    string
		strip   func(*server.RealEstateWriteRequest)
		missing string
	}{
		{"complete", func(*server.RealEstateWriteRequest) {}, ""},
		{"no issuer", func(r *server.RealEstateWriteRequest) { r.IssuerId = 0 }, "issuer_id"},
		{"no description", func(r *server.RealEstateWriteRequest) { r.Description = nil }, "description"},
		{"blank description", func(r *server.RealEstateWriteRequest) { r.Description = strPtrT("   ") }, "description"},
		{"no visual", func(r *server.RealEstateWriteRequest) { r.Media = nil }, "media"},
		{"no configuration", func(r *server.RealEstateWriteRequest) { r.Configuration = nil }, "configuration.yield"},
		{"no yield", func(r *server.RealEstateWriteRequest) { r.Configuration.Yield = nil }, "configuration.yield"},
		{"no payment rhythm", func(r *server.RealEstateWriteRequest) { r.Configuration.PaymentFrequency = 0 }, "configuration.payment_frequency"},
		{"no payment period", func(r *server.RealEstateWriteRequest) { r.Configuration.PaymentFrequencyTypeId = 0 }, "configuration.payment_frequency_type_id"},
		{"blank yield", func(r *server.RealEstateWriteRequest) { r.Configuration.Yield = strPtrT("") }, "configuration.yield"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := completeRequest()
			tc.strip(&req)
			missing := publicationRequirements(mustWriteDTO(t, req))

			if tc.missing == "" {
				if len(missing) != 0 {
					t.Fatalf("a complete asset is reported incomplete: %v", missing)
				}
				return
			}

			if !contains(missing, tc.missing) {
				t.Fatalf("want %q reported missing, got %v", tc.missing, missing)
			}
		})
	}
}

// An asset created before the gallery existed carries imageurl instead, and
// that is enough: demanding media would make those assets unpublishable.
func TestImageurlCountsAsAVisual(t *testing.T) {
	req := completeRequest()
	req.Media = nil
	req.Imageurl = strPtrT("https://x/legacy.jpg")

	if missing := publicationRequirements(mustWriteDTO(t, req)); len(missing) != 0 {
		t.Fatalf("imageurl refused as a visual: %v", missing)
	}
}

// Surface area is comparison material, not a condition to understand an offer.
// Demanding it would push a manager to invent a number to get past the check.
func TestSurfaceIsNotRequiredToPublish(t *testing.T) {
	req := completeRequest()
	req.Specification = nil

	if missing := publicationRequirements(mustWriteDTO(t, req)); len(missing) != 0 {
		t.Fatalf("surface demanded to publish: %v", missing)
	}
}

func TestCheckStillPublishable(t *testing.T) {
	complete := mustWriteDTO(t, completeRequest())

	stripped := completeRequest()
	stripped.Description = nil
	incomplete := mustWriteDTO(t, stripped)

	t.Run("a draft can lose anything", func(t *testing.T) {
		if err := checkStillPublishable(database.StatusDraft, complete, incomplete); err != nil {
			t.Fatalf("a draft was refused: %v", err)
		}
	})

	t.Run("a published asset cannot lose a requirement", func(t *testing.T) {
		err := checkStillPublishable(database.StatusPublished, complete, incomplete)
		if err == nil {
			t.Fatal("a published asset lost its description silently")
		}
		if !strings.Contains(err.Error(), "description") {
			t.Errorf("the lost field is not named: %v", err)
		}
	})

	// The comparison is before/after, not the result alone: an asset published
	// before these rules existed must stay editable, including to complete it.
	t.Run("what was already missing is not held against the patch", func(t *testing.T) {
		if err := checkStillPublishable(database.StatusPublished, incomplete, incomplete); err != nil {
			t.Fatalf("an already incomplete published asset was frozen: %v", err)
		}
		if err := checkStillPublishable(database.StatusPublished, incomplete, complete); err != nil {
			t.Fatalf("completing a published asset was refused: %v", err)
		}
	})
}

// Regression: applyPatch used to write through the pointers of the request it
// was given, so the caller's "before" was already patched when compared to the
// "after". Clearing a required field then looked like no change at all, and a
// published asset silently lost it.
func TestApplyPatchDoesNotTouchItsInput(t *testing.T) {
	stored := completeRequest()

	blank := ""
	merged := applyPatch(stored, server.RealEstatePatchRequest{
		Configuration: &server.RealEstateConfigurationPatch{Yield: &blank},
		Specification: &server.RealEstateSpecificationInput{EnergyClass: strPtrT("A")},
		Media:         &[]server.RealEstateMediaInput{{Url: "https://x/2.jpg"}},
	})

	if stored.Configuration.Yield == nil || *stored.Configuration.Yield != "4.25" {
		t.Errorf("the stored configuration was modified: %q", derefString(stored.Configuration.Yield))
	}
	if merged.Configuration.Yield == nil || *merged.Configuration.Yield != "" {
		t.Errorf("the patch was not applied: %q", derefString(merged.Configuration.Yield))
	}
	if (*stored.Media)[0].Url != "https://x/1.jpg" {
		t.Errorf("the stored gallery was modified: %v", *stored.Media)
	}
}

func contains(values []string, want string) bool {
	for _, v := range values {
		if v == want {
			return true
		}
	}

	return false
}

func TestSamePtrInt(t *testing.T) {
	one, otherOne, two := 1, 1, 2

	cases := []struct {
		name string
		a, b *int
		want bool
	}{
		{"both nil", nil, nil, true},
		{"nil and set", nil, &one, false},
		{"set and nil", &one, nil, false},
		{"same value, distinct pointers", &one, &otherOne, true},
		{"different values", &one, &two, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := samePtrInt(tc.a, tc.b); got != tc.want {
				t.Fatalf("samePtrInt = %v, want %v", got, tc.want)
			}
		})
	}
}

// checkIssuerAttachment must decide without reading the database in every case
// where the answer cannot depend on the issuer's state. These tests run without
// a connection on purpose: a query in any of these paths would panic here,
// which is exactly the regression to catch — a round trip paid on every patch
// that does not move the asset onto another vehicle.
//
// A patch that DOES change the issuer now always queries, draft included: a
// dissolved vehicle no longer exists, so attaching even an invisible asset to it
// is refused. That path needs a database and is covered by the integration
// suite, not here.
func TestCheckIssuerAttachmentSkipsTheQuery(t *testing.T) {
	withIssuer := func(id *int) database.RealEstateWriteDTO {
		in := mustWriteDTO(t, completeRequest())
		in.IssuerId = id

		return in
	}

	one, otherOne := 1, 1

	t.Run("an unchanged issuer is not re-judged", func(t *testing.T) {
		// The value matters, not the pointer: applyPatch rebuilds the request,
		// so the merged DTO never shares its pointers with the stored one.
		if err := checkIssuerAttachment(t.Context(), database.StatusPublished, withIssuer(&one), withIssuer(&otherOne)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("an unchanged issuer on a draft is not re-judged either", func(t *testing.T) {
		if err := checkIssuerAttachment(t.Context(), database.StatusDraft, withIssuer(&one), withIssuer(&otherOne)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("clearing the issuer is checkStillPublishable's business", func(t *testing.T) {
		if err := checkIssuerAttachment(t.Context(), database.StatusPublished, withIssuer(&one), withIssuer(nil)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// issuerStandsBehind answers for an absent issuer without a query, so that a
// missing issuer is reported once, by publicationRequirements, and not a second
// time in words the caller cannot act on.
func TestIssuerStandsBehindWithoutIssuer(t *testing.T) {
	standing, err := issuerStandsBehind(t.Context(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !standing {
		t.Fatal("an absent issuer must not be reported as withdrawn")
	}
}
