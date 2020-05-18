FROM nginx:alpine

# Copy in custom reverse proxy config
RUN rm /etc/nginx/conf.d/*
COPY build/api-gateway.conf /etc/nginx/conf.d/api-gateway.conf
