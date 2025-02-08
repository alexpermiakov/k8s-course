FROM --platform=linux/amd64 nginx:latest

COPY nginx/nginx.conf /etc/nginx/nginx.conf

COPY index.html /usr/share/nginx/html

EXPOSE 9090 

CMD ["nginx", "-g", "daemon off;"]