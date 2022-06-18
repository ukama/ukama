FROM scratch 
ENTRYPOINT []
CMD []
WORKDIR /

# Copy virtual node related bin and config
COPY virtualNode /sbin/
COPY capps /capps

CMD ["/sbin/virtualNode", "--target", "alpine:latest", --exec "create", --c", "/capps/"]
