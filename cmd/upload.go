package cmd

import (
	"os"

	"github.com/ip75/rutubeuploader/internal/log"
	"github.com/ip75/rutubeuploader/internal/pusher"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var (
	uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Start upload process to Rutube video server",
		Long: `Start uploading process. Upload information to Rutube API to start upload video to platform.
For example:
  json data from file
  $ rutubeuploader upload video_list.json

  or you can use stdin to pass data
  $ cat video_list.json | rutubeuploader upload

  video_list.json input JSON format:
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
		Run: func(cmd *cobra.Command, args []string) {
			log.Logger.Info().Msg("Start uploading video to Rutube...")

			inputReader := cmd.InOrStdin()

			if len(args) > 0 {
				file, err := os.Open(args[0])
				if err != nil {
					log.Logger.Panic().Err(err).Msg("open file")
				}
				inputReader = file
			} else {
				log.Logger.Info().Msg("get data from stdin")
			}

			p, err := pusher.New(tasksCount, wait, inputReader, tokenFile)
			if err != nil {
				log.Logger.Panic().Err(err).Msg("create pusher")
			}
			p.Run()
		},
		Args: cobra.MaximumNArgs(1),
	}

	wait       bool
	tasksCount int
	tokenFile  string
)

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().BoolVarP(&wait, "wait", "w", false, "Do not exit immediately. Wait processing video on Rutube")
	uploadCmd.Flags().IntVarP(&tasksCount, "tasks", "n", pusher.DefaultTasksCount, "Tasks count run simultaneously")
	uploadCmd.Flags().StringVarP(&tokenFile, "token", "t", defaultTokenFile, "File where token was dumped as a token command result")
}
