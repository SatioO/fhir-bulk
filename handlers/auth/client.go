package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/satioO/fhir/v2/domain"
)

type client struct{}

func NewAuthClient() *client {
	return &client{}
}

func (c *client) GenerateToken(request *domain.FHIRAuthServer) (string, error) {
	body := `grant_type=client_credentials&scope=` + request.Scopes
	req, err := http.NewRequest(http.MethodPost, request.TokenURL, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(request.ClientID+":"+request.ClientSecret)))

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

	return response.AccessToken, nil
}
