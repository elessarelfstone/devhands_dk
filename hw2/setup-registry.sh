#!/bin/bash
# Run this ONCE before the build race to pre-cache everything locally.
# After this, docker compose build will never touch the internet.

set -e

REGISTRY="localhost:5000"

echo "==> Starting local registry..."
if ! docker ps --format '{{.Names}}' | grep -q '^registry$'; then
  docker run -d -p 5000:5000 --restart=always --name registry registry:3
  echo "    Registry started."
else
  echo "    Registry already running."
fi

echo ""
echo "==> Pulling and pushing base images to local registry..."

images=(
  "golang:1.24-alpine"
  "python:3.13-slim-trixie"
  "nginx:1.24-alpine"
  "postgres:15-alpine"
  "redis:7-alpine"
  "alpine:3.21"
)

for img in "${images[@]}"; do
  echo "  -> $img"
  docker pull "$img"
  docker tag "$img" "$REGISTRY/$img"
  docker push "$REGISTRY/$img"
done

echo ""
echo "==> All base images cached. You're ready to run:"
echo "    docker compose build"
