package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/satioO/fhir/v2/domain"
)

type authClient struct{}

func NewAuthClient() *authClient {
	return &authClient{}
}

func (a *authClient) GenerateToken(request *domain.FhirAuthServer) (string, error) {
	body := `grant_type=client_credentials&scope=system/Observation.read system/Patient.read system/Encounter.read`
	req, err := http.NewRequest(http.MethodPost, request.TokenURL, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic M2FiNTdiNDQtOGNjNi00ZWQ3LTg0MjktNDY5OTk5Mzc0Zjg0OlU4NU90bHVlUXB6MEstVmoycHNDYlJ4X2ZYd3JVc3h4")

	if err != nil {
		return "", fmt.Errorf("client: could not create request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("client: error making http request: %s", err)
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("client: could not read response body: %s", err)
	}

	var response AuthClientResponse
	if err := json.Unmarshal(resBody, &response); err != nil {
		return "", fmt.Errorf("client: could not parse response body: %s", err)
	}

	log.Println(response)
	return response.AccessToken, nil
}
