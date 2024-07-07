#!/bin/bash

# Define a list of ports
ports=(4000 4002 4003 4004 4005 4006 4007 4008 4009 4010 4011 4012 4013 4014 4015 4016 4017 4018)

# Loop through each port in the list
for PORT in "${ports[@]}"; do
    # Check if the port is in use and get the PID of the process using it
    PID=$(lsof -ti:$PORT)

    if [ -z "$PID" ]; then
        echo "Port $PORT is available."
    else
        echo "Port $PORT is in use by PID $PID. Attempting to kill..."
        kill -9 $PID
        if [ $? -eq 0 ]; then
            echo "Successfully killed process on port $PORT."
        else
            echo "Failed to kill process on port $PORT. You might need to run the script as root."
        fi
    fi
done

sleep 5