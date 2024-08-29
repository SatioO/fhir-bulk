package resource

import (
	"log"

	"github.com/satioO/fhir/v2/domain"
	"gorm.io/gorm"
)

type FHIRResourceRepo struct {
	db *gorm.DB
}

func NewFHIRResourceRepo(db *gorm.DB) *FHIRResourceRepo {
	return &FHIRResourceRepo{db}
}

func (r *FHIRResourceRepo) GetFHIRResourcesForJob(jobId string) ([]domain.FHIRResource, error) {
	var resources []domain.FHIRResource
	result := r.db.Find(&resources)

	return resources, result.Error
}

func (r *FHIRResourceRepo) CreateFHIRResource(entity *domain.FHIRResource) error {
	result := r.db.Save(entity)

	return result.Error
}

func (r *FHIRResourceRepo) CreateFHIRResources(entities []domain.FHIRResource) error {
	log.Println("slice: ", entities)

	result := r.db.Save(entities)

	return result.Error
}
