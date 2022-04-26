FROM golang:1.16

COPY test/integration /go/src/test
WORKDIR /go/src/test
RUN env CGO_ENABLED=0 go test . -tags integration -v -c -o bin/integration

FROM 003664043471.dkr.ecr.us-east-1.amazonaws.com/e2e-base-image:latest

ENV PROJECT_NAME=bootstrap-api
COPY --from=0 /go/src/test/bin/integration /usr/bin/integration