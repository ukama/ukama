# Node Registry

Node operations within the Registry System.

## Description
The Node Registry is a sub system that handle various directory operations regarding Nodes within the Registry System. Nodes incoming requests to the Registry System are forwarded to this module for fulfillment. 

## Service definition
The Node Registry exposes the following RPC definitions:

``` proto
service NodeService {
    rpc AttachNodes(AttachNodesRequest) returns (AttachNodesResponse);
    rpc DetachNode(DetachNodeRequest) returns (DetachNodeResponse);
    rpc UpdateNodeState(UpdateNodeStateRequest) returns (UpdateNodeStateResponse);
    rpc UpdateNode(UpdateNodeRequest) returns (UpdateNodeResponse);
    rpc GetNode(GetNodeRequest) returns (GetNodeResponse);
    rpc AddNode(AddNodeRequest) returns (AddNodeResponse);
    rpc Delete(DeleteRequest) returns (DeleteResponse);
}
```


## How to use
### From within the Registry System
Just grab and instrument the Node Registry client stub and make the desired service calls.

### From outside the Registry System
Use the Registry System's API Gateway interface to perform the desired RESTful operations. See the Registry System API Gateway documentatiion for more.
