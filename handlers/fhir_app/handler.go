package fhir_app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/satioO/fhir/v2/api"
	"github.com/satioO/fhir/v2/handlers/auth"
)

type fhirApp struct {
	fhirAppRepo *FHIRAppRepo
	auth        auth.AuthService
}

func NewFHIRAppHandler(fhirAppRepo *FHIRAppRepo, auth auth.AuthService) *fhirApp {
	return &fhirApp{fhirAppRepo: fhirAppRepo, auth: auth}
}

func (p *fhirApp) GetApps(w http.ResponseWriter, r *http.Request) {
	result, err := p.fhirAppRepo.GetApps()

	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to get apps: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, result)
}

func (p *fhirApp) GetAppById(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("appId")

	result, err := p.fhirAppRepo.GetAppById(appId)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to fetch app: %v", err), http.StatusBadRequest)
		return
	}

	api.SuccessJson(w, r, result)
}

func (p *fhirApp) RegisterApp(w http.ResponseWriter, r *http.Request) {
	jsonObj, err := io.ReadAll(r.Body)

	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to read body: %v", err), http.StatusBadRequest)
		return
	}

	var body CreateFHIRAppRequest
	if err := json.Unmarshal(jsonObj, &body); err != nil {
		api.Error(w, r, fmt.Errorf("failed to parse body: %v", err), http.StatusBadRequest)
		return
	}

	app, err := p.fhirAppRepo.CreateApp(&body)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to create app: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = p.auth.RegisterAuthServer(app.ID, body.RegisterAuthServerRequest)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to create auth server: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, app)
}
