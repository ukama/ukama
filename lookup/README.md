# Node Lookup Service for Bootstrap system

This service provides a Rest api to query Lookup DB. 

## Interface 

### Get Device 

```
curl --request GET \
  --url http://LOOKUP-URL/devices/51fbba62-c79f-11eb-b8bc-0242ac130003
```
Response:
```
{
  "uuid": "51fbba62-c79f-11eb-b8bc-0242ac130003",
  "orgName": "test-org",
  "certificate": "test-org cert",
  "ip": "192.124.23.1"
}
```

### Add or update Organisation

```
curl --request POST \
  --url http://LOOKUP-URL/orgs/test-org \
  --header 'Content-Type: application/json' \
  --data '{
	"certificate":"test-org cert",
	"ip": "192.124.23.1"
}'
```
Response:
```
{
  "status": "Organisation added or updated"
}
```


### Add or update Node-Organisation mapping 

```
curl --request POST \
  --url http://LOOKUP-URL/devices/51fbba62-c79f-11eb-b8bc-0242ac130005 \
  --header 'Content-Type: application/json' \
  --data '{
	"org":"test-org"
}'
```
Response:
```
{
  "status": "Mapping added"
}
```