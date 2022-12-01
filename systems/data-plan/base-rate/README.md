# Base rate sub-system

Base rate sub-system is responsibe of:

- Populating rates in DB from CSV
- Provide rates with some optional and require filters
- Provide functionality to get record by id

## File Structure

    .
    └── systems
        └── data-plan
            │── base-rate
            │   ├── cmd
            │   │   ├── server
            │   │   │   └── main.go
            │   │   └── version
            │   │       └── version.go
            │   ├── mocks
            │   │   └── BaseRateRepo.go
            │   ├── pb
            │   │   ├── gen
            │   │   │   └── mocks
            │   │   │       ├── BaseRatesServiceClient.go
            │   │   │       ├── BaseRatesServiceServer.go
            │   │   │       └── UnsafeBaseRatesServiceServer.go
            │   │   ├── rate.pb.go
            │   │   ├── rate_grpc.pb.go
            │   │   └── rate.proto
            │   ├── pkg
            │   │   ├── db
            │   │   │   ├── model.go
            │   │   │   └── rate_repo.go
            │   │   ├── server
            │   │   │   └── rate.go
            │   │   ├── utils
            │   │   │   └── utils.go
            │   │   └── validations
            │   │       └── validations.go
            |   ├── Dockerfile
            │   ├── go.mod
            │   ├── go.sum
            │   └── Makefile
            └── README

- **cmd**: Contains the server and system/sub-system version. Purpose of this file is to initialize the DB and start server. We use `make server` command to run this file.
- **mocks**: This directory contains the auto generagted file which get generated based on `*.proto`. It contains functions which we can use to write test cases.
- **pb**: This directory contains the `*.proto` file. In proto file we define service with all the rpc's and messages.
- **pkg/db**: DB directory under pkg contains 2 files.
`model.go` file contains the db model structure/s.
`*_repo.go` is reponsible of communicating with db using [gorm](https://gorm.io/docs/).
- **pkg/server** This directory contains the file in which all the RPC functions logic is implemented. Those functions call `pkg\*_repo.go` functions to perform db operations.

## RPC Functions

### UploadBaseRates

<img src="https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/digrams/baserate/UploadBaseRates.png" alt="J" width="500"/>

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

**Demo**

<img src="https://user-images.githubusercontent.com/83802574/198561831-0efe13de-0e7e-465f-a6b9-58244296bca5.gif" alt="uploadBaseRates" width="720"/>

---

### GetBaseRates

<img src="https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/digrams/baserate/GetBaseRates.png" alt="J" width="500"/>

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
    [optional] effectiveAt => String
}
```

**Demo**

<img src="https://user-images.githubusercontent.com/83802574/198694692-abed26f1-2ed1-4f4a-8e81-f67a9d0c7270.gif" alt="getBaseRates" width="720"/>

---

### GetBaseRate

<img src="https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/digrams/baserate/GetBaseRate.png" alt="J" width="500"/>

Get base rate service provide functionality to fetch base rate by base rate id.

```proto
service BaseRatesService {
    rpc GetBaseRate(GetBaseRateRequest) returns (GetBaseRateResponse){}
}
```

Function take below argument:

```js
{
    [required] rateId => uint64
}
```

**Demo**

<img src="https://user-images.githubusercontent.com/83802574/198693504-ea7339cb-1795-4c1e-9156-6d383471091a.gif" alt="getBaseRate" width="720"/>

---

## How to use?

Before using the repo make sure below tools are installed:

- Go 1.18
- PostgreSQL
- gRPC client
Then navigate into base-rate directory and run below command:

**To Test**

For unit tests run below commands:

```
make test
```

This commnd will run unit tests under all base-rate directories.

![BaseRateTest](https://user-images.githubusercontent.com/83802574/203062237-1cf8e774-8e0f-4ea6-9e2b-638a9807057e.gif)

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from `pb/*.proto`.

**To Run Server & Test RPC**

```
make server
```

This command will run the server on port `9090` ,and craeate a database name `baserate` with `rates` table under it.

Server is running, Now we can use any gRPC client to intract with RPC handlers. We're using [Evans](https://github.com/ktr0731/evans) here:

```
evans --path /path/to --path . --proto pb/rates.proto --host localhost --port 9090
```

Next run:

```
show rpc
```

This command will show all the available RPC calls under base-rate sub-system.

**UploadBaseRates**

Let's first populate data in our newly created DB using UploadBaseRates RPC.

```
call UploadBaseRates
```

Service takes 3 aurguments **fileURL**, **effectiveAt** & **simType**. For fileURL we can use url of template file existing under `data-plan/docs/template/template.csv`, effectiveAt can be any future UTC formate date and then choose simType.

**GetBaseRates**

To verify that our records are populated we can use GetBaseRates RPC which will return list of rates base on filters provided.
This rpc function takes `country` as required param and some optional arguments **country**,**effectiveAt**,**network**,**simType**,**from**,**to**.

```
call GetBaseRates
```

**GetBaseRate**

To get base rate by id one can user GetBaseRate RPC which return single base rate record.
This rpc takes required argument of `id`.

```
call GetBaseRate
```
