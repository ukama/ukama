FROM alpine:3.13.5

COPY bin/device-gateway /usr/bin/device-gateway
COPY swagger-ui /var/www/swagger-ui
ENV SWAGGERASSETS=/var/www/swagger-ui


CMD ["/usr/bin/device-gateway"]