KRATOS_URL="https://auth.dev.ukama.com/.api"

actionUrl=$(curl -s -X GET -H "Accept: application/json" "$KRATOS_URL/self-service/login/api" | jq -r '.ui.action')


echo $actionUrl
curl -s -X POST -H  "Accept: application/json" -H "Content-Type: application/json" \
    -d '{"password_identifier": "denis@ukama.com", "password": "Pass2021!", "method": "password"}' \
    "$actionUrl" | jq -r '.session_token'

