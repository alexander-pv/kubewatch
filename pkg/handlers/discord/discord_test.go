package discord

import (
	"fmt"
	"github.com/bitnami-labs/kubewatch/config"
	"reflect"
	"testing"
)

func TestWebhookInit(t *testing.T) {
	s := &config.Discord{}

	var Tests = []struct {
		discordwebhook config.Discord
		err            error
	}{
		{config.Discord{Url: "foo", Username: "bar"}, nil},
		{config.Discord{Url: "foo"}, fmt.Errorf(discordErrMsg, "Missing Discord bot Username")},
		{config.Discord{Username: "bar"}, fmt.Errorf(discordErrMsg, "Missing Discord webhook URL")},
		{config.Discord{}, fmt.Errorf(discordErrMsg, "Missing Discord webhook URL")},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.Discord = tt.discordwebhook
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}
