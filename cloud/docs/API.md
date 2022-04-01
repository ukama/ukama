# Ukama API 


## Authentication 

Here is that script to generate an Session token that could be use to authorize API request.

```
AUTH_URL="https://auth.ukama.com/.api"
# ADD YOU ACCOUNT EMAIL HERE 
EMAIL=YOUR EMAIL HERE
# AND PASSWORD
PASSWORD=YOUR PASSWORK HERT


actionUrl=$(curl -s -X GET -H "Accept: application/json" "$AUTH_URL/self-service/login/api" | jq -r '.ui.action')

echo $actionUrl
TOKEN=$(curl -s -X POST -H  "Accept: application/json" -H "Content-Type: application/json" \
    -d "{\"password_identifier\": \"${EMAIL}\", \"password\": \"${PASSWORD}\", \"method\": \"password\"}" \
    "$actionUrl" | jq -r '.session_token')


# Get default org ID for a user
ORG=$(curl --request GET \
  --url $AUTH_URL/sessions/whoami \
  --header "X-Session-Token: ${TOKEN}" | jq -r '.identity.id')

echo TOKEN: $TOKEN
echo ORG: $ORG

```


Then you can use TOKEN and ORG env vars in you requests.  For example, getting a list of nodes
``` 
curl --request GET \
  --url https://api.dev.ukama.com/orgs/$ORG/nodes \
  --header "Authorization: Bearer ${TOKEN}"
```



### Adding a User 

You can provision an new user on your account via API using `orgs/a32485e4-d842-45da-bf3e-798889c68ad0/users` endpoint.

You will need authentication token and org ID. Refer to [Authentication](#authentication) to find out how to get them.

User in Ukama is a sim owner so every user have a SIM attached. In order to provision a sim we need to have a SIM token. 

```
curl --request POST \
  --url https://api.ukama.com/orgs/$ORG/users \
  --header 'Authorization: Bearer ${TOKEN}' \
  --header 'Content-Type: application/json' \
  --data '{  
	"name": "Joe Doe",	
	"email":"joe@example.com",
	"simToken": "sim token goes here"
}'
```

If request is succeeded then you will get below response:
```
{
	"user": {
		"name": "Joe Doe",
		"email": "joe@example.com",
		"phone": "",
		"uuid": "3732be84-5ce8-4ba1-b8b1-f18480060edd"
	},
	"iccid": "010100001648242526"
}
```
