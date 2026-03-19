FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY main.go .
RUN go mod init deoployment-info;CGO_ENABLED=0 GOOS=linux go build -o deployment-info .

FROM alpine:3.19
WORKDIR /root/

COPY --from=builder /app/deployment-info .

ENV API_KEY="" \
  DEPLOYMENT_TYPE="unknown" \
  ENDPOINT_PATH="/api/deployment_type" \
  HOST="0.0.0.0" \
  PORT="8080"

CMD ["./deployment-info"]
