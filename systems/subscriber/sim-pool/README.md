# Sim Pool sub-system

Sim pool sub-system is responsible of:

- Populating sims data in DB from CSV
- Provide sim stats
- Allows to add slice of sims

## File Structure

    .
    └── systems
        └── subscriber
            │── sim-pool
            │   ├── cmd
            │   │   ├── server
            │   │   │   └── main.go
            │   │   └── version
            │   │       └── version.go
            │   ├── mocks
            │   │   └── SimRepo.go
            │   ├── pb
            │   │   ├── gen
            │   │   │   └── mocks
            │   │   │       ├── SimServiceClient.go
            │   │   │       ├── SimServiceServer.go
            │   │   │       └── UnsafeSimServiceServer.go
            │   │   ├── sim.pb.go
            │   │   ├── sim_grpc.pb.go
            │   │   └── sim.proto
            │   ├── pkg
            │   │   ├── db
            │   │   │   ├── model.go
            │   │   │   └── sim_pool_repo.go
            │   │   ├── server
            │   │   │   └── sim.go
            │   │   └── utils
            │   │       └── utils.go
            |   ├── Dockerfile
            │   ├── go.mod
            │   ├── go.sum
            |   ├── README   
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

### GetStats

GetStats rpc provides return statistics based on sims data.

```proto
service SimService {
    rpc GetStats (GetStatsRequest) returns (GetStatsResponse) {}
}
```

Function takes one optional arguments:

```js
{
    [optional] simType => SimType
}
```

**Demo**

---

### Add

Add rpc provides functionality of adding sims.

```proto
service SimService {
    rpc Add(AddRequest) returns (AddResponse){}
}
```

Function takes these arguments:

```js
{
    [required] sim => []AddSim
}
```

**Demo**

---

### Upload

Upload rpc provides functionality populate batch of sims in DB.

```proto
service SimService {
    rpc Upload(UploadRequest) returns (UploadResponse){}
}
```

Function takes these arguments:

```js
{
    [required] sim => []byte
    [required] simType => SimType
}
```

**Demo**

---

### Delete

Delete service provides functionality of Deleting multiple sims.

```proto
service SimService {
    rpc Delete(DeleteRequest) returns (DeleteResponse){}
}
```

Function takes below argument:

```js
{
    [required] id => []uint64
}
```

**Demo**

---

## How to use?

Before using the repo make sure below tools are installed:

- Go 1.18
- PostgreSQL
- gRPC client
Then navigate into sim-pool directory and run below command:

**To Test**

For unit tests run below commands:

```
make test
```

This command will run unit tests under all sim-pool directories.

**To Generate PB file**

```
make gen
```

This command will generate protobuf files from `pb/*.proto`.

**To Run Server & Test RPC**

```
make server
```

This command will run the server on port `9090`, and create a database named `sim` with `sims` table under it.

Server is running, Now we can use any gRPC client to interact with RPC handlers. We're using [Evans](https://github.com/ktr0731/evans) here:

```
evans --path /path/to --path . --proto pb/sim.proto --host localhost --port 9090
```

Next run:

```
show rpc
```
