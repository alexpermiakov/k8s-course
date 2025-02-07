FROM nginx:latest

ENV NGINX_ROOT_PATH=/usr/share/nginx/html

COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf

COPY index.html /usr/share/nginx/html/

EXPOSE 9090 