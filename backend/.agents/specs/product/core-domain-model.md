## Spec: Product Catalog — Core Domain Model & Ports

- **Bounded Context**: product
- **Goal**: Establish the complete domain model, port interfaces, domain events, application services, and errors for the Product BC as pure Go — zero infrastructure dependencies. This is the foundation all future adapters (REST, DB, storage) will plug into.
- **Context source**: `.agents/discover/product.md`
- **Decisions**: see `.agents/specs/product/adr.md` (ADR 1)

---

### In Scope

- Domain entities: `Product`, `ProductVariant`, `Category`, `ProductAsset`
- Value types: `Money`, `Attribute`, `ProductStatus`, `ProductId`, `VariantId`, `CategoryId`, `AssetId`
- Repository ports: `ProductRepo`, `VariantRepo`, `CategoryRepo`, `AssetRepo`
- Infrastructure ports: `ImageStore`, `OrderQueryPort`, `EventPublisher`
- Domain event struct types: `ProductPublished`, `ProductArchived`, `StockAdjusted`, `VariantPriceChanged`
- Domain services: `SlugGenerator`, `StockAdjuster`
- Application services: `ProductSvc`, `ProductAssetSvc`, `CategorySvc`, `ProductQuerySvc`
- Domain errors: all `Err*` sentinel values in `errors.go`
- Unit tests for state transitions, stock invariant, slug generation, SKU uniqueness, RemoveVariant guard

### Out of Scope

- REST handlers and GraphQL resolvers
- Any concrete repository implementation (DB adapters)
- Image file storage implementation (`LocalImageStore`)
- Event bus implementation (`InProcessEventBus`)
- `shared.Pagination` — belongs to the shared module; assumed available here as a cross-module dependency

---

### Application Services

#### `ProductSvc`
**Purpose**: The primary command handler for the product lifecycle. Orchestrates business rules across `Product` and `ProductVariant` aggregates — from creation through publishing, archiving, variant management, stock adjustment, and asset ordering. All catalog write operations except file uploads flow through this service.

| Method | Summary |
|---|---|
| `CreateProduct(ctx, input)` | Creates a new product in `Draft` status. Generates a unique slug via `SlugGenerator`. |
| `UpdateProduct(ctx, productId, updater func(*Product))` | Applies field updates via the updater lambda. Re-generates slug if title changes; checks for slug conflicts. |
| `PublishProduct(ctx, productId)` | Transitions product `Draft → Active`. Blocked if the product has no active variants. |
| `ArchiveProduct(ctx, productId)` | Transitions product `Active → Archived`. |
| `AddVariant(ctx, productId, input)` | Appends a new variant to a product. Enforces SKU uniqueness within the product. |
| `UpdateVariant(ctx, variantId, updater func(*ProductVariant))` | Applies variant field updates via the updater lambda. Publishes `VariantPriceChanged` if price changes. |
| `RemoveVariant(ctx, variantId)` | Deletes a variant. **Blocked** (hard error) if `OrderQueryPort` reports the variant has active orders. |
| `AdjustStock(ctx, variantId, delta)` | Increases or decreases `stockQty`. `StockAdjuster` enforces non-negative invariant. Publishes `StockAdjusted` event on success. |
| `ReorderAssets(ctx, productId, orderedIds)` | Updates `order` on each `ProductAsset` to match the supplied ID sequence. |
| `RemoveAsset(ctx, assetId)` | Deletes asset record and instructs `ImageStore` to remove the file. |

> ADR 1: updaters use `func(*Entity)` lambdas instead of per-field input structs — syntactic sugar, not a technical requirement; a dedicated update struct may be introduced later if field granularity is needed.

#### `ProductAssetSvc`
**Purpose**: Handles file upload concerns, which are not a responsibility of `ProductSvc`. Delegates storage to the `ImageStore` port and persists asset metadata via `AssetRepo`. "Asset" covers images and videos (the only types currently supported).

| Method | Summary |
|---|---|
| `UploadImage(ctx, productId, file, altText)` | Stores the file via `ImageStore`, persists a `ProductAsset` record with returned URI and sort order. |
| `UploadVideo(ctx, productId, file)` | Same as above for video assets. |

> ADR 1: image uploading extracted from `ProductSvc` into this service behind a driving port (`AssetRepo`) + application service.

#### `CategorySvc`
**Purpose**: Manages the category tree used to organise products. Categories are owned entirely by this BC; they are referenced (by ID) from products but their lifecycle is independent.

| Method | Summary |
|---|---|
| `CreateCategory(ctx, input)` | Creates a new category. Validates that `parentId` (if supplied) exists. |
| `UpdateCategory(ctx, categoryId, updater func(*Category))` | Applies category field updates via the updater lambda. |
| `DeleteCategory(ctx, categoryId)` | Deletes a category. **Blocked** if any product is currently assigned to it (`ErrCategoryInUse`). |
| `ListCategories(ctx)` | Returns the full category tree (flat list with parentId links). |

#### `ProductQuerySvc`
**Purpose**: The read side of the Product BC. Provides all query operations used by both the admin REST API and the storefront GraphQL endpoint. Contains no business logic — pure data retrieval through repository ports.

| Method | Summary |
|---|---|
| `GetProduct(ctx, productId)` | Fetches a single product by ID. |
| `GetProducts(ctx, filter, page shared.Pagination)` | Returns a paginated, filtered list of products. Filter includes status, categoryId, tags. |
| `SearchProducts(ctx, query, filter)` | Full-text search over product titles/descriptions; delegates to repo. |
| `GetVariant(ctx, variantId)` | Fetches a single variant by ID. |
| `GetStockLevel(ctx, variantId)` | Returns current `stockQty` for a variant. |

---

### Acceptance Criteria

- [ ] All entity types compile cleanly under `internal/product/`
- [ ] Port interfaces defined in `internal/product/ports.go`: `ProductRepo`, `VariantRepo`, `CategoryRepo`, `AssetRepo`, `ImageStore`, `OrderQueryPort`, `EventPublisher`
- [ ] Domain errors defined in `internal/product/errors.go`
- [ ] Domain events defined in `internal/product/events.go`
- [ ] `Money` stores amounts in configured minor units; single currency (per ADR 1)
- [ ] `ProductAsset.uri` holds relative path or URL — no `file:///` prefix (per ADR 1)
- [ ] `ProductSvc.PublishProduct` returns `ErrInvalidTransition` when product has no active variants
- [ ] `ProductSvc.PublishProduct` returns `ErrInvalidTransition` when called on an already-Active or Archived product
- [ ] `ProductSvc.AdjustStock` returns `ErrInsufficientStock` when delta would push qty below zero
- [ ] `ProductSvc.AdjustStock` publishes `StockAdjusted` event on success
- [ ] `ProductSvc.RemoveVariant` returns `ErrVariantHasActiveOrders` when `OrderQueryPort` returns true
- [ ] `ProductSvc.AddVariant` returns `ErrDuplicateSKU` when SKU already exists on the product
- [ ] `ProductSvc.UpdateVariant` publishes `VariantPriceChanged` when price differs from stored value
- [ ] `ProductSvc.ReorderAssets` persists the supplied ID ordering
- [ ] `ProductAssetSvc.UploadImage` / `UploadVideo` store files via `ImageStore` and persist assets via `AssetRepo`
- [ ] `CategorySvc.DeleteCategory` returns `ErrCategoryInUse` when products are assigned
- [ ] `SlugGenerator` produces URL-safe slugs; appends suffix on collision
- [ ] All unit tests pass: `go test ./internal/product/...`

---

### Affected Context

- **Entities / services touched**: `Product`, `ProductVariant`, `Category`, `ProductAsset`, `ProductSvc`, `ProductAssetSvc`, `CategorySvc`, `ProductQuerySvc`, `SlugGenerator`, `StockAdjuster`
- **Cross-context dependencies**: `OrderQueryPort` (read-only check into Order BC — interface only; no concrete dependency at this spec stage); `shared.Pagination` from shared module
