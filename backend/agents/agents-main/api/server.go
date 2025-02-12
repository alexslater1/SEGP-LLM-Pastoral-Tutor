package api

import (
	"net/http"
	"os"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/api/handlers"
	"github.com/segp/agents-main/storage"
)

type Server struct {
	listenAddr string
	router     *http.ServeMux
}

func NewServer(listenAddr string) *Server {
	s := &Server{
		listenAddr: listenAddr,
		router:     http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")                            // Frontend URL
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")   // Allowed methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Include Authorization header

		if r.Method == http.MethodOptions {
			// Respond to preflight requests
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) routes() {
	agent := agent.NewDefaultEventStoringLoggingFastAgent()

	s.router.HandleFunc("POST /completion", handlers.ChatCompletion(agent, storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))))
	s.router.HandleFunc("POST /completion/v2", handlers.ChatCompletionV2(agent, storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))))
	s.router.HandleFunc("GET /completion/v2/status/{request_id}", handlers.ChatCompletionV2Status(storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))))
	s.router.HandleFunc("POST /completion/v3", handlers.ChatCompletionV3(agent, storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))))
}

func (s *Server) Start() error {
	stack := CreateMiddlewareStack(
		s.corsMiddleware, // CORS middleware should be first
		// Auth,
	)

	return http.ListenAndServe(s.listenAddr, stack(s.router))
}
