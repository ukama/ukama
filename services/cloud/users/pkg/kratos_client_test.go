package pkg

import (
    "log"
    "testing"
    "io/ioutil"
	"net/http"
    "fmt"
)

type kratosClient struct {
   
}

  
type UserOject struct {
    Name string
    Email  string
}
func (i kratosClient) GetAccountName(networkOwnerId string) (string, error) {

resp, err := http.Get(`https://kratos-admin.dev.ukama.com/admin/identities/`+networkOwnerId)
	if err != nil {
	   log.Fatal(err)
      
	}
bytes ,errRead :=ioutil.ReadAll(resp.Body)

if errRead!=nil{
	log.Fatal(errRead)
}
fmt.Println(string(bytes.schema_url))

return string(bytes),errRead

}

func Test_GetAccounName(t *testing.T)  {
   
    c:= kratosClient {}
    usr, err := c.GetAccountName("a32485e4-d842-45da-bf3e-798889c68ad0")
  
    if err !=nil{
        log.Fatal(err)
    }
   
   if usr !="Ukama Dev Team"{
    log.Fatal("Failed to find the user")
   }
}

