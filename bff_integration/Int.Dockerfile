FROM 003664043471.dkr.ecr.us-east-1.amazonaws.com/e2e-base-image:latest

ENV PROJECT_NAME=bff-api
COPY bff_integration/test/integration/bin/integration /usr/bin/integration