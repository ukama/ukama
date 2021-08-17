KRATOS_URL="https://ukama-test.dev.ukama.com/.ory/kratos/public"

actionUrl=$(curl -s -X GET -H "Accept: application/json" "$KRATOS_URL/self-service/login/api" | jq -r '.ui.action')


echo $actionUrl
curl -s -X POST -H  "Accept: application/json" -H "Content-Type: application/json" \
    -d '{"password_identifier": "hz911@mail.ru", "password": "4j4A1t0pXsBY", "method": "password"}' \
    "$actionUrl" | jq -r '.session_token'
