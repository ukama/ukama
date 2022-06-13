# Notifyd [WIP]

Notification service running on the ukama node to report alarm and notification to consile apps.

## Architecture Diagram

Notifyd consist of a rest based client and server. Server for recieving events from the services and client to report these events to cloud servcie using meshd. Notifyd enacapsulate the event recived from the services and add information regarding nodes the event for easily identifying the node on cloud.


![notify.d](./docs/Notifyd.jpg)


## Build
```
make
```

## Test

Start a remote web server to receive alerts

```
TBU
``

Start noded service

```
./build/noded
```

Start NotifyD service

```
./build/notifyd
```

Post event to NotifyD:

```
curl --request POST \
  --url http://localhost:8085/notify/v1/event/core \
  --header 'content-type: application/json' \
  --data '{\n   "notify":{\n        "serviceName":"noded",\n        "time":1654566750,\n        "severity":"high",\n        "description":"UserOverload",\n     "reason": "Too Many users",\n       "details": "Error: user exceeding limit.",\n        "attribute": {\n            "name":"ActiveUser",\n          "units":"",\n           "value":64\n        },\n        "units":""\n    }\n}'
```

Post alert to NotifyD:

```
curl --request POST \
  --url http://localhost:8085/notify/v1/alert/core \
  --header 'content-type: application/json' \
  --data '{\n   "notify":{\n        "serviceName":"core",\n     "time":1654566750,\n        "severity":"high",\n        "description":"UserOverload",\n     "reason": "Too Many users",\n       "details": "Error: user exceeding limit.",\n        "attribute": {\n            "name":"ActiveUser",\n          "units":"",\n           "value":64\n        },\n        "units":""\n    }\n}'
```
