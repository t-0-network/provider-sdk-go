package payment_intent

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/api/gen/proto/common"
	. "github.com/t-0-network/provider-sdk-go/api/gen/proto/payment_intent/provider"
	. "github.com/t-0-network/provider-sdk-go/api/gen/proto/payment_intent/provider/providerconnect"
)

func ExamplePaymentIntentPayInProvider() {
	// PayIn provider will implement the ProviderServiceHandler interface
	// which has only 2 methods:
	// 1. CreatePaymentIntent - to create a payment intent and return the list of
	//    available payment methods along with the URL to redirect user to make the payment.
	// 2. ConfirmPayout - to confirm the payout after the payment is completed successfully.
	path, handler := NewProviderServiceHandler(&PayInProviderServiceHandler{})
	startServer(path, handler)

	// PayIn provider will interact with the network using the NetworkServiceClient interface.
	// It will use ConfirmPayment/RejectPaymentIntent rpcs to notify the network about the payment intent status.
	// ConfirmSettlement rpc should be used to notify the network about the settlement transfer (in case of pre-settlement).
	networkClient := createNetworkClient()

	// The flow starts when the network call the CreatePaymentIntent method of the PayIn provider.

	// Pay-in provider will return the list of available payment methods, and when it receives the payment from the payer,
	// it will call the ConfirmPayout method to confirm the payout.
	_, err := networkClient.ConfirmPayment(context.Background(), connect.NewRequest(&ConfirmPaymentRequest{
		PaymentIntentId: 123,
		PaymentMethod:   common.PaymentMethodType_PAYMENT_METHOD_TYPE_CARD,
	}))
	if err != nil {
		// Handle error
	}

	// if the payment collection was not successful, the provider will call RejectPaymentIntent method to notify
	//the network about the failure.
	_, err = networkClient.RejectPaymentIntent(context.Background(), connect.NewRequest(&RejectPaymentIntentRequest{
		PaymentIntentId: 123,
		Reason:          "Payment collection failed",
	}))
	if err != nil {
		// Handle error
	}

	// Next step would be to transfer the settlement amount to the pay-out provider, and
	// then call the ConfirmSettlement endpoint
	_, err = networkClient.ConfirmSettlement(context.Background(), connect.NewRequest(&ConfirmSettlementRequest{
		Blockchain:      common.Blockchain_BLOCKCHAIN_TRON,
		TxHash:          "tx-hash-of-the-pre-settlement-transfer",
		PaymentIntentId: []uint64{123}, // one settlement may include several payment intents
	}))
	if err != nil {
		// Handle error
	}

	// And the last step - ConfirmPayout rpc will be called by Network to finalize the process.
}

func createNetworkClient() NetworkServiceClient {
	httpClient := http.DefaultClient
	networkClient := NewNetworkServiceClient(httpClient, "tzero-network-url")
	return networkClient
}

func startServer(path string, handler http.Handler) {
	_, _ = path, handler
}

type PayInProviderServiceHandler struct {
	// Add any necessary fields for the service handler
}

var _ ProviderServiceHandler = (*PayInProviderServiceHandler)(nil)

func (p *PayInProviderServiceHandler) CreatePaymentIntent(
	ctx context.Context,
	req *connect.Request[CreatePaymentIntentRequest],
) (*connect.Response[CreatePaymentIntentResponse], error) {
	// payment intent id is the idempotency key for the payment
	_ = req.Msg.PaymentIntentId
	// pay-in amount to be collected from the payer
	_ = req.Msg.Amount
	// amount is expressed in the pay-in currency
	_ = req.Msg.Currency

	// payment intent should be saved in the database or parameters can be embedded in the URL

	// provider will generate the list of available payment methods along with the URL to redirect user
	methods := []*CreatePaymentIntentResponse_PaymentMethod{
		{
			// This is the URL where the client should be redirected to make the payment.
			PaymentUrl: fmt.Sprintf("https://example.com/pay/%d", req.Msg.PaymentIntentId),
			// Enum of available payment methods includes SEPA, SWIFT, CARD.
			PaymentMethod: common.PaymentMethodType_PAYMENT_METHOD_TYPE_CARD,
		},
	}

	return connect.NewResponse(&CreatePaymentIntentResponse{
		PaymentMethods: methods,
	}), nil
}

func (p *PayInProviderServiceHandler) ConfirmPayout(
	ctx context.Context,
	req *connect.Request[ConfirmPayoutRequest],
) (*connect.Response[ConfirmPayoutResponse], error) {
	// confirm payout is just a notification that the payment was completed successfully. Nothing to return in the response here.
	_ = req.Msg.PaymentIntentId

	return connect.NewResponse(&ConfirmPayoutResponse{}), nil
}
