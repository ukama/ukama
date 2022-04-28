package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/services/bootstrap/bootstrap/pkg"
	"github.com/ukama/ukama/services/bootstrap/bootstrap/pkg/nmr"
	"github.com/ukama/ukama/services/common/rest"

	sr "github.com/ukama/ukama/services/common/srvcrouter"
)

func init() {
	pkg.IsDebugMode = true
}

var defaultCongif = &pkg.Config{
	Server: rest.HttpConfig{
		Cors: cors.Config{
			AllowAllOrigins: true,
		},
	},
}

func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	rs := sr.ServiceRouter{}

	r := NewRouter(defaultCongif, &rs).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_GetbootstrapGetHandler(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	nmrReply := &nmr.NodeStatus{
		Status: "StatusNodeIntransit",
	}

	// lookupReply := &lookup.RespOrgCredentials{
	// 	Node:    "abcd",
	// 	OrgName: "org",
	// 	Ip:      "127.0.0.1",
	// 	OrgCred: []byte("aGVscG1lCg=="),
	// }

	// arrange
	ServiceRouter := "http://localhost:8091"

	//lsFakeUrl := ServiceRouter + "/service"
	fsFakeUrl := ServiceRouter + "/service"

	rs := sr.NewServiceRouter(ServiceRouter)
	r := NewRouter(defaultCongif, rs)
	re := r.fizz.Engine()
	httpmock.Activate()
	httpmock.ActivateNonDefault(rs.C.GetClient())

	//httpmock.RegisterResponder("GET", lsFakeUrl, httpmock.NewJsonResponderOrPanic(200, lookupReply))
	httpmock.RegisterResponder("GET", fsFakeUrl, httpmock.NewJsonResponderOrPanic(200, nmrReply))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?node=abcd&looking_for=validation", nil)

	// act
	re.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "")

}
