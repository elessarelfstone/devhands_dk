skopeo copy \
  --dest-tls-verify=false \
  docker-archive:./alpine-latest-image.tar \
  docker://127.0.0.1:5000/library/alpine:latest
