package database

// ReservedOrders is THE definition of the orders that hold shares of an asset.
// Every read of reserved shares builds on it: the stock
// checked when an order is placed, tokens_sold shown on the asset, and the
// guard state a manager's write is checked against. Three copies of this rule
// drifted once (the lapsed reservations were counted by two of them): keep one.
//
// It is a FROM ... WHERE fragment over two aliases, o for iss.issuance_orders
// and s for its status. Callers select what they need and may add conditions
// with AND, or a GROUP BY.
//
// An order holds shares when:
//   - its status says so (iss.issuance_order_statuses.counts_as_reserved), so
//     adding a status to the flow needs no change here;
//   - and, if it is still AWAITING_PAYMENT, its reservation has not lapsed.
//     Nothing writes EXPIRED (D6): the lapse is read from the date, with the
//     database clock, the one reservation_expires_at was written with (U2).
const ReservedOrders = `
FROM iss.issuance_orders o
JOIN iss.issuance_order_statuses s
    ON s.id = o.status_id
WHERE s.counts_as_reserved
  AND (s.code <> 'AWAITING_PAYMENT' OR o.reservation_expires_at > now())`

// reservedByAsset is ReservedOrders summed per asset, to be LEFT JOINed as
// "sold" on the asset (sold.asset_id, sold.reserved).
//
// The aggregate is computed ONCE and then joined: a correlated subquery would
// be re-executed per asset (N+1), measured at 196 ms against 9 ms on 500
// assets / 50,000 orders. On a detail query PostgreSQL pushes the re.id = $1
// predicate into the aggregate.
const reservedByAsset = `
LEFT JOIN (
    SELECT o.asset_id, SUM(o.quantity) AS reserved` + ReservedOrders + `
    GROUP BY o.asset_id
) sold
    ON sold.asset_id = re.id
`
