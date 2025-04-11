package httpserver

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	e := echo.New()
	duration := time.Hour * 24
	port := uint16(1)

	s := New(e, WriteTimeout(duration), ReadTimeout(duration), IdleTimeout(duration), Port(port))
	require.NotNil(t, s)
	require.Equal(t, duration, s.readTimeout)
	require.Equal(t, duration, s.writeTimeout)
	require.Equal(t, duration, s.idleTimeout)
	require.Equal(t, ":"+strconv.Itoa(int(port)), s.address)

	s.Start()

	err := s.Shutdown(context.Background())
	require.NoError(t, err)
}
