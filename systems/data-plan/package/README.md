# Package sub-system

Package sub-system provide CRUD options to organization. Sub-system provide following rpc's:

- Add package under organization
- Update organization package
- Get packages by `id` & `orgId`
- Delete package uder organization

## File Structure

    .
    └── systems
        └── package
            │── base-rate
            │   ├── cmd
            │   │   ├── server
            │   │   │   └── main.go
            │   │   └── version
            │   │       └── version.go
            │   ├── mocks
            │   │   └── PackageRepo.go
            │   ├── pb
            │   │   ├── gen
            │   │   │   └── mocks
            │   │   │       ├── PackageServiceClient.go
            │   │   │       ├── PackageServiceServer.go
            │   │   │       └── UnsafePackageServiceServer.go
            │   │   ├── package.pb.go
            │   │   ├── package_grpc.pb.go
            │   │   └── package.proto
            │   ├── pkg
            │   │   ├── db
            │   │   │   ├── model.go
            │   │   │   └── package_repo.go
            │   │   ├── server
            │   │   │   └── package.go
            │   │   └── validations
            │   │   │   └── validations.go
            │   │   ├── config.go
            │   │   └── global.go
            |   ├── Dockerfile
            │   ├── go.mod
            │   ├── go.sum
            │   └── Makefile
            └── README

- **cmd**: Contains the server and system/sub-system version. Purpose of this file is to initialize the DB and start server. We use `make server` command to run this file.
- **mocks**: This directory contains the auto generated file which get generated based on `*.proto`. It contains functions which we can use to write test cases.
- **pb**: This directory contains the `*.proto` file. In proto file we define service with all the rpc's and messages.
- **pkg/db**: DB directory under pkg contains 2 files.
`model.go` file contains the db model structure/s.
`*_repo.go` is responsible of communicating with db using [gorm](https://gorm.io/docs/).
- **pkg/server** This directory contains the file in which all the RPC functions logic is implemented. Those functions call `pkg\*_repo.go` functions to perform db operations.

## RPC Functions

### Get Packages

<img src="https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/digrams/package/GetPackages.png" alt="J" width="500"/>

```proto
service PackagesService {
    rpc Get(GetPackagesRequest) returns (GetPackagesResponse) {}
}
```

Function takes below argument:

```js
{
    [required] orgId => UInt64
    [optional] id => UInt64
}
```

**Demo**

![getPackage](https://user-images.githubusercontent.com/15526332/202410788-44f74507-94e7-42f5-98a5-adab48e7d50c.gif)

---

### Add Package

<img src="https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/digrams/package/AddPackage.png" alt="J" width="500"/>

```proto
service PackagesService {
    rpc Add(CreatePackageRequest) returns (CreatePackageResponse){}
}
```

Function takes below argument:

```js
{
    [required] orgId => UInt64
    [required] name => String
    [required] duration => UInt64
    [required] org_rates_id => UInt64
    [optional] active => Boolean
    [optional] sim_type => String
    [optional] sms_volume => Int64
    [optional] data_volume => Int64
    [optional] voice_volume => Int64
}
```

**Demo**

![AddPackage](https://user-images.githubusercontent.com/15526332/202414408-fa235cd8-a750-40c4-84b4-4d4d46efe01d.gif)

---

### UpdatePackage

<img src="https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/digrams/package/UpdatePackage.png" alt="J" width="500"/>

```proto
service PackagesService {
    rpc Update(UpdatePackageRequest) returns (UpdatePackageResponse){}
}
```

Function takes below argument:

```js
{
    [required] id => UInt64
    [optional] name => String
    [optional] duration => UInt64
    [optional] org_rates_id => UInt64
    [optional] active => Boolean
    [optional] sim_type => String
    [optional] sms_volume => Int64
    [optional] data_volume => Int64
    [optional] voice_volume => Int64
}
```

**Demo**

![UpdatePackage](https://user-images.githubusercontent.com/15526332/202415517-2cb9a2d1-39bb-4912-93cd-e3d5631c52ea.gif)

---

### DeletePackage

<img src="https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/digrams/package/DeletePackage.png" alt="J" width="500"/>

```proto
service PackagesService {
    rpc Delete(DeletePackageRequest) returns (DeletePackageResponse){}
}
```

Function takes below argument:

```js
{
    [required] id => UInt64
    [required] orgId => UInt64
}
```

**Demo**

![deletePackage](https://user-images.githubusercontent.com/83802574/202478420-1f144efa-e356-480d-ac27-d09a389820e9.gif)

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

This command will run unit tests under all base-rate directories.

![pkg_test](https://user-images.githubusercontent.com/83802574/203065220-496256f6-1ec4-4a78-8a0c-36d1f200f862.gif)

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from `pb/*.proto`.

**To Run Server & Test RPC**

```
make server
```

This command will run the server on port `9090`, and create a database name `package` with `packages` table under it.

Server is running, Now we can use any gRPC client to interact with RPC handlers. We're using [Evans](https://github.com/ktr0731/evans) here:

```
evans --path /path/to --path . --proto pb/package.proto --host localhost --port 9090
```

Next run:

```
show rpc
```

This command will show all the available RPC calls under package sub-system.

→ **Add**

Let's first populate data in our newly created DB using Add RPC.

```
call Add
```

→ **Get**

Get package

```
call Get
```

→ **Update**

Update package

```
call Update
```

→ **Delete**

Delete package

```
call Delete
```