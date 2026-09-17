# planimpl Output — Customer BC

> Date: 2026-09-18 | BC: customer | Plan: `.agents/progress/customer/implementation-plan.md`

---

## Summary

All 5 ADRs for the Customer BC were reviewed and resolved. No `[Needs Human]` flags remain open in the Customer-specific ambiguous questions. The implementation plan is committed in 7 task files delegated from the master plan.

---

## Architectural Decisions Made During Planning

### ADR-001 Applied — Email lives in Auth BC
`CustomerProfile` has NO email field. The `CustomerRepo.GetEmailByUserId` method is defined on the driven port and implemented by the adapter by joining against the Auth user table. This keeps the Customer BC from owning or duplicating Auth's credential data.

**Impact**: Driving adapters (REST handlers) that need to display email will call `GetEmailByUserId` via the repo adapter and enrich the response DTO. The domain layer never touches email.

### ADR-003 Applied — Three-way delete logic
`DeleteAddress` in `CustomerSvc` applies conditional logic:
1. Active orders → block with `ErrCannotDeleteOnlyAddress`
2. Historical orders only → `SoftDelete()` + save
3. No orders → hard delete via `AddressRepo.Delete`

The `OrderQueryPort` interface defines two separate query methods (`HasActiveOrdersForAddress`, `HasHistoricalOrdersForAddress`) to support this three-way branching without a single ambiguous flag.

### ADR-004 Applied — Guest profiles
`NewGuestCustomerProfile(shopId)` constructor added. Sets `isGuest=true`. Driving adapters create guest profiles when merchant has guest checkout enabled. Whether merchant has it enabled is read from Settings BC (out of scope for this plan — the driving adapter checks the setting before calling this service).

### Deviation from domain-model.md — `createdBy` on CustomerProfile
The spec lists `createdBy AdminId` on `CustomerProfile`. For shopper self-registration (the primary flow), there is no admin. The constructor accepts `createdBy shared.AdminId` and the calling service passes zero value when the shopper registers themselves. This is noted here as a minor deviation from strict semantics but is the simplest approach that avoids over-engineering an optional field.

### Cross-BC boundary — `variantId` as plain string
`WishlistItem.variantId` is typed as `string`, not a Product BC type. This is intentional to avoid importing the Product package and creating a compile-time coupling between Customer and Product. The driving adapter is responsible for validating that the variantId exists (by calling Product BC's query service) before calling `AddToWishlist`.

---

## Unresolved Risks

| # | Risk | Severity | Notes |
|---|---|---|---|
| R1 | **Guest checkout merchant toggle not implemented** | Medium | ADR-004 says guest checkout requires merchant configuration in Settings BC. The Customer BC plan implements the `NewGuestCustomerProfile` constructor, but the enforcement of the merchant toggle lives in the driving adapter — which is out of scope for this domain layer plan. Must be addressed when implementing `CustomerStorefrontHandler`. |
| R2 | **Email retrieval performance** | Low | `CustomerRepo.GetEmailByUserId` joins against Auth's user table in the adapter. In a high-traffic storefront, this adds a join per customer fetch. Mitigations (caching, denormalization) deferred to V2. |
| R3 | **OrderQueryPort availability** | Medium | `HasActiveOrdersForAddress` and `HasHistoricalOrdersForAddress` require Order BC to expose query methods. The Order BC is not yet implemented. The driven adapter for `OrderQueryPort` will be a stub returning `false, nil` until Order BC is ready. This means address deletion will always hard-delete until the real adapter is wired. Document this stub assumption in the Order BC's todo. |
| R4 | **Event bus not wired** | Medium | `EventPublisher.Publish` is an in-process interface. The concrete `InProcessEventBus` is a shared adapter — not implemented in this plan. Until wired, `CustomerRegistered`, `AccountDeactivated`, etc. events will not reach Reporting or Notifications BCs. Domain layer is correct; adapter wiring is a separate concern. |
| R5 | **shopId not validated** | Low | Per spec, the Customer BC does not own or validate `ShopId`. The value is trusted as it comes from the JWT (Auth's responsibility). If a malformed shopId is passed by a misbehaving client, data corruption could occur at the DB level (wrong store scoping). This is acceptable per the multi-store model ADR. |

---

## Plan Files Generated

- [`implementation-plan.md`](./implementation-plan.md) — Master plan
- [`task-01-value-types.md`](./task-01-value-types.md) — Value types
- [`task-02-entities.md`](./task-02-entities.md) — Entities
- [`task-03-errors-ports.md`](./task-03-errors-ports.md) — Errors & ports
- [`task-04-domain-services.md`](./task-04-domain-services.md) — Domain services & events
- [`task-05-customer-svc.md`](./task-05-customer-svc.md) — CustomerSvc
- [`task-06-query-svc.md`](./task-06-query-svc.md) — CustomerQuerySvc
- [`task-07-wishlist-svc.md`](./task-07-wishlist-svc.md) — WishlistSvc

---

## Awaiting Architect Review

Please review the implementation plan and this output summary. Confirm or correct before a child agent begins implementation.
