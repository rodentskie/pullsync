package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slack-pr-lambda/logger"
	"syscall"
)

type Response struct {
	Message string `json:"message"`
}

func IndexRequestHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.LoggerConfig()
	zapLog, _ := l.Build()

	defer func() {
		if err := zapLog.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			log.Fatalf("error closing the logger. %v\n", err)
		}
	}()

	bodyBytes := Response{
		Message: "Welcome to Slack PR Lamba demo.",
	}

	j, err := json.Marshal(bodyBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
