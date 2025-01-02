package rutube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ip75/rutubeuploader/internal/log"
)

type Video struct {
	Url                string  `json:"url"`
	SuccessCompleteUrl *string `json:"success_complete_url,omitempty"`
	ErrorCompleteUrl   *string `json:"error_complete_url,omitempty"`
	QualityReport      *bool   `json:"quality_report,omitempty"`
	Author             *int    `json:"author,omitempty"`
	Title              *string `json:"title,omitempty"`
	Description        *string `json:"description,omitempty"`
	Hidden             *bool   `json:"hidden,omitempty"`
	CategoryID         *int    `json:"category_id,omitempty"`
	ConverterParams    *string `json:"converter_params,omitempty"`
	Protected          *bool   `json:"protected,omitempty"`
}

type uploadResponse struct {
	VideoID string `json:"video_id"`
}

func (v *Video) Upload(tok string, waitForCompletion bool) error {
	body := bytes.Buffer{}

	d := url.Values{
		"url": []string{v.Url},
	}

	if v.QualityReport != nil {
		d.Add("quality_report", strconv.FormatBool(*v.QualityReport))
	}

	if v.Hidden != nil {
		d.Add("is_hidden", strconv.FormatBool(*v.Hidden))
	}

	if v.Title != nil {
		d.Add("title", *v.Title)
	}

	if v.Description != nil {
		d.Add("description", *v.Description)
	}

	if v.CategoryID != nil {
		d.Add("category_id", strconv.Itoa(*v.CategoryID))
	}

	if v.Author != nil {
		d.Add("author", strconv.Itoa(*v.Author))
	}

	if v.Protected != nil {
		d.Add("protected", strconv.FormatBool(*v.Protected))
	}

	if waitForCompletion {
		// TODO: not implemented yet.
		log.Logger.Warn().Msg("waiting for completion is not implemented yet")
		// d.Add("callback_url", "")
		// d.Add("errback_url", "")
	}

	body.WriteString(d.Encode())

	r, err := http.NewRequest(http.MethodPost, UploadUrl, &body)
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", fmt.Sprintf("%d", body.Len()))
	r.Header.Add("Authorization", tok)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("do upload video: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload status: %s: %w", resp.Status, ParseErr(resp.Body))
	}

	response := uploadResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Logger.Error().Err(err).Msg("decode response")
	}

	log.Logger.Info().Str("video ID", response.VideoID).Msg("video uploaded successfully")

	return nil
}
