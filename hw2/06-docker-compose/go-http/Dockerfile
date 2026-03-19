FROM golang:1.24-alpine

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev curl

# Set working directory
WORKDIR /app

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o http-server

RUN addgroup -g 1000 appuser && \
  adduser -u 1000 -G appuser -D appuser
USER appuser

HEALTHCHECK --interval=30s --timeout=5s CMD curl -f http://localhost:8080/health
CMD ["/app/http-server"]
