FROM ubuntu:22.04

#ENV DEBUGMODE=true
ENV PORT=8080
COPY initClient /usr/bin/initClient
COPY libs /usr/libs
RUN apt-get update

CMD ["/usr/bin/initClient"]
