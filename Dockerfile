# Use Go 1.24 image
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy .env file
COPY .env .env

# Copy all source code
COPY . .

# Run tests (fail build if tests fail)
RUN go test -v -tags=integration ./...

# Build the binary
RUN go build -o main ./cmd/app

# Final lightweight image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

CMD ["./main"]
