package lookup

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	sr "github.com/ukama/ukama/services/common/srvcrouter"
)

func Test_LookupRequestOrgCredentialForNodePass(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	reply := &RespOrgCredentials{
		Node:    "abcd",
		OrgName: "org",
		Ip:      "127.0.0.1",
		OrgCred: []byte("aGVscG1lCg=="),
	}

	// arrange
	ServiceRouter := "http://localhost:8091"

	fakeUrl := ServiceRouter + "/service"

	// fetch the article into struct
	rs := sr.NewServiceRouter(ServiceRouter)

	l := NewLookUp(rs)

	httpmock.Activate()
	httpmock.ActivateNonDefault(l.S.C.GetClient())
	httpmock.RegisterResponder("GET", fakeUrl, httpmock.NewJsonResponderOrPanic(200, &reply))

	val, cred, err := l.LookupRequestOrgCredentialForNode("abcd")

	// do stuff with the article object ...
	assert.Nil(t, err)

	assert.Equal(t, val, true)

	assert.NotNil(t, cred)

	assert.Equal(t, cred.Node, reply.Node)
	assert.Equal(t, cred.Ip, reply.Ip)
	assert.Equal(t, cred.OrgName, reply.OrgName)
	assert.Equal(t, cred.OrgCred, reply.OrgCred)
}

func Test_LookupRequestOrgCredentialForNodePassFailure(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	reply := &RespOrgCredentials{
		Node:    "abcd",
		OrgName: "org",
		Ip:      "127.0.0.1",
		OrgCred: []byte("aGVscG1lCg=="),
	}

	// arrange
	ServiceRouter := "http://localhost:8091"

	fakeUrl := ServiceRouter + "/service"

	// fetch the article into struct
	rs := sr.NewServiceRouter(ServiceRouter)

	l := NewLookUp(rs)

	httpmock.Activate()
	httpmock.ActivateNonDefault(l.S.C.GetClient())
	httpmock.RegisterResponder("GET", fakeUrl, httpmock.NewJsonResponderOrPanic(400, &reply))

	val, cred, err := l.LookupRequestOrgCredentialForNode("abcd")

	assert.NotNil(t, err)

	assert.Equal(t, val, false)

	assert.Contains(t, err.Error(), "failed to get credentials:")

	assert.Nil(t, cred)
}
