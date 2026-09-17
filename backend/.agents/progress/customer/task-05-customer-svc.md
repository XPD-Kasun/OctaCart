# Task 05 — CustomerSvc Application Service

> BC: customer | File: `internal/customer/customer_service.go`

---

## 5.1 — CustomerSvc Struct & Methods

Create `internal/customer/customer_service.go` with file header.

Define `CustomerSvc` struct with dependencies:
```
customerRepo CustomerRepo
addressRepo  AddressRepo
orderQuery   OrderQueryPort
eventPub     EventPublisher
registrar    *CustomerRegistrar
logger       *zerolog.Logger
```

Constructor: `NewCustomerSvc(customerRepo CustomerRepo, addressRepo AddressRepo, orderQuery OrderQueryPort, eventPub EventPublisher, registrar *CustomerRegistrar, logger *zerolog.Logger) *CustomerSvc`

Implement the following methods. Extract `actorId` and `actorType` from `ctx` via `shared.Claim` where needed for events.

---

**Profile Methods**:

`RegisterCustomer(ctx context.Context, shopId ShopId, userId shared.UserId, firstName, lastName string) (CustomerProfile, error)`:
- Call `registrar.EnsureUnique(ctx, userId, shopId, customerRepo)` — return `ErrCustomerAlreadyExists` if duplicate
- Extract `createdBy` (AdminId) from ctx claims (or use zero value if triggered by the user themselves)
- Call `NewCustomerProfile(shopId, userId, firstName, lastName, createdBy)`
- Save via `customerRepo.Save`
- Publish `CustomerRegisteredEvent`
- Log: `logger.Info().Str("customerId", string(profile.Id())).Msg("customer registered")`
- Return profile by value

`UpdateProfile(ctx context.Context, customerId CustomerId, cmd UpdateProfileCmd) (CustomerProfile, error)`:
- Load profile via `customerRepo.FindById` → return `ErrCustomerNotFound` if not found
- If profile is deactivated, return `ErrCustomerDeactivated`
- Track which fields change (for event)
- Call `profile.Update(cmd)`
- Save via `customerRepo.Save`
- If any field changed, publish `ProfileUpdatedEvent` with changed field names
- Log: `logger.Info().Str("customerId", string(customerId)).Msg("profile updated")`
- Return profile by value

`DeactivateAccount(ctx context.Context, customerId CustomerId) error`:
- Load profile
- Call `profile.Deactivate()` → return error if `ErrInvalidStatus`
- Save via `customerRepo.Save`
- Publish `AccountDeactivatedEvent`
- Log: `logger.Info().Str("customerId", string(customerId)).Msg("account deactivated")`

`ReactivateAccount(ctx context.Context, customerId CustomerId) error`:
- Load profile
- Call `profile.Reactivate()` → return error if `ErrInvalidStatus`
- Save via `customerRepo.Save`
- Log: `logger.Info().Str("customerId", string(customerId)).Msg("account reactivated")`

---

**Address Methods**:

`AddAddress(ctx context.Context, customerId CustomerId, input AddAddressInput) (Address, error)`:
- Load profile → return `ErrCustomerNotFound` / `ErrCustomerDeactivated` if applicable
- Call `NewAddress(customerId, input)` → return error if validation fails
- Load existing addresses via `addressRepo.FindByCustomerId`
- If no existing addresses: set `address.SetDefault(true)` — first address auto-becomes default
- Save via `addressRepo.Save`
- Log: `logger.Info().Str("customerId", string(customerId)).Msg("address added")`
- Return address by value

`UpdateAddress(ctx context.Context, customerId CustomerId, addressId AddressId, cmd UpdateAddressCmd) (Address, error)`:
- Load address via `addressRepo.FindById` → return `ErrAddressNotFound` if not found or soft-deleted
- Verify `address.CustomerId() == customerId` (ownership check) — return `ErrAddressNotFound` if mismatch
- Call `address.Update(cmd)`
- Save via `addressRepo.Save`
- Return address by value

`DeleteAddress(ctx context.Context, customerId CustomerId, addressId AddressId) error`:
- Load address → verify ownership
- Check `orderQuery.HasActiveOrdersForAddress(ctx, addressId)` — if true, return `ErrCannotDeleteOnlyAddress`
- Check `orderQuery.HasHistoricalOrdersForAddress(ctx, addressId)` — if true, call `address.SoftDelete()` then save (soft-delete per ADR-003)
- Otherwise: call `addressRepo.Delete(ctx, addressId)` (hard delete)
- If the deleted address was the default, attempt to promote another address:
  - Load remaining active addresses; if any, call `SetDefaultAddress` on the first one
- Log: `logger.Info().Str("addressId", string(addressId)).Msg("address deleted")`

`SetDefaultAddress(ctx context.Context, customerId CustomerId, addressId AddressId) error`:
- Load address → verify ownership and not soft-deleted
- Call `addressRepo.ClearDefaultForCustomer(ctx, customerId)` — clears all existing defaults
- Call `address.SetDefault(true)`
- Save via `addressRepo.Save`
- Log: `logger.Info().Str("addressId", string(addressId)).Msg("default address set")`

`GetAddressSnapshot(ctx context.Context, customerId CustomerId, addressId *AddressId) (AddressSnapshot, error)`:
- If `addressId` is nil: load via `addressRepo.FindDefaultByCustomerId` → return `ErrAddressNotFound` if none
- Else: load via `addressRepo.FindById` → verify ownership
- Call `address.ToSnapshot()` and return

**>>> STOP. Tell the user: "Task 05 Step 5.1 complete: CustomerSvc struct and methods. Ready for review." Wait for the user to say "proceed".**

---

## 5.2 — CustomerSvc Unit Tests

Create `internal/customer/customer_service_test.go` with file header.

Create mock implementations (inline in the test file):
- `mockCustomerRepo` — in-memory `map[CustomerId]*CustomerProfile`; implement all `CustomerRepo` methods
- `mockAddressRepo` — in-memory `map[AddressId]*Address`; implement all `AddressRepo` methods
  - `FindByCustomerId` returns only non-soft-deleted addresses
  - `FindDefaultByCustomerId` returns the first address with `isDefault=true`
  - `ClearDefaultForCustomer` iterates map and sets `isDefault=false` for all matching customerId
- `mockOrderQueryPort` — configurable booleans for `HasActiveOrdersForAddress` and `HasHistoricalOrdersForAddress`
- `mockEventPublisher` — slice to record published events
- Helper `newTestCustomerSvc(t)` → returns a `*CustomerSvc` wired with all mocks + a `zerolog.Nop()` logger

**Test cases**:

`TestCustomerSvc_RegisterCustomer`:
- `"new customer/should register and publish event"` — verify profile saved, `CustomerRegisteredEvent` published
- `"duplicate userId+shopId/should return ErrCustomerAlreadyExists"` — mock returns existing profile from `FindByUserIdAndShop`
- `"empty firstName/should return error"`

`TestCustomerSvc_UpdateProfile`:
- `"valid update/should update and publish ProfileUpdatedEvent"` — update firstName, verify event fired with changed fields
- `"deactivated customer/should return ErrCustomerDeactivated"`
- `"no changes/should not publish event"` — all-nil cmd, verify no event

`TestCustomerSvc_DeactivateAccount`:
- `"active account/should deactivate and publish AccountDeactivatedEvent"`
- `"already deactivated/should return ErrInvalidStatus"`

`TestCustomerSvc_ReactivateAccount`:
- `"deactivated account/should reactivate"`
- `"already active/should return ErrInvalidStatus"`

`TestCustomerSvc_AddAddress`:
- `"first address/should become default automatically"`
- `"second address/should not override default"`
- `"deactivated customer/should return ErrCustomerDeactivated"`
- `"invalid address (empty city)/should return error"`

`TestCustomerSvc_SetDefaultAddress`:
- `"valid address/should clear old default and set new"` — add 2 addresses, set second as default, verify first's isDefault=false

`TestCustomerSvc_DeleteAddress`:
- `"no orders linked/should hard delete"` — mock orderQuery returns false for both → verify `Delete` called
- `"active orders linked/should return ErrCannotDeleteOnlyAddress"`
- `"historical orders only/should soft delete"` — mock historical=true, active=false → verify `SoftDelete` and `Save` called, not `Delete`
- `"deleting default address/should promote another as default"` — add 2 addresses, first is default, delete first → second becomes default

`TestCustomerSvc_GetAddressSnapshot`:
- `"addressId provided/should return snapshot for that address"`
- `"nil addressId/should return snapshot for default address"`
- `"wrong customer ownership/should return ErrAddressNotFound"`

**>>> STOP. Wait for user to say "proceed".**

---

## 5.3 — Run Tests

- Run: `go test ./internal/customer/...`
- If tests pass, continue.
- If tests fail: attempt a fix within spec scope and re-run, up to 2 attempts. If still failing, stop and report.

**>>> STOP. Tell the user: "Task 05 complete: CustomerSvc. Next item: CustomerQuerySvc." Wait for user to say "proceed".**
