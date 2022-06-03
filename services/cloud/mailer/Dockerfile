FROM alpine:3.13.5

COPY bin/mailer /usr/bin/mailer
COPY templates/ /usr/bin/templates/
WORKDIR /usr/bin

CMD ["/usr/bin/mailer"]