package token

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ip75/rutubeuploader/internal/log"
)

type Token struct {
	Token string `json:"token"`
}

func (t *Token) SaveToken(path string) error {
	log.Logger.Info().Msgf("save token to %s...", path)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("open/create token file: %w", err)
	}
	defer f.Close()

	d, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}
	if _, err := f.Write(d); err != nil {
		return fmt.Errorf("save token: %w", err)
	}

	return nil
}

func (t *Token) Authorization() string {
	return t.Token
}

func (t *Token) LoadToken(path string) error {
	log.Logger.Info().Msgf("save token to %s...", path)

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("check token file exist: %w", err)
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("open token file: %w", err)
	}

	err = json.Unmarshal(f, t)
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}

	return nil
}
