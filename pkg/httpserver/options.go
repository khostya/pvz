package httpserver

import (
	"net"
	"strconv"
	"time"
)

type Option func(*Server)

func Port(port uint16) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", strconv.Itoa(int(port)))
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

func IdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.idleTimeout = timeout
	}
}
