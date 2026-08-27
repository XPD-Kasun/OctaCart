# Product ADR 1

todo: Expand this to series of real ADR formats. we'll go like this until we fully resolve issues with discovered bc review.

- For `Product` how are we going to store this in db ? `attributes []Attribute`
- `ProductMedia` has `uri` which can include relative paths for a configured directory, even its url, we dont use `file:///` prefix for brevity and this stage
- `Category` uses single root heirarchy tree with path persistance at the db level. Finalized. Considered alternatives:
Recursive CTE, separete parent table, nested sets.

  Lets say following shoes(id:1) has men(id:2) and women(id:3). then paths are `1`, `1.2`, `1.3`. We use dots to leverage any future use of postgres ltree support.

  Also category can specify attribute spec (attr name, attr type, required) which becomes a part of products that is attached to this category. Also, a product can add more attrs. This is mainly used in filtering.

- `Money` is always single currency in this version. That currency can be configured per tenent. Stores the values in configured least minor units.
- For updates we use `func(*Entity)` to set the entity values instead of AI suggested input. Fields as args makes the signature lengthier. Lambda simply allows to set public fields. But field by field granularity is missed. If this is needed we introduce a separate struct only containing those. We'll evaluate this further. This is a syntatic sugar and makes it easier for the caller not a technical requirement. **This is not possible with entity encapsulation anymore. Therefore we need to use xcommand pattern with short param list whichever suit best.**

- Image uploading is not a responsibility of `ProductSvc`. Therefore, we use a driving port `MediaRepo` with a application service `ProductMediaSvc`. We use asset to include images, videos etc. Currently images and videos are supported.
- `Pagination` is a core need for any application service and adapters, therefore, we need to consider this in shared module and belongs to shared module.
- Domain entities should not be leaked into the adapters. Even at the initial draft we found OK to do it. Therefore, after careful thoughts and considering all alternatives, we use unexported fields in entity with public getters(). Then we return this entity from app service to the driving adapters by value (creating a copy). This aligns well with GO's encapuslation model and reducing boilerplate of having lots of dto stuff.