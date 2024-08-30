package auth

import validation "github.com/go-ozzo/ozzo-validation"

func validateAppId(appId string) error {
	return validation.Validate(appId, validation.Required)
}
