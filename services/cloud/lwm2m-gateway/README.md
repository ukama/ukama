# LwM2M Gateway

LwM2M Gateway receives meassges from Controller and translates those to LwM2M objects.

## Build

`make all`

## Starting service
LwM2M gateway can be configured by using a JSON config. Sample config is present in configs/lwm2m-gateway.json.

This config can be used either by copying to /etc/config on target or by using -cfgPath.

`./bin/lwm2mGateway`

## Test

### Start LwM2M Server

Start lwM2m Server at "0.0.0.0:3000"
[LwM2M server and client](https://github.com/ukama/ukamaOS/tree/main/distro/system/ukamaDM/dm)

### Using Test mock

``` Usage
Usage: testService <Command> <ARGS> 
Commands:
         READ_CONFIG 
         UPDATE_CONFIG 
         EXEC
```

Example

```
./bin/testService READ_CONFIG
```
## Preparing Container Image

Before building the container image for LwM2M gateway make sure build is done.

```
docker build . -t lwm2mgateway:v0.0.1
```

## Starting Container Image

```
docker run --network host -p 3100:3100 lwm2mgateway:v0.0.1
```