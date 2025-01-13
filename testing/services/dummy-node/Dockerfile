FROM ubuntu:latest

RUN apt-get update && apt-get install -y supervisor && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/log/supervisor

WORKDIR /usr/bin/
COPY bin/ /usr/bin/

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

EXPOSE 8080 8081

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]