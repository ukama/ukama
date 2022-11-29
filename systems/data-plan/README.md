# Data Plan System

Data Plan system manages all the sim packages and base rates functionality. Data Plan system will include 2 micro services:

- Base rate sub-system
- Packages sub-system

## Base rate sub-system

Base rate sub-system is responsibe of:

- Populating rates in DB from CSV
- Provide rates with some optional and required filters
- Provide functionality to get record by id

### Directory Structure

    .
    └── systems
        └── data-plan
            │── base-rate
            │   │── cmd
            |   |   ├── listner
            |   |   ├── server
            │   │   └── version
            |   │── mocks
            │   │── pb
            |   |   └──gen
            │   └── pkg
            │       ├── db
            │       ├── queue
            │       ├── server
            │       ├── utils
            │       └── validations
            ├──  docs
            |    ├── digrams
            |    └── template
            └── README

### RPC Function

**UploadBaseRates**

<img src="https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/digrams/UploadBaseRates.png" alt="upload-base-rate" width="500"/>

Upload base rates service provides functionality to populate rates from CSV file to DB.

```proto
service BaseRatesService {
    rpc UploadBaseRates(UploadBaseRatesRequest) returns (UploadBaseRatesResponse){}
}
```

Function takes these arguments:

```js
{
    [required] fileUrl => String
    [required] simType => String
    [required] effectiveAt => String
}
```

**Demo**

<img src="https://user-images.githubusercontent.com/83802574/198561831-0efe13de-0e7e-465f-a6b9-58244296bca5.gif" alt="uploadBaseRates" width="720"/>

---

**GetBaseRates**

<img src="https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/digrams/GetBaseRates.png" alt="get-base-rates" width="500"/>

Get base rates service provides functionality to fetch base rates, and filter data on some required and optional arguments.

```proto
service BaseRatesService {
    rpc GetBaseRates (GetBaseRatesRequest) returns (GetBaseRatesResponse) {}
}
```

Function takes these arguments:

```js
{
    [required] country => String
    [optional] provider => String
    [optional] to => DateTime
    [optional] from => DateTime
    [optional] simType => String
    [optional] effectiveAt => String
}
```

**Demo**

<img src="https://user-images.githubusercontent.com/83802574/198694692-abed26f1-2ed1-4f4a-8e81-f67a9d0c7270.gif" alt="getBaseRates" width="720"/>

---

**GetBaseRate**

<img src="https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/digrams/GetBaseRate.png" alt="get-base-rate" width="500"/>

Get base rate service provides functionality to fetch base rate by base rate id.

```proto
service BaseRatesService {
    rpc GetBaseRate(GetBaseRateRequest) returns (GetBaseRateResponse){}
}
```

Function takes below argument:

```js
{
    [required] rateId => Int32
}
```

**Demo**

<img src="https://user-images.githubusercontent.com/83802574/198693504-ea7339cb-1795-4c1e-9156-6d383471091a.gif" alt="getBaseRate" width="720"/>

---

### How to use?

Before using the repo make sure below tools are installed:

- Go 1.18
- PostgreSQL
- gRPC client

Then navigate into base-rate directory and run below command:

```
make gen
```

This command will generate protobuf and mocks from `pb/rate.proto`.

```
make test
```

This command will run the unit tests.

```
make server
```

This command will run the server on port `9090`, and create a database named `baserate` with `rates` table under it.

Server is running, Now we can use any gRPC client to interact with RPC handlers. We're using [Evans](https://github.com/ktr0731/evans) here:

```
evans --path /path/to --path . --proto pb/rate.proto --host localhost --port 9090
```

Next run:

```
show rpc
```

This command will show all the available RPC calls under base-rate service.

**UploadBaseRates**

Let's first populate data in our newly created DB using UploadBaseRates RPC.

```
call UploadBaseRates
```

Service takes 3 aurguments **fileURL**, **effectiveAt** & **simType**. For fileURL we can use url of template file existing under `data-plan/docs/template/template.csv`, effectiveAt can be any future UTC formated date and then choose simType.

**GetBaseRates**

To verify that our records are populated we can use GetBaseRates RPC which will return a list of rates based on filters provided.
This rpc function takes `country` as required param and some optional arguments **country**,**effectiveAt**,**network**,**simType**,**from**,**to**.

```
call GetBaseRates
```

**GetBaseRate**

To get base rate by id, one can use GetBaseRate RPC which returns a single base rate record.
This rpc takes required argument of `id`.

```
call GetBaseRate
```
