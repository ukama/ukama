FROM ubuntu:latest

LABEL maintainer="hello@ukama.com"

RUN apt-get update && apt-get install -y \
    software-properties-common \
    && add-apt-repository universe \
    && apt-get update && apt-get install -y \
    build-essential \
    git \
    wget \
    autoconf \
    automake \
    libtool \
    pkg-config \
    libssl-dev \
    texinfo \
    cmake \
    tcl \
    zlib1g-dev \
    texlive \
    texlive-latex-extra \
    ghostscript \
    gperf \
    gtk-doc-tools \
    libev-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace