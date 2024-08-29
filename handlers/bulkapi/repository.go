package bulkapi

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

func NewFHIRJobRepo(db *gorm.DB) *repository {
	return &repository{db}
}
