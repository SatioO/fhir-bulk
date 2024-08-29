package bulkapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/satioO/fhir/v2/api"
)

type handler struct {
	bulkAPIService BulkAPIService
}

func NewBulkAPIHandler(bulkAPIService BulkAPIService) *handler {
	return &handler{bulkAPIService}
}

func (h *handler) GetFHIRJobsForApp(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("appId")
	result, err := h.bulkAPIService.GetJobsByApp(appId)

	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to fetch jobs: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, result)
}

func (h *handler) CreateNewFHIRJob(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("appId")

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to read request body: %s", err), http.StatusBadRequest)
		return
	}

	var request TriggerFHIRJobRequest
	if err := json.Unmarshal(body, &request); err != nil {
		api.Error(w, r, fmt.Errorf("failed to parse request body:%s", err), http.StatusInternalServerError)
		return
	}

	result, err := h.bulkAPIService.CreateNewFHIRJob(appId, &request)
	if err != nil {
		api.Error(w, r, fmt.Errorf("failed to register job: %v", err), http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, result)
}

func (h *handler) DeleteFHIRJob(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("appId")
	jobId := r.PathValue("jobId")

	h.bulkAPIService.DeleteFHIRJob(appId, jobId)

	api.SuccessJson(w, r, nil)
}
