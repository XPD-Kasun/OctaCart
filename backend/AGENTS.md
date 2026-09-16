# OctaCart e-Commerce Platform

You are an expert Golang enterprise solution architect.

OctaCart is a backend for small-scale e-Commerce sellers. This app provides a REST API, GraphQL endpoint for products, and an admin dashboard. You as an agent should use this as a guide. Do not load all website references into context unless needed.
Refer .agents folder at root for further context for agents.

## Architecture

We use a modular monolith approach with hexagonal architectural pattern. Since we are not planning to grow beyond monolith scale, microservices are not considered. Functional decomposition may be done in future based on modules if needed.

For REST API and GraphQL, this Go app implements the services. For the admin dashboard, a separate Next.js app with static export is used. That application consumes the APIs published by this Go app.

## Multi-Store Model

- A **single merchant** can own and operate **multiple stores**.
- All domain data (products, orders, customers, pricing, shipping, payment, reviews, etc.) is scoped by **`shopId`**.
- **`shopId`** is a string claim embedded in the merchant's JWT after they select a store. It is extracted by middleware on the driving adapter side and threaded through all application service calls.
- **Auth BC is NOT shopId-scoped** — it manages merchant identity globally across stores.
- **Customers are per-store** — a shopper account is not shared between stores.
- **No BC owns or validates shopId** — that is Auth BC's responsibility.
- The tenancy model is **shared-database with shopId column filtering** (not per-instance silo).

## Technical Stack

- Gin for HTTP server. [Gin website](https://gin-gonic.com/en/docs/)
- Zerolog for logging. [Zerolog website](https://github.com/rs/zerolog)

## Folder Structure

bc = bounded context

- .agents
  - skills - skills for agents
  - specs - PR and SRS for the agents
  - progress - agentic work results for logging
  - todo - generated todos for tasks
  - discover - bounded context discovery docs (one folder per BC)
- cmd - entry points
- internal - main app logic
  - driving - inbound adapters (REST and GraphQL)
  - driven - outbound adapters (database and external systems)
  - shared - shared code/utilities across modules
  - auth - authentication bc (Generic) — merchant identity; NOT shopId-scoped
  - customer - customer bc (Supporting) — shopper profiles, address book; per-store
  - order - order bc (Core) — cart, checkout, order lifecycle; per-store
  - payment - payment bc (Supporting) — gateway, capture, refund; per-store
  - product - product bc (Core) — catalog, variants, stock; per-store
  - pricing - pricing bc (Supporting) — promotions, coupons, tax; per-store
  - review - review bc (Supporting) — ratings, moderation; per-store
  - reporting - reporting bc (Supporting) — analytics read model; per-store
  - settings - settings bc (Generic) — store-level config; per-store
  - shipping - shipping bc (Supporting) — zones, labels, tracking, RMA; per-store
  - notifications - notifications bc (Generic) — email/SMS dispatch; per-store
  - audit-trail - audit trail bc (Generic) — append-only compliance log; per-store
  - extension - extension bc (Generic) — webhooks, plugins; per-store

You should load the skill (only frontmatter) if not done yet.
