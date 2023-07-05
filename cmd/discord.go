package cmd

import (
	"github.com/bitnami-labs/kubewatch/config"
)

// webhookConfigCmd represents the webhook subcommand
var discordConfigCmd = &cobra.Command{
	Use:   "discord",
	Short: "specific discord configuration",
	Long:  `specific discord configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		url, err := cmd.Flags().GetString("url")
		if err == nil {
			if len(url) > 0 {
				conf.Handler.Discord.Url = url
			}
		} else {
			logrus.Fatal(err)
		}
		username, err := cmd.Flags().GetString("username")
		if err == nil {
			if len(url) > 0 {
				conf.Handler.Discord.Username = username
			}
		} else {
			logrus.Fatal(err)
		}

		if err = conf.Write(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	webhookConfigCmd.Flags().StringP("url", "u", "", "Specify Discord webhook url")
	webhookConfigCmd.Flags().StringP("username", "", "", "Specify Discord bot username")
}
