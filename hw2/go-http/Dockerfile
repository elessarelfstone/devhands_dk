# syntax=docker/dockerfile:1

# ── Stage 1: build ────────────────────────────────────────────────────────────
FROM localhost:5000/golang:1.24-alpine AS builder

# Only what's needed to compile
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Cache dependency downloads separately from source code.
# This layer is reused on every rebuild as long as go.mod/go.sum don't change.
COPY go.mod go.sum ./
RUN go mod download

# Now copy source and build
COPY . .
RUN CGO_ENABLED=0 go build -o http-server

# ── Stage 2: minimal runtime ──────────────────────────────────────────────────
# alpine instead of scratch because we need curl for the healthcheck
FROM localhost:5000/alpine:3.21

RUN apk add --no-cache curl

RUN addgroup -g 1000 appuser && \
    adduser -u 1000 -G appuser -D appuser

COPY --from=builder /app/http-server /app/http-server

USER appuser

HEALTHCHECK --interval=30s --timeout=5s CMD curl -f http://localhost:8080/health
CMD ["/app/http-server"]
