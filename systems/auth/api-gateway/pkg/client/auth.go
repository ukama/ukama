package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/ory/client-go"
	ory "github.com/ory/client-go"
	"github.com/ukama/ukama/systems/common/uuid"
)

var SESSION_KEY = "ukama_session"
var namespace ="9982-23984-349389"
var object ="org"
var relation ="member"
var subjectId ="user"
type AuthManager struct {
	client    *ory.APIClient
	serverUrl string
	timeout   time.Duration
	orgRegistry  OrgMemberRoleClient
	ketoClient *client.APIClient

}

type UIErrorResp struct {
	Id string          `json:"id"`
	Ui ory.UiContainer `json:"ui"`
}

func NewAuthManager(serverUrl string, timeout time.Duration, orgRegistry OrgMemberRoleClient) *AuthManager {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	configuration := ory.NewConfiguration()
	configuration.Servers = []ory.ServerConfiguration{
		{
			URL: serverUrl,
		},
	}
	

	jar, _ := cookiejar.New(nil)
	oc := ory.NewAPIClient(configuration)
	oc.GetConfig().HTTPClient = &http.Client{
		Jar: jar,
	}
	client := oc

	return &AuthManager{
		client:    client,
		timeout:   timeout,
		serverUrl: serverUrl,
		orgRegistry:   orgRegistry,
	}
}

func NewAuthManagerFromClient() *AuthManager {
	return &AuthManager{
		serverUrl: "localhost",
		timeout:   1 * time.Second,
		client:    nil,
		orgRegistry:  nil,
	}
}

func (am *AuthManager) ValidateSession(ss string, t string) (*client.Session, error) {

	if t == "cookie" {
		urlObj, _ := url.Parse(am.client.GetConfig().Servers[0].URL)
		cookie := &http.Cookie{
			Name:  SESSION_KEY,
			Value: ss,
		}
		am.client.GetConfig().HTTPClient.Jar.SetCookies(urlObj, []*http.Cookie{cookie})
	} else if t == "header" {
		am.client.GetConfig().AddDefaultHeader("X-Session-Token", ss)

	}
	resp, r, err := am.client.FrontendApi.ToSession(context.Background()).Execute()
	if err != nil {
		return nil, err
	}
	if r.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("no valid session cookie found")
	}
	
	userId, _ := uuid.FromString("28669952-d8ba-47d1-b8a7-59371a667f9b")
	orgId, _ := uuid.FromString("39a2138f-d6e2-42fb-ae92-9056ad125df4")

	role,err := am.orgRegistry.GetMemberRole(userId,orgId )
	if err != nil {
		return nil, fmt.Errorf("error while getting member role: %s", err.Error())
	}
	fmt.Println("Role", role)
	
	return resp, nil
}

func (am *AuthManager) LoginUser(email string, password string) (*client.SuccessfulNativeLogin, error) {
	flow, _, err := am.client.FrontendApi.CreateNativeLoginFlow(context.Background()).Execute()
	if err != nil {
		return nil, err
	}
	b := client.UpdateLoginFlowWithPasswordMethod{
		Password:           password,
		Method:             "password",
		Identifier:         email,
		PasswordIdentifier: &email,
	}
	body := client.UpdateLoginFlowBody{
		UpdateLoginFlowWithPasswordMethod: &b,
	}
	flow1, resp, err := am.client.FrontendApi.UpdateLoginFlow(context.Background()).Flow(flow.Id).UpdateLoginFlowBody(body).Execute()

	if err != nil {
		if resp.StatusCode == http.StatusBadRequest {
			u := UIErrorResp{}
			buf := &bytes.Buffer{}
			_, e := buf.ReadFrom(resp.Body)
			if e != nil {
				return nil, e
			}
			e = json.Unmarshal(buf.Bytes(), &u)
			if e != nil {
				return nil, e
			}
			return nil, fmt.Errorf("%v", u.Ui.Messages[0].Text)
		}
		return nil, err
	}

	return flow1, nil
}
