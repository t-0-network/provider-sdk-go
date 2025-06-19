package examples_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	networkproto "github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network"
	providerClient "github.com/t-0-network/provider-sdk-go/pkg/provider/client"
	providerService "github.com/t-0-network/provider-sdk-go/pkg/provider/service"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	networkPublicKey  = "0x049bb924680bfba3f64d924bf9040c45dcc215b124b5b9ee73ca8e32c050d042c0bbd8dbb98e3929ed5bc2967f28c3a3b72dd5e24312404598bbf6c6cc47708dc7"
	networkPrivateKey = "691db48202ca70d83cc7f5f3aa219536f9bb2dfe12ebb78a7bb634544858ee92"
)

func ExampleNewProviderHandler() {
	// Start the provider service
	providerServiceHandler, err := providerService.NewProviderHandler(
		// Your provider service implementation
		&ProviderServiceImplementation{},
		// Only signed requests coming from this network public key will be accepted
		providerService.WithNetworkPublicKey(networkPublicKey),
	)
	if err != nil {
		log.Fatalf("Failed to create provider service handler: %v", err)
	}

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Duration(10) * time.Second,
		WriteTimeout: time.Duration(10) * time.Second,
		Handler:      h2c.NewHandler(providerServiceHandler, &http2.Server{}),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Create client to consume the service
	providerClient, err := providerClient.NewServiceClient(
		providerClient.WithBaseURL("http://localhost:8080"),
		providerClient.WithNetworkPrivateKeyHexed(networkPrivateKey),
	)
	if err != nil {
		log.Fatalf("Failed to create provider client: %v", err)
	}

	// Make a request to the service
	req := connect.NewRequest(&networkproto.CreatePayInDetailsRequest{
		PaymentId: uuid.New().String(),
	})

	_, err = providerClient.CreatePayInDetails(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to create pay in details: %v", err)
	}

	fmt.Println("Successfully created pay in details")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	// Output:
	// Successfully created pay in details
}

type ProviderServiceImplementation struct{}

func (s *ProviderServiceImplementation) AppendLedgerEntries(
	ctx context.Context, req *connect.Request[networkproto.AppendLedgerEntriesRequest],
) (*connect.Response[networkproto.AppendLedgerEntriesResponse], error) {
	return connect.NewResponse(&networkproto.AppendLedgerEntriesResponse{}), nil
}

func (s *ProviderServiceImplementation) CreatePayInDetails(
	ctx context.Context, req *connect.Request[networkproto.CreatePayInDetailsRequest],
) (*connect.Response[networkproto.CreatePayInDetailsResponse], error) {
	return connect.NewResponse(&networkproto.CreatePayInDetailsResponse{}), nil
}

func (s *ProviderServiceImplementation) PayOut(ctx context.Context, req *connect.Request[networkproto.PayoutRequest],
) (*connect.Response[networkproto.PayoutResponse], error) {
	return connect.NewResponse(&networkproto.PayoutResponse{}), nil
}

func (s *ProviderServiceImplementation) UpdateLimit(
	ctx context.Context, req *connect.Request[networkproto.UpdateLimitRequest],
) (*connect.Response[networkproto.UpdateLimitResponse], error) {
	return connect.NewResponse(&networkproto.UpdateLimitResponse{}), nil
}

func (s *ProviderServiceImplementation) UpdatePayment(
	ctx context.Context, req *connect.Request[networkproto.UpdatePaymentRequest],
) (*connect.Response[networkproto.UpdatePaymentResponse], error) {
	return connect.NewResponse(&networkproto.UpdatePaymentResponse{}), nil
}
