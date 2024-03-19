package types

type TablePullRequestData struct {
	ID             string `json:"id"`
	PullRequestId  int    `json:"pullRequestId"`
	SlackTimeStamp string `json:"slackTimeStamp"`
}

type OpenPullRequest struct {
	Action      string                `json:"action"`
	Number      int                   `json:"number"`
	PullRequest pullRequest           `json:"pull_request"`
	Repository  pullRequestRepository `json:"repository"`
	Sender      sender                `json:"sender"`
}

type ReviewRequestPullRequest struct {
	Action            string               `json:"action"`
	Number            int                  `json:"number"`
	PullRequest       pullRequest          `json:"pull_request"`
	RequestedReviewer pullRequestReviewers `json:"requested_reviewer"`
}

type CommentPullRequest struct {
	Action     string                `json:"action"`
	Issue      issue                 `json:"issue"`
	Comment    comment               `json:"comment"`
	Repository pullRequestRepository `json:"repository"`
}

type ClosedPullRequest struct {
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	PullRequest pullRequest `json:"pull_request"`
	Sender      sender      `json:"sender"`
}

type SubmitReviewPullRequest struct {
	Action      string      `json:"action"`
	PullRequest pullRequest `json:"pull_request"`
	Review      review      `json:"review"`
}

type PushPullRequestSync struct {
	Action      string                `json:"action"`
	Number      int                   `json:"number"`
	PullRequest pullRequest           `json:"pull_request"`
	Repository  pullRequestRepository `json:"repository"`
	After       string                `json:"after"`
	Sender      sender                `json:"sender"`
}

type CheckRunPullRequest struct {
	Action     string                `json:"action"`
	Repository pullRequestRepository `json:"repository"`
	Sender     sender                `json:"sender"`
	CheckRun   checkRun              `json:"check_run"`
}
