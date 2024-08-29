package bulkapi

import "gorm.io/gorm"

type FHIRJobRepo struct {
	db *gorm.DB
}

func NewFHIRJobRepo(db *gorm.DB) *FHIRJobRepo {
	return &FHIRJobRepo{db}
}
