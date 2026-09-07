## Spec: Product Catalog

- **Bounded Context**: product
- **Goal**: Manages the store's product catalog — creating, publishing, and maintaining products, their variants, pricing, and stock levels so that shoppers can discover and purchase them.

---

### In Scope
- Create, update, publish, and archive products
- Manage product variants (SKU, pricing, dynamic attributes)
- Upload and manage product media (images and videos — local storage initially, multi-tenancy ready)
- Adjust and query stock levels for variants (via domain events)
- Organize products into hierarchical categories (materialized path model)
- Tag-based product grouping and discovery
- Expose product listings and search for storefront (GraphQL)
- Expose product management APIs for admin dashboard (REST)

### Out of Scope
- Multi-currency support (single configured currency only)
- Seller/tenant scoping (single-seller platform)
- External carrier or shipping rate logic
- Order lifecycle management
- Customer management

---

### Entities

#### Product *(Aggregate Root)*

The central catalog item. A product must have at least one active variant before it can be published. All fields are unexported; driving adapters read values via public getters.

| Property      | Description                                               | DB Type          |
|---------------|-----------------------------------------------------------|------------------|
| `id`          | Opaque identifier (DB-generated)                          | `string`         |
| `title`       | Display name of the product                               | `varchar`        |
| `slug`        | URL-safe identifier; auto-generated from title            | `varchar UNIQUE` |
| `description` | Long-form text description                                | `text`           |
| `status`      | Lifecycle state: Draft/Active/Archived                    | `varchar`        |
| `categoryId`  | FK to owning Category (leaf node only)                    | `string`         |
| `tags`        | Flexible non-hierarchical labels for discovery/filtering  | `text[]`/JSON    |
| `isDigital`   | Whether this is a downloadable/digital product            | `bool`           |
| `createdAt`   | Creation timestamp                                        | `timestamptz`    |
| `updatedAt`   | Last-update timestamp                                     | `timestamptz`    |

| Method (public getter)     | Description              |
|----------------------------|--------------------------|
| `Id() string`            | Returns product ID       |
| `Title() string`         | Returns title            |
| `Slug() string`          | Returns slug             |
| `Status() ProductStatus` | Returns lifecycle status |
| `Tags() []string`        | Returns tags             |
| `IsDigital() bool`       | Returns digital flag     |

> **Note (ADR-008)**: All struct fields are unexported. Driving adapters receive a copy of the entity by value. JSON serialization uses getter methods.

---

#### ProductVariant

A purchasable SKU within a product (e.g., Size=L, Color=Red). Stored in a **separate** collection/table from `products`.

| Property         | Description                                                 | DB Type       |
|------------------|-------------------------------------------------------------|---------------|
| `id`             | Opaque identifier (DB-generated)                            | `string`      |
| `productId`      | FK to owning Product                                        | `string`      |
| `sku`            | Admin-assigned stock-keeping unit code (unique per product) | `varchar`     |
| `price`          | Current selling price (Money value object)                  | `int8`        |
| `compareAtPrice` | Optional original/crossed-out price (Money)                 | `int8 NULL`   |
| `stockQty`       | Current on-hand stock count; never goes below zero          | `int`         |
| `attrs`          | Dynamic variant attributes (ProAttributes)                  | `jsonb`       |
| `mediaId`        | Optional linked media asset ID                              | `string NULL` |

> **Note**: SKU uniqueness is enforced per product, not globally. Removing a variant that has active orders is **blocked** with `ErrVariantHasActiveOrders`.

---

#### Category

A hierarchical grouping of products using a **materialized path** strategy. Only **leaf nodes** may have products assigned. Uses dot-separated path strings (e.g., `1`, `1.2`, `1.3`) to support future PostgreSQL `ltree` indexing.

| Property | Description                                            | DB Type   |
|----------|--------------------------------------------------------|-----------|
| `id`   | Opaque identifier (DB-generated)                       | `string` |
| `name` | Display name                                           | `varchar` |
| `slug` | URL-friendly name                                      | `varchar` |
| `path` | Dot-separated materialized path (e.g., `"1.2.3"`)   | `varchar` |
| `attrs` | Attribute specifications (ProAttributes) defining required product attrs | `jsonb` |

> **Note (ADR-003)**: The `CategoryRepo` returns flat `[]*Category` lists. The `CatTree(cats []*Cat) *CatNode` domain helper constructs a tree for UI/breadcrumb use cases.

---

#### ProductMedia

Ordered gallery entry for images or videos associated with a product.

| Property    | Description                                                  | DB Type   |
|-------------|--------------------------------------------------------------|-----------|
| `id`        | Opaque identifier (DB-generated)                             | `string` |
| `productId` | FK to owning Product                                         | `string` |
| `uri`       | Relative path or direct URL (no `file:///` prefix — ADR-002) | `varchar` |
| `altText`   | Accessibility/SEO alt text for image media                   | `varchar` |
| `order`     | Integer ordering position in the gallery                     | `int`    |

---

### Value Types

#### ProId
`string` — Opaque product identifier; DB-generated.

#### ProVariantId
`string` — Opaque variant identifier; DB-generated.

#### CatId
`string` — Opaque category identifier; DB-generated.

#### MediaId
`string` — Opaque media identifier; DB-generated.

#### Money
*(ADR-004)* Represents a monetary amount in the **configured least minor units** (e.g., cents/paise). Single-currency per platform configuration.

| Property | Description                                  | Type    |
|----------|----------------------------------------------|---------|
| `amount` | Amount in minor currency units (e.g., cents) | `int64` |

> **Note**: Immutable. No floating-point representation. Currency is configured per tenant in Settings BC.

#### Attribute
A single named key-value pair for a variant or product (e.g., `{name:"Color", value:"Red"}`).

| Property | Type     |
|----------|----------|
| `name` | `string` |
| `value` | `string` |

#### AttrType
*(ADR-001)* A discriminated string representing the datatype of a dynamic attribute. One of:
- `Num` — numeric value
- `Str` — free-form string
- `Enum:XL|M|SX` — enumeration of allowed values (pipe-separated)
- `Range:1|10` — min/max integer range

Requires a **parser method** to decode from persisted string form.

#### ProAttributes
*(ADR-001)* A typed map of attribute name to AttrType. Used on `Category` to specify attribute schema, and on `ProductVariant` for variant-level attributes. Serialized as JSONB.

| Method                          | Description                                   |
|---------------------------------|-----------------------------------------------|
| `New() ProAttributes`         | Creates an empty ProAttributes instance       |
| `Add(attr string, dt AttrType)` | Adds or updates an attribute spec           |
| `Attrs() []string`            | Returns all attribute names (keys of the map) |
| `Remove(attr string)`         | Removes an attribute specification            |

#### ProductStatus
Enum: `Draft | Active | Archived`

| Value      | Meaning                                                      |
|------------|--------------------------------------------------------------|
| `Draft`  | Not yet visible to shoppers; default on creation             |
| `Active` | Published; visible on storefront; requires >=1 active variant |
| `Archived` | Retired; hidden from shoppers but retained for history     |

> **Note**: Valid transitions: `Draft -> Active` (publish), `Active -> Archived` (archive).

---

### Application Services

#### ProductSvc

Purpose: Primary application service for catalog management. Handles product lifecycle, variant management, stock adjustments, and image ordering. Called by `ProductRestHandler` and GraphQL resolvers.

| Method | Summary |
|--------|---------|
| `CreateProduct(ctx, input ProductCreateInput) (Product, error)` | Creates a product in Draft status; auto-generates slug |
| `UpdateProduct(ctx, productId ProId, cmd UpdateProductCmd) (Product, error)` | Updates mutable fields; re-generates slug if title changes |
| `PublishProduct(ctx, productId ProId) error` | Transitions Draft -> Active; blocked if no active variants |
| `ArchiveProduct(ctx, productId ProId) error` | Transitions Active -> Archived |
| `AddVariant(ctx, productId ProId, input AddVariantInput) (ProductVariant, error)` | Appends variant; enforces SKU uniqueness per product |
| `UpdateVariant(ctx, variantId ProVariantId, cmd UpdateVariantCmd) error` | Updates mutable variant fields |
| `RemoveVariant(ctx, variantId ProVariantId) error` | Deletes variant; blocked if HasActiveOrdersForVariant = true |
| `AdjustStock(ctx, variantId ProVariantId, delta int) error` | Adjusts stockQty; publishes StockAdjusted event; blocked if result < 0 |
| `ReorderImages(ctx, productId ProId, orderedIds []MediaId) error` | Resets order field across a product's media gallery |
| `RemoveImage(ctx, imageId MediaId) error` | Deletes image from MediaRepo and ImageStore |

**Errors**: `ErrProductNotFound`, `ErrVariantNotFound`, `ErrDuplicateSKU`, `ErrInvalidTransition`, `ErrInsufficientStock`, `ErrCategoryNotFound`, `ErrSlugConflict`, `ErrVariantHasActiveOrders`

> **Notes (ADR-005)**: Use xcommand pattern — dedicated input structs for >=4 params, bare parameter lists for <4 params.

---

#### ProductMediaSvc

Purpose: Dedicated application service for asset management (images/videos). Keeps file/IO concerns separate from catalog domain rules (ADR-006).

| Method | Summary |
|--------|---------|
| `UploadImage(ctx, productId ProId, fileStream io.Reader, altText string) (ProductMedia, error)` | Saves image via ImageStore; persists metadata in MediaRepo |
| `UploadVideo(ctx, productId ProId, fileStream io.Reader) (ProductMedia, error)` | Saves video via ImageStore; persists metadata in MediaRepo |

---

#### CategorySvc

Purpose: Manages hierarchical category tree operations. Enforces materialized path invariants on create/delete/move.

| Method | Summary |
|--------|---------|
| `CreateCat(ctx, catName string, parentCategory CatId, subCategories []CatId) (Category, error)` | Creates a category node; computes materialized path |
| `UpdateCat(ctx, categoryId CatId, input UpdateCatInput) error` | Updates name/slug/attrs; recomputes descendant paths if slug changes |
| `DeleteCat(ctx, categoryId CatId) error` | Blocked if products are assigned (ErrCategoryInUse) |
| `GetCatById(ctx, catId CatId) (Category, error)` | Fetch by ID |
| `GetCatBySlug(ctx, slug string) (Category, error)` | Fetch by slug |
| `GetDescendents(ctx, cat Category) ([]Category, error)` | Returns all descendant categories (flat list) |

**Errors**: `ErrCategoryNotFound`, `ErrCategoryInUse`

---

#### ProductQuerySvc

Purpose: Read-side service for REST (admin) and GraphQL (storefront) queries.

| Method | Summary |
|--------|---------|
| `GetProduct(ctx, productId ProId) (*Product, error)` | Fetch single product by ID |
| `GetProducts(ctx, filter ProductFilter, page shared.Pagination) ([]*Product, error)` | Paginated product listing with filter |
| `GetVariant(ctx, variantId ProVariantId) (*ProductVariant, error)` | Fetch single variant by ID |
| `GetStockLevel(ctx, variantId ProVariantId) (int, error)` | Returns current stockQty for a variant |

> **Notes (ADR-007)**: `shared.Pagination` is defined in `internal/shared` and reused across all BCs and adapters.

---

### Domain Services

#### SlugGenerator

Derives a URL-safe slug from a product/category title; guarantees uniqueness by appending a short suffix on collision.

Participating Entities: Product, Category

| Method | Description |
|--------|-------------|
| `Generate(ctx, title string, repo ProductRepo) (string, error)` | Generates a unique slug from title |

---

#### StockAdjuster

Enforces the invariant that stockQty never goes below zero before publishing StockAdjusted event.

Participating Entities: ProductVariant

| Method | Description |
|--------|-------------|
| `Adjust(ctx, variant *ProductVariant, delta int, pub EventPublisher) error` | Validates delta, applies adjustment, publishes event |

---

### Driven Ports

#### ProductRepo

| Method | Description |
|--------|-------------|
| `Save(ctx, product *Product) error` | Insert or update a product |
| `FindById(ctx, id ProId) (*Product, error)` | Find by ID |
| `FindBySlug(ctx, slug string) (*Product, error)` | Find by slug |
| `List(ctx, filter ProductFilter, page shared.Pagination) ([]*Product, error)` | Paginated list with filter |
| `Delete(ctx, id ProId) error` | Delete a product record |

Database mapping: `products` table. Dynamic attributes as JSONB.

---

#### VariantRepo

| Method | Description |
|--------|-------------|
| `Save(ctx, variant *ProductVariant) error` | Insert or update a variant |
| `FindById(ctx, id ProVariantId) (*ProductVariant, error)` | Find by variant ID |
| `FindByProductId(ctx, productId ProId) ([]*ProductVariant, error)` | List all variants for a product |
| `FindBySKU(ctx, sku string, productId ProId) (*ProductVariant, error)` | Find by SKU within a product |
| `Delete(ctx, id ProVariantId) error` | Delete a variant record |

Database mapping: `product_variants` table. `attrs` as JSONB.

---

#### CategoryRepo

| Method | Description |
|--------|-------------|
| `Save(ctx, cat *Category) error` | Insert or update a category |
| `GetCatBySlug(ctx, slug string) (*Category, error)` | Find by slug |
| `GetCatById(ctx, id CatId) (*Category, error)` | Find by ID |
| `GetDescendCats(ctx, catId CatId) ([]*Category, error)` | Returns all descendants via path prefix query |
| `GetChildren(ctx, id CatId) ([]*Category, error)` | Returns immediate children |
| `GetAncestors(ctx, id CatId) ([]*Category, error)` | Returns all ancestors |
| `GetSubtree(ctx, id CatId) ([]*Category, error)` | Returns full subtree rooted at ID |
| `Move(ctx, id CatId, parentId CatId) error` | Moves a category and recomputes subtree paths |
| `GetDepth(ctx, depth int) ([]*Category, error)` | Returns all categories at a given depth |
| `Delete(ctx, id CatId) error` | Delete a category (only if no products assigned) |

Database mapping: `categories` table. `attrs` (ProAttributes) as JSONB.

---

#### MediaRepo

| Method | Description |
|--------|-------------|
| `Save(ctx, media *ProductMedia) error` | Insert or update a media record |
| `FindById(ctx, id MediaId) (*ProductMedia, error)` | Find by media ID |
| `FindByProId(ctx, productId ProId) ([]*ProductMedia, error)` | List all media for a product, ordered by `order` |
| `Delete(ctx, id MediaId) error` | Delete a media record |

Database mapping: `product_media` table. `uri` stored without `file:///` prefix (ADR-002).

---

#### ImageStore

Filesystem abstraction for storing and deleting asset files. Multi-tenancy ready.

| Method | Description |
|--------|-------------|
| `Store(ctx, file io.Reader, filename string) (filePath string, error)` | Saves asset to configured storage; returns relative path/URI |
| `Delete(ctx, filePath string) error` | Removes asset from storage |

---

#### OrderQueryPort

Read-only cross-BC port into the Order BC. Used only for variant-removal invariant enforcement.

| Method | Description |
|--------|-------------|
| `HasActiveOrdersForVariant(ctx, variantId ProVariantId) (bool, error)` | Returns true if any active order references this variant |

---

#### EventPublisher

In-process event publishing port.

| Method | Description |
|--------|-------------|
| `Publish(ctx, event shared.DomainEvent) error` | Publishes a domain event to subscribed handlers |

---

### Adapters

| Adapter | Side | Description |
|---------|------|-------------|
| `ProductRepoImpl` | Driven | Implements `ProductRepo` against the configured database (DB TBD) |
| `VariantRepoImpl` | Driven | Implements `VariantRepo` against the configured database |
| `CategoryRepoImpl` | Driven | Implements `CategoryRepo`; uses materialized path queries |
| `MediaRepoImpl` | Driven | Implements `MediaRepo` against the configured database |
| `LocalImageStore` | Driven | Implements `ImageStore`; saves assets to a local folder with tenant-placeholder path prefix |
| `InProcessEventBus` | Driven | Implements `EventPublisher`; synchronous in-process fan-out to registered handlers |
| `ProductRestHandler` | Driving | Gin REST handlers under `/api/v1/products` and `/api/v1/categories` for admin operations |
| `ProductGraphQLResolver` | Driving | GraphQL resolvers for storefront product queries |

---

### Inbound

Commands handled by this BC (via REST or GraphQL driving adapters):

`CreateProduct`, `UpdateProduct`, `PublishProduct`, `ArchiveProduct`, `AddVariant`, `UpdateVariant`, `RemoveVariant`, `AdjustStock`, `UploadImage`, `UploadVideo`, `ReorderImages`, `RemoveImage`, `CreateCategory`, `UpdateCategory`, `DeleteCategory`

**Events Consumed**: None at launch.

---

### Outbound

**Events Published**:

| Event | Payload | Consumers |
|-------|---------|-----------|
| `ProductPublished` | `{productId, title, slug}` | Reporting BC |
| `ProductArchived` | `{productId}` | Reporting BC |
| `StockAdjusted` | `{variantId, delta, newQty}` | Reporting BC; Order BC (async mechanism TBD) |
| `VariantPriceChanged` | `{variantId, oldPrice, newPrice}` | Order BC (invalidate active carts) |

**REST Endpoints (admin)**:

| Method | Path |
|--------|------|
| POST | /api/v1/products |
| GET | /api/v1/products |
| GET | /api/v1/products/:id |
| PATCH | /api/v1/products/:id |
| POST | /api/v1/products/:id/publish |
| POST | /api/v1/products/:id/archive |
| POST | /api/v1/products/:id/variants |
| PATCH | /api/v1/products/:id/variants/:vid |
| DELETE | /api/v1/products/:id/variants/:vid |
| POST | /api/v1/products/:id/images |
| DELETE | /api/v1/products/:id/images/:iid |
| PATCH | /api/v1/products/:id/images/reorder |
| GET | /api/v1/categories |
| POST | /api/v1/categories |
| PATCH | /api/v1/categories/:id |
| DELETE | /api/v1/categories/:id |

**GraphQL (storefront)**: `products(filter, page)`, `product(id)`, `productVariant(id)`

---

### Dependencies

| Direction | Bounded Context | Pattern | Notes |
|-----------|----------------|---------|-------|
| Upstream | Auth (Generic) | OHS / Published Language | JWT claims provide authenticated admin/user identity |
| Upstream | Settings (Generic) | OHS / Published Language | Currency, tax, and store-level config consumed by Product |
| Downstream | Order BC | Customer-Supplier | Order reads Product/Variant/Stock; consumes VariantPriceChanged & StockAdjusted events |
| Downstream | Reporting BC | Customer-Supplier | Reporting consumes ProductPublished, ProductArchived, StockAdjusted events |

---

### Ubiquitous Language

| Term | Meaning in this BC |
|------|--------------------|
| Product | A catalog item with a title, description, and one or more variants |
| Variant | A purchasable SKU — the unit that carries a price and stock count |
| SKU | Stock Keeping Unit — admin-assigned unique code for a variant within a product |
| Category | A hierarchical label grouping products; only leaf nodes may hold products |
| Materialized Path | Dot-separated string encoding of a category's full ancestry (e.g., "1.2.3") |
| Slug | URL-friendly identifier derived from title; guaranteed unique |
| Draft | A product not yet visible to shoppers (default status on creation) |
| Active | A published product visible on the storefront; requires >=1 active variant |
| Archived | A retired product; hidden from shoppers but retained for history |
| Stock Adjust | Increase or decrease a variant's quantity; never below zero |
| Money | Amount in configured minor currency units (int64); no floating-point |
| ProAttributes | Typed map of attribute name -> AttrType used on categories and variants |
| AttrType | Discriminated string encoding a dynamic attribute's datatype |
| Tag | Free-form label attached to a product for orthogonal grouping and discovery |
| Listing | A read projection of a product as shown in search/browse results |
| Media | An ordered gallery asset (image or video) associated with a product |

---

### Polysemes

- **`Price`**: In the Product context, the current display/selling price on a `ProductVariant`. In the Order context, the price **snapshot captured at time of purchase** — immutable. Orders must copy the value, never reference Product pricing directly.
- **`Stock`**: In the Product context, the current on-hand quantity (`stockQty` on a variant). In the Order context, the availability state at the time the order was placed.

---

### Acceptance Criteria

- [ ] A product can be created in Draft status with a title, description, category, and optional tags
- [ ] Slug is auto-generated from the title and guaranteed unique (suffix appended on collision)
- [ ] A product cannot be published without at least one active variant
- [ ] Publishing a Draft product transitions it to Active; archiving an Active product transitions it to Archived
- [ ] Invalid status transitions return ErrInvalidTransition
- [ ] A variant with a duplicate SKU within the same product returns ErrDuplicateSKU
- [ ] Removing a variant that has active orders returns ErrVariantHasActiveOrders
- [ ] Stock adjustment that would result in stockQty < 0 returns ErrInsufficientStock
- [ ] A successful AdjustStock call publishes a StockAdjusted domain event
- [ ] Publishing a product publishes a ProductPublished domain event
- [ ] Archiving a product publishes a ProductArchived domain event
- [ ] A variant price change publishes a VariantPriceChanged domain event
- [ ] Categories form a single-root hierarchy; only leaf categories may be assigned to products
- [ ] Category paths use dot-separated materialized paths (e.g., "1.2.3")
- [ ] Deleting a category assigned to products returns ErrCategoryInUse
- [ ] Product media URIs are stored without the file:/// scheme prefix
- [ ] Image and video upload is handled exclusively by ProductMediaSvc, not ProductSvc
- [ ] All domain entities use unexported fields with public getter methods; adapters receive entities by value
- [ ] All update operations use the xcommand pattern (dedicated structs for >=4 params, param lists for <4)
- [ ] Pagination is sourced from internal/shared and used consistently across all product query operations
- [ ] Money values are stored as int64 in configured least minor units; no floating-point representations
- [ ] ProAttributes on a category define the attribute schema inherited by associated products
- [ ] Tags can be added/removed from products and are used in storefront filtering and discovery
