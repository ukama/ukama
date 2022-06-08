go run cmd/server/main.go

echo '{ "name": "org-1", "owner": "24eb9fd7-17c7-4b21-880f-6cefe4c4bb0a" }' | evans cli call  --path ../../common/pb/dep,./pb/  --proto network.proto  ukama.network.v1.NetworkService.AddOrg --port 9090
echo '{ "name": "org-1" }' | evans cli call  --path ../../common/pb/dep,./pb/  --proto network.proto  ukama.network.v1.NetworkService.GetOrg --port 9090
echo '{ "orgName": "org-1", "node": { "nodeId":"uk-ab0001-hnode-a1-0001", "name":"test"  }}' | evans cli call  --path ../../common/pb/dep,./pb/  --proto network.proto  ukama.network.v1.NetworkService.AddNode --port 9090
echo '{ "orgName": "org-1", "node": { "nodeId":"uk-sa2209-comv1-a1-ee58", "name":"test1"  }}' | evans cli call  --path ../../common/pb/dep,./pb/  --proto network.proto  ukama.network.v1.NetworkService.AddNode --port 9090
echo '{ "orgName": "org-1", "node": { "nodeId":"uk-sa2209-anode-a1-070d", "name":"test2"  }}' | evans cli call  --path ../../common/pb/dep,./pb/  --proto network.proto  ukama.network.v1.NetworkService.AddNode --port 9090
echo '{ "parentNodeId": "uk-sa2209-comv1-a1-ee58", "attachedNodeIds": [ "uk-sa2209-anode-a1-070d" ] }' | evans cli call  --path ../../common/pb/dep,./pb/  --proto network.proto  ukama.network.v1.NetworkService.AttachNodes  --port 9090