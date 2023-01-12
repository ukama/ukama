
export ENV_SYSTEM_ORG="ukama001" 
export ENV_SYSTEM_NAME="messaging"
export ENV_SYSTEM_ADDR="192.168.0.0"
export ENV_SYSTEM_PORT="8888"
export ENV_SYSTEM_CERT="This is a certificate"
export ENV_INIT_SYSTEM_ADDR="localhost"
export ENV_INIT_SYSTEM_PORT=8081
export ENV_INIT_CLIENT_ADDR="localhost"
export ENV_INIT_CLIENT_PORT=9091

# to add org:
#curl -X PUT http://localhost:8081/v1/orgs/ukama001 \
#	 -H 'Content-Type: application/json' \
#	 -d '{ "certificate": "Hello certificate", "ip": "192.168.0.1", "port": 333
#}'

# then execute the initClient
