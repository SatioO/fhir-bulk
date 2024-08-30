package router

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/satioO/fhir/v2/api"
	"github.com/satioO/fhir/v2/domain"
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

	r.HandleFunc("GET /auth-servers", authServerHandler.GetAuthServerForApp)

	r.HandleFunc("GET /jobs", bulkApiHandler.GetFHIRJobsForApp)
	r.HandleFunc("GET /jobs/{jobId}", bulkApiHandler.GetFHIRJobStatus)
	r.HandleFunc("GET /jobs/{jobId}/resources", resourceHandler.GetFHIRResourcesByJobID)
	r.HandleFunc("GET /jobs/{jobId}/resources/{resourceId}", resourceHandler.GetFHIRResource)
	r.HandleFunc("POST /jobs", bulkApiHandler.CreateNewFHIRJob)
	r.HandleFunc("DELETE /jobs/{jobId}", bulkApiHandler.DeleteFHIRJob)
}

func AppMiddleware(next http.HandlerFunc, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appId := r.Header.Get("X-App-Id")

		var app domain.FHIRApp
		result := db.Where(domain.FHIRApp{ID: appId}).First(&app)

		if result.Error != nil {
			api.Error(w, r, fmt.Errorf("could not find app: %s", appId), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "app", app)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
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
