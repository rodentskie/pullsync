package slack

import (
	"slack-pr-lambda/env"
	"slack-pr-lambda/types"

	"github.com/slack-go/slack"
)

func SlackSendMessage(input types.OpenPullRequest, msg string) (string, error) {
	token := env.GetEnv("SLACK_TOKEN", "")
	channel := env.GetEnv("SLACK_CHANNEL", "")
	api := slack.New(token)

	_, timestamp, err := api.PostMessage(
		channel,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAsUser(false),
	)

	if err != nil {
		return "", err
	}

	return timestamp, nil
}

func SlackSendMessageThread(timeStamp string, message string) error {
	token := env.GetEnv("SLACK_TOKEN", "")
	channel := env.GetEnv("SLACK_CHANNEL", "")
	api := slack.New(token)

	_, _, err := api.PostMessage(
		channel,
		slack.MsgOptionText(message, false),
		slack.MsgOptionTS(timeStamp),
	)
	if err != nil {
		return err
	}
	return nil
}
