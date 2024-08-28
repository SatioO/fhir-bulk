package fhir_app

import "github.com/satioO/fhir/v2/handlers/auth"

type CreateFHIRAppRequest struct {
	ID                             string `json:"id"`
	BaseUrl                        string `json:"baseUrl"`
	auth.RegisterAuthServerRequest `json:"auth"`
}
