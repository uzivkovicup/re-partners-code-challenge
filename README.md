# Pack Calculator Application

This application calculates the optimal number of packs needed to fulfill an order based on the following rules:

1. Only whole packs can be sent. Packs cannot be broken open.
2. Within the constraints of Rule 1 above, send out the least amount of items to fulfill the order.
3. Within the constraints of Rules 1 & 2 above, send out as few packs as possible to fulfill each order.

## Features

- REST API for pack calculation and pack size management
- Flexible pack size configuration (add, remove, update pack sizes)
- Optimized algorithm for calculating the minimum number of packs
- Modern UI for interacting with the application
- Containerized for easy deployment
- API documentation with Swagger
- Robust configuration system using Viper

## Application Screenshots

### Homepage - Pack Calculator Interface

![Pack Calculator Homepage](docs/images/homepage.png)

## Technology Stack

- Backend: Go (Golang)
- API: REST using Gin framework
- Database: PostgreSQL
- Frontend: HTML, CSS, JavaScript
- Containerization: Docker
- API Documentation: Swagger
- Configuration: Viper

## Hexagonal Architecture

This project follows the hexagonal architecture pattern (also known as ports and adapters), which provides a clean separation of concerns and makes the system more maintainable, testable, and adaptable to changes.

### Architecture Overview

![Hexagonal Architecture](https://miro.medium.com/max/1400/1*NfFzI7Z-E3ypn8ahESbDzw.png)

The architecture is organized into the following layers:

#### 1. Domain Layer

The core of the application that contains the business logic and domain entities. It is independent of any external concerns.

- `internal/domain/entities`: Domain entities that represent the core business objects
- `internal/domain/errors`: Domain-specific errors
- `internal/domain/services`: Domain services that implement core business logic

#### 2. Application Layer

The application layer coordinates the interactions between the domain layer and the external world. It implements use cases by orchestrating domain objects.

- `internal/application/usecases`: Use cases that implement specific application features
- `internal/application/services`: Services that implement the primary ports (input ports)

#### 3. Ports Layer

Interfaces that define how the application interacts with the outside world.

- `internal/ports/primary`: Input ports (interfaces implemented by the application layer)
- `internal/ports/secondary`: Output ports (interfaces implemented by adapters)

#### 4. Adapters Layer

Implementations of the ports that connect the application to external systems.

- `internal/adapters/primary`: Primary adapters (REST API, CLI, etc.)
- `internal/adapters/secondary`: Secondary adapters (PostgreSQL, in-memory repository, etc.)

### Flow of Control

1. External requests come through primary adapters (e.g., REST API)
2. Primary adapters use primary ports (interfaces) to communicate with the application layer
3. Application layer implements use cases by orchestrating domain objects
4. Application layer uses secondary ports (interfaces) to communicate with external systems
5. Secondary adapters implement secondary ports to connect to external systems (e.g., database)

### Benefits of Hexagonal Architecture

1. **Separation of Concerns**: Each layer has a specific responsibility
2. **Testability**: Domain and application layers can be tested in isolation
3. **Flexibility**: Easy to replace adapters without changing the core logic
4. **Maintainability**: Clear boundaries between layers make the code easier to understand and maintain
5. **Domain-Driven Design**: Focus on the domain model and business rules

### Implementation Details

#### Domain Entities

- `PackSize`: Represents a pack size with validation rules
- `CalculationResult`: Represents the result of a pack calculation

#### Use Cases

- `PackSizeUseCase`: Manages pack size operations (CRUD)
- `CalculationUseCase`: Calculates optimal packs for orders

#### Ports

- Primary Ports:
  - `PackSizeService`: Interface for pack size operations
  - `CalculationService`: Interface for calculation operations

- Secondary Ports:
  - `PackSizeRepository`: Interface for pack size persistence

#### Adapters

- Primary Adapters:
  - `rest.PackCalculatorHandler`: Handles HTTP requests

- Secondary Adapters:
  - `postgres.PackSizeRepository`: PostgreSQL implementation
  - `inmemory.PackSizeRepository`: In-memory implementation for testing

#### Dependency Flow

The dependencies flow inward, with the domain layer at the center having no dependencies on other layers:

1. Domain layer has no external dependencies
2. Application layer depends only on the domain layer
3. Ports layer defines interfaces for the application layer
4. Adapters layer depends on the ports layer

This ensures that the core business logic remains isolated from external concerns and can evolve independently.

## Configuration

The application uses Viper for configuration management, which provides several benefits:

- Environment variables take precedence over .env file values
- Sensible defaults for development
- Support for multiple environments (development, production)
- Automatic reading of .env files when present

Configuration can be provided in two ways:

1. **Environment Variables**: Set the following environment variables:
   ```
   ENVIRONMENT=production
   PROTOCOL=https
   HOST=0.0.0.0
   PORT=8080
   POSTGRES_DB_HOST=your-db-host
   POSTGRES_DB_PORT=5432
   POSTGRES_DB_USER=your-username
   POSTGRES_DB_PASSWORD=your-password
   POSTGRES_DB_NAME=your-db-name
   POSTGRES_DB_SSLMODE=require
   ```

2. **.env File**: Create a .env file in the project root (a sample is provided in .env.sample)

## Running the Application

### Prerequisites

- Docker and Docker Compose installed
- Make

### Using Make

The project includes a Makefile for convenience. Follow these steps:

1. Clone the repository
2. Run the application (this will automatically create `.env` from `.env.sample` if needed):

```bash
make local-up
```

3. Access the application:
   - Pack Calculator UI: http://localhost:8080
   - REST API: http://localhost:8080/api
   - Swagger Documentation: http://localhost:8080/swagger/index.html

### Available Make Commands

Here are all the available Make commands:

```bash
# Create Docker network if it doesn't exist
make network-up

# Start the application (creates .env from .env.sample if needed)
make local-up

# Stop the application
make local-down

# View application logs
make logs-app

# Run tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Generate Swagger documentation
make generate-swagger-docs

# Build the production Docker image
make build-prod-image

# Clean up Docker volumes
make clean
```

### Using Docker Compose Directly

1. Clone the repository
2. Create a `.env` file based on the `.env.sample` file
3. Create the Docker network:

```bash
docker network create --driver bridge app-network
```

## Deployment

### Deploying to Render

This application can be deployed to Render using the provided `render.yaml` configuration file:

1. Create a Render account at https://render.com
2. Connect this GitHub repository to Render
3. Click "New" and select "Blueprint" from the dropdown
4. Select this repository and click "Apply Blueprint"
5. Render will automatically create the web service and PostgreSQL database using the production Dockerfile
6. Deploy the application

Once deployed, this application will be available at the URL provided by Render.

### Production Docker Setup

The application includes a production-optimized Dockerfile (`Dockerfile.prod`) that:

- Uses a multi-stage build to create a smaller final image
- Compiles the application with optimizations
- Includes only the necessary files in the final image
- Sets appropriate environment variables for production
- Uses .dockerignore to exclude unnecessary files

To build and run the production Docker image locally:

```bash
# Build the production image using Make
make build-prod-image

# Or build directly with Docker
docker build -t go-pack-calculator:prod -f Dockerfile.prod .

# Run the production container
docker run -p 8080:8080 go-pack-calculator:prod
```

## API Documentation

The application provides a REST API with the following endpoints:

### Endpoints

#### Pack Sizes

- `GET /api/pack-sizes`: Get a paginated list of pack sizes
  - Query parameters: `size`, `limit`, `page`
- `GET /api/pack-sizes/:id`: Get a pack size by ID
- `POST /api/pack-sizes`: Create a new pack size
  - Request body: `{ "size": 250 }`
- `PUT /api/pack-sizes/:id`: Update an existing pack size
  - Request body: `{ "size": 500 }`
- `DELETE /api/pack-sizes/:id`: Delete a pack size

#### Pack Calculation

- `POST /api/calculate-packs`: Calculate the optimal packs for an order
  - Request body: `{ "itemsOrdered": 501 }`

For detailed API documentation, visit the Swagger UI at `/swagger/index.html` when the application is running.

## Algorithm

The pack calculation algorithm uses dynamic programming to find the optimal solution that satisfies all the requirements. It works as follows:

1. First, it finds all possible combinations that satisfy the order with the minimum number of items (rule 2)
2. Then, among those combinations, it selects the one with the fewest packs (rule 3)

For very large orders or specific pack size configurations, the algorithm uses a greedy approach as a fallback.

## Testing

The application includes unit tests for the core algorithm and API endpoints. To run the tests:

```bash
make test
make test-verbose
make test-coverage
```

## Recent Refactoring to Hexagonal Architecture

The project was recently refactored to follow a fully hexagonal architecture pattern. This refactoring improved the code organization, testability, and maintainability by clearly separating concerns and dependencies.

### Key Changes

1. **Restructured Project Layout**:
   - Organized code into domain, application, ports, and adapters layers
   - Removed redundant files and consolidated functionality
   - Created clear separation between business logic and infrastructure concerns

2. **Domain Layer**:
   - Created clean domain entities with proper validation
   - Implemented domain services with core business logic
   - Added domain-specific error types for better error handling

3. **Application Layer**:
   - Implemented use cases that orchestrate domain objects
   - Created application services that implement the primary ports
   - Added support for pagination in the application layer

4. **Ports Layer**:
   - Defined clear interfaces for primary (input) and secondary (output) ports
   - Established contracts between layers

5. **Adapters Layer**:
   - Implemented REST API handlers as primary adapters
   - Created PostgreSQL and in-memory repository implementations as secondary adapters

6. **Testing**:
   - Added comprehensive tests for domain entities and services
   - Ensured all tests pass with the new architecture

### Benefits Achieved

1. **Improved Testability**: Domain logic can be tested in isolation without dependencies on external systems
2. **Better Separation of Concerns**: Each component has a single responsibility
3. **Enhanced Maintainability**: Changes in one layer don't affect others
4. **Flexibility**: Easy to swap implementations (e.g., database, API) without changing business logic
5. **Cleaner Code**: Clear boundaries and dependencies make the code easier to understand

## Future Enhancements

1. **Add More Tests**: Expand test coverage to include application use cases and adapters
2. **Monitoring**: Add metrics and logging to track performance and issues
3. **Feature Expansion**: Implement additional features using the new architecture
4. **UI Improvements**: Enhance the user interface with more features and better visualization

## Recent Updates

### GORM Integration

The project has been updated to use GORM (Go Object Relational Mapper) for database operations instead of raw SQL queries. This provides several benefits:

1. **Cleaner Code**: GORM provides a higher-level abstraction over database operations, making the code more readable and maintainable
2. **Type Safety**: GORM maps database records to Go structs, providing type safety and reducing runtime errors
3. **Automatic Migrations**: GORM can automatically create tables and manage schema changes
4. **Pagination**: Built-in support for pagination makes it easier to implement paginated APIs
5. **Transaction Support**: GORM provides simple transaction management with Begin, Commit, and Rollback methods
6. **Hooks**: GORM supports hooks for custom logic before/after create, update, delete operations

The implementation includes:
- Mapping domain entities to GORM models
- Using GORM's query builder for database operations
- Implementing pagination using GORM's Offset and Limit methods
- Using transactions for migrations
