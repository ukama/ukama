FROM surnet/alpine-wkhtmltopdf:3.13.5-0.12.6-small

RUN mkdir -p /srv/static

COPY bin/report /usr/bin/report

ENTRYPOINT ["/usr/bin/report"]
