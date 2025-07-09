# Development image
FROM golang:1.23-alpine

# Install ca-certificates and tools needed for development
RUN apk --no-cache add ca-certificates curl

# Install CompileDaemon for hot reloading
RUN go install github.com/githubnemo/CompileDaemon@latest

# Install Swagger tools
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set working directory
WORKDIR /app

# Create necessary directories
RUN mkdir -p ./build

# Set environment variables for color output
ENV TERM=xterm-256color
ENV COLORTERM=truecolor

# Expose the application port
EXPOSE 8080

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/http/main.go -o docs

# Command to run the executable
CMD CompileDaemon --build="go build -o ./build/main ./cmd/http" --command="./build/main" --pattern="(.+\.go|.+\.env)" --directory="." --graceful-kill=true --color=true
