# Users Registry

Users operations within the Registry System.

## Description
The User Registry is the sub system that handles various directory operations regarding Users within the Registry System. Users related incoming requests to the Registry System are forwarded to this service for fulfillment.

## How to use from within the Registry System?
### Prerequisites
Before using this repository, make sure the followng tools are installed:

- Go 1.18
- PostgreSQL
- gRPC client of your choice

### Running
From the root directory, run the following command:

```
make server
```

This command will run the server (by default on `localhost:9090`) and create a database named `users` with the relevant tables.

Now that the server is running, you can use any gRPC client to interact with the RPC handlers. For the rest of this document, we'll be using [gRPCurl](https://github.com/fullstorydev/grpcurl).

### Service discovery
By default, all running services support server reflection, which makes describing and testing RPC endpoints even easier. Please, refer to your gRPC client documentation for alternatives introspection methods (such as compiled protoset files or source proto files) if it does not support server reflection.

#### Describing services
Use the describing feature of your gRPC client to get a description of each service exposed by the server.
```shell
# Server supports reflection
> grpcurl -plaintext localhost:9090 describe ukama.users.v1.UserService

service UserService {
  rpc Add ( .ukama.users.v1.AddRequest ) returns ( .ukama.users.v1.AddResponse );
  rpc Deactivate ( .ukama.users.v1.DeactivateRequest ) returns ( .ukama.users.v1.DeactivateResponse );
  rpc Delete ( .ukama.users.v1.DeleteRequest ) returns ( .ukama.users.v1.DeleteResponse );
  rpc Get ( .ukama.users.v1.GetRequest ) returns ( .ukama.users.v1.GetResponse );
  rpc Update ( .ukama.users.v1.UpdateRequest ) returns ( .ukama.users.v1.UpdateResponse );
}
```
Alternative introspection methods when server reflection is not supported:
```shell
# Using compiled protoset files
grpcurl -plaintext -protoset my-protos.bin describe ukama.users.v1.UserService

# Using proto sources
grpcurl -plaintext -import-path ./path-to-protos -proto my-files.proto describe ukama.users.v1.UserService
```

> **Note**
>
>  When using .proto source or protoset files instead of server reflection, this will describe all services defined in the source or protoset files, not the actual services running on the server.


Each RPC can also be described:
```shell
> grpcurl -plaintext localhost:9090 describe ukama.users.v1.UserService.Get

ukama.users.v1.UserService.Get is a method:
rpc Get ( .ukama.users.v1.GetRequest ) returns ( .ukama.users.v1.GetResponse );
```

As well as each proto message:
```shell
> grpcurl -plaintext localhost:9090 describe ukama.users.v1.GetRequest

ukama.users.v1.GetRequest is a message:
message GetRequest {
  string userUuid = 1 [json_name = "user_uuid"];
}
```
#### Calling RPCs
RPCs can be invoked as long as each required input parameter is provided as JSON payload.

Example: **UserService.Get**
```shell
> echo '{"user_uuid":"bcee352c-91e8-46fb-a5a2-1b144543f327"}' | grpcurl -plaintext -d @ localhost:9090 ukama.users.v1.UserService.Get

{
  "user": {
    "name": "Foo Cole",
    "email": "foo@example.com",
    "phone": "0000000000",
    "uuid": "bcee352c-91e8-46fb-a5a2-1b144543f327",
    "registered_since": "2022-11-28T17:33:36.850917Z"
  }
}
```

## How to use from outside the Registry System?
Use the Registry System's API Gateway interface to perform the desired RESTful operations from outside the Registry System. See the Registry System API Gateway documentatiion for more.
