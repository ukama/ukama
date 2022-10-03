# Node Gateway

Restful API Gateway exposed to the Nodes for bootstrapping requests.

## Description
A Node Gateway is a single entry point dedicated to the Nodes. It is the one which communicate with Nodes and Units (and their devices within). For the `Init` system, the Node Gateway allows only the `GET` method over a `nodes` resource for bootstrapping any specific node, given its `nodeId`.

## Interface

### Get Node

```
curl --request GET \
  --url http://INIT_SYSTEM_NODE_GATEWAY-URL/nodes/some-ukama-node-id
```
Response:
```
{
  "node": "some-ukama-node-id",
  "org": "test-org",
  "certificate": "test-org cert",
  "ip": "192.124.23.1"
}
```
