package router

import (
	"log"
	"net/http"

	"github.com/satioO/fhir/v2/handlers/auth"
	"github.com/satioO/fhir/v2/handlers/fhir_app"
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

	fhirAppRepo := fhir_app.NewFHIRAppRepo(conn)
	authRepo := auth.NewAuthRepo(conn)

	authService := auth.NewAuthService(authRepo, authClient)

	fhirAppHandler := fhir_app.NewFHIRAppHandler(fhirAppRepo, authService)
	authServerHandler := auth.NewAuthHandler(authService)

	r.HandleFunc("GET /api/v1/fhir/app", fhirAppHandler.GetApps)
	r.HandleFunc("GET /api/v1/fhir/app/{appId}", fhirAppHandler.GetAppById)
	r.HandleFunc("POST /api/v1/fhir/app", fhirAppHandler.RegisterApp)

	r.HandleFunc("GET /api/v1/fhir/app/{appId}/auth", authServerHandler.GetAuthServerForApp)
	r.HandleFunc("POST /api/v1/fhir/app/{appId}/auth", authServerHandler.RegisterAuthServer)
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
