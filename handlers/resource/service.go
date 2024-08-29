package resource

import (
	"github.com/satioO/fhir/v2/domain"
	"github.com/satioO/fhir/v2/repositories"
)

type FHIRResourceService interface {
	GetFHIRResourcesByJobID(jobId string) ([]domain.FHIRResource, error)
	GetFHIRResource(jobId, resourceId string) ([]byte, error)
}

type service struct {
	fhirResourceRepo   *repositories.FHIRResourceRepo
	fhirJobRepo        *repositories.FHIRJobRepo
	fhirAppRepo        *repositories.FHIRAppRepo
	fhirResourceClient *client
}

func NewFHIRResourceService(fhirResourceRepo *repositories.FHIRResourceRepo, fhirJobRepo *repositories.FHIRJobRepo, fhirAppRepo *repositories.FHIRAppRepo, client *client) FHIRResourceService {
	return &service{fhirResourceRepo, fhirJobRepo, fhirAppRepo, client}
}

func (s *service) GetFHIRResourcesByJobID(jobId string) ([]domain.FHIRResource, error) {
	return s.fhirResourceRepo.GetFHIRResourcesForJob(jobId)
}

func (s *service) GetFHIRResource(jobId, resourceId string) ([]byte, error) {
	foundJob, err := s.fhirJobRepo.GetJobByID(jobId)
	if err != nil {
		return nil, err
	}

	foundApp, err := s.fhirAppRepo.GetAppById(foundJob.AppID)
	if err != nil {
		return nil, err
	}

	return s.fhirResourceClient.GetFHIRResource(&foundApp, resourceId)
}
