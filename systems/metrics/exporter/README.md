# Protocol Documentation
<a name="top"></a>

## Table of Contents

 - [ExporterService](#ukama.metrics.exporter.v1.ExporterService)
- [Directory structure](#directory-structure)
- [How to use?](#how-to)
 - [exporter.proto](#exporter.proto)
   - [DummyParameter Messages](#ukama.metrics.exporter.v1.DummyParameter)
  
- [Scalar Value Types](#scalar-value-types)



<a name="exporter.proto"></a>
<p align="right"><a href="#top">Top</a></p>


<a name="ukama.metrics.exporter.v1.ExporterService"></a>

# ExporterService


## RPC Functions

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Dummy | [DummyParameter](#ukama.metrics.exporter.v1.DummyParameter) | [DummyParameter](#ukama.metrics.exporter.v1.DummyParameter) |  |




<a name="#directory-structure"></a>

## Directory structure



- **cmd**: Contains the server and system/sub-system version. Purpose of this file is to initialize the DB and start server. We use `make` command to build this service.
- **mocks**: This directory contains the auto generated file which get generated based on `*.proto`. It contains functions which we can use to write test cases.
- **pb**: This directory contains the `*.proto` file. In proto file we define service with all the rpc's and messages.
- **pkg/db**: DB directory under pkg contains 2 files.
  `model.go` file contains the db model structure/s.
  `*_repo.go` is reponsible of communicating with db using [gorm](https://gorm.io/docs/).
- **pkg/client** This dir contains the REST client interfaces to other system like factory, PCRF and Network. 
- **pkg/server** This dir contains the logic of RPC handlers.

<a name="#how-to"></a>

## How to use?

Before using the repo make sure below tools are installed:

- Go 1.18
- PostgreSQL
- gRPC client
Then navigate into base-rate directory and run below command:

**To Generate PB file**

```
make gen
```
This command will generate protobuf files from exporter.proto and mocks for the test.


**To Test**

For unit tests run below commands:

```
make test
```
This command will run unit tests under all exporter.proto directories.


**Build**

```
make
```

**Run**
```
./bin/exporter
```

## exporter.proto



<a name="ukama.metrics.exporter.v1.DummyParameter"></a>

### DummyParameter








 





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
