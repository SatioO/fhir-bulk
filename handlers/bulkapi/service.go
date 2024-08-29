package bulkapi

import (
	"fmt"

	"github.com/satioO/fhir/v2/domain"
	"github.com/satioO/fhir/v2/handlers/fhirapp"
)

type BulkAPIService interface {
	GetJobsByApp(appId string) ([]domain.FHIRJob, error)
	GetFHIRJobStatus(appId, jobId string) (TriggerFHIRJobResponse, error)
	CreateNewFHIRJob(appId string, body *TriggerFHIRJobRequest) (domain.FHIRJob, error)
	DeleteFHIRJob(appId, jobId string) error
}

type service struct {
	fhirJobRepo    *FHIRJobRepo
	fhirAppRepo    *fhirapp.FHIRAppRepo
	bulkFHIRClient *client
}

func NewBulkAPIService(fhirJobRepo *FHIRJobRepo, fhirAppRepo *fhirapp.FHIRAppRepo, bulkFHIRClient *client) BulkAPIService {
	return &service{fhirJobRepo, fhirAppRepo, bulkFHIRClient}
}

func (s *service) GetJobsByApp(appId string) ([]domain.FHIRJob, error) {
	return s.fhirJobRepo.GetJobsByApp(appId)
}

func (s *service) GetFHIRJobStatus(appId, jobId string) (TriggerFHIRJobResponse, error) {
	app, err := s.fhirAppRepo.GetAppById(appId)

	if err != nil {
		return TriggerFHIRJobResponse{}, err
	}

	return s.bulkFHIRClient.GetFHIRJobStatus(&app, jobId)
}

func (s *service) CreateNewFHIRJob(appId string, body *TriggerFHIRJobRequest) (domain.FHIRJob, error) {
	app, err := s.fhirAppRepo.GetAppById(appId)

	if err != nil {
		return domain.FHIRJob{}, err
	}

	jobId, err := s.bulkFHIRClient.CreateNewJob(&app, body)
	if err != nil {
		return domain.FHIRJob{}, err
	}

	job, err := s.fhirJobRepo.CreateJob(&domain.FHIRJob{ID: jobId, AppID: appId, Status: "submitted"})
	if err != nil {
		return domain.FHIRJob{}, fmt.Errorf("error creating fhir job: %s", err)
	}

	return job, nil
}

func (s *service) DeleteFHIRJob(appId, jobId string) error {
	_, err := s.fhirAppRepo.GetAppById(appId)

	if err != nil {
		return err
	}

	s.bulkFHIRClient.DeleteFHIRJob(jobId)
	return nil
}
