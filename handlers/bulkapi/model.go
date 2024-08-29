package bulkapi

type TriggerFHIRJobRequest struct {
	GroupID string `json:"groupId"`
}

type TriggerFHIRJobResponse struct {
	Status string
	JobID  string
}
