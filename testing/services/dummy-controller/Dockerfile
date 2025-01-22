FROM alpine:latest

WORKDIR /app

COPY bin/dummy-controller .

ENV MOCK_BATTERY=true

EXPOSE 2112

CMD ["./dummy-controller"]