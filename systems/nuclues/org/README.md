# Org Registry

Org operations within the Registry System.

## Description
The Org Registry is the sub system that handles various directory operations regarding Organizations within the Registry System. Orgs related incoming requests to the Registry System are forwarded to this service for fulfillment.

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

This command will run the server (by default on `localhost:9090`) and create a database named `org` with the relevant tables.

Now that the server is running, you can use any gRPC client to interact with the RPC handlers. For the rest of this document, we'll be using [gRPCurl](https://github.com/fullstorydev/grpcurl).

### Service discovery
By default, all running services support server reflection, which makes describing and testing RPC endpoints even easier. Please, refer to your gRPC client documentation for alternatives introspection methods (such as compiled protoset files or source proto files) if it does not support server reflection.

#### Describing services
Use the describing feature of your gRPC client to get a description of each service exposed by the server.
```shell
# Server supports reflection
> grpcurl -plaintext localhost:9090 describe ukama.org.v1.OrgService

ukama.org.v1.OrgService is a service:
service OrgService {
  rpc Add ( .ukama.org.v1.AddRequest ) returns ( .ukama.org.v1.AddResponse );
  rpc AddMember ( .ukama.org.v1.MemberRequest ) returns ( .ukama.org.v1.MemberResponse );
  rpc Get ( .ukama.org.v1.GetRequest ) returns ( .ukama.org.v1.GetResponse );
  rpc GetByName ( .ukama.org.v1.GetByNameRequest ) returns ( .ukama.org.v1.GetByNameResponse );
  rpc GetByOwner ( .ukama.org.v1.GetByOwnerRequest ) returns ( .ukama.org.v1.GetByOwnerResponse );
  rpc GetMember ( .ukama.org.v1.MemberRequest ) returns ( .ukama.org.v1.MemberResponse );
  rpc GetMembers ( .ukama.org.v1.GetMembersRequest ) returns ( .ukama.org.v1.GetMembersResponse );
  rpc RegisterUser ( .ukama.org.v1.RegisterUserRequest ) returns ( .ukama.org.v1.MemberResponse );
  rpc RemoveMember ( .ukama.org.v1.MemberRequest ) returns ( .ukama.org.v1.MemberResponse );
  rpc UpdateMember ( .ukama.org.v1.UpdateMemberRequest ) returns ( .ukama.org.v1.MemberResponse );
  rpc UpdateUser ( .ukama.org.v1.UpdateUserRequest ) returns ( .ukama.org.v1.UpdateUserResponse );
}
```
**Demo:**
![orgs_services](https://user-images.githubusercontent.com/10562122/205147195-251aff9c-9ff2-4f80-98f7-5e4e1aeace0d.gif)

Alternative introspection methods when server reflection is not supported:
```shell
# Using compiled protoset files
grpcurl -plaintext -protoset my-protos.bin describe ukama.org.v1.OrgService

# Using proto sources
grpcurl -plaintext -import-path ./path-to-protos -proto my-files.proto describe ukama.org.v1.OrgService
```

> **Note**
>
>  When using .proto source or protoset files instead of server reflection, this will describe all services defined in the source or protoset files, not the actual services running on the server.


Each RPC can also be described:
```shell
> grpcurl -plaintext localhost:9090 describe ukama.org.v1.OrgService.Get

ukama.org.v1.OrgService.Get is a method:
rpc Get ( .ukama.org.v1.GetRequest ) returns ( .ukama.org.v1.GetResponse );
```
**Demo:**
![orgs_rpc](https://user-images.githubusercontent.com/10562122/205147189-5a4211cd-ec89-470a-94a3-5c9e09555367.gif)


As well as each proto message:
```shell
> grpcurl -plaintext localhost:9090 describe ukama.org.v1.Organization

ukama.org.v1.Organization is a message:
message Organization {
  uint64 id = 1;
  string name = 2;
  string owner = 3;
  string certificate = 4;
  bool isDeactivated = 5 [json_name = "is_deactivated"];
  .google.protobuf.Timestamp created_at = 6 [json_name = "created_at"];
}
```
**Demo:**
![orgs_proto](https://user-images.githubusercontent.com/10562122/205147187-7dc1a8bd-68ef-4fb8-9788-6668930d4288.gif)


#### Calling RPCs
RPCs can be invoked as long as each required input parameter (proto message) is provided as a valid JSON payload to the gRPC client.

Example: **OrgService.Get**

You can manually construct your JSON payload with plain `echo` shell statements and pipe it to the standard input of the gRPC client (assuming that your gRPC client of choice can read from standard input), like the following:
```shell
> echo '{"org_id":1}' | grpcurl -plaintext -d @ localhost:9091 ukama.org.v1.OrgService.Get

{
  "org": {
    "id": "1",
    "name": "ukama",
    "owner": "bcee352c-91e8-46fb-a5a2-1b144543f327",
    "created_at": "2022-11-28T17:33:36.922984Z"
  }
}
```

Or you can use JSON printing tools (like [jo](https://github.com/jpmens/jo)) to make the JSON payload construction easier:
```shell
> jo org_id=1 | grpcurl -plaintext -d @ localhost:9091 ukama.org.v1.OrgService.Get

{
  "org": {
    "id": "1",
    "name": "ukama",
    "owner": "bcee352c-91e8-46fb-a5a2-1b144543f327",
    "created_at": "2022-11-28T17:33:36.922984Z"
  }
}
```

**Demo:**
![orgs_call](https://user-images.githubusercontent.com/10562122/205152211-7de5e9af-430b-4b10-8f4c-e9b08750abf6.gif)


## How to use from outside the Registry System?
Use the Registry System's API Gateway interface to perform the equivalent RESTful operations from outside the Registry System. See the Registry System API Gateway documentatiion for more.
