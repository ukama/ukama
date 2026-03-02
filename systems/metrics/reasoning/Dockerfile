FROM alpine:3.13.5

COPY bin/reasoning /usr/bin/reasoning
COPY metric-key-map.json /etc/reasoning/metric-key-map.json
COPY metric-rules.json /etc/reasoning/metric-rules.json

CMD ["/usr/bin/reasoning"]