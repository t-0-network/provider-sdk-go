package examples_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/api/tzero/v1/common"
	"github.com/t-0-network/provider-sdk-go/network"

	"github.com/t-0-network/provider-sdk-go/api/tzero/v1/payment"
	"github.com/t-0-network/provider-sdk-go/api/tzero/v1/payment/paymentconnect"
	"github.com/t-0-network/provider-sdk-go/crypto"
	"github.com/t-0-network/provider-sdk-go/provider"
)

var (
	dummyNetworkPublicKey  = "0x049bb924680bfba3f64d924bf9040c45dcc215b124b5b9ee73ca8e32c050d042c0bbd8dbb98e3929ed5bc2967f28c3a3b72dd5e24312404598bbf6c6cc47708dc7"
	dummyNetworkPrivateKey = "691db48202ca70d83cc7f5f3aa219536f9bb2dfe12ebb78a7bb634544858ee92"
)

func ExampleProviderServiceHandler() {
	// Initialize a provider service handler using your implementation of the
	// networkconnect.ProviderServiceHandler interface.
	providerServiceHandler, err := provider.NewHttpHandler(
		// Provide the T-ZERO Network Public Key in hex format. This key is used to verify
		// the signatures of incoming requests.
		provider.NetworkPublicKeyHexed(dummyNetworkPublicKey),
		// Your provider service implementation
		provider.Handler(paymentconnect.NewProviderServiceHandler, paymentconnect.ProviderServiceHandler(&ProviderServiceImplementation{})),
	)
	if err != nil {
		log.Fatalf("Failed to create provider service handler: %v", err)
	}

	// Start an HTTP server with the provider service handler,
	shutdownFunc, err := provider.StartServer(
		providerServiceHandler,
		provider.WithAddr(":8080"),
	)
	if err != nil {
		log.Fatalf("Failed to start provider server: %v", err)
	}

	// Create a provider client to connect to the provider service.
	providerClient, err := newProviderClient(dummyNetworkPrivateKey)
	if err != nil {
		log.Fatalf("Failed to create provider client: %v", err)
	}

	// Build a CreatePayInDetails request
	req := connect.NewRequest(&payment.UpdateLimitRequest{
		Limits: []*payment.UpdateLimitRequest_Limit{
			{
				Version:       1,
				CounterpartId: 3,
				PayoutLimit: &common.Decimal{
					Unscaled: 100,
					Exponent: 0,
				},
				CreditLimit: &common.Decimal{
					Unscaled: 1000,
					Exponent: 0,
				},
				CreditUsage: &common.Decimal{
					Unscaled: 900,
					Exponent: 0,
				},
			}},
	})

	// Use the providerClient to call the UpdateLimit method
	_, err = providerClient.UpdateLimit(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to update limits: %v", err)
	}

	fmt.Println("Successfully updated limits")

	if err := shutdownFunc(context.Background()); err != nil {
		log.Fatalf("Failed to shutdown provider service: %v", err)
	}

	// Output:
	// Successfully updated limits
}

func newProviderClient(privateKey string) (paymentconnect.ProviderServiceClient, error) {
	// Create an http client with a custom transport which signs the raw
	// request body using the dummy network private key.
	signFn, err := crypto.NewSignerFromHex(privateKey)
	if err != nil {
		return nil, fmt.Errorf("creating signer function: %w", err)
	}

	// Create a custom HTTP client with a custom transport to sign requests.
	httpClient := http.Client{
		Timeout:   15 * time.Second,
		Transport: network.NewSigningTransport(signFn, time.Now),
	}

	// Initialize the provider service client using custom HTTP client.
	return paymentconnect.NewProviderServiceClient(&httpClient, "http://127.0.0.1:8080"), nil
}

type ProviderServiceImplementation struct{}

func (s *ProviderServiceImplementation) ApprovePaymentQuotes(ctx context.Context, c *connect.Request[payment.ApprovePaymentQuoteRequest]) (*connect.Response[payment.ApprovePaymentQuoteResponse], error) {
	return connect.NewResponse(&payment.ApprovePaymentQuoteResponse{}), nil
}

var _ paymentconnect.ProviderServiceHandler = (*ProviderServiceImplementation)(nil)

func (s *ProviderServiceImplementation) AppendLedgerEntries(
	ctx context.Context, req *connect.Request[payment.AppendLedgerEntriesRequest],
) (*connect.Response[payment.AppendLedgerEntriesResponse], error) {
	return connect.NewResponse(&payment.AppendLedgerEntriesResponse{}), nil
}

func (s *ProviderServiceImplementation) PayOut(ctx context.Context, req *connect.Request[payment.PayoutRequest],
) (*connect.Response[payment.PayoutResponse], error) {
	return connect.NewResponse(&payment.PayoutResponse{}), nil
}

func (s *ProviderServiceImplementation) UpdateLimit(
	ctx context.Context, req *connect.Request[payment.UpdateLimitRequest],
) (*connect.Response[payment.UpdateLimitResponse], error) {
	return connect.NewResponse(&payment.UpdateLimitResponse{}), nil
}

func (s *ProviderServiceImplementation) UpdatePayment(
	ctx context.Context, req *connect.Request[payment.UpdatePaymentRequest],
) (*connect.Response[payment.UpdatePaymentResponse], error) {
	return connect.NewResponse(&payment.UpdatePaymentResponse{}), nil
}
