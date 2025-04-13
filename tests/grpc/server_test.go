//go:build integration.grpc

package grpc

import (
	pvz_v1 "github.com/khostya/pvz/pkg/api/v1/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"testing"
)

const (
	ADDRESS    = "localhost:3000"
	AddressKey = "GRPC_ADDRESS"
)

func TestServer(t *testing.T) {
	address := os.Getenv(AddressKey)
	if address == "" {
		address = ADDRESS
	}

	clientConn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer clientConn.Close()

	pvzClient := pvz_v1.NewPVZServiceClient(clientConn)
	_, err = pvzClient.GetPVZList(t.Context(), &pvz_v1.GetPVZListRequest{})
	require.NoError(t, err)
}
