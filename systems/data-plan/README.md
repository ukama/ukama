# Data Plan System

Data Plan system manage all the sim packages and base rates functionality. Data Plan system will include 2 micro services:

- Base rate microservice
- Packages microservice

## Base rate microservice

Base rate microservice is responsibe of:

- Populating rates in DB from CSV
- Provide rates with some optional and require filters
- Provide functionality to get record by id
  
### Structure

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
    │   │   │   └──  template
    └── README

### RPC Handlers

- UploadBaseRates
- GetBaseRates
- GetBaseRate
