#!/bin/bash
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#USAGE: kisckstart.sh  

# Deal with on-boot processes.

echo "Starting on-boot."
supervisorctl start on-boot:*

# Check for noded to move in running state.
while ! supervisorctl status on-boot:noded_latest | grep -q 'RUNNING'; do sleep 2; done

# Start the bootstrap process
supervisorctl start bootstrap_latest

# Check for the oneshot bootstrap process to complete.
while ! supervisorctl status bootstrap_latest | grep -q 'EXITED'; do sleep 10; done

# Start Other processes
supervisorctl start meshd_latest

#Add new services under service group
#supervisorctl start service:*

