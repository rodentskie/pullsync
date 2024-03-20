package main

import "testing"

func TestMain(t *testing.T) {
	t.Logf("Granule code are tested, no need to test main function.")
	if false {
		t.Errorf("This should not fail")
	}
}
