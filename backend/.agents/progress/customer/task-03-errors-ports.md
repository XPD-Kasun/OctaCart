# Task 03 — Errors & Ports

> BC: customer | Files: `internal/customer/errors.go`, `internal/customer/ports.go`

---

## 3.1 — Error Sentinels & Port Interfaces

**`internal/customer/errors.go`**:
- Create with file header
- Define all sentinel errors:

```go
var (
    ErrCustomerNotFound      = errors.New("customer not found")
    ErrCustomerAlreadyExists = errors.New("customer already exists for this userId and shopId")
    ErrCustomerDeactivated   = errors.New("customer account is deactivated")
    ErrInvalidStatus         = errors.New("invalid status transition")
    ErrAddressNotFound       = errors.New("address not found")
    ErrCannotDeleteOnlyAddress = errors.New("cannot delete address linked to active orders")
    ErrWishlistItemNotFound  = errors.New("wishlist item not found")
)
```

**`internal/customer/ports.go`**:
- Create with file header
- Define `CustomerFilter` struct (used by repo list query):
  - `Status *CustomerStatus`, `SearchQuery *string` (matches on name/displayName)
- Define all driven port interfaces:

**CustomerRepo** interface:
- `Save(ctx context.Context, profile *CustomerProfile) error`
- `FindById(ctx context.Context, id CustomerId) (*CustomerProfile, error)`
- `FindByUserIdAndShop(ctx context.Context, userId shared.UserId, shopId ShopId) (*CustomerProfile, error)` — per shopId-scoping note in spec
- `List(ctx context.Context, shopId ShopId, filter CustomerFilter, page shared.Pagination) ([]*CustomerProfile, error)`
- `Delete(ctx context.Context, id CustomerId) error`
- `GetEmailByUserId(ctx context.Context, userId shared.UserId) (string, error)` — resolves email from Auth BC (ADR-001); implemented in the driven adapter by querying Auth's user table

**AddressRepo** interface:
- `Save(ctx context.Context, address *Address) error`
- `FindById(ctx context.Context, id AddressId) (*Address, error)`
- `FindByCustomerId(ctx context.Context, customerId CustomerId) ([]*Address, error)` — excludes soft-deleted
- `FindDefaultByCustomerId(ctx context.Context, customerId CustomerId) (*Address, error)`
- `ClearDefaultForCustomer(ctx context.Context, customerId CustomerId) error` — clears isDefault on all addresses before setting a new default
- `Delete(ctx context.Context, id AddressId) error` — hard delete (adapter handles soft-delete vs hard-delete per ADR-003)

**OrderQueryPort** interface — allows Customer BC to check if an address is linked to active orders (ADR-003):
- `HasActiveOrdersForAddress(ctx context.Context, addressId AddressId) (bool, error)`
- `HasHistoricalOrdersForAddress(ctx context.Context, addressId AddressId) (bool, error)`

**WishlistRepo** interface:
- `Save(ctx context.Context, item *WishlistItem) error`
- `FindById(ctx context.Context, id WishlistItemId) (*WishlistItem, error)`
- `FindByCustomerAndVariant(ctx context.Context, customerId CustomerId, variantId string) (*WishlistItem, error)` — for idempotency check
- `FindByCustomerId(ctx context.Context, customerId CustomerId) ([]*WishlistItem, error)`
- `Delete(ctx context.Context, id WishlistItemId) error`
- `DeleteAllByCustomerId(ctx context.Context, customerId CustomerId) error`

**EventPublisher** interface:
- `Publish(ctx context.Context, event shared.DomainEvent) error`
- Note: Use the `shared.DomainEvent` interface already defined in `internal/shared/types.go`. Do NOT redefine it.

**>>> STOP. Tell the user: "Task 03 Step 3.1 complete: errors and ports. Ready for review." Wait for the user to say "proceed".**

---

## 3.2 — Compile Check

- No behavioral tests for interfaces/sentinels.
- Run: `go build ./internal/customer/...`
- Confirm no compilation errors.
- Run: `go test ./internal/customer/...` to confirm existing tests still pass.

**>>> STOP. Wait for user to say "proceed".**

---

## 3.3 — Verify

- Run: `go test ./internal/customer/...`

**>>> STOP. Tell the user: "Task 03 complete: errors and ports. Next item: Domain Services." Wait for user to say "proceed".**
