ARG GO_VERSION=1.21
ARG ALPINE_VERSION=3.19
ARG APPNAME=deployment-info
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app

# Copy and build
COPY main.go .
RUN go mod init deoployment-info;CGO_ENABLED=0 GOOS=linux go build -o deployment-info .

# Stage 2: Runtime
FROM alpine:${ALPINE_VERSION}
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/${APPNAME} .

# Set environment variables with defaults
ENV API_KEY="" \
  DEPLOYMENT_TYPE="unknown" \
  ENDPOINT_PATH="/v1/deployment-type" \
  HOST="0.0.0.0" \
  PORT="8080"

# Expose port and run
EXPOSE 8080
CMD ["./"${APPNAME}]
