# Data Plan System

Data Plan system manages all the sim packages and base rates functionality. Data Plan system will include 2 sub-systems:

- API Gateway
- Base rate sub-system
- Package sub-system

## Directory Structure

    .
    └── systems
        └── data-plan
            │── api-gateway
            │   ├── cmd
            │   │   ├── version
            │   ├── mocks
            │   └── pkg
            │       ├── client
            │       └── rest
            │
            │── base-rate
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

[API Gateway](https://github.com/ukama/ukama/tree/main/systems/data-plan/api-gateway)

[Base rate sub-system](https://github.com/ukama/ukama/tree/main/systems/data-plan/base-rate)

[Package sub-system](https://github.com/ukama/ukama/tree/main/systems/data-plan/packge)
