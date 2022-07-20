FROM alpine:3.13.5

COPY mesh /usr/bin/mesh
RUN mkdir /config/ 
COPY config/config.toml /config

CMD ["/usr/bin/mesh -f /config/config.toml"]
