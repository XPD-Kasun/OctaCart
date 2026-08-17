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

- .agents
    - skills - skills for agents
    - specs - PR and SRS for the agents
    - progress - agentic work results for logging
    - todo - generated todos for tasks
- cmd - entry points
- internal - main app logic
    - driving - inbound adapters (REST and GraphQL)
    - driven - outbound adapters (database and external systems)
    - usecases - business use cases orchestration
    - services - core domain services/business logic (auth, product, order, payment, customer, shipping, reporting, settings)



