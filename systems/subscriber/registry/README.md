# Protocol Documentation
<a name="top"></a>

## Table of Contents

 - [RegistryService](#ukama.subscriber.registry.v1.RegistryService)
- [Directory structure](#directory-structure)
- [How to use?](#how-to)
 - [pb/subscriber.proto](#pb/subscriber.proto)
   - [AddSubscriberRequest Messages](#ukama.subscriber.registry.v1.AddSubscriberRequest)
   - [AddSubscriberResponse Messages](#ukama.subscriber.registry.v1.AddSubscriberResponse)
   - [DeleteSubscriberRequest Messages](#ukama.subscriber.registry.v1.DeleteSubscriberRequest)
   - [DeleteSubscriberResponse Messages](#ukama.subscriber.registry.v1.DeleteSubscriberResponse)
   - [GetByNetworkRequest Messages](#ukama.subscriber.registry.v1.GetByNetworkRequest)
   - [GetByNetworkResponse Messages](#ukama.subscriber.registry.v1.GetByNetworkResponse)
   - [GetSubscriberRequest Messages](#ukama.subscriber.registry.v1.GetSubscriberRequest)
   - [GetSubscriberResponse Messages](#ukama.subscriber.registry.v1.GetSubscriberResponse)
   - [ListSubscribersRequest Messages](#ukama.subscriber.registry.v1.ListSubscribersRequest)
   - [ListSubscribersResponse Messages](#ukama.subscriber.registry.v1.ListSubscribersResponse)
   - [Package Messages](#ukama.subscriber.registry.v1.Package)
   - [Sim Messages](#ukama.subscriber.registry.v1.Sim)
   - [Subscriber Messages](#ukama.subscriber.registry.v1.Subscriber)
   - [UpdateSubscriberRequest Messages](#ukama.subscriber.registry.v1.UpdateSubscriberRequest)
   - [UpdateSubscriberResponse Messages](#ukama.subscriber.registry.v1.UpdateSubscriberResponse)
  
- [Scalar Value Types](#scalar-value-types)



<a name="pb/subscriber.proto"></a>
<p align="right"><a href="#top">Top</a></p>


<a name="ukama.subscriber.registry.v1.RegistryService"></a>

# RegistryService
Defines the service for subscriber registry operations

## RPC Functions

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Get | [GetSubscriberRequest](#ukama.subscriber.registry.v1.GetSubscriberRequest) | [GetSubscriberResponse](#ukama.subscriber.registry.v1.GetSubscriberRequest) | Get method to retrieve a subscriber by subscriber ID |
| Add | [AddSubscriberRequest](#ukama.subscriber.registry.v1.AddSubscriberRequest) | [AddSubscriberResponse](#ukama.subscriber.registry.v1.AddSubscriberRequest) | Add method to add a new subscriber |
| Update | [UpdateSubscriberRequest](#ukama.subscriber.registry.v1.UpdateSubscriberRequest) | [UpdateSubscriberResponse](#ukama.subscriber.registry.v1.UpdateSubscriberRequest) | Update method to update an existing subscriber |
| Delete | [DeleteSubscriberRequest](#ukama.subscriber.registry.v1.DeleteSubscriberRequest) | [DeleteSubscriberResponse](#ukama.subscriber.registry.v1.DeleteSubscriberRequest) | Delete method to delete a subscriber by subscriber ID |
| GetByNetwork | [GetByNetworkRequest](#ukama.subscriber.registry.v1.GetByNetworkRequest) | [GetByNetworkResponse](#ukama.subscriber.registry.v1.GetByNetworkRequest) | GetByNetwork method to retrieve subscribers by network ID |
| ListSubscribers | [ListSubscribersRequest](#ukama.subscriber.registry.v1.ListSubscribersRequest) | [ListSubscribersResponse](#ukama.subscriber.registry.v1.ListSubscribersRequest) | ListSubscribers method to retrieve a list of all subscribers |




<a name="#directory-structure"></a>

## Directory structure

      ├── cmd
      │      ├── server
      │      │      ├── main.go
      │      ├── version
      │      │      ├── version.go
      ├── dockerfile
      ├── go.mod
      ├── go.sum
      ├── makefile
      ├── mocks
      │      ├── NetworkInfoClient.go
      │      ├── SimManagerClientProvider.go
      │      ├── SubscriberRepo.go
      │      ├── network.go
      ├── pb
      │      ├── gen
      │      │      ├── mocks
      │      │      │      ├── RegistryServiceClient.go
      │      │      │      ├── RegistryServiceServer.go
      │      │      │      ├── SubscriberRegistryServiceClient.go
      │      │      │      ├── SubscriberRegistryServiceServer.go
      │      │      │      ├── UnsafeRegistryServiceServer.go
      │      │      │      ├── UnsafeSubscriberRegistryServiceServer.go
      │      │      ├── subscriber.pb.go
      │      │      ├── subscriber.validator.pb.go
      │      │      ├── subscriber_grpc.pb.go
      │      ├── subscriber.proto
      ├── pkg
      │      ├── client
      │      │      ├── network.go
      │      │      ├── network_test.go
      │      │      ├── sim_manager_client.go
      │      ├── config.go
      │      ├── db
      │      │      ├── model.go
      │      │      ├── subscriber_repo.go
      │      │      ├── subscriber_repo_test.go
      │      ├── global.go
      │      ├── server
      │      │      ├── subscriber.go
      │      │      ├── subscriber_test.go
      ├── test
      │      ├── integration
      │      │      ├── susbcriber_test.go

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

This command will run unit tests under all pb/subscriber.proto directories.

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from pb/subscriber.proto.

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

## pb/subscriber.proto



<a name="ukama.subscriber.registry.v1.AddSubscriberRequest"></a>

### AddSubscriberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| firstName | [string](#string) |  | first name of the subscriber, must not be empty and have a length greater than 1 |
| lastName | [string](#string) |  | last name of the subscriber, must not be empty and have a length greater than 1 |
| email | [string](#string) |  | email of the subscriber, must be in email format |
| phoneNumber | [string](#string) |  | phone number of the subscriber, must be in phone number format |
| address | [string](#string) |  | address of the subscriber |
| idSerial | [string](#string) |  | idSerial of the subscriber |
| networkID | [string](#string) |  | network ID of the subscriber, must be a UUID and not empty |
| proofOfIdentification | [string](#string) |  | proof of identification of the subscriber, must not be empty and have a length greater than 1 |
| dateOfBirth | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | date of birth of the subscriber |
| gender | [string](#string) |  | gender of the subscriber, must not be empty and have a length greater than 1 |
| orgID | [string](#string) |  | org ID of the subscriber, must be a UUID and not empty |




<a name="ukama.subscriber.registry.v1.AddSubscriberResponse"></a>

### AddSubscriberResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Subscriber | [Subscriber](#ukama.subscriber.registry.v1.Subscriber) |  |  |




<a name="ukama.subscriber.registry.v1.DeleteSubscriberRequest"></a>

### DeleteSubscriberRequest
DeleteSubscriberRequest message is used to delete a subscriber


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscriberID | [string](#string) |  | subscriberID field to be deleted |




<a name="ukama.subscriber.registry.v1.DeleteSubscriberResponse"></a>

### DeleteSubscriberResponse
DeleteSubscriberResponse message is used to delete a subscriber




<a name="ukama.subscriber.registry.v1.GetByNetworkRequest"></a>

### GetByNetworkRequest
GetByNetworkRequest message is used to get all subscribers by network id


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| networkID | [string](#string) |  | networkID field is used to specify the network id for getting subscribers |




<a name="ukama.subscriber.registry.v1.GetByNetworkResponse"></a>

### GetByNetworkResponse
GetByNetworkResponse message is used to get all subscribers by network id


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscribers | [Subscriber](#ukama.subscriber.registry.v1.Subscriber) | repeated | Repeated field of Subscriber message |




<a name="ukama.subscriber.registry.v1.GetSubscriberRequest"></a>

### GetSubscriberRequest
GetSubscriberRequest message is used to get a subscriber by subscriber ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscriberID | [string](#string) |  | subscriberID field is used to specify the subscriber id for getting |




<a name="ukama.subscriber.registry.v1.GetSubscriberResponse"></a>

### GetSubscriberResponse
GetSubscriberResponse message is used to get a subscriber by subscriber ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscriber | [Subscriber](#ukama.subscriber.registry.v1.Subscriber) |  | Subscriber field contains the subscriber information |




<a name="ukama.subscriber.registry.v1.ListSubscribersRequest"></a>

### ListSubscribersRequest
ListSubscribersRequest message is used to list all subscribers




<a name="ukama.subscriber.registry.v1.ListSubscribersResponse"></a>

### ListSubscribersResponse
ListSubscribersResponse message is used to list all subscribers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscribers | [Subscriber](#ukama.subscriber.registry.v1.Subscriber) | repeated | Repeated field of Subscriber message |




<a name="ukama.subscriber.registry.v1.Package"></a>

### Package
Package message defines the structure for a package object


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id field is a string that must be a valid UUID version 0 and cannot be empty |
| startDate | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | startDate field is a google.protobuf.Timestamp and its json representation is "start_date" |
| endDate | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | endDate field is a google.protobuf.Timestamp and its json representation is "end_date" |




<a name="ukama.subscriber.registry.v1.Sim"></a>

### Sim



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| subscriberID | [string](#string) |  | subscriberID field is a string that must be a valid UUID version 0 and cannot be empty |
| networkID | [string](#string) |  | networkID field is a string that must be a valid UUID version 0 and cannot be empty |
| orgID | [string](#string) |  | orgID field is a string that must be a valid UUID version 0 and cannot be empty |
| package | [Package](#ukama.subscriber.registry.v1.Package) |  | package field is a package object |
| iccid | [string](#string) |  | iccid field is a string |
| msisdn | [string](#string) |  | msisdn field is a string that must match phone number format or be empty. |
| imsi | [string](#string) |  | imsi field is a string |
| type | [string](#string) |  | type field is a string |
| status | [string](#string) |  | status field is a string |
| isPhysical | [bool](#bool) |  | isPhysical field is a boolean, and its json representation is "is_physical" |
| firstActivatedOn | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | this field stores the timestamp of the first activation of the SIM card |
| lastActivatedOn | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | this field stores the timestamp of the last activation of the SIM card |
| activationsCount | [uint64](#uint64) |  | This field stores the number of times the SIM card has been activated |
| deactivationsCount | [uint64](#uint64) |  | This field stores the number of times the SIM card has been deactivated |
| allocatedAt | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | This field stores the timestamp when the SIM card was allocated. |




<a name="ukama.subscriber.registry.v1.Subscriber"></a>

### Subscriber



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| orgID | [string](#string) |  |  |
| firstName | [string](#string) |  |  |
| lastName | [string](#string) |  |  |
| subscriberID | [string](#string) |  |  |
| networkID | [string](#string) |  |  |
| email | [string](#string) |  |  |
| phoneNumber | [string](#string) |  |  |
| address | [string](#string) |  |  |
| proofOfIdentification | [string](#string) |  |  |
| created_at | [string](#string) |  |  |
| deleted_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |
| sim | [Sim](#ukama.subscriber.registry.v1.Sim) | repeated |  |
| date_of_birth | [string](#string) |  |  |
| idSerial | [string](#string) |  |  |
| gender | [string](#string) |  |  |




<a name="ukama.subscriber.registry.v1.UpdateSubscriberRequest"></a>

### UpdateSubscriberRequest
UpdateSubscriberRequest defines the request to update a subscriber


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscriberID | [string](#string) |  | subscriberID is the ID of the subscriber to update |
| email | [string](#string) |  | email is the new email address of the subscriber |
| phoneNumber | [string](#string) |  | phoneNumber is the new phone number of the subscriber |
| address | [string](#string) |  | address is the new address of the subscriber |
| idSerial | [string](#string) |  | idSerial is the new idSerial of the subscriber |
| proofOfIdentification | [string](#string) |  | proofOfIdentification is the new proofOfIdentification of the subscriber |




<a name="ukama.subscriber.registry.v1.UpdateSubscriberResponse"></a>

### UpdateSubscriberResponse
UpdateSubscriberResponse defines the response when updating a subscriber

Return nothing when subscriber has been updated







 





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
