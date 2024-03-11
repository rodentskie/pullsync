package handlers

import (
	"errors"
	"fmt"
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

	fmt.Println(string(body))
}
