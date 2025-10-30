FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy & download dependencies
COPY go.mod go.sum ./
RUN go mod download -x

# Copy semua kode
COPY . .

# Build binary
RUN go build -o /bin/app ./main.go

# Stage final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/app .

COPY database ./database

ENTRYPOINT ["./app"]
