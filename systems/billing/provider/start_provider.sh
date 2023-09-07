#!/usr/bin/env sh
#
# This script is used to bootstrap the billing provider
# with the needed configuration. You need to have both
# openssl and base64 installed in order to run this script.

# Set up environment configuration
if [ ! -f .env ]
then
  echo "LAGO_RSA_PRIVATE_KEY=\"`openssl genrsa 2048 | base64`\"" >> .env
fi

source .env

# Start
docker-compose -f provider.yml up
