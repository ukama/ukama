AUTH_URL="https://auth.dev.ukama.com/.api"
PASSWORD=Pass2020!
EMAIL=test@ukama.com

actionUrl=$(curl -s -X GET -H "Accept: application/json" "$AUTH_URL/self-service/login/api" | jq -r '.ui.action')

echo $actionUrl
TOKEN=$(curl -s -X POST -H  "Accept: application/json" -H "Content-Type: application/json" \
    -d "{\"password_identifier\": \"${EMAIL}\", \"password\": \"${PASSWORD}\", \"method\": \"password\"}" \
    "$actionUrl" | jq -r '.session_token')

# Get default org ID for a user
ORG=$(curl --request GET \
  --url $AUTH_URL/sessions/whoami \
  --header "X-Session-Token: ${TOKEN}" | jq -r '.identity.id')

# Call nodes endpoint
curl --request GET \
  --url https://api.dev.ukama.com/orgs/$ORG/nodes \
  --header "Authorization: Bearer ${TOKEN}"