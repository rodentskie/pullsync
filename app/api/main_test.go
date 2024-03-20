package main

import "testing"

func TestMain(t *testing.T) {
	t.Logf("Running minimal example test")
	if false {
		t.Errorf("This should not fail")
	}
}
