FROM golang:1.24-alpine AS builder

# Install build dependencies for CGO and swag
RUN apk add --no-cache gcc musl-dev sqlite-dev git

# Install swag for API documentation generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY . .

RUN go mod download

# Generate Swagger documentation
RUN swag init -g cmd/docs.go -o swagger --parseDependency --parseInternal

# Update dependencies after generating docs
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=1 go build -o gh-inspector .

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite-libs curl

WORKDIR /app

COPY --from=builder /app/gh-inspector .
COPY configs configs/

ENTRYPOINT ["./gh-inspector"]
CMD ["--help"]
