# Use a base image with an updated glibc version
FROM ubuntu:22.04 AS builder

# Set the working directory
WORKDIR /app

# Install build tools, Go, and dependencies
RUN apt-get update && apt-get install -y \
    wget \
    build-essential \
    golang \
    ca-certificates \
    && apt-get clean

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Create a smaller runtime image
FROM ubuntu:22.04

# Set the working directory
WORKDIR /app

# Install required runtime libraries, including CA certificates
RUN apt-get update && apt-get install -y \
    libc6 \
    ca-certificates \
    && apt-get clean

# Copy the built application from the builder stage
COPY --from=builder /app/main /app/main

# Expose the application's port
EXPOSE 8080

# Start the application
CMD ["/app/main"]
