package bulkapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/satioO/fhir/v2/domain"
)

type client struct{}

func NewBulkFHIRClient() *client {
	return &client{}
}

func (c *client) CreateNewJob(app *domain.FHIRApp, request *TriggerFHIRJobRequest) (string, error) {
	var URL string

	if request.GroupID != "" {
		URL = app.BaseUrl + "/Group/" + request.GroupID + "/$export"
	} else {
		URL = app.BaseUrl + "/Patient/$export"
	}

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	req.Header.Set("Accept", "application/fhir+json")
	req.Header.Set("Prefer", "respond-async")
	req.Header.Set("Authorization", "Bearer "+app.Token)

	if err != nil {
		return "", fmt.Errorf("client: could not create request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("client: error making http request: %s", err)
	}

	if res.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("client: request failed with unknown error: %v", res)
	}

	contentUrl := res.Header.Get("Content-Location")

	jobId, err := c.parseJobID(contentUrl)
	if err != nil {
		return "", fmt.Errorf("client: error parsing content url: %s", err)
	}

	return jobId, nil
}

func (c *client) DeleteFHIRJob(jobId string) error {
	return nil
}

func (c *client) parseJobID(contentUrl string) (string, error) {
	parsedUrl, err := url.Parse(contentUrl)
	if err != nil {
		return "", err
	}

	segments := strings.Split(parsedUrl.Path, "/")
	return segments[len(segments)-1], nil
}
