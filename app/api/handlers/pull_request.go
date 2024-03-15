package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slack-pr-lambda/constants"
	db "slack-pr-lambda/dynamodb"
	"slack-pr-lambda/env"
	"slack-pr-lambda/github"
	"slack-pr-lambda/logger"
	"slack-pr-lambda/mapstruct"
	"slack-pr-lambda/slack"
	"slack-pr-lambda/types"
	"syscall"

	"go.uber.org/zap"
)

func PullRequestHandler(w http.ResponseWriter, r *http.Request) {
	env := env.GetEnv("ENV", "local")

	l := logger.LoggerConfig()
	zapLog, _ := l.Build()

	slackUsers := constants.SlackUsers()
	slackUsersMap := mapstruct.StructToMapInterface(*slackUsers)

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatalf("error close req body. %v\n", err)
		}
	}()

	defer func() {
		if err := zapLog.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			log.Fatalf("error closing the logger. %v\n", err)
		}
	}()

	// read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		zapLog.Fatal("error read request body",
			zap.Error(err),
		)
	}

	if env != "local" {
		fmt.Printf("Payload %v", string(body))
	}

	// partial parse into map string JSON
	var result map[string]json.RawMessage
	if err := json.Unmarshal(body, &result); err != nil {
		zapLog.Error("error unmarshal JSON raw message",
			zap.Error(err),
		)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// get unique action key
	var action string
	if err := json.Unmarshal(result["action"], &action); err != nil {
		log.Fatal(err)
	}

	// Opened new pull request
	if action == "opened" {
		// parse request
		var input types.OpenPullRequest
		err = json.Unmarshal(body, &input)
		if err != nil {
			zapLog.Error("error unmarshal JSON",
				zap.Error(err),
			)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		messageText := fmt.Sprintf("<@%s> opened new <%s|pull request> in `%s`.", slackUsersMap[input.PullRequest.User.Login], input.PullRequest.HtmlUrl, input.Repository.Name)
		timeStamp, err := slack.SlackSendMessage(input, messageText)
		if err != nil {
			zapLog.Error("error slack send message",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if len(input.PullRequest.RequestedReviewers) > 0 {
			reviewers := []string{}
			for _, reviewer := range input.PullRequest.RequestedReviewers {
				reviewers = append(reviewers, reviewer.Login)
			}

			var slackMention string = "Please review: "
			for _, user := range reviewers {
				slackMention += fmt.Sprintf("<@%s>", slackUsersMap[user])
			}
			err = slack.SlackSendMessageThread(timeStamp, slackMention)
			if err != nil {
				zapLog.Error("error slack send message",
					zap.Error(err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

		}

		svc := db.DynamoDbConnection()
		item := &types.TablePullRequestData{
			ID:             fmt.Sprintf("%d", input.PullRequest.ID),
			PullRequestId:  input.Number,
			SlackTimeStamp: timeStamp,
		}

		err = db.InsertItem(svc, item)
		if err != nil {
			zapLog.Error("error insert data",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Add new reviewer
	if action == "review_requested" {
		// parse request
		var input types.ReviewRequestPullRequest
		err = json.Unmarshal(body, &input)
		if err != nil {
			zapLog.Error("error unmarshal JSON",
				zap.Error(err),
			)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		svc := db.DynamoDbConnection()
		timeStamp, err := db.GetSlackTimeStamp(svc, input.PullRequest.ID, input.Number)
		if err != nil {
			zapLog.Error("error slack send message",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if timeStamp != "" {
			reviewers := []string{input.RequestedReviewer.Login}
			var slackMention string = "Please review: "
			for _, user := range reviewers {
				slackMention += fmt.Sprintf("<@%s>", slackUsersMap[user])
			}
			err = slack.SlackSendMessageThread(timeStamp, slackMention)
			if err != nil {
				zapLog.Error("error slack send message",
					zap.Error(err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

		}
	}

	// Directly commented in the PR issue
	if action == "created" {
		// parse request
		var input types.CommentPullRequest
		err = json.Unmarshal(body, &input)
		if err != nil {
			zapLog.Error("error unmarshal JSON",
				zap.Error(err),
			)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		svc := db.DynamoDbConnection()
		prId, err := github.GetPullRequestId(input.Repository.Name, input.Issue.Number)
		if err != nil {
			zapLog.Error("error get pull request id",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		timeStamp, err := db.GetSlackTimeStamp(svc, int(prId), input.Issue.Number)
		if err != nil {
			zapLog.Error("error slack send message",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if timeStamp != "" {
			message := fmt.Sprintf("<@%s> submitted an issue comment.", slackUsersMap[input.Comment.User.Login])
			message += fmt.Sprintf("```%s```\n", input.Comment.Body)
			message += fmt.Sprintf("View <%s|here>", input.Comment.HtmlUrl)
			err = slack.SlackSendMessageThread(timeStamp, message)
			if err != nil {
				zapLog.Error("error slack send message",
					zap.Error(err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}

	// closed / merged PR
	if action == "closed" {
		// parse request
		var input types.ClosedPullRequest
		err = json.Unmarshal(body, &input)
		if err != nil {
			zapLog.Error("error unmarshal JSON",
				zap.Error(err),
			)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		svc := db.DynamoDbConnection()
		timeStamp, err := db.GetSlackTimeStamp(svc, input.PullRequest.ID, input.Number)
		if err != nil {
			zapLog.Error("error slack send message",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if timeStamp != "" {
			var message string = fmt.Sprintf("<@%s> closed the pull request. ", slackUsersMap[input.PullRequest.User.Login])
			if len(input.PullRequest.MergedAt) > 0 {
				message = fmt.Sprintf("<@%s> merged the pull request. ", slackUsersMap[input.PullRequest.User.Login])
			}

			err := slack.SlackSendMessageThread(timeStamp, message)
			if err != nil {
				zapLog.Error("error slack send message",
					zap.Error(err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Delete PR in dynamodb Table
			svc := db.DynamoDbConnection()
			err = db.DeleteItem(svc, input.PullRequest.ID, input.Number)
			if err != nil {
				zapLog.Error("error delete data",
					zap.Error(err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}

	// submitted a PR review
	if action == "submitted" {
		// parse request
		var input types.SubmitReviewPullRequest
		err = json.Unmarshal(body, &input)
		if err != nil {
			zapLog.Error("error unmarshal JSON",
				zap.Error(err),
			)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		svc := db.DynamoDbConnection()
		timeStamp, err := db.GetSlackTimeStamp(svc, input.PullRequest.ID, input.PullRequest.Number)
		if err != nil {
			zapLog.Error("error slack send message",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if timeStamp != "" {
			if input.Review.State == "commented" {
				message := fmt.Sprintf("<@%s> submitted a review comment. ", slackUsersMap[input.PullRequest.User.Login])
				if len(input.Review.Body) > 0 {
					message += fmt.Sprintf("```%s```\n", input.Review.Body)
				}
				message += fmt.Sprintf("View <%s|here>", input.Review.HtmlUrl)
				err := slack.SlackSendMessageThread(timeStamp, message)
				if err != nil {
					zapLog.Error("error slack send message",
						zap.Error(err),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			if input.Review.State == "approved" {
				message := fmt.Sprintf("<@%s> approved the pull request. ", slackUsersMap[input.PullRequest.User.Login])
				if len(input.Review.Body) > 0 {
					message += fmt.Sprintf("```%s```\n", input.Review.Body)
				}
				message += fmt.Sprintf("View <%s|here>", input.Review.HtmlUrl)
				err := slack.SlackSendMessageThread(timeStamp, message)
				if err != nil {
					zapLog.Error("error slack send message",
						zap.Error(err),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// added commits to the PR branch
	if action == "synchronize" {
		// parse request
		var input types.PushPullRequestSync
		err = json.Unmarshal(body, &input)
		if err != nil {
			zapLog.Error("error unmarshal JSON",
				zap.Error(err),
			)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		svc := db.DynamoDbConnection()
		timeStamp, err := db.GetSlackTimeStamp(svc, input.PullRequest.ID, input.PullRequest.Number)
		if err != nil {
			zapLog.Error("error slack send message",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if timeStamp != "" {
			commitLink := fmt.Sprintf("%s/commits/%s", input.PullRequest.HtmlUrl, input.After)
			message := fmt.Sprintf("<@%s> committed a change. See it <%s|here>.", slackUsersMap[input.PullRequest.User.Login], commitLink)
			err = slack.SlackSendMessageThread(timeStamp, message)
			if err != nil {
				zapLog.Error("error slack send message",
					zap.Error(err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}

	// check run completed
	if action == "completed" {
		// parse request
		var input types.CheckRunPullRequest
		err = json.Unmarshal(body, &input)
		if err != nil {
			zapLog.Error("error unmarshal JSON",
				zap.Error(err),
			)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		svc := db.DynamoDbConnection()
		var pullRequestNumber int
		var pullRequestId int

		// should always only have one element
		for _, e := range input.CheckRun.PullRequests {
			pullRequestId = e.ID
			pullRequestNumber = e.Number
		}

		timeStamp, err := db.GetSlackTimeStamp(svc, pullRequestId, pullRequestNumber)
		if err != nil {
			zapLog.Error("error slack send message",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if timeStamp != "" {
			if input.CheckRun.Status == "completed" && len(input.CheckRun.CompletedAt) > 0 {
				message := fmt.Sprintf("Check run <%s|%s> has completed.", input.CheckRun.HtmlUrl, input.CheckRun.Name)
				err := slack.SlackSendMessageThread(timeStamp, message)
				if err != nil {
					zapLog.Error("error slack send message",
						zap.Error(err),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			if input.CheckRun.Status == "failed" {
				message := fmt.Sprintf("Check run <%s|%s> has failed.", input.CheckRun.HtmlUrl, input.CheckRun.Name)
				err := slack.SlackSendMessageThread(timeStamp, message)
				if err != nil {
					zapLog.Error("error slack send message",
						zap.Error(err),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// PR reopened
	if action == "reopened" {
		// parse request
		var input types.OpenPullRequest
		err = json.Unmarshal(body, &input)
		if err != nil {
			zapLog.Error("error unmarshal JSON",
				zap.Error(err),
			)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		messageText := fmt.Sprintf("<@%s> Reopened <%s|pull request> in `%s`.", slackUsersMap[input.PullRequest.User.Login], input.PullRequest.HtmlUrl, input.Repository.Name)

		timeStamp, err := slack.SlackSendMessage(input, messageText)
		if err != nil {
			zapLog.Error("error slack send message",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if len(input.PullRequest.RequestedReviewers) > 0 {
			reviewers := []string{}
			for _, reviewer := range input.PullRequest.RequestedReviewers {
				reviewers = append(reviewers, reviewer.Login)
			}

			var slackMention string = "Please review: "
			for _, user := range reviewers {
				slackMention += fmt.Sprintf("<@%s>", slackUsersMap[user])
			}
			err = slack.SlackSendMessageThread(timeStamp, slackMention)
			if err != nil {
				zapLog.Error("error slack send message",
					zap.Error(err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

		}

		svc := db.DynamoDbConnection()
		item := &types.TablePullRequestData{
			ID:             fmt.Sprintf("%d", input.PullRequest.ID),
			PullRequestId:  input.Number,
			SlackTimeStamp: timeStamp,
		}

		err = db.InsertItem(svc, item)
		if err != nil {
			zapLog.Error("error insert data",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	bodyBytes := Response{
		Message: "Webhook done.",
	}

	j, err := json.Marshal(bodyBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
