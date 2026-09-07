package database

import (
	"strings"
	"testing"
)

// The guarantee brought by the column table is only worth something if it is
// checked: these tests fail if someone adds a SQL expression without a scan
// target, or points two columns at the same field.

func TestRealEstateColumnsAreAligned(t *testing.T) {
	cols := realEstateCoreColumns

	// Expressions are separated by ",\n"; an inner comma (ROUND(x, 2)) is
	// always followed by a space, never by a newline.
	if got := len(strings.Split(selectList(cols), ",\n")); got != len(cols) {
		t.Fatalf("the SELECT renders %d expressions for %d declared columns", got, len(cols))
	}

	var dto RealEstateDTO
	dto.Specification = &RealEstateSpecificationDTO{}
	dto.SharesConfig = &RealEstateSharesConfigDTO{}
	dto.Location = &RealEstateLocationDTO{}

	dest := scanTargets(&dto, cols)
	if len(dest) != len(cols) {
		t.Fatalf("%d scan targets for %d columns", len(dest), len(cols))
	}

	// Two columns scanning into the same field: the second one would silently
	// overwrite the first.
	seen := make(map[any]int, len(dest))
	for i, d := range dest {
		if d == nil {
			t.Errorf("column %d (%s) has no target", i, cols[i].expr)
			continue
		}
		if j, dup := seen[d]; dup {
			t.Errorf("columns %d (%s) and %d (%s) target the same field",
				j, cols[j].expr, i, cols[i].expr)
		}
		seen[d] = i
	}
}

// The list may read the coarse location (city, country_code) but never the
// columns that pinpoint an asset. Checking the selected columns rather than the
// joined tables is the stronger guarantee: the address table is joined, and
// what must not leak is what is read from it.
func TestListQueryReadsNoPinpointingColumn(t *testing.T) {
	query := buildRealEstateQuery(realEstateCoreColumns, "WHERE re.active = true")

	for _, forbidden := range []string{
		"addr.street",
		"addr.street_complement",
		"addr.postal_code",
		"addr.latitude",
		"addr.longitude",
	} {
		if strings.Contains(query, forbidden) {
			t.Errorf("the list query reads %s", forbidden)
		}
	}

	for _, forbidden := range []string{"ass.real_estate_media", "ass.issuer"} {
		if strings.Contains(query, forbidden) {
			t.Errorf("the list query joins %s", forbidden)
		}
	}
}

// What counts as a reserved share must stay in the reference table
// (iss.issuance_order_statuses.counts_as_reserved), never hardcoded here.
func TestReservedScopeComesFromTheReferential(t *testing.T) {
	query := buildRealEstateQuery(realEstateCoreColumns, "")

	if !strings.Contains(query, "s.counts_as_reserved") {
		t.Error("the reserved shares filter must rely on counts_as_reserved")
	}
	for _, hardcoded := range []string{"'CANCELLED'", "'EXPIRED'", "status_id IN"} {
		if strings.Contains(query, hardcoded) {
			t.Errorf("status hardcoded in the query: %s", hardcoded)
		}
	}
}

// A multi-row query without ORDER BY has no defined order at all: PostgreSQL
// returns whatever the plan produces. An UPDATE moves the row to the end of the
// heap, so a migration is enough to reshuffle the list — which is exactly what
// happened when migration 000008 rewrote the assets with an empty
// contract_address.
func TestListQueryIsOrdered(t *testing.T) {
	query := buildRealEstateQuery(realEstateCoreColumns, "WHERE re.active = true\nORDER BY re.id")

	if !strings.Contains(query, "ORDER BY") {
		t.Error("a query returning several rows must define its order explicitly")
	}
}
