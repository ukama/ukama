# Data Plan System

Data Plan system manages all the sim packages and base rates functionality. Data Plan system will include 2 sub-systems:

- API Gateway
- Base rate sub-system
- Package sub-system

## Directory Structure

    .
    └── systems
        └── data_plan
            │── api-gateway
            │   ├── cmd
            │   │   ├── version
            │   ├── mocks
            │   └── pkg
            │       ├── client
            │       └── rest
            │
            │── base_rate
            │   ├── cmd
            │   │   ├── server
            │   │   └── version
            │   ├── mocks
            │   ├── pb
            │   │   └──gen
            │   └── pkg
            │       ├── db
            │       ├── queue
            │       ├── server
            │       ├── utils
            │       └── validations
            │
            ├── docs
            │    ├── diagrams
            │    └── template
            │
            ├── package
            │   │── cmd
            │   │   ├── server
            │   │   └── version
            │   ├── mocks
            │   ├── pb
            │   │   └──gen
            │   └── pkg
            │       ├── db
            │       ├── server
            │       └── validations
            │
            └── README

## Learn more about sub-systems

[API Gateway](https://github.com/ukama/ukama/tree/main/systems/data_plan/api-gateway)

[Base rate sub-system](https://github.com/ukama/ukama/tree/main/systems/data_plan/base_rate)

[Package sub-system](https://github.com/ukama/ukama/tree/main/systems/data_plan/packge)
