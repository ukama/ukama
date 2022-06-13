package pkg

import (
    "log"
)
type KratosClient interface{
	GetAccountName(networkOwnerId string )(string ,error) 
   
}
type kratosClient struct{

}
func (k *kratosClient) GetAccountName(networkOwnerId string )(string ,error) {
    c:= kratosClient {}
    usr, err := c.GetAccountName("a32485e4-d842-45da-bf3e-798889c68ad0")
   if usr !="Ukama Dev Team"{
    log.Fatal(err)
   }
  return usr, err
}

