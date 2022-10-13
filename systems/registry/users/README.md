# Users Service

Users service is responsible for managing users and serve as a proxy for managing sim cards. It stores the user details and ICCID assigned to them. 


## Sim Manager
Users service relies on [Sim Card Manager service](pb/client/sim_manager.proto) to source the information about sim cards(IMSI, data usage) and control their properties such as services availability.
