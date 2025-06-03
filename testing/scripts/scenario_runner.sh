#!/bin/bash
set -euo pipefail
# Function to display all scenarios with descriptions
show_scenarios() {
    echo "Available Scenarios:"
    echo "----------------------------------------"
    echo "1) Node Online"
    echo "   Description: Turns node online"
    echo "   Required: Node ID"
    echo
    echo "2) Node On"
    echo "   Description: Turns on the node"
    echo "   Required: Node ID"
    echo
    echo "3) Node Off"
    echo "   Description: Turns off the node"
    echo "   Required: Node ID"
    echo
    echo "4) Node Restart"
    echo "   Description: Simulates node restart"
    echo "   Required: Node ID"
    echo
    echo "5) Node RF Off"
    echo "   Description: Simulates node RF being turned off"
    echo "   Required: Node ID, Subscriber ID"
    echo
    echo "6) Node RF On"
    echo "   Description: Simulates node RF being turned on"
    echo "   Required: Node ID, Subscriber ID"
    echo
    echo "q) Quit"
    echo "----------------------------------------"
}

# Function to get user input
get_input() {
    local prompt="$1"
    local input=""
    while [ -z "$input" ]; do
        read -p "$prompt" input
        if [ -z "$input" ]; then
            echo "Input cannot be empty. Please try again."
        fi
    done
    echo "$input"
}

API_BASE_URL="http://localhost:8086"  # Adjust this based on your API gateway address
NODE_API_URL="http://localhost:8085"  # Node API URL
NOTIFY_API_URL="http://localhost:8036"
CURL_TIMEOUT="--max-time 120"  # 2 minutes timeout

while true; do
    # Show scenarios at the start and after each execution
    show_scenarios
    
    # Get scenario selection
    read -p "Select a scenario (1-6) or 'q' to quit: " choice
    
    case $choice in
        1) SCENARIO="node_online" ;;
        2) SCENARIO="node_on" ;;
        3) SCENARIO="node_off" ;;
        4) SCENARIO="node_restart" ;;
        5) SCENARIO="node_rf_off" ;;
        6) SCENARIO="node_rf_on" ;;
        q|Q) echo "Exiting..."; exit 0 ;;
        *) echo "Invalid selection. Please try again."; continue ;;
    esac

    echo "Selected scenario: $SCENARIO"
    
    # Get required inputs based on scenario
    NODE_ID=$(get_input "Enter Node ID: ")
    
    if [[ "$SCENARIO" == "node_rf_off" || "$SCENARIO" == "node_rf_on" ]]; then
        SUBSCRIBER_ID=$(get_input "Enter Subscriber ID: ")
    fi

    case "$SCENARIO" in
        "node_online")
            echo "Running node online scenario..."
            echo "Making API call with Node ID: $NODE_ID"
            response=$(curl -s $CURL_TIMEOUT --location "$NODE_API_URL/online?nodeid=$NODE_ID")
            echo "API Response: $response"
            sleep 10
            nr=$(curl -X 'POST' \
                "$NOTIFY_API_URL/v1/notify" \
                -H 'accept: application/json' \
                -H 'Content-Type: application/json' \
                -d "{
                    \"details\": {
                        \"latitude\": 37.7781135,
                        \"longitude\": -121.983609
                    },
                    \"node_id\": \"$NODE_ID\",
                    \"service_name\": \"health\", 
                    \"severity\": \"low\",
                    \"status\": 8100,
                    \"time\": 1733753967,
                    \"type\": \"event\"
                }")
            echo "API Response: $nr"
            ;;
        "node_on"|"node_off"|"node_restart")
            echo "Running node $SCENARIO scenario..."
            echo "Making API call with Node ID: $NODE_ID"
            response=$(curl -s $CURL_TIMEOUT --location "$NODE_API_URL/update?nodeid=$NODE_ID&profile=NORMAL&scenario=$SCENARIO")
            echo "API Response: $response"
            ;;
            
        "node_rf_off"|"node_rf_on")
            echo "Running node $SCENARIO scenario..."
            echo "Making API call with Node ID: $NODE_ID"
            response=$(curl -s $CURL_TIMEOUT --location "$NODE_API_URL/update?nodeid=$NODE_ID&profile=NORMAL&scenario=$SCENARIO")
            echo "Node API Response: $response"
            
            echo "Making API call with Subscriber ID: $SUBSCRIBER_ID"
            sr=$(curl -s $CURL_TIMEOUT --location --request PUT "${API_BASE_URL}/v1/dsubscriber/update" \
                --header 'Content-Type: application/json' \
                --data "{
                    \"iccid\": \"$SUBSCRIBER_ID\",
                    \"profile\": \"normal\",
                    \"scenario\": \"$SCENARIO\"
                }")
            echo "Subscriber API Response: $sr"
            ;;
    esac

    echo "Scenario $SCENARIO completed."
    echo "----------------------------------------"
    echo "Press Enter to continue..."
    read
    clear
done 
