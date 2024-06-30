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
    bison \
    jq

RUN apt-get update && apt-get install -y \
    qemu-system \
    virt-manager \
    virt-viewer \
    libvirt-daemon-system \
    libvirt-clients \
    bridge-utils \
    debootstrap \
    extlinux \
    kpartx \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace