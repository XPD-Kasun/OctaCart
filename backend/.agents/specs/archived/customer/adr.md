---
Reviewed on: 15/09/2026
Reviewed By: XPD
---

# Customer Adrs

- Only name is stored additionally in Customer bc. auth bc already contains email and customer bc should read it from there. currently the repo adapter will do it. In a microservice, a grpc or rest api can do this. Also at V1, we dont need customer specific contact email. Having multiple emails for a customer is cumbersome. We prefer least information from the customer.
- Register customer action should sequentially and synchronously call createCustomer after creating auth profile. Using event bus overcomplicate the process and also adds error handling overhead with a service bus.
- Should allow to delete addresses that doesnt have any active/open orders. Those assigned addresses deletion should fail with ErrCannotDeleteOnlyAddress. If an address that doesnt have active orders but has past orders should only be soft deleted. Others can be permanantly deleted.
- Merchant should configure the option to enable guest checkout. If this is true only we allow inline checkout. For this we should create special customer guest and attached with checkout inline addresses.
- V1 should only have a single Wishlist. In the future, we may think about multiple wishlists for a single customer. Also, merchant can configure No Wishlist, Single Wishlist, Multiple Wishlists by admin in the future. V1 default is single wishlist only.
