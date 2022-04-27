package kratos

// those methods are copy-paste from kratos example repo https://github.com/ory/kratos/tree/master/examples/go/pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/google/uuid"
	ory "github.com/ory/kratos-client-go"
	"github.com/sirupsen/logrus"
)

func PrintJSONPretty(v interface{}) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}

func MustReadAll(r io.Reader) []byte {
	all, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return all
}

func NewSDKForSelfHosted(endpoint string) *ory.APIClient {
	conf := ory.NewConfiguration()
	conf.Servers = ory.ServerConfigurations{{URL: endpoint}}
	cj, _ := cookiejar.New(nil)
	conf.HTTPClient = &http.Client{Jar: cj}
	return ory.NewAPIClient(conf)
}

func LogKratosSdkError(err error, res *http.Response) {
	if err == nil {
		return
	}
	body, _ := json.MarshalIndent(json.RawMessage(MustReadAll(res.Body)), "", "  ")
	out, _ := json.MarshalIndent(err, "", "  ")
	logrus.Printf("%s\n\nAn error occurred: %+v\nbody: %s\n", out, err, body)
}

func RandomCredentials() (email, password string) {
	email = "dev+" + uuid.New().String() + "@dev.ukama.com"
	password = strings.ReplaceAll(uuid.New().String(), "-", "")
	return
}

// CreateIdentityWithSession creates an identity and an Ory Session Token for it.
func CreateIdentityWithSession(c *ory.APIClient, email, password string) (*ory.Session, string, error) {
	ctx := context.Background()

	if email == "" && password == "" {
		panic("empty username or password")
	}

	// Initialize a registration flow
	flow, _, err := c.V0alpha2Api.InitializeSelfServiceRegistrationFlowWithoutBrowser(ctx).Execute()
	if err != nil {
		return nil, "", err
	}

	// Submit the registration flow
	result, res, err := c.V0alpha2Api.SubmitSelfServiceRegistrationFlow(ctx).Flow(flow.Id).SubmitSelfServiceRegistrationFlowBody(
		ory.SubmitSelfServiceRegistrationFlowWithPasswordMethodBodyAsSubmitSelfServiceRegistrationFlowBody(&ory.SubmitSelfServiceRegistrationFlowWithPasswordMethodBody{
			Method:   "password",
			Password: password,
			Traits:   map[string]interface{}{"email": email},
		}),
	).Execute()

	LogKratosSdkError(err, res)

	if err != nil {
		return nil, "", err
	}

	if result.Session == nil {
		log.Fatalf("The server is expected to create sessions for new registrations.")
	}

	return result.Session, *result.SessionToken, nil
}
