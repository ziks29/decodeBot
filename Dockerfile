FROM golang:1.24-alpine AS builder
WORKDIR /src

# Install git and CA certs for fetching modules
RUN apk add --no-cache git ca-certificates

# Cache modules
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct
RUN go mod download

# Copy source and build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot cmd/bot/main.go

# Runtime stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# Copy binary from builder
COPY --from=builder /src/bot /app/bot

# Set timezone for scheduler/notification functionality
ENV TZ=Europe/Warsaw

ENTRYPOINT ["/app/bot"]
