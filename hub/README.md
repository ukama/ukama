# Artifacts Hub
Ukama Artifacts Hub is a web application that allows users to upload and share artifacts. 
It exposes rest interface which can be examined via `/swagger` or `/openapi.json` endpoint

### Upload artifact
``` bash
curl --request PUT \
  --url http://$HUB_HOST/capps/test-capp/0.0.3 \
  --header 'Content-Type: application/gzip' \
   --data-binary "@path/to/file"
```
### Download artifact

``` 
 curl --request GET \
  --url http://localhost:8080/capps/test-capp/0.0.3 \
  --output test-capp-v-0.0.3.tar.gz
```

# Contribute
[/docker-compose.yaml](/docker-compose.yaml) start Hub with all required dependencies.

