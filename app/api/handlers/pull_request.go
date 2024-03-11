package handlers

import (
	"fmt"
	"net/http"
)

func PullRequestHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Printf("item id is %s\n", id)
}
