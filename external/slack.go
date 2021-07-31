package external

import "github.com/slack-go/slack"

type SlackClient struct {
	client *slack.Client
}

type SlackService interface {
	PostMessage(channel string, option ...slack.MsgOption) error
}

func NewSlack(token string) SlackService {
	c := slack.New(token)
	return SlackClient{
		client: c,
	}
}

func (s SlackClient) PostMessage(channel string, option ...slack.MsgOption) error {
	_, _, _, err := s.client.SendMessage(channel, option...)
	return err
}
