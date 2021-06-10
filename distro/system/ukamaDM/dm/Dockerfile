FROM alpine:3.13.5

RUN apk update && apk add --no-cache supervisor

RUN mkdir -p /var/log/supervisor

ENTRYPOINT []
WORKDIR /

COPY build/bin/server /sbin/lwm2mServer
COPY config /etc/configs

CMD ["/usr/bin/supervisord", "-c", "/etc/configs/ssupervisord.conf"]
