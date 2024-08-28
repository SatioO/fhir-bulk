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
	fhirAppService FHIRAppService
	auth           auth.AuthService
}

func NewFHIRAppHandler(fhirAppService FHIRAppService, auth auth.AuthService) *fhirApp {
	return &fhirApp{fhirAppService, auth}
}

func (p *fhirApp) GetApps(w http.ResponseWriter, r *http.Request) {
	result, err := p.fhirAppService.GetApps()

	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to get apps: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, result)
}

func (p *fhirApp) GetAppById(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("appId")

	result, err := p.fhirAppService.GetAppById(appId)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to fetch app: %v", err), http.StatusBadRequest)
		return
	}

	api.SuccessJson(w, r, result)
}

func (p *fhirApp) RegisterApp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

	app, err := p.fhirAppService.CreateApp(&body)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to create app: %v", err), http.StatusInternalServerError)
		return
	}

	token, err := p.auth.RegisterAuthServer(app.ID, &body.RegisterAuthServerRequest)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to create auth server: %v", err), http.StatusInternalServerError)
		return
	}
	app.Token = token

	if err := p.fhirAppService.UpdateToken(app.ID, token); err != nil {
		api.Error(w, r, fmt.Errorf("failed to update token: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, app)
}
