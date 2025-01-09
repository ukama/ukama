ARG TARGETPLATFORM=ubuntu:latest
FROM ${TARGETPLATFORM} as builder

LABEL maintainer="hello@ukama.com"

# Set working directory
WORKDIR /workspace

# Make TARGETPLATFORM available in the RUN commands
ARG TARGETPLATFORM
ENV TARGETPLATFORM=${TARGETPLATFORM}

# Install dependencies based on platform
RUN if echo "${TARGETPLATFORM}" | grep -q "alpine"; then \
        apk add --no-cache \
        build-base \
        make \
        cmake \
        git \
        bash \
        curl \
        wget \
        curl-dev \
        libmicrohttpd-dev \
        gnutls-dev \
        nettle-dev \
        p11-kit-dev \
        zlib-dev \
        jansson-dev \
        sqlite-dev \
        util-linux-dev \
        linux-headers \
        bsd-compat-headers \
        musl-dev; \
    elif echo "${TARGETPLATFORM}" | grep -q "ubuntu"; then \
        apt-get update && apt-get install -y \
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
        jq \
        qemu-system \
        virt-manager \
        virt-viewer \
        libvirt-daemon-system \
        libvirt-clients \
        bridge-utils \
        debootstrap \
        kpartx \
        fdisk \
        util-linux \
        libgnutls28-dev \
        libmicrohttpd-dev \
        libcurl4-openssl-dev \
        libjansson-dev \
        gcc-aarch64-linux-gnu \
        g++-aarch64-linux-gnu \
        uuid-dev \
        libsqlite3-dev \
        && rm -rf /var/lib/apt/lists/*; \
    else \
        echo "Unsupported platform: ${TARGETPLATFORM}" && exit 1; \
    fi

# Clone and build Orcania (common for both platforms)
RUN git clone https://github.com/babelouest/orcania.git && \
    cd orcania && \
    mkdir build && \
    cd build && \
    cmake .. && \
    make && \
    make install && \
    cd ../../ && rm -rf orcania

# Clone and build Ulfius with Yder disabled
RUN git clone https://github.com/babelouest/ulfius.git && \
    cd ulfius && \
    mkdir build && \
    cd build && \
    cmake -DWITH_YDER=off -DWITH_JOURNALD=off -DCMAKE_C_FLAGS="-Wno-stringop-overflow" .. && \
    make YDERFLAG=1 && \
    make install && \
    cd ../../ && rm -rf ulfius

# Install tomlc99
RUN git clone https://github.com/cktan/tomlc99.git && \
    cd tomlc99 && \
    make && \
    cp toml.h /usr/include/ && \
    cp libtoml.* /usr/lib/ && \
    cd .. && rm -rf tomlc99

WORKDIR /

# Clone Prometheus client repository and build it
RUN git clone https://github.com/ukama/prometheus-client.git && \
    cd prometheus-client && \
    mkdir -p prom/build && \
    cd prom/build && \
    cmake ../ -DCMAKE_C_COMPILER=gcc -DCMAKE_C_FLAGS="-I/usr/include" \
        -DCMAKE_LD_FLAGS="-L/usr/lib" && \
    make install DESTDIR=/usr && \
    mkdir -p ../../promhttp/build && \
    cd ../../promhttp/build && \
    cmake ../ -DCMAKE_C_COMPILER=gcc -DCMAKE_C_FLAGS="-I/usr/include" \
        -DCMAKE_LD_FLAGS="-L/usr/lib" && \
    make install DESTDIR=/usr && \
    cd ../../.. && \
    rm -rf prometheus-client

# Final stage
FROM builder as final

# Conditional logic for final stage
#RUN if echo "${TARGETPLATFORM}" | grep -q "ubuntu"; then \
#        apt-get update && apt-get install -y extlinux; \
#    elif echo "${TARGETPLATFORM}" | grep -q "alpine"; then \
#        echo "No additional packages for Alpine."; \
#    fi

# Optionally set the working directory again if it was overridden
WORKDIR /workspace