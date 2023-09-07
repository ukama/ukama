# LwM2M 

## Build

Follow the instructions mentioned at [Build ukamaOS](https://github.com/ukama/ukamaOS#readme)

Make sure you have ukama tool chain build for OS. After that we can just do a make for Lwm2m without building complete ukamaOS
 
```
 make CC=<Path to musl-gcc>
```

## Preparing Container Image

Before building the container image for LwM2M server make sure build is done.

```
docker build . -t lwm2mserver:v0.0.1
```

## Starting Container Image

Server:

```
 docker run --network host -p 3000:3000 lwm2mserver:v0.0.1
```

Client:

```
docker run --network host -v ${PWD}/container/lwm2m/clientconf:/etc/lwclient lwm2mClient:v0.0.1
```
