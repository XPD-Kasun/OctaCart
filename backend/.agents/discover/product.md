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

Add public getters for the below for Dto boundaries appservice -> handler.

- **Product** — The central aggregate root.
  - `id` (ProductId), `title`, `slug`, `description`, `status` (Draft | Active | Archived), `categoryId`, `tags []string`, `createdAt`, `updatedAt, isDigital`
- **ProductVariant** — A purchasable SKU within a product (e.g., Size=L, Color=Red). Stored in a separate collection.
  - `id` (VariantId), `productId`, `sku`, `price` (Money), `compareAtPrice` (Money?), `stockQty`, `attrs ProAttributes`, `mediaId? `
- **Category** — A hierarchical grouping of products.  We use single heirarchy with single parent. only the leaves can have products. We use materialized paths. 
  - `id` (CategoryId), `name`, `slug`, `path`, `attrs ProAttributes`
  Lets say following shoes(id:1) has men(id:2) and women(id:3). then paths are `1`, `1.2`, `1.3`. We use dots to leverage any future use of postgres ltree support.

- **ProductMedia** — Ordered gallery media for a product. Supports for images and videos currently. 
  - `id` (MediaId), `productId`, `uri`, `altText`, `order`

---

### Value Types
- **shared.KV** - `struct{k:string,v:string}` (Add this to shared package)
- **ProId** — `string` (opaque; DB-generated)
- **ProVariantId** — `string` (opaque; DB-generated)
- **CatId** — `string` (opaque; DB-generated)
- **MediaId** — `string` (opaque; DB-generated)
- **Money** — `amount int64` — amount in minor units (cents/paise); immutable; single-currency for now
- **Attribute** — `{name string, value string}` (e.g., `{name:"Color", value:"Red"}`)
- **ProductStatus** — enum: `Draft | Active | Archived`
- **ProAttributes** - map[string]AttrType - A tyoe collection of key, value for product attr name and datatype(see below).
  This should have methods: Add(attr, datatype string), Attrs() : keys of the map, Remove(attr), New(...shared.KV)
- **AttrType** - enum: `Num | Str | Enum (eg: Enum:XL|M|SX) | Range (eg: Range:1|10)` 

Note: datatype for `ProAttributes` is defined by us and is an enum. These are used as a kind of specification that is adhered by products attached to this category. Useful for filtering.
---

### Application Services

- **ProductSvc** (primary application service)
  - `CreateProduct(ctx, input)` → `(Product, error)` — creates a Draft product
  - `UpdateProduct(ctx, productId, UpdateProductCmd)` → `(Product, error)` — updates mutable fields; slug re-gen on title change
  - `PublishProduct(ctx, productId)` → `error` — Draft → Active; requires >=1 active variant
  - `ArchiveProduct(ctx, productId)` → `error` — Active → Archived
  - `AddVariant(ctx, productId, input)` → `(ProductVariant, error)` — appends a variant with a unique SKU per product
  - `UpdateVariant(ctx, variantId, UpdateVariantCmd)` → `error`
  - `RemoveVariant(ctx, variantId)` → `error` — blocked if variant has active orders (checks via OrderRepo port)
  - `AdjustStock(ctx, variantId, deltaAmount int)` → `error` — publishes StockAdjusted domain event; delta negative (sale) or positive (restock)
  - `ReorderImages(ctx, productId, orderedIds)` → `error`
  - `RemoveImage(ctx, imageId)` → `error`

  - **Errors**: `ErrProductNotFound`, `ErrVariantNotFound`, `ErrDuplicateSKU`, `ErrInvalidTransition`, `ErrInsufficientStock`, `ErrCategoryNotFound`, `ErrSlugConflict`, `ErrVariantHasActiveOrders`

- **ProductMediaSvc**
  - `UploadImage(ctx, productId, fileStream io.Reader, altText)` → `(ProductMedia, error)` — saves to local folder; path stored in MediaRepo  
  - `UploadVideo(ctx, productId, fileStream io.Reader)` → `(ProductMedia, error)` — saves to local folder; path stored in MediaRepo

**CategorySvc**
  - `CreateCat(ctx, catName, parentCategory, subCategories)` → `(Category, error)`
  - `UpdateCat(ctx, categoryId, input)` → `error`
  - `DeleteCat(ctx, categoryId)` → `error` — blocked if products assigned
  - `GetCatById(ctx, catId)` → `(Category, error)`
  - `GetCatBySlug(ctx, slug)` → `(Category, error)`
  - `GetDescendents(ctx, Category)` -> `([]Category, error)`
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
| `CategoryRepo`   | `Save`, `GetCatBySlug(ctx, slug)`, `GetCatById(ctx, id)`, `GetDescendCats(ctx, catId)`      |
| `MediaRepo`      | `Save`, `FindById`, `FindByProId`, `Delete`                                             |
| `ImageStore`     | `Store(file) → filePath`, `Delete(filePath)` — filesystem abstraction (multi-tenancy ready) |
| `OrderQueryPort` | `HasActiveOrdersForVariant(ctx, variantId) → bool` — read-only check into Order BC          |
| `EventPublisher` | `Publish(ctx, event shared.DomainEvent)` — in-process initially                                    |
In above CategoryRepo use ProAttributes as a jsonb.
---

### Adapters (to be implemented later - not implemented along with above currently)

| Side    | Adapter                  | Purpose                                                                           |
| ------- | ------------------------ | --------------------------------------------------------------------------------- |
| Driven  | `ProductRepoImpl`        | Implements `ProductRepo` (DB TBD)                                                 |
| Driven  | `VariantRepoImpl`        | Implements `VariantRepo` (DB TBD)                                                 |
| Driven  | `CategoryRepoImpl`       | Implements `CategoryRepo` (DB TBD)                                                |
| Driven  | `MediaRepoImpl`          | Implements `MediaRepo` (DB TBD)                                                   |
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

- **Upstream**:
  - Auth (Generic) — JWT claims provide authenticated user identity. Pattern: **OHS/PL**.
  - Settings (Generic) — Currency, tax, and store-level config consumed by Product. Pattern: **OHS/PL**.
- **Downstream**:
  - Order BC — reads Product/Variant/Stock; consumes VariantPriceChanged & StockAdjusted events. Product is upstream. Pattern: **Customer-Supplier**.
  - Reporting BC — consumes Product domain events. Pattern: **Customer-Supplier**.
  - Shipping — no direct dependency at catalog level.

---

### Data Owned

- `products` collection/table
- `product_variants` collection/table (separate from products)
- `categories` collection/table
- `product_media` collection/table + local asset files

---

### All Decisions Resolved

All architectural questions from discovery are resolved. No open [Needs Human] items remain.
