package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethanhosier/web-crawler-coordinator/coordinator_client"
	"github.com/ethanhosier/web-crawler-coordinator/utils"
	"github.com/google/uuid"
)

const (
	maxUrls = 500
)

type ScraperWorkerParams struct {
	URL string `json:"url"`
}

type CreateScrapeRagTaskRequest struct {
	URLs []string `json:"urls"`
}

type CreatedTask struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Error string `json:"error,omitempty"`
}

type CreateScrapeRagTaskResponse struct {
	CreatedTasks []CreatedTask `json:"created_tasks"`
}

func ScrapeRagTask(coordinatorClient coordinator_client.CoordinatorClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateScrapeRagTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(req.URLs) == 0 {
			WriteJSONError(w, "No URLs provided", http.StatusBadRequest)
			return
		}

		if len(req.URLs) > maxUrls {
			WriteJSONError(w, fmt.Sprintf("Maximum number of URLs is %d", maxUrls), http.StatusBadRequest)
			return
		}

		tasks := make([]*coordinator_client.Task, 0, len(req.URLs))
		createdTasks := make([]CreatedTask, 0, len(req.URLs))

		for _, url := range req.URLs {
			task, createdTask := processURL(url)
			if task != nil {
				tasks = append(tasks, task)
			}

			createdTasks = append(createdTasks, createdTask)
		}

		err := coordinatorClient.CreateTasks(r.Context(), coordinator_client.CoordinatorClientTaskTopicUrls, tasks)
		if err != nil {
			WriteJSONError(w, "Failed to create tasks", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, CreateScrapeRagTaskResponse{
			CreatedTasks: createdTasks,
		})
	}
}

type TasksStatusResponse struct {
	NumUrlsTasks          int `json:"num_urls_tasks"`
	NumProcessingUrlTasks int `json:"num_processing_url_tasks"`
	NumRagTasks           int `json:"num_rag_tasks"`
	NumProcessingRagTasks int `json:"num_processing_rag_tasks"`
}

func TasksStatus(coordinatorClient coordinator_client.CoordinatorClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		numUrlsTasks, err := coordinatorClient.NumTasks(r.Context(), coordinator_client.CoordinatorClientTaskTopicUrls)
		if err != nil {
			WriteJSONError(w, "Failed to get number of tasks", http.StatusInternalServerError)
			return
		}

		numProcessingUrlTasks, err := coordinatorClient.NumProcessingTasks(r.Context(), coordinator_client.CoordinatorClientTaskTopicUrls)
		if err != nil {
			WriteJSONError(w, "Failed to get number of processing tasks", http.StatusInternalServerError)
			return
		}

		numRagTasks, err := coordinatorClient.NumTasks(r.Context(), coordinator_client.CoordinatorClientTaskTopicRag)
		if err != nil {
			WriteJSONError(w, "Failed to get number of tasks", http.StatusInternalServerError)
			return
		}

		numProcessingRagTasks, err := coordinatorClient.NumProcessingTasks(r.Context(), coordinator_client.CoordinatorClientTaskTopicRag)
		if err != nil {
			WriteJSONError(w, "Failed to get number of processing tasks", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, TasksStatusResponse{
			NumUrlsTasks:          numUrlsTasks,
			NumProcessingUrlTasks: numProcessingUrlTasks,
			NumRagTasks:           numRagTasks,
			NumProcessingRagTasks: numProcessingRagTasks,
		})
	}
}

func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func WriteJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func processURL(url string) (*coordinator_client.Task, CreatedTask) {
	formattedUrl, err := utils.FormatUrl(url)
	if err != nil {
		return nil, CreatedTask{
			ID:    "",
			URL:   url,
			Error: err.Error(),
		}
	}

	params := ScraperWorkerParams{
		URL: formattedUrl,
	}

	task, err := coordinator_client.NewTask(uuid.New().String(), "coordinator-client", params)
	if err != nil {
		return nil, CreatedTask{
			ID:    "",
			URL:   url,
			Error: err.Error(),
		}
	}

	return task, CreatedTask{
		ID:    task.ID,
		URL:   url,
		Error: "",
	}
}
