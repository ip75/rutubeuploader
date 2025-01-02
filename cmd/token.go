package cmd

import (
	"github.com/ip75/rutubeuploader/internal/log"
	"github.com/ip75/rutubeuploader/internal/token"
	"github.com/spf13/cobra"
)

const (
	defaultTokenFile = "token.json"
)

var (
	tokenCmd = &cobra.Command{
		Use:   "token",
		Short: "generate token to access rutube platform API",
		Long: `Use credentials to authorize to perform media operations on rutube platform.
File token.json with access token will be created at current directory and then will be used to perform API calls for authentication.
For example:
	rutubeuploader token --user=compas --password=123
		`,
		Run: func(cmd *cobra.Command, args []string) {
			t := token.Token{}
			if err := t.Generate(username, password, regenate); err != nil {
				log.Logger.Err(err).Msg("token: generate")
				return
			}
			if err := t.SaveToken(defaultTokenFile); err != nil {
				log.Logger.Err(err).Msg("token: save to file")
				return
			}
		},
	}

	username, password string
	regenate           bool
)

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().StringVarP(&username, "user", "u", "", "username of your account on rutube.ru")
	tokenCmd.MarkFlagRequired("user")
	tokenCmd.Flags().StringVarP(&password, "password", "p", "", "password to access your account on rutube.ru")
	tokenCmd.MarkFlagRequired("password")
	tokenCmd.Flags().BoolVarP(&regenate, "regenerate", "r", false, "regenerate token to access rutube.ru API")
}
