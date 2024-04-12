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
	"strings"
	"syscall"

	"go.uber.org/zap"
)

func PullRequestHandler(w http.ResponseWriter, r *http.Request) {
	env := env.GetEnv("ENV", "local")

	l := logger.LoggerConfig()
	zapLog, _ := l.Build()

	slackUsers := constants.SlackUsers()
	slackUsersMap := mapstruct.StructToMap(*slackUsers)

	emoji := constants.Emoji()

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

		zapLog.Error("error read request body",
			zap.Error(err),
		)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
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
		zapLog.Error("error parse action from req body",
			zap.Error(err),
		)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
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

		user := slackUsersMap[input.Sender.Login]
		if input.Sender.Login == "dependabot[bot]" {
			user = "dependabot[bot]"
		}

		messageText := fmt.Sprintf("<@%s> %s opened new <%s|pull request> in `%s`.", user, emoji.Opened, input.PullRequest.HtmlUrl, input.Repository.Name)
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
				slackMention += fmt.Sprintf("<@%s> %s", slackUsersMap[user], emoji.RequestReview)
			}
			if err = slack.SlackSendMessageThread(timeStamp, slackMention); err != nil {
				zapLog.Error("error slack send message",
					zap.Error(err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		if err := slack.SlackAddReaction(timeStamp, strings.ReplaceAll(emoji.Opened, ":", "")); err != nil {
			zapLog.Error("error slack add reaction",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
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
				slackMention += fmt.Sprintf("<@%s> %s", slackUsersMap[user], emoji.RequestReview)
			}
			if err = slack.SlackSendMessageThread(timeStamp, slackMention); err != nil {
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
			message := fmt.Sprintf("<@%s> %s submitted an issue <%s|comment>. \n", slackUsersMap[input.Comment.User.Login], emoji.Comment, input.Comment.HtmlUrl)
			message += fmt.Sprintf("```%s```\n", input.Comment.Body)
			if err = slack.SlackSendMessageThread(timeStamp, message); err != nil {
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
			closeEmoji := emoji.Closed
			message := fmt.Sprintf("<@%s> closed the pull request %s. ", slackUsersMap[input.Sender.Login], emoji.Closed)
			if len(input.PullRequest.MergedAt) > 0 {
				closeEmoji = emoji.Merged
				message = fmt.Sprintf("<@%s> merged the pull request %s. ", slackUsersMap[input.Sender.Login], emoji.Merged)
			}

			if err := slack.SlackAddReaction(timeStamp, strings.ReplaceAll(closeEmoji, ":", "")); err != nil {
				zapLog.Error("error slack add reaction",
					zap.Error(err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if err := slack.SlackSendMessageThread(timeStamp, message); err != nil {
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
				message := fmt.Sprintf("<@%s> submitted a review <%s|comment> %s. \n ", slackUsersMap[input.Review.User.Login], input.Review.HtmlUrl, emoji.Reviewed)
				if len(input.Review.Body) > 0 {
					message += fmt.Sprintf("```%s```\n", input.Review.Body)
				}
				if err := slack.SlackSendMessageThread(timeStamp, message); err != nil {
					zapLog.Error("error slack send message",
						zap.Error(err),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			if input.Review.State == "approved" {
				message := fmt.Sprintf("<@%s> approved the pull <%s|request> %s. \n", slackUsersMap[input.Review.User.Login], input.Review.HtmlUrl, emoji.Approved)
				if len(input.Review.Body) > 0 {
					message += fmt.Sprintf("```%s```\n", input.Review.Body)
				}

				if err := slack.SlackAddReaction(timeStamp, strings.ReplaceAll(emoji.Approved, ":", "")); err != nil {
					zapLog.Error("error slack add reaction",
						zap.Error(err),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				if err := slack.SlackSendMessageThread(timeStamp, message); err != nil {
					zapLog.Error("error slack send message",
						zap.Error(err),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			if input.Review.State == "changes_requested" {
				message := fmt.Sprintf("<@%s> requested a change <%s|comment> %s. \n ", slackUsersMap[input.Review.User.Login], input.Review.HtmlUrl, emoji.RequestedChanges)
				if len(input.Review.Body) > 0 {
					message += fmt.Sprintf("```%s```\n", input.Review.Body)
				}
				if err := slack.SlackSendMessageThread(timeStamp, message); err != nil {
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
			message := fmt.Sprintf("<@%s> %s pushed a <%s|change>.", slackUsersMap[input.Sender.Login], emoji.Pushed, commitLink)
			if err = slack.SlackSendMessageThread(timeStamp, message); err != nil {
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
				message := fmt.Sprintf("Check run <%s|%s> %s.", input.CheckRun.HtmlUrl, input.CheckRun.Name, emoji.CheckPassed)
				if input.CheckRun.Conclusion == "failure" {
					message = fmt.Sprintf("Check run <%s|%s> %s.", input.CheckRun.HtmlUrl, input.CheckRun.Name, emoji.CheckFailed)
				}

				if input.CheckRun.Conclusion == "cancelled" {
					message = fmt.Sprintf("Check run <%s|%s> %s.", input.CheckRun.HtmlUrl, input.CheckRun.Name, emoji.CheckCanceled)
				}

				if err := slack.SlackSendMessageThread(timeStamp, message); err != nil {
					zapLog.Error("error slack send message",
						zap.Error(err),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			if input.CheckRun.CheckSuite.Status == "completed" && input.CheckRun.CheckSuite.Conclusion == "success" {
				message := fmt.Sprintf("All checks have passed. %s", emoji.CheckPassed)
				if err := slack.SlackSendMessageThread(timeStamp, message); err != nil {
					zapLog.Error("error slack send message",
						zap.Error(err),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			if input.CheckRun.CheckSuite.Status == "completed" && input.CheckRun.CheckSuite.Conclusion == "failure" {
				message := fmt.Sprintf("Some checks were not successful. %s", emoji.CheckFailed)
				if err := slack.SlackSendMessageThread(timeStamp, message); err != nil {
					zapLog.Error("error slack send message",
						zap.Error(err),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			if input.CheckRun.CheckSuite.Status == "completed" && input.CheckRun.CheckSuite.Conclusion == "cancelled" {
				message := fmt.Sprintf("Some checks were cancelled. %s", emoji.CheckCanceled)
				if err := slack.SlackSendMessageThread(timeStamp, message); err != nil {
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

		messageText := fmt.Sprintf("<@%s> %s Reopened <%s|pull request> in `%s`.", slackUsersMap[input.Sender.Login], emoji.Opened, input.PullRequest.HtmlUrl, input.Repository.Name)

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
				slackMention += fmt.Sprintf("<@%s> %s", slackUsersMap[user], emoji.RequestReview)
			}
			if err = slack.SlackSendMessageThread(timeStamp, slackMention); err != nil {
				zapLog.Error("error slack send message",
					zap.Error(err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

		}

		if err := slack.SlackAddReaction(timeStamp, strings.ReplaceAll(emoji.Opened, ":", "")); err != nil {
			zapLog.Error("error slack add reaction",
				zap.Error(err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
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
