# Org Metrics

The purpose of `org-metrics` service is to expose and endpoint with org and network level metrics for the Promethus to scrape.

Metrics for an organization are exposed at `/:org_name` path and for network  on `/:org_name/:network_name`. 