# Protocol Documentation
<a name="top"></a>

## Table of Contents


- [Directory structure](#directory-structure)
- [How to use?](#how-to)
 - [](#)
  

 - [PackagesService](#ukama.data_plan.package.v1.PackagesService)
- [Directory structure](#directory-structure)
- [How to use?](#how-to)
 - [package.proto](#package.proto)
   - [AddPackageRequest Messages](#ukama.data_plan.package.v1.AddPackageRequest)
   - [AddPackageResponse Messages](#ukama.data_plan.package.v1.AddPackageResponse)
   - [DeletePackageRequest Messages](#ukama.data_plan.package.v1.DeletePackageRequest)
   - [DeletePackageResponse Messages](#ukama.data_plan.package.v1.DeletePackageResponse)
   - [GetByOrgPackageRequest Messages](#ukama.data_plan.package.v1.GetByOrgPackageRequest)
   - [GetByOrgPackageResponse Messages](#ukama.data_plan.package.v1.GetByOrgPackageResponse)
   - [GetPackageRequest Messages](#ukama.data_plan.package.v1.GetPackageRequest)
   - [GetPackageResponse Messages](#ukama.data_plan.package.v1.GetPackageResponse)
   - [Package Messages](#ukama.data_plan.package.v1.Package)
   - [PackageMarkup Messages](#ukama.data_plan.package.v1.PackageMarkup)
   - [PackageRates Messages](#ukama.data_plan.package.v1.PackageRates)
   - [UpdatePackageRequest Messages](#ukama.data_plan.package.v1.UpdatePackageRequest)
   - [UpdatePackageResponse Messages](#ukama.data_plan.package.v1.UpdatePackageResponse)
  
- [Scalar Value Types](#scalar-value-types)



<a name=""></a>
<p align="right"><a href="#top">Top</a></p>




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

This command will run unit tests under all  directories.

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from .

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

## 






 





<a name="package.proto"></a>
<p align="right"><a href="#top">Top</a></p>


<a name="ukama.data_plan.package.v1.PackagesService"></a>

# PackagesService


## RPC Functions

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Get | [GetPackageRequest](#ukama.data_plan.package.v1.GetPackageRequest) | [GetPackageResponse](#ukama.data_plan.package.v1.GetPackageRequest) |  |
| GetDetails | [GetPackageRequest](#ukama.data_plan.package.v1.GetPackageRequest) | [GetPackageResponse](#ukama.data_plan.package.v1.GetPackageRequest) |  |
| Add | [AddPackageRequest](#ukama.data_plan.package.v1.AddPackageRequest) | [AddPackageResponse](#ukama.data_plan.package.v1.AddPackageRequest) |  |
| Delete | [DeletePackageRequest](#ukama.data_plan.package.v1.DeletePackageRequest) | [DeletePackageResponse](#ukama.data_plan.package.v1.DeletePackageRequest) |  |
| Update | [UpdatePackageRequest](#ukama.data_plan.package.v1.UpdatePackageRequest) | [UpdatePackageResponse](#ukama.data_plan.package.v1.UpdatePackageRequest) |  |
| GetByOrg | [GetByOrgPackageRequest](#ukama.data_plan.package.v1.GetByOrgPackageRequest) | [GetByOrgPackageResponse](#ukama.data_plan.package.v1.GetByOrgPackageRequest) |  |




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

This command will run unit tests under all package.proto directories.

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from package.proto.

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

## package.proto



<a name="ukama.data_plan.package.v1.AddPackageRequest"></a>

### AddPackageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| orgId | [string](#string) |  |  |
| active | [bool](#bool) |  |  |
| duration | [uint64](#uint64) |  |  |
| simType | [string](#string) |  |  |
| smsVolume | [int64](#int64) |  |  |
| dataVolume | [int64](#int64) |  |  |
| voiceVolume | [int64](#int64) |  |  |
| dlbr | [int64](#int64) |  |  |
| ulbr | [int64](#int64) |  |  |
| markup | [double](#double) |  |  |
| type | [string](#string) |  |  |
| dataUnit | [string](#string) |  |  |
| voiceUnit | [string](#string) |  |  |
| messageunit | [string](#string) |  |  |
| flatrate | [bool](#bool) |  |  |
| amount | [double](#double) |  |  |
| effectiveAt | [string](#string) |  |  |
| endAt | [string](#string) |  |  |
| country | [string](#string) |  |  |
| provider | [string](#string) |  |  |
| apn | [string](#string) |  |  |
| baserate | [string](#string) |  |  |
| ownerId | [string](#string) |  |  |
| from | [string](#string) |  |  |
| to | [string](#string) |  |  |




<a name="ukama.data_plan.package.v1.AddPackageResponse"></a>

### AddPackageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| package | [Package](#ukama.data_plan.package.v1.Package) |  |  |




<a name="ukama.data_plan.package.v1.DeletePackageRequest"></a>

### DeletePackageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| orgId | [string](#string) |  |  |




<a name="ukama.data_plan.package.v1.DeletePackageResponse"></a>

### DeletePackageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| orgId | [string](#string) |  |  |




<a name="ukama.data_plan.package.v1.GetByOrgPackageRequest"></a>

### GetByOrgPackageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| orgId | [string](#string) |  |  |




<a name="ukama.data_plan.package.v1.GetByOrgPackageResponse"></a>

### GetByOrgPackageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| packages | [Package](#ukama.data_plan.package.v1.Package) | repeated |  |




<a name="ukama.data_plan.package.v1.GetPackageRequest"></a>

### GetPackageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |




<a name="ukama.data_plan.package.v1.GetPackageResponse"></a>

### GetPackageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| package | [Package](#ukama.data_plan.package.v1.Package) |  |  |




<a name="ukama.data_plan.package.v1.Package"></a>

### Package



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| orgId | [string](#string) |  |  |
| active | [bool](#bool) |  |  |
| duration | [uint64](#uint64) |  |  |
| simType | [string](#string) |  |  |
| createdAt | [string](#string) |  |  |
| deletedAt | [string](#string) |  |  |
| updatedAt | [string](#string) |  |  |
| smsVolume | [int64](#int64) |  |  |
| dataVolume | [int64](#int64) |  |  |
| voiceVolume | [int64](#int64) |  |  |
| dlbr | [int64](#int64) |  |  |
| ulbr | [int64](#int64) |  |  |
| rate | [PackageRates](#ukama.data_plan.package.v1.PackageRates) |  |  |
| markup | [PackageMarkup](#ukama.data_plan.package.v1.PackageMarkup) |  |  |
| type | [string](#string) |  |  |
| dataUnit | [string](#string) |  |  |
| voiceUnit | [string](#string) |  |  |
| messageunit | [string](#string) |  |  |
| flatrate | [bool](#bool) |  |  |
| currency | [string](#string) |  |  |
| from | [string](#string) |  |  |
| to | [string](#string) |  |  |
| country | [string](#string) |  |  |
| provider | [string](#string) |  |  |
| apn | [string](#string) |  |  |
| OwnerId | [string](#string) |  |  |
| amount | [double](#double) |  |  |




<a name="ukama.data_plan.package.v1.PackageMarkup"></a>

### PackageMarkup



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| baserate | [string](#string) |  |  |
| markup | [double](#double) |  |  |




<a name="ukama.data_plan.package.v1.PackageRates"></a>

### PackageRates



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| SmsMo | [double](#double) |  |  |
| SmsMt | [double](#double) |  |  |
| Data | [double](#double) |  |  |
| Amount | [double](#double) |  |  |




<a name="ukama.data_plan.package.v1.UpdatePackageRequest"></a>

### UpdatePackageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| orgId | [string](#string) |  |  |
| name | [string](#string) |  |  |
| active | [bool](#bool) |  |  |




<a name="ukama.data_plan.package.v1.UpdatePackageResponse"></a>

### UpdatePackageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| package | [Package](#ukama.data_plan.package.v1.Package) |  |  |







 





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
