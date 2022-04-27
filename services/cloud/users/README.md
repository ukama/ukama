# Users Service

Users service is responsible for managing users and serve as a proxy for managing sim cards. It stores the user details and ICCID assigned to them. 
It relies on [Sim Card Manager](pb/client/sim_manager.proto) to source the information about sim cards and control their properties such as available service.

