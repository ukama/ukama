# Org Registry

Org operations within the Registry System.

## Description
The Org Registry is a sub system that handle various directory operations regarding Orgs within the Registry System. Orgs incoming requests to the Registry System are forwarded to this module for fulfillment.

## Service definition
The Org Registry exposes the following RPC definitions:

``` proto
service OrgService {
    rpc Get(GetRequest) returns (GetResponse);
    rpc Add(AddRequest) returns (AddResponse);
    rpc Delete(DeleteRequest) returns (DeleteResponse);
}
```


## How to use
### From within the Registry System
Just grab and instrument the Org Registry client stub and make the desired service calls.

### From outside the Registry System
Use the Registry System's API Gateway interface to perform the desired RESTful operations. See the Registry System API Gateway documentatiion for more.
