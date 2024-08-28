package auth

import "github.com/satioO/fhir/v2/domain"

type AuthService interface {
	GetAuthServerForApp(string) (domain.FhirAuthServer, error)
	RegisterAuthServer(string, RegisterAuthServerRequest) (domain.FhirAuthServer, error)
}

type authServiceImpl struct {
	authRepo *AuthRepo
}

func NewAuthService(authRepo *AuthRepo) AuthService {
	return &authServiceImpl{authRepo}
}

func (a *authServiceImpl) GetAuthServerForApp(appId string) (domain.FhirAuthServer, error) {
	return a.authRepo.GetAuthServerForApp(appId)
}

func (a *authServiceImpl) RegisterAuthServer(appId string, body RegisterAuthServerRequest) (domain.FhirAuthServer, error) {
	entity := domain.FhirAuthServer{
		TokenURL:     body.TokenUrl,
		ClientID:     body.ClientID,
		ClientSecret: body.ClientSecret,
		AppID:        appId,
		Status:       "active",
	}

	return a.authRepo.RegisterAuthServer(entity)
}
