package pkg

import (
    "log"
    "testing"
    "io/ioutil"
	"net/http"
    "fmt"
    "encoding/json"
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

func Test_GetAccounName(t *testing.T)  {
   
    c:= kratosClient {}
    usr, err := c.GetAccountName("a32485e4-d842-45da-bf3e-798889c68ad0")
   
    if len(usr) == 0{
        fmt.Println("Cannot find user with that ID")
    }
    if err !=nil{
        log.Fatal(err)
    }
   if usr !="Ukama Dev Team"{
    log.Fatal("Failed to find the user")
   }
   
}

