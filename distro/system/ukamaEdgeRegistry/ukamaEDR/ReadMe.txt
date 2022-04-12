Build docker image:
docker build . -t unode:v0.0.1

Run docker image:
docker run --network host -p 5001:5001 -p 56830:56830 -p 5683:5683 -p 7001:7001 unode:v0.0.1

Run lwM2M Server:
./container/lwm2m/server -4 -l 5683



