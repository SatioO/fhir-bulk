package bulkapi

import (
	"encoding/json"
	"fmt"
	"io"
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
		return "", fmt.Errorf("client: request failed with error: %v", res)
	}

	contentUrl := res.Header.Get("Content-Location")

	jobId, err := c.parseId(contentUrl)
	if err != nil {
		return "", fmt.Errorf("client: error parsing content url: %s", err)
	}

	return jobId, nil
}

func (c *client) GetFHIRJobStatus(app *domain.FHIRApp, jobId string) (TriggerFHIRJobResponse, error) {
	req, err := http.NewRequest(http.MethodGet, app.BaseUrl+"/bulk-export/jobs/"+jobId, nil)
	req.Header.Set("Accept", "application/fhir+json")
	req.Header.Set("Authorization", "Bearer "+app.Token)

	if err != nil {
		return TriggerFHIRJobResponse{}, fmt.Errorf("client: could not create request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return TriggerFHIRJobResponse{}, fmt.Errorf("client: error making http request: %s", err)
	}

	if res.StatusCode == http.StatusOK {
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return TriggerFHIRJobResponse{}, fmt.Errorf("client: could not read response body: %s", err)
		}

		var response FHIRJobStatusClientResponse
		if err := json.Unmarshal(resBody, &response); err != nil {
			return TriggerFHIRJobResponse{}, fmt.Errorf("client: could not parse response body: %s", err)
		}

		resources := []FHIRResourceResponse{}
		for _, resource := range response.Output {
			resourceId, err := c.parseId(resource.URL)
			if err != nil {
				return TriggerFHIRJobResponse{}, fmt.Errorf("client: error parsing resource url: %s", err)
			}

			resources = append(resources, FHIRResourceResponse{Type: resource.Type, ResourceID: resourceId})
		}

		return TriggerFHIRJobResponse{AppID: app.ID, JobID: jobId, Status: "completed", Resources: resources}, nil
	} else if res.StatusCode == http.StatusAccepted {
		resources := []FHIRResourceResponse{}
		return TriggerFHIRJobResponse{AppID: app.ID, JobID: jobId, Status: "submitted", Resources: resources}, nil
	} else {
		return TriggerFHIRJobResponse{}, fmt.Errorf("client: request failed with error: %v", res)
	}
}

func (c *client) DeleteFHIRJob(jobId string) error {
	return nil
}

func (c *client) parseId(contentUrl string) (string, error) {
	parsedUrl, err := url.Parse(contentUrl)
	if err != nil {
		return "", err
	}

	segments := strings.Split(parsedUrl.Path, "/")
	return segments[len(segments)-1], nil
}
