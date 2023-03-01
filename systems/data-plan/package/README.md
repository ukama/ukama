# Protocol Documentation

<a name="top"></a>

## Table of Contents

- [PackagesService](#ukama.data_plan.package.v1.PackagesService)
- [Directory structure](#directory-structure)
- [How to use?](#how-to)
- [pb/package.proto](#pb/package.proto)
  - [AddPackageRequest Messages](#ukama.data_plan.package.v1.AddPackageRequest)
  - [AddPackageResponse Messages](#ukama.data_plan.package.v1.AddPackageResponse)
  - [DeletePackageRequest Messages](#ukama.data_plan.package.v1.DeletePackageRequest)
  - [DeletePackageResponse Messages](#ukama.data_plan.package.v1.DeletePackageResponse)
  - [GetByOrgPackageRequest Messages](#ukama.data_plan.package.v1.GetByOrgPackageRequest)
  - [GetByOrgPackageResponse Messages](#ukama.data_plan.package.v1.GetByOrgPackageResponse)
  - [GetPackageRequest Messages](#ukama.data_plan.package.v1.GetPackageRequest)
  - [GetPackageResponse Messages](#ukama.data_plan.package.v1.GetPackageResponse)
  - [Package Messages](#ukama.data_plan.package.v1.Package)
  - [UpdatePackageRequest Messages](#ukama.data_plan.package.v1.UpdatePackageRequest)
  - [UpdatePackageResponse Messages](#ukama.data_plan.package.v1.UpdatePackageResponse)
  
- [Scalar Value Types](#scalar-value-types)

<a name="pb/package.proto"></a>
<p align="right"><a href="#top">Top</a></p>

<a name="ukama.data_plan.package.v1.PackagesService"></a>

# PackagesService

Defines the service for packages  operations

## RPC Functions

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Get | [GetPackageRequest](#ukama.data_plan.package.v1.GetPackageRequest) | [GetPackageResponse](#ukama.data_plan.package.v1.GetPackageRequest) |  |
| Add | [AddPackageRequest](#ukama.data_plan.package.v1.AddPackageRequest) | [AddPackageResponse](#ukama.data_plan.package.v1.AddPackageRequest) |  |
| Delete | [DeletePackageRequest](#ukama.data_plan.package.v1.DeletePackageRequest) | [DeletePackageResponse](#ukama.data_plan.package.v1.DeletePackageRequest) |  |
| Update | [UpdatePackageRequest](#ukama.data_plan.package.v1.UpdatePackageRequest) | [UpdatePackageResponse](#ukama.data_plan.package.v1.UpdatePackageRequest) |  |
| GetByOrg | [GetByOrgPackageRequest](#ukama.data_plan.package.v1.GetByOrgPackageRequest) | [GetByOrgPackageResponse](#ukama.data_plan.package.v1.GetByOrgPackageRequest) |  |

<a name="#directory-structure"></a>

## Directory structure

      ├── Dockerfile
      ├── Int.Dockerfile
      ├── Makefile
      ├── bin
      │      ├── integration
      │      ├── package
      ├── cmd
      │      ├── server
      │      │      ├── main.go
      │      ├── version
      │      │      ├── version.go
      ├── go.mod
      ├── go.sum
      ├── mocks
      │      ├── PackageRepo.go
      ├── pb
      │      ├── gen
      │      │      ├── mocks
      │      │      │      ├── PackagesServiceClient.go
      │      │      │      ├── PackagesServiceServer.go
      │      │      │      ├── UnsafeEventNotificationServiceServer.go
      │      │      │      ├── UnsafePackagesServiceServer.go
      │      │      ├── package.pb.go
      │      │      ├── package.validator.pb.go
      │      │      ├── package_grpc.pb.go
      │      ├── package.proto
      ├── pkg
      │      ├── config.go
      │      ├── db
      │      │      ├── model.go
      │      │      ├── package_repo.go
      │      │      ├── package_repo_test.go
      │      ├── global.go
      │      ├── server
      │      │      ├── package.go
      │      │      ├── package_test.go
      │      ├── validations
      │      │      ├── validations.go
      │      │      ├── validations_test.go
      ├── test
      │      ├── integration
      │      │      ├── package_test.go

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

This command will run unit tests under all pb/package.proto directories.

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from pb/package.proto.

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

## pb/package.proto

<a name="ukama.data_plan.package.v1.AddPackageRequest"></a>

### AddPackageRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name of the package being added |
| orgID | [string](#string) |  | The ID of the organization the package belongs to |
| active | [bool](#bool) |  | Whether or not the package is currently active |
| duration | [uint64](#uint64) |  | The duration of the package in days |
| simType | [string](#string) |  | The type of SIM card required for the package |
| smsVolume | [int64](#int64) |  | The volume of SMS messages included in the package |
| dataVolume | [int64](#int64) |  | The volume of data included in the package |
| voiceVolume | [int64](#int64) |  | The volume of voice minutes included in the package |
| orgRatesID | [uint64](#uint64) |  | The ID of the organization's rate plan |

<a name="ukama.data_plan.package.v1.AddPackageResponse"></a>

### AddPackageResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| package | [Package](#ukama.data_plan.package.v1.Package) |  | The response message for the Add RPC that returns the added package |

<a name="ukama.data_plan.package.v1.DeletePackageRequest"></a>

### DeletePackageRequest

define a message named DeletePackageRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| packageID | [string](#string) |  | define a string field named packageID with a tag for field validation and JSON name |

<a name="ukama.data_plan.package.v1.DeletePackageResponse"></a>

### DeletePackageResponse

<a name="ukama.data_plan.package.v1.GetByOrgPackageRequest"></a>

### GetByOrgPackageRequest

define a message named GetByOrgPackageRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| orgID | [string](#string) |  | define a string field named orgID with a tag for field validation and JSON name |

<a name="ukama.data_plan.package.v1.GetByOrgPackageResponse"></a>

### GetByOrgPackageResponse

define a message named GetByOrgPackageResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| packages | [Package](#ukama.data_plan.package.v1.Package) | repeated | define a repeated field of type Package named packages |

<a name="ukama.data_plan.package.v1.GetPackageRequest"></a>

### GetPackageRequest

define a message named GetPackageRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| packageID | [string](#string) |  | define a string field named packageID with a tag for field validation and JSON name |

<a name="ukama.data_plan.package.v1.GetPackageResponse"></a>

### GetPackageResponse

define a message named GetPackageResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| package | [Package](#ukama.data_plan.package.v1.Package) |  | define a field of type Package named package |

<a name="ukama.data_plan.package.v1.Package"></a>

### Package

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| packageID | [string](#string) |  | The unique ID of the package |
| name | [string](#string) |  | The name of the package |
| orgID | [string](#string) |  | The ID of the organization the package belongs to |
| active | [bool](#bool) |  | Whether or not the package is currently active |
| duration | [uint64](#uint64) |  | The duration of the package in days |
| simType | [string](#string) |  | The type of SIM card required for the package |
| createdAt | [string](#string) |  | The date and time the package was created |
| deletedAt | [string](#string) |  | The date and time the package was deleted, if applicable |
| updatedAt | [string](#string) |  | The date and time the package was last updated |
| smsVolume | [int64](#int64) |  | The volume of SMS messages included in the package |
| dataVolume | [int64](#int64) |  | The volume of data included in the package |
| voiceVolume | [int64](#int64) |  | The volume of voice minutes included in the package |
| orgRatesID | [uint64](#uint64) |  | The ID of the organization's rate plan |

<a name="ukama.data_plan.package.v1.UpdatePackageRequest"></a>

### UpdatePackageRequest

define a message named UpdatePackageRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| packageID | [string](#string) |  | define a string field named packageID with a tag for field validation and JSON name |
| name | [string](#string) |  | define a string field named name |
| active | [bool](#bool) |  | define a boolean field named active |
| duration | [uint64](#uint64) |  | define a uint64 field named duration |
| simType | [string](#string) |  | define a string field named simType |
| smsVolume | [int64](#int64) |  | define a int64 field named smsVolume |
| dataVolume | [int64](#int64) |  | define a int64 field named dataVolume |
| voiceVolume | [int64](#int64) |  | define a int64 field named voiceVolume |
| orgRatesID | [uint64](#uint64) |  | define a uint64 field named |

<a name="ukama.data_plan.package.v1.UpdatePackageResponse"></a>

### UpdatePackageResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| package | [Package](#ukama.data_plan.package.v1.Package) |  | The response message for the Update RPC that returns the updated package |

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
