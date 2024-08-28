package auth

type RegisterAuthServerRequest struct {
	TokenUrl     string `json:"tokenUrl"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Scopes       string `json:"scopes"`
}

type AuthClientResponse struct {
	AccessToken string `json:"access_token"`
}
