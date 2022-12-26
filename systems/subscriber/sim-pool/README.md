# Protocol Documentation
<a name="top"></a>

## Table of Contents

 - [SimService](#ukama.sim.v1.SimService)
- [Directory structure](#directory-structure)
- [How to use?](#how-to)
 - [pb/sim.proto](#pb/sim.proto)
   - [AddRequest Messages](#ukama.sim.v1.AddRequest)
   - [AddResponse Messages](#ukama.sim.v1.AddResponse)
   - [AddSim Messages](#ukama.sim.v1.AddSim)
   - [DeleteRequest Messages](#ukama.sim.v1.DeleteRequest)
   - [DeleteResponse Messages](#ukama.sim.v1.DeleteResponse)
   - [GetRequest Messages](#ukama.sim.v1.GetRequest)
   - [GetResponse Messages](#ukama.sim.v1.GetResponse)
   - [GetStatsRequest Messages](#ukama.sim.v1.GetStatsRequest)
   - [GetStatsResponse Messages](#ukama.sim.v1.GetStatsResponse)
   - [Sim Messages](#ukama.sim.v1.Sim)
   - [UploadRequest Messages](#ukama.sim.v1.UploadRequest)
   - [UploadResponse Messages](#ukama.sim.v1.UploadResponse)
    - [SimType](#ukama.sim.v1.SimType)
- [Scalar Value Types](#scalar-value-types)



<a name="pb/sim.proto"></a>
<p align="right"><a href="#top">Top</a></p>


<a name="ukama.sim.v1.SimService"></a>

# SimService
Sim pool sub-system is responsible of:

- Populating sims data in DB from CSV
- Provide sim stats
- Provide sim on request
- Allows to add slice of sims

## RPC Functions

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Get | [GetRequest](#ukama.sim.v1.GetRequest) | [GetResponse](#ukama.sim.v1.GetRequest) |  |
| GetStats | [GetStatsRequest](#ukama.sim.v1.GetStatsRequest) | [GetStatsResponse](#ukama.sim.v1.GetStatsRequest) |  |
| Add | [AddRequest](#ukama.sim.v1.AddRequest) | [AddResponse](#ukama.sim.v1.AddRequest) |  |
| Delete | [DeleteRequest](#ukama.sim.v1.DeleteRequest) | [DeleteResponse](#ukama.sim.v1.DeleteRequest) |  |
| Upload | [UploadRequest](#ukama.sim.v1.UploadRequest) | [UploadResponse](#ukama.sim.v1.UploadRequest) |  |




<a name="#directory-structure"></a>

## Directory structure

      ├── Dockerfile
      ├── Int.Dockerfile
      ├── Makefile
      ├── README.md
      ├── bin
      │      ├── sim-pool
      ├── cmd
      │      ├── server
      │      │      ├── main.go
      │      ├── version
      │      │      ├── version.go
      ├── go.mod
      ├── go.sum
      ├── mocks
      │      ├── SimRepo.go
      ├── pb
      │      ├── gen
      │      │      ├── mocks
      │      │      │      ├── SimServiceClient.go
      │      │      │      ├── SimServiceServer.go
      │      │      │      ├── UnsafeSimServiceServer.go
      │      │      ├── sim.pb.go
      │      │      ├── sim.validator.pb.go
      │      │      ├── sim_grpc.pb.go
      │      ├── sim.proto
      ├── pkg
      │      ├── config.go
      │      ├── db
      │      │      ├── model.go
      │      │      ├── sim_repo.go
      │      │      ├── sim_repo_test.go
      │      ├── global.go
      │      ├── server
      │      │      ├── sim.go
      │      │      ├── sim_test.go
      │      ├── utils
      │      │      ├── utils.go
      │      │      ├── utils_test.go

- **cmd**: Contains the server and system/sub-system version. Purpose of this file is to initialize the DB and start server. We use `make server` command to run this file.
- **mocks**: This directory contains the auto generated file which get generated based on `*.proto`. It contains functions which we can use to write test cases.
- **pb**: This directory contains the `*.proto` file. In proto file we define service with all the rpc's and messages.
- **pkg/db**: DB directory under pkg contains 2 files.
  `model.go` file contains the db model structure/s.
  `*_repo.go` is reponsible of communicating with db using [gorm](https://gorm.io/docs/).
- **pkg/server** This directory contains the file in which all the RPC functions logic is implemented. Those functions call `pkg\*_repo.go` functions to perform db operations.

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

This command will run unit tests under all pb/sim.proto directories.

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from pb/sim.proto.

**To Run Server & Test RPC**

```
make server
```

This command will run the server on port `9090`. It'll also create the database and table under it.

Server is running, Now we can use any gRPC client to interact with RPC handlers. We're using [Evans](https://github.com/ktr0731/evans). Run below command in new terminal tab:

```
evans --path /path/to --path . --proto pb/*.proto --host localhost --port 9090
```

Next run:

```
show rpc
```

This command will show all the available RPC calls under base-rate sub-system. To call any RPC function run `call FUNCATION_NAME`.

## pb/sim.proto



<a name="ukama.sim.v1.AddRequest"></a>

### AddRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sim | [AddSim](#ukama.sim.v1.AddSim) | repeated |  |




<a name="ukama.sim.v1.AddResponse"></a>

### AddResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sim | [Sim](#ukama.sim.v1.Sim) | repeated |  |




<a name="ukama.sim.v1.AddSim"></a>

### AddSim



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| iccid | [string](#string) |  |  |
| simType | [SimType](#ukama.sim.v1.SimType) |  |  |
| msisdn | [string](#string) |  |  |
| smDpAddress | [string](#string) |  |  |
| activationCode | [string](#string) |  |  |
| qrCode | [string](#string) |  |  |
| isPhysicalSim | [bool](#bool) |  |  |




<a name="ukama.sim.v1.DeleteRequest"></a>

### DeleteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) | repeated |  |




<a name="ukama.sim.v1.DeleteResponse"></a>

### DeleteResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) | repeated |  |




<a name="ukama.sim.v1.GetRequest"></a>

### GetRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| IsPhysicalSim | [bool](#bool) |  |  |




<a name="ukama.sim.v1.GetResponse"></a>

### GetResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sim | [Sim](#ukama.sim.v1.Sim) |  |  |




<a name="ukama.sim.v1.GetStatsRequest"></a>

### GetStatsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| simType | [SimType](#ukama.sim.v1.SimType) |  |  |




<a name="ukama.sim.v1.GetStatsResponse"></a>

### GetStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [uint64](#uint64) |  |  |
| available | [uint64](#uint64) |  |  |
| consumed | [uint64](#uint64) |  |  |
| failed | [uint64](#uint64) |  |  |




<a name="ukama.sim.v1.Sim"></a>

### Sim



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| iccid | [string](#string) |  |  |
| msisdn | [string](#string) |  |  |
| isAllocated | [bool](#bool) |  |  |
| simType | [SimType](#ukama.sim.v1.SimType) |  |  |
| smDpAddress | [string](#string) |  |  |
| activationCode | [string](#string) |  |  |
| created_at | [string](#string) |  |  |
| deleted_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |
| isPhysicalSim | [bool](#bool) |  |  |
| qrCode | [string](#string) |  |  |




<a name="ukama.sim.v1.UploadRequest"></a>

### UploadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| simData | [bytes](#bytes) |  |  |
| simType | [SimType](#ukama.sim.v1.SimType) |  |  |




<a name="ukama.sim.v1.UploadResponse"></a>

### UploadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sim | [Sim](#ukama.sim.v1.Sim) | repeated |  |






<a name="ukama.sim.v1.SimType"></a>

### SimType


| Name | Number | Description |
| ---- | ------ | ----------- |
| ANY | 0 |  |
| INTER_NONE | 1 |  |
| INTER_MNO_DATA | 2 |  |
| INTER_MNO_ALL | 3 |  |
| INTER_UKAMA_ALL | 4 |  |




 





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
