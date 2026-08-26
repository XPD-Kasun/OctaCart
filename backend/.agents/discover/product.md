## Bounded Context: Product (Catalog)

- **Purpose**: Manages the store's product catalog — creating, publishing, and maintaining products, their variants, pricing, and stock levels so that shoppers can discover and purchase them.
- **Classification**: Core

> **Architecture notes (decided)**
> - Storage is fully deferred behind repository ports — no specific DB assumed.
> - Single-seller/single-admin platform (no sellerId scoping needed).
> - Variants stored in a separate collection (via a dedicated port).
> - Stock adjustment via async domain event (`StockAdjusted`); consumption mechanism TBD.
> - Image upload to local folder initially; designed for multi-tenancy when it arrives.
> - Single-currency platform.
> - Removing a variant that has active orders is **blocked** with an error.

---

### Entities

- **Product** — The central aggregate root.
  - `id` (ProductId), `title`, `slug`, `description`, `status` (Draft | Active | Archived), `categoryId`, `tags []string`, `createdAt`, `updatedAt`
- **ProductVariant** — A purchasable SKU within a product (e.g., Size=L, Color=Red). Stored in a separate collection.
  - `id` (VariantId), `productId`, `sku`, `price` (Money), `compareAtPrice` (Money?), `stockQty`, `attributes []Attribute`, `assetId? `
- **Category** — A hierarchical grouping of products.  We use single heirarchy with single parent. only the leaves can have products.
  - `id` (CategoryId), `name`, `slug`, `parentId?`, `description` 
  - Need more info on methods and repr.
- **ProductAsset** — Ordered gallery asset for a product. Supports for images and videos currently. 
  - `id` (AssetId), `productId`, `uri`, `altText`, `order`

---

### Value Types

- **ProductId** — `string` (opaque; DB-generated)
- **VariantId** — `string` (opaque; DB-generated)
- **CategoryId** — `string` (opaque; DB-generated)
- **ImageId** — `string` (opaque; DB-generated)
- **Money** — `amount int64` — amount in minor units (cents/paise); immutable; single-currency for now
- **Attribute** — `{name string, value string}` (e.g., `{name:"Color", value:"Red"}`)
- **ProductStatus** — enum: `Draft | Active | Archived`

---

### Application Services

- **ProductSvc** (primary application service)
  - `CreateProduct(ctx, input)` → `(*Product, error)` — creates a Draft product
  - `UpdateProduct(ctx, productId, fieldUpdater func(*Product))` → `(*Product, error)` — updates mutable fields; slug re-gen on title change
  - `PublishProduct(ctx, productId)` → `error` — Draft → Active; requires >=1 active variant
  - `ArchiveProduct(ctx, productId)` → `error` — Active → Archived
  - `AddVariant(ctx, productId, input)` → `(VariantId, error)` — appends a variant with a unique SKU per product
  - `UpdateVariant(ctx, variantId, func(*ProductVariant))` → `error`
  - `RemoveVariant(ctx, variantId)` → `error` — blocked if variant has active orders (checks via OrderRepo port)
  - `AdjustStock(ctx, variantId, delta int)` → `error` — publishes StockAdjusted domain event; delta negative (sale) or positive (restock)
  - `ReorderImages(ctx, productId, orderedIds)` → `error`
  - `RemoveImage(ctx, imageId)` → `error`
  - **Errors**: `ErrProductNotFound`, `ErrVariantNotFound`, `ErrDuplicateSKU`, `ErrInvalidTransition`, `ErrInsufficientStock`, `ErrCategoryNotFound`, `ErrSlugConflict`, `ErrVariantHasActiveOrders`
- **ProductAssestSvc**
  - `UploadImage(ctx, productId, file MultipartFile, altText)` → `(AssetId, error)` — saves to local folder; path stored in ImageRepo  
  `UploadVideo(ctx, productId, file MultipartFile)` → `(AssetId, error)` — saves to local folder; path stored in ImageRepo
- **CategorySvc**
  - `CreateCategory(ctx, input)` → `(CategoryId, error)`
  - `UpdateCategory(ctx, categoryId, input)` → `error`
  - `DeleteCategory(ctx, categoryId)` → `error` — blocked if products assigned
  - `ListCategories(ctx)` → `([]Category, error)`
  - **Errors**: `ErrCategoryNotFound`, `ErrCategoryInUse`
- **ProductQuerySvc** (read side — REST for admin and GraphQL for shop fronts)
  - `GetProduct(ctx, productId)` → `(*Product, error)`
  - `GetProducts(ctx, filter ProductFilter, page shared.Pagination)` → `([]*Product, error)`
  - `GetVariant(ctx, variantId)` → `(*ProductVariant, error)`
  - `GetStockLevel(ctx, variantId)` → `(int, error)`

---

### Domain Services

- **SlugGenerator** — derives a URL-safe slug from a product title; guarantees uniqueness within the catalog by appending a short suffix on collision.
- **StockAdjuster** — enforces the invariant that stockQty never goes below zero before publishing StockAdjusted event.

---

### Ports (Interfaces - storage agnostic)

| Port             | Methods                                                                                     |
| ---------------- | ------------------------------------------------------------------------------------------- |
| `ProductRepo`    | `Save`, `FindById`, `FindBySlug`, `List`, `Delete`                                          |
| `VariantRepo`    | `Save`, `FindById`, `FindByProductId`, `FindBySKU`, `Delete`                                |
| `CategoryRepo`   | `Save`, `FindById`, `List`, `Delete`                                                        |
| `AssetRepo`      | `Save`, `FindById`, `FindByProductId`, `Delete`                                             |
| `ImageStore`     | `Store(file) → filePath`, `Delete(filePath)` — filesystem abstraction (multi-tenancy ready) |
| `OrderQueryPort` | `HasActiveOrdersForVariant(ctx, variantId) → bool` — read-only check into Order BC          |
| `EventPublisher` | `Publish(ctx, event DomainEvent)` — in-process initially                                    |

---

### Adapters (to be implemented)

| Side    | Adapter                  | Purpose                                                                           |
| ------- | ------------------------ | --------------------------------------------------------------------------------- |
| Driven  | `ProductRepoImpl`        | Implements `ProductRepo` (DB TBD)                                                 |
| Driven  | `VariantRepoImpl`        | Implements `VariantRepo` (DB TBD)                                                 |
| Driven  | `CategoryRepoImpl`       | Implements `CategoryRepo` (DB TBD)                                                |
| Driven  | `ImageRepoImpl`          | Implements `ImageRepo` (DB TBD)                                                   |
| Driven  | `LocalImageStore`        | Implements `ImageStore` — saves to local folder; path includes tenant placeholder |
| Driven  | `InProcessEventBus`      | Implements `EventPublisher` — in-process fan-out                                  |
| Driving | `ProductRestHandler`     | Gin REST handlers under `/api/v1/products`                                        |
| Driving | `ProductGraphQLResolver` | GraphQL resolvers for storefront product queries                                  |

---

### Ubiquitous Language

| Term         | Definition in THIS context                                                       |
| ------------ | -------------------------------------------------------------------------------- |
| Product      | A catalog item with a title, description, and one or more variants               |
| Variant      | A purchasable SKU — the unit that has a price and stock count                    |
| SKU          | Stock Keeping Unit — admin-assigned unique code for a variant                    |
| Category     | A hierarchical label grouping products by type                                   |
| Slug         | URL-friendly identifier derived from the product title                           |
| Draft        | A product not yet visible to shoppers                                            |
| Active       | A published product visible on the storefront                                    |
| Archived     | A retired product; hidden from shoppers but retained for history                 |
| Stock Adjust | Increase or decrease a variant quantity (sale, restock, correction)              |
| Money        | A value object: amount in minor currency units + currency code (single currency) |
| Listing      | A read projection of a product as shown in search/browse results                 |

**Polysemes**

- `Price` — Product context: current display price. Order context: price snapshot at purchase time (immutable). Order must capture the value, not reference Product.
- `Stock` — Product context: current on-hand quantity. Order context: availability at time order was placed.

---

### Capabilities

- Create, update, publish, and archive products
- Manage product variants (SKU, pricing, attributes)
- Upload and manage product images (local storage initially, multi-tenancy ready)
- Adjust and query stock levels for variants (via domain events)
- Organize products into hierarchical categories
- Expose product listings and search for the storefront (GraphQL)
- Expose product management APIs for the admin dashboard (REST)

---

### Inbound

- **Commands handled**: `CreateProduct`, `UpdateProduct`, `PublishProduct`, `ArchiveProduct`, `AddVariant`, `UpdateVariant`, `RemoveVariant`, `AdjustStock`, `UploadImage`, `ReorderImages`, `RemoveImage`, `CreateCategory`, `UpdateCategory`, `DeleteCategory`
- **Events consumed**: None at launch

---

### Outbound

- **Events published**:
  - `ProductPublished {productId, title, slug}` — Reporting
  - `ProductArchived {productId}` — Reporting
  - `StockAdjusted {variantId, delta, newQty}` — Reporting; Order (consumed via async mechanism TBD)
  - `VariantPriceChanged {variantId, oldPrice, newPrice}` — Order (to invalidate active carts)
- **REST endpoints** (admin):
  - `POST /api/v1/products`, `GET /api/v1/products`, `GET /api/v1/products/:id`
  - `PATCH /api/v1/products/:id`, `POST /api/v1/products/:id/publish`, `POST /api/v1/products/:id/archive`
  - `POST /api/v1/products/:id/variants`, `PATCH /api/v1/products/:id/variants/:vid`, `DELETE /api/v1/products/:id/variants/:vid`
  - `POST /api/v1/products/:id/images`, `DELETE /api/v1/products/:id/images/:iid`, `PATCH /api/v1/products/:id/images/reorder`
  - `GET /api/v1/categories`, `POST /api/v1/categories`, `PATCH /api/v1/categories/:id`, `DELETE /api/v1/categories/:id`
- **GraphQL** (storefront): `products(filter, page)`, `product(id)`, `productVariant(id)`

---

### Dependencies

- **Upstream**: Auth (Generic) — JWT claims provide authenticated user identity. Pattern: **OHS/PL**.
- **Downstream**:
  - Order BC — reads Product/Variant/Stock; Product is upstream. Pattern: **Customer-Supplier**.
  - Reporting BC — consumes Product domain events. Pattern: **Customer-Supplier**.
  - Shipping — no direct dependency at catalog level.

---

### Data Owned

- `products` collection/table
- `product_variants` collection/table (separate from products)
- `categories` collection/table
- `product_images` collection/table + local image files

---

### All Decisions Resolved

All architectural questions from discovery are resolved. No open [Needs Human] items remain.
