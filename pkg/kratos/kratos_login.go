package kratos

import (
	"context"
	"net/http"

	ory "github.com/ory/kratos-client-go"
	"github.com/paralus/cli/pkg/log"
)

type kratosLogin struct {
	ory.SuccessfulNativeLogin
}

type KratosLoginClient interface {
	HttpGet(url string) (*http.Response, error)
}

func Login(clientURL, email, password string) (KratosLoginClient, error) {
	ctx := context.Background()
	client := getNewAPIClient(clientURL)

	log.GetLogger().Debug("Initializing the login flow.")
	flow, _, err := client.FrontendApi.CreateNativeLoginFlow(ctx).Execute()
	if err != nil {
		return nil, err
	}
	log.GetLogger().Debugf("Flow id fetched successfully issued_at: %v, expires_at: %v", flow.GetIssuedAt(), flow.GetExpiresAt())

	log.GetLogger().Debug("Logging in using user credentials.")
	body := ory.UpdateLoginFlowBody{
		UpdateLoginFlowWithPasswordMethod: &ory.UpdateLoginFlowWithPasswordMethod{
			Method:             "password",
			Password:           password,
			Identifier:         email,
			PasswordIdentifier: &email,
		},
	}
	login, _, err := client.FrontendApi.UpdateLoginFlow(ctx).UpdateLoginFlowBody(body).Flow(flow.GetId()).Execute()
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"session_id":  login.Session.GetId(),
		"issued_at":   login.Session.GetIssuedAt(),
		"expires_at":  login.Session.GetExpiresAt(),
		"identity_id": login.Session.Identity.GetId(),
		"user_state":  login.Session.GetActive(),
	}
	log.GetLogger().Debugf("User logged in successfully. User info: %v", info)
	return kratosLogin{*login}, nil
}

func (k kratosLogin) HttpGet(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Session-Token", k.GetSessionToken())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, err
}
