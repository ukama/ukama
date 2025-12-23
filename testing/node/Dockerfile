FROM ubuntu:22.04
ENTRYPOINT []

WORKDIR /

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update

RUN apt install software-properties-common -y

RUN add-apt-repository ppa:deadsnakes/ppa -y

RUN apt-get update && apt-get install -y \
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
    libncurses5-dev \
    libssl-dev \
    libelf-dev \
    zlib1g-dev \
    libisl-dev \
    libjansson-dev \
    libmicrohttpd-dev \
    libidn2-dev \
    libcap-dev \
    libuuid1 \
    libprotobuf-dev \
    sqlite3 \
    libsqlite3-dev \
    python3.9 \
    python3.9-dev \
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
    libmicrohttpd-dev \
    libjansson-dev \
    libcurl4-openssl-dev \
    libgnutls28-dev \
    libsqlite3-dev \
    \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y software-properties-common \
 && add-apt-repository -y universe \
 && apt-get update \
 && apt-get install -y podman fuse-overlayfs \
 && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/lib/shared/overlay-images /var/lib/shared/overlay-layers; touch /var/lib/shared/overlay-images/images.lock; touch /var/lib/shared/overlay-layers/layers.lock

RUN mkdir -p /ukama

RUN mkdir -p /scripts

ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/tmp/virtnode/ukamaOS/distro/vendor/build/lib:/tmp/virtnode/ukamaOS/distro/vendor/build/lib64:/tmp/virtnode/ukamaOS/distro/platform/build"

ENV _BUILDAH_STARTED_IN_USERNS="" BUILDAH_ISOLATION=chroot

# Copy virtual node related bin and config
COPY virtualNode /sbin/

COPY capps /capps

COPY ./scripts/*.sh /scripts/

COPY supervisord.conf /

# /sbin/virtualNode --target "ubuntu:22.04" --exec create --c "/capps"
CMD ["/sbin/virtualNode", "--target", "ubuntu:22.04", "--exec", "create", "--c", "/capps/"]
