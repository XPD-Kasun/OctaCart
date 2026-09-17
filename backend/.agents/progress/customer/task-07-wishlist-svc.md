# Task 07 — WishlistSvc Application Service

> BC: customer | File: `internal/customer/wishlist_service.go`

---

## 7.1 — WishlistSvc Struct & Methods

Create `internal/customer/wishlist_service.go` with file header.

Define `WishlistSvc` struct with dependencies:
```
wishlistRepo  WishlistRepo
customerRepo  CustomerRepo
eventPub      EventPublisher
logger        *zerolog.Logger
```

Constructor: `NewWishlistSvc(wishlistRepo WishlistRepo, customerRepo CustomerRepo, eventPub EventPublisher, logger *zerolog.Logger) *WishlistSvc`

Implement methods:

`AddToWishlist(ctx context.Context, customerId CustomerId, variantId string) (WishlistItem, error)`:
- Load customer profile → return `ErrCustomerNotFound` or `ErrCustomerDeactivated` if applicable
- Check for existing item: call `wishlistRepo.FindByCustomerAndVariant(ctx, customerId, variantId)`:
  - If found (not `ErrWishlistItemNotFound`), return the existing item by value (idempotent — per spec)
- Create via `NewWishlistItem(customerId, variantId)`
- Save via `wishlistRepo.Save`
- Publish `WishlistItemAddedEvent`
- Log: `logger.Info().Str("customerId", string(customerId)).Str("variantId", variantId).Msg("added to wishlist")`
- Return item by value

`RemoveFromWishlist(ctx context.Context, customerId CustomerId, wishlistItemId WishlistItemId) error`:
- Load item via `wishlistRepo.FindById` → return `ErrWishlistItemNotFound` if not found
- Verify `item.CustomerId() == customerId` (ownership check) — return `ErrWishlistItemNotFound` if mismatch
- Delete via `wishlistRepo.Delete`
- Log: `logger.Info().Str("wishlistItemId", string(wishlistItemId)).Msg("removed from wishlist")`

`ListWishlist(ctx context.Context, customerId CustomerId) ([]WishlistItem, error)`:
- Load customer → return `ErrCustomerNotFound` if not found
- Delegate to `wishlistRepo.FindByCustomerId(ctx, customerId)`
- Return slice by value (convert `[]*WishlistItem` to `[]WishlistItem`)

`ClearWishlist(ctx context.Context, customerId CustomerId) error`:
- Load customer → return `ErrCustomerNotFound` if not found
- Call `wishlistRepo.DeleteAllByCustomerId(ctx, customerId)`
- Log: `logger.Info().Str("customerId", string(customerId)).Msg("wishlist cleared")`

**>>> STOP. Tell the user: "Task 07 Step 7.1 complete: WishlistSvc struct and methods. Ready for review." Wait for the user to say "proceed".**

---

## 7.2 — WishlistSvc Unit Tests

Create `internal/customer/wishlist_service_test.go` with file header.

Create a minimal inline mock for `WishlistRepo` (all `WishlistRepo` methods backed by in-memory map).

Reuse `mockCustomerRepo` and `mockEventPublisher` from the same package.

Helper `newTestWishlistSvc(t)` → returns a `*WishlistSvc` wired with all mocks + `zerolog.Nop()` logger.

**Test cases**:

`TestWishlistSvc_AddToWishlist`:
- `"new item/should save and publish event"` — verify item saved, `WishlistItemAddedEvent` published
- `"duplicate variantId/should return existing item without duplicate save"` — call twice with same variantId, verify only one item in mock store, verify event published only once
- `"deactivated customer/should return ErrCustomerDeactivated"`
- `"non-existing customer/should return ErrCustomerNotFound"`

`TestWishlistSvc_RemoveFromWishlist`:
- `"existing item/should delete"` — seed item, remove, verify gone
- `"wrong customer/should return ErrWishlistItemNotFound"` — seed item for customerA, try remove as customerB
- `"non-existing item/should return ErrWishlistItemNotFound"`

`TestWishlistSvc_ListWishlist`:
- `"customer with items/should return all"` — seed 3 items
- `"non-existing customer/should return ErrCustomerNotFound"`
- `"customer with empty wishlist/should return empty slice"`

`TestWishlistSvc_ClearWishlist`:
- `"has items/should delete all"` — seed 3 items, clear, verify `FindByCustomerId` returns empty
- `"non-existing customer/should return ErrCustomerNotFound"`

**>>> STOP. Wait for user to say "proceed".**

---

## 7.3 — Run Tests (Final)

- Run: `go test ./internal/customer/...`
- If tests pass, continue.
- If tests fail: attempt a fix within spec scope and re-run, up to 2 attempts. If still failing, stop and report.

**>>> STOP. Tell the user: "Task 07 complete: WishlistSvc. All Customer BC domain code is implemented. Running final full test suite." Wait for user to say "proceed".**

---

## 7.4 — Final Verification

- Run: `go test ./internal/customer/...` — all tests must pass
- Run: `go build ./internal/customer/...` — must compile cleanly
- Run: `go vet ./internal/customer/...` — no vet errors

**>>> STOP. Tell the user: "Customer BC implementation complete. All tasks done. Awaiting architect review."**
