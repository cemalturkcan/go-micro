# Microservices Project with Go

This project is a microservices example built with Go. It consists of the following microservices:

- **Authentication**
- **Product Catalog**
- **Product Category**

Each microservice is registered with Consul. Services communicate with each other via gRPC, and authentication is managed using Keycloak.

## Microservice Architecture

Microservices are structured using `common/app` for a generic setup:

![Microservices Structure](https://github.com/user-attachments/assets/68a19ccf-87bb-4d3c-b612-ecf6463ee5c9)

### Configuration

Each microservice has its own `.env` file where all service configurations are defined:

![Configuration File](https://github.com/user-attachments/assets/144c19cf-6ad9-42d6-a50a-2f40dbfb8810)

### Database and Migrations

Each service has its own dedicated database and migration files, which can be found under the `migrations` directory inside the respective service.

## Authentication

### Register

To register a new user:

```
POST http://authentication.localhost/register
```

Response will include a token:

![Register Response](https://github.com/user-attachments/assets/6be099de-0290-4092-ba63-58b62d7047a9)

### Login

Upon login, a token is returned:

![Login Response](https://github.com/user-attachments/assets/5e994a2b-463f-4620-9bd0-629b4266bf87)

## Role-Based Access Control (RBAC)

Services register the `AuthenticationMiddleware`:

![Middleware Registration](https://github.com/user-attachments/assets/92ede1ea-4a16-492a-87cd-83897e67647f)

Unauthorized routes are defined during registration. Each request is validated by the middleware, which checks the token with Keycloak.

Example role validation:

![Role Validation](https://github.com/user-attachments/assets/35999f6a-af72-4daf-a128-f8f87252ce91)

Role module:

![Role Module](https://github.com/user-attachments/assets/5074bbc5-724c-43f1-a7c9-34fece30aca9)

## gRPC Communication

Microservices communicate via gRPC. gRPC connections can be established as follows:

![gRPC Connection](https://github.com/user-attachments/assets/f612b8c4-75a7-49c2-bf5a-9c3399303d85)

Load balancing is handled through Consul registry, ensuring healthy services are selected every 10 seconds.

Resolver location: `common/registry/resolver`

![Resolver](https://github.com/user-attachments/assets/4b7f74f8-08e1-4845-b689-4edd84101b57)

## Running the Project

1. Start the services with Docker Compose:
   ```sh
   docker-compose up -d
   ```
2. Run the initial SQL setup.
3. Configure Keycloak settings through its admin interface.
4. Install dependencies and start services:
   ```sh
   make tidy
   make watch # Starts all microservices
   ```
5. Access microservices via:
   ```
   http://{microservice-name}.localhost
   
