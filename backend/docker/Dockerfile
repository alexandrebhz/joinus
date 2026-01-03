FROM golang:1.25-alpine AS builder

# Install git and ca-certificates (needed for go mod download with private repos)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go.mod first
COPY go.mod ./

# Download dependencies (this will work even without go.sum)
RUN go mod download

# Copy the rest of the application
COPY . .

# Generate go.sum and verify dependencies
RUN go mod tidy && go mod verify

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]


