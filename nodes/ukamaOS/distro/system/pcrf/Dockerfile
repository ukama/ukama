FROM alpine:3.13.5

RUN addgroup -S nonroot \
    && adduser -S ukama -G nonroot

USER ukama

COPY bin/pcrf /usr/bin/pcrf

CMD ["/usr/bin/pcrf"]
