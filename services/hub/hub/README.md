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

#### Get artifact in tar.gz format
``` 
 curl --request GET \
  --url http://localhost:8080/capps/test-capp/0.0.3.tar.gz \
  --output test-capp-v-0.0.3.tar.gz
```

#### Get chunk index 
``` 
 curl --request GET \
  --url http://localhost:8080/capps/test-capp/0.0.3.caidx \
  --output test-capp-v-0.0.3.caidx
```
#### Get chunk 
```
curl --request GET \
--url https://localhost:8080/chunks/0001/00016cf7c1a372d113c4ba64b56dbd387661d44864a04f59742e3f25a57c594d.cacnk
```
# Contribute
[/docker-compose.yaml](/docker-compose.yaml) start Hub with all required dependencies.

