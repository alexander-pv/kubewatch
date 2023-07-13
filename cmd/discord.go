package cmd

import (
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

		avatar, err := cmd.Flags().GetString("avatar_url")
		if err == nil {
			if len(avatar) > 0 {
				conf.Handler.Discord.AvatarURL = avatar
			}
		}

		if err = conf.Write(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	discordConfigCmd.Flags().StringP("url", "u", "", "Specify Discord webhook url")
	discordConfigCmd.Flags().StringP("username", "n", "", "Specify Discord bot username")
	discordConfigCmd.Flags().StringP("avatar_url", "a", "", "Specify Discord bot avatar url")
}
