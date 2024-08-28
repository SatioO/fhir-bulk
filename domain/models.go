package domain

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

type Status string

func (p *Status) Scan(value interface{}) error {
	*p = Status(value.([]byte))
	return nil
}

func (p Status) Value() (driver.Value, error) {
	return string(p), nil
}

const (
	Active   Status = "active"
	InActive Status = "inactive"
)

type FHIRApp struct {
	ID        string
	BaseUrl   string
	Token     string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (FHIRApp) TableName() string {
	return "fhir_app"
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type FhirAuthServer struct {
	gorm.Model
	TokenURL     string
	ClientID     string
	ClientSecret string
	AppID        string
}

func (FhirAuthServer) TableName() string {
	return "auth_server"
}
