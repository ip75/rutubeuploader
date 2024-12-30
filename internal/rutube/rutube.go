package rutube

import (
	"encoding/json"
	"errors"
	"io"
)

const (
	TokenUrl  = "https://rutube.ru/api/accounts/token_auth/"
	UploadUrl = "https://rutube.ru/api/video/"
)

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type errResponse struct {
	NonFieldErrors []string `json:"non_field_errors"`
}

func ParseErr(errBody io.Reader) error {
	desc := errResponse{}
	d := json.NewDecoder(errBody)
	d.Decode(&desc)
	res := []error{}
	for _, e := range desc.NonFieldErrors {
		res = append(res, errors.New(e))
	}
	return errors.Join(res...)
}

type TokenResponse struct {
	Token string `json:"token"`
}

func ParseToken(token io.Reader) *TokenResponse {
	t := TokenResponse{}
	d := json.NewDecoder(token)
	d.Decode(&t)
	return &t
}
