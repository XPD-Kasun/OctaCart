# OctaCart e-Commerce Platform

OctaCart is a modular monolith backend application tailored for small-scale e-commerce sellers. It provides:
- A REST API (powered by Gin)
- A GraphQL endpoint for products
- Integration capabilities for an admin dashboard

## Project Structure

- Backend: The main Go backend application.
  - cmd: Application entry point.
  - internal: Core application logic using Hexagonal Architecture.
    - driving: Inbound adapters (REST, GraphQL).
    - driven: Outbound adapters.
    - usecases: Business logic orchestration.
    - services: Core domain services.
- Frontend The nextjs frontend
## License

This project is licensed under the Apache License 2.0. See the \LICENSE\ file for details.
