package types

type pullRequestUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

type pullRequestReviewers struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

type pullRequest struct {
	ID                 int                    `json:"id"`
	Number             int                    `json:"number"`
	HtmlUrl            string                 `json:"html_url"`
	State              string                 `json:"state"`
	Locked             bool                   `json:"locked"`
	Title              string                 `json:"title"`
	User               pullRequestUser        `json:"user"`
	RequestedReviewers []pullRequestReviewers `json:"requested_reviewers"`
	MergedAt           string                 `json:"merged_at"`
}

type sender struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

type pullRequestRepository struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	HtmlUrl string `json:"html_url"`
}

type issue struct {
	ID      int    `json:"id"`
	Number  int    `json:"number"`
	HtmlUrl string `json:"html_url"`
}

type comment struct {
	ID      int             `json:"id"`
	HtmlUrl string          `json:"html_url"`
	Body    string          `json:"body"`
	User    pullRequestUser `json:"user"`
}

type review struct {
	ID      int             `json:"id"`
	HtmlUrl string          `json:"html_url"`
	User    pullRequestUser `json:"user"`
	Body    string          `json:"body"`
	State   string          `json:"state"`
}
