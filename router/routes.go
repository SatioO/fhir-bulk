package router

import (
	"log"
	"net/http"

	"github.com/satioO/fhir/v2/handlers/auth"
	"github.com/satioO/fhir/v2/handlers/bulkapi"
	"github.com/satioO/fhir/v2/handlers/fhirapp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	addRoutes(r)

	return r
}

func addRoutes(r *http.ServeMux) {
	config := DBServerConfig{
		Username: "root",
		Password: "password",
		Name:     "fhir",
		Host:     "127.0.0.1",
		Port:     "3306",
	}
	conn := ConnectToDB(&config)

	authClient := auth.NewAuthClient()

	fhirAppRepo := fhirapp.NewFHIRAppRepo(conn)
	authRepo := auth.NewAuthRepo(conn)
	fhirJobRepo := bulkapi.NewFHIRJobRepo(conn)

	authService := auth.NewAuthService(authRepo, authClient)
	fhirAppService := fhirapp.NewFHIRAppService(fhirAppRepo)
	bulkApiService := bulkapi.NewBulkAPIService(fhirJobRepo)

	fhirAppHandler := fhirapp.NewFHIRAppHandler(fhirAppService, authService)
	authServerHandler := auth.NewAuthHandler(authService)
	bulkApiHandler := bulkapi.NewBulkAPIHandler(bulkApiService)

	r.HandleFunc("GET /api/v1/fhir/apps", fhirAppHandler.GetApps)
	r.HandleFunc("GET /api/v1/fhir/apps/{appId}", fhirAppHandler.GetAppById)
	r.HandleFunc("POST /api/v1/fhir/apps", fhirAppHandler.RegisterApp)

	r.HandleFunc("GET /api/v1/fhir/apps/{appId}/auth", authServerHandler.GetAuthServerForApp)
	r.HandleFunc("POST /api/v1/fhir/apps/{appId}/auth", authServerHandler.RegisterAuthServer)

	r.HandleFunc("GET /api/v1/fhir/apps/{appId}/job", bulkApiHandler.GetFHIRJobsForApp)
	r.HandleFunc("POST /api/v1/fhir/apps/{appId}/job", bulkApiHandler.CreateNewFHIRJob)
}

type DBServerConfig struct {
	Username string
	Password string
	Name     string
	Host     string
	Port     string
}

type DBOption func(*DBServerConfig)

func ConnectToDB(config *DBServerConfig) *gorm.DB {
	dsn := config.Username + ":" + config.Password + "@tcp" + "(" + config.Host + ":" + config.Port + ")/" + config.Name + "?" + "parseTime=true&loc=Local"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return conn
}
