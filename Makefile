.PHONY: network-up local-up local-down clean test test-verbose test-coverage generate-swagger-docs build-prod-image

network-up:
	@docker network inspect app-network >/dev/null 2>&1 || docker network create --driver bridge app-network

local-up: network-up
	@if [ ! -f .env ]; then \
		echo "Creating .env file from .env.sample..."; \
		cp .env.sample .env; \
	fi
	docker compose up -d --build

local-down: network-up
	docker compose down

logs-app:
	docker logs -f go-app

test:
	docker exec -it go-app go test ./...

test-verbose:
	docker exec -it go-app go test ./... -v

test-coverage:
	docker exec -it go-app go test ./... -coverprofile=coverage.out
	docker exec -it go-app go tool cover -html=coverage.out -o coverage.html
	docker cp go-app:/app/coverage.html ./coverage.html
	@echo "Coverage report generated at ./coverage.html"

generate-swagger-docs:
	@echo "Checking if swag is installed..."
	@if ! command -v swag &> /dev/null; then \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@echo "Generating Swagger documentation..."
	swag init -g cmd/http/main.go -o docs
	@echo "Swagger documentation generated successfully"

build-prod-image:
	@echo "Building production Docker image..."
	docker build -t go-pack-calculator:prod -f Dockerfile.prod .
	@echo "Production Docker image built successfully"

clean: local-down
	docker volume rm -f $(shell docker volume ls -q | grep go-task-assessment)
	@echo "Docker volumes cleaned successfully"