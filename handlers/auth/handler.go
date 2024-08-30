package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/satioO/fhir/v2/api"
	"github.com/satioO/fhir/v2/domain"
)

type handler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *handler {
	return &handler{authService}
}

func (h *handler) GetAuthServerForApp(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("appId")

	result, err := h.authService.GetAuthServerForApp(appId)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to fetch auth server details: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, result)
}

func (h *handler) RegisterAuthServer(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("appId")

	jsonObj, err := io.ReadAll(r.Body)
	if err != nil {
		api.Error(w, r, domain.ErrReadingRequestBody, http.StatusBadRequest)
		return
	}

	var body RegisterAuthServerRequest
	if err := json.Unmarshal(jsonObj, &body); err != nil {
		api.Error(w, r, domain.ErrParsingRequestBody, http.StatusBadRequest)
		return
	}

	result, err := h.authService.RegisterAuthServer(appId, &body)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to register auth server: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, result)

}
