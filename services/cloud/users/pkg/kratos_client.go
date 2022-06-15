package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type kratosClient struct {
   
}

type Response struct {
	Traits         struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"traits"`
}
func (i kratosClient) GetAccountName(networkOwnerId string) (string, error) {

resp, err := http.Get(`https://kratos-admin.dev.ukama.com/admin/identities/`+networkOwnerId)
	if err != nil {
	   log.Fatal(err)
      
	}
    
 defer resp.Body.Close()

dataByte,erroBytes:=ioutil.ReadAll(resp.Body)

var result Response

    if err := json.Unmarshal(dataByte, &result); err != nil {  
        fmt.Println("Can not unmarshal JSON")
    }
   
return result.Traits.Name,erroBytes
}



