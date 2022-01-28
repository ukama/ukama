# Network service
Network service is a service that stores node-id to IP mapping. 
It exposes several interfaces to query and update the mapping: 
1. [GRPC interface](/pb/net.proto)
2. [Prometheus HTTP](https://prometheus.io/docs/prometheus/2.31/configuration/configuration/#http_sd_config) targets config that can be used by 
Prometheus to scrape metrics from nodes. Available at `/prometheus`
3. DNS GRPC. Implementation of [CoreDNS grpc plugin](https://coredns.io/plugins/grpc/) 
that could be integrated with CoreDNS to resolve node-id on DNS level. Useful for running in Kubernetes  


## Running Locally

1. Build the app:
`make build`
2. Run docker-compose
`docker-compose up --build`
3. Add DNS record to etcd:
` docker exec -it etcd etcdctl put uk-sa2203-hnode-a1-0a16 172.10.0.1`
4. Login to bastion: `docker exec -it bastion bash`
5. Check the resolution of node's hostname `dig -p 53 @coredns uk-sa2203-hnode-a1-0a16.node.mesh`
```
...
;; ANSWER SECTION:
uk-sa2203-hnode-a1-0a16.node.mesh. 0 IN A       172.10.0.1
...
```