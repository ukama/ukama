# Node Metrics Server

Node metrics server is a simple HTTP server that exposes node metrics from prometheus.
Service expose Swagger interface and open API specification. 


### Metrics Configuration 

List of metrics that node-metrics server expose is configured via config file. 
Here is a basic example:
``` yaml
  server: 
    port: 8080
  debugMode: true
  nodeMetrics:
    metricsServer: "http://localhost:8080"
    metrics:
      cpu: { needRate: true, metric: trx_soc_cpu_usage }
      memory: { needRate: true, metric: trx_memory_ddr_used, rateInterval: 5m }      
```

`metrics` tag contains the list of expose metrics. In our case node-metrics exposes two metrics cpu and memory 
that could be requested using `nodes/[NODE_ID]/metrics/[METRIC_NAME]?from=1643108704&to=1643195104&step=3600` format. 

In our case we can send two requests: 
```
GET http://localhost:8080/nodes/uk-test36-hnode-a1-30df/metrics/cpu?from=1643108704&to=1643195104&step=3600
``` 
and
```
GET http:// localhost:8080/nodes/uk-test36-hnode-a1-30df/metrics/memory?from=1643108704&to=1643195104&step=3600
```


Adding a new "metric" will make `node-metrics` service serve it from the path that correspond to the name of the metric.
- `metric`(required)  - metric name in prometheus 
- `needRate`(optional, default=false) -  flag indicates whether `rate` prometheus function is applied to the metric query
- `rateInterval`(optional, default=1h)  - could be set for metrics that have `needRate` enabled. It has prometheus duration format for ex 2h, 10m etc. If it's empty then default is `1h` 