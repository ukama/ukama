package kratos

import (
	"context"

	ory "github.com/ory/kratos-client-go"
)

func Login(kratosUrl string, email string, password string) (*ory.SuccessfulSelfServiceLoginWithoutBrowser, error) {
	var client = NewSDKForSelfHosted(kratosUrl)
	ctx := context.Background()

	// Initialize the flow
	flow, res, err := client.V0alpha2Api.InitializeSelfServiceLoginFlowWithoutBrowser(ctx).Execute()
	LogKratosSdkError(err, res)
	if err != nil {
		return nil, err
	}

	// If you want, print the flow here:
	PrintJSONPretty(flow)

	// Submit the form
	result, res, err := client.V0alpha2Api.SubmitSelfServiceLoginFlow(ctx).Flow(flow.Id).SubmitSelfServiceLoginFlowBody(
		ory.SubmitSelfServiceLoginFlowWithPasswordMethodBodyAsSubmitSelfServiceLoginFlowBody(&ory.SubmitSelfServiceLoginFlowWithPasswordMethodBody{
			Method:             "password",
			Password:           password,
			PasswordIdentifier: email,
		}),
	).Execute()
	LogKratosSdkError(err, res)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CreateTemporaryUser(client *ory.APIClient) (email string, password string, err error) {
	// Create a temporary user
	email, password = RandomCredentials()
	err = CreateIdentityWithSession(client, email, password)
	if err != nil {
		return "", "", err
	}
	return email, password, nil
}
