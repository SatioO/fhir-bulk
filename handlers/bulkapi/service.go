package bulkapi

import (
	"fmt"

	"github.com/satioO/fhir/v2/domain"
	"github.com/satioO/fhir/v2/repositories"
)

type BulkAPIService interface {
	GetJobsByApp(appId string) ([]domain.FHIRJob, error)
	GetFHIRJobStatus(appId, jobId string) (TriggerFHIRJobResponse, error)
	CreateNewFHIRJob(appId string, body *TriggerFHIRJobRequest) (domain.FHIRJob, error)
	DeleteFHIRJob(appId, jobId string) error
}

type service struct {
	fhirJobRepo      *repositories.FHIRJobRepo
	fhirAppRepo      *repositories.FHIRAppRepo
	fhirResourceRepo *repositories.FHIRResourceRepo
	bulkFHIRClient   *client
}

func NewBulkAPIService(fhirJobRepo *repositories.FHIRJobRepo, fhirAppRepo *repositories.FHIRAppRepo, fhirResourceRepo *repositories.FHIRResourceRepo, bulkFHIRClient *client) BulkAPIService {
	return &service{fhirJobRepo, fhirAppRepo, fhirResourceRepo, bulkFHIRClient}
}

func (s *service) GetJobsByApp(appId string) ([]domain.FHIRJob, error) {
	return s.fhirJobRepo.GetJobsByApp(appId)
}

func (s *service) GetFHIRJobStatus(appId, jobId string) (TriggerFHIRJobResponse, error) {
	foundApp, err := s.fhirAppRepo.GetAppById(appId)
	if err != nil {
		return TriggerFHIRJobResponse{}, domain.ErrNotFound
	}

	foundJob, err := s.fhirJobRepo.GetJobByID(jobId)
	if err != nil {
		return TriggerFHIRJobResponse{}, domain.ErrNotFound
	}

	job, err := s.bulkFHIRClient.GetFHIRJobStatus(&foundApp, jobId)
	if err != nil {
		return TriggerFHIRJobResponse{}, err
	}

	if foundJob.Status == "submitted" {
		var resources []domain.FHIRResource
		for _, resource := range job.Resources {
			resources = append(resources, domain.FHIRResource{
				JobID:      jobId,
				AppID:      appId,
				ResourceID: resource.ResourceID,
				Type:       resource.Type})
		}

		if len(resources) > 0 {
			foundJob.Status = domain.JobStatus(job.Status)
			s.fhirJobRepo.CreateOrUpdateJob(&foundJob)

			if err := s.fhirResourceRepo.CreateFHIRResources(resources); err != nil {
				return TriggerFHIRJobResponse{}, err
			}
		}
	}

	return job, nil
}

func (s *service) CreateNewFHIRJob(appId string, body *TriggerFHIRJobRequest) (domain.FHIRJob, error) {
	app, err := s.fhirAppRepo.GetAppById(appId)
	if err != nil {
		return domain.FHIRJob{}, domain.ErrNotFound
	}

	jobId, err := s.bulkFHIRClient.CreateNewJob(&app, body)
	if err != nil {
		return domain.FHIRJob{}, domain.ErrNotFound
	}

	job, err := s.fhirJobRepo.CreateOrUpdateJob(&domain.FHIRJob{ID: jobId, AppID: appId, Status: "submitted", Resources: []domain.FHIRResource{}})
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

	return s.bulkFHIRClient.DeleteFHIRJob(jobId)
}
