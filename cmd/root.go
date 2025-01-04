package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rutubeuploader",
	Short: "A CLI to upload video to rutube.ru",
	Long: `This is a command line application to upload video content to video hoster rutube.ru.
	For example:
		$ rutubeuploader token --user=compas --password=123  # to generate token.json
		$ rutubeuploader upload video_list.json              # to start upload video to rutube from file
	or
		$ cat video_list.json | rutubeuploader upload        # --||-- from stdin

# video_list.json input JSON format
[
    {
        "url": "https://user:password@videoserver.su/video/alien.mp4",  // url to video where Rutube will get from. Set necessary credentials to url to download video successfully. Available schema: https/http, ftp
        "quality_report": true,       // if true, notification will be called every time when video will be converted to every step of quality, if false notification will be called once when all convertions will be completed
        "author": 123,                // Identity of author. This author has to have access to upload video to specified channel
        "title": "Alien",             // video title, max 100 runes
        "description": "Alien movie", // description of video, max 5000 runes
        "hidden": false,              // true: private, false: public video 
        "category_id": 13,            // Identity number of category for video. Default 13
        "converter_params": "",       // additional video convertion parameters. Example xml-tagging: converter_params=%7B%22editor_xml%22%3A%22ftp%3A%5C%2F%5C%2Frutube%3pass%4010.122.50.222%5C%2FPR291117-A.xml%22%7D
        "protected": false            // true: video will be queued to DRM checking
    },
    /// etc.
]
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rutubeuploader.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initLogger() {
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".rutubeuploader" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".rutubeuploader")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
