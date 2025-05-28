# 4IT428GoProject

This project is part of the course **Vývoj mikroslužeb v jazyce Go** at the University of Economics, Prague.

## How to run?

**1. Clone the repository**

**2. Update your environment variables for each service and docker-compose**

**3. Inside the services folder execute:**
    ```
    make build
    ```


## Project Structure
Project is divided into multiple services, each service is in its own directory. Each service has its own .env file, where you can set environment variables for the service.

- nginx: main entry point
- editor-service: handles register and login requests
- newsletter-service: handles CRUD operations for newsletters and posts
- subscriber-service: handles subscriptions and sending of published posts to subscribers

    ```
    .
    └── services/
        ├── nginx/
        │   ├── nginx.conf
        │   └── Dockerfile
        ├── editor-service/
        │   ├── .env
        │   ├── Dockerfile
        │   ├── firebase-admin-sdk.json
        │   ├── cmd/
        │   ├── pkg/
        │   ├── models/
        │   ├── repository/
        │   ├── service/
        │   └── transport/
        ├── newsletter-service/
        │   ├── .env
        │   ├── Dockerfile
        │   ├── firebase-admin-sdk.json
        │   ├── cmd/
        │   ├── middleware/
        │   ├── repository/
        │   ├── service/
        │   └── transport/
        ├── subscriber-service/
        │   ├── .env
        │   ├── Dockerfile
        │   ├── firebase-admin-sdk.json
        │   ├── cmd/
        │   ├── pkg/
        │   ├── repository/
        │   ├── service/
        │   └── transport/
        └── docker-compose.yaml
    ```