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
    │   │   │   └──  docs
    |   |   |        ├── digrams
    │   │   │        └── template
    └── README

### RPC Function

**UploadBaseRates**

<img src="https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/docs/digrams/UploadBaseRates.png" alt="J" width="500"/>

Upload base rates service provide functionality to populate rates from CV file to DB.

```proto
service BaseRatesService {
    rpc UploadBaseRates(UploadBaseRatesRequest) returns (UploadBaseRatesResponse){}
}
```

Function take these arguments:

```js
{
    [required] fileUrl => String
    [required] simType => String
    [required] effectiveAt => String
}
```

________________
**GetBaseRates**

<img src="https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/docs/digrams/GetBaseRates.png" alt="J" width="500"/>

Get base rates service provide functionality to fetch base rates, and filter data on some required and optional arguments.

```proto
service BaseRatesService {
    rpc GetBaseRates (GetBaseRatesRequest) returns (GetBaseRatesResponse) {}
}
```

Function take these arguments:

```js
{
    [required] country => String
    [optional] provider => String
    [optional] to => DateTime
    [optional] from => DateTime
    [optional] simType => String
}
```

________________
**GetBaseRate**

<img src="https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/docs/digrams/GetBaseRate.png" alt="J" width="500"/>

Get base rate service provide functionality to fetch base rate by base rate id.

```proto
service BaseRatesService {
    rpc GetBaseRate(GetBaseRateRequest) returns (GetBaseRateResponse){}
}
```

Function take below argument:

```js
{
    [required] rateId => Int32
}
```
________________

### How to use?

Before using the repo make sure below tools are installed:

- Go 1.18
- PostgreSQL

Then navigate into base-rate directory and run below command:

```
make server
```

This command will run the server and create database named `baserate` with `rates` table under it.

Server is running, Now we can use any gRPC client to intract with RPC handlers. I'm using [Evans](https://github.com/ktr0731/evans) here:

```
evans --path /path/to --path . --proto pb/rate.proto --host localhost --port 9090
```

Next run:

```
show rpc
```

This command will show all the available RPC calls under base-rate service.

**UploadBaseRates**

Let's first populate data in out newly created DB using UploadBaseRates RPC.

```
call UploadBaseRates
````

Service takes 3 aurguments **fileURL**, **effectiveAt** & **simType**. For fileURL we can use url of template file existing under `data-plan/docs/template/template.csv`, effectiveAt can be any future UTC formate date and then choose simType.

**GetBaseRates**

To verify that our records are populated we can use GetBaseRates RPC which will return list of rates base on filters provided.

```
call GetBaseRates
```

**GetBaseRate**

To get base rate by id one can user GetBaseRate RPC which will return base rate record.
This GetBaseRate takes required argument of `id`.

```
call GetBaseRate
```