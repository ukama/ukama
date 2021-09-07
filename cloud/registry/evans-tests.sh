go run cmd/server/main.go

echo '{ "name": "org-1", "owner": "24eb9fd7-17c7-4b21-880f-6cefe4c4bb0a" }' | evans cli call --path ./pb/  --proto registry.proto registry.RegistryService.AddOrg --port 9090
echo '{ "name": "org-1" }' | evans cli call --path ./pb/  --proto registry.proto registry.RegistryService.GetOrg --port 9090
echo '{ "name": "org-1" }' | evans cli call --path ./pb/  --proto registry.proto registry.RegistryService.AddNode --port 9090
