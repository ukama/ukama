# Protocol Documentation
<a name="top"></a>

## Table of Contents

 - [ASR](#ukama.subscriber.asr.v1.AsrRecordService)
- [Directory structure](#directory-structure)
- [How to use?](#how-to)
 - [asr.proto](#asr.proto)
   - [ActivateReq Messages](#ukama.subscriber.asr.v1.ActivateReq)
   - [ActivateResp Messages](#ukama.subscriber.asr.v1.ActivateResp)
   - [Apn Messages](#ukama.subscriber.asr.v1.Apn)
   - [Guti Messages](#ukama.subscriber.asr.v1.Guti)
   - [InactivateReq Messages](#ukama.subscriber.asr.v1.InactivateReq)
   - [InactivateResp Messages](#ukama.subscriber.asr.v1.InactivateResp)
   - [ReadReq Messages](#ukama.subscriber.asr.v1.ReadReq)
   - [ReadResp Messages](#ukama.subscriber.asr.v1.ReadResp)
   - [Record Messages](#ukama.subscriber.asr.v1.Record)
   - [UpdateGutiReq Messages](#ukama.subscriber.asr.v1.UpdateGutiReq)
   - [UpdateGutiResp Messages](#ukama.subscriber.asr.v1.UpdateGutiResp)
   - [UpdatePackageReq Messages](#ukama.subscriber.asr.v1.UpdatePackageReq)
   - [UpdatePackageResp Messages](#ukama.subscriber.asr.v1.UpdatePackageResp)
   - [UpdateTaiReq Messages](#ukama.subscriber.asr.v1.UpdateTaiReq)
   - [UpdateTaiResp Messages](#ukama.subscriber.asr.v1.UpdateTaiResp)
  
- [Scalar Value Types](#scalar-value-types)



<a name="asr.proto"></a>
<p align="right"><a href="#top">Top</a></p>


<a name="ukama.subscriber.asr.v1.AsrRecordService"></a>

# ASR aka Active Subscriber Registry
ASR maintains the record of all the active subscribers of a organization. All the network with in the organization share same ASR.
Subscriber is added to ASr as soon as its activated and removed from ASR as soon as its deactivated.

ASR has REST ineterfaces to the servics like 
-  Factory, for reading sim data
-  PCRF, for setting policies for subscriber
-  Organization registry for validating network.enum

For now subscriber can only be a part on one network under organization. If he needs to join other network a new sim needs to be allocated.

RPC exposed by ASR
- Activate
- Inactivate
- UpdatePackage
- UpdateGuti
- UpdateTai
- Read

## RPC Functions

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Activate | [ActivateReq](#ukama.subscriber.asr.v1.ActivateReq) | [ActivateResp](#ukama.subscriber.asr.v1.ActivateReq) | Use this RPC to activate or add a new subscriber to ASR |
| Inactivate | [InactivateReq](#ukama.subscriber.asr.v1.InactivateReq) | [InactivateResp](#ukama.subscriber.asr.v1.InactivateReq) | Use this RPC to inactivate or remove a subscriber to ASR |
| UpdatePackage | [UpdatePackageReq](#ukama.subscriber.asr.v1.UpdatePackageReq) | [UpdatePackageResp](#ukama.subscriber.asr.v1.UpdatePackageReq) | Use this RPC to update a subscriber package in ASR |
| UpdateGuti | [UpdateGutiReq](#ukama.subscriber.asr.v1.UpdateGutiReq) | [UpdateGutiResp](#ukama.subscriber.asr.v1.UpdateGutiReq) | This RPC is called when a Update GUTI message is sent by node |
| UpdateTai | [UpdateTaiReq](#ukama.subscriber.asr.v1.UpdateTaiReq) | [UpdateTaiResp](#ukama.subscriber.asr.v1.UpdateTaiReq) | This RPC is called when a Update TAI message is sent by node |
| Read | [ReadReq](#ukama.subscriber.asr.v1.ReadReq) | [ReadResp](#ukama.subscriber.asr.v1.ReadReq) | This RPC is used to read the subscriber data from ASR based on IMSI or ICCID |




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
This command will generate protobuf files from asr.proto and mocks for the test.


**To Test**

For unit tests run below commands:

```
make test
```
This command will run unit tests under all asr.proto directories.


**Build**

```
make
```

**Run**
```
./bin/asr
```

## asr.proto



<a name="ukama.subscriber.asr.v1.ActivateReq"></a>

### ActivateReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network | [string](#string) |  |  |
| Iccid | [string](#string) |  |  |
| PackageId | [string](#string) |  |  |




<a name="ukama.subscriber.asr.v1.ActivateResp"></a>

### ActivateResp
Empty




<a name="ukama.subscriber.asr.v1.Apn"></a>

### Apn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Name | [string](#string) |  |  |




<a name="ukama.subscriber.asr.v1.Guti"></a>

### Guti



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| PlmnId | [string](#string) |  |  |
| Mmegi | [uint32](#uint32) |  |  |
| Mmec | [uint32](#uint32) |  |  |
| Mtmsi | [uint32](#uint32) |  |  |




<a name="ukama.subscriber.asr.v1.InactivateReq"></a>

### InactivateReq
Could be called by subscriber manager with ICCID and by billing service with imsi


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Imsi | [string](#string) |  | oneof |
| Iccid | [string](#string) |  | oneof |




<a name="ukama.subscriber.asr.v1.InactivateResp"></a>

### InactivateResp
Empty




<a name="ukama.subscriber.asr.v1.ReadReq"></a>

### ReadReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Imsi | [string](#string) |  | oneof |
| Iccid | [string](#string) |  | oneof |




<a name="ukama.subscriber.asr.v1.ReadResp"></a>

### ReadResp



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Record | [Record](#ukama.subscriber.asr.v1.Record) |  |  |




<a name="ukama.subscriber.asr.v1.Record"></a>

### Record



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Imsi | [string](#string) |  |  |
| SimId | [string](#string) |  |  |
| Iccid | [string](#string) |  |  |
| Key | [bytes](#bytes) |  |  |
| Op | [bytes](#bytes) |  |  |
| Amf | [bytes](#bytes) |  |  |
| Apn | [Apn](#ukama.subscriber.asr.v1.Apn) |  |  |
| AlgoType | [uint32](#uint32) |  |  |
| UeDlAmbrBps | [uint32](#uint32) |  |  |
| UeUlAmbrBps | [uint32](#uint32) |  |  |
| Sqn | [uint64](#uint64) |  |  |
| CsgIdPrsent | [bool](#bool) |  |  |
| CsgId | [uint32](#uint32) |  |  |
| PackageId | [string](#string) |  |  |




<a name="ukama.subscriber.asr.v1.UpdateGutiReq"></a>

### UpdateGutiReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Imsi | [string](#string) |  |  |
| Guti | [Guti](#ukama.subscriber.asr.v1.Guti) |  |  |
| UpdatedAt | [uint32](#uint32) |  | unix timestamp |




<a name="ukama.subscriber.asr.v1.UpdateGutiResp"></a>

### UpdateGutiResp
Empty




<a name="ukama.subscriber.asr.v1.UpdatePackageReq"></a>

### UpdatePackageReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Iccid | [string](#string) |  |  |
| PackageId | [string](#string) |  |  |




<a name="ukama.subscriber.asr.v1.UpdatePackageResp"></a>

### UpdatePackageResp
Empty




<a name="ukama.subscriber.asr.v1.UpdateTaiReq"></a>

### UpdateTaiReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Imsi | [string](#string) |  |  |
| PlmnId | [string](#string) |  |  |
| Tac | [uint32](#uint32) |  | 16 bit max |
| UpdatedAt | [uint32](#uint32) |  | unix timestamp |




<a name="ukama.subscriber.asr.v1.UpdateTaiResp"></a>

### UpdateTaiResp
Empty







 





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
