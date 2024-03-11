package constants

import (
	"reflect"
	"testing"
)

func TestPort(t *testing.T) {
	expected := &Ports{
		MainApi: 8080,
	}

	result := Port()

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("struct does not match the expected.")
	}
}
