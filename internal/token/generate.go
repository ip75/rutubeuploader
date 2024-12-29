package token

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ip75/rutubeuploader/internal/log"
	"github.com/ip75/rutubeuploader/internal/rutube"
)

func (t *Token) Generate(user, pass string, regenerate bool) error {
	log.Logger.Info().Msg("generate token...")

	body := bytes.Buffer{}

	d := url.Values{
		"username": []string{user},
		"password": []string{pass},
	}.Encode()

	body.WriteString(d)

	method := http.MethodPost
	if regenerate {
		method = http.MethodPut
	}

	r, err := http.NewRequest(method, rutube.TokenUrl, &body)
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", fmt.Sprintf("%d", body.Len()))

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Logger.Error().Msgf("do request token: %s", err)
		return fmt.Errorf("do request token: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Logger.Error().Msgf("error when generate token: %s", resp.Status)
		return fmt.Errorf("error when generate token: %s", resp.Status)
	}

	return nil
}

func (t *Token) SaveToken(path string) error {
	log.Logger.Info().Msgf("save token to %s...", path)
	return nil
}
