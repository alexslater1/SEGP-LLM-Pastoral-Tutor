package api

import (
	"log"
	"net/http"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/api/handlers"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/jobs"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
)

type Server struct {
	listenAddr string
	router     *http.ServeMux

	storage       storage.Storage
	history       history.History
	agent         agent.Agent
	jobManager    *jobs.JobManager
	llm           llm.LLM
	agentProvider *agent.AgentProvider
}

func NewServer(listenAddr string, storage storage.Storage, agent agent.Agent, history history.History, jobManager *jobs.JobManager, llm llm.LLM, agentProvider *agent.AgentProvider) *Server {
	s := &Server{
		listenAddr:    listenAddr,
		router:        http.NewServeMux(),
		storage:       storage,
		agent:         agent,
		history:       history,
		jobManager:    jobManager,
		llm:           llm,
		agentProvider: agentProvider,
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("POST /completion/v2", handlers.ChatCompletionV2(s.agentProvider, s.storage, s.history, s.llm))
	s.router.HandleFunc("GET /completion/v2/status/{request_id}", handlers.ChatCompletionV2Status(s.storage))

	s.router.HandleFunc("GET /sessions/{session_id}/history", handlers.ChatHistory(s.history))
	s.router.HandleFunc("GET /sessions", handlers.SessionIdsForUser(s.storage))
	s.router.HandleFunc("GET /sessions/{session_id}", handlers.SessionFromId(s.storage))
 	s.router.HandleFunc("DELETE /sessions/{session_id}", handlers.SetDeletedSessionFromId(s.storage))

	s.router.HandleFunc("GET /agents", handlers.GetAllAgents(s.storage))
	s.router.HandleFunc("GET /agents/config", handlers.GetConfigOptions())
	s.router.HandleFunc("POST /agents", handlers.SetConfigs(s.agentProvider))
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
