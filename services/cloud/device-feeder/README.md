# Device Feeder

Device Feeder sends HTTP request to list of devices.



## Using Device Feeder
To send request to a group of devices send a message to `amq.topic` exchange with routing key `request.cloud.device-feeder`.
The body should have the following format:

```
{
"target": "MY_ORG.*",  # to send request to all devices in MY_ORG or "MY_ORG.MY_DEVICE" if we want to send request to a single device
"httpMethod":"POST",  
"body": "THE BODY OF THE REQUEST", 
"path": "/ping"  # path 
}
```

In case when wildcarded target is used, device feeder creates messages for each node that matches the pattern. 

In case when request fails with http code greater or equal to 500 or connection is refused then Device Feeder will retry the  attempt. The configuration of retry mechanism could be changed by setting below env vars 
- `LISTENER_EXECUTIONRETRYCOUNT` - number of retries
- `LISTENER_RETRYPERIODSEC` - period between retries in seconds

## Implementation

Device Feeder uses queuse to manage requests lifecycle. Here is the queues topology 
![image](https://user-images.githubusercontent.com/154290/147089205-14058d8a-ec92-4c43-b777-7e9f3fc42af0.png)

When request is sent to 'device-feeder' queue it is consumed by 'device-feeder' service. If request fails with 'connection' issue or with http code greater or equal to 500 then it is sent to "waiting" queue where all messages have default TTL. When TTL expires those messages are sent back to device-feeder queue and consumed by service again.
