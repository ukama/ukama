FROM ubuntu:21.04
ENTRYPOINT []
CMD []
WORKDIR /

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update
RUN apt-get install bc build-essential git libncurses5-dev lzop perl libssl-dev gnat flex wget zlib1g-dev gcc-arm-linux-gnueabihf automake bison python libelf-dev cmake curl libtool tcl pkg-config tcl pkg-config autopoint wget libisl-dev g++ texinfo texlive ghostscript gawk libjansson-dev buildah patchelf vim -y

RUN mkdir -p /ukama

RUN mkdir -p /scripts

ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/tmp/virtnode/ukamaOS/distro/vendor/build/lib:/tmp/virtnode/ukamaOS/distro/platform/build"

ENV _BUILDAH_STARTED_IN_USERNS="" BUILDAH_ISOLATION=chroot

# Copy virtual node related bin and config
COPY virtualNode /sbin/

COPY capps /capps

COPY ukamaOS_*.tgz /ukama

COPY ./scripts/*.sh /scripts

COPY supervisord.conf /

CMD ["/sbin/virtualNode", "--target", "ubuntu:20.04:latest", "--exec", "create", "--c", "/capps/"]
