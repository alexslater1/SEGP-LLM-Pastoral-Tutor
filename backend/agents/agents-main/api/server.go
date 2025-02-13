package api

import (
	"net/http"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/api/handlers"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/storage"
)

type Server struct {
	listenAddr string
	router     *http.ServeMux

	storage storage.Storage
	agent   agent.Agent
	history history.History
}

func NewServer(listenAddr string, storage storage.Storage, agent agent.Agent, history history.History) *Server {
	s := &Server{
		listenAddr: listenAddr,
		router:     http.NewServeMux(),
		storage:    storage,
		agent:      agent,
		history:    history,
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("POST /completion/v2", handlers.ChatCompletionV2(s.agent, s.storage, s.history))
	s.router.HandleFunc("GET /completion/v2/status/{request_id}", handlers.ChatCompletionV2Status(s.storage))
	s.router.HandleFunc("GET /chats/{chat_id}", handlers.ChatHistory(s.history))
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

	return http.ListenAndServe(s.listenAddr, stack(s.router))
}
