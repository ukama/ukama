# Data-Plan Gateway

Restful API Gateway exposed basic crud operations for packages sub system.

## Description

A Data-plan Gateway is a single entry point dedicated to the `packages`,`base-rates` etc... It is the one which communicate with `OrgRate` service to get base-rate per region.

## Project structure

```bash
├── cmd (server initialization)
│   ├── main.go
│   ├── version ( Api version )
└── pkg
    ├── client
    │   └── DataPlan
    ├── rest
       ├── api (Api request schema structs)
       └── router ( Api routes and handlers)
```

## Getting started

```bash
https://github.com/ukama/ukama.git

cd /systems/data-plan

```

### Prerequisites

- [Go lang](https://go.dev//) >= v1.18
- An API platform for rest(postman,insomnia ,etc...)

## Interface

### Upload baseRates

```bash
curl -X 'POST' \
  'http://localhost:8080/v1/baseRates' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "effective_at": "string",
  "file_url": "string",
  "sim_type": "string"
}'
```

Response:

```
{
	"rates": [
		{
			"id": "2",
			"name": "",
			"orgId": "",
			"active": true,
			"duration": "43",
			"simType": "INTER_MNO_DATA",
			"createdAt": "",
			"deletedAt": "",
			"updatedAt": "",
			"smsVolume": "0",
			"dataVolume": "23",
			"voiceVolume": "23",
			"orgRatesId": "55"
		}
	]
}
```

### Get baseRate

```bash
curl -X 'GET' \
  'http://localhost:8080/v1/baseRates/3' \
  -H 'accept: application/json'
```

Response:

```
{
  "rate": {
    "id": "3",
    "country": "",
    "network": "",
    "vpmn": "",
    "imsi": "",
    "smsMo": "$0.1",
    "smsMt": "$0.1",
    "data": "$0.3",
    "2g": "2G",
    "3g": "3G",
    "5g": "",
    "lte": "",
    "lteM": "",
    "apn": "",
    "createdAt": "2022-11-23 13:36:45.108238 +0200 CAT",
    "deletedAt": "0001-01-01 00:00:00 +0000 UTC",
    "updatedAt": "2022-11-23 13:36:45.108238 +0200 CAT",
    "effectiveAt": "2023-06-23T09:16:50Z",
    "endAt": "",
    "simType": "inter_none"
  }
}
```

### Get baseRates

```bash
curl -X 'GET' \
  'http://localhost:8080/v1/baseRates' \
  -H 'accept: application/json'
```

```
GetBaseRates end point accepts query params (country,effectiveAt,simType,to,from)

```

Response:

```
{
  "rates": [
    {
      "id": "1",
      "country": "",
      "network": "",
      "vpmn": "",
      "imsi": "1",
      "smsMo": "$0.1",
      "smsMt": "$0.1",
      "data": "$0.1",
      "2g": "2G",
      "3g": "",
      "5g": "",
      "lte": "",
      "lteM": "",
      "apn": "",
      "createdAt": "2022-11-23 13:36:45.108238 +0200 CAT",
      "deletedAt": "0001-01-01 00:00:00 +0000 UTC",
      "updatedAt": "2022-11-23 13:36:45.108238 +0200 CAT",
      "effectiveAt": "2023-06-23T09:16:50Z",
      "endAt": "",
      "simType": "inter_none"
    },
    {
      "id": "2",
      "country": "",
      "network": "",
      "vpmn": "",
      "imsi": "1",
      "smsMo": "$0.2",
      "smsMt": "$0.1",
      "data": "$0.2",
      "2g": "2G",
      "3g": "",
      "5g": "",
      "lte": "",
      "lteM": "",
      "apn": "",
      "createdAt": "2022-11-23 13:36:45.108238 +0200 CAT",
      "deletedAt": "0001-01-01 00:00:00 +0000 UTC",
      "updatedAt": "2022-11-23 13:36:45.108238 +0200 CAT",
      "effectiveAt": "2023-06-23T09:16:50Z",
      "endAt": "",
      "simType": "inter_none"
    },
  ]
}
```

### Get Package

```bash
 curl -X 'GET' \
  'http://localhost:8080/v1/packages/{package}' \
  -H 'accept: application/json'

```

Response:

```
{
	"packages": [
		{
			"id": "2",
			"name": "",
			"orgId": "",
			"active": true,
			"duration": "43",
			"simType": "INTER_MNO_DATA",
			"createdAt": "2022-11-22 11:38:47.953352 +0200 CAT",
			"deletedAt": "0001-01-01 00:00:00 +0000 UTC",
			"updatedAt": "2022-11-22 11:38:47.953352 +0200 CAT",
			"smsVolume": "0",
			"dataVolume": "23",
			"voiceVolume": "23",
			"orgRatesId": "55"
		}
	]
}
```

### Add Package

```bash
curl -X 'PUT' \
  'http://localhost:8080/v1/packages' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "active": true,
  "data_volume": 0,
  "duration": 0,
  "name": "string",
  "org_id": 0,
  "org_rates_id": 0,
  "sim_type": "string",
  "sms_volume": 0,
  "voice_volume": 0
}'

```

Response:

```
{
	"package": {
		"id": "3",
		"name": "daily",
		"orgId": "12345",
		"active": true,
		"duration": "43",
		"simType": "INTER_MNO_DATA",
		"createdAt": "2022-11-22 11:40:28.118044 +0200 CAT",
		"deletedAt": "0001-01-01 00:00:00 +0000 UTC",
		"updatedAt": "2022-11-22 11:40:28.118044 +0200 CAT",
		"smsVolume": "0",
		"dataVolume": "23",
		"voiceVolume": "23",
		"orgRatesId": "55"
	}
}
```

### Update Package

```bash
 curl -X 'PATCH' \
  'http://localhost:8080/v1/packages' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "active": true,
  "data_volume": 0,
  "duration": 0,
  "id": 0,
  "name": "string",
  "org_id": 0,
  "org_rates_id": 0,
  "sim_type": "string",
  "sms_volume": 0,
  "voice_volume": 0
}'

```

Response:

```
{
	"packages": [
		{
			"id": "2",
			"name": "some-package-name",
			"orgId": "orgID",
			"active": true,
			"duration": "55",
			"simType": "INTER_NONE",
			"createdAt": "2022-11-22 11:38:47.953352 +0200 CAT",
			"deletedAt": "0001-01-01 00:00:00 +0000 UTC",
			"updatedAt": "2022-11-22 11:38:47.953352 +0200 CAT",
			"smsVolume": "65",
			"dataVolume": "40",
			"voiceVolume": "44",
			"orgRatesId": "40"
		}
	]
}
```

### Delete Package

```bash
curl -X 'DELETE' \
  'http://localhost:8080/v1/packages/{package}' \
  -H 'accept: application/json'

```

Response:

```
{
	"id": "id-of-deleted-package"
}
```
