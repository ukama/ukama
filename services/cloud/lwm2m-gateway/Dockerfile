FROM alpine:3.13.5

RUN apk update && apk add --no-cache supervisor

RUN mkdir -p /var/log/supervisor

ENTRYPOINT []
WORKDIR /

COPY bin/lwm2mGateway /sbin/lwm2mGateway

COPY configs /etc/configs

CMD ["/usr/bin/supervisord", "-c", "/etc/configs/supervisord.conf"]

