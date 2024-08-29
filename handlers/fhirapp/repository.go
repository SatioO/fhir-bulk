package fhirapp

import (
	"github.com/satioO/fhir/v2/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewFHIRAppRepo(db *gorm.DB) *repository {
	return &repository{db}
}

func (d *repository) GetApps() ([]domain.FHIRApp, error) {
	var entity []domain.FHIRApp
	result := d.db.Find(&entity)

	return entity, result.Error
}

func (d *repository) GetAppById(appId string) (domain.FHIRApp, error) {
	var entity domain.FHIRApp
	result := d.db.Where(domain.FHIRApp{ID: appId}).First(&entity)

	return entity, result.Error
}

func (d *repository) CreateApp(entity domain.FHIRApp) (domain.FHIRApp, error) {
	result := d.db.Create(&entity)
	return entity, result.Error
}

func (d *repository) UpdateToken(appId, token string) error {
	var entity domain.FHIRApp
	result := d.db.Where(domain.FHIRApp{ID: appId}).First(&entity)
	if result.Error != nil {
		return result.Error
	}

	entity.Token = token
	result = d.db.Save(&entity)

	return result.Error
}
