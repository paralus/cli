package kratos

import (
	"net/http"
	"net/http/cookiejar"

	ory "github.com/ory/kratos-client-go"
)

func getNewAPIClient(endpoint string) *ory.APIClient {
	conf := ory.NewConfiguration()
	conf.Servers = ory.ServerConfigurations{{URL: endpoint}}
	cj, _ := cookiejar.New(nil)
	conf.HTTPClient = &http.Client{Jar: cj}
	return ory.NewAPIClient(conf)
}
