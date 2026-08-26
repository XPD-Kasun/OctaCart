# Todo: Product Catalog — Core Domain Model & Ports

Spec: `.agents/specs/product/core-domain-model.md`
Decisions: `.agents/specs/product/adr.md`
Status: **Awaiting "proceed"**

---

## Phase 1 — Value Types & Entities

- [ ] Create `internal/product/product.go`
  - [ ] Define `ProductId`, `ProductStatus` (Draft/Active/Archived) types
  - [ ] Define `Product` struct with all fields
  - [ ] Add `CanPublish() bool` method — true when product has >=1 active variant
  - [ ] Add `Publish() error` / `Archive() error` state-transition methods with guard logic

- [ ] Create `internal/product/variant.go`
  - [ ] Define `VariantId`, `Attribute`, `Money` value types (Money: minor units, single currency)
  - [ ] Define `ProductVariant` struct

- [ ] Create `internal/product/category.go`
  - [ ] Define `CategoryId` type
  - [ ] Define `Category` struct

- [ ] Create `internal/product/asset.go`
  - [ ] Define `AssetId` type
  - [ ] Define `ProductAsset` struct (`uri` relative path/URL, no `file:///` prefix)

---

## Phase 2 — Ports, Errors & Events

- [ ] Create `internal/product/ports.go`
  - [ ] `ProductRepo` interface (Save, FindById, FindBySlug, List, Delete)
  - [ ] `VariantRepo` interface (Save, FindById, FindByProductId, FindBySKU, Delete)
  - [ ] `CategoryRepo` interface (Save, FindById, List, Delete, HasProducts)
  - [ ] `AssetRepo` interface (Save, FindById, FindByProductId, UpdateOrder, Delete)
  - [ ] `ImageStore` interface (Store, Delete)
  - [ ] `OrderQueryPort` interface (HasActiveOrdersForVariant)
  - [ ] `EventPublisher` interface (Publish)

- [ ] Create `internal/product/errors.go`
  - [ ] `ErrProductNotFound`, `ErrVariantNotFound`, `ErrCategoryNotFound`, `ErrAssetNotFound`
  - [ ] `ErrDuplicateSKU`, `ErrInvalidTransition`, `ErrInsufficientStock`
  - [ ] `ErrSlugConflict`, `ErrCategoryInUse`, `ErrVariantHasActiveOrders`

- [ ] Create `internal/product/events.go`
  - [ ] `DomainEvent` interface (EventName() string)
  - [ ] `ProductPublished` struct
  - [ ] `ProductArchived` struct
  - [ ] `StockAdjusted` struct
  - [ ] `VariantPriceChanged` struct

---

## Phase 3 — Domain Services

- [ ] Create `internal/product/slug.go`
  - [ ] `SlugGenerator` type
  - [ ] `Generate(title string) string` — URL-safe, lowercase, hyphenated
  - [ ] `GenerateUnique(ctx, title string, repo ProductRepo) string` — appends suffix until no conflict

- [ ] Create `internal/product/stock.go`
  - [ ] `StockAdjuster` type
  - [ ] `Adjust(variant *ProductVariant, delta int) error` — returns `ErrInsufficientStock` if result < 0

---

## Phase 4 — Application Services

- [ ] Create `internal/product/product_service.go`
  - [ ] `ProductSvc` struct with injected ports
  - [ ] Implement all methods per spec (CreateProduct → RemoveAsset); updates use `func(*Entity)` updaters per ADR 1

- [ ] Create `internal/product/asset_service.go`
  - [ ] `ProductAssetSvc` struct with injected ports (`ImageStore`, `AssetRepo`)
  - [ ] `UploadImage(ctx, productId, file, altText)`
  - [ ] `UploadVideo(ctx, productId, file)`

- [ ] Create `internal/product/category_service.go`
  - [ ] `CategorySvc` struct with injected ports
  - [ ] Implement all 4 methods per spec (updater lambda style)

- [ ] Create `internal/product/query_service.go`
  - [ ] `ProductQuerySvc` struct with injected ports
  - [ ] Implement all 5 query methods per spec (uses `shared.Pagination`)

---

## Phase 5 — Unit Tests

- [ ] Create `internal/product/product_test.go`
  - [ ] Test `Publish()` — success when variant exists
  - [ ] Test `Publish()` — error when no variants
  - [ ] Test `Publish()` — error when already Active
  - [ ] Test `Archive()` — success from Active
  - [ ] Test `Archive()` — error from Draft

- [ ] Create `internal/product/stock_test.go`
  - [ ] Test `Adjust()` with positive delta (restock)
  - [ ] Test `Adjust()` with negative delta within bounds
  - [ ] Test `Adjust()` with negative delta that would go below zero → `ErrInsufficientStock`

- [ ] Create `internal/product/slug_test.go`
  - [ ] Test slug generation: title → expected URL-safe slug
  - [ ] Test collision suffix appended

- [ ] Create `internal/product/product_service_test.go`
  - [ ] Test `AddVariant` — `ErrDuplicateSKU` on duplicate SKU
  - [ ] Test `RemoveVariant` — `ErrVariantHasActiveOrders` when port returns true
  - [ ] Test `AdjustStock` — publishes `StockAdjusted` event on success
  - [ ] Test `UpdateVariant` — publishes `VariantPriceChanged` when price differs
  - [ ] Test `ReorderAssets` — persists supplied ordering

- [ ] Create `internal/product/asset_service_test.go`
  - [ ] Test `UploadImage` — stores via `ImageStore`, persists via `AssetRepo`

- [ ] Create `internal/product/category_service_test.go`
  - [ ] Test `DeleteCategory` — `ErrCategoryInUse` when products assigned

---

## Phase 6 — Verify

- [ ] Run `go build ./internal/product/...` — zero errors
- [ ] Run `go test ./internal/product/...` — all tests pass
- [ ] Run `go vet ./internal/product/...` — zero warnings
