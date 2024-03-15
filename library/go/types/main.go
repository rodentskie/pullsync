package types

type TablePullRequestData struct {
	ID             string `json:"id"`
	PullRequestId  int32  `json:"pullRequestId"`
	SlackTimeStamp string `json:"slackTimeStamp"`
}

type OpenPullRequest struct {
	Action      string                `json:"action"`
	Number      int32                 `json:"number"`
	PullRequest pullRequest           `json:"pull_request"`
	Repository  pullRequestRepository `json:"repository"`
}

type ReviewRequestPullRequest struct {
	Action            string               `json:"action"`
	Number            int32                `json:"number"`
	PullRequest       pullRequest          `json:"pull_request"`
	RequestedReviewer pullRequestReviewers `json:"requested_reviewer"`
}

type CommentPullRequest struct {
	Action  string  `json:"action"`
	Issue   issue   `json:"issue"`
	Comment comment `json:"comment"`
}

type ClosedPullRequest struct {
	Action      string      `json:"action"`
	Number      int32       `json:"number"`
	PullRequest pullRequest `json:"pull_request"`
}
