FROM alpine:3.14.0 
ENTRYPOINT []
WORKDIR /

# Prepare FS for ukamaDM
COPY container/root /
COPY container/start.sh /start.sh

# Copy device registry related bins
COPY build/bin/client /sbin/lwm2mClient

CMD [ "/start.sh" ]
