package bulkapi

import "net/http"

type handler struct {
	bulkAPIService BulkAPIService
}

func NewBulkAPIHandler(bulkAPIService BulkAPIService) *handler {
	return &handler{bulkAPIService}
}

func (h *handler) CreateNewFHIRJob(w http.ResponseWriter, r *http.Request) {}

func (h *handler) GetFHIRJobsForApp(w http.ResponseWriter, r *http.Request) {}
