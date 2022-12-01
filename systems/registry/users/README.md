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
**Demo:**
![users_services](https://user-images.githubusercontent.com/10562122/205072320-da6c4e55-b49b-4820-8281-f6d09c43bf46.gif)

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
**Demo:**
![users_rpc](https://user-images.githubusercontent.com/10562122/205072314-316fffd1-d042-4b39-b688-bf52121619e8.gif)

As well as each proto message:
```shell
> grpcurl -plaintext localhost:9090 describe ukama.users.v1.User

ukama.users.v1.User is a message:
message User {
  string name = 1;
  string email = 2;
  string phone = 3;
  string uuid = 4;
  bool isDeactivated = 5 [json_name = "is_deactivated"];
  .google.protobuf.Timestamp created_at = 6 [json_name = "registered_since"];
}
```
**Demo:**
![users_proto](https://user-images.githubusercontent.com/10562122/205072308-4af74f4b-c0d8-4146-820b-a30601ed6316.gif)

#### Calling RPCs
RPCs can be invoked as long as each required input parameter (proto message) is provided as a valid JSON payload to the gRPC client.

Example: **UserService.Get**

You can manually construct your JSON payload with plain `echo` shell statements and pipe it to the standard input of the gRPC client (assuming that your gRPC client of choice can read from standard input), like the following:
```shell
> echo '{"user_uuid":"23829813-0eaf-4586-8176-52e50318ab9b"}' | grpcurl -plaintext -d @ localhost:9090 ukama.users.v1.UserService.Get

{
  "user": {
    "name": "Foo Cole",
    "email": "foo@example.com",
    "phone": "1111111111",
    "uuid": "23829813-0eaf-4586-8176-52e50318ab9b",
    "registered_since": "2022-12-01T12:32:19.422218Z"
  }
}
```

Or you can use JSON printing tools (like [jo](https://github.com/jpmens/jo)) to make the JSON payload construction easier:
```shell
> jo user_uuid=23829813-0eaf-4586-8176-52e50318ab9b | grpcurl -plaintext -d @ localhost:9090 ukama.users.v1.UserService.Get

{
  "user": {
    "name": "Foo Cole",
    "email": "foo@example.com",
    "phone": "1111111111",
    "uuid": "23829813-0eaf-4586-8176-52e50318ab9b",
    "registered_since": "2022-12-01T12:32:19.422218Z"
  }
}
```
**Demo:**
![users_call](https://user-images.githubusercontent.com/10562122/205072298-27087a22-725e-4c16-b087-b147c4a90ba5.gif)

## How to use from outside the Registry System?
Use the Registry System's API Gateway interface to perform the equivalent RESTful operations from outside the Registry System. See the Registry System API Gateway documentatiion for more.
