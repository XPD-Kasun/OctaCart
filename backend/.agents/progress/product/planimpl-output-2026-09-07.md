# Product BC — Plan Implementation Output

> Date: 2026-09-07 | BC: product

## Summary

The implementation plan covers 10 phases following the spec order: shared types → value types → entities → errors & ports → domain services → application services (ProductSvc, ProductMediaSvc, CategorySvc, ProductQuerySvc) → final verification.

**Total files to create/modify**: 17 Go files + 1 shared file modification
**Total stop points for review**: 22

## Spec Deviations Recorded

| # | Item | Deviation | Rationale |
|---|------|-----------|-----------|
| 1 | `SlugGenerator.Generate` | Uses `checkSlug` callback instead of `ProductRepo` | Keeps domain service decoupled from full repo interface |
| 2 | `Category` constructor | Path computed post-save, not at construction | ID is DB-assigned; path needs the ID |
| 3 | `ProductFilter` struct | Not in original domain-model spec | Required parameter type for `ProductRepo.List`; exact filter fields are a concern for the repo |
| 4 | `ErrMediaNotFound` | Not in original spec error list | Needed for `RemoveImage` and media upload error paths |

None of these deviations violate the spec's invariants or cross BC boundaries.

## Unresolved Risks

### Risk 1: Category Path Computation Timing (Medium)
**Description**: The spec says categories use materialized paths with the category ID, but the ID is DB-assigned (opaque `string`). The constructor cannot compute the path because the ID doesn't exist yet. The plan instructs the service to save first, then compute the path and re-save.
**Impact**: Two saves per category creation. The repo adapter will need to support this pattern (save returns/sets ID).
**Mitigation**: The `CategoryRepo.Save` implementation must assign the `id` on the category object before returning. Document this contract in the port interface documentation.

### Risk 2: Slug Collision Retry Bound (Low)
**Description**: The `SlugGenerator` appends `-2`, `-3`, etc. on collision. In theory, this could loop indefinitely if all suffixed slugs are taken.
**Impact**: Extremely unlikely in a small-scale e-commerce catalog.
**Mitigation**: Add a max-retry cap (e.g., 100 attempts) and return an error if exhausted. The plan mentions this pattern but doesn't enforce a specific cap; the child should add one.

### Risk 3: UpdateVariantCmd Double Pointer Pattern (Low)
**Description**: The `CompareAtPrice` and `MediaId` fields use `**shared.Money` / `**MediaId` (double pointer) to distinguish "don't touch" from "set to null". This is idiomatic Go but can confuse the child model.
**Impact**: Subtle implementation bugs if the child model doesn't understand the pattern.
**Mitigation**: Plan includes detailed explanation of the semantics in Phase 3.1.

### Risk 4: Event Publishing Failure Handling (Medium)
**Description**: The `EventPublisher.Publish` can return an error. The plan instructs the domain service and app service to propagate this error, but does not specify whether the stock adjustment or product state change should be rolled back.
**Impact**: If event publishing fails after a successful domain state mutation, the system is in an inconsistent state (stock changed but no event published).
**Mitigation**: This is an infrastructure concern for the adapter layer. For in-process event bus, failures are unlikely. For future message broker adoption, this becomes a transactional outbox concern. Document this as a future consideration. No action needed in this plan.

### Risk 5: Mock Complexity in Service Tests (Low)
**Description**: Phase 6 requires 6 mock implementations (one per port). This is a significant amount of test infrastructure.
**Impact**: The child model might produce incomplete or incorrect mocks, causing test failures unrelated to the actual service logic.
**Mitigation**: Mocks are in-memory maps — simple to implement. The plan provides explicit guidance on mock structure.

## No Architectural Decisions Made

No new ADRs were generated during this planning phase. All decisions were already resolved in ADR-001 through ADR-009 and the discovered BC spec.

---

**Awaiting architect review. Do not begin implementation until review is complete.**
