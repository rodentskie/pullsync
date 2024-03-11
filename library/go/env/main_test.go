package env

import (
	"os"
	"testing"
)

type TestEnvData struct {
	key      string
	fallback string
}

func TestEnv(t *testing.T) {
	data := []TestEnvData{
		{"ONE_TEST", "onetest"},
		{"TWO_TEST", "twotest"},
		{"THREE_TEST", "threetest"},
		{"FOUR_TEST", "fourtest"},
		{"FIVE_TEST", "fivetest"},
	}

	for index, e := range data {
		if mod := index % 2; mod == 0 {
			os.Setenv(e.key, e.fallback)
		}
		result := GetEnv(e.key, e.fallback)

		if result != e.fallback {
			t.Errorf("FAIL: Expected env variable. Expected: %s, Got: %s\n", result, e.fallback)
		}
	}

}
