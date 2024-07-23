#!/bin/bash

# Define a list of ports
ports=(8080 8081 5041 5042 5043 5044 5045 5046 5047 5048 5049 5050 5051 5052 5053 5054 5055 5056 5057)

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