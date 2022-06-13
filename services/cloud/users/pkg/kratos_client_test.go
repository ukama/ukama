package pkg

import (
    "log"
    "testing"
)
type KratosClient interface{
	GetAccountName(networkOwnerId string) (string, error)
   
}
type kratosClient struct{
   
}
func (k *kratosClient) GetAccountName(networkOwnerId string) (string, error){
 
}

func Test_GetAccounName(t *testing.T)  {
    c:= kratosClient {}
    usr, err := c.GetAccountName("a32485e4-d842-45da-bf3e-798889c68ad0")
   if usr !="Ukama Dev Team"{
    log.Fatal(err)
   }
}

