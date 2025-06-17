package internal

import (
	"context"
	"errors"
	"fmt"
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

// AppContainer holds all dependencies for the app.
type AppContainer struct {
	Cfg             config.Config
	NetworkClient   networkconnect.NetworkServiceClient
	ProviderHandler networkconnect.ProviderServiceHandler
}

// NewAppContainer wires all dependencies.
func NewAppContainer() (*AppContainer, error) {
	var err error

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load app config: %w", err)
	}

	networkClient, err := network.NewNetworkServiceClient(cfg.NetworkClient)
	if err != nil {
		log.Fatalf("Failed to create network service client: %v", err)
	}

	providerHandler := api.NewProviderService(networkClient)

	return &AppContainer{
		Cfg:             cfg,
		NetworkClient:   networkClient,
		ProviderHandler: providerHandler,
	}, nil
}

// Run starts the application.
func (*AppContainer) Run() error {
	app, err := NewAppContainer()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	path, provideServiceHandler := networkconnect.NewProviderServiceHandler(app.ProviderHandler)
	mux.Handle(path, provideServiceHandler)

	httpServer := &http.Server{
		Addr:         app.Cfg.Server.Address,
		ReadTimeout:  app.Cfg.Server.ReadTimeout,
		WriteTimeout: app.Cfg.Server.WriteTimeout,
		Handler:      h2c.NewHandler(mux, &http2.Server{}),
	}

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

	log.Printf("Starting HTTP server on %s", app.Cfg.Server.Address)
	if err := httpServer.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}
