FROM alpine:3.13.5

RUN addgroup -S nonroot \
    && adduser -S ukama -G nonroot

USER ukama

COPY bin/node-gateway /usr/bin/node-gateway

CMD ["/usr/bin/node-gateway"]
