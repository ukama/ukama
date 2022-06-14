package pkg

import (
	"io/ioutil"
	"log"
	"net/http"
)

type KratosClient interface {
	GetAccountName(networkOwnerId string) (string, error)
}

type kratosClient struct {

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
return string(bytes),errRead
}





