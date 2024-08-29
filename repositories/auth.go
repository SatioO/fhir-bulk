package repositories

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

func (d *AuthRepo) GetAuthServerForApp(appId string) (domain.FHIRAuthServer, error) {
	var entity domain.FHIRAuthServer
	result := d.db.Where(domain.FHIRAuthServer{AppID: appId}).First(&entity)

	return entity, result.Error
}

func (d *AuthRepo) RegisterAuthServer(entity domain.FHIRAuthServer) (domain.FHIRAuthServer, error) {
	result := d.db.Create(&entity)

	return entity, result.Error
}
