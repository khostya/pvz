package grpcserver

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	ctx  context.Context
	port int
	grpc *grpc.Server

	err     chan error
	started atomic.Bool
}

func New(ctx context.Context, port int, interceptor ...grpc.ServerOption) *Server {
	grpc := grpc.NewServer(interceptor...)

	server := &Server{
		grpc: grpc,
		port: port,
		ctx:  ctx,
	}

	return server
}

func (s *Server) Wait() <-chan error {
	return s.err
}

func (s *Server) Start() error {
	if s.started.Swap(true) {
		return ErrAlreadyStarted
	}

	lis, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(s.port)))
	if err != nil {
		return err
	}

	errorChan := make(chan error)
	go func() {
		defer close(errorChan)
		err := s.grpc.Serve(lis)
		if err != nil {
			errorChan <- err
		}
	}()
	s.err = errorChan

	go func(ctx context.Context) {
		defer lis.Close()
		defer s.grpc.GracefulStop()

		select {
		case <-ctx.Done():
			return
		}
	}(s.ctx)

	return nil
}

func (s *Server) GetRegistrar() grpc.ServiceRegistrar {
	return s.grpc
}
