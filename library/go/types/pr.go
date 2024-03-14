package types

type pullRequestUser struct {
	Login string `json:"login"`
	ID    int32  `json:"id"`
}

type pullRequestReviewers struct {
	Login string `json:"login"`
	ID    int32  `json:"id"`
}

type pullRequest struct {
	ID                 int32                  `json:"id"`
	HtmlUrl            string                 `json:"html_url"`
	State              string                 `json:"state"`
	Locked             bool                   `json:"locked"`
	Title              string                 `json:"title"`
	User               pullRequestUser        `json:"user"`
	RequestedReviewers []pullRequestReviewers `json:"requested_reviewers"`
}

type pullRequestRepository struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	HtmlUrl string `json:"html_url"`
}
