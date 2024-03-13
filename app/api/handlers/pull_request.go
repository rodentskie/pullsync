package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"slack-pr-lambda/logger"
	"syscall"

	"go.uber.org/zap"
)

func PullRequestHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.LoggerConfig()
	zapLog, _ := l.Build()

	defer func() {
		err := r.Body.Close()
		if err != nil {
			zapLog.Fatal("error close req body",
				zap.Error(err),
			)
		}

		if err := zapLog.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			log.Fatalf("error closing the logger. %v\n", err)
		}
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zapLog.Fatal("error read request body",
			zap.Error(err),
		)
	}

	zapLog.Info("PR body",
		zap.String("req", string(body)),
	)

	bodyBytes := Response{
		Message: "Just a response.",
	}

	j, err := json.Marshal(bodyBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
