package examples_test

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	networkproto "github.com/t-0-network/provider-sdk-go/api/gen/proto/network"
	"github.com/t-0-network/provider-sdk-go/pkg/network"
)

// ExampleNewServiceClient demonstrates how to create a new network service client
// to interact with the T-0 Network.
func ExampleNewServiceClient() {
	// Replace with your actual private key in hex format.
	yourPrivateKey := network.PrivateKeyHexed("0x7795db2f4499c04d80062c1f1614ff1e427c148e47ed23e387d62829f437b5d8")

	networkClient, err := network.NewServiceClient(
		yourPrivateKey,
		network.WithBaseURL("http://0.0.0.0:8080"), // No need to set, defaults to t-zero network
	)
	if err != nil {
		log.Fatalf("Failed to create network service client: %v", err)
	}

	_, err = networkClient.UpdateQuote(context.Background(), connect.NewRequest(&networkproto.UpdateQuoteRequest{
		// You actual quote data here
	}))
	if err != nil {
		fmt.Println(err)
	}

	// Example will fail as it tries to connect to a fake address using an invalid key.
	// Output:
	// unavailable: dial tcp 0.0.0.0:8080: connect: connection refused
}
