# Users Registry

Users operations within the Registry System.

## Description
The Users Registry is a sub system that handle various directory operations regarding users within the Registry System. It is responsible for managing users and serve as a proxy for managing sim cards. It stores the user details and ICCID assigned to them.
The Users service relies on the Sim Manager service to source the information about sim cards(IMSI, data usage) and control their properties such as services availability.

## Service definition
### Users Service
The User service exposes the following RPC definitions:

``` proto
service UserService {
    rpc Add (AddRequest) returns (AddResponse);
    rpc AddInternal (AddInternalRequest) returns (AddInternalResponse);
    rpc Delete (DeleteRequest) returns (DeleteResponse);
    rpc List(ListRequest) returns (ListResponse);
    rpc Get(GetRequest) returns (GetResponse);
    rpc Update(UpdateRequest) returns (UpdateResponse);
    rpc GenerateSimToken(GenerateSimTokenRequest) returns (GenerateSimTokenResponse);
    rpc SetSimStatus (SetSimStatusRequest) returns (SetSimStatusResponse);
    rpc DeactivateUser(DeactivateUserRequest) returns (DeactivateUserResponse);
    rpc GetQrCode(GetQrCodeRequest) returns (GetQrCodeResponse);
}
```

### Sim Manager Service
The Sim Manager exposes the following RPC definitions:

``` proto
service SimManagerService {
    rpc GetSimStatus(GetSimStatusRequest) returns (GetSimStatusResponse);
    rpc SetServiceStatus(SetServiceStatusRequest) returns (SetServiceStatusResponse);
    rpc GetSimInfo(GetSimInfoRequest) returns (GetSimInfoResponse);
    rpc TerminateSim(TerminateSimRequest) returns (TerminateSimResponse);
    rpc GetUsage(GetUsageRequest) returns (GetUsageResponse);
    rpc GetQrCode(GetQrCodeRequest) returns (GetQrCodeResponse);
}
```


## How to use
### From within the Registry System
Just grab and instrument any client stub and make the desired service calls.

### From outside the Registry System
Use the Registry System's API Gateway interface to perform the desired RESTful operations. See the Registry System API Gateway documentatiion for more.

