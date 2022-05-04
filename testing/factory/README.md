# Virtual Factory 

## Setup



## Virtual Node Builder

```
TBU
```

## Build

```
make
```

### Unit Test

Update the enviornment variables

GITUSERNAME, GITPASS, DOCKERUSER, DOCKERPASS, BUILDERIMAGE, VMIMAGE, KUBECONFIG, RABBITURI

```
cd virtualfactory/internal/db;
mockgen -source vnode_repo.go -destination mocks/vnode_repo.go -package mockdb
make test
```

## Run

### Prerequisites

**PostgresSQL**

```
docker run --network host --name postgres -p5432:5432 -e POSTGRES_PASSWORD=dev@ukama -d postgres
```

**PGADMIN**

```
 docker pull dpage/pgadmin4
 
 docker run -p 8090:80 -e 'PGADMIN_DEFAULT_EMAIL=vishal@ukama.com' -e 'PGADMIN_DEFAULT_PASSWORD=<password>' -d dpage/pgadmin4
 ```

**RabbitMQ**

```
docker run -d --name rabbit -p 5672:5672 -p 5673:5673 -p 15672:15672 rabbitmq:3-management
```

**NMR**

```
cd ukama/services/factory/nmr
./bin/nmr
```

**ServiceRouter**
```
cd ukama/transport/router/
./servicerouter --host localhost --p 8091 --l debug
```

**VirtualFactory**
```
./factory
```

**RabbitMQ Mock**
```
./mock
```

### Request to build a virtualnode
```
URL: localhost:8091/service?looking_to=create_node&type=hnode&count=1
HTTP Method: POST
```

**Response**
```
Code: 202
Content-Type: application/json"
Body: 
{
    "nodeList": {[
        "uk-sa2130-hnode-a1-f5a4",
    ]}
}
```
