package rutube

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ip75/rutubeuploader/internal/log"
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
	Details        string   `json:"detail"`
}

func ParseErr(errBody io.Reader) error {
	desc := errResponse{}

	buf := bytes.Buffer{}
	tee := io.TeeReader(errBody, &buf)

	d := json.NewDecoder(tee)
	if err := d.Decode(&desc); err != nil {
		return fmt.Errorf("decode error response: %w", err)
	}

	log.Logger.Debug().Str("raw error respose", buf.String()).Msg("body with error")

	res := []error{}
	for _, e := range desc.NonFieldErrors {
		res = append(res, errors.New(e))
	}

	if desc.Details != "" {
		res = append(res, errors.New(desc.Details))
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
