package api

import (
	"context"
	"net/http"
	"os"

	"github.com/ethanhosier/web-crawler-coordinator/api/handlers"
	"github.com/ethanhosier/web-crawler-coordinator/coordinator_client"
	"github.com/ethanhosier/web-crawler-coordinator/utils"
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
	var (
		redisAddress  = utils.Required(os.Getenv("REDIS_ADDRESS"), "REDIS_ADDRESS")
		redisDB       = utils.RequiredInt(os.Getenv("REDIS_DB"), "REDIS_DB")
		redisPassword = os.Getenv("REDIS_PASSWORD")
	)

	coordinatorClient := coordinator_client.NewRedisCoordinatorClient(context.Background(), redisAddress, redisPassword, redisDB)

	s.router.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	s.router.HandleFunc("POST /scrape-rag-task", handlers.ScrapeRagTask(coordinatorClient))
	s.router.HandleFunc("GET /tasks-status", handlers.TasksStatus(coordinatorClient))
}

func (s *Server) Start() error {
	stack := CreateMiddlewareStack(
		s.corsMiddleware, // CORS middleware should be first
		// Auth,
	)

	return http.ListenAndServe(s.listenAddr, stack(s.router))
}
