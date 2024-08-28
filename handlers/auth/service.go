package auth

import (
	"github.com/satioO/fhir/v2/domain"
)

type AuthService interface {
	GetAuthServerForApp(string) (domain.FHIRAuthServer, error)
	RegisterAuthServer(string, *RegisterAuthServerRequest) (string, error)
}

type service struct {
	authRepo   *AuthRepo
	authClient *authClient
}

func NewAuthService(authRepo *AuthRepo, authClient *authClient) AuthService {
	return &service{authRepo, authClient}
}

func (a *service) GetAuthServerForApp(appId string) (domain.FHIRAuthServer, error) {
	return a.authRepo.GetAuthServerForApp(appId)
}

func (a *service) RegisterAuthServer(appId string, body *RegisterAuthServerRequest) (string, error) {
	entity := domain.FHIRAuthServer{
		TokenURL:     body.TokenUrl,
		ClientID:     body.ClientID,
		ClientSecret: body.ClientSecret,
		Scopes:       body.Scopes,
		AppID:        appId,
	}

	auth, err := a.authRepo.RegisterAuthServer(entity)

	if err != nil {
		return "", err
	}

	return a.authClient.GenerateToken(&auth)
}
