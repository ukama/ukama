//go:build integration
// +build integration

package db

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/sql"
	"testing"
)

func Test_netRepo_GetNetwork(t *testing.T) {
	dbConf := config.DefaultDatabase()
	db := sql.NewDb(dbConf, true)
	db.Connect()
	r := NewNetRepo(db)
	resp, err := r.Get("network-listener-integration-test-org", "net-1")
	assert.NoError(t, err)
	fmt.Printf("%+v\n", resp)
}
