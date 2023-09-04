FROM 003664043471.dkr.ecr.us-east-1.amazonaws.com/e2e-base-image:latest

RUN addgroup -S nonroot \
    && adduser -S ukama -G nonroot

USER ukama

ENV PROJECT_NAME=node-gateway
COPY bin/node-gateway /usr/bin/node-gateway
