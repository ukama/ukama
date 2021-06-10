FROM ubuntu:18.04

# Install supervisord
RUN apt-get update && apt-get install -y supervisor vim htop
RUN mkdir -p /var/log/supervisor
COPY container/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Prepare FS for ukamaEDR
RUN mkdir -p /etc/certs
RUN mkdir -p /etc/ip
RUN mkdir -p /etc/ukamaEDR/schema
RUN mkdir -p /etc/ukamaEDR/config 
RUN mkdir -p /tmp/json/
RUN mkdir /history

# Copy device registry related bins
COPY build/ukamaEDR /sbin/ukamaEDR
COPY lib/ubsp/mfgdata/property/property.json /etc/ukamaEDR
COPY lib/ubsp/build/ubsp /sbin/testubsp
COPY lib/ubsp/build/schema /sbin/schema
COPY lib/ubsp/build/mfgutil /sbin/mfgutil
COPY lib/ubsp/mfgdata/schema/*.json /etc/ukamaEDR/schema/
COPY lib/ubsp/utils/prepare_env.sh /sbin/prepare_env.sh

# Copy the device File system
COPY container/sys /tmp/sys

#lwm2m Client
COPY container/lwm2m/client /sbin/lwm2mClient
#COPY container/lwm2m/client.toml /etc/ukamaEDR/config/client.toml

#LTE mediator, acsMifClient(stub), random metric genrator
COPY container/lte/* /sbin/

CMD ["/usr/bin/supervisord","-c", "/etc/supervisor/conf.d/supervisord.conf"]
