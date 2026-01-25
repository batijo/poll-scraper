package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/batijo/poll-scraper/api/handlers"
)

type Server struct {
	*http.Server
	mux *http.ServeMux
}

func New() *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Data)
	handler := withMiddleware(mux)
	return &Server{
		Server: &http.Server{
			Handler: handler,
		},
		mux: mux,
	}
}

func withMiddleware(h http.Handler) http.Handler {
	origins := "*"
	if os.Getenv("PS_DOMAINS") != "" {
		origins = os.Getenv("PS_DOMAINS")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "poll-scraper")
		w.Header().Set("Access-Control-Allow-Origin", origins)
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"message": message})
}
