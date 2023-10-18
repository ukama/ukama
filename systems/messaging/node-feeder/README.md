# Node Feeder

Node Feeder sends HTTP request to list of nodes.

## Using Node Feeder

To send request to a group of Nodes send a message to `amq.topic` exchange with routing key `request.cloud.node-feeder`.
The message body should follow the Proto Buffers format::

```
message NodeUpdateRequest {
  string Target = 1;
  string HTTPMethod = 2;
  string Path = 3;
  google.protobuf.Any msg = 4;
}

```

```
{
"target": "MY_ORG.*",  # to send request to all nodes in MY_ORG or "MY_ORG.MY_NODES" if we want to send request to a single node
"httpMethod":"POST",
"Msg": "PROTO BUF MSG",
"path": "/ping"  # path
}
```

In case when wildcarded target is used, node feeder creates messages for each node that matches the pattern.

In case when request fails with http code greater or equal to 500 or connection is refused then Node Feeder will retry the attempt. The configuration of retry mechanism could be changed by setting below env vars

- `LISTENER_EXECUTIONRETRYCOUNT` - number of retries
- `LISTENER_RETRYPERIODSEC` - period between retries in seconds

## Implementation

Node Feeder uses queuse to manage requests lifecycle. Here is the queues topology
![image](https://user-images.githubusercontent.com/154290/147089205-14058d8a-ec92-4c43-b777-7e9f3fc42af0.png)

When request is sent to 'node-feeder' queue it is consumed by 'node-feeder' service. If request fails with 'connection' issue or with http code greater or equal to 500 then it is sent to "waiting" queue where all messages have default TTL. When TTL expires those messages are sent back to node-feeder queue and consumed by service again.
