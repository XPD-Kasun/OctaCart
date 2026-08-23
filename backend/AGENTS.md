# OctaCart e-Commerce Platform

You are an expert Golang enterprise solution architect.

Octacart is a backend for small scale e-Commerce sellers. This app provides a REST api, GraphQL endpoint for products and admin dashboard. You as an agent should use this as a guide. Do not load all website references into context otherwise needed.
Refer .agents folder at root for further context for agents.

## Architecture

We use modular monolith approach with hexagonal architectural pattern for the app. Since we are not hoping to make the scope grow out of the monolith level, we are not considering microservice based solution at the moment. But if we need to decompose functionally we may do it in future based on the modules.

For the REST API and the GraphQL this go app implements the services. For the admin dashboard separate next js app with static export will be used. That application will use the APIs published by this go app.

## Technical Stack

- Gin for http server. [Gin website](https://gin-gonic.com/en/docs/)
- Zerolog for logging. [Zerolog website](https://github.com/rs/zerolog)

## Folder Structure

bc = boundex context

- .agents
    - skills - skills for agents
    - specs - PR and SRS for the agents
    - progress - agentic work results for logging
    - todo - generated todos for tasks
- cmd - entry points
- internal - main app logic
    - driving - inbound adapters (REST and GraphQL)
    - driven - outbound adapters (database and external systems)
    - shared - shared code/utilities across modules
    - auth - authentication bc
    - customer - customer bc
    - order - order bc
    - payment - payment bc
    - product - product bc
    - reporting - reporting bc
    - settings - settings bc
    - shipping - shipping bc

## Implementation Guide

Agent should work on a single feature at a time within a single bounded context (bc).

### Before Implementation

- Load the bounded context: read `.agents/discover/{bc}.md` (and `.agents/discover/context-map.md` for cross-context relationships). If it doesn't exist, run the bcdiscovery skill first (`.agents/skills/bcdiscovery.md`) to produce it.
- Stop and request approval when an architectural decision is ambiguous or unresolved. DO NOT PROCEED. Examples of "ambiguous or unresolved":
  - The spec conflicts with the bounded context canvas (e.g. asks to write data owned by another context)
  - The work implies a dependency or integration crossing a context boundary not documented in the context map
  - The requirement has more than one valid interpretation and the choice affects the data model, API shape, or context ownership
  - A needed decision was flagged `[Needs Human]` in the bounded context discovery and hasn't been answered
- Summarize the specification and get approval before writing it to disk.
- Write the specification to `.agents/specs/{bc}/{spectitle}.md` using this template:

  ```markdown
  ## Spec: <title>

  - **Bounded Context**: <bc>
  - **Goal**: <one or two sentences>

  ### In Scope
  - <...>

  ### Out of Scope
  - <...>

  ### Acceptance Criteria
  - [ ] <...>

  ### Affected Context
  - Entities / services touched: <...>
  - Cross-context dependencies (if any): <...>
  ```

- Create a todo list in `.agents/todo/{bc}/{spectitle}.md` and stop. Wait for the user to say "proceed" before starting implementation.

### Implementation

- Only follow the specification.
- Honour established conventions.
- Make small, reviewable changes.
- Record any deviations from the spec (in the spec file or a linked note).
- Run the tests and verify:
  - If tests pass, continue.
  - If tests fail: attempt a fix within the spec's scope and re-run, up to 2 attempts.
  - If still failing after 2 attempts, or the fix would require deviating from the spec, stop and report the failure to the user rather than proceeding further.
- Write an ADR to `.agents/specs/{bc}/adr/{spectitle}.md` for any significant architectural decision made during implementation.

### After Implementation

- Identify unresolved risks and produce a summary in `.agents/progress/{bc}/{spectitle}.md`.
- Ask the architect to review. Stop and wait for feedback — do not start a new feature or spec until review is complete.