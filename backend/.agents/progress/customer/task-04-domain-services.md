# Task 04 — Domain Services

> BC: customer | Files: `internal/customer/registrar_service.go`, `internal/customer/events.go`

---

## 4.1 — Domain Service Structs & Events

**Domain Events** — create `internal/customer/events.go` with file header:

Define all domain events (they all implement `shared.DomainEvent`):

```go
type CustomerRegisteredEvent struct {
    ShopId     ShopId
    CustomerId CustomerId
    UserId     shared.UserId
    // no email — per ADR-001, email is read from Auth BC; we don't duplicate it in events
    actorId    shared.UserId
    actorType  shared.ActorType
    occurredAt time.Time
}
// Implement EventName() string → "CustomerRegistered"
// Implement OccurredAt() time.Time
// Implement ActorId() shared.UserId
// Implement ActorType() shared.ActorType

type ProfileUpdatedEvent struct {
    CustomerId    CustomerId
    ChangedFields []string // list of field names that changed
    actorId       shared.UserId
    actorType     shared.ActorType
    occurredAt    time.Time
}
// EventName → "ProfileUpdated"

type AccountDeactivatedEvent struct {
    ShopId     ShopId
    CustomerId CustomerId
    UserId     shared.UserId
    actorId    shared.UserId
    actorType  shared.ActorType
    occurredAt time.Time
}
// EventName → "AccountDeactivated"

type WishlistItemAddedEvent struct {
    CustomerId CustomerId
    VariantId  string
    actorId    shared.UserId
    actorType  shared.ActorType
    occurredAt time.Time
}
// EventName → "WishlistItemAdded"
```

Add constructor functions for each event (e.g., `NewCustomerRegisteredEvent(...)`) to keep construction logic out of services.

**CustomerRegistrar** — create `internal/customer/registrar_service.go` with file header:

```
type CustomerRegistrar struct{}
func NewCustomerRegistrar() *CustomerRegistrar
```

Method `EnsureUnique(ctx context.Context, userId shared.UserId, shopId ShopId, repo CustomerRepo) error`:
- Calls `repo.FindByUserIdAndShop(ctx, userId, shopId)`
- If a profile is found, return `ErrCustomerAlreadyExists`
- If repo returns `ErrCustomerNotFound`, return nil (no existing profile — safe to register)
- Propagate any other repo error

Note: CustomerRegistrar is a pure domain service — no logger, no eventPub. Stateless.

**>>> STOP. Tell the user: "Task 04 Step 4.1 complete: domain service and event types. Ready for review." Wait for the user to say "proceed".**

---

## 4.2 — Domain Service Unit Tests

Create `internal/customer/registrar_service_test.go` with file header.

Create a minimal inline mock for `CustomerRepo` (only the `FindByUserIdAndShop` method needs to be mocked for this test file — implement all other interface methods as stubs returning nil/zero).

`TestCustomerRegistrar_EnsureUnique`:
- `"no existing customer/should return nil"` — mock returns `ErrCustomerNotFound` → `EnsureUnique` returns nil
- `"existing customer/should return ErrCustomerAlreadyExists"` — mock returns a profile → `EnsureUnique` returns `ErrCustomerAlreadyExists`
- `"repo error/should propagate"` — mock returns an unexpected error → `EnsureUnique` returns that error

Add to `internal/customer/customer_test.go`:

`TestCustomerRegisteredEvent`:
- Construct a `CustomerRegisteredEvent` via its constructor, verify `EventName()`, `OccurredAt()`, `ActorId()`, `ActorType()` return correct values.

**>>> STOP. Wait for user to say "proceed".**

---

## 4.3 — Run Tests

- Run: `go test ./internal/customer/...`
- If tests pass, continue.
- If tests fail: attempt a fix within spec scope and re-run, up to 2 attempts. If still failing, stop and report.

**>>> STOP. Tell the user: "Task 04 complete: domain services. Next item: CustomerSvc (Application Service)." Wait for user to say "proceed".**
