# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o project main.go

# Stage 2: Create a smaller image for the final application
FROM alpine:latest

# Install certificates (useful for HTTPS requests if needed)
# RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/project .

# Run the Go application
CMD ["./project"]
