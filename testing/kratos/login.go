package kratos

import (
	"context"

	ory "github.com/ory/kratos-client-go"
)

func Login(kratosUrl string) (*ory.SuccessfulSelfServiceLoginWithoutBrowser, error) {
	var client = NewSDKForSelfHosted(kratosUrl)

	ctx := context.Background()

	// Create a temporary user
	email, password := RandomCredentials()
	_, _, err := CreateIdentityWithSession(client, email, password)
	if err != nil {
		return nil, err
	}
	// Initialize the flow
	flow, res, err := client.V0alpha2Api.InitializeSelfServiceLoginFlowWithoutBrowser(ctx).Execute()
	LogKratosSdkError(err, res)
	if err != nil {
		return nil, err
	}

	// If you want, print the flow here:
	//
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
