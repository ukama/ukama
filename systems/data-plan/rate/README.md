# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [Protocol Documentation](#protocol-documentation)
  - [Table of Contents](#table-of-contents)
- [Rate Service](#rate-service)
  - [RPC Functions](#rpc-functions)
  - [Directory structure](#directory-structure)
  - [How to use?](#how-to-use)
  - [rate.proto](#rateproto)
    - [DeleteMarkupRequest](#deletemarkuprequest)
    - [DeleteMarkupResponse](#deletemarkupresponse)
    - [GetDefaultMarkupHistoryRequest](#getdefaultmarkuphistoryrequest)
    - [GetDefaultMarkupHistoryResponse](#getdefaultmarkuphistoryresponse)
    - [GetDefaultMarkupRequest](#getdefaultmarkuprequest)
    - [GetDefaultMarkupResponse](#getdefaultmarkupresponse)
    - [GetMarkupHistoryRequest](#getmarkuphistoryrequest)
    - [GetMarkupHistoryResponse](#getmarkuphistoryresponse)
    - [GetMarkupRequest](#getmarkuprequest)
    - [GetMarkupResponse](#getmarkupresponse)
    - [GetRateRequest](#getraterequest)
    - [GetRateResponse](#getrateresponse)
    - [GetRatesRequest](#getratesrequest)
    - [GetRatesResponse](#getratesresponse)
    - [MarkupRates](#markuprates)
    - [UpdateDefaultMarkupRequest](#updatedefaultmarkuprequest)
    - [UpdateDefaultMarkupResponse](#updatedefaultmarkupresponse)
    - [UpdateMarkupRequest](#updatemarkuprequest)
    - [UpdateMarkupResponse](#updatemarkupresponse)
  - [Scalar Value Types](#scalar-value-types)



<a name="rate.proto"></a>
<p align="right"><a href="#top">Top</a></p>


<a name="ukama.dataplan.rate.v1.RateService"></a>

# Rate Service
Rate service adds the markup on the base rates and provides the final applicable rates for the user.

If the user doesn't have any custom rates for him then the default markup is applied otherwise custom markup is considered


## RPC Functions

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetMarkup | [GetMarkupRequest](#ukama.dataplan.rate.v1.GetMarkupRequest) | [GetMarkupResponse](#ukama.dataplan.rate.v1.GetMarkupRequest) |  |
| UpdateMarkup | [UpdateMarkupRequest](#ukama.dataplan.rate.v1.UpdateMarkupRequest) | [UpdateMarkupResponse](#ukama.dataplan.rate.v1.UpdateMarkupRequest) |  |
| DeleteMarkup | [DeleteMarkupRequest](#ukama.dataplan.rate.v1.DeleteMarkupRequest) | [DeleteMarkupResponse](#ukama.dataplan.rate.v1.DeleteMarkupRequest) |  |
| GetMarkupHistory | [GetMarkupHistoryRequest](#ukama.dataplan.rate.v1.GetMarkupHistoryRequest) | [GetMarkupHistoryResponse](#ukama.dataplan.rate.v1.GetMarkupHistoryRequest) |  |
| GetDefaultMarkup | [GetDefaultMarkupRequest](#ukama.dataplan.rate.v1.GetDefaultMarkupRequest) | [GetDefaultMarkupResponse](#ukama.dataplan.rate.v1.GetDefaultMarkupRequest) |  |
| UpdateDefaultMarkup | [UpdateDefaultMarkupRequest](#ukama.dataplan.rate.v1.UpdateDefaultMarkupRequest) | [UpdateDefaultMarkupResponse](#ukama.dataplan.rate.v1.UpdateDefaultMarkupRequest) |  |
| GetDefaultMarkupHistory | [GetDefaultMarkupHistoryRequest](#ukama.dataplan.rate.v1.GetDefaultMarkupHistoryRequest) | [GetDefaultMarkupHistoryResponse](#ukama.dataplan.rate.v1.GetDefaultMarkupHistoryRequest) |  |
| GetRates | [GetRatesRequest](#ukama.dataplan.rate.v1.GetRatesRequest) | [GetRatesResponse](#ukama.dataplan.rate.v1.GetRatesRequest) |  |
| GetRate | [GetRateRequest](#ukama.dataplan.rate.v1.GetRateRequest) | [GetRateResponse](#ukama.dataplan.rate.v1.GetRateRequest) |  |




<a name="#directory-structure"></a>

## Directory structure



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

This command will run unit tests under all rate.proto directories.

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from rate.proto.

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

## rate.proto



<a name="ukama.dataplan.rate.v1.DeleteMarkupRequest"></a>

### DeleteMarkupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ownerId | [string](#string) |  |  |




<a name="ukama.dataplan.rate.v1.DeleteMarkupResponse"></a>

### DeleteMarkupResponse





<a name="ukama.dataplan.rate.v1.GetDefaultMarkupHistoryRequest"></a>

### GetDefaultMarkupHistoryRequest





<a name="ukama.dataplan.rate.v1.GetDefaultMarkupHistoryResponse"></a>

### GetDefaultMarkupHistoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| markupRates | [MarkupRates](#ukama.dataplan.rate.v1.MarkupRates) | repeated |  |




<a name="ukama.dataplan.rate.v1.GetDefaultMarkupRequest"></a>

### GetDefaultMarkupRequest





<a name="ukama.dataplan.rate.v1.GetDefaultMarkupResponse"></a>

### GetDefaultMarkupResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Markup | [double](#double) |  |  |




<a name="ukama.dataplan.rate.v1.GetMarkupHistoryRequest"></a>

### GetMarkupHistoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ownerId | [string](#string) |  |  |




<a name="ukama.dataplan.rate.v1.GetMarkupHistoryResponse"></a>

### GetMarkupHistoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ownerId | [string](#string) |  |  |
| markupRates | [MarkupRates](#ukama.dataplan.rate.v1.MarkupRates) | repeated |  |




<a name="ukama.dataplan.rate.v1.GetMarkupRequest"></a>

### GetMarkupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ownerId | [string](#string) |  |  |




<a name="ukama.dataplan.rate.v1.GetMarkupResponse"></a>

### GetMarkupResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ownerId | [string](#string) |  |  |
| Markup | [double](#double) |  |  |




<a name="ukama.dataplan.rate.v1.GetRateRequest"></a>

### GetRateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ownerId | [string](#string) |  |  |
| country | [string](#string) |  |  |
| provider | [string](#string) |  |  |
| to | [uint64](#uint64) |  |  |
| from | [uint64](#uint64) |  |  |
| simType | [string](#string) |  |  |
| effectiveAt | [string](#string) |  |  |




<a name="ukama.dataplan.rate.v1.GetRateResponse"></a>

### GetRateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rates | [ukama.dataplan.baserate.v1.Rate](#ukama.dataplan.baserate.v1.Rate) | repeated |  |




<a name="ukama.dataplan.rate.v1.GetRatesRequest"></a>

### GetRatesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| country | [string](#string) |  |  |
| provider | [string](#string) |  |  |
| to | [uint64](#uint64) |  |  |
| from | [uint64](#uint64) |  |  |
| simType | [string](#string) |  |  |
| effectiveAt | [string](#string) |  |  |




<a name="ukama.dataplan.rate.v1.GetRatesResponse"></a>

### GetRatesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rates | [ukama.dataplan.baserate.v1.Rate](#ukama.dataplan.baserate.v1.Rate) | repeated |  |




<a name="ukama.dataplan.rate.v1.MarkupRates"></a>

### MarkupRates



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| createdAt | [string](#string) |  |  |
| deletedAt | [string](#string) |  |  |
| Markup | [double](#double) |  |  |




<a name="ukama.dataplan.rate.v1.UpdateDefaultMarkupRequest"></a>

### UpdateDefaultMarkupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Markup | [double](#double) |  |  |




<a name="ukama.dataplan.rate.v1.UpdateDefaultMarkupResponse"></a>

### UpdateDefaultMarkupResponse





<a name="ukama.dataplan.rate.v1.UpdateMarkupRequest"></a>

### UpdateMarkupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ownerId | [string](#string) |  |  |
| Markup | [double](#double) |  |  |




<a name="ukama.dataplan.rate.v1.UpdateMarkupResponse"></a>

### UpdateMarkupResponse








 





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
