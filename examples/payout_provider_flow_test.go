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
	networkreq "github.com/t-0-network/provider-sdk-go/api/gen/proto/network"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ExamplePayoutProviderBasicFlow() {
	// Start the provider service, which listens for incoming requests from the network.
	shutdownFunc := startTheProviderService(&PayOutProviderImplementation{})

	// Create a network service client to interact with the T-ZERO Network.
	networkClient := createClientToInteractWithNetwork()

	// ----- Step 0 (optional): Submit periodically pay-out quotes to the network. The pay-out quotes are needed to convert the
	// payment amount from network's settlement currency (USD) to the pay-out currency. For example, if the pay-out currency is BRL,
	// the pay-out quote should be for USD/BRL. Since the base currency is always USD for all the quotes, only the quote currency needs
	// to be specified in the request.
	resp, err := networkClient.UpdateQuote(context.Background(), payOutQuotesRequestExample())
	if err != nil {
		log.Fatalf("Failed to submit pay-out quotes: %v", err)
	}

	// ------ Step 1: The network will call the provider to request a pay-out for a specific payment.
	// See the PayOutProviderImplementation.PayOut method for the implementation details.

	// ------ Step 2: The provider notifies the network about the pay-out status using the UpdatePayout RPC.
	networkClient.UpdatePayout(context.Background(), connect.NewRequest(&networkreq.UpdatePayoutRequest{
		PaymentId: 1, // This is the payment ID that the network provided in the pay-out request.
		PayoutId:  1, // This is the pay-out ID that the network provided in the pay-out request.
		// either success or failure result
		Result: &networkreq.UpdatePayoutRequest_Success_{
			Success: &networkreq.UpdatePayoutRequest_Success{},
		},
		//Result:    &networkreq.UpdatePayoutRequest_Failure_{},
	}))

	log.Printf("Pay-out quotes submitted successfully: %v", resp.Msg)

	if err := shutdownFunc(context.Background()); err != nil {
		log.Fatalf("Failed to shutdown provider service: %v", err)
	}
}

func payOutQuotesRequestExample() *connect.Request[networkreq.UpdateQuoteRequest] {
	return connect.NewRequest(&networkreq.UpdateQuoteRequest{
		// There are 2 repeated fields in the request, one for the pay-in quotes and one for the pay-out quotes.
		// So the provider can either submit pay-in quotes, pay-out quotes or both.
		PayOut: []*networkreq.UpdateQuoteRequest_Quote{
			{
				// specify the currency for the pay-out quote, e.g. BRL. In this case the rate is for USD/BRL.
				Currency: payoutCurrency,
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
						//rate for the band, USD/BRL = 5.56
						// This means, that the provider is willing to pay out 5.56 BRL for each USD
						Rate: utils.DecimalToProto(decimal.NewFromFloat(5.56)),
					},
					{
						ClientQuoteId: uuid.NewString(),
						// band of the quote, e.g. this rate is up to 5000 USD payment amount
						MaxAmount: utils.DecimalToProto(decimal.NewFromFloat(5000.0)),
						// rate for this band, USD/EUR = 0.88. Rate for the bigger bands includes risk premium,
						// so the rate is lower for the bigger bands, if it's a payout quote.
						// This means, that the provider is willing to pay out 5.46 BRL for each USD
						Rate: utils.DecimalToProto(decimal.NewFromFloat(5.46)),
					},
				},
			},
		},
	})
}

type PayOutProviderImplementation struct{}

func (p *PayOutProviderImplementation) PayOut(ctx context.Context, c *connect.Request[networkreq.PayoutRequest]) (*connect.Response[networkreq.PayoutResponse], error) {
	// This function is called by the network to request a pay-out for a specific payment.
	// The provider should implement the logic to process the pay-out and initiate the transfer to the recipient
	log.Printf("Received pay-out request: %v", c.Msg)

	// The request contains the paymentID, payoutID, and the amount to be paid out.
	// Provider should store the paymentID and payoutID in its database to track the payment status, and notify the
	// network about the pay-out status later using these IDs.

	// Here you can implement your logic to process the pay-out, e.g. initiate a transfer to the recipient.
	// For now, we just return a success response.

	// If the payment is processed by the payout provider banking system immediately, the provider can notify the network
	// by calling the NetworkService.UpdatePayout RPC inside this handler.
	// Otherwise, the provider can notify the network later, when the payment is processed.
	// In this case just return here a success response to the network.
	return connect.NewResponse(&networkreq.PayoutResponse{}), nil
}

func (p *PayOutProviderImplementation) UpdateLimit(ctx context.Context, c *connect.Request[networkreq.UpdateLimitRequest]) (*connect.Response[networkreq.UpdateLimitResponse], error) {
	// This function is called by the network to notify about the changes in the limits between providers.
	// This is not required to be implemented by the pay-in provider, but it can be useful to keep track of the limits.
	log.Printf("Received limit update: %v", c.Msg)

	// Here you can implement your logic to handle the limit update, e.g. update the limits in your database.
	// For now, we just return a success response.
	return connect.NewResponse(&networkreq.UpdateLimitResponse{}), nil
}

func (p *PayOutProviderImplementation) AppendLedgerEntries(ctx context.Context, c *connect.Request[networkreq.AppendLedgerEntriesRequest]) (*connect.Response[networkreq.AppendLedgerEntriesResponse], error) {
	// Alternatively to the UpdateLimit, the provider can handle all the changes in the ledger entries via this rpc.
	// This is not required to be implemented by the pay-in provider, but if the provider wants to keep track of all the changes
	// in the ledger, it can implement this rpc.
	log.Printf("Appending ledger entries: %v", c.Msg)

	// Here you can implement your logic to handle the ledger entries
	// For now, we just return a success response.
	return connect.NewResponse(&networkreq.AppendLedgerEntriesResponse{}), nil
}

func (p *PayOutProviderImplementation) CreatePayInDetails(ctx context.Context, c *connect.Request[networkreq.CreatePayInDetailsRequest]) (*connect.Response[networkreq.CreatePayInDetailsResponse], error) {
	// this function is not required for the pay-out provider flow, but it can be implemented if provider wants to participate as pay-in provider as well.
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("CreatePayInDetails is not implemented for PayOutProviderImplementation"))
}

func (p *PayOutProviderImplementation) UpdatePayment(ctx context.Context, c *connect.Request[networkreq.UpdatePaymentRequest]) (*connect.Response[networkreq.UpdatePaymentResponse], error) {
	// this function is not required for the pay-out provider flow, but it can be implemented if provider wants to participate as pay-in provider as well.
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("UpdatePayment is not implemented for PayOutProviderImplementation"))
}
