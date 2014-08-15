package slack

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

var (
	SLACK_API   = "https://curtmfg.slack.com/services/hooks/incoming-webhook"
	SLACK_TOKEN = "yqtFfAU6BFA8AaAsEaYcaf4D"
)

type Message struct {
	Channel  string `json:"channel"`
	Username string `json:"username,omitempty"`
	Text     string `json:"text"`
	Icon     string `json:"icon_emoji,omitempty"`
}

func (m *Message) Send() error {
	if len(m.Channel) == 0 {
		return errors.New("Must specify a slack channel!")
	}

	if len(m.Text) == 0 {
		return errors.New("Must specifty text for slack message!")
	}

	js, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := http.PostForm(SLACK_API, url.Values{
		"token":   {SLACK_TOKEN},
		"payload": {string(js)},
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
