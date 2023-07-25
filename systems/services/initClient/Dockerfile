FROM ubuntu:22.04

#ENV DEBUGMODE=true
ENV PORT=8080
COPY initClient /usr/bin/initClient
COPY lib /usr/lib
RUN apt-get update

CMD ["/usr/bin/initClient"]
