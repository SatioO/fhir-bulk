package domain

import (
	"gorm.io/gorm"
)

type FHIRApp struct {
	gorm.Model
	ID      string
	BaseUrl string
	Token   string
}

func (FHIRApp) TableName() string {
	return "fhir_app"
}

type FHIRAuthServer struct {
	gorm.Model
	TokenURL     string
	ClientID     string
	ClientSecret string
	AppID        string
	Scopes       string
}

func (FHIRAuthServer) TableName() string {
	return "auth_server"
}
