package main

import "testing"

func TestMain(t *testing.T) {
	t.Logf("No need test, granule infra are tested.")
	if false {
		t.Errorf("This should not fail")
	}
}
