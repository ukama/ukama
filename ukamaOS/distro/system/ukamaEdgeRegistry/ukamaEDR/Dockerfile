FROM scratch 
ENTRYPOINT []
CMD []
WORKDIR /
# Prepare FS for ukamaEDR
COPY container/root /

# Copy device registry related bins
COPY build/ukamaEDR /sbin/ukamaEDR
COPY lib/ubsp/build/ubsp /sbin/testubsp
COPY lib/ubsp/build/schema /sbin/schema
COPY lib/ubsp/build/mfgutil /sbin/mfgutil
COPY lib/ubsp/utils/prepare_env.sh /sbin/prepare_env.sh

CMD ["ukamaEDR", "--p", "/etc/ukamaEDR/property.json", "--s", "/tmp/sys/cnode-systemdb", "--l", "DEBUG"]
