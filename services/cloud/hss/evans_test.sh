echo '{ "org": "org-1",  "imsi": { "imsi": "313212" } }' | evans cli call --path ./pb/  --proto hss.proto hss.v1.ImsiService.Add --port 9090
echo '{ "org": "org-1",  "imsi": "31321" }' | evans cli call --path ./pb/  --proto hss.proto hss.v1.ImsiService.Get --port 9090
echo '{ "org": "org-1",  "imsi": "100047265624299" }' | evans cli call --path ./pb/  --proto hss.proto hss.v1.UserService.Add --port 9090
