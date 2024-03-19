package constants

import (
	"reflect"
	"testing"
)

func TestEmoji(t *testing.T) {
	expected := &Emojis{
		CheckPassed:      ":check-passed:",
		CheckFailed:      ":check-failed:",
		CheckCanceled:    ":octagonal_sign:",
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

	result := Emoji()

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("struct does not match the expected.")
	}
}
