package httpclient

import (
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"os"
)

const (
	defaultSeverHost = "http://localhost:8080"
	serverHostKey    = "SERVER_HOST"
)

func NewClient() api.ClientInterface {
	serverHost := os.Getenv(serverHostKey)
	if serverHost == "" {
		serverHost = defaultSeverHost
	}

	client, _ := api.NewClient(serverHost)

	return client
}
