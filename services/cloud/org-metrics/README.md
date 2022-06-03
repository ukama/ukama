# Org Metrics

The purpose of `org-metrics` service is to expose and endpoint with org and network level metrics for the Prometheus to scrape.

Metrics for all organizations are exposed at `/` and `/metrics` paths. 
Port is configured by `server.port` key or `SERVER_PORT` env var. Default is 10251.  
