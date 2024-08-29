package bulkapi

import (
	"fmt"

	"github.com/satioO/fhir/v2/domain"
	"github.com/satioO/fhir/v2/handlers/fhirapp"
)

type BulkAPIService interface {
	GetJobsByApp(appId string) ([]domain.FHIRJob, error)
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

func (s *service) CreateNewFHIRJob(appId string, body *TriggerFHIRJobRequest) (domain.FHIRJob, error) {
	app, err := s.fhirAppRepo.GetAppById(appId)

	if err != nil {
		return domain.FHIRJob{}, err
	}

	jobId, err := s.bulkFHIRClient.CreateNewJob(app.BaseUrl, body)
	if err != nil {
		return domain.FHIRJob{}, err
	}

	job, err := s.fhirJobRepo.CreateJob(&domain.FHIRJob{JobID: jobId, AppID: appId, Status: "inprogress"})
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
