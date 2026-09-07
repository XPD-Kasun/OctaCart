# OctaCart Coding Conventions

> Extracted from the auth bounded context reference implementation and ADRs.

## General Go Style

- Module name: `octacart`
- Go version: 1.26.1
- Use `gofmt` / `goimports` for formatting (default Go formatting rules)
- Prefer standard library over third-party when possible

## File Header

Every `.go` file must start with:

```go
// Copyright 2026 OctaCart. All rights reserved.
// Author: <author>
// Created: <date YYYY-MM-DD>
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.
```

## Package Documentation

The main entity file in each BC package should contain a package-level doc comment explaining the bounded context, its key types, and where to find ports/errors.

## Package Structure (per BC)

Within `internal/{bc}/`:

| File | Purpose |
|------|---------|
| `{entity}.go` | Entity structs, value objects, constructors, domain methods |
| `ports.go` | All driven port interfaces for this BC |
| `errors.go` | All sentinel errors (`var ErrX = errors.New("...")`) |
| `{service_name}_service.go` | Each application service in its own file |
| `{entity}_test.go` | Unit tests for entity and value object logic |
| `{service_name}_service_test.go` | Unit tests for service logic |

## Entity Encapsulation (ADR-008)

- All entity struct fields are **unexported** (lowercase)
- Expose **public getter methods** for each field needed by driving adapters
- Return entities **by value** from application services to driving adapters (creates an immutable copy)
- Getters use the field name in PascalCase without `Get` prefix (e.g., `Id()`, `Title()`, not `GetId()`)
- Example: `func (p *Product) Title() string { return p.title }`

## Constructor Pattern

- Use `New{Type}(...)` factory functions for entities and services
- Constructors validate required fields and return `({Type}, error)` when validation is needed
- For application services: `func New{Svc}(deps...) *{Svc}`

## Update Commands (ADR-005 — xcommand pattern)

- If a method has **fewer than 4 parameters**: use a bare parameter list
- If a method has **4 or more parameters**: use a dedicated input/command struct
- Naming: `{Entity}CreateInput`, `Update{Entity}Cmd`, `Add{Sub}Input`
- Command structs live in the same package as the entity they modify

## Error Conventions

- Sentinel errors as package-level vars: `var ErrXxx = errors.New("descriptive message")`
- Use `fmt.Errorf("context: %w", err)` for wrapping
- Callers check with `errors.Is(err, pkg.ErrXxx)`
- Error variable names: `Err` prefix + PascalCase description (e.g., `ErrProductNotFound`)

## Port Conventions

- All driven ports are **interfaces** defined in `ports.go` inside the BC package
- Method signatures use `context.Context` as the first parameter
- Return `error` as the last return value
- Port names: `{Entity}Repo`, `{Concept}Store`, `{Concept}Port`, `EventPublisher`

## Service Conventions

- Application services are **structs** holding port interfaces as dependencies (dependency injection via constructor)
- Services use `*zerolog.Logger` for logging (may change to abstraction later)
- Service methods take `context.Context` as the first parameter

## Testing Conventions

- Use Go standard `testing` package — **no third-party test frameworks**
- Tests live in the **same package** as the code under test (white-box testing)
- Use `t.Run("description", func(t *testing.T) {...})` for subtests
- Test function names: `Test{Type}_{Method}` (e.g., `TestProduct_Publish`)
- Subtest names: `"condition/expected outcome"` with slashes (e.g., `"draft product/should transition to active"`)
- Assertions use `t.Error` / `t.Errorf` — no `t.Fatal` unless the test cannot continue
- Interface compliance checks: `var _ InterfaceName = &ConcreteType{}`
- Run tests with: `go test ./internal/{bc}/...`

## Shared Types (`internal/shared`)

Types already in `internal/shared/types.go`:
- `Money int64` — monetary amount in minor units
- `KV[Key comparable, Val any]` — generic key-value pair
- `UserId int` — user identifier
- `Claim string` — auth claim type

Types that need to be added (as part of implementation):
- `Pagination` — reusable pagination struct (ADR-007)
- `DomainEvent` — base event type for event publishing

## Logging

- Use `zerolog` (`github.com/rs/zerolog`)
- Follow structured logging: `.Str("key", val).Msg("message")` style

## Value Types

- Opaque ID types are `type XxxId string` aliases
- `Money` is `int64` (already in shared)
- Enums use typed string constants with `const` block

## Imports

- Group imports: stdlib, then blank line, then project imports, then blank line, then third-party
- Example:

```go
import (
    "context"
    "errors"

    "octacart/internal/shared"

    "github.com/rs/zerolog"
)
```
