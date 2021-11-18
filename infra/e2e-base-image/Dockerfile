FROM golang:1.16

RUN go get -u github.com/jstemmer/go-junit-report

FROM alpine:3.13.5

RUN apk add --no-cache \
        libc6-compat curl

RUN  curl https://github.com/testhub-io/testhub-cli/releases/download/v0.14/testhub-cli-v0.14-linux-amd64.tar.gz --output testhub-cli.tar.gz -L && tar -xzf testhub-cli.tar.gz && \
        mv ./testhub-cli /usr/bin/testhub-cli


COPY --from=0 /go/bin/go-junit-report /usr/bin/go-junit-report
COPY --from=0 /usr/local/go/pkg/tool/linux_amd64/test2json /usr/bin/test2json


RUN  curl https://github.com/gotestyourself/gotestsum/releases/download/v1.7.0/gotestsum_1.7.0_linux_amd64.tar.gz --output gotestsum.tar.gz -L && tar -xzf gotestsum.tar.gz && \
        mv ./gotestsum /usr/bin/gotestsum

ENV ON_PREMISE=true
ENV TESTHUB_DOMAIN=testhub-api.dev.ukama.com

RUN echo "gotestsum --raw-command --junitfile report.xml --  test2json /usr/bin/integration -test.v" > run.sh
RUN echo "export TESTS_RESULT=\$?" >> run.sh
RUN echo "testhub-cli upload -p ukama/\$PROJECT_NAME -f report.xml  --build \$(date '+%Y-%m-%d-%H-%M')" >> run.sh
RUN echo "exit \$TESTS_RESULT" >> run.sh

CMD ["/bin/ash",  "run.sh" ]