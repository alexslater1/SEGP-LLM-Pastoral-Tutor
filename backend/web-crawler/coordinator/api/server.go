package api

import "net/http"

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

	s.router.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
}

func (s *Server) Start() error {
	stack := CreateMiddlewareStack(
		s.corsMiddleware, // CORS middleware should be first
		// Auth,
	)

	return http.ListenAndServe(s.listenAddr, stack(s.router))
}