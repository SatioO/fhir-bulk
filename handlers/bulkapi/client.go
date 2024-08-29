package bulkapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type client struct{}

func NewBulkFHIRClient() *client {
	return &client{}
}

func (c *client) CreateNewJob(baseUrl string, request *TriggerFHIRJobRequest) (string, error) {
	var URL string

	if request.GroupID != "" {
		URL = baseUrl + "/Group/" + request.GroupID + "/$export"
	} else {
		URL = baseUrl + "/Patient/$export"
	}

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	req.Header.Set("Accept", "application/fhir+json")
	req.Header.Set("Prefer", "respond-async")

	if err != nil {
		return "", fmt.Errorf("client: could not create request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("client: error making http request: %s", err)
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
