# Protocol Documentation

<a name="top"></a>

## Table of Contents

- [MsgClientService](#ukama.msgClient.v1.MsgClientService)
- [Directory structure](#directory-structure)
- [How to use?](#how-to)
- [pb/msgClient.proto](#pb/msgClient.proto)
  - [PublishMsgRequest Messages](#ukama.msgClient.v1.PublishMsgRequest)
  - [PublishMsgResponse Messages](#ukama.msgClient.v1.PublishMsgResponse)
  - [RegisterServiceReq Messages](#ukama.msgClient.v1.RegisterServiceReq)
  - [RegisterServiceResp Messages](#ukama.msgClient.v1.RegisterServiceResp)
  - [StartMsgBusHandlerReq Messages](#ukama.msgClient.v1.StartMsgBusHandlerReq)
  - [StartMsgBusHandlerResp Messages](#ukama.msgClient.v1.StartMsgBusHandlerResp)
  - [StopMsgBusHandlerReq Messages](#ukama.msgClient.v1.StopMsgBusHandlerReq)
  - [StopMsgBusHandlerResp Messages](#ukama.msgClient.v1.StopMsgBusHandlerResp)
  - [UnregisterServiceReq Messages](#ukama.msgClient.v1.UnregisterServiceReq)
  - [UnregisterServiceResp Messages](#ukama.msgClient.v1.UnregisterServiceResp)
  - [REGISTRAION_STATUS](#ukama.msgClient.v1.REGISTRAION_STATUS)
- [Scalar Value Types](#scalar-value-types)

<a name="pb/msgClient.proto"></a>
<p align="right"><a href="#top">Top</a></p>

<a name="ukama.msgClient.v1.MsgClientService"></a>

# MsgClientService

## RPC Functions

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RegisterService | [RegisterServiceReq](#ukama.msgClient.v1.RegisterServiceReq) | [RegisterServiceResp](#ukama.msgClient.v1.RegisterServiceReq) | Use this rpc to register system to MsgClient |
| StartMsgBusHandler | [StartMsgBusHandlerReq](#ukama.msgClient.v1.StartMsgBusHandlerReq) | [StartMsgBusHandlerResp](#ukama.msgClient.v1.StartMsgBusHandlerReq) | Call this rpc to StartMsgBus after registration |
| StopMsgBusHandler | [StopMsgBusHandlerReq](#ukama.msgClient.v1.StopMsgBusHandlerReq) | [StopMsgBusHandlerResp](#ukama.msgClient.v1.StopMsgBusHandlerReq) | Call this rpc to StopMsgBus |
| UnregisterService | [UnregisterServiceReq](#ukama.msgClient.v1.UnregisterServiceReq) | [UnregisterServiceResp](#ukama.msgClient.v1.UnregisterServiceReq) | Unregister service from MsgClient |
| PublishMsg | [PublishMsgRequest](#ukama.msgClient.v1.PublishMsgRequest) | [PublishMsgResponse](#ukama.msgClient.v1.PublishMsgRequest) | Call this rpc to publisg events |

<a name="#directory-structure"></a>

## Directory structure

      ├── Dockerfile
      ├── Int.Dockerfile
      ├── Makefile
      ├── README.md
      ├── cmd
      │      ├── msgClient
      │      │      ├── main.go
      │      ├── version
      │      │      ├── version.go
      ├── go.mod
      ├── go.sum
      ├── internal
      │      ├── config.go
      │      ├── db
      │      │      ├── model.go
      │      │      ├── route_repo.go
      │      │      ├── route_repo_test.go
      │      │      ├── service_repo.go
      │      │      ├── service_repo_test.go
      │      ├── global.go
      │      ├── queue
      │      │      ├── listener.go
      │      │      ├── listener_test.go
      │      │      ├── msgbus.go
      │      │      ├── publisher.go
      │      │      ├── publisher_test.go
      │      ├── server
      │      │      ├── msgClient.go
      │      │      ├── msgClient_test.go
      ├── mocks
      │      ├── MsgBusHandlerInterface.go
      │      ├── RouteRepo.go
      │      ├── ServiceRepo.go
      ├── pb
      │      ├── gen
      │      │      ├── mocks
      │      │      │      ├── MsgClientServiceClient.go
      │      │      │      ├── MsgClientServiceServer.go
      │      │      │      ├── UnsafeMsgClientServiceServer.go
      │      │      ├── msgClient.pb.go
      │      │      ├── msgClient.validator.pb.go
      │      │      ├── msgClient_grpc.pb.go
      │      ├── msgClient.proto
      ├── template.tmpl
      ├── test
      │      ├── integration
      │      │      ├── msgclient_test.go

- **cmd**: Contains the server and system/sub-system version. Purpose of this file is to initialize the DB and start server. We use `make server` command to run this file.
- **mocks**: This directory contains the auto generated file which get generated based on `*.proto`. It contains functions which we can use to write test cases.
- **pb**: This directory contains the `*.proto` file. In proto file we define service with all the rpc's and messages.
- **internal/db**: DB directory under pkg contains 2 files.
  `model.go` file contains the db model structure/s.
  `*_repo.go` is reponsible of communicating with db using [gorm](https://gorm.io/docs/).
- **internal/queue** This dir contains the logic of Queue's. like: Listner queue, Publisher queue.
- **internal/server** This dir contains the logic of RPC handlers.

<a name="#how-to"></a>

## How to use?

Before using the repo make sure below tools are installed:

- Go 1.18
- PostgreSQL
- gRPC client
Then navigate into base-rate directory and run below command:

**To Test**

For unit tests run below commands:

```
make test
```

This command will run unit tests under all pb/msgClient.proto directories.

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from pb/msgClient.proto.

**Make sure rabbitMq is running on docker**

```
docker run -d --name rabbit -p 5672:5672 -p 5673:5673 -p 15672:15672 rabbitmq:3-management
```

**To Run Server**

```
make server
```

## pb/msgClient.proto

<a name="ukama.msgClient.v1.PublishMsgRequest"></a>

### PublishMsgRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceUuid | [string](#string) |  | Uuid of service |
| routingKey | [string](#string) |  | Unique routing key |
| msg | [google.protobuf.Any](#google.protobuf.Any) |  | Msg proto |

<a name="ukama.msgClient.v1.PublishMsgResponse"></a>

### PublishMsgResponse

<a name="ukama.msgClient.v1.RegisterServiceReq"></a>

### RegisterServiceReq

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| systemName | [string](#string) |  |  |
| serviceName | [string](#string) |  |  |
| instanceId | [string](#string) |  |  |
| msgBusURI | [string](#string) |  |  |
| serviceURI | [string](#string) |  |  |
| listQueue | [string](#string) |  |  |
| publQueue | [string](#string) |  |  |
| exchange | [string](#string) |  |  |
| grpcTimeout | [uint32](#uint32) |  |  |
| routes | [string](#string) | repeated |  |

<a name="ukama.msgClient.v1.RegisterServiceResp"></a>

### RegisterServiceResp

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [REGISTRAION_STATUS](#ukama.msgClient.v1.REGISTRAION_STATUS) |  |  |
| serviceUuid | [string](#string) |  |  |

<a name="ukama.msgClient.v1.StartMsgBusHandlerReq"></a>

### StartMsgBusHandlerReq

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceUuid | [string](#string) |  | Uuid of service to start |

<a name="ukama.msgClient.v1.StartMsgBusHandlerResp"></a>

### StartMsgBusHandlerResp

<a name="ukama.msgClient.v1.StopMsgBusHandlerReq"></a>

### StopMsgBusHandlerReq

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceUuid | [string](#string) |  | Uuid of service to stop |

<a name="ukama.msgClient.v1.StopMsgBusHandlerResp"></a>

### StopMsgBusHandlerResp

<a name="ukama.msgClient.v1.UnregisterServiceReq"></a>

### UnregisterServiceReq

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceUuid | [string](#string) |  | Uuid of service to unregister |

<a name="ukama.msgClient.v1.UnregisterServiceResp"></a>

### UnregisterServiceResp

<a name="ukama.msgClient.v1.REGISTRAION_STATUS"></a>

### REGISTRAION_STATUS

Registration status enums

| Name | Number | Description |
| ---- | ------ | ----------- |
| REGISTERED | 0 | System registered status |
| ALLREADY_REGISTERED | 1 | System already registered status |
| NOT_REGISTERED | 2 | System not registered status |
| LISTNENING | 3 | Listening to event |
| LISTNENING_FAILURE | 4 | Listening failed |

## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" ></a> double |  | double | double | float |
| <a name="float" ></a> float |  | float | float | float |
| <a name="int32" ></a> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" ></a> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" ></a> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" ></a> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" ></a> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" ></a> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" ></a> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" ></a> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" ></a> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" ></a> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" ></a> bool |  | bool | boolean | boolean |
| <a name="string" ></a> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" ></a> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |
