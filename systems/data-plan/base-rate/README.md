# Protocol Documentation
<a name="top"></a>

## Table of Contents

 - [BaseRatesService](#ukama.data_plan.rate.v1.BaseRatesService)
- [Directory structure](#directory-structure)
- [How to use?](#how-to)
 - [pb/rate.proto](#pb/rate.proto)
   - [GetBaseRateRequest Messages](#ukama.data_plan.rate.v1.GetBaseRateRequest)
   - [GetBaseRateResponse Messages](#ukama.data_plan.rate.v1.GetBaseRateResponse)
   - [GetBaseRatesRequest Messages](#ukama.data_plan.rate.v1.GetBaseRatesRequest)
   - [GetBaseRatesResponse Messages](#ukama.data_plan.rate.v1.GetBaseRatesResponse)
   - [Rate Messages](#ukama.data_plan.rate.v1.Rate)
   - [UploadBaseRatesRequest Messages](#ukama.data_plan.rate.v1.UploadBaseRatesRequest)
   - [UploadBaseRatesResponse Messages](#ukama.data_plan.rate.v1.UploadBaseRatesResponse)
  
- [Scalar Value Types](#scalar-value-types)



<a name="pb/rate.proto"></a>
<p align="right"><a href="#top">Top</a></p>


<a name="ukama.data_plan.rate.v1.BaseRatesService"></a>

# BaseRatesService
Define BaseRatesService service and its methods

## RPC Functions

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetBaseRates | [GetBaseRatesRequest](#ukama.data_plan.rate.v1.GetBaseRatesRequest) | [GetBaseRatesResponse](#ukama.data_plan.rate.v1.GetBaseRatesRequest) | Method to get base rates based on various parameters |
| GetBaseRate | [GetBaseRateRequest](#ukama.data_plan.rate.v1.GetBaseRateRequest) | [GetBaseRateResponse](#ukama.data_plan.rate.v1.GetBaseRateRequest) | Method to get base rate for a specific rate ID |
| UploadBaseRates | [UploadBaseRatesRequest](#ukama.data_plan.rate.v1.UploadBaseRatesRequest) | [UploadBaseRatesResponse](#ukama.data_plan.rate.v1.UploadBaseRatesRequest) | Method to upload base rates |




<a name="#directory-structure"></a>

## Directory structure

      ├── Dockerfile
      ├── Int.Dockerfile
      ├── Makefile
      ├── bin
      │      ├── base-rate
      │      ├── integration
      ├── cmd
      │      ├── server
      │      │      ├── main.go
      │      ├── version
      │      │      ├── version.go
      ├── go.mod
      ├── go.sum
      ├── mocks
      │      ├── BaseRateRepo.go
      ├── pb
      │      ├── gen
      │      │      ├── mocks
      │      │      │      ├── BaseRatesServiceClient.go
      │      │      │      ├── BaseRatesServiceServer.go
      │      │      │      ├── UnsafeBaseRatesServiceServer.go
      │      │      ├── rate.pb.go
      │      │      ├── rate.validator.pb.go
      │      │      ├── rate_grpc.pb.go
      │      ├── rate.proto
      ├── pkg
      │      ├── config.go
      │      ├── db
      │      │      ├── model.go
      │      │      ├── rate_repo.go
      │      │      ├── rate_repo_test.go
      │      ├── global.go
      │      ├── server
      │      │      ├── rate.go
      │      │      ├── rate_test.go
      │      ├── utils
      │      │      ├── utils.go
      │      │      ├── utils_test.go
      │      ├── validations
      │      │      ├── validations.go
      │      │      ├── validations_test.go
      ├── test
      │      ├── integration
      │      │      ├── rate_test.go

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

This command will run unit tests under all pb/rate.proto directories.

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from pb/rate.proto.

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

## pb/rate.proto
Define syntax and package name


<a name="ukama.data_plan.rate.v1.GetBaseRateRequest"></a>

### GetBaseRateRequest
Define GetBaseRateRequest message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rateID | [string](#string) |  | Rate ID to retrieve |




<a name="ukama.data_plan.rate.v1.GetBaseRateResponse"></a>

### GetBaseRateResponse
Define GetBaseRateResponse message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rate | [Rate](#ukama.data_plan.rate.v1.Rate) |  | Single rate |




<a name="ukama.data_plan.rate.v1.GetBaseRatesRequest"></a>

### GetBaseRatesRequest
Define GetBaseRatesRequest message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| country | [string](#string) |  | Country name |
| provider | [string](#string) |  | Provider name |
| to | [uint64](#uint64) |  | End time in Unix timestamp format |
| from | [uint64](#uint64) |  | Start time in Unix timestamp format |
| simType | [string](#string) |  | SIM type (e.g.unkown, ukama-data) |
| effectiveAt | [string](#string) |  | Effective date in "YYYY-MM-DD" format |




<a name="ukama.data_plan.rate.v1.GetBaseRatesResponse"></a>

### GetBaseRatesResponse
Define GetBaseRatesResponse message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rates | [Rate](#ukama.data_plan.rate.v1.Rate) | repeated | List of rates |




<a name="ukama.data_plan.rate.v1.Rate"></a>

### Rate
Define Rate message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rateID | [string](#string) |  | Rate ID |
| country | [string](#string) |  | Country name |
| network | [string](#string) |  | Network name |
| vpmn | [string](#string) |  | Virtual private mobile network (VPMN) name |
| imsi | [string](#string) |  | International Mobile Subscriber Identity (IMSI) |
| smsMo | [string](#string) |  | Short Message Service - Mobile Originated (SMS-MO) rate |
| smsMt | [string](#string) |  | Short Message Service - Mobile Terminated (SMS-MT) rate |
| data | [string](#string) |  | Data rate |
| _2g | [string](#string) |  | 2G rate |
| _3g | [string](#string) |  | 3G rate |
| _5g | [string](#string) |  | 5G rate |
| lte | [string](#string) |  | Long-Term Evolution (LTE) rate |
| lteM | [string](#string) |  | Machine-to-Machine (M2M) LTE rate |
| apn | [string](#string) |  | Access Point Name (APN) |
| createdAt | [string](#string) |  | Creation date and time in "YYYY-MM-DD HH:mm:ss" format |
| deletedAt | [string](#string) |  | Deletion date and time in "YYYY-MM-DD HH:mm:ss" format |
| updatedAt | [string](#string) |  | Update date and time in "YYYY-MM-DD HH:mm:ss" format |
| effectiveAt | [string](#string) |  | Effective |
| endAt | [string](#string) |  | endAt the date the rate will end |
| simType | [string](#string) |  | SIM type (e.g. unkown, ukama-data) |




<a name="ukama.data_plan.rate.v1.UploadBaseRatesRequest"></a>

### UploadBaseRatesRequest
Define UploadBaseRatesRequest message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fileURL | [string](#string) |  | URL of the file containing base rates data |
| effectiveAt | [string](#string) |  | Effective date in "YYYY-MM-DD" format |
| simType | [string](#string) |  | SIM type (e.g. unkown, ukama-data) |




<a name="ukama.data_plan.rate.v1.UploadBaseRatesResponse"></a>

### UploadBaseRatesResponse
Define UploadBaseRatesResponse message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rate | [Rate](#ukama.data_plan.rate.v1.Rate) | repeated | List of rates |







 





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
