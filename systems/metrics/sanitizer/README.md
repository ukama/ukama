# Sanitizer service


This service is providing the following features:

    Add support for Prometheus remote_write for metrics pushes
    Sanitize node scraped metrics by appending missing node_id and network_id values

The goal for the sanitizer is to get values of metric `trx_lte_core_active_ue` with incomplete labels such as **network** and **site**, append the missing values for these labels and re-publish the metric under a new name. The metric is re-published under the following new name: `active_subscribers_per_node`
