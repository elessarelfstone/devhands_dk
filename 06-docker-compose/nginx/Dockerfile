FROM nginx:1.24-alpine

# Create necessary directories and set permissions
RUN mkdir -p /var/cache/nginx && \
  chown -R nginx:nginx /var/cache/nginx && \
  chmod -R 755 /var/cache/nginx

# Copy configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Set proper permissions for nginx directories
RUN chown -R nginx:nginx /var/log/nginx /var/run

CMD ["nginx", "-c", "/etc/nginx/conf.d/default.conf", "-g", "daemon off;"]
