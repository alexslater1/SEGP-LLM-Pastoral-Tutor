package api

import (
	"log"
	"net/http"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/api/handlers"
	"github.com/segp/agents-main/crew"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/jobs"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
)

type Server struct {
	listenAddr string
	router     *http.ServeMux

	storage    storage.Storage
	agent      agent.Agent
	crew       *crew.Crew
	history    history.History
	jobManager *jobs.JobManager
	llm        llm.LLM
}

func NewServer(listenAddr string, storage storage.Storage, agent agent.Agent, crew *crew.Crew, history history.History, jobManager *jobs.JobManager, llm llm.LLM) *Server {
	s := &Server{
		listenAddr: listenAddr,
		router:     http.NewServeMux(),
		storage:    storage,
		agent:      agent,
		crew:       crew,
		history:    history,
		jobManager: jobManager,
		llm:        llm,
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("POST /completion/v2", handlers.ChatCompletionV2(s.crew, s.agent, s.storage, s.history, s.llm))
	s.router.HandleFunc("GET /completion/v2/status/{request_id}", handlers.ChatCompletionV2Status(s.storage))

	s.router.HandleFunc("GET /sessions/{session_id}/history", handlers.ChatHistory(s.history))
	s.router.HandleFunc("GET /sessions", handlers.SessionIdsForUser(s.storage))
	s.router.HandleFunc("GET /sessions/{session_id}", handlers.SessionFromId(s.storage))
}

func (s *Server) Start() error {
	requestIdMiddlewareClosure := func(next http.Handler) http.Handler {
		return requestIdMiddleware(next, s.storage)
	}

	stack := CreateMiddlewareStack(
		corsMiddleware, // CORS middleware should be first
		authMiddleware,
		requestIdMiddlewareClosure,
		sessionIDMiddleware,
	)

	go func() {
		errCh := s.jobManager.Start()
		for err := range errCh {
			log.Printf("Job manager error: %v", err)
		}
	}()

	return http.ListenAndServe(s.listenAddr, stack(s.router))
}
