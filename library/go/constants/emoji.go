package constants

// github user name and slack user id
type Emojis struct {
	CheckPassed      string
	CheckFailed      string
	Approved         string
	Closed           string
	Merged           string
	Opened           string
	PullRequest      string
	Pushed           string
	RequestedChanges string
	Reviewed         string
	RequestReview    string
	Comment          string
}

func Emoji() *Emojis {
	return &Emojis{
		CheckPassed:      ":check-passed:",
		CheckFailed:      ":check-failed:",
		Approved:         ":approved:",
		Closed:           ":closed:",
		Merged:           ":merged:",
		Opened:           ":opened:",
		PullRequest:      ":pull-request:",
		Pushed:           ":pushed:",
		RequestedChanges: ":requested-changes:",
		Reviewed:         ":reviewed:",
		RequestReview:    ":eyes:",
		Comment:          ":writing_hand:",
	}
}
