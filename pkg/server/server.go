package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	notify chan error
}

func New(handler http.Handler, host, port string) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         net.JoinHostPort(host, port),
	}

	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
	}

	go s.start()

	return s
}

func (s *Server) start() {
	s.notify <- s.server.ListenAndServe()
	close(s.notify)
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Fatal("server - Close - s.server.Shutdown")
	}
}
