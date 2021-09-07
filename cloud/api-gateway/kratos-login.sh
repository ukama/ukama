KRATOS_URL="https://app.dev.ukama.com/.ory/kratos/public"

actionUrl=$(curl -s -X GET -H "Accept: application/json" "$KRATOS_URL/self-service/login/api" | jq -r '.ui.action')


echo $actionUrl
curl -s -X POST -H  "Accept: application/json" -H "Content-Type: application/json" \
    -d '{"password_identifier": "denis@ukama.com", "password": "123456", "method": "password"}' \
    "$actionUrl" | jq -r '.session_token'

