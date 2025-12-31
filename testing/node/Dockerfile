# ---------- Stage 1: builder (contains full repo) ----------
FROM ubuntu:22.04 AS builder
ARG DEBIAN_FRONTEND=noninteractive
ARG VNODE_DIR=testing/node

RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    g++ \
    make \
    cmake \
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
    python3 \
    python3-dev \
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
    && rm -rf /var/lib/apt/lists/*

WORKDIR /work

# Copy entire repo (build context = repo root)
COPY . /work/ukama

WORKDIR /work/ukama/${VNODE_DIR}
RUN make clean && make all


# ---------- Stage 2: runtime (small; no repo) ----------
FROM ubuntu:22.04
ARG DEBIAN_FRONTEND=noninteractive
ARG VNODE_DIR=testing/node

ENTRYPOINT []
WORKDIR /

# ---------- Stage 2: runtime (builds happen here; no repo) ----------
FROM ubuntu:22.04
ARG DEBIAN_FRONTEND=noninteractive
ARG VNODE_DIR=testing/node

ENTRYPOINT []
WORKDIR /

RUN apt-get update && apt-get install -y \
    software-properties-common \
 && add-apt-repository -y universe \
 && apt-get update && apt-get install -y \
    podman \
    fuse-overlayfs \
    buildah \
    ca-certificates \
    make \
    build-essential \
    pkg-config \
    cmake \
    autoconf \
    automake \
    libtool \
    git \
    curl \
    wget \
    rsync \
    tar \
    pigz \
    xz-utils \
    unzip \
    iproute2 \
    iputils-ping \
    openssh-client \
    tcpdump \
    ethtool \
    iperf3 \
    dosfstools \
    e2fsprogs \
    util-linux \
    tree \
    htop \
    tzdata \
    perl \
    vim \
    libjansson4 \
    libmicrohttpd12 \
    libidn2-0 \
    libcap2 \
    libsqlite3-0 \
    libssl3 \
    zlib1g \
    libcurl4 \
    libgnutls28-dev \
    libmicrohttpd-dev \
    libcurl4-openssl-dev \
    libidn2-dev \
    zlib1g-dev \
    libssl-dev \  
    libunistring-dev \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/lib/shared/overlay-images
RUN mkdir -p /var/lib/shared/overlay-layers
RUN touch /var/lib/shared/overlay-images/images.lock
RUN touch /var/lib/shared/overlay-layers/layers.lock

RUN mkdir -p /ukama
RUN mkdir -p /scripts

ENV _BUILDAH_STARTED_IN_USERNS=""
ENV BUILDAH_ISOLATION=chroot

ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/tmp/virtnode/ukamaOS/distro/vendor/build/lib:/tmp/virtnode/ukamaOS/distro/vendor/build/lib64:/tmp/virtnode/ukamaOS/distro/platform/build"

# Copy only runtime artifacts
COPY --from=builder /work/ukama/${VNODE_DIR}/virtualNode /sbin/virtualNode
COPY --from=builder /work/ukama/${VNODE_DIR}/apps /apps
COPY --from=builder /work/ukama/${VNODE_DIR}/scripts /scripts
COPY --from=builder /work/ukama/${VNODE_DIR}/supervisord.conf /supervisord.conf

CMD ["/sbin/virtualNode", "--target", "ubuntu:22.04", "--exec", "create", "--c", "/apps/"]
