package pkg

import (
    "log"
    "testing"
    "fmt"
)



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

