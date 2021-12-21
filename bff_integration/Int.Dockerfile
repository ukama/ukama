FROM golang:1.16

COPY test/integration /go/src/test
WORKDIR /go/src/test

FROM 003664043471.dkr.ecr.us-east-1.amazonaws.com/e2e-base-image:latest

ENV PROJECT_NAME=bff-api
COPY bin/integration /usr/bin/integration