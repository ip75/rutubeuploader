package rutube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ip75/rutubeuploader/internal/log"
)

type Video struct {
	Url                *url.URL // mandatory
	SuccessCompleteUrl *string
	ErrorCompleteUrl   *string
	QualityReport      *bool
	Author             *int
	Title              *string
	Description        *string
	Hidden             *bool
	CategoryID         *int
	ConverterParams    *string
	Protected          *bool
}

type uploadResponse struct {
	VideoID string `json:"video_id"`
}

func (v *Video) Upload(waitForCompletion bool) error {
	body := bytes.Buffer{}

	d := url.Values{
		"url": []string{v.Url.String()},
	}

	if waitForCompletion {
		d.Add("callback_url", "") // TODO: not implemented yet
		d.Add("errback_url", "")  // TODO: not implemented yet
	}

	body.WriteString(d.Encode())

	r, err := http.NewRequest(http.MethodPost, UploadUrl, &body)
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", fmt.Sprintf("%d", body.Len()))

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("do upload video: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload status: %s: %w", resp.Status, ParseErr(resp.Body))
	}

	response := uploadResponse{}
	json.NewDecoder(resp.Body).Decode(&response)

	log.Logger.Info().Str("video ID", response.VideoID).Msg("upload video")

	return nil
}
