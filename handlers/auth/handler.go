package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/satioO/fhir/v2/api"
)

type authHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *authHandler {
	return &authHandler{authService}
}

func (h *authHandler) GetAuthServerForApp(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("appId")
	result, err := h.authService.GetAuthServerForApp(appId)

	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to fetch client: %v", err), http.StatusBadRequest)
		return
	}

	api.SuccessJson(w, r, result)
}

func (h *authHandler) RegisterAuthServer(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("appId")
	jsonObj, err := io.ReadAll(r.Body)

	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to read body: %v", err), http.StatusBadRequest)
		return
	}

	var body RegisterAuthServerRequest
	if err := json.Unmarshal(jsonObj, &body); err != nil {
		api.Error(w, r, fmt.Errorf("failed to parse body: %v", err), http.StatusBadRequest)
		return
	}

	result, err := h.authService.RegisterAuthServer(appId, &body)

	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to fetch client: %v", err), http.StatusBadRequest)
		return
	}

	api.SuccessJson(w, r, result)

}
