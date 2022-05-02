# Ukama Common 
Common code for all Ukama projects 

##How to use

``` shell
# bypass go proxy
go env -w GOPRIVATE=$(go env GOPRIVATE),github.com/ukama
# bypass git login prompt
git config --global url."https://noneed:{GITHUB_TOKEN}@github.com/ukama".insteadOf /"https://github.com/ukama"
# import ukama, you can specify a git commit by adding @COMMIT_HASH to the end
go get github.com/ukama/ukama/services/common
```

Then import in go file like
``` go
package main

import (
	"fmt"

	conf "github.com/ukama/ukama/services/common/config"
)

func main() {
	db := conf.Database{}
	fmt.Printf("Db config %v", db)
}
```
