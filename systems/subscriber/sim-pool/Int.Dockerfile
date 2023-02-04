FROM 003664043471.dkr.ecr.us-east-1.amazonaws.com/e2e-base-image:latest 

ENV PROJECT_NAME=subscriber-sim-pool
<<<<<<< HEAD
COPY bin/integration /usr/bin/integration
=======
COPY bin/base-rate /usr/bin/sim-pool
>>>>>>> subscriber-sys_sim-manager
