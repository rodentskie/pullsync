package github

import (
	"context"
	"slack-pr-lambda/env"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v39/github"
)

func GetPullRequestId(repo string, prNumber int) (int64, error) {
	token := env.GetEnv("GITHUB_TOKEN", "token")
	owner := env.GetEnv("GITHUB_OWNER", "owner")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return 0, err
	}

	return pr.GetID(), nil
}
