# Data Plan System

Data Plan system manage all the sim packages and base rates functionality. Data Plan system will include 2 micro services:

- Base rate microservice
- Packages microservice

## Base rate microservice

Base rate microservice is responsibe of:

- Populating rates in DB from CSV
- Provide rates with some optional and require filters
- Provide functionality to get record by id
  
### Directory Structure

    .
    ├── systems
    │   ├── data-plan
    │   │   │── base-rate
    │   │   │   │── cmd
    │   │   │   │   └── server
    │   │   │   │── pb
    │   │   │   │── pkg
    │   │   │   │   ├── db
    │   │   │   │   ├── config
    │   │   │   │   ├── models
    │   │   │   │   ├── utils
    │   │   │   │   ├── services
    │   │   │   │   └── validations
    │   │   │   └──  template
    └── README

### RPC Function

<u>**UploadBaseRates**</u>
Upload base rates service provide functionality to populate rates from CV file to DB.

```proto
service BaseRatesService {
    rpc UploadBaseRates(UploadBaseRatesRequest) returns (UploadBaseRatesResponse){}
}
```

Function take these arguments:

```js
{
    [required] fileUrl string
    [required] simType string
    [required] effectiveAt string
}
```

<u>**GetBaseRates**</u>
Get base rates service provide functionality to fetch base rates, and filter data on some required and optional arguments.

```proto
service BaseRatesService {
    rpc GetBaseRates (GetBaseRatesRequest) returns (GetBaseRatesResponse) {}
}
```

Function take these arguments:

```js
{
    [required] country string
    [optional] provider string
    [optional] to DateTime
    [optional] from DateTime
    [optional] simType string
}
```

<u>**GetBaseRate**</u>
Get base rate service provide functionality to fetch base rate by base rate id.

```proto
service BaseRatesService {
    rpc GetBaseRate(GetBaseRateRequest) returns (GetBaseRateResponse){}
}
```

Function take below argument:

```js
{
    [required] rateId int32
}
```