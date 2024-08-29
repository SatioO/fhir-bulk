package fhirapp

import "github.com/satioO/fhir/v2/domain"

type FHIRAppService interface {
	GetApps() ([]domain.FHIRApp, error)
	GetAppById(appId string) (domain.FHIRApp, error)
	CreateApp(payload *CreateFHIRAppRequest) (domain.FHIRApp, error)
	UpdateToken(appId, token string) error
}

type service struct {
	appRepo *FHIRAppRepo
}

func NewFHIRAppService(appRepo *FHIRAppRepo) FHIRAppService {
	return &service{appRepo}
}

func (s *service) CreateApp(body *CreateFHIRAppRequest) (domain.FHIRApp, error) {
	payload := domain.FHIRApp{ID: body.ID, BaseUrl: body.BaseUrl}
	return s.appRepo.CreateApp(payload)
}

func (s *service) GetAppById(appId string) (domain.FHIRApp, error) {
	return s.appRepo.GetAppById(appId)
}

func (s *service) GetApps() ([]domain.FHIRApp, error) {
	return s.appRepo.GetApps()
}

func (s *service) UpdateToken(appId string, token string) error {
	return s.appRepo.UpdateToken(appId, token)
}
