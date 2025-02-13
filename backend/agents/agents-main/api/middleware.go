package api

import (
	"bytes"

	"encoding/json"
	"io"

	"net/http"

	"github.com/segp/agents-main/context_keys"
	"github.com/segp/agents-main/storage"
)

type Middleware func(http.Handler) http.Handler

func CreateMiddlewareStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}
		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *wrappedWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		// If WriteHeader was never called, default to status 200
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(data)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// var jwtSecret = os.Getenv("SUPABASE_JWT_SECRET")

		// authHeader := r.Header.Get("Authorization")
		// if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, jwt.ErrSignatureInvalid
		// 	}
		// 	return []byte(jwtSecret), nil
		// })

		// if err != nil || !token.Valid {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// // Extract user ID from token claims
		// claims, ok := token.Claims.(jwt.MapClaims)
		// if !ok || !token.Valid {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// userID, ok := claims["sub"].(string)
		// log.Println("id: ", userID)
		// if !ok {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// // Add user ID to context
		// ctx := context.WithValue(r.Context(), "USER_ID", userID)
		// next.ServeHTTP(w, r.WithContext(ctx))

		userID := "2fc9c0d0-833d-462b-9cd5-670c5d075bf6"
		ctx := context_keys.SetUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TODO: add user id to the request table? (quick)
func requestIdMiddleware(next http.Handler, store storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the request path (if needed)
		path := r.URL.Path
		_ = path // currently unused but retained as in the original code

		var extractedData map[string]interface{}

		switch r.Method {
		case http.MethodPost:
			// Extract the request body as a map[string]interface{}
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}
			// Restore the request body for downstream handlers
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Only attempt to unmarshal if the body is non-empty
			if len(bodyBytes) > 0 {
				if err := json.Unmarshal(bodyBytes, &extractedData); err != nil {
					http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
					return
				}
			}

		case http.MethodGet:
			// Extract query parameters as a map[string]interface{}
			extractedData = make(map[string]interface{})
			for key, values := range r.URL.Query() {
				if len(values) == 1 {
					extractedData[key] = values[0]
				} else {
					extractedData[key] = values
				}
			}
		default:
			// do nothign
		}

		createdReq, err := storage.Store(store, storage.NewAgentRequest(path, extractedData, ""))
		if err != nil {
			http.Error(w, "Error storing request", http.StatusInternalServerError)
			return
		}

		ctx := context_keys.SetRequestID(r.Context(), createdReq.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func sessionIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request has a body to read
		if r.Body != nil {
			// Read the request body
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}
			// Restore the request body so that downstream handlers can access it
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Only attempt to unmarshal if the body is non-empty
			if len(bodyBytes) > 0 {
				var payload map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &payload); err == nil {
					if sessionID, ok := payload["session_id"].(string); ok {
						// If session_id is found, add it to the context
						r = r.WithContext(context_keys.SetSessionID(r.Context(), sessionID))
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
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
