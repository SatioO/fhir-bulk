package bulkapi

import (
	"fmt"

	"github.com/satioO/fhir/v2/domain"
	"github.com/satioO/fhir/v2/handlers/fhirapp"
	"github.com/satioO/fhir/v2/handlers/resource"
)

type BulkAPIService interface {
	GetJobsByApp(appId string) ([]domain.FHIRJob, error)
	GetFHIRJobStatus(appId, jobId string) (TriggerFHIRJobResponse, error)
	CreateNewFHIRJob(appId string, body *TriggerFHIRJobRequest) (domain.FHIRJob, error)
	DeleteFHIRJob(appId, jobId string) error
}

type service struct {
	fhirJobRepo      *FHIRJobRepo
	fhirAppRepo      *fhirapp.FHIRAppRepo
	fhirResourceRepo *resource.FHIRResourceRepo
	bulkFHIRClient   *client
}

func NewBulkAPIService(fhirJobRepo *FHIRJobRepo, fhirAppRepo *fhirapp.FHIRAppRepo, fhirResourceRepo *resource.FHIRResourceRepo, bulkFHIRClient *client) BulkAPIService {
	return &service{fhirJobRepo, fhirAppRepo, fhirResourceRepo, bulkFHIRClient}
}

func (s *service) GetJobsByApp(appId string) ([]domain.FHIRJob, error) {
	return s.fhirJobRepo.GetJobsByApp(appId)
}

func (s *service) GetFHIRJobStatus(appId, jobId string) (TriggerFHIRJobResponse, error) {
	app, err := s.fhirAppRepo.GetAppById(appId)

	if err != nil {
		return TriggerFHIRJobResponse{}, err
	}

	foundJob, err := s.fhirJobRepo.GetJobByID(jobId)

	if err != nil {
		return TriggerFHIRJobResponse{}, err
	}

	job, err := s.bulkFHIRClient.GetFHIRJobStatus(&app, jobId)
	if err != nil {
		return TriggerFHIRJobResponse{}, err
	}

	if foundJob.Status == "submitted" {
		foundJob.Status = domain.JobStatus(job.Status)
		s.fhirJobRepo.UpdateJob(&foundJob)

		foundResources, err := s.fhirResourceRepo.GetFHIRResourcesForJob(jobId)
		if err != nil {
			return TriggerFHIRJobResponse{}, err
		}

		if len(foundResources) == 0 {
			var resources []domain.FHIRResource
			for _, resource := range job.Resources {
				resources = append(resources, domain.FHIRResource{
					JobID:      jobId,
					AppID:      appId,
					ResourceID: resource.ResourceID,
					Type:       resource.Type})
			}

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
