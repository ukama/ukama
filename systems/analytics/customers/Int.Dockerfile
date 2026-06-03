FROM 003664043471.dkr.ecr.us-east-1.amazonaws.com/e2e-base-image:latest

ENV PROJECT_NAME=analytics-customers
COPY bin/customers /usr/bin/customers
