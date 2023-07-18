# API Gateway for  Registry System

This service provides a Rest APIs to interact with the Registry System.

## Interface 

### Org
#### Get Org

```curl
curl -X 'GET' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/{org}' \
  -H 'accept: application/json'
```

#### Add Org

```curl
curl -X 'PUT' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG' \
  -H 'accept: application/json'
```

### Nodes
#### Get Nodes

```curl
curl -X 'GET' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/{org}/nodes' \
  -H 'accept: application/json''
```

#### Get Node

```curl
curl -X 'GET' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/nodes/NODE' \
  -H 'accept: application/json'
```

#### Put Node

```curl
curl -X 'PUT' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/nodes/NODE' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "node": {
    "attached": [
      {
        "nodeId": "string"
      }
    ],
    "name": "string"
  }
}'
```
#### Delete Node

```curl
curl -X 'DELETE' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/nodes/NODE' \
  -H 'accept: application/json'
```

#### Patch Node

```curl
curl -X 'PATCH' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/nodes/NODE' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "node": {
    "attached": [
      {
        "nodeId": "string"
      }
    ],
    "name": "string"
  }
}'
```

#### Delete Attached Node

```curl
curl -X 'DELETE' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/nodes/NODE/attached/ATTACHED_ID' \
  -H 'accept: application/json'
```

### Network Users
#### Get Users

```curl
curl -X 'GET' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/{org}/users' \
  -H 'accept: application/json'
```

#### Post User

```curl
curl -X 'POST' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/users' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "string",
  "name": "string",
  "phone": "string",
  "simToken": "string"
}'
```
#### Get User

```curl
curl -X 'GET' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/users/USER' \
  -H 'accept: application/json'
```
#### Delete User

```curl
curl -X 'DELETE' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/nodes/NODE/attached/ATTACHED_ID' \
  -H 'accept: application/json'
```
#### Patch User

```curl
curl -X 'PATCH' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/users/USER' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "string",
  "isDeactivated": true,
  "name": "string",
  "phone": "string"
}'
```
#### Get E-Sim QR Code

```curl
curl -X 'GET' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/users/USER/sims/ICCID/qr' \
  -H 'accept: application/json'
```
#### Put Sim card Services (enable/disable)

```curl
curl -X 'PUT' \
  'http://REGISTRY-SYSTEM-URL/v1/orgs/ORG/users/USER/sims/ICCID/services' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "carrier": {
    "data": true,
    "sms": true,
    "voice": true
  },
  "ukama": {
    "data": true,
    "sms": true,
    "voice": true
  }
}'
```
