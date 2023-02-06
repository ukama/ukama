# Subscriber System

Subscriber system manages all the Subscribers and Sim related flows. Subscriber system is further divided into 4 sub-systems:

- API Gateway
- Sim Pool
- Subscriber
- Sim Manager

## Directory Structure

    .
    └── systems
        └── subscriber
            │── api-gateway
            │   ├── cmd
            │   │   ├── version
            │   ├── mocks
            │   └── pkg
            │       ├── client
            │       └── rest
            │
            │── Sim Pool
            │   ├── cmd
            │   │   ├── server
            │   │   └── version
            │   ├── mocks
            │   ├── pb
            │   │   └──gen
            │   └── pkg
            │       ├── db
            │       ├── server
            │       └── utils
            │
            ├── docs
            │    └── template
            │
            └── README

## Learn more about sub-systems

[API Gateway](https://github.com/ukama/ukama/tree/main/systems/subscriber/api-gateway)

[Sim Pool sub-system](https://github.com/ukama/ukama/tree/main/systems/subscriber/sim-pool)

[Subscriber sub-system](https://github.com/ukama/ukama/tree/main/systems/subscriber/subscriber)

[Sim Manager sub-system](https://github.com/ukama/ukama/tree/main/systems/subscriber/sim-manager)
