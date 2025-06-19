package examples_test

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	networkproto "github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network"
	"github.com/t-0-network/provider-sdk-go/pkg/network"
)

func ExampleNewServiceClient() {
	providerPrivateKey := "0x7795db2f4499c04d80062c1f1614ff1e427c148e47ed23e387d62829f437b5d8"

	networkClient, err := network.NewServiceClient(
		network.WithBaseURL("http://0.0.0.0:8080"),
		network.WithProviderPrivateKeyHexed(providerPrivateKey),
	)
	if err != nil {
		log.Fatalf("Failed to create network service client: %v", err)
	}

	_, err = networkClient.UpdateQuote(context.Background(), connect.NewRequest(&networkproto.UpdateQuoteRequest{}))
	if err != nil {
		fmt.Println("Need to use a valid public key")
	}

	// Output:
	// Need to use a valid public key
}
