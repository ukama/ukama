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

	ory "github.com/ory/client-go"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg"
)

var SESSION_KEY = "ukama_session"

type AuthManager struct {
	client        *ory.APIClient
	serverUrl     string
	timeout       time.Duration
	ketoClientUrl string
	ketoc         *ory.APIClient
}

var object = "/v1/package"
var relation = "owner"

type UIErrorResp struct {
	Id string          `json:"id"`
	Ui ory.UiContainer `json:"ui"`
}

func NewAuthManager(serverUrl string, timeout time.Duration, ketoClientUrl string) *AuthManager {
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
	configuration = ory.NewConfiguration()
	configuration.Servers = []ory.ServerConfiguration{
		{
			URL: ketoClientUrl,
		},
	}
	kc := ory.NewAPIClient(configuration)
	ketoc := kc

	return &AuthManager{
		client:        client,
		timeout:       timeout,
		serverUrl:     serverUrl,
		ketoc:         ketoc,
		ketoClientUrl: ketoClientUrl,
	}
}

func NewAuthManagerFromClient() *AuthManager {
	return &AuthManager{
		serverUrl: "localhost",
		timeout:   1 * time.Second,
		client:    nil,
	}
}

func (am *AuthManager) ValidateSession(ss string, t string) (*ory.Session, error) {

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
	fmt.Println("KRTAOS_REQ", r.Request.URL.String())

	if err != nil {
		return nil, err
	}
	if r.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("no valid session cookie found")
	}

	return resp, nil
}

func (am *AuthManager) LoginUser(email string, password string) (*ory.SuccessfulNativeLogin, error) {
	flow, _, err := am.client.FrontendApi.CreateNativeLoginFlow(context.Background()).Execute()
	if err != nil {
		return nil, err
	}
	b := ory.UpdateLoginFlowWithPasswordMethod{
		Password:           password,
		Method:             "password",
		Identifier:         email,
		PasswordIdentifier: &email,
	}
	body := ory.UpdateLoginFlowBody{
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

type TRole struct {
	name         string
	organization string
}

func (am *AuthManager) UpdateRole(ss, t, orgId, role string, user *pkg.UserTraits) error {
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
	roles := []TRole{
		{
			name:         role,
			organization: orgId,
		},
	}
	_, r, err := am.client.IdentityApi.UpdateIdentity(
		context.Background(), user.Id,
	).UpdateIdentityBody(
		ory.UpdateIdentityBody(
			ory.UpdateIdentityBody{
				Traits: map[string]interface{}{
					"name":        user.Name,
					"email":       user.Email,
					"first_visit": user.FirstVisit,
				},
				MetadataPublic: map[string]interface{}{
					"roles": roles,
				},
			},
		),
	).Execute()

	if err != nil {
		return err
	}
	if r.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("no valid session cookie found")
	}

	return nil
}

func (am *AuthManager) AuthorizeUser(ss string, t string, role string, orgId string) (*ory.Session, error) {

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

	check, r, err := am.ketoc.PermissionApi.CheckPermission(context.Background()).
		Namespace(*&orgId).
		Object(*&object).
		Relation(*&relation).
		SubjectId(*&role).Execute()

	if err != nil {
		logrus.Errorf("Encountered error: %v\n", err)
		return nil, err
	}
	if check.Allowed {
		logrus.Infof(*&role + " can " + *&relation + " the " + *&object)
		return resp, nil
	}
	return nil, fmt.Errorf(role + " is not authorized to " + relation + " the " + object)
}
