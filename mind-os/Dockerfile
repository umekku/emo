# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Install git for fetch dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary
# CGO_ENABLED=0 is default for modernc.org/sqlite compatibility in pure Go
RUN CGO_ENABLED=0 GOOS=linux go build -o mind-os .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/mind-os .

# Expose the application port
EXPOSE 8081

# Default environment variables
ENV NODE_ENV=production
ENV GIN_MODE=release
ENV PORT=8081

# Run the binary
CMD ["./mind-os"]
