package router

import (
	"log"
	"net/http"

	"github.com/satioO/fhir/v2/handlers/auth"
	"github.com/satioO/fhir/v2/handlers/bulkapi"
	"github.com/satioO/fhir/v2/handlers/fhirapp"
	"github.com/satioO/fhir/v2/handlers/resource"
	"github.com/satioO/fhir/v2/repositories"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func RegisterRoutes() *http.ServeMux {
	v1 := http.NewServeMux()
	addRoutes(v1)

	main := http.NewServeMux()
	main.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	return main
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
	bulkFHIRClient := bulkapi.NewBulkFHIRClient()
	fhirResourceClient := resource.NewFHIRResourceClient()

	fhirAppRepo := repositories.NewFHIRAppRepo(conn)
	authRepo := repositories.NewAuthRepo(conn)
	fhirJobRepo := repositories.NewFHIRJobRepo(conn)
	fhirResourceRepo := repositories.NewFHIRResourceRepo(conn)

	authService := auth.NewAuthService(authRepo, authClient)
	fhirAppService := fhirapp.NewFHIRAppService(fhirAppRepo)
	bulkApiService := bulkapi.NewBulkAPIService(fhirJobRepo, fhirAppRepo, fhirResourceRepo, bulkFHIRClient)
	resourceService := resource.NewFHIRResourceService(fhirResourceRepo, fhirJobRepo, fhirAppRepo, fhirResourceClient)

	fhirAppHandler := fhirapp.NewFHIRAppHandler(fhirAppService, authService)
	authServerHandler := auth.NewAuthHandler(authService)
	bulkApiHandler := bulkapi.NewBulkAPIHandler(bulkApiService)
	resourceHandler := resource.NewFHIRResourceHandler(resourceService)

	r.HandleFunc("GET /apps", fhirAppHandler.GetApps)
	r.HandleFunc("GET /apps/{appId}", fhirAppHandler.GetAppById)
	r.HandleFunc("POST /apps", fhirAppHandler.RegisterApp)

	r.HandleFunc("GET /apps/{appId}/auth", authServerHandler.GetAuthServerForApp)
	r.HandleFunc("POST /apps/{appId}/auth", authServerHandler.RegisterAuthServer)

	r.HandleFunc("GET /apps/{appId}/jobs", bulkApiHandler.GetFHIRJobsForApp)
	r.HandleFunc("GET /apps/{appId}/jobs/{jobId}", bulkApiHandler.GetFHIRJobStatus)
	r.HandleFunc("POST /apps/{appId}/jobs", bulkApiHandler.CreateNewFHIRJob)
	r.HandleFunc("DELETE /apps/{appId}/jobs/{jobId}", bulkApiHandler.DeleteFHIRJob)

	r.HandleFunc("GET /jobs/{jobId}/resources", resourceHandler.GetFHIRResourcesByJobID)
	r.HandleFunc("GET /jobs/{jobId}/resources/{resourceId}", resourceHandler.GetFHIRResource)
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
