package fhirapp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/satioO/fhir/v2/api"
	"github.com/satioO/fhir/v2/domain"
	"github.com/satioO/fhir/v2/handlers/auth"
)

type handler struct {
	fhirAppService FHIRAppService
	authService    auth.AuthService
}

func NewFHIRAppHandler(fhirAppService FHIRAppService, authService auth.AuthService) *handler {
	return &handler{fhirAppService, authService}
}

func (h *handler) GetApps(w http.ResponseWriter, r *http.Request) {
	result, err := h.fhirAppService.GetApps()
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to get apps: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, result)
}

func (h *handler) GetAppById(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("appId")

	result, err := h.fhirAppService.GetAppById(appId)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to fetch app: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, result)
}

func (h *handler) RegisterApp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	jsonObj, err := io.ReadAll(r.Body)
	if err != nil {
		api.Error(w, r, domain.ErrReadingRequestBody, http.StatusBadRequest)
		return
	}

	var body CreateFHIRAppRequest
	if err := json.Unmarshal(jsonObj, &body); err != nil {
		api.Error(w, r, domain.ErrParsingRequestBody, http.StatusBadRequest)
		return
	}

	app, err := h.fhirAppService.CreateApp(&body)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to create app: %v", err), http.StatusInternalServerError)
		return
	}

	token, err := h.authService.RegisterAuthServer(app.ID, &body.RegisterAuthServerRequest)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to create auth server: %v", err), http.StatusInternalServerError)
		return
	}
	app.Token = token

	if err := h.fhirAppService.UpdateToken(app.ID, token); err != nil {
		api.Error(w, r, fmt.Errorf("failed to update token: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, app)
}
