# Dockerfile specifically for Swagger documentation generation
FROM golang:1.25.0-alpine AS swagger-generator

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate swagger documentation
RUN swag init -g cmd/main.go --parseDependency --parseInternal -o ./docs/swagger

# Final stage - just the generated docs
FROM scratch AS docs-output
COPY --from=swagger-generator /app/docs/swagger /docs/swagger

# Development stage - includes the generator for live updates
FROM golang:1.25.0-alpine AS swagger-dev
RUN go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
CMD ["sh", "-c", "swag init -g cmd/main.go --parseDependency --parseInternal -o ./docs/swagger && echo 'Swagger docs generated successfully!'"]
