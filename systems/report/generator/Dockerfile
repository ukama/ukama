FROM surnet/alpine-wkhtmltopdf:3.13.5-0.12.6-small

RUN addgroup -S nonroot \
    && adduser -S ukama -G nonroot

USER ukama

RUN mkdir -p /srv/static

COPY bin/generator /usr/bin/generator

ENTRYPOINT ["/usr/bin/generator"]
