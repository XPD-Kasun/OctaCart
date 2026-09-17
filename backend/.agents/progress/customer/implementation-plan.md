# Customer Bounded Context — Implementation Plan

> Generated: 2026-09-18 | BC: customer | Spec: `.agents/discover/customer/domain-model.md`

---

## Preamble: Coding Instructions for the Implementer

Before you write any code, read the following instructions in full. These are your operating rules for the entire plan.

### Step 0 — Read Context

- [ ] Read `AGENTS.md` at the project root to understand the architecture, tech stack, and folder structure.
- [ ] Read `.agents/discover/customer/domain-model.md` — this is your **authoritative spec**. Every type, method, and error you create must match it.
- [ ] Read **all ADR files** `.agents/specs/customer/ADR001.md` through `ADR005.md` — these are reviewed architectural decisions that inform implementation details. **Do NOT read** archived ADR files.
- [ ] Read `.agents/discover/conventions.md` — this defines all coding style rules.
- [ ] Read `.agents/discover/context-map.md` — this tells you the cross-BC boundaries. The Customer BC must **never** write to data owned by another BC.

### Coding Conventions Summary

1. **Module**: `octacart`
2. **File header**: Every `.go` file starts with the copyright block (see conventions.md).
3. **Package**: All customer domain code goes in `internal/customer/` under `package customer`.
4. **Entity fields**: Always **unexported** (lowercase). Provide **public getters** (PascalCase, no `Get` prefix).
5. **Constructors**: `New{Type}(...)` factory functions. Return `({Type}, error)` when validation is needed.
6. **xcommand pattern**: Use input structs for methods with ≥ 4 params; bare param lists for < 4.
7. **Errors**: Sentinel vars in `errors.go`: `var ErrXxx = errors.New("...")`.
8. **Ports**: Interfaces in `ports.go`. First param is always `context.Context`. Last return is always `error`.
9. **Tests**: Standard `testing` package only. Same package (white-box). Use `t.Run` subtests. Name pattern: `Test{Type}_{Method}`.
10. **Imports**: stdlib → blank line → project (`octacart/internal/...`) → blank line → third-party.
11. **Shared types**: `shared.Money`, `shared.UserId`, `shared.AdminId`, `shared.Claim`, `shared.ActorType`, `shared.Pagination`, `shared.DomainEvent` already exist in `internal/shared/types.go`.
12. **Logging (zerolog)**: App services receive `*zerolog.Logger` via constructor. Domain services have no logger. Log key mutations at `Info` level. Include entity IDs in log events.
13. **shopId scoping**: `ShopId` is a value type defined in the Customer BC package. It is threaded through all repo calls. The repo adapter enforces it at the database level — never skip it.
14. **Email (ADR-001)**: Do NOT store email in CustomerProfile. Email is read from Auth BC via `CustomerRepo.GetEmailByUserId` in the driven adapter. The driving adapter may call this after fetching a profile.
15. **Guest checkout (ADR-004)**: Guests are represented by a `CustomerProfile` with `isGuest=true`. Created by the driving adapter when guest checkout is toggled on by the merchant (Settings BC decision).
16. **Address deletion (ADR-003)**: Active/open order → block (return error). Historical order only → soft-delete. No order reference → hard-delete. This logic lives in `CustomerSvc.DeleteAddress`.
17. **Wishlist (ADR-005)**: Single wishlist per customer. No multi-wishlist support in V1.
18. **Cross-BC references**: `WishlistItem.variantId` is a plain `string` — do NOT import Product BC package. `AddressSnapshot` is passed to Order BC by value.
19. **Marketing consent**: Excluded from V1 per spec. Do not implement.

### Test Running Protocol

- After writing tests and implementation for each item, run: `go test ./internal/customer/...`
- If tests pass, continue to the next item.
- If tests fail: attempt a fix **within the spec's scope** and re-run, up to **2 attempts**.
- If still failing after 2 attempts, or the fix would require **deviating from the spec**: **STOP** and report the failure to the user. Do not proceed further.

### File Layout

Create these files in `internal/customer/`:

| File | Contents |
|---|---|
| `customer.go` | All value types, entity structs (CustomerProfile, Address, WishlistItem), constructors, domain methods, command structs. Package doc comment. |
| `events.go` | All domain event types implementing `shared.DomainEvent` |
| `errors.go` | All sentinel error vars |
| `ports.go` | All driven port interfaces (CustomerRepo, AddressRepo, WishlistRepo, OrderQueryPort, EventPublisher) plus CustomerFilter struct |
| `registrar_service.go` | CustomerRegistrar domain service |
| `customer_service.go` | CustomerSvc application service (profile + address methods) |
| `query_service.go` | CustomerQuerySvc application service |
| `wishlist_service.go` | WishlistSvc application service |
| `customer_test.go` | Tests for value types, entities, domain events |
| `registrar_service_test.go` | Tests for CustomerRegistrar |
| `customer_service_test.go` | Tests for CustomerSvc (includes shared mock infrastructure) |
| `query_service_test.go` | Tests for CustomerQuerySvc |
| `wishlist_service_test.go` | Tests for WishlistSvc |

---

## Task Delegation

This master plan delegates implementation to numbered sub-task files. Read each task file, complete it, then return here for the next task. **Do not read ahead** — work one task at a time.

---

### Task 01 — Value Types

📄 Read and follow: [`.agents/progress/customer/task-01-value-types.md`](../customer/task-01-value-types.md)

**Scope**: `CustomerId`, `AddressId`, `WishlistItemId`, `ShopId`, `CustomerStatus`, `AddressSnapshot`

After completing Task 01, return here and proceed to Task 02.

---

### Task 02 — Entities

📄 Read and follow: [`.agents/progress/customer/task-02-entities.md`](../customer/task-02-entities.md)

**Scope**: `CustomerProfile` (aggregate root), `Address`, `WishlistItem` — structs, constructors, domain methods, command structs

After completing Task 02, return here and proceed to Task 03.

---

### Task 03 — Errors & Ports

📄 Read and follow: [`.agents/progress/customer/task-03-errors-ports.md`](../customer/task-03-errors-ports.md)

**Scope**: `errors.go` sentinel vars, `ports.go` interfaces (CustomerRepo, AddressRepo, WishlistRepo, OrderQueryPort, EventPublisher), `CustomerFilter`

After completing Task 03, return here and proceed to Task 04.

---

### Task 04 — Domain Services

📄 Read and follow: [`.agents/progress/customer/task-04-domain-services.md`](../customer/task-04-domain-services.md)

**Scope**: `CustomerRegistrar` domain service, `events.go` (CustomerRegisteredEvent, ProfileUpdatedEvent, AccountDeactivatedEvent, WishlistItemAddedEvent)

After completing Task 04, return here and proceed to Task 05.

---

### Task 05 — CustomerSvc Application Service

📄 Read and follow: [`.agents/progress/customer/task-05-customer-svc.md`](../customer/task-05-customer-svc.md)

**Scope**: `CustomerSvc` — all profile methods (`RegisterCustomer`, `UpdateProfile`, `DeactivateAccount`, `ReactivateAccount`) and all address methods (`AddAddress`, `UpdateAddress`, `DeleteAddress`, `SetDefaultAddress`, `GetAddressSnapshot`)

After completing Task 05, return here and proceed to Task 06.

---

### Task 06 — CustomerQuerySvc Application Service

📄 Read and follow: [`.agents/progress/customer/task-06-query-svc.md`](../customer/task-06-query-svc.md)

**Scope**: `CustomerQuerySvc` — `GetCustomer`, `GetCustomerByUserId`, `ListCustomers`, `ListAddresses`

After completing Task 06, return here and proceed to Task 07.

---

### Task 07 — WishlistSvc Application Service

📄 Read and follow: [`.agents/progress/customer/task-07-wishlist-svc.md`](../customer/task-07-wishlist-svc.md)

**Scope**: `WishlistSvc` — `AddToWishlist`, `RemoveFromWishlist`, `ListWishlist`, `ClearWishlist` + final full-suite verification

After completing Task 07, the Customer BC domain layer implementation is **complete**.

---

## Completion Checklist

When all tasks are done, verify:
- [ ] `go test ./internal/customer/...` passes with no failures
- [ ] `go build ./internal/customer/...` compiles cleanly
- [ ] `go vet ./internal/customer/...` reports no issues
- [ ] All 7 task files show all checkboxes ticked

Then tell the architect: **"Customer BC implementation complete. All tasks done. Ready for architect review."**
