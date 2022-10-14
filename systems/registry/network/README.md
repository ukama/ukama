# Node Network  Registry

user->org->nodes relationships.

## Description
The Node Network Registry is a sub system that handle various directory operations regarding the user->org->nodes relationships within the Registry System. The relationships is as follows: every node belongs to an organization and every organization belongs to a user (owner).

## Service definition
The Org Registry exposes the following RPC definitions:

``` proto
service NetworkService {
    rpc Add(AddRequest) returns (AddResponse);
    // list all orgs and networks in the network
    rpc List(ListRequest) returns (ListResponse);
    rpc Delete(DeleteRequest) returns (DeleteResponse);

    rpc AddNode(AddNodeRequest) returns (AddNodeResponse);
    rpc DeleteNode(DeleteNodeRequest) returns (DeleteNodeResponse);
    rpc GetNodes(GetNodesRequest) returns (GetNodesResponse);
    rpc UpdateNode(UpdateNodeRequest) returns (UpdateNodeResponse);
    rpc GetNode(GetNodeRequest) returns (GetNodeResponse);
}
```

## How to use
### From within the Registry System
Just grab and instrument the Network client stub and make the desired service calls.

### From outside the Registry System
Use the Registry System's API Gateway interface to perform the desired RESTful operations. See the Registry System API Gateway documentatiion for more.

