# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [Protocol Documentation](#protocol-documentation)
  - [Table of Contents](#table-of-contents)
- [ProfileService](#profileservice)
  - [RPC Functions](#rpc-functions)
  - [Directory structure](#directory-structure)
  - [How to use?](#how-to-use)
  - [profile.proto](#profileproto)
    - [AddReq](#addreq)
    - [AddResp](#addresp)
    - [Apn](#apn)
    - [Package](#package)
    - [Profile](#profile)
    - [ReadReq](#readreq)
    - [ReadResp](#readresp)
    - [RemoveReq](#removereq)
    - [RemoveResp](#removeresp)
    - [SyncReq](#syncreq)
    - [SyncResp](#syncresp)
    - [UpdatePackageReq](#updatepackagereq)
    - [UpdatePackageResp](#updatepackageresp)
    - [UpdateUsageReq](#updateusagereq)
    - [UpdateUsageResp](#updateusageresp)
  - [Scalar Value Types](#scalar-value-types)



<a name="profile.proto"></a>
<p align="right"><a href="#top">Top</a></p>


<a name="ukama.subscriber.profile.v1.ProfileService"></a>

# ProfileService
Profile registry maintains the record of all the active subscribers of a organization and thier config based on the active package id.
Profile also provide the input to policy control for enforcing the policies.
Initial implementation of Profile service itself will handle the policy control.   

RPC exposed by Profile
- Add
- Remove
- UpdatePackage
- UpdateUsage
- Sync
- Read

## RPC Functions

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Add | [AddReq](#ukama.subscriber.profile.v1.AddReq) | [AddResp](#ukama.subscriber.profile.v1.AddReq) | Use this RPC to add a new subscriber profile |
| Remove | [RemoveReq](#ukama.subscriber.profile.v1.RemoveReq) | [RemoveResp](#ukama.subscriber.profile.v1.RemoveReq) | Use this RPC to remove a subscriber profile |
| UpdatePackage | [UpdatePackageReq](#ukama.subscriber.profile.v1.UpdatePackageReq) | [UpdatePackageResp](#ukama.subscriber.profile.v1.UpdatePackageReq) | Use this RPC to update a active package of the subscriber |
| UpdateUsage | [UpdateUsageReq](#ukama.subscriber.profile.v1.UpdateUsageReq) | [UpdateUsageResp](#ukama.subscriber.profile.v1.UpdateUsageReq) | Use this RPC to update a usage of the subscriber |
| Sync | [SyncReq](#ukama.subscriber.profile.v1.SyncReq) | [SyncResp](#ukama.subscriber.profile.v1.SyncReq) | Use this RPC to sync records from cloud to node profile registry |
| Read | [ReadReq](#ukama.subscriber.profile.v1.ReadReq) | [ReadResp](#ukama.subscriber.profile.v1.ReadReq) | This RPC is used to read the active profiles |




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
This command will generate protobuf files from profile.proto and mocks for the test.


**To Test**

For unit tests run below commands:

```
make test
```
This command will run unit tests under all profile.proto directories.


**Build**

```
make
```

**Run**
```
./bin/profile
```

## profile.proto



<a name="ukama.subscriber.profile.v1.AddReq"></a>

### AddReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Profile | [Profile](#ukama.subscriber.profile.v1.Profile) |  |  |




<a name="ukama.subscriber.profile.v1.AddResp"></a>

### AddResp
Empty




<a name="ukama.subscriber.profile.v1.Apn"></a>

### Apn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Name | [string](#string) |  |  |




<a name="ukama.subscriber.profile.v1.Package"></a>

### Package



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| UeDlBps | [uint64](#uint64) |  |  |
| UeUlBps | [uint64](#uint64) |  |  |
| apn | [Apn](#ukama.subscriber.profile.v1.Apn) |  |  |
| PackageId | [string](#string) |  |  |
| AllowedTimeOfService | [int64](#int64) |  |  |
| TotalDataBytes | [uint64](#uint64) |  |  |
| ConsumedDataBytes | [uint64](#uint64) |  |  |




<a name="ukama.subscriber.profile.v1.Profile"></a>

### Profile



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Imsi | [string](#string) |  |  |
| Iccid | [string](#string) |  |  |
| UeDlBps | [uint64](#uint64) |  |  |
| UeUlBps | [uint64](#uint64) |  |  |
| Apn | [Apn](#ukama.subscriber.profile.v1.Apn) |  |  |
| NetworkId | [string](#string) |  |  |
| PackageId | [string](#string) |  |  |
| AllowedTimeOfService | [int64](#int64) |  |  |
| TotalDataBytes | [uint64](#uint64) |  |  |
| ConsumedDataBytes | [uint64](#uint64) |  |  |
| UpdatedAt | [int64](#int64) |  |  |
| LastChange | [string](#string) |  |  |
| LastChangeAt | [int64](#int64) |  |  |




<a name="ukama.subscriber.profile.v1.ReadReq"></a>

### ReadReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Imsi | [string](#string) |  | oneof |
| Iccid | [string](#string) |  | oneof |




<a name="ukama.subscriber.profile.v1.ReadResp"></a>

### ReadResp



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Profile | [Profile](#ukama.subscriber.profile.v1.Profile) |  |  |




<a name="ukama.subscriber.profile.v1.RemoveReq"></a>

### RemoveReq
Could be called by subscriber manager with ICCID and by billing service with imsi


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Imsi | [string](#string) |  | oneof |
| Iccid | [string](#string) |  | oneof |




<a name="ukama.subscriber.profile.v1.RemoveResp"></a>

### RemoveResp
Empty




<a name="ukama.subscriber.profile.v1.SyncReq"></a>

### SyncReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Iccid | [string](#string) | repeated |  |




<a name="ukama.subscriber.profile.v1.SyncResp"></a>

### SyncResp



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Iccid | [string](#string) | repeated |  |




<a name="ukama.subscriber.profile.v1.UpdatePackageReq"></a>

### UpdatePackageReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Iccid | [string](#string) |  |  |
| Package | [Package](#ukama.subscriber.profile.v1.Package) |  |  |




<a name="ukama.subscriber.profile.v1.UpdatePackageResp"></a>

### UpdatePackageResp
Empty




<a name="ukama.subscriber.profile.v1.UpdateUsageReq"></a>

### UpdateUsageReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Imsi | [string](#string) |  |  |
| ConsumedDataBytes | [uint64](#uint64) |  |  |




<a name="ukama.subscriber.profile.v1.UpdateUsageResp"></a>

### UpdateUsageResp








 





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
