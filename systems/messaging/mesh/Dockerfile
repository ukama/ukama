FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install libprotobuf-c-dev librabbitmq-dev libmicrohttpd-dev pkg-config libcurl4-openssl-dev libjansson-dev libmicrohttpd-dev pkg-config libcurl4-openssl-dev libjansson-dev -y

COPY mesh /usr/bin/mesh
RUN mkdir /config/ 
COPY configs/config.toml /config

CMD ["/usr/bin/mesh", "-c", "/config/config.toml"]
