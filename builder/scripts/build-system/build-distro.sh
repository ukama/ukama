#!/bin/ash

set -x

MAJOR_VERSION="v3.21"

# ===== Set up resolv.conf =====
echo "nameserver 8.8.8.8" > /etc/resolv.conf

# ===== Set Up Package Repositories =====
echo "https://dl-cdn.alpinelinux.org/alpine/${MAJOR_VERSION}/main" > /etc/apk/repositories
echo "https://dl-cdn.alpinelinux.org/alpine/${MAJOR_VERSION}/community" >> /etc/apk/repositories

# ===== Update Package Manager =====
apk update
apk upgrade

# ===== Install Essential System Packages =====
apk add alpine-base openrc busybox bash sudo shadow tzdata

apk add acpid busybox-openrc busybox-extras busybox-mdev-openrc

apk add readline bash autoconf automake libmicrohttpd-dev gnutls-dev openssl-dev iptables libuuid sqlite dhcpcd protobuf iproute2 zlib curl-dev nettle libcap libidn2   libmicrohttpd gnutls openssl-dev curl-dev  linux-headers bsd-compat-headers tree libtool sqlite-dev openssl-dev readline cmake autoconf automake alpine-sdk build-base git tcpdump ethtool iperf3 htop vim doas

#build apps
/ukama/builder/scripts/build-all-apps.sh /ukama
if [ $? -eq 0 ]; then
  echo "Apps build:"
  ls -ltr /ukama/build/pkgs/*

  echo "Package vendor libs"
  cd /ukama/nodes/ukamaOS/distro/vendor	
  ls -ltr /build/lib/*
  mkdir -p /ukama/build/libs
  tar -zcvf /ukama/build/libs/vendor_libs.tgz build/*

else 
  exit 1	
fi
