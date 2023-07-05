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

You need to set Webhook url using "--url/-u, --username/-n" or using environment variables:

export KW_DISCORD_WEBHOOK_URL=discord_webhook_url
export KW_DISCORD_USERNAME=discord_username

Command line flags will override environment variables

`

type DiscordWebhook struct {
	Url      string
	Username string
}

// DiscordMessage
type DiscordMessage struct {
	Username string `json:"username"`
	Content  string `json:"content"`
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

	if url == "" {
		url = os.Getenv("KW_DISCORD_WEBHOOK_URL")
	}
	if username == "" {
		username = os.Getenv("KW_DISCORD_USERNAME")
	}

	m.Url = url
	m.Username = username
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return checkMissingDiscordVars(m)
}

// Handle handles an event.
func (m *DiscordWebhook) Handle(e event.Event) {
	DiscordMessage := prepareDiscordMessage(e, m)

	err := postMessage(m.Url, DiscordMessage)
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

	return nil
}

func formatEventContent(e event.Event) string {
	eventMeta := fmt.Sprintf("Kind: %s\nName: %s\nNamespace: %s\nReason: %s\n",
		e.Kind, e.Name, e.Namespace, e.Reason)
	eventText := fmt.Sprintf("Text: %s\n", e.Message())
	eventTime := fmt.Sprintf("Time: %s\n", time.Now().Format(time.RFC3339))

	return eventMeta + eventText + eventTime
}

func prepareDiscordMessage(e event.Event, m *DiscordWebhook) *DiscordMessage {

	eventContent := formatEventContent(e)
	formattedContent := fmt.Sprintf("```md\n%s\n```", eventContent)

	return &DiscordMessage{
		Username: m.Username,
		Content:  formattedContent,
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

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Unexpected response status: %d", resp.StatusCode)
		responseBody, _ := ioutil.ReadAll(resp.Body)
		logrus.Debugf("Response body: %s", responseBody)
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	return nil
}
