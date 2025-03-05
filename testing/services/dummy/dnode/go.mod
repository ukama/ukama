module github.com/ukama/ukama/testing/services/dummy/dnode

go 1.24.0

replace github.com/ukama/ukama/testing/services/dummy/dsubscriber => ./

replace github.com/ukama/ukama/testing/services/dummy/dnode => ./

replace github.com/ukama/ukama/testing/common => ../../../common

replace github.com/ukama/ukama/systems/common => ../../../../systems/common

replace github.com/ukama/ukama/systems/services/msgClient => ../../../../systems/services/msgClient
