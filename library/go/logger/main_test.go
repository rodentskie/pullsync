package logger

import (
	"os"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

type TestEnvData struct {
	key      string
	fallback string
}

func TestLoggerConfig(t *testing.T) {
	data := []TestEnvData{
		{"ENV", "dev"},
		{"ENV", "test"},
	}

	for _, e := range data {
		os.Setenv(e.key, e.fallback)

		config := LoggerConfig()

		expectedType := reflect.TypeOf(zap.Config{})

		if reflect.TypeOf(config) != expectedType {
			t.Errorf("FAIL: Unexpected config type. Expected: %s, Got: %s", expectedType, reflect.TypeOf(config))
		}

		if os.Getenv("ENV") == "dev" && config.Level.Level().CapitalString() != "INFO" {
			t.Errorf("FAIL: Unexpected config level. Expected: %s, Got: %s", "info", config.Level)
		}

		if os.Getenv("ENV") == "test" && config.Level.Level().CapitalString() != "FATAL" {
			t.Errorf("FAIL: Unexpected config type. Expected: %s, Got: %s", "fatal", config.Level)
		}
	}

}
