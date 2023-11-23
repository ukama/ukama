# Ukama Node System

The Ukama Node System is comprised of four key services: Software Manager, Health, Controller, and Configurator. These services work together to ensure the smooth operation of the Ukama Node.

## Health Service

The Health service is responsible for gathering information about the running applications (capps) from the node via the Node Gateway using REST. This information is then stored in the Health service and subsequently broadcasted through events.

## Software Manager

The Software Manager listens to Health and Software Hub events. When it receives an event, it populates the database with information from the Hub event. It then compares the capps' tags. If a discrepancy is found, an event is triggered with the updated tag, allowing the Node Feeder to instruct the device to update to the latest tag.

## Node Gateway Endpoints

The Node Gateway exposes two endpoints:

1. **Get Capps**:

   ```bash

   curl -X 'GET' 'http://localhost:8036/v1/health/nodes/uk-983794-hnode-78-7830/performance' -H 'accept: application/json'
   ```

2. **Store Capps**:

   ```bash
   curl -X 'POST' 'http://localhost:8036/v1/health/nodes/uk-983794-hnode-78-7830/performance' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{
     "timestamp": "1697867721",
     "capps": [
       {
         "space": "reboot",
         "name": "mesh",
         "tag": "latest",
         "status": "Unknown",
         "resources": [
           {
             "name": "memory",
             "value": "0"
           },
           {
             "name": "disk",
             "value": "0"
           },
           {
             "name": "cpu",
             "value": "0.000000"
           }
         ]
       },
       {
         "space": "boot",
         "name": "mesh",
         "tag": "latest",
         "status": "Unknown",
         "resources": [
           {
             "name": "memory",
             "value": "0"
           },
           {
             "name": "disk",
             "value": "0"
           },
           {
             "name": "cpu",
             "value": "0.000000"
           }
         ]
       }
     ],
     "system": [
       {
         "name": "radio",
         "value": "on"
       }
     ]
   }'
   ```

## Controller and Configurator

The API Gateway exposes Controller and Configurator APIs:

### Controller

- Restart a Site in an Organization:

  - `POST /controllers/networks/:network_id/sites/:site_name/restart`

- Restart a Node:

  - `POST /controllers/nodes/:node_id/restart`

- Restart Multiple Nodes within a Network:
  - `POST /controllers/networks/:network_id/restart-nodes`

### Configurator

- Push Event in Config Store:

  - `POST /configurator/config`

- Apply Config Version:

  - `POST /configurator/config/apply/:commit`

- Get Current Running Config:
  - `GET /configurator/config/node/:node_id`
