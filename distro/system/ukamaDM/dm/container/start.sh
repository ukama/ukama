#!/bin/sh

echo "Starting script."

sleep 45;

/sbin/lwm2mClient -4 -f /etc/lwclient/client.toml
