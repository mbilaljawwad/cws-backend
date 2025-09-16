# Development Dockerfile with Hot Reloading using Air
# Multi-stage build for optimization

# Stage 1: Base image with Go and development tools
FROM golang:1.25-alpine AS base

# Install git and other necessary tools
RUN apk add --no-cache git curl

# Set working directory
WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Stage 2: Dependencies
FROM base AS deps

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Stage 3: Development environment
FROM deps AS development

# Copy Air configuration
COPY .air.toml ./

# Copy source code
COPY . .

# Expose the port
EXPOSE 9000

# Create a non-root user for security
# This command creates a new group named 'golang' with group ID 1001 (using the -S flag for a system group),
# and then creates a new user named 'golang' with user ID 1001 (using the -S flag for a system user).
# This non-root user will be used to run the application for improved security.
RUN addgroup -g 1001 -S golang && \
    adduser -S golang -u 1001

# Change ownership of the app directory and Go cache
RUN chown -R golang:golang /app && \
    chown -R golang:golang /go

# Switch to non-root user
USER golang

# Command to run Air for hot reloading
CMD ["air", "-c", ".air.toml"]
