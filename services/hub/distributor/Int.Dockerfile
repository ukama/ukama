FROM 003664043471.dkr.ecr.us-east-1.amazonaws.com/e2e-base-image:latest 

ENV PROJECT_NAME=distributor

RUN mkdir -p /usr/bin/testdata
COPY bin/integration /usr/bin/integration
COPY test/integration/testdata /testdata
