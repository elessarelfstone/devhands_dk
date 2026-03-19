mkdir -p storage
chmod og+rwX storage
docker run -d --name harbor-registry \
  -p 5000:5000 \
  -v $PWD/config.yml:/etc/registry/config.yml \
  -v $PWD/storage:/storage \
  bitnami/harbor-registry:latest
