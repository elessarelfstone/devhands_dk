ARG GO_VERSION=1.21
ARG ALPINE_VERSION=3.19

# Stage 1: Build
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go mod init deployment-info && \
	CGO_ENABLED=0 GOOS=linux go build -o deployment-info .

# Stage 2: Runtime
FROM alpine:${ALPINE_VERSION}
WORKDIR /app

# Create non-root user
RUN addgroup -S www-data ; \
	adduser -S www-data -G www-data --home /app --no-create-home ; true

# Create secrets directory (owned by www-data)
RUN mkdir -p /run/secrets && \
	chown www-data:www-data /run/secrets && \
	chmod 750 /run/secrets

# Copy binary
COPY --from=builder --chown=www-data:www-data /app/deployment-info .

# Environment variables (NO API_KEY!)
ENV DEPLOYMENT_TYPE="unknown" \
	ENDPOINT_PATH="/v1/deployment-type" \
	HOST="0.0.0.0" \
	PORT="8080"

EXPOSE 8080
USER www-data:www-data
CMD ["./deployment-info"]
