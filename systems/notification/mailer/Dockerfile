FROM alpine:3.13.5

COPY bin/mailer /usr/bin/mailer
COPY templates/test-template.tmpl /templates/test-template.tmpl
COPY templates/member-invite.tmpl /templates/member-invite.tmpl
COPY templates/sim-allocation.tmpl /templates/sim-allocation.tmpl


CMD ["/usr/bin/mailer"]
