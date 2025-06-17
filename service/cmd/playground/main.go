package main

import (
	"context"
	"log"

	"connectrpc.com/connect"
	v1 "github.com/t-0-network/provider-sdk-go/service/gen/proto/network"
	"github.com/t-0-network/provider-sdk-go/service/internal/config"
	"github.com/t-0-network/provider-sdk-go/service/internal/network"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	networkClient, err := network.NewNetworkServiceClient(cfg.NetworkClient)
	if err != nil {
		log.Fatalf("Failed to create network service client: %v", err)
	}

	_, err = networkClient.UpdateQuote(context.Background(), connect.NewRequest(&v1.UpdateQuoteRequest{}))
	if err != nil {
		log.Fatalf("Failed to get KYC data: %v", err)
	}

	log.Println("KYC data retrieved successfully")
}
