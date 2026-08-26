# Product ADR 1

- For `Product` how are we going to store this in db ? `attributes []Attribute`
- `ProductImage` has `uri` which can include relative paths for a configured directory, even its url, we dont use `file:///` prefix for brevity and this stage
- `Category` architecture needed to plan
- `Money` is always single currency in this version. That currency can be configured per tenent. Stores the values in configured least minor units.
- For updates we use `func(*Entity)` to set the entity values instead of AI suggested input. Fields as args makes the signature lengthier. Lambda simply allows to set public fields. But field by field granularity is missed. If this is needed we introduce a separate struct only containing those. We'll evaluate this further. This is a syntatic sugar and makes it easier for the caller not a technical requirement.
- Image uploading is not a responsibility of `ProductSVC`. Therefore, we use a driving port `ImageRepo` with a application service `ProducAssetSvc`. We use asset to include images, videos etc. Currently images and videos are supported.
- `Pagination` is a core need for any application service and adapters, therefore, we need to consider this in shared module and belongs to shared module.
