package examples_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	networkproto "github.com/t-0-network/provider-sdk-go/api/gen/proto/network"
	"github.com/t-0-network/provider-sdk-go/api/gen/proto/network/networkconnect"
	"github.com/t-0-network/provider-sdk-go/pkg/constant"
	"github.com/t-0-network/provider-sdk-go/pkg/crypto"
	"github.com/t-0-network/provider-sdk-go/pkg/provider"
)

var (
	dummyNetworkPublicKey  = "0x049bb924680bfba3f64d924bf9040c45dcc215b124b5b9ee73ca8e32c050d042c0bbd8dbb98e3929ed5bc2967f28c3a3b72dd5e24312404598bbf6c6cc47708dc7"
	dummyNetworkPrivateKey = "691db48202ca70d83cc7f5f3aa219536f9bb2dfe12ebb78a7bb634544858ee92"
)

func ExampleNewProviderHandler() {
	// Initialize a provider service handler using your implementation of the
	// networkconnect.ProviderServiceHandler interface.
	providerServiceHandler, err := provider.NewProviderHandler(
		// Provide the T-ZERO Network Public Key in hex format. This key is used to verify
		// the signatures of incoming requests.
		provider.NetworkPublicKeyHexed(dummyNetworkPublicKey),
		// Your provider service implementation
		&ProviderServiceImplementation{},
	)
	if err != nil {
		log.Fatalf("Failed to create provider service handler: %v", err)
	}

	// Start an HTTP server with the provider service handler,
	shutdownFunc := provider.StartServer(
		providerServiceHandler,
		provider.WithAddr(":8080"),
	)

	// Create a provider client to connect to the provider service.
	providerClient, err := newProviderClient(dummyNetworkPrivateKey)
	if err != nil {
		log.Fatalf("Failed to create provider client: %v", err)
	}

	// Build a CreatePayInDetails request
	req := connect.NewRequest(&networkproto.CreatePayInDetailsRequest{
		PaymentId: "a-payment-unique-identifier",
	})

	// Use the providerClient to call the CreatePayInDetails method
	_, err = providerClient.CreatePayInDetails(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to create pay in details: %v", err)
	}

	fmt.Println("Successfully created pay in details")

	if err := shutdownFunc(context.Background()); err != nil {
		log.Fatalf("Failed to shutdown provider service: %v", err)
	}

	// Output:
	// Successfully created pay in details
}

func newProviderClient(privateKey string) (networkconnect.ProviderServiceClient, error) {
	// Create an http client with a custom transport which signs the raw
	// request body using the dummy network private key.
	signFn, err := crypto.NewSignerFromHex(privateKey)
	if err != nil {
		return nil, fmt.Errorf("creating signer function: %w", err)
	}

	// Create a custom HTTP client with a custom transport to sign requests.
	httpClient := http.Client{
		Timeout: 15 * time.Second,
		Transport: &signingTransport{
			transport: http.DefaultTransport,
			signFn:    signFn,
		},
	}

	// Initialize the provider service client using custom HTTP client.
	return networkconnect.NewProviderServiceClient(&httpClient, "http://127.0.0.1:8080"), nil
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

type signingTransport struct {
	transport http.RoundTripper
	signFn    crypto.SignFn
}

func (t *signingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read and restore request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request body: %w", err)
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(body))

	// Hash the request body using LegacyKeccak256
	digest := crypto.LegacyKeccak256(body)

	signature, pubKeyBytes, err := t.signFn(digest)
	if err != nil {
		return nil, fmt.Errorf("signing request body: %w", err)
	}

	// Set headers
	req.Header.Set(constant.PublicKeyHeader, "0x"+hex.EncodeToString(pubKeyBytes))
	req.Header.Set(constant.SignatureHeader, "0x"+hex.EncodeToString(signature))

	return t.transport.RoundTrip(req)
}
