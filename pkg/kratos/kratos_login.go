package kratos

import (
	"context"
	"net/http"

	ory "github.com/ory/kratos-client-go"
	"github.com/paralus/cli/pkg/log"
)

type kratosLogin struct {
	ory.SuccessfulSelfServiceLoginWithoutBrowser
}

type KratosLoginClient interface {
	HttpGet(url string) (*http.Response, error)
}

func Login(clientURL, email, password string) (KratosLoginClient, error) {
	ctx := context.Background()
	client := getNewAPIClient(clientURL)

	// Initialize a registration flow
	log.GetLogger().Debug("Initializing the login flow.")
	flow, _, err := client.V0alpha2Api.InitializeSelfServiceLoginFlowWithoutBrowser(ctx).Execute()
	if err != nil {
		return nil, err
	}
	log.GetLogger().Debugf("Flow id fetched successfully issued_at: %v, expires_at: %v", flow.IssuedAt, flow.ExpiresAt)
	log.GetLogger().Debug("Logging in using user credentials.")
	result, _, err := client.V0alpha2Api.SubmitSelfServiceLoginFlow(ctx).Flow(flow.Id).SubmitSelfServiceLoginFlowBody(
		ory.SubmitSelfServiceLoginFlowWithPasswordMethodBodyAsSubmitSelfServiceLoginFlowBody(&ory.SubmitSelfServiceLoginFlowWithPasswordMethodBody{
			Method:             "password",
			Password:           password,
			Identifier:         email,
			PasswordIdentifier: &email,
		}),
	).Execute()
	if err != nil {
		return nil, err
	}

	log.GetLogger().Debug("User credentials validated successfully.")
	log.GetLogger().Debugf("User logged in successfully, session_id: %v, issued_at: %v, expires_at: %v, identity_id: %v, state: %v", result.Session.Id, result.Session.IssuedAt, result.Session.Identity.Id, result.Session.Identity.State)

	return kratosLogin{*result}, nil
}

func (k kratosLogin) HttpGet(url string) (*http.Response, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Session-Token", *k.SessionToken)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, err
}
