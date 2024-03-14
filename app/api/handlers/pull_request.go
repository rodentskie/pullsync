package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slack-pr-lambda/logger"
	"syscall"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func PullRequestHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.LoggerConfig()
	zapLog, _ := l.Build()
	defer func() {
		err := r.Body.Close()
		if err != nil {
			zapLog.Fatal("error close req body",
				zap.Error(err),
			)
		}

		if err := zapLog.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			log.Fatalf("error closing the logger. %v\n", err)
		}
	}()
	api := slack.New("token")
	inlineCode := "code"
	codeBlock := `func main() {
    fmt.Println("Hello, world!")
}`

	// Use fmt.Sprintf to construct the message text with dynamic content
	messageText := fmt.Sprintf(`Here are some bullet points:
• Item 1
• Item 2
• Item 3

And here is an inline code: `+"`%s`"+`

And a code block:
`+"```%s```"+`

`, inlineCode, codeBlock)

	channelID, timestamp, err := api.PostMessage(
		"general",                               // Channel name. Ensure your bot is a member of this channel.
		slack.MsgOptionText(messageText, false), // Passing the message with bullet points
		slack.MsgOptionAsUser(true),             // Send as a user, not as a bot
	)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

	threadTimestamp := "1710406869.549219"
	_, _, err = api.PostMessage(
		channelID,
		slack.MsgOptionText("This is a reply in a thread! v2", false),
		slack.MsgOptionTS(threadTimestamp),
	)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zapLog.Fatal("error read request body",
			zap.Error(err),
		)
	}

	bodyBytes := Response{
		Message: string(body),
	}

	j, err := json.Marshal(bodyBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
