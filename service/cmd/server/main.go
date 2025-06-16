package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/t-0-network/provider-sdk-go/service/gen/proto/network/networkconnect"
	"github.com/t-0-network/provider-sdk-go/service/internal/api"
	"github.com/t-0-network/provider-sdk-go/service/internal/config"
	"github.com/t-0-network/provider-sdk-go/service/internal/network"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	networkClient, err := network.NewNetworkServiceClient(cfg.NetworkClient)
	if err != nil {
		log.Fatalf("Failed to create network service client: %v", err)
	}

	providerService := api.NewProviderService(networkClient)
	path, provideHandler := networkconnect.NewProviderServiceHandler(providerService)

	mux := http.NewServeMux()
	mux.Handle(path, provideHandler)

	httpServer := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
	}

	// Listen for interrupt signal for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		log.Println("Received interrupt signal, shutting down server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	if err := httpServer.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Failed to start server:", err)
	}
}
