package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ory/client-go"
	ory "github.com/ory/client-go"
)

var SESSION_KEY = "ukama_session"

func ValidateSession(ss string, t string, o *ory.APIClient) (*client.Session, error) {
	if t == "cookie" {
		urlObj, _ := url.Parse(o.GetConfig().Servers[0].URL)
		cookie := &http.Cookie{
			Name:  SESSION_KEY,
			Value: ss,
		}
		o.GetConfig().HTTPClient.Jar.SetCookies(urlObj, []*http.Cookie{cookie})
	} else if t == "header" {
		o.GetConfig().AddDefaultHeader("X-Session-Token", ss)
	}
	resp, r, err := o.FrontendApi.ToSession(context.Background()).Execute()
	if err != nil {
		return nil, err
	}
	if r.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("no valid session cookie found")
	}
	return resp, nil
}

func LoginUser(email string, password string, o *ory.APIClient) (*client.SuccessfulNativeLogin, error) {
	flow, _, err := o.FrontendApi.CreateNativeLoginFlow(context.Background()).Execute()
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
	flow1, _, err := o.FrontendApi.UpdateLoginFlow(context.Background()).Flow(flow.Id).UpdateLoginFlowBody(body).Execute()
	if err != nil {
		return nil, err
	}
	return flow1, nil
}
