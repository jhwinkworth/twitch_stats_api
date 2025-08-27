# Use Go 1.24 image
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy all source code
COPY . .

# run tests (fail build if tests fail)
RUN go test -v -tags=integration ./...

# build the binary
RUN go build -o main ./cmd/app

# final lightweight image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
