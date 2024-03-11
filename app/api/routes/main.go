package routes

import (
	"net/http"
	"slack-pr-lambda/api/handlers"
)

func MainRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handlers.IndexRequestHandler)
	mux.HandleFunc("POST /pull-request", handlers.PullRequestHandler)
}
