package slack

import (
	"fmt"
	"slack-pr-lambda/env"
	"slack-pr-lambda/types"

	"github.com/slack-go/slack"
)

func SlackSendMessage(input types.OpenPullRequest) (string, error) {
	token := env.GetEnv("SLACK_TOKEN", "")
	channel := env.GetEnv("SLACK_CHANNEL", "")
	api := slack.New(token)

	messageText := fmt.Sprintf("New <%s|pull request> in `%s`.", input.PullRequest.HtmlUrl, input.Repository.Name)

	_, timestamp, err := api.PostMessage(
		channel,
		slack.MsgOptionText(messageText, false),
		slack.MsgOptionAsUser(false),
	)

	if err != nil {
		return "", err
	}

	return timestamp, nil
}

func SlackSendMessageThreadClosed(timeStamp string) error {
	token := env.GetEnv("SLACK_TOKEN", "")
	channel := env.GetEnv("SLACK_CHANNEL", "")
	api := slack.New(token)

	messageText := "PR is closed."

	_, _, err := api.PostMessage(
		channel,
		slack.MsgOptionText(messageText, false),
		slack.MsgOptionTS(timeStamp),
	)
	if err != nil {
		return err
	}
	return nil
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
