package discord

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var discordErrMsg = `
%s

You need to set Discord parameters in k8s configmap or using "--url/-u, --username/-n, --avatar_url/-a" or using environment variables:

export KW_DISCORD_WEBHOOK_URL=discord_webhook_url
export KW_DISCORD_USERNAME=discord_username
export KW_DISCORD_AVATAR_URL=avatar_url

Command line flags will override environment variables


` // Green, Yellow, Red
var colorMap = map[string]string{"Info": "5763719", "Warning": "16776960", "Critical": "15548997"}

type DiscordEmbed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

type DiscordWebhook struct {
	Url       string
	Username  string
	AvatarURL string
}

// DiscordMessage
type DiscordMessage struct {
	Username  string          `json:"username"`
	Content   string          `json:"content"`
	AvatarURL string          `json:"avatar_url"`
	Embeds    *[]DiscordEmbed `json:"embeds,omitempty"`
}

// EventMessage for DiscordMessage
type EventMessage struct {
	EventMeta EventMeta `json:"eventmeta"`
	Text      string    `json:"text"`
	Time      time.Time `json:"time"`
}

// EventMeta containes the meta data about the event occurred
type EventMeta struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Reason    string `json:"reason"`
}

// Init prepares Discord configuration
func (m *DiscordWebhook) Init(c *config.Config) error {
	url := c.Handler.Discord.Url
	username := c.Handler.Discord.Username
	avatar := c.Handler.Discord.AvatarURL

	if url == "" {
		url = os.Getenv("KW_DISCORD_WEBHOOK_URL")
	}
	if username == "" {
		username = os.Getenv("KW_DISCORD_USERNAME")
	}
	if avatar == "" {
		avatar = os.Getenv("KW_DISCORD_AVATAR_URL")
	}

	m.Url = url
	m.Username = username
	m.AvatarURL = avatar
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return checkMissingDiscordVars(m)
}

// Handle handles an event.
func (m *DiscordWebhook) Handle(e event.Event) {
	discordMessage := prepareDiscordMessage(e, m)

	err := postMessage(m.Url, discordMessage)
	if err != nil {
		logrus.Printf("%s\n", err)
		return
	}

	logrus.Printf("Message successfully sent to %s at %s ", m.Url, time.Now())
}

func checkMissingDiscordVars(s *DiscordWebhook) error {
	if s.Url == "" {
		return fmt.Errorf(discordErrMsg, "Missing Discord webhook url")
	}
	if s.Username == "" {
		return fmt.Errorf(discordErrMsg, "Missing Discord bot username")
	}
	if s.AvatarURL == "" {
		logrus.Debugf("Missing Discord bot avatar url. Setting default Woof-woof")
		s.AvatarURL = "https://i.imgur.com/oBPXx0D.png"

	}

	return nil
}

func formatEventContent(e event.Event) string {
	eventMeta := fmt.Sprintf("Kind: %s\nName: %s\nNamespace: %s\nReason: %s\n",
		e.Kind, e.Name, e.Namespace, e.Reason)
	eventText := fmt.Sprintf("Text: %s\n", e.Message())
	eventTime := fmt.Sprintf("UTC Time: %s\n", time.Now().Format(time.RFC3339))
	return eventMeta + eventText + eventTime
}

func prepareDiscordEmbeds(e event.Event) *DiscordEmbed {
	return &DiscordEmbed{Title: e.Status, Description: e.InfoMessage, Color: colorMap[e.Status]}
}

func prepareDiscordMessage(e event.Event, m *DiscordWebhook) *DiscordMessage {

	eventContent := formatEventContent(e)
	formattedContent := fmt.Sprintf("```md\n%s\n```", eventContent)
	embeds := prepareDiscordEmbeds(e)
	return &DiscordMessage{
		Username:  m.Username,
		Content:   formattedContent,
		AvatarURL: m.AvatarURL,
		Embeds:    &[]DiscordEmbed{*embeds},
	}

}

func postMessage(url string, discordMessage *DiscordMessage) error {
	message, err := json.Marshal(discordMessage)
	logrus.Debugf("Marshaled JSON message: %s", string(message))
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error sending HTTP request: %s", err)
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		logrus.Errorf("Unexpected response status: %d", resp.StatusCode)
		responseBody, _ := ioutil.ReadAll(resp.Body)
		logrus.Debugf("Response body: %s", responseBody)
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	return nil
}
