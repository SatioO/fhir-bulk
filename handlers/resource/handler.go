package resource

import (
	"fmt"
	"net/http"

	"github.com/satioO/fhir/v2/api"
)

type handler struct {
	resourceService FHIRResourceService
}

func NewFHIRResourceHandler(resourceService FHIRResourceService) *handler {
	return &handler{resourceService}
}

func (h *handler) GetFHIRResourcesByJobID(w http.ResponseWriter, r *http.Request) {
	jobId := r.PathValue("jobId")

	result, err := h.resourceService.GetFHIRResourcesByJobID(jobId)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to get resources: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, result)
}

func (h *handler) GetFHIRResource(w http.ResponseWriter, r *http.Request) {
	jobId := r.PathValue("jobId")
	resourceId := r.PathValue("resourceId")

	result, err := h.resourceService.GetFHIRResource(jobId, resourceId)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to get resource details: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/fhir+ndjson")
	w.WriteHeader(http.StatusOK)

	// Write the response body
	_, err = w.Write(result)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
