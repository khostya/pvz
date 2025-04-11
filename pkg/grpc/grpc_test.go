package grpcserver

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()

	port := 33333
	s := New(ctx, port)

	err := s.Start()
	require.NoError(t, err)
	require.Equal(t, port, s.port)

	err = <-s.Wait()
	require.NoError(t, err)
}
