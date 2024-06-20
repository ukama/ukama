FROM ubuntu:latest

LABEL maintainer="hello@ukama.com"

RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    wget \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace