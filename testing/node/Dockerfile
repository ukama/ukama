FROM ubuntu:24.04
ENTRYPOINT []

WORKDIR /

ARG DEBIAN_FRONTEND=noninteractive

# Enable Universe + install everything in one go
RUN apt-get update && apt-get install -y --no-install-recommends \
        software-properties-common \
    && add-apt-repository -y universe \
    && apt-get update && apt-get install -y --no-install-recommends \
        build-essential \
        gcc-arm-linux-gnueabihf \
        g++ \
        cmake \
        make \
        automake \
        autoconf \
        libtool \
        pkg-config \
        bison \
        flex \
        gawk \
        bc \
        pigz \
        lzop \
        unzip \
        patchelf \
        libncurses-dev \
        libssl-dev \
        libelf-dev \
        zlib1g-dev \
        libisl-dev \
        libjansson-dev \
        libmicrohttpd-dev \
        libidn2-dev \
        libcap-dev \
        uuid-dev \
        libprotobuf-dev \
        protobuf-compiler \
        sqlite3 \
        libsqlite3-dev \
        python3 \
        python3-dev \
        python3-venv \
        curl \
        wget \
        rsync \
        iproute2 \
        iputils-ping \
        openssh-client \
        tcpdump \
        ethtool \
        iperf3 \
        dosfstools \
        e2fsprogs \
        util-linux \
        git \
        buildah \
        tree \
        htop \
        tzdata \
        perl \
        gnat \
        tcl \
        vim \
        libunistring-dev \
        libcurl4-openssl-dev \
        libgnutls28-dev \
        nettle-dev \
        libp11-kit-dev \
        libtasn1-6-dev \
        podman \
        fuse-overlayfs \
        uidmap \
        slirp4netns \
    && rm -rf /var/lib/apt/lists/*

# Buildah/Podman storage placeholders (your original behavior)
RUN mkdir -p /var/lib/shared/overlay-images /var/lib/shared/overlay-layers \
    && touch /var/lib/shared/overlay-images/images.lock \
    && touch /var/lib/shared/overlay-layers/layers.lock

RUN mkdir -p /ukama /scripts

ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/tmp/virtnode/ukamaOS/distro/vendor/build/lib:/tmp/virtnode/ukamaOS/distro/vendor/build/lib64:/tmp/virtnode/ukamaOS/distro/platform/build"
ENV _BUILDAH_STARTED_IN_USERNS="" \
    BUILDAH_ISOLATION=chroot

# Copy virtual node related bin and config
COPY virtualNode /sbin/
COPY apps /apps
COPY ./scripts/*.sh /scripts/
COPY supervisord.conf /

CMD ["/sbin/virtualNode", "--target", "ubuntu:24.04", "--exec", "create", "--c", "/apps/"]
