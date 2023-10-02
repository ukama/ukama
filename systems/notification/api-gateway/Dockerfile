FROM alpine:3.13.5

RUN addgroup -S nonroot \
    && adduser -S ukama -G nonroot

USER ukama

COPY bin/api-gateway /usr/bin/api-gateway

CMD ["/usr/bin/api-gateway"]
