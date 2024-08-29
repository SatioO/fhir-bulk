package bulkapi

import "time"

type TriggerFHIRJobRequest struct {
	GroupID string `json:"groupId"`
}

type TriggerFHIRJobResponse struct {
	AppID     string                 `json:"appId"`
	JobID     string                 `json:"jobId"`
	Status    string                 `json:"status"`
	Resources []FHIRResourceResponse `json:"resources"`
}

type FHIRResourceResponse struct {
	Type       string `json:"type"`
	ResourceID string `json:"resourceId"`
}

type FHIRJobStatusClientResponse struct {
	TransactionTime     time.Time                    `json:"transactionTime"`
	Request             string                       `json:"request"`
	RequiresAccessToken bool                         `json:"requiresAccessToken"`
	Output              []FHIRResourceClientResponse `json:"output"`
}

type FHIRResourceClientResponse struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}
