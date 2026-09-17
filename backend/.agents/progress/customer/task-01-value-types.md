# Task 01 — Value Types

> BC: customer | File: `internal/customer/customer.go`

---

## 1.1 — Define Value Types

Create `internal/customer/customer.go` with the required file header and package doc comment explaining the Customer BC.

Define all value types:

**ID types** (opaque string aliases — per conventions):
- `type CustomerId string`
- `type AddressId string`
- `type WishlistItemId string`
- `type ShopId string` — received from JWT claim, not owned/validated by this BC

**CustomerStatus** (typed string enum):
- `type CustomerStatus string`
- Constants: `StatusActive CustomerStatus = "Active"`, `StatusDeactivated CustomerStatus = "Deactivated"`

**AddressSnapshot** (read-only value object for Order BC):
- `type AddressSnapshot struct` with **exported** fields (this is a DTO crossing a context boundary):
  - `RecipientName string`
  - `Line1 string`
  - `Line2 string` (empty string when absent)
  - `City string`
  - `Province string` (empty string when absent)
  - `PostalCode string`
  - `Country string`
  - `Phone string` (empty string when absent)
- Note: Does NOT include `AddressId` or `customerId` — it is intentionally stripped per spec.

**>>> STOP. Tell the user: "Task 01 Step 1.1 complete: value type structs defined. Ready for review." Wait for the user to say "proceed".**

---

## 1.2 — Value Type Unit Tests

Create `internal/customer/customer_test.go` with the file header.

Write `TestCustomerStatus_Constants`:
- Verify `StatusActive` has string value `"Active"`
- Verify `StatusDeactivated` has string value `"Deactivated"`

Write `TestAddressSnapshot_Fields`:
- Construct an `AddressSnapshot` with all fields populated, verify each field is accessible via direct struct access (exported fields)
- Construct one with only required fields (Line2, Province, Phone empty) — verify it compiles and no nil panic

**>>> STOP. Wait for user to say "proceed".**

---

## 1.3 — Run Tests

- Run: `go test ./internal/customer/...`
- If tests pass, continue.
- If tests fail: attempt a fix within spec scope and re-run, up to 2 attempts. If still failing, stop and report.

**>>> STOP. Tell the user: "Task 01 complete: value types. Next item: Entities." Wait for user to say "proceed".**
