package examples_test

import (
	"context"
	"errors"
	"log"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/t-0-network/provider-sdk-go-examples/utils"
	"github.com/t-0-network/provider-sdk-go/api/gen/proto/common"
	networkreq "github.com/t-0-network/provider-sdk-go/api/gen/proto/network"
	"github.com/t-0-network/provider-sdk-go/api/gen/proto/network/networkconnect"
	"github.com/t-0-network/provider-sdk-go/pkg/network"
	"github.com/t-0-network/provider-sdk-go/pkg/provider"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	payInCurrency  = "EUR"
	payoutCurrency = "BRL"
)

func ExamplePayinProviderBasicFlow() {
	shutdownFunc := startTheProviderService() // Start the provider service, which listens for incoming requests from the network.

	networkClient := createClientToInteractWithNetwork() // Create a network service client to interact with the T-ZERO Network.

	// ----- Step 0 (optional): Submit periodically pay-in quotes to the network. The pay-in quotes needs to convert the
	// payment amount from local pay-in currency to the network's settlement currency (USD). For example, if the pay-in currency is EUR,
	// the pay-in quote should be for USD/EUR. Since the base currency is always USD for all the quotes, only the quote currency needs to be specified.
	//
	// It's also possible to submit an amount in payment requests in USD directly, in which case the pay-in quote is not needed.
	resp, err := networkClient.UpdateQuote(context.Background(), payInQuotesRequestExample())
	if err != nil {
		log.Fatalf("Failed to submit pay-in quotes: %v", err)
	}

	log.Printf("Pay-in quotes submitted successfully: %v", resp.Msg)

	// ----- Step 1: When the pay-in provider wants to show the fx rates to the user, it can call the GetPayoutQuote API to
	// get the best available quote for the pay-out currency. For example, the user wants to pay in EUR and send the BRL to the recipient.
	// Then the BRL quote needs to be requested from the network. This represents the USD/BRL rate. Pay-in providers already
	// knows the rate for EUR/USD, so it can calculate the USD amount and then multiply it by the USD/BRL rate to get the final BRL amount.

	payoutQuote, err := networkClient.GetPayoutQuote(context.Background(), connect.NewRequest(&networkreq.GetPayoutQuoteRequest{
		PayoutCurrency: payoutCurrency,
		Amount:         utils.DecimalToProto(decimal.NewFromFloat(500.0)), // Amount in USD, which is the settlement currency of the network.
		QuoteType:      networkreq.QuoteType_QUOTE_TYPE_REALTIME,
	}))
	if err != nil {
		log.Fatalf("Failed to get payout quote: %v", err)
	}
	log.Printf("Payout quote received successfully: %v", payoutQuote.Msg)

	// ----- Step 2: When the pay-in provider receives a payment from the user, it calls the CreatePayment API to initiate the process
	// of paying out to the recipient. It can also specify the pay-out quote to be used for this payment (please pay attention to the expiration time of the quote)
	paymentResp, err := networkClient.CreatePayment(context.Background(), createPaymentRequestExample(payoutQuote))
	if err != nil {
		log.Fatalf("Failed to create payment: %v", err)
	}
	log.Printf("Payment created successfully: %v", paymentResp.Msg)

	// ----- Step 3: The network will process the payment and pay out to the recipient. The pay-in provider will receive
	// a webhook notification via ProviderService.UpdatePayment rpc (should be implemented by the provider). The result
	// can be either success or failure, depending on the payment processing outcome.
	// see the PayInProviderImplementation.UpdatePayment method for more details.

	if err := shutdownFunc(context.Background()); err != nil {
		log.Fatalf("Failed to shutdown provider service: %v", err)
	}
}

func createPaymentRequestExample(payoutQuote *connect.Response[networkreq.GetPayoutQuoteResponse]) *connect.Request[networkreq.CreatePaymentRequest] {
	return connect.NewRequest(&networkreq.CreatePaymentRequest{
		// Unique identifier for the payment, can be used to reference it later.
		PaymentClientId: uuid.NewString(),
		// The currency in which the payout will be made, e.g. BRL.
		PayoutCurrency: payoutCurrency,
		// if the pay-in is not specified, the amount is in the settlement currency (USD). Otherwise, the amount is in the pay-in currency (EUR).
		// In this case it's required that the pay-in quotes are submitted to the network before, so the network can convert the amount to USD.
		Amount: utils.DecimalToProto(decimal.NewFromFloat(500.0)),
		// Pay-in currency is optional, if not specified, the amount is in the settlement currency (USD).
		PayinCurrency: &payInCurrency,
		Sender: &networkreq.CreatePaymentRequest_Sender{
			Sender: &networkreq.CreatePaymentRequest_Sender_PrivatePerson{
				PrivatePerson: &networkreq.CreatePaymentRequest_PrivatePerson{
					PrivatePersonClientId: uuid.NewString(),
					FirstName:             "Daniel",
					LastName:              "Carter",
				},
			},
		},
		Recipient: &networkreq.CreatePaymentRequest_Recipient{
			Recipient: &networkreq.CreatePaymentRequest_Recipient_PrivatePerson{
				PrivatePerson: &networkreq.CreatePaymentRequest_PrivatePerson{
					PrivatePersonClientId: uuid.NewString(),
					FirstName:             "John",
					LastName:              "Doe",
				},
			},
		},
		QuoteId: payoutQuote.Msg.QuoteId,
	})
}

func payInQuotesRequestExample() *connect.Request[networkreq.UpdateQuoteRequest] {

	return connect.NewRequest(&networkreq.UpdateQuoteRequest{
		PayIn: []*networkreq.UpdateQuoteRequest_Quote{
			{
				Currency: payInCurrency,
				//right now only realtime quotes are supported
				QuoteType: networkreq.QuoteType_QUOTE_TYPE_REALTIME,
				// Set the expiration time for the quote
				Expiration: timestamppb.New(time.Now().Add(10 * time.Minute)),
				Timestamp:  timestamppb.Now(),
				Bands: []*networkreq.UpdateQuoteRequest_Quote_Band{
					{
						// ClientQuoteId is a unique identifier for each quote of this provider, which can be used to reference it later.
						ClientQuoteId: uuid.NewString(),
						//band of the quote, e.g. this rate is up to 1000 USD
						MaxAmount: utils.DecimalToProto(decimal.NewFromFloat(1000.0)),
						//rate for the band, USD/EUR = 0.86
						Rate: utils.DecimalToProto(decimal.NewFromFloat(0.86)),
					},
					{
						ClientQuoteId: uuid.NewString(),
						//band of the quote, e.g. this rate is up to 5000 USD payment amount
						MaxAmount: utils.DecimalToProto(decimal.NewFromFloat(5000.0)),
						//rate for this band, USD/EUR = 0.88. Rate for the bigger bands includes risk premium, so the rate is higher than for the smaller bands.
						Rate: utils.DecimalToProto(decimal.NewFromFloat(0.88)),
					},
				},
			},
		},
	})
}

func createClientToInteractWithNetwork() networkconnect.NetworkServiceClient {
	// Replace with your actual private key in hex format.
	yourPrivateKey := network.PrivateKeyHexed("0x7795db2f4499c04d80062c1f1614ff1e427c148e47ed23e387d62829f437b5d8")

	networkClient, err := network.NewServiceClient(
		yourPrivateKey,
		// Optional configuration for the network service client.
		network.WithBaseURL("http://0.0.0.0:8080"), // No need to set, defaults to t-zero network
	)
	if err != nil {
		log.Fatalf("Failed to create network service client: %v", err)
	}
	return networkClient
}

func startTheProviderService() provider.ServerShutdownFn {
	providerServiceHandler, err := provider.NewProviderHandler(
		provider.NetworkPublicKeyHexed(dummyNetworkPublicKey),
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
	return shutdownFunc
}

type PayInProviderImplementation struct{}

func (p *PayInProviderImplementation) UpdatePayment(ctx context.Context, c *connect.Request[networkreq.UpdatePaymentRequest]) (*connect.Response[networkreq.UpdatePaymentResponse], error) {
	// This function is called by the network to notify the provider about the payment status.
	// The provider can use this to update its internal state and notify the user about the payment status.
	// The message contains the payment client ID, which was specified in the CreatePaymentRequest by this provider.
	log.Printf("Received payment update: %v, %s", c.Msg, c.Msg.PaymentClientId)

	// Here you can implement your logic to handle the payment update, e.g. update the payment status in your database.
	// For now, we just return a success response.
	return connect.NewResponse(&networkreq.UpdatePaymentResponse{}), nil
}

func (p *PayInProviderImplementation) UpdateLimit(ctx context.Context, c *connect.Request[networkreq.UpdateLimitRequest]) (*connect.Response[networkreq.UpdateLimitResponse], error) {
	// This function is called by the network to notify about the changes in the limits between providers.
	// This is not required to be implemented by the pay-in provider, but it can be useful to keep track of the limits.
	log.Printf("Received limit update: %v", c.Msg)

	// Here you can implement your logic to handle the limit update, e.g. update the limits in your database.
	// For now, we just return a success response.
	return connect.NewResponse(&networkreq.UpdateLimitResponse{}), nil
}

func (p *PayInProviderImplementation) AppendLedgerEntries(ctx context.Context, c *connect.Request[networkreq.AppendLedgerEntriesRequest]) (*connect.Response[networkreq.AppendLedgerEntriesResponse], error) {
	// Alternatively to the UpdateLimit, the provider can handle all the changes in the ledger entries via this rpc.
	// This is not required to be implemented by the pay-in provider, but if the provider wants to keep track of all the changes
	// in the ledger, it can implement this rpc.
	log.Printf("Appending ledger entries: %v", c.Msg)

	// Here you can implement your logic to handle the ledger entries
	// For now, we just return a success response.
	return connect.NewResponse(&networkreq.AppendLedgerEntriesResponse{}), nil

}

func (p *PayInProviderImplementation) CreatePayInDetails(ctx context.Context, c *connect.Request[networkreq.CreatePayInDetailsRequest]) (*connect.Response[networkreq.CreatePayInDetailsResponse], error) {
	// This function is called by the network to create pay-in details for the payment.
	// This is not necessarily required, this rpc will be called if the network will initiate the pay-in flow.
	// The provider then response with the payment details for the user to pay in. Specifically, the
	log.Printf("Creating pay-in details for payment: %s", c.Msg.PaymentIntentId)

	return connect.NewResponse(&networkreq.CreatePayInDetailsResponse{
		PayInMethod: []*common.PaymentMethod{
			{
				Details: &common.PaymentMethod_Sepa{
					Sepa: &common.SepaPaymentMethod{
						// Example IBAN to which the user should pay in.
						Iban: "IE29AIBK93115212345678",
						// Reference for the user to use when making the payment, this is needed for pay-in provider to identify the payment from user.
						PaymentReference: "paymentRef-1234",
						// This is name associated with the IBAN to initate the payment
						Name: "PAY-IN PROVIDER",
					},
				},
			},
		},
	}), nil

}

func (p *PayInProviderImplementation) PayOut(ctx context.Context, c *connect.Request[networkreq.PayoutRequest]) (*connect.Response[networkreq.PayoutResponse], error) {
	// this function is not required for the pay-in provider flow, but it can be implemented if provider wants to participate in the payout flow as well.
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("PayOut is not implemented for PayInProviderImplementation"))
}
