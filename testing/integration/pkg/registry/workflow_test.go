package pkg

import (
	"strings"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/tj/assert"
	api "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
)

func TestAddMemberWorkflows(t *testing.T) {
	host := "http://localhost:8082"
	registryClient := NewRegistryClient(host)

	t.Run("Owner adds valid member should succeed", func(t *testing.T) {
		orgName := "milky-way"
		nomMemberId := "c9647e7a-8967-4978-b512-38a35899f32d"

		reqGetOrg := api.GetOrgRequest{OrgName: orgName}

		orgResp, err := registryClient.GetOrg(reqGetOrg)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, orgName, orgResp.Org.Name)

		reqGetMember := api.GetMemberRequest{
			OrgName:  orgName,
			UserUuid: nomMemberId}

		gmResp, err := registryClient.GetMember(reqGetMember)
		assert.Error(t, err)
		assert.Nil(t, gmResp)

		// make sure the owner is the request executor
		reqAddMember := api.MemberRequest{
			OrgName:  orgName,
			UserUuid: nomMemberId}

		amResp, err := registryClient.AddMember(reqAddMember)
		assert.NoError(t, err)
		assert.NotNil(t, amResp)
	})

	t.Run("Non owner adds member should fail", func(t *testing.T) {
		orgName := "saturn"
		member := "c9647e7a-8967-4978-b512-38a35899f32d"
		nonMemberId := "ec4c897e-cc78-43c7-aee3-871a956808c4"

		reqGetOrg := api.GetOrgRequest{OrgName: orgName}

		orgResp, err := registryClient.GetOrg(reqGetOrg)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, orgName, orgResp.Org.Name)

		owner := orgResp.Org.Owner
		assert.NotEqual(t, owner, member)

		reqGetMember := api.GetMemberRequest{
			OrgName:  orgName,
			UserUuid: member}

		gmResp, err := registryClient.GetMember(reqGetMember)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, member, gmResp.Member.Uuid)

		reqGetNonMember := api.GetMemberRequest{
			OrgName:  orgName,
			UserUuid: nonMemberId}

		gnResp, err := registryClient.GetMember(reqGetNonMember)
		assert.Error(t, err)
		assert.Nil(t, gnResp)

		// make sure the member is the request executor
		reqAddMember := api.MemberRequest{
			OrgName:  orgName,
			UserUuid: nonMemberId}

		amResp, err := registryClient.AddMember(reqAddMember)
		assert.Error(t, err)
		assert.Nil(t, amResp)
	})
}

func TestUpdateMemberWorkflows(t *testing.T) {
	host := "http://localhost:8082"
	registryClient := NewRegistryClient(host)

	t.Run("Owner updates valid member should succeed", func(t *testing.T) {
		orgName := "saturn"
		member := "c9647e7a-8967-4978-b512-38a35899f32d"

		reqGetOrg := api.GetOrgRequest{OrgName: orgName}

		orgResp, err := registryClient.GetOrg(reqGetOrg)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, orgName, orgResp.Org.Name)

		reqGetMember := api.GetMemberRequest{
			OrgName:  orgName,
			UserUuid: member}

		gmResp, err := registryClient.GetMember(reqGetMember)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, member, gmResp.Member.Uuid)

		// make sure the owner is the request executor
		reqUpdateMember := api.UpdateMemberRequest{
			OrgName:       orgName,
			UserUuid:      member,
			IsDeactivated: true}

		err = registryClient.UpdateMember(reqUpdateMember)
		assert.NoError(t, err)
	})

	t.Run("Non owner updates member should fail", func(t *testing.T) {
		orgName := "saturn"
		member := "c9647e7a-8967-4978-b512-38a35899f32d"
		nonOwner := "022586f0-2d0f-4b30-967d-2156574fece4"

		reqGetOrg := api.GetOrgRequest{OrgName: orgName}

		orgResp, err := registryClient.GetOrg(reqGetOrg)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, orgName, orgResp.Org.Name)

		owner := orgResp.Org.Owner
		assert.NotEqual(t, owner, nonOwner)

		reqGetMember := api.GetMemberRequest{
			OrgName:  orgName,
			UserUuid: member}

		gmResp, err := registryClient.GetMember(reqGetMember)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, member, gmResp.Member.Uuid)

		reqGetNonOwner := api.GetMemberRequest{
			OrgName:  orgName,
			UserUuid: nonOwner}

		gnResp, err := registryClient.GetMember(reqGetNonOwner)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, nonOwner, gnResp.Member.Uuid)

		// make sure non owner is the request executor
		reqUpdateMember := api.UpdateMemberRequest{
			OrgName:       orgName,
			UserUuid:      member,
			IsDeactivated: true}

		err = registryClient.UpdateMember(reqUpdateMember)
		assert.Error(t, err)
	})
}

func TestAddNetworkWorkflows(t *testing.T) {
	host := "http://localhost:8082"
	registryClient := NewRegistryClient(host)

	t.Run("Owner-admin adds new network to eshould succeed", func(t *testing.T) {
		orgName := "saturn"

		// or use an admin member_id instead of owner_id
		owner := "08a594d7-a292-43cf-9652-54785b03f48f"

		netName := strings.ToLower(faker.FirstName()) + "-net"

		reqGetOrg := api.GetOrgRequest{OrgName: orgName}

		orgResp, err := registryClient.GetOrg(reqGetOrg)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, orgName, orgResp.Org.Name)
		assert.Equal(t, owner, orgResp.Org.Owner)

		// make sure the owner or admin is the request executor
		reqAddNetwork := api.AddNetworkRequest{
			OrgName: orgName,
			NetName: netName}

		anResp, err := registryClient.AddNetwork(reqAddNetwork)
		assert.NoError(t, err)
		assert.NotNil(t, anResp)
	})

	t.Run("Non owner-admin adds new network should fail", func(t *testing.T) {
		orgName := "saturn"
		member := "c9647e7a-8967-4978-b512-38a35899f32d"
		netName := strings.ToLower(faker.FirstName()) + "-net"

		reqGetOrg := api.GetOrgRequest{OrgName: orgName}

		orgResp, err := registryClient.GetOrg(reqGetOrg)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, orgName, orgResp.Org.Name)

		// make sure the member is not owner
		owner := orgResp.Org.Owner
		assert.NotEqual(t, owner, member)

		reqGetMember := api.GetMemberRequest{
			OrgName:  orgName,
			UserUuid: member}

		gmResp, err := registryClient.GetMember(reqGetMember)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, member, gmResp.Member.Uuid)

		// make sure the member is not admin
		// something like assertNotEqual(member.Role.String(), "admin")

		// make sure the member is the request executor
		reqAddNetwork := api.AddNetworkRequest{
			OrgName: orgName,
			NetName: netName}

		amResp, err := registryClient.AddNetwork(reqAddNetwork)
		assert.Error(t, err)
		assert.Nil(t, amResp)
	})

	t.Run("Owner-admin adds new network to non existing org should fail", func(t *testing.T) {
		orgName := "saturn"
		missingOrgName := "non-existing-org"

		// or use an admin member_id instead of owner_id
		owner := "08a594d7-a292-43cf-9652-54785b03f48f"

		netName := strings.ToLower(faker.FirstName()) + "-net"

		reqGetOrg := api.GetOrgRequest{OrgName: missingOrgName}

		orgResp, err := registryClient.GetOrg(reqGetOrg)
		assert.Error(t, err)
		assert.Nil(t, orgResp)

		reqGetOrg = api.GetOrgRequest{OrgName: orgName}

		orgResp, err = registryClient.GetOrg(reqGetOrg)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, orgName, orgResp.Org.Name)
		assert.Equal(t, owner, orgResp.Org.Owner)

		// make sure the owner or admin is the request executor
		reqAddNetwork := api.AddNetworkRequest{
			OrgName: missingOrgName,
			NetName: netName}

		anResp, err := registryClient.AddNetwork(reqAddNetwork)
		assert.Error(t, err)
		assert.Nil(t, anResp)
	})
}

func TestUpdateNetworkWorkflows(t *testing.T) {
	host := "http://localhost:8082"
	registryClient := NewRegistryClient(host)

	t.Run("Owner-admin updates network name should succeed", func(t *testing.T) {
		orgName := "saturn"

		oldNetId := "b884485f-cb43-44b1-be57-0b777b154ff2"
		newNetName := strings.ToLower(faker.FirstName()) + "-net"

		owner := "08a594d7-a292-43cf-9652-54785b03f48f"

		reqGetOrg := api.GetOrgRequest{OrgName: orgName}

		orgResp, err := registryClient.GetOrg(reqGetOrg)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, orgName, orgResp.Org.Name)
		assert.Equal(t, owner, orgResp.Org.Owner)

		reqGetNet := api.GetNetworkRequest{NetworkId: oldNetId}

		netResp, err := registryClient.GetNetwork(reqGetNet)
		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, orgResp.Org.Id, netResp.Network.OrgId)
		assert.NotEqual(t, newNetName, netResp.Network.Name)

		// make sure the owner or admin is the request executor
		// in order for this test to even compile, we need to implement
		// missing endppins (PATCH /v1/networks) and missing APIs
		// requests wrappers (RegistryClient.UpdateNetwork and
		// systems/registry/api-gateway/pkg/rest/api.go for
		// api.UpdateNetworkRequest)
		reqUpdateNetwork := api.UpdateNetworkRequest{
			NetId:   oldNetId,
			NetName: newNetName}

		unResp, err := registryClient.UpdateNetwork(reqUpdateNetwork)
		assert.NoError(t, err)
		assert.NotNil(t, unResp)
	})

	t.Run("Non owner-admin updates network name should fail", func(t *testing.T) {
		orgName := "saturn"

		oldNetId := "b884485f-cb43-44b1-be57-0b777b154ff2"
		newNetName := strings.ToLower(faker.FirstName()) + "-net"

		member := "c9647e7a-8967-4978-b512-38a35899f32d"

		reqGetOrg := api.GetOrgRequest{OrgName: orgName}

		orgResp, err := registryClient.GetOrg(reqGetOrg)
		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, orgName, orgResp.Org.Name)

		// make sure the member is not admin, nor owner
		// something like assertNotEqual(member.Role.String(), "admin")
		assert.NotEqual(t, member, orgResp.Org.Owner)

		reqGetNet := api.GetNetworkRequest{NetworkId: oldNetId}

		netResp, err := registryClient.GetNetwork(reqGetNet)
		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, orgResp.Org.Id, netResp.Network.OrgId)
		assert.NotEqual(t, newNetName, netResp.Network.Name)

		// make sure the member is the request executor
		// in order for this test to even compile, we need to implement
		// missing endppins (PATCH /v1/networks) and missing APIs
		// requests wrappers (RegistryClient.UpdateNetwork and
		// systems/registry/api-gateway/pkg/rest/api.go for
		// api.UpdateNetworkRequest)

		reqUpdateNetwork := api.UpdateNetworkRequest{
			NetId:   oldNetId,
			NetName: newNetName}

		unResp, err := registryClient.UpdateNetwork(reqUpdateNetwork)
		assert.Error(t, err)
		assert.Nil(t, unResp)
	})
}
