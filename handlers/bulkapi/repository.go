package bulkapi

import (
	"github.com/satioO/fhir/v2/domain"
	"gorm.io/gorm"
)

type FHIRJobRepo struct {
	db *gorm.DB
}

func NewFHIRJobRepo(db *gorm.DB) *FHIRJobRepo {
	return &FHIRJobRepo{db}
}

func (r *FHIRJobRepo) GetJobsByApp(appId string) ([]domain.FHIRJob, error) {
	var jobs []domain.FHIRJob
	result := r.db.Where(domain.FHIRJob{AppID: appId}).Find(&jobs)

	return jobs, result.Error
}

func (r *FHIRJobRepo) CreateJob(entity *domain.FHIRJob) (domain.FHIRJob, error) {
	result := r.db.Save(entity)
	return *entity, result.Error
}
