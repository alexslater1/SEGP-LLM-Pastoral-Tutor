package handlers

import (
	"net/http"

	"github.com/ethanhosier/web-crawler-coordinator/coordinator_client"
)

type ScraperWorkerParams struct {
	URL string `json:"url"`
}

type ScrapeRagCreationRequest struct {
	URLs []string `json:"urls"`
}

func ScrapeRagCreation(coordinatorClient *coordinator_client.CoordinatorClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ScrapeRagCreationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
