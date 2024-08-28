package fhir_app

import (
	"github.com/satioO/fhir/v2/domain"
	"gorm.io/gorm"
)

type FHIRAppRepo struct {
	db *gorm.DB
}

func NewFHIRAppRepo(db *gorm.DB) *FHIRAppRepo {
	return &FHIRAppRepo{db}
}

func (d *FHIRAppRepo) GetApps() ([]domain.FHIRApp, error) {
	var entity []domain.FHIRApp
	result := d.db.Find(&entity)

	return entity, result.Error
}

func (d *FHIRAppRepo) GetAppById(appId string) (domain.FHIRApp, error) {
	var entity domain.FHIRApp
	result := d.db.First(&entity, appId)

	return entity, result.Error
}

func (d *FHIRAppRepo) CreateApp(body *CreateFHIRAppRequest) (domain.FHIRApp, error) {
	entity := domain.FHIRApp{ID: body.ID, BaseUrl: body.BaseUrl}
	result := d.db.Create(&entity)

	return entity, result.Error
}

func (d *FHIRAppRepo) UpdateToken(appId, token string) error {
	var entity domain.FHIRApp
	result := d.db.First(&entity, appId)
	if result.Error != nil {
		return result.Error
	}

	entity.Token = token
	result = d.db.Save(&entity)

	return result.Error
}
