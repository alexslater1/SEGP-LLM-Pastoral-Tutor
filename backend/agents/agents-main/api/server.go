package api

import (
	"net/http"
	"os"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/api/handlers"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/storage"
)

type Server struct {
	listenAddr string
	router     *http.ServeMux

	storage storage.Storage
}

func NewServer(listenAddr string, storage storage.Storage) *Server {
	s := &Server{
		listenAddr: listenAddr,
		router:     http.NewServeMux(),
		storage:    storage,
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	agent := agent.NewDefaultEventStoringLoggingFastAgent()
	store := storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
	history := history.NewStoreHistory(store)

	s.router.HandleFunc("POST /completion/v2", handlers.ChatCompletionV2(agent, store, history))
	s.router.HandleFunc("GET /completion/v2/status/{request_id}", handlers.ChatCompletionV2Status(store))

	s.router.HandleFunc("GET /chats/{chat_id}", handlers.ChatHistory(store))
}

func (s *Server) Start() error {
	requestIdMiddlewareClosure := func(http.Handler) http.Handler {
		return requestIdMiddleware(s.router, s.storage)
	}

	stack := CreateMiddlewareStack(
		corsMiddleware, // CORS middleware should be first
		requestIdMiddlewareClosure,
		// Auth,
	)

	return http.ListenAndServe(s.listenAddr, stack(s.router))
}
