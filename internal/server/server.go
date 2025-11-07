package server

import "net/http"

type Server struct{}

func New() *Server {
    return &Server{}
}

func (s *Server) Router() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })
    mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ready"))
    })
    return mux
}
