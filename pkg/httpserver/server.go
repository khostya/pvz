package httpserver

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	_defaultAddr         = ":80"
	_defaultReadTimeout  = 5 * time.Second
	_defaultWriteTimeout = 5 * time.Second
	_defaultIdleTimeout  = 5 * time.Second
)

// Server -.
type Server struct {
	server *http.Server

	notify       chan error
	address      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

func New(h *echo.Echo, opts ...Option) *Server {
	s := &Server{
		notify:       make(chan error, 1),
		address:      _defaultAddr,
		readTimeout:  _defaultReadTimeout,
		writeTimeout: _defaultWriteTimeout,
		idleTimeout:  _defaultIdleTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      h,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
	}

	return s
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
