FROM surnet/alpine-wkhtmltopdf:3.13.5-0.12.6-small

RUN mkdir -p /srv/static

RUN addgroup -S nonroot \
    && adduser -S ukama -G nonroot

USER ukama


COPY bin/generator /usr/bin/generator

ENTRYPOINT ["/usr/bin/generator"]
