# Product Bounded Context — Implementation Plan

> Generated: 2026-09-07 | BC: product | Spec: `.agents/specs/product/domain-model.md`

---

## Preamble: Coding Instructions for the Implementer

Before you write any code, read the following instructions in full. These are your operating rules for the entire plan.

### Step 0 — Read Context

- [ ] Read `AGENTS.md` at the project root to understand the architecture, tech stack, and folder structure.
- [ ] Read `.agents/specs/product/domain-model.md` — this is your **authoritative spec**. Every type, method, and error you create must match it.
- [ ] Read **all ADR files** `.agents/specs/product/ADR001.md` through `ADR009.md` — these are reviewed architectural decisions that inform implementation details. **Do NOT read `adr-archived.md**` — it is obsolete.
- [ ] Read `.agents/discover/conventions.md` — this defines all coding style rules.
- [ ] Read `.agents/discover/context-map.md` — this tells you the cross-BC boundaries. The product BC must **never** write to data owned by another BC.

### Coding Conventions Summary

1. **Module**: `octacart`
2. **File header**: Every `.go` file starts with the copyright block (see conventions.md).
3. **Package**: All product domain code goes in `internal/product/` under `package product`.
4. **Entity fields**: Always **unexported** (lowercase). Provide **public getters** (PascalCase, no `Get` prefix).
5. **Constructors**: `New{Type}(...)` factory functions. Return `({Type}, error)` when validation is needed.
6. **xcommand pattern (ADR-005)**: Use input structs for methods with >= 4 params. Use bare param lists for < 4 params.
7. **Errors**: Sentinel vars in `errors.go`: `var ErrXxx = errors.New("...")`.
8. **Ports**: Interfaces in `ports.go`. First param is always `context.Context`. Last return is always `error`.
9. **Tests**: Standard `testing` package only. Same package (white-box). Use `t.Run` subtests. Name pattern: `Test{Type}_{Method}`.
10. **Imports**: stdlib → project → third-party, separated by blank lines.
11. **Shared types**: `shared.Money` (int64), `shared.KV`, `shared.UserId`, `shared.Claim` already exist. You will add `shared.Pagination` and `shared.DomainEvent`.
12. **Money (ADR-004)**: All monetary values are stored as `int64` in the **configured least minor units** (e.g., cents/paise). Single-currency per tenant. **No floating-point**. Currency is configured in Settings BC.
13. **Media URIs (ADR-002)**: Product media URIs use **relative paths** for a configured storage directory or **direct URLs without `file:///` prefix**. The `MediaRepo`/`ImageStore` adapters resolve locations relative to environment config.
14. **Logging (zerolog)**: All app services receive `*zerolog.Logger` via constructor. Use structured logging at key points: `Info` for successful mutations (create, publish, archive, stock adjust), `Warn` for business rule blocks (e.g., no variants), `Error` for unexpected failures. Always include entity IDs: `log.Info().Int("productId", int(id)).Msg("product published")`. Keep domain services log-free — logging is an app service concern.
15. **OTel-readiness**: All service methods already take `context.Context` first. When OTel is added later, tracing spans attach to ctx at the Gin middleware layer. **No domain code changes needed** — just add middleware + optional `otel.Tracer.Start(ctx, ...)` in app service entry points.

### Test Running Protocol

- After writing tests and implementation for each item, run: `go test ./internal/product/...`
- If tests pass, continue to the next item.
- If tests fail: attempt a fix **within the spec's scope** and re-run, up to **2 attempts**.
- If still failing after 2 attempts, or the fix would require **deviating from the spec**: **STOP** and report the failure to the user. Do not proceed further.

### File Layout

You will create these files in `internal/product/`:

| File                       | Contents                                                                                                                                         |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| `product.go`               | Product entity, ProductVariant entity, Category entity, ProductMedia entity, all value types, constructors, domain methods. Package doc comment. |
| `ports.go`                 | All driven port interfaces (ProductRepo, VariantRepo, CategoryRepo, MediaRepo, ImageStore, OrderQueryPort, EventPublisher)                       |
| `errors.go`                | All sentinel error vars                                                                                                                          |
| `slug_service.go`          | SlugGenerator domain service                                                                                                                     |
| `stock_service.go`         | StockAdjuster domain service                                                                                                                     |
| `product_service.go`       | ProductSvc application service                                                                                                                   |
| `media_service.go`         | ProductMediaSvc application service                                                                                                              |
| `category_service.go`      | CategorySvc application service                                                                                                                  |
| `query_service.go`         | ProductQuerySvc application service                                                                                                              |
| `product_test.go`          | Tests for entities and value types                                                                                                               |
| `slug_service_test.go`     | Tests for SlugGenerator                                                                                                                          |
| `stock_service_test.go`    | Tests for StockAdjuster                                                                                                                          |
| `product_service_test.go`  | Tests for ProductSvc                                                                                                                             |
| `media_service_test.go`    | Tests for ProductMediaSvc                                                                                                                        |
| `category_service_test.go` | Tests for CategorySvc                                                                                                                            |
| `query_service_test.go`    | Tests for ProductQuerySvc                                                                                                                        |

And in `internal/shared/`:

| File       | Contents                                                |
| ---------- | ------------------------------------------------------- |
| `types.go` | Add `Pagination` and `DomainEvent` to the existing file |

### Important — Existing Code

The file `internal/product/product.go` already exists with a **placeholder** `Product` struct. You must **replace it entirely** with the spec-compliant version. The placeholder has incorrect field types (`id int`, `name string`, `sku string`) that do not match the spec.

---

## Phase 1: Shared Types

### 1.1 — Shared Type Structs

- [x] Open `internal/shared/types.go`
- [x] Add the `Pagination` struct with these fields:
  - `Page int` — 1-based page number
  - `PerPage int` — items per page
  - Add a method `Offset() int` that returns `(Page - 1) * PerPage`
- [x] Add the `DomainEvent` interface:
  - `EventName() string` — returns the event's name
  - `OccurredAt() time.Time` — returns when the event occurred
  - `ActorId() shared.UserId` — returns the admin/user who triggered the event (audit trail)
- [x] Do NOT remove existing types (`Money`, `KV`, `UserId`, `Claim`)
- [x] Reviewed(2026-09-07): Added ActorType for differentiating the actor. This becomes part of the domain event iface.

**>>> STOP. Tell the user: "Phase 1.1 complete: shared types structs. Ready for review." Wait for the user to say "proceed".**

### 1.2 — Shared Type Tests

- [x] Create `internal/shared/types_test.go`
- [x] Test `Pagination.Offset()`:
  - `Page=1, PerPage=10` → Offset = 0
  - `Page=3, PerPage=20` → Offset = 40
  - `Page=1, PerPage=1` → Offset = 0

**>>> STOP. Wait for user to say "proceed".**

### 1.3 — Shared Type Implementation Verification

- [x] The `Pagination` struct and `Offset()` method should already be implemented from 1.1. Verify tests pass.
- [x] Run: `go test ./internal/shared/...`

**>>> STOP. Tell the user: "Phase 1 complete: shared types. Next item: Product BC Value Types." Wait for user to say "proceed".**
Reviewed: 2026-09-07
---

## Phase 2: Value Types

### 2.1 — Value Type Structs

- [x] Create (or replace entirely) `internal/product/product.go` with the file header and package doc comment.
- [x] Define the following value types in `product.go`:
  **ID types** (type aliases — ADR-003 uses `int64`):
  - `type ProId int`
  - `type ProVariantId int`
  - `type CatId int`
  - `type MediaId int`  
  **ProductStatus** (typed string enum):
  - `type ProductStatus string`
  - Constants: `StatusDraft ProductStatus = "Draft"`, `StatusActive ProductStatus = "Active"`, `StatusArchived ProductStatus = "Archived"`  
  **Money**:
  - `Money` is already `shared.Money` (`int64`) in the shared package. Use `shared.Money` throughout the product package. Do NOT redefine it.  
  **Attributes** (ADR-001 — opaque dynamic attribute map):
  - Unexported helper: `type val struct { Type string; Value any }`
  - `type Attributes struct` with unexported field `attrs map[string]*val`
  - Factory: `NewAttributes() *Attributes` — initialises empty map
  - Typed Add methods (each sets `Type` + `Value` in the map):
    - `AddStr(name, value string)` — Type=`"str"`
    - `AddNum(name string, value float64)` — Type=`"num"`
    - `AddRange(name string, min, max int)` — Type=`"range"`, Value=`"min,max"` formatted string
    - `AddEnum(name string, values ...string)` — Type=`"enum"`, Value=`"A|B|C"` pipe-joined
  - Typed Get methods (return typed value + error if missing/wrong type):
    - `GetStr(name string) (string, error)`
    - `GetNum(name string) (float64, error)`
    - `GetRange(name string) (min, max int, err error)` — parses `"min,max"` back
    - `GetEnum(name string) ([]string, error)` — splits pipe-separated value
  - General methods:
    - `Remove(name string)`
    - `Names() []string` — sorted keys
    - `Has(name string) bool`
  - Serialised as JSON for DB storage. Repo uses Add methods to hydrate from `(name, type, value)` rows or JSONB.

**>>> STOP. Tell the user: "Phase 2.1 complete: value type structs and constructors. Ready for review." Wait for the user to say "proceed".**

#### 2.2 — Value Type Tests

- [x] Create `internal/product/product_test.go` with file header
- [x] Write tests for:
  **TestAttributes**:
  - `"new/should be empty"` — `NewAttributes().Names()` → empty slice
  - `"AddStr and GetStr"` — `AddStr("Color","Red")` → `GetStr("Color")` returns `"Red"`, nil
  - `"AddNum and GetNum"` — `AddNum("Weight",1.5)` → `GetNum("Weight")` returns `1.5`, nil
  - `"AddRange and GetRange"` — `AddRange("Size",1,10)` → `GetRange("Size")` returns `1, 10, nil`
  - `"AddEnum and GetEnum"` — `AddEnum("Size","S","M","L")` → `GetEnum("Size")` returns `["S","M","L"]`, nil
  - `"GetStr missing/should error"` — `GetStr("nope")` returns error
  - `"Remove/should delete"` — Add then Remove → `Has()` returns false
  - `"Names/should return sorted"` — Add `"Z"`, `"A"`, `"M"` → `Names()` returns `["A","M","Z"]`
  - `"overwrite/should replace"` — `AddStr("X","a")` then `AddNum("X",1)` → `GetNum("X")` works, `GetStr("X")` errors
  **TestProductStatus constants**:
  - Verify `StatusDraft`, `StatusActive`, `StatusArchived` have correct string values
### 2.3 — Value Type Implementation & Test Run

- [x] Implementations should already exist from 2.1. Make any adjustments needed to pass the tests.
- [x] Run: `go test ./internal/product/...`

**>>> STOP. Tell the user: "Phase 2 complete: value types. Next item: Entities." Wait for user to say "proceed".**

---

## Phase 3: Entities

#### 3.1 — Entity Structs

All entities go in `internal/product/product.go` (append after value types).

**Product** (Aggregate Root):

- [x] Struct with unexported fields matching the spec: `id ProId`, `title string`, `slug string`, `description string`, `status ProductStatus`, `categoryId CatId`, `tags []string`, `isDigital bool`, `createdAt time.Time`, `updatedAt time.Time`
- [x] Public getters for all fields: `Id()`, `Title()`, `Slug()`, `Description()`,`SKU()`, `Status()`, `CategoryId()`, `Tags()`, `IsDigital()`, `CreatedAt()`, `UpdatedAt()`
- [x] Constructor: `NewProduct(title, description string, categoryId CatId, isDigital bool, tags []string) (*Product, error)`:
  - Sets `status` to `StatusDraft`
  - Sets `createdAt` and `updatedAt` to `time.Now()`
  - Validates that `title` is not empty (return error if so)
  - `id` is left as zero value (DB will assign)
  - `slug` is left empty (SlugGenerator will assign before save)
- [x] Domain method: `Publish() error` — transitions Draft → Active. Returns `ErrInvalidTransition` otherwise.
- [x] Domain method: `Archive() error` — transitions Active → Archived. Returns `ErrInvalidTransition` otherwise.
- [x] Domain method: `SetSlug(slug string)` — sets the slug field.
- [x] Domain method: `Update(cmd UpdateProductCmd)` — updates mutable fields from the command struct. Refreshes `updatedAt`.
- [x] Command struct `UpdateProductCmd`:
  - `Title *string` (optional — pointer means "if non-nil, update")
  - `Description *string`
  - `CategoryId *CatId`
  - `Tags *[]string`
  - `IsDigital *bool`

**ProductVariant**:

- [x] Struct with unexported fields: `id ProVariantId`, `productId ProId`, `sku string`, `price shared.Money`, `compareAtPrice *shared.Money` (pointer for nullable), `stockQty int`, `attrs *Attributes`, `mediaId *MediaId` (pointer for nullable)
- [x] Public getters for all fields
- [x] Constructor: `NewProductVariant(productId ProId, sku string, price shared.Money, attrs *Attributes) (*ProductVariant, error)`:
  - Validates `sku` is not empty
  - Validates `price >= 0`
  - Sets `stockQty` to 0
  - `id` is left as zero value (DB assigns)
- [x] Domain method: `AdjustStock(delta int) error` — adds delta to stockQty; returns `ErrInsufficientStock` if result < 0
- [x] Domain method: `Update(cmd UpdateVariantCmd)` — updates mutable fields from command struct
- [x] Command struct `UpdateVariantCmd`:
  - `SKU *string`
  - `Price *shared.Money`
  - `CompareAtPrice **shared.Money` (double pointer: outer nil = "don't touch", inner nil = "set to null")
  - `Attrs *Attributes`
  - `MediaId **MediaId`
- [x] Helper to check if price changed: `SetPrice(newPrice shared.Money) (oldPrice shared.Money, changed bool)` — sets price and returns old + whether it changed

**Category**:

- [x] Struct with unexported fields: `id CatId`, `name string`, `slug string`, `parentId *CatId` (nullable for root — ADR-003), `path string`, `depth int` (ADR-003), `sortOrder int` (ADR-003), `attrs *Attributes`
- [x] Public getters for all fields: `Id()`, `Name()`, `Slug()`, `ParentId()`, `Path()`, `Depth()`, `SortOrder()`, `Attrs()`
- [x] Constructor: `NewCategory(name, slug string, parentPath string, parentId *CatId) (*Category, error)`:
  - Validates `name` is not empty
  - Computes `depth` from `parentPath`: if empty → depth=0 (root); else depth = count dots in parentPath + 1
  - `path` is set to empty at construction (computed post-save when ID is known — see Spec Deviations)
  - Sets `parentId` to the given parent pointer (nil for root)
  - Sets `sortOrder` to 0 (caller/service can adjust)
  - Sets `attrs` to `NewAttributes()` (empty)
- [x] Domain method: `Update(cmd UpdateCatCmd)` — updates mutable fields
- [x] Domain method: `SetPath(path string)` — sets the path after DB assigns the ID
- [x] Domain method: `SetSortOrder(order int)` — sets sort order
- [x] Command struct `UpdateCatCmd`: `Name *string`, `Slug *string`, `Attrs *Attributes`

> **Note (ADR-003)**: `depth` and `parentId` can also be derived from `path`, but storing them enables efficient queries without parsing. `sortOrder` controls sibling display order in category menus.

**CatNode** (tree helper, per ADR-003):

- [x] `type CatNode struct` with exported fields: `Cat Category`, `Parent *CatNode` (ADR-003 algorithm), `Children []*CatNode`
- [x] Function `BuildCatTree(cats []Category) []*CatNode` — converts a flat list of categories into a tree. Use the algorithm from ADR-003: build a cache map, link each node to its parent via `parentId`, and collect children. Returns root nodes (those with nil parent).

**ProductMedia**:

- [x] Struct with unexported fields: `id MediaId`, `productId ProId`, `uri string`, `altText string`, `order int`
- [x] Public getters for all fields
- [x] Constructor: `NewProductMedia(productId ProId, uri, altText string, order int) *ProductMedia`
  - **Note (ADR-002)**: `uri` must be a relative path or direct URL. Never use `file:///` prefix. The `ImageStore.Store` adapter returns paths in this format.

**>>> STOP. Tell the user: "Phase 3.1 complete: entity structs. Ready for review." Wait for the user to say "proceed".**

### 3.2 — Entity Tests

- [x] Add to `internal/product/product_test.go`:

**TestNewProduct**:

- [x] `"valid input/should create draft product"` — verify status=Draft, title set, createdAt not zero
- [x] `"empty title/should return error"` — verify error returned

**TestProduct_Publish**:

- [x] `"draft product/should transition to active"` — NewProduct → Publish → status=Active
- [x] `"active product/should return ErrInvalidTransition"` — NewProduct → Publish → Publish → error
- [x] `"archived product/should return ErrInvalidTransition"` — NewProduct → Publish → Archive → Publish → error

**TestProduct_Archive**:

- [x] `"active product/should transition to archived"` — NewProduct → Publish → Archive → status=Archived
- [x] `"draft product/should return ErrInvalidTransition"` — NewProduct → Archive → error

**TestProduct_Update**:

- [x] `"update title/should set title and refresh updatedAt"` — create, sleep briefly, update with Title ptr, verify new title and updatedAt changed
- [x] `"nil fields/should not change"` — create, update with all-nil UpdateProductCmd, verify nothing changed

**TestProduct_SetSlug**:

- [x] `"should set slug"` — NewProduct → SetSlug("my-slug") → Slug() == "my-slug"

**TestNewProductVariant**:

- [x] `"valid input/should create variant"` — verify sku, price, stockQty=0
- [x] `"empty sku/should return error"`
- [x] `"negative price/should return error"`

**TestProductVariant_AdjustStock**:

- [x] `"positive delta/should increase stock"` — start 0, adjust +10 → stockQty=10
- [x] `"negative delta within range/should decrease"` — stockQty=10, adjust -5 → stockQty=5
- [x] `"negative delta below zero/should return ErrInsufficientStock"` — stockQty=3, adjust -5 → error

**TestProductVariant_SetPrice**:

- [x] `"changed price/should return old and true"` — price=100, SetPrice(200) → old=100, changed=true
- [x] `"same price/should return old and false"` — price=100, SetPrice(100) → old=100, changed=false

**TestNewCategory**:

- [x] `"valid input/should create category"` — verify name, slug, empty attrs, depth=0 for root, parentId=nil for root
- [x] `"child category/should compute depth from parent path"` — parentPath="1.2" → depth=2
- [x] `"empty name/should return error"`

**TestCategory_SetPath**:

- [x] `"should set path"` — NewCategory → SetPath("1.2.3") → Path() == "1.2.3"

**TestCategory_SetSortOrder**:

- [x] `"should set sort order"` — NewCategory → SetSortOrder(5) → SortOrder() == 5

**TestBuildCatTree**:

- [x] `"flat list/should build correct hierarchy with parent links"` — Given categories with paths "1", "1.2", "1.3" → root has 2 children, each child's Parent points to root
- [x] `"single root/should have nil parent"` — verify root CatNode has Parent == nil

**TestNewProductMedia**:

- [x] `"should create media with all fields"` — verify uri, altText, order

**>>> STOP. Wait for user to say "proceed".**

### 3.3 — Entity Implementation & Test Run

- [x] Ensure all entity code from 3.1 is complete. Adjust to pass tests.
- [x] Run: `go test ./internal/product/...`

**>>> STOP. Tell the user: "Phase 3 complete: entities. Next item: Errors & Ports." Wait for user to say "proceed".**

---

## Phase 4: Errors & Ports

### 4.1 — Error Sentinels & Port Interfaces

**errors.go**:

- [ ] Create `internal/product/errors.go` with file header
- [ ] Define all sentinel errors:
  ```
  var (
      ErrProductNotFound       = errors.New("product not found")
      ErrVariantNotFound       = errors.New("variant not found")
      ErrDuplicateSKU          = errors.New("duplicate SKU within product")
      ErrInvalidTransition     = errors.New("invalid status transition")
      ErrInsufficientStock     = errors.New("insufficient stock")
      ErrCategoryNotFound      = errors.New("category not found")
      ErrSlugConflict          = errors.New("slug already exists")
      ErrVariantHasActiveOrders = errors.New("variant has active orders")
      ErrCategoryInUse         = errors.New("category is assigned to products")
      ErrMediaNotFound         = errors.New("media not found")
  )
  ```

**ports.go**:

- [ ] Create `internal/product/ports.go` with file header
- [ ] Define all driven port interfaces exactly as specified in the domain-model.md:
  **ProductRepo** interface:
  - `Save(ctx context.Context, product *Product) error`
  - `FindById(ctx context.Context, id ProId) (*Product, error)`
  - `FindBySlug(ctx context.Context, slug string) (*Product, error)`
  - `List(ctx context.Context, filter ProductFilter, page shared.Pagination) ([]*Product, error)`
  - `Delete(ctx context.Context, id ProId) error`  
  **Note**: You will also need to define `ProductFilter` struct. Define it here or in `product.go`. It should include optional filter fields like `Status *ProductStatus`, `CategoryId *CatId`, `Tags []string`, `SearchQuery *string`.  
  **VariantRepo** interface:
  - `Save(ctx context.Context, variant *ProductVariant) error`
  - `FindById(ctx context.Context, id ProVariantId) (*ProductVariant, error)`
  - `FindByProductId(ctx context.Context, productId ProId) ([]*ProductVariant, error)`
  - `FindBySKU(ctx context.Context, sku string, productId ProId) (*ProductVariant, error)`
  - `Delete(ctx context.Context, id ProVariantId) error`  
  **CategoryRepo** interface:
  - `Save(ctx context.Context, cat *Category) error`
  - `GetCatBySlug(ctx context.Context, slug string) (*Category, error)`
  - `GetCatById(ctx context.Context, id CatId) (*Category, error)`
  - `GetDescendCats(ctx context.Context, catId CatId) ([]*Category, error)`
  - `GetChildren(ctx context.Context, id CatId) ([]*Category, error)`
  - `GetAncestors(ctx context.Context, id CatId) ([]*Category, error)`
  - `GetSubtree(ctx context.Context, id CatId) ([]*Category, error)`
  - `Move(ctx context.Context, id CatId, parentId CatId) error`
  - `GetDepth(ctx context.Context, depth int) ([]*Category, error)`
  - `Delete(ctx context.Context, id CatId) error`  
  **MediaRepo** interface:
  - `Save(ctx context.Context, media *ProductMedia) error`
  - `FindById(ctx context.Context, id MediaId) (*ProductMedia, error)`
  - `FindByProId(ctx context.Context, productId ProId) ([]*ProductMedia, error)`
  - `Delete(ctx context.Context, id MediaId) error`  
  **ImageStore** interface:
  - `Store(ctx context.Context, file io.Reader, filename string) (string, error)`
  - `Delete(ctx context.Context, filePath string) error`  
  **OrderQueryPort** interface:
  - `HasActiveOrdersForVariant(ctx context.Context, variantId ProVariantId) (bool, error)`  
  **EventPublisher** interface:
  - `Publish(ctx context.Context, event shared.DomainEvent) error`

**>>> STOP. Tell the user: "Phase 4.1 complete: errors and port interfaces. Ready for review." Wait for the user to say "proceed".**

### 4.2 — Port Interface Tests (Compile Check)

- [ ] No behavioral tests for interfaces, but add compile-time interface compliance checks at the bottom of `product_test.go`:
  ```go
  // Compile-time interface checks — these will be satisfied when adapters are implemented.
  // For now, they document the expected contracts.
  ```
  This is optional since adapters are out of scope for this plan.

**>>> STOP. Wait for user to say "proceed".**

### 4.3 — Verify Compilation

- [ ] Run: `go build ./internal/product/...` to confirm no compilation errors
- [ ] Run: `go test ./internal/product/...` to confirm all existing tests still pass

**>>> STOP. Tell the user: "Phase 4 complete: errors and ports. Next item: Domain Services." Wait for user to say "proceed".**

---

## Phase 5: Domain Services

### 5.1 — Domain Service Structs

**SlugGenerator** — `internal/product/slug_service.go`:

- [ ] Create with file header
- [ ] Define `type SlugGenerator struct{}` (stateless — no dependencies stored; repo passed per call)
- [ ] Function `NewSlugGenerator() *SlugGenerator`
- [ ] Method `Generate(ctx context.Context, title string, checkSlug func(ctx context.Context, slug string) (*Product, error)) (string, error)`:
  - Convert title to lowercase
  - Replace spaces with hyphens
  - Remove non-alphanumeric characters (keep hyphens)
  - Collapse consecutive hyphens
  - Trim leading/trailing hyphens
  - Check uniqueness using the `checkSlug` callback (which wraps `ProductRepo.FindBySlug`)
  - If a product with that slug exists, append a short suffix (e.g., `-2`, `-3`, etc., incrementally trying until unique)
  - Return the unique slug
  - Note: We use a function callback `checkSlug` rather than the full `ProductRepo` interface to keep the domain service decoupled. The application service passes `repo.FindBySlug` as the callback.

**StockAdjuster** — `internal/product/stock_service.go`:

- [ ] Create with file header
- [ ] Define `type StockAdjuster struct{}` (stateless)
- [ ] Function `NewStockAdjuster() *StockAdjuster`
- [ ] Method `Adjust(ctx context.Context, variant *ProductVariant, delta int, pub EventPublisher) error`:
  - Call `variant.AdjustStock(delta)` — if error, return it
  - Create a `StockAdjustedEvent` struct (implements `shared.DomainEvent`):
    - Fields: `variantId ProVariantId`, `delta int`, `newQty int`, `occurredAt time.Time`
    - `EventName() string` → `"StockAdjusted"`
    - `OccurredAt() time.Time` → the timestamp
  - Publish the event via `pub.Publish(ctx, event)`
  - Return any error from publishing

**Domain Events** (define in `product.go` or a separate `events.go` — your choice, but keep in same package):

- [ ] `type StockAdjustedEvent struct`:
  - `ProductId ProId`, `VariantId ProVariantId`, `Delta int`, `NewQty int`, `actorId shared.UserId`, `occurredAt time.Time`
  - Implements `shared.DomainEvent`
- [ ] `type ProductPublishedEvent struct`:
  - `ProductId ProId`, `Title string`, `Slug string`, `actorId shared.UserId`, `occurredAt time.Time`
- [ ] `type ProductArchivedEvent struct`:
  - `ProductId ProId`, `actorId shared.UserId`, `occurredAt time.Time`
- [ ] `type VariantPriceChangedEvent struct`:
  - `VariantId ProVariantId`, `OldPrice shared.Money`, `NewPrice shared.Money`, `actorId shared.UserId`, `occurredAt time.Time`

All events implement `shared.DomainEvent` interface. The `actorId`  and `actorType` is extracted from `ctx` (via `shared.Claim`) by the app service before constructing the event.

**>>> STOP. Tell the user: "Phase 5.1 complete: domain service structs and event types. Ready for review." Wait for the user to say "proceed".**

### 5.2 — Domain Service Tests

**`internal/product/slug_service_test.go`**:

- [ ] Create with file header
- [ ] `TestSlugGenerator_Generate`:
  - `"simple title/should slugify"` — `"My Product"` → `"my-product"` (with checkSlug returning nil/not-found)
  - `"special characters/should strip"` — `"Hello! World @#$"` → `"hello-world"`
  - `"consecutive spaces/should collapse hyphens"` — `"A   B"` → `"a-b"`
  - `"slug collision/should append suffix"` — First call to checkSlug returns a product (exists), second returns nil → slug = `"my-product-2"` or similar
  - `"empty title/should return error"`
  - For the mock `checkSlug` function: create a simple closure that returns a `*Product` on first call and `nil, ErrProductNotFound` on subsequent calls (or use a counter).

**`internal/product/stock_service_test.go`**:

- [ ] Create with file header
- [ ] Create a mock `EventPublisher` (a struct implementing the interface that records published events)
- [ ] `TestStockAdjuster_Adjust`:
  - `"positive delta/should increase stock and publish event"` — variant with stockQty=5, delta=+3 → stockQty=8, mock publisher has 1 event with correct fields
  - `"negative delta within range/should decrease and publish"` — stockQty=10, delta=-3 → stockQty=7, event published
  - `"negative delta below zero/should return error and not publish"` — stockQty=2, delta=-5 → `ErrInsufficientStock`, mock publisher has 0 events
  - `"zero delta/should publish event"` — stockQty=5, delta=0 → stockQty=5, event published

**>>> STOP. Wait for user to say "proceed".**

### 5.3 — Domain Service Implementation & Test Run

- [ ] Implementations from 5.1 should already exist. Adjust to pass tests.
- [ ] Run: `go test ./internal/product/...`

**>>> STOP. Tell the user: "Phase 5 complete: domain services. Next item: Application Services." Wait for user to say "proceed".**

---

## Phase 6: Application Services

### 6.1 — ProductSvc Struct

**`internal/product/product_service.go`**:

- [ ] Create with file header
- [ ] Define `ProductSvc` struct with dependencies:
  ```
  productRepo  ProductRepo
  variantRepo  VariantRepo
  mediaRepo    MediaRepo
  imageStore   ImageStore
  orderQuery   OrderQueryPort
  eventPub     EventPublisher
  slugGen      *SlugGenerator
  stockAdj     *StockAdjuster
  logger       *zerolog.Logger
  ```
- [ ] Constructor: `NewProductSvc(productRepo ProductRepo, variantRepo VariantRepo, mediaRepo MediaRepo, imageStore ImageStore, orderQuery OrderQueryPort, eventPub EventPublisher, slugGen *SlugGenerator, stockAdj *StockAdjuster, logger *zerolog.Logger) *ProductSvc`
- [ ] Define `ProductCreateInput` struct: `Title string`, `Description string`, `CategoryId CatId`, `IsDigital bool`, `Tags []string`
- [ ] Define `AddVariantInput` struct: `SKU string`, `Price shared.Money`, `Attrs *Attributes`
- [ ] Implement methods:
  **CreateProduct(ctx, input ProductCreateInput) (Product, error)**:
  - Call `NewProduct(input.Title, input.Description, input.CategoryId, input.IsDigital, input.Tags)`
  - Generate slug via `slugGen.Generate(ctx, input.Title, productRepo.FindBySlug)`
  - Set slug on product
  - Save via `productRepo.Save`
  - Log: `logger.Info().Int("productId", int(product.Id())).Str("slug", product.Slug()).Msg("product created")`
  - Return product by value  
  **UpdateProduct(ctx, productId ProId, cmd UpdateProductCmd) (Product, error)**:
  - Load product via `productRepo.FindById`
  - If cmd.Title is set and changed, re-generate slug
  - Call `product.Update(cmd)`
  - Save via `productRepo.Save`
  - Log: `logger.Info().Int("productId", int(productId)).Msg("product updated")`
  - Return product by value  
  **PublishProduct(ctx, productId ProId) error**:
  - Load product
  - Check at least 1 variant exists via `variantRepo.FindByProductId`
  - If no variants: `logger.Warn().Int("productId", int(productId)).Msg("publish blocked: no variants")`; return error
  - Call `product.Publish()`
  - Save product
  - Publish `ProductPublishedEvent` (extract actorId from ctx)
  - Log: `logger.Info().Int("productId", int(productId)).Msg("product published")`  
  **ArchiveProduct(ctx, productId ProId) error**:
  - Load product
  - Call `product.Archive()`
  - Save product
  - Publish `ProductArchivedEvent` (extract actorId from ctx)
  - Log: `logger.Info().Int("productId", int(productId)).Msg("product archived")`  
  **AddVariant(ctx, productId ProId, input AddVariantInput) (ProductVariant, error)**:
  - Check product exists
  - Check SKU uniqueness via `variantRepo.FindBySKU`
  - Create variant via `NewProductVariant`
  - Save via `variantRepo.Save`
  - Return variant by value  
  **UpdateVariant(ctx, variantId ProVariantId, cmd UpdateVariantCmd) error**:
  - Load variant
  - Track if price changed (for event)
  - Call `variant.Update(cmd)`
  - Save
  - If price changed, publish `VariantPriceChangedEvent` (extract actorId from ctx)
  - Log: `logger.Info().Int("variantId", int(variantId)).Msg("variant updated")`  
  **RemoveVariant(ctx, variantId ProVariantId) error**:
  - Check `orderQuery.HasActiveOrdersForVariant` → if true, return `ErrVariantHasActiveOrders`
  - Delete via `variantRepo.Delete`
  - Log: `logger.Info().Int("variantId", int(variantId)).Msg("variant removed")`  
  **AdjustStock(ctx, variantId ProVariantId, delta int) error**:
  - Load variant
  - Call `stockAdj.Adjust(ctx, variant, delta, eventPub)` (StockAdjuster constructs event with actorId from ctx)
  - Save variant
  - Log: `logger.Info().Int("variantId", int(variantId)).Int("delta", delta).Msg("stock adjusted")`  
  **ReorderImages(ctx, productId ProId, orderedIds []MediaId) error**:
  - Load all media for product
  - For each orderedId, find the matching media and update its `order` field to the index position
  - Save each updated media  
  **RemoveImage(ctx, imageId MediaId) error**:
  - Load media by ID
  - Delete from `imageStore.Delete(media.uri)`
  - Delete from `mediaRepo.Delete`

**>>> STOP. Tell the user: "Phase 6.1 complete: ProductSvc struct and methods. Ready for review." Wait for the user to say "proceed".**

### 6.2 — ProductSvc Tests

**`internal/product/product_service_test.go`**:

- [ ] Create with file header
- [ ] Create mock implementations for all ports:
  - `mockProductRepo` — in-memory map[ProId]*Product
  - `mockVariantRepo` — in-memory map[ProVariantId]*ProductVariant
  - `mockMediaRepo` — in-memory map[MediaId]*ProductMedia
  - `mockImageStore` — records Store/Delete calls
  - `mockOrderQueryPort` — configurable return value
  - `mockEventPublisher` — records published events
- [ ] Helper function to create a pre-configured `ProductSvc` with all mocks

**Test cases**:

`TestProductSvc_CreateProduct`:

- `"valid input/should create draft product with slug"` — verify returned product has Draft status, non-empty slug
- `"empty title/should return error"` — verify error

`TestProductSvc_PublishProduct`:

- `"draft with variants/should publish"` — create product, add variant, then publish → status Active, event published
- `"draft without variants/should return error"` — create product, publish → error (no variants)

`TestProductSvc_ArchiveProduct`:

- `"active product/should archive"` — create, add variant, publish, archive → status Archived

`TestProductSvc_AddVariant`:

- `"valid input/should add variant"` — verify variant saved, SKU correct
- `"duplicate SKU/should return ErrDuplicateSKU"` — add variant with same SKU twice

`TestProductSvc_RemoveVariant`:

- `"no active orders/should remove"` — mock OrderQueryPort returns false → variant deleted
- `"has active orders/should return ErrVariantHasActiveOrders"` — mock returns true → error

`TestProductSvc_AdjustStock`:

- `"positive delta/should adjust and publish event"` — verify stock changed, event published
- `"insufficient stock/should return error"` — verify error, no event

`TestProductSvc_ReorderImages`:

- `"valid order/should reorder"` — create 3 media, reorder, verify new order values

`TestProductSvc_RemoveImage`:

- `"existing image/should delete from store and repo"` — verify both mock calls

**>>> STOP. Wait for user to say "proceed".**

### 6.3 — ProductSvc Implementation & Test Run

- [ ] Adjust implementation to pass all tests
- [ ] Run: `go test ./internal/product/...`

**>>> STOP. Tell the user: "Phase 6 (ProductSvc) complete. Next item: ProductMediaSvc." Wait for user to say "proceed".**

---

## Phase 7: ProductMediaSvc

### 7.1 — ProductMediaSvc Struct

**`internal/product/media_service.go`**:

- [ ] Create with file header
- [ ] Define `ProductMediaSvc` struct with dependencies: `mediaRepo MediaRepo`, `imageStore ImageStore`, `productRepo ProductRepo`, `logger *zerolog.Logger`
- [ ] Constructor: `NewProductMediaSvc(mediaRepo, imageStore, productRepo, logger)`
- [ ] Implement methods:
  **UploadImage(ctx, productId ProId, fileStream io.Reader, altText string) (ProductMedia, error)**:
  - Verify product exists via `productRepo.FindById`
  - Store file via `imageStore.Store`
  - Determine next order number by counting existing media via `mediaRepo.FindByProId`
  - Create `ProductMedia` via `NewProductMedia`
  - Save via `mediaRepo.Save`
  - Return by value  
  **UploadVideo(ctx, productId ProId, fileStream io.Reader) (ProductMedia, error)**:
  - Same as UploadImage but with empty altText (or a default)
  - Store file via `imageStore.Store`
  - Create and save media

**>>> STOP. Tell the user: "Phase 7.1 complete: ProductMediaSvc. Ready for review." Wait for the user to say "proceed".**

### 7.2 — ProductMediaSvc Tests

**`internal/product/media_service_test.go`**:

- [ ] Reuse the mock implementations from Phase 6
- [ ] `TestProductMediaSvc_UploadImage`:
  - `"valid product/should store and save media"` — verify imageStore.Store called, mediaRepo.Save called, returned media has correct URI
  - `"invalid product/should return error"` — product not found
- [ ] `TestProductMediaSvc_UploadVideo`:
  - `"valid product/should store and save"` — verify same as image

**>>> STOP. Wait for user to say "proceed".**

### 7.3 — Implementation & Test Run

- [ ] Run: `go test ./internal/product/...`

**>>> STOP. Tell the user: "Phase 7 complete: ProductMediaSvc. Next item: CategorySvc." Wait for user to say "proceed".**

---

## Phase 8: CategorySvc

### 8.1 — CategorySvc Struct

**`internal/product/category_service.go`**:

- [ ] Create with file header
- [ ] Define `CategorySvc` struct with dependencies: `catRepo CategoryRepo`, `productRepo ProductRepo`, `logger *zerolog.Logger`
- [ ] Constructor: `NewCategorySvc(catRepo, productRepo, logger)`
- [ ] Define `UpdateCatInput` struct if not already defined: `Name *string`, `Slug *string`, `Attrs *Attributes`
- [ ] Implement methods:
  **CreateCat(ctx, catName string, parentCategory CatId, subCategories []CatId) (Category, error)**:
  - If parentCategory is not empty, load parent via `catRepo.GetCatById` (verify exists)
  - Create category via `NewCategory`
  - Save via `catRepo.Save` (repo assigns ID)
  - Compute path: if parent exists, path = `parentPath + "." + string(newCat.Id())`; if root, path = `string(newCat.Id())`
  - Update category path and re-save
  - Handle subCategories: move them under the new category (optional — depends on whether repo assigns IDs on Save)  
  **UpdateCat(ctx, categoryId CatId, input UpdateCatInput) error**:
  - Load category
  - Apply updates
  - Save  
  **DeleteCat(ctx, categoryId CatId) error**:
  - Check if any products use this category (via `productRepo.List` with category filter)
  - If products exist, return `ErrCategoryInUse`
  - Delete via `catRepo.Delete`  
  **GetCatById(ctx, catId CatId) (Category, error)**:
  - Delegate to `catRepo.GetCatById`
  - Return by value  
  **GetCatBySlug(ctx, slug string) (Category, error)**:
  - Delegate to `catRepo.GetCatBySlug`
  - Return by value  
  **GetDescendents(ctx, cat Category) ([]Category, error)**:
  - Delegate to `catRepo.GetDescendCats(ctx, cat.Id())`
  - Return by value (convert []*Category to []Category)

**>>> STOP. Tell the user: "Phase 8.1 complete: CategorySvc. Ready for review." Wait for the user to say "proceed".**

### 8.2 — CategorySvc Tests

**`internal/product/category_service_test.go`**:

- [ ] Create mock `CategoryRepo` (in-memory)
- [ ] `TestCategorySvc_CreateCat`:
  - `"root category/should create with path"` — verify category created, path set
  - `"child category/should compute path from parent"`
- [ ] `TestCategorySvc_DeleteCat`:
  - `"no products assigned/should delete"` — mock productRepo returns empty list
  - `"products assigned/should return ErrCategoryInUse"` — mock returns non-empty list
- [ ] `TestCategorySvc_GetCatById`:
  - `"existing category/should return"` — verify
  - `"non-existing/should return error"`

**>>> STOP. Wait for user to say "proceed".**

### 8.3 — Implementation & Test Run

- [ ] Run: `go test ./internal/product/...`

**>>> STOP. Tell the user: "Phase 8 complete: CategorySvc. Next item: ProductQuerySvc." Wait for user to say "proceed".**

---

## Phase 9: ProductQuerySvc

### 9.1 — ProductQuerySvc Struct

**`internal/product/query_service.go`**:

- [ ] Create with file header
- [ ] Define `ProductQuerySvc` struct with dependencies: `productRepo ProductRepo`, `variantRepo VariantRepo`
- [ ] Constructor: `NewProductQuerySvc(productRepo, variantRepo)`
- [ ] Implement methods (all are thin delegations to repos):
  **GetProduct(ctx, productId ProId) (Product, error)* — delegates to `productRepo.FindById`  
  **GetProducts(ctx, filter ProductFilter, page shared.Pagination) ([]Product, error)* — delegates to `productRepo.List`  
  **GetVariant(ctx, variantId ProVariantId) (ProductVariant, error)* — delegates to `variantRepo.FindById`  
  **GetStockLevel(ctx, variantId ProVariantId) (int, error)** — loads variant, returns `variant.StockQty()`

**>>> STOP. Tell the user: "Phase 9.1 complete: ProductQuerySvc. Ready for review." Wait for the user to say "proceed".**

### 9.2 — ProductQuerySvc Tests

**`internal/product/query_service_test.go`**:

- [ ] Reuse mocks
- [ ] `TestProductQuerySvc_GetProduct`:
  - `"existing product/should return"` — verify
  - `"non-existing/should return error"`
- [ ] `TestProductQuerySvc_GetStockLevel`:
  - `"existing variant/should return stock"` — create variant with known stock, verify returned value

**>>> STOP. Wait for user to say "proceed".**

### 9.3 — Implementation & Test Run

- [ ] Run: `go test ./internal/product/...`

**>>> STOP. Tell the user: "Phase 9 complete: ProductQuerySvc. All application services done. Running final verification." Wait for user to say "proceed".**

---

## Phase 10: Final Verification

- [ ] Run: `go build ./internal/product/...` — must compile with zero errors
- [ ] Run: `go test ./internal/product/...` — all tests must pass
- [ ] Run: `go test ./internal/shared/...` — shared tests must still pass
- [ ] Run: `go vet ./internal/product/...` — no vet warnings
- [ ] Verify file layout matches the table in the Preamble
- [ ] Verify all acceptance criteria from the domain-model.md spec are addressed by the implementation (cross-reference each checkbox)

**>>> STOP. Tell the user: "Implementation plan execution complete. All phases done. Ready for architect review."**

---

## Spec Deviations

| Item                               | Deviation                                                                             | Rationale                                                                                                             |
| ---------------------------------- | ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| `SlugGenerator.Generate` signature | Uses `checkSlug` callback function instead of `ProductRepo` interface                 | Keeps domain service decoupled from the full repo interface; application service passes `repo.FindBySlug` as callback |
| `Category` constructor             | Path is computed post-save (since ID is DB-assigned) rather than at construction time | ID is not known until after `Save`; path computation requires the ID                                                  |
| `ProductFilter` struct             | Not in spec — added as a required parameter type for `ProductRepo.List`               | Needed to compile the repo interface; exact fields are a filtering concern                                            |
| `ErrMediaNotFound`                 | Not in spec — added                                                                   | Needed for `RemoveImage` and `UploadImage` error paths                                                                |

---

## ADR Revision Log (2026-09-07)

The initial plan was generated reading `adr-archived.md` (legacy) instead of the individual reviewed ADR files (ADR001–ADR009). The following corrections were applied:

| ADR     | Change Applied                                                                                                  |
| ------- | --------------------------------------------------------------------------------------------------------------- |
| All     | Added Step 0 instruction to read ADR001–ADR009 individually; explicit `adr-archived.md` exclusion               |
| ADR-001 | Added note: category `attrs` define attribute specifications inherited by products (filtering)                  |
| ADR-002 | Added convention #13: media URIs use relative paths or direct URLs, no `file:///` prefix                        |
| ADR-003 | Added `parentId *CatId`, `depth int`, `sortOrder int` to Category entity; `SetPath()`, `SetSortOrder()` methods |
| ADR-003 | Updated `CatNode` to include `Parent *CatNode` link alongside `Children`; updated `BuildCatTree` algorithm      |
| ADR-003 | Added tests: `TestCategory_SetPath`, `TestCategory_SetSortOrder`, child depth test, parent-link tree tests      |
| ADR-004 | Added convention #12: money values in configured least minor units (cents/paise), no floating-point             |
