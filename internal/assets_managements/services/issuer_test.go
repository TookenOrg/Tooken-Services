package services

import (
	"errors"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/assets_managements/database"
	"github.com/lib/pq"
)

func validIssuerRequest() server.IssuerWriteRequest {
	return server.IssuerWriteRequest{
		Name:        "Tooken Real Estate S.C.Sp.",
		LegalForm:   "SCSp",
		CountryCode: "LU",
	}
}

// Every value a client can get wrong has to come back named. An issuer stored
// with a jurisdiction nobody typed is exactly the class of plausible-and-false
// record the currency default produced (§11.20).
func TestToIssuerWriteDTORejectsInvalidPayloads(t *testing.T) {
	tests := []struct {
		name string
		mut  func(*server.IssuerWriteRequest)
	}{
		{"blank name", func(r *server.IssuerWriteRequest) { r.Name = "   " }},
		{"blank legal form", func(r *server.IssuerWriteRequest) { r.LegalForm = "" }},
		{"missing country", func(r *server.IssuerWriteRequest) { r.CountryCode = "" }},
		{"three letter country", func(r *server.IssuerWriteRequest) { r.CountryCode = "LUX" }},
		{"name too long", func(r *server.IssuerWriteRequest) { r.Name = strings.Repeat("é", 256) }},
		{"legal form too long", func(r *server.IssuerWriteRequest) { r.LegalForm = strings.Repeat("a", 101) }},
		{"registration number too long", func(r *server.IssuerWriteRequest) {
			r.RegistrationNumber = strPtrT(strings.Repeat("B", 65))
		}},
		{"short lei", func(r *server.IssuerWriteRequest) { r.LeiCode = strPtrT("ABC123") }},
		{"lei with punctuation", func(r *server.IssuerWriteRequest) {
			r.LeiCode = strPtrT("2138001-2345678901A")
		}},
		{"unknown status", func(r *server.IssuerWriteRequest) { r.StatusId = intPtrT(9) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := validIssuerRequest()
			tc.mut(&req)

			if _, err := toIssuerWriteDTO(req); !errors.Is(err, ErrInvalidIssuer) {
				t.Fatalf("expected ErrInvalidIssuer, got %v", err)
			}
		})
	}
}

// A LEI is case-insensitive on the wire but stored upper case: the CHECK and
// the unique index both read the stored form, so normalising in one place is
// what makes "same LEI" mean the same thing everywhere.
func TestToIssuerWriteDTONormalises(t *testing.T) {
	req := validIssuerRequest()
	req.CountryCode = " lu "
	req.Name = "  Tooken RE  "
	req.LeiCode = strPtrT(" 2138001234567890abcd ")
	req.RegistrationNumber = strPtrT("  B123456  ")

	in, err := toIssuerWriteDTO(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if in.CountryCode != "LU" {
		t.Errorf("country_code = %q, want LU", in.CountryCode)
	}
	if in.Name != "Tooken RE" {
		t.Errorf("name = %q, want trimmed", in.Name)
	}
	if in.LeiCode == nil || *in.LeiCode != "2138001234567890ABCD" {
		t.Errorf("lei_code = %v, want upper case and trimmed", in.LeiCode)
	}
	if in.RegistrationNumber == nil || *in.RegistrationNumber != "B123456" {
		t.Errorf("registration_number = %v, want trimmed", in.RegistrationNumber)
	}
}

// An empty optional value is a NULL, never an empty string: lei_code is
// CHAR(20), so "" would be stored as twenty spaces and fail its own CHECK.
func TestToIssuerWriteDTOClearsEmptyOptionals(t *testing.T) {
	req := validIssuerRequest()
	req.LeiCode = strPtrT("   ")
	req.RegistrationNumber = strPtrT("")

	in, err := toIssuerWriteDTO(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if in.LeiCode != nil {
		t.Errorf("lei_code = %v, want nil", *in.LeiCode)
	}
	if in.RegistrationNumber != nil {
		t.Errorf("registration_number = %v, want nil", *in.RegistrationNumber)
	}
}

// A manager declaring a vehicle expects it usable, not parked in a draft he
// then has to find again; but asking for a draft must still be honoured.
func TestToIssuerWriteDTOStatusDefaultsToActive(t *testing.T) {
	in, err := toIssuerWriteDTO(validIssuerRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if in.StatusId != database.IssuerStatusActive {
		t.Errorf("status_id = %d, want %d", in.StatusId, database.IssuerStatusActive)
	}

	req := validIssuerRequest()
	req.StatusId = intPtrT(database.IssuerStatusDraft)

	in, err = toIssuerWriteDTO(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if in.StatusId != database.IssuerStatusDraft {
		t.Errorf("status_id = %d, want %d", in.StatusId, database.IssuerStatusDraft)
	}
}

// A patch changes what it names and nothing else — and an optional value sent
// empty is cleared, otherwise a LEI typed by mistake could never be removed.
func TestApplyIssuerPatch(t *testing.T) {
	current := server.IssuerWriteRequest{
		Name:               "Tooken RE",
		LegalForm:          "SCSp",
		RegistrationNumber: strPtrT("B123456"),
		CountryCode:        "LU",
		LeiCode:            strPtrT("2138001234567890ABCD"),
		StatusId:           intPtrT(database.IssuerStatusActive),
	}

	merged := applyIssuerPatch(current, server.IssuerPatchRequest{Name: strPtrT("Tooken RE II")})
	if merged.Name != "Tooken RE II" {
		t.Errorf("name = %q, want patched", merged.Name)
	}
	if merged.LegalForm != "SCSp" || merged.CountryCode != "LU" {
		t.Errorf("untouched fields changed: %+v", merged)
	}
	if merged.LeiCode == nil || *merged.LeiCode != "2138001234567890ABCD" {
		t.Errorf("lei_code lost by an unrelated patch: %v", merged.LeiCode)
	}

	cleared := applyIssuerPatch(current, server.IssuerPatchRequest{LeiCode: strPtrT("")})
	in, err := toIssuerWriteDTO(cleared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if in.LeiCode != nil {
		t.Errorf("lei_code = %v, want cleared", *in.LeiCode)
	}
}

// The merged result is validated as a whole: a patch cannot leave the issuer in
// a state a creation would have refused.
func TestPatchValidatedAsAWhole(t *testing.T) {
	current := server.IssuerWriteRequest{
		Name:        "Tooken RE",
		LegalForm:   "SCSp",
		CountryCode: "LU",
	}

	merged := applyIssuerPatch(current, server.IssuerPatchRequest{Name: strPtrT("  ")})
	if _, err := toIssuerWriteDTO(merged); !errors.Is(err, ErrInvalidIssuer) {
		t.Fatalf("expected ErrInvalidIssuer, got %v", err)
	}
}

// Withdrawing a vehicle from service while investors are being shown its assets
// would leave an offer standing on nothing — and activating one that is not
// identifiable would put an unverifiable vehicle behind an offer.
func TestCheckIssuerStatusChange(t *testing.T) {
	registered := func(status int) database.IssuerWriteDTO {
		return database.IssuerWriteDTO{
			StatusId:           status,
			RegistrationNumber: strPtrT("B123456"),
		}
	}
	unregistered := func(status int) database.IssuerWriteDTO {
		return database.IssuerWriteDTO{StatusId: status}
	}

	tests := []struct {
		name    string
		state   database.IssuerGuardStateDTO
		in      database.IssuerWriteDTO
		wantErr bool
	}{
		{
			name:  "no change is always allowed",
			state: database.IssuerGuardStateDTO{StatusId: database.IssuerStatusDissolved},
			in:    registered(database.IssuerStatusDissolved),
		},
		{
			// The freeze trap of §11.20: a legacy row stored without a
			// registration number must stay editable, or it could never be
			// corrected — not even its name.
			name:  "patching a legacy active row without a registration number",
			state: database.IssuerGuardStateDTO{StatusId: database.IssuerStatusActive},
			in:    unregistered(database.IssuerStatusActive),
		},
		{
			name:    "dissolved is final",
			state:   database.IssuerGuardStateDTO{StatusId: database.IssuerStatusDissolved},
			in:      registered(database.IssuerStatusActive),
			wantErr: true,
		},
		{
			name:    "activating a draft that is not identifiable",
			state:   database.IssuerGuardStateDTO{StatusId: database.IssuerStatusDraft},
			in:      unregistered(database.IssuerStatusActive),
			wantErr: true,
		},
		{
			name:  "activating a draft that supplies its number in the same request",
			state: database.IssuerGuardStateDTO{StatusId: database.IssuerStatusDraft},
			in:    registered(database.IssuerStatusActive),
		},
		{
			name: "suspending an issuer with published assets",
			state: database.IssuerGuardStateDTO{
				StatusId: database.IssuerStatusActive, AttachedCount: 2, PublishedCount: 1,
			},
			in:      registered(database.IssuerStatusSuspended),
			wantErr: true,
		},
		{
			name: "suspending an issuer whose assets are all drafts",
			state: database.IssuerGuardStateDTO{
				StatusId: database.IssuerStatusActive, AttachedCount: 2,
			},
			in: registered(database.IssuerStatusSuspended),
		},
		{
			name: "dissolving an issuer that still carries a draft",
			state: database.IssuerGuardStateDTO{
				StatusId: database.IssuerStatusActive, AttachedCount: 1,
			},
			in:      registered(database.IssuerStatusDissolved),
			wantErr: true,
		},
		{
			name:  "dissolving an issuer with nothing attached",
			state: database.IssuerGuardStateDTO{StatusId: database.IssuerStatusDraft},
			in:    registered(database.IssuerStatusDissolved),
		},
		{
			name: "reactivating a suspended issuer that carries published assets",
			state: database.IssuerGuardStateDTO{
				StatusId: database.IssuerStatusSuspended, AttachedCount: 3, PublishedCount: 3,
			},
			in: registered(database.IssuerStatusActive),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkIssuerStatusChange(tc.state, tc.in)
			if tc.wantErr && !errors.Is(err, ErrIssuerConflict) {
				t.Fatalf("expected ErrIssuerConflict, got %v", err)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// An active vehicle is one assets are published under: it has to be
// identifiable in its register, which is an investor's only way to check that
// it exists at all. A vehicle still being incorporated has no number yet — that
// is what the draft status is for.
func TestCheckIssuerCreatable(t *testing.T) {
	tests := []struct {
		name    string
		in      database.IssuerWriteDTO
		wantErr bool
	}{
		{
			name:    "active without a registration number",
			in:      database.IssuerWriteDTO{StatusId: database.IssuerStatusActive},
			wantErr: true,
		},
		{
			name: "active with one",
			in: database.IssuerWriteDTO{
				StatusId:           database.IssuerStatusActive,
				RegistrationNumber: strPtrT("B123456"),
			},
		},
		{
			// A vehicle being incorporated gets its number when the register
			// assigns it, days after the deed. Refusing the draft would push a
			// manager to invent one.
			name: "draft without a registration number",
			in:   database.IssuerWriteDTO{StatusId: database.IssuerStatusDraft},
		},
		{
			name:    "dissolved at birth",
			in:      database.IssuerWriteDTO{StatusId: database.IssuerStatusDissolved},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkIssuerCreatable(tc.in)
			if tc.wantErr && !errors.Is(err, ErrInvalidIssuer) {
				t.Fatalf("expected ErrInvalidIssuer, got %v", err)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// The LEI is never demanded: it is mandatory only for entities under market
// reporting obligations, and it lapses every year. Demanding it would turn a
// renewal delay into an unusable issuer.
func TestActiveRequirementsIgnoresTheLei(t *testing.T) {
	in := database.IssuerWriteDTO{
		StatusId:           database.IssuerStatusActive,
		RegistrationNumber: strPtrT("B123456"),
	}

	if missing := activeRequirements(in); len(missing) > 0 {
		t.Fatalf("nothing should be missing without a LEI, got %v", missing)
	}
}

// A duplicate registration number is a state conflict, not a malformed payload:
// the value is well formed, it simply already designates another vehicle. A
// referenced status that does not exist is the opposite — the caller's typo.
func TestTranslateIssuerWriteError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		sentine error
	}{
		{
			name:    "duplicate registration number",
			err:     &pq.Error{Code: "23505", Constraint: "issuer_registration_uk"},
			sentine: ErrIssuerConflict,
		},
		{
			name:    "duplicate lei",
			err:     &pq.Error{Code: "23505", Constraint: "issuer_lei_uk"},
			sentine: ErrIssuerConflict,
		},
		{
			name:    "unknown status",
			err:     &pq.Error{Code: "23503", Constraint: "issuer_status_id_fkey"},
			sentine: ErrInvalidIssuer,
		},
		{
			name:    "malformed country code",
			err:     &pq.Error{Code: "23514", Constraint: "issuer_country_code_ck"},
			sentine: ErrInvalidIssuer,
		},
		{
			name:    "malformed lei",
			err:     &pq.Error{Code: "23514", Constraint: "issuer_lei_ck"},
			sentine: ErrInvalidIssuer,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := translateIssuerWriteError(tc.err); !errors.Is(got, tc.sentine) {
				t.Fatalf("expected %v, got %v", tc.sentine, got)
			}
		})
	}
}

// A primary key collision is the sequence lagging behind the data, never
// something the caller sent: answering 409 sent a manager hunting for a
// duplicate value in a payload that never carried one (§11.17).
func TestPrimaryKeyCollisionIsNotAConflict(t *testing.T) {
	err := translateIssuerWriteError(&pq.Error{Code: "23505", Constraint: "issuer_pkey"})

	if errors.Is(err, ErrIssuerConflict) || errors.Is(err, ErrInvalidIssuer) {
		t.Fatalf("a sequence lag must stay a 500, got %v", err)
	}
	if !strings.Contains(err.Error(), "000013") {
		t.Errorf("the message must point at the fix, got %q", err.Error())
	}
}

// An issuer that carries published assets must not reach ANY non-active status.
// The guard used to name suspended and dissolved, which left active -> draft
// wide open: the same forbidden state, reached through the status nobody
// thought of as a withdrawal.
func TestCheckIssuerStatusChangeRefusesEveryWithdrawal(t *testing.T) {
	carrying := database.IssuerGuardStateDTO{
		Id:             1,
		StatusId:       database.IssuerStatusActive,
		AttachedCount:  4,
		PublishedCount: 4,
	}

	identified := database.IssuerWriteDTO{
		Name:               "SCI Dupont",
		LegalForm:          "SCI",
		CountryCode:        "LU",
		RegistrationNumber: strPtrT("B12345"),
	}

	for _, target := range []int{
		database.IssuerStatusDraft,
		database.IssuerStatusSuspended,
		database.IssuerStatusDissolved,
	} {
		t.Run(issuerStatusLabels[target], func(t *testing.T) {
			in := identified
			in.StatusId = target

			err := checkIssuerStatusChange(carrying, in)
			if !errors.Is(err, ErrIssuerConflict) {
				t.Fatalf("moving to %s must be refused, got %v", issuerStatusLabels[target], err)
			}
		})
	}
}

// The mirror of TestCheckStillPublishable on the issuer side: a patch that
// keeps the issuer active never reaches checkIssuerStatusChange, so the loss of
// an identification has to be caught on its own.
func TestCheckIssuerStillIdentifiable(t *testing.T) {
	active := func(registration *string) database.IssuerWriteDTO {
		return database.IssuerWriteDTO{
			Name:               "SCI Dupont",
			LegalForm:          "SCI",
			CountryCode:        "LU",
			RegistrationNumber: registration,
			StatusId:           database.IssuerStatusActive,
		}
	}

	identified := active(strPtrT("B12345"))
	anonymous := active(nil)

	t.Run("an active issuer cannot be stripped of its number", func(t *testing.T) {
		err := checkIssuerStillIdentifiable(identified, anonymous)
		if !errors.Is(err, ErrIssuerConflict) {
			t.Fatalf("expected a conflict, got %v", err)
		}
	})

	t.Run("a legacy row missing it stays editable", func(t *testing.T) {
		if err := checkIssuerStillIdentifiable(anonymous, anonymous); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("completing a legacy row is accepted", func(t *testing.T) {
		if err := checkIssuerStillIdentifiable(anonymous, identified); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("a suspended issuer is not held to it", func(t *testing.T) {
		suspended := anonymous
		suspended.StatusId = database.IssuerStatusSuspended

		if err := checkIssuerStillIdentifiable(identified, suspended); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// storedIssuerWriteDTO must not validate: a row stored before today's rules
// existed has to stay correctable, and running it through toIssuerWriteDTO
// would refuse the very patch that fixes it.
func TestStoredIssuerWriteDTODoesNotValidate(t *testing.T) {
	legacy := database.IssuerDTO{
		Id:          7,
		Name:        "",
		LegalForm:   "",
		CountryCode: "lu",
		StatusId:    database.IssuerStatusActive,
	}

	got := storedIssuerWriteDTO(legacy)
	if got.CountryCode != "lu" || got.Name != "" || got.StatusId != database.IssuerStatusActive {
		t.Fatalf("stored row must be mapped verbatim, got %+v", got)
	}
}
