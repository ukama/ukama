FROM alpine:3.13.5

COPY bin/analytics /usr/bin/analytics
COPY bin/integration /usr/bin/integration

CMD ["/usr/bin/analytics"]
