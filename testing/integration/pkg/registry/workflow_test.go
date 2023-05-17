package pkg

import (
	"testing"

	"github.com/tj/assert"
	api "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
)

var OrgName string
var validUserId string

func init() {
	initializeData()
}

func initializeData() {
	OrgName = "milky-way"
	validUserId = "c9647e7a-8967-4978-b512-38a35899f32d"
}

func TestOwnerAddsValidMember(t *testing.T) {
	host := "http://localhost:8082"
	registryClient := NewRegistryClient(host)

	reqGetOrg := api.GetOrgRequest{OrgName: OrgName}

	orgResp, err := registryClient.GetOrg(reqGetOrg)
	assert.NoError(t, err)
	assert.NotNil(t, orgResp)
	assert.Equal(t, OrgName, orgResp.Org.Name)

	reqGetMember := api.GetMemberRequest{
		OrgName:  OrgName,
		UserUuid: validUserId}

	gmResp, err := registryClient.GetMember(reqGetMember)
	assert.Error(t, err)
	assert.Nil(t, gmResp)

	reqAddMember := api.MemberRequest{
		OrgName:  OrgName,
		UserUuid: validUserId}

	amResp, err := registryClient.AddMember(reqAddMember)
	assert.NoError(t, err)
	assert.NotNil(t, amResp)
}
