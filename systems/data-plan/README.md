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
    │   │   │   └──  static
    |   |   |        ├── digrams
    │   │   │        └── template
    └── README

### RPC Function

**UploadBaseRates**

<img src="https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/static/digrams/UploadBaseRates.png" alt="J" width="500"/>

Upload base rates service provide functionality to populate rates from CV file to DB.

```proto
service BaseRatesService {
    rpc UploadBaseRates(UploadBaseRatesRequest) returns (UploadBaseRatesResponse){}
}
```

Function take these arguments:

```js
{
    [required] fileUrl String
    [required] simType String
    [required] effectiveAt String
}
```

________________
**GetBaseRates**

<img src="https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/static/digrams/GetBaseRates.png" alt="J" width="500"/>

Get base rates service provide functionality to fetch base rates, and filter data on some required and optional arguments.

```proto
service BaseRatesService {
    rpc GetBaseRates (GetBaseRatesRequest) returns (GetBaseRatesResponse) {}
}
```

Function take these arguments:

```js
{
    [required] country String
    [optional] provider String
    [optional] to DateTime
    [optional] from DateTime
    [optional] simType String
}
```

________________
**GetBaseRate**

<img src="https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/static/digrams/GetBaseRate.png" alt="J" width="500"/>

Get base rate service provide functionality to fetch base rate by base rate id.

```proto
service BaseRatesService {
    rpc GetBaseRate(GetBaseRateRequest) returns (GetBaseRateResponse){}
}
```

Function take below argument:

```js
{
    [required] rateId Int32
}
```
________________

### How to use?

