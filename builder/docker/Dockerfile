ARG TARGETPLATFORM
FROM ${TARGETPLATFORM} as base

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
    kpartx \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && \
    apt-get install -y \
    fdisk \
    util-linux \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

# Final stage with platform-specific logic
FROM base

ARG TARGETPLATFORM

# Install extlinux conditionally
RUN if [ "$TARGETPLATFORM" = "amd64/ubuntu:latest" ]; then \
        apt-get update && apt-get install -y extlinux; \
    else \
        wget http://archive.ubuntu.com/ubuntu/pool/universe/s/syslinux/extlinux_6.04~git20190206.bf6db5b4+dfsg1-2_amd64.deb \
        && sha256sum extlinux_6.04~git20190206.bf6db5b4+dfsg1-2_amd64.deb \
        && dpkg -i extlinux_6.04~git20190206.bf6db5b4+dfsg1-2_amd64.deb \
        && apt-get install -f \
        && rm extlinux_6.04~git20190206.bf6db5b4+dfsg1-2_amd64.deb; \
    fi

# Optionally set the working directory again if it was overridden
WORKDIR /workspace
