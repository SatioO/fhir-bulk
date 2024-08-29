package resource

import (
	"fmt"
	"io"
	"net/http"

	"github.com/satioO/fhir/v2/domain"
)

type client struct{}

func NewFHIRResourceClient() *client {
	return &client{}
}

func (c *client) GetFHIRResource(app *domain.FHIRApp, resourceId string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, app.BaseUrl+"/bulk-export/files/"+resourceId, nil)
	req.Header.Set("Accept", "application/fhir+ndjson")
	req.Header.Set("Authorization", "Bearer "+app.Token)

	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("client: error making http request: %s", err)
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %s", err)
	}

	return resBody, nil
}
