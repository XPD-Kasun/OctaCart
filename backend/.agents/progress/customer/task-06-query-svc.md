# Task 06 — CustomerQuerySvc Application Service

> BC: customer | File: `internal/customer/query_service.go`

---

## 6.1 — CustomerQuerySvc Struct & Methods

Create `internal/customer/query_service.go` with file header.

Define `CustomerQuerySvc` struct with dependencies:
```
customerRepo CustomerRepo
addressRepo  AddressRepo
```

Constructor: `NewCustomerQuerySvc(customerRepo CustomerRepo, addressRepo AddressRepo) *CustomerQuerySvc`

> Note: Query services hold no logger and no eventPub — they are read-only and return clean values. All query methods receive `shopId` from the application layer (extracted from JWT middleware), passed directly as a parameter.

Implement methods:

`GetCustomer(ctx context.Context, customerId CustomerId) (CustomerProfile, error)`:
- Delegate to `customerRepo.FindById(ctx, customerId)`
- Return profile by value

`GetCustomerByUserId(ctx context.Context, userId shared.UserId, shopId ShopId) (CustomerProfile, error)`:
- Delegate to `customerRepo.FindByUserIdAndShop(ctx, userId, shopId)`
- Return profile by value
- Used by driving adapters post-login to resolve customerId from the Auth UserId in the JWT

`ListCustomers(ctx context.Context, shopId ShopId, filter CustomerFilter, page shared.Pagination) ([]CustomerProfile, error)`:
- Delegate to `customerRepo.List(ctx, shopId, filter, page)` — shopId scoping is enforced at repo level
- Return slice of profiles by value (convert `[]*CustomerProfile` to `[]CustomerProfile`)

`ListAddresses(ctx context.Context, customerId CustomerId) ([]Address, error)`:
- Delegate to `addressRepo.FindByCustomerId(ctx, customerId)` — returns only non-soft-deleted
- Return slice of addresses by value (convert `[]*Address` to `[]Address`)

**>>> STOP. Tell the user: "Task 06 Step 6.1 complete: CustomerQuerySvc struct and methods. Ready for review." Wait for the user to say "proceed".**

---

## 6.2 — CustomerQuerySvc Unit Tests

Create `internal/customer/query_service_test.go` with file header.

Reuse mock implementations from `customer_service_test.go` (same package, white-box testing).

Helper `newTestCustomerQuerySvc(t)` → returns a `*CustomerQuerySvc` wired with `mockCustomerRepo` and `mockAddressRepo`.

**Test cases**:

`TestCustomerQuerySvc_GetCustomer`:
- `"existing customer/should return profile"` — seed a profile in mock, verify returned
- `"non-existing/should return ErrCustomerNotFound"`

`TestCustomerQuerySvc_GetCustomerByUserId`:
- `"existing/should resolve profile by userId and shopId"`
- `"wrong shopId/should return ErrCustomerNotFound"` — same userId, different shopId

`TestCustomerQuerySvc_ListCustomers`:
- `"no filter/should return all customers for shopId"` — seed 2 profiles with same shopId + 1 with different shopId (mock handles scoping); verify count
- `"status filter/should return only matching"` — seed active and deactivated; filter by Active status

`TestCustomerQuerySvc_ListAddresses`:
- `"customer with addresses/should return all active addresses"` — seed 2 active + 1 soft-deleted; verify 2 returned
- `"customer with no addresses/should return empty slice"`

**>>> STOP. Wait for user to say "proceed".**

---

## 6.3 — Run Tests

- Run: `go test ./internal/customer/...`
- If tests pass, continue.
- If tests fail: attempt a fix within spec scope and re-run, up to 2 attempts. If still failing, stop and report.

**>>> STOP. Tell the user: "Task 06 complete: CustomerQuerySvc. Next item: WishlistSvc." Wait for user to say "proceed".**
