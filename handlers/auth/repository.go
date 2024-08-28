package auth

import (
	"github.com/satioO/fhir/v2/domain"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db}
}

func (d *AuthRepo) GetAuthServerForApp(appId string) (domain.FhirAuthServer, error) {
	var entity domain.FhirAuthServer
	result := d.db.Where(domain.FhirAuthServer{AppID: appId}).First(&entity)

	return entity, result.Error
}

func (d *AuthRepo) RegisterAuthServer(entity domain.FhirAuthServer) (domain.FhirAuthServer, error) {
	result := d.db.Create(&entity)

	return entity, result.Error
}
