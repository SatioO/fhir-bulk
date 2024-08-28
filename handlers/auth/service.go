package auth

import (
	"github.com/satioO/fhir/v2/domain"
)

type AuthService interface {
	GetAuthServerForApp(string) (domain.FhirAuthServer, error)
	RegisterAuthServer(string, *RegisterAuthServerRequest) (string, error)
}

type authServiceImpl struct {
	authRepo   *AuthRepo
	authClient *authClient
}

func NewAuthService(authRepo *AuthRepo, authClient *authClient) AuthService {
	return &authServiceImpl{authRepo, authClient}
}

func (a *authServiceImpl) GetAuthServerForApp(appId string) (domain.FhirAuthServer, error) {
	return a.authRepo.GetAuthServerForApp(appId)
}

func (a *authServiceImpl) RegisterAuthServer(appId string, body *RegisterAuthServerRequest) (string, error) {
	entity := domain.FhirAuthServer{
		TokenURL:     body.TokenUrl,
		ClientID:     body.ClientID,
		ClientSecret: body.ClientSecret,
		AppID:        appId,
	}

	auth, err := a.authRepo.RegisterAuthServer(entity)

	if err != nil {
		return "", err
	}

	return a.authClient.GenerateToken(&auth)
}
