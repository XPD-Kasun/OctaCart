# Task 02 — Entities

> BC: customer | File: `internal/customer/customer.go` (append after value types)

---

## 2.1 — Entity Structs

Append all entity structs to `internal/customer/customer.go`.

**CustomerProfile** (Aggregate Root):
- Struct with **unexported** fields:
  - `id CustomerId`, `shopId ShopId`, `userId shared.UserId`
  - `firstName string`, `lastName string`, `displayName string`
  - `phone string` (empty string when absent), `avatarUri string` (empty string when absent)
  - `status CustomerStatus`
  - `isGuest bool` — true when this is a special guest-checkout profile (ADR-004)
  - `createdAt time.Time`, `updatedAt time.Time`, `createdBy shared.AdminId`
- Public getters: `Id()`, `ShopId()`, `UserId()`, `FirstName()`, `LastName()`, `DisplayName()`, `Phone()`, `AvatarUri()`, `Status()`, `IsGuest()`, `CreatedAt()`, `UpdatedAt()`, `CreatedBy()`
- Constructor: `NewCustomerProfile(shopId ShopId, userId shared.UserId, firstName, lastName string, createdBy shared.AdminId) (*CustomerProfile, error)`:
  - Sets `status` to `StatusActive`
  - Sets `isGuest` to `false`
  - Sets `createdAt` and `updatedAt` to `time.Now()`
  - Sets `displayName` to `firstName + " " + lastName` as default
  - Validates `firstName` and `lastName` are not empty (return error if so)
  - `id` left as zero value (DB assigns)
- Constructor for guest: `NewGuestCustomerProfile(shopId ShopId) *CustomerProfile`:
  - Sets `status` to `StatusActive`, `isGuest` to `true`
  - Sets `firstName` to `"Guest"`, `lastName` to `""`, `displayName` to `"Guest"`
  - `userId` left as zero value (guest has no Auth identity)
  - Sets timestamps to `time.Now()`
- Domain method: `Deactivate() error` — transitions Active → Deactivated. Returns `ErrInvalidStatus` if already Deactivated.
- Domain method: `Reactivate() error` — transitions Deactivated → Active. Returns `ErrInvalidStatus` if already Active.
- Domain method: `Update(cmd UpdateProfileCmd)` — updates mutable fields from command struct; refreshes `updatedAt`.
- Command struct `UpdateProfileCmd`:
  - `FirstName *string`, `LastName *string`, `DisplayName *string`
  - `Phone *string`, `AvatarUri *string`
  - Apply pointer-nil-means-skip semantics for each field.
  - When `FirstName` or `LastName` is updated, if `DisplayName` pointer is nil, auto-refresh `displayName` to `firstName + " " + lastName`.

**Address**:
- Struct with **unexported** fields:
  - `id AddressId`, `customerId CustomerId`
  - `label string` (empty string when absent, e.g., "Home", "Office")
  - `recipientName string`, `line1 string`, `line2 string` (empty string when absent)
  - `city string`, `province string` (empty string when absent), `postalCode string`, `country string`
  - `phone string` (empty string when absent)
  - `isDefault bool`
  - `isDeleted bool` — soft-delete flag (ADR-003)
- Public getters for all fields: `Id()`, `CustomerId()`, `Label()`, `RecipientName()`, `Line1()`, `Line2()`, `City()`, `Province()`, `PostalCode()`, `Country()`, `Phone()`, `IsDefault()`, `IsDeleted()`
- Constructor: `NewAddress(customerId CustomerId, input AddAddressInput) (*Address, error)`:
  - Validates `recipientName`, `line1`, `city`, `postalCode`, `country` are not empty
  - Sets `isDefault` and `isDeleted` to false
  - `id` left as zero value (DB assigns)
- Input struct `AddAddressInput` (≥ 4 fields, use struct per xcommand ADR):
  - `Label string`, `RecipientName string`, `Line1 string`, `Line2 string`
  - `City string`, `Province string`, `PostalCode string`, `Country string`, `Phone string`
- Domain method: `SetDefault(v bool)` — sets `isDefault`
- Domain method: `SoftDelete()` — sets `isDeleted` to `true`
- Domain method: `Update(cmd UpdateAddressCmd)`:
  - Applies pointer-nil-means-skip for all fields.
  - Command struct `UpdateAddressCmd`: same optional pointer fields as `AddAddressInput` fields but all as `*string`.
- Helper method: `ToSnapshot() AddressSnapshot` — returns an `AddressSnapshot` from this address's fields.

**WishlistItem**:
- Struct with **unexported** fields:
  - `id WishlistItemId`, `customerId CustomerId`, `variantId string` (cross-BC reference — plain string, no Product type imported), `addedAt time.Time`
- Public getters: `Id()`, `CustomerId()`, `VariantId()`, `AddedAt()`
- Constructor: `NewWishlistItem(customerId CustomerId, variantId string) (*WishlistItem, error)`:
  - Validates `variantId` is not empty
  - Sets `addedAt` to `time.Now()`
  - `id` left as zero value (DB assigns)

> **Cross-BC note**: `variantId` is stored as a plain `string` — the Customer BC does NOT import the Product BC package. This keeps the boundary clean per the context map.

**>>> STOP. Tell the user: "Task 02 Step 2.1 complete: entity structs. Ready for review." Wait for the user to say "proceed".**

---

## 2.2 — Entity Unit Tests

Add to `internal/customer/customer_test.go`.

**TestNewCustomerProfile**:
- `"valid input/should create active profile"` — verify status=Active, displayName set, timestamps not zero
- `"empty firstName/should return error"`
- `"empty lastName/should return error"`

**TestNewGuestCustomerProfile**:
- `"should create guest profile"` — verify `IsGuest()=true`, `Status()=Active`, `FirstName()="Guest"`

**TestCustomerProfile_Deactivate**:
- `"active profile/should deactivate"` — verify status=Deactivated
- `"already deactivated/should return ErrInvalidStatus"`

**TestCustomerProfile_Reactivate**:
- `"deactivated profile/should reactivate"` — verify status=Active
- `"already active/should return ErrInvalidStatus"`

**TestCustomerProfile_Update**:
- `"update firstName/should refresh displayName automatically"` — update firstName only (no displayName ptr), verify displayName changed
- `"update displayName explicitly/should override auto-refresh"` — set both firstName and displayName ptr, verify displayName matches explicit value
- `"nil fields/should not change anything"`

**TestNewAddress**:
- `"valid input/should create address"` — verify required fields, isDefault=false, isDeleted=false
- `"empty recipientName/should return error"`
- `"empty line1/should return error"`
- `"empty country/should return error"`

**TestAddress_SetDefault**:
- `"should set isDefault to true"`, `"should set isDefault to false"`

**TestAddress_SoftDelete**:
- `"should mark address as deleted"` — verify `IsDeleted()=true`

**TestAddress_ToSnapshot**:
- `"should map fields to AddressSnapshot without id or customerId"`

**TestNewWishlistItem**:
- `"valid/should create"` — verify variantId, addedAt not zero
- `"empty variantId/should return error"`

**>>> STOP. Wait for user to say "proceed".**

---

## 2.3 — Run Tests

- Run: `go test ./internal/customer/...`
- If tests pass, continue.
- If tests fail: attempt a fix within spec scope and re-run, up to 2 attempts. If still failing, stop and report.

**>>> STOP. Tell the user: "Task 02 complete: entities. Next item: Errors & Ports." Wait for user to say "proceed".**
