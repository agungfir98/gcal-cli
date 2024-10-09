package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	*http.Server
	codeCh chan string
}

func NewServer(listenAddr string, codeCh chan string) *Server {

	s := &Server{
		Server: &http.Server{Addr: listenAddr},
	}
	s.codeCh = codeCh

	http.HandleFunc("/callback", s.handleCallback)

	return s

}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code != "" {
		s.codeCh <- code
	}
	fmt.Fprintln(w, "Auth successful, you may close this page.")
}

func (s *Server) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
		}
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Server.Shutdown(ctx)
}
