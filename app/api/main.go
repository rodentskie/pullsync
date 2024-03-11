package main

import (
	"fmt"
	"net/http"
	"slack-pr-lambda/api/handlers"
	"slack-pr-lambda/constants"
	"slack-pr-lambda/env"
	"slack-pr-lambda/logger"

	"go.uber.org/zap"
)

func main() {
	l := logger.LoggerConfig()
	log, _ := l.Build()
	ports := constants.Port()

	portString := fmt.Sprintf(":%d", ports.MainApi)

	port := env.GetEnv("PORT", portString)
	host := env.GetEnv("HOST", "localhost")

	mux := http.NewServeMux()

	mux.HandleFunc("POST /items/{id}", handlers.PullRequestHandler)

	sandboxLink := fmt.Sprintf("http://%s%s", host, port)
	log.Info("running at 🚀⚙️",
		zap.String("link", sandboxLink),
	)

	if err := http.ListenAndServe(port, mux); err != nil && err != http.ErrServerClosed {
		log.Fatal("error serve api",
			zap.String("port", port),
			zap.Error(err),
		)
	}
}
