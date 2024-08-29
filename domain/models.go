package domain

import (
	"gorm.io/gorm"
)

type JobStatus string

type FHIRApp struct {
	gorm.Model
	ID      string
	BaseUrl string
	Token   string
}

func (FHIRApp) TableName() string {
	return "fhir_apps"
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
	return "auth_servers"
}

type FHIRJob struct {
	gorm.Model
	JobID  string
	AppID  string
	Status JobStatus
}

func (FHIRJob) TableName() string {
	return "fhir_jobs"
}
