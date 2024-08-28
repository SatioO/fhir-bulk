package fhir_app

import (
	"github.com/satioO/fhir/v2/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FHIRAppRepo struct {
	db *gorm.DB
}

func NewFHIRAppRepo(conn *gorm.DB) *FHIRAppRepo {
	return &FHIRAppRepo{db: conn}
}

func (d *FHIRAppRepo) GetApps() ([]domain.FHIRApp, error) {
	var entity []domain.FHIRApp
	result := d.db.Find(&entity)

	return entity, result.Error
}

func (d *FHIRAppRepo) GetAppById(appId string) (domain.FHIRApp, error) {
	var entity domain.FHIRApp
	result := d.db.Where(domain.FHIRApp{ID: appId, Status: "active"}).First(&entity)

	return entity, result.Error
}

func (d *FHIRAppRepo) CreateApp(body *CreateFHIRAppRequest) (domain.FHIRApp, error) {
	entity := domain.FHIRApp{ID: body.ID, BaseUrl: body.BaseUrl, Status: "active"}
	result := d.db.Omit(clause.Associations).Create(&entity)

	return entity, result.Error
}
