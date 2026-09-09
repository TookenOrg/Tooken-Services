package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/globals"
)

// Read and write side of ass.issuer, the legal vehicle a token represents
// shares of.
//
// The table was created by migration 000007 but stayed reachable only as a
// sub-object of an asset. Since publication demands an issuer_id (§11.18), the
// only way to fill it was a direct SQL INSERT: this file closes that gap.

// The four statuses seeded by migration 000007. They are hardcoded here for the
// same reason the asset statuses are: the flow is decided by the code, and a
// lookup on every write would buy nothing.
const (
	IssuerStatusDraft     = 1
	IssuerStatusActive    = 2
	IssuerStatusSuspended = 3
	IssuerStatusDissolved = 4
)

// IssuerDTO maps a row of ass.issuer, joined with its status label and with how
// many assets depend on it.
//
// RealEstateCount is part of the row rather than a separate call: it is what
// decides whether the issuer can be withdrawn from service, and a manager
// reading the list needs it on every line.
type IssuerDTO struct {
	Id                 int
	Name               string
	LegalForm          string
	RegistrationNumber *string
	CountryCode        string
	LeiCode            *string
	StatusId           int
	StatusCode         string
	RealEstateCount    int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// IssuerWriteDTO is an issuer as accepted by a write.
//
// The three optional columns are pointers: nil is a NULL, never an empty
// string. `lei_code CHAR(20)` would pad "" to twenty spaces and then fail its
// own CHECK, so the distinction is not cosmetic.
type IssuerWriteDTO struct {
	Name               string
	LegalForm          string
	RegistrationNumber *string
	CountryCode        string
	LeiCode            *string
	StatusId           int
}

// IssuerGuardStateDTO is what a business rule needs before changing an issuer.
//
// The two counts are distinct on purpose: assets merely attached to an issuer
// forbid dissolving it, while assets visible to investors additionally forbid
// suspending it — an offer must not outlive the vehicle standing behind it.
type IssuerGuardStateDTO struct {
	Id             int
	StatusId       int
	AttachedCount  int
	PublishedCount int
}

// issuerSelect is shared by the list and the detail so the two can never
// disagree on what an issuer is.
const issuerSelect = `
SELECT i.id, i.name, i.legal_form, i.registration_number,
       i.country_code, i.lei_code, i.status_id, ist.code,
       COALESCE(re.attached, 0)::int AS real_estate_count,
       i.created_at, i.updated_at
FROM ass.issuer i
JOIN ass.issuer_status ist
    ON ist.id = i.status_id
LEFT JOIN (
    SELECT issuer_id, COUNT(*) AS attached
    FROM ass.real_estate
    WHERE issuer_id IS NOT NULL AND deleted_at IS NULL
    GROUP BY issuer_id
) re
    ON re.issuer_id = i.id
`

func scanIssuer(scan func(dest ...any) error) (IssuerDTO, error) {
	var i IssuerDTO
	err := scan(
		&i.Id, &i.Name, &i.LegalForm, &i.RegistrationNumber,
		&i.CountryCode, &i.LeiCode, &i.StatusId, &i.StatusCode,
		&i.RealEstateCount, &i.CreatedAt, &i.UpdatedAt,
	)

	return i, err
}

// GetIssuers lists every issuer, dissolved ones included.
//
// Hiding a dissolved vehicle would leave a manager unable to understand why an
// identifier he remembers is refused. The ordering is by name so the list is
// stable between two calls, which an unordered SELECT does not guarantee.
func GetIssuers(ctx context.Context) ([]IssuerDTO, error) {
	rows, err := globals.DB.QueryContext(ctx, issuerSelect+"ORDER BY i.name, i.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	issuers := []IssuerDTO{}
	for rows.Next() {
		issuer, err := scanIssuer(rows.Scan)
		if err != nil {
			return nil, err
		}
		issuers = append(issuers, issuer)
	}

	return issuers, rows.Err()
}

// GetIssuerById reports sql.ErrNoRows when there is no such issuer, so the
// caller can answer 404 without a second round trip.
func GetIssuerById(ctx context.Context, id int) (IssuerDTO, error) {
	row := globals.DB.QueryRowContext(ctx, issuerSelect+"WHERE i.id = $1", id)

	return scanIssuer(row.Scan)
}

// GetIssuerStatus reads just the status of an issuer.
//
// A dedicated query rather than GetIssuerById: the publication path needs one
// integer, and the full row carries an aggregate over ass.real_estate that
// every publication would then pay for without reading it.
//
// It reports sql.ErrNoRows when there is no such issuer.
func GetIssuerStatus(ctx context.Context, id int) (int, error) {
	var status int
	err := globals.DB.QueryRowContext(ctx,
		"SELECT status_id FROM ass.issuer WHERE id = $1", id).Scan(&status)

	return status, err
}

// GetIssuerGuardState reads the state a write must be checked against.
//
// "Published" is counted with re.active on purpose, not with
// real_estate_status.is_public. The two disagree: the seed of migration 000006
// marks status 6 (closed) is_public, while the trigger of the same migration
// sets active := status_id IN (3,4,5). What decides whether an investor
// actually sees the asset is re.active, because that is what the listing and
// the detail queries filter on (real_estate_db.go). Counting is_public here
// would block an issuer over assets nobody can see.
//
// That divergence is inherited debt, not a decision taken here — see
// PROGRESS.md §11.23.
func GetIssuerGuardState(ctx context.Context, id int) (IssuerGuardStateDTO, error) {
	const query = `
SELECT i.id,
       i.status_id,
       COUNT(re.id) FILTER (WHERE re.deleted_at IS NULL)::int AS attached,
       COUNT(re.id) FILTER (WHERE re.deleted_at IS NULL AND re.active)::int AS published
FROM ass.issuer i
LEFT JOIN ass.real_estate re
    ON re.issuer_id = i.id
WHERE i.id = $1
GROUP BY i.id, i.status_id
`

	var s IssuerGuardStateDTO
	err := globals.DB.QueryRowContext(ctx, query, id).Scan(
		&s.Id, &s.StatusId, &s.AttachedCount, &s.PublishedCount,
	)
	if err != nil {
		return IssuerGuardStateDTO{}, err
	}

	return s, nil
}

// CreateIssuer inserts the vehicle and returns its identifier.
func CreateIssuer(ctx context.Context, in IssuerWriteDTO) (id int, err error) {
	const query = `
INSERT INTO ass.issuer (name, legal_form, registration_number, country_code, lei_code, status_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
`

	if err := globals.DB.QueryRowContext(ctx, query,
		in.Name, in.LegalForm, in.RegistrationNumber, in.CountryCode, in.LeiCode, in.StatusId,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert issuer: %w", err)
	}

	return id, nil
}

// UpdateIssuer replaces the issuer with the merged state computed by the
// caller, exactly as an asset update does: one definition of what a stored
// issuer looks like, whether it was created or patched.
//
// The withdrawal guard travels INSIDE the statement rather than being checked
// before it. Read on one connection and written on another, it left a window in
// which an asset could be published between the two, and the foreign key says
// nothing about a status change — it is declared ON DELETE RESTRICT, not ON
// UPDATE. Postgres evaluates the NOT EXISTS and the write as one statement, so
// the window does not exist. A mutex would not do: the two requests may well be
// served by two instances of the API.
//
// The condition is skipped when the target status is active, because becoming
// active is never a withdrawal.
//
// It reports sql.ErrNoRows when nothing was written — the issuer is gone, or the
// guard declined. The caller re-reads the state to tell the two apart.
func UpdateIssuer(ctx context.Context, id int, in IssuerWriteDTO) error {
	const query = `
UPDATE ass.issuer
SET name = $2,
    legal_form = $3,
    registration_number = $4,
    country_code = $5,
    lei_code = $6,
    status_id = $7,
    updated_at = now()
WHERE id = $1
  AND ($7::smallint = $8::smallint
       OR NOT EXISTS (
           SELECT 1
           FROM ass.real_estate
           WHERE issuer_id = $1
             AND deleted_at IS NULL
             AND active))
`

	res, err := globals.DB.ExecContext(ctx, query,
		id, in.Name, in.LegalForm, in.RegistrationNumber, in.CountryCode, in.LeiCode, in.StatusId,
		IssuerStatusActive,
	)
	if err != nil {
		return fmt.Errorf("update issuer: %w", err)
	}

	return expectOneRow(res)
}

// DissolveIssuer is the deletion of an issuer: the row stays, the status
// becomes final.
//
// Assets, orders and on-chain tokens point at it through a FK declared
// ON DELETE RESTRICT — a real DELETE could not even run once a single asset
// referenced it, and would silently erase history if none did.
//
// As in UpdateIssuer, the "carries no asset" condition is part of the statement:
// checking it first would leave a window in which an asset gets attached, and
// the result would be an offer standing on a liquidated vehicle, unreachable by
// the API afterwards.
//
// It reports sql.ErrNoRows when nothing was written. That covers three cases —
// the row is gone, it was already dissolved, or an asset was attached in the
// meantime — which only the caller can tell apart, and must: the second is a
// success, not a 404.
func DissolveIssuer(ctx context.Context, id int) error {
	const query = `
UPDATE ass.issuer
SET status_id = $2,
    updated_at = now()
WHERE id = $1
  AND status_id <> $2
  AND NOT EXISTS (
      SELECT 1
      FROM ass.real_estate
      WHERE issuer_id = $1
        AND deleted_at IS NULL)
`

	res, err := globals.DB.ExecContext(ctx, query, id, IssuerStatusDissolved)
	if err != nil {
		return err
	}

	return expectOneRow(res)
}

func expectOneRow(res sql.Result) error {
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
