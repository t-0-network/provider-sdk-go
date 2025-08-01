package examples_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/api/gen/proto/tzero/v1/common"
	"github.com/t-0-network/provider-sdk-go/api/gen/proto/tzero/v1/payment"
	"github.com/t-0-network/provider-sdk-go/pkg/network"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ExampleNewServiceClient demonstrates how to create a new network service client
// to interact with the T-0 Network.
func ExampleNewServiceClient() {
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

	// Example request
	req := payment.UpdateQuoteRequest{
		PayOut: []*payment.UpdateQuoteRequest_Quote{
			{
				Currency:  "BRL",
				QuoteType: payment.QuoteType_QUOTE_TYPE_REALTIME,
				Bands: []*payment.UpdateQuoteRequest_Quote_Band{
					{
						ClientQuoteId: "quote-id-1",
						MaxAmount: &common.Decimal{
							Unscaled: 100000, // 1,000.00 BRL
							Exponent: -2,
						},
						Rate: &common.Decimal{
							Unscaled: 551, // 5.51 USD/BRL
							Exponent: -2,
						},
					},
					{
						ClientQuoteId: "quote-id-2",
						MaxAmount: &common.Decimal{
							Unscaled: 500000, // 5,000.00 BRL
							Exponent: -2,
						},
						Rate: &common.Decimal{
							Unscaled: 550, // 5.50 USD/BRL (slightly better rate)
							Exponent: -2,
						},
					},
					{
						ClientQuoteId: "quote-id-3",
						MaxAmount: &common.Decimal{
							Unscaled: 1500000, // 15,000.00 BRL
							Exponent: -2,
						},
						Rate: &common.Decimal{
							Unscaled: 549, // 5.49 USD/BRL (best rate)
							Exponent: -2,
						},
					},
				},
				Expiration: timestamppb.New(time.Now().Add(15 * time.Minute)),
				Timestamp:  timestamppb.Now(),
			},
		},
		PayIn: []*payment.UpdateQuoteRequest_Quote{
			{
				Currency:  "EUR",
				QuoteType: payment.QuoteType_QUOTE_TYPE_REALTIME,
				Bands: []*payment.UpdateQuoteRequest_Quote_Band{
					{
						ClientQuoteId: "eur-quote-id-1",
						MaxAmount: &common.Decimal{
							Unscaled: 100000, // 1000.00 EUR
							Exponent: -2,
						},
						Rate: &common.Decimal{
							Unscaled: 8069, // 0.8069 USD/EUR
							Exponent: -4,
						},
					},
					{
						ClientQuoteId: "eur-quote-id-2",
						MaxAmount: &common.Decimal{
							Unscaled: 500000, // 5000.00 EUR
							Exponent: -2,
						},
						Rate: &common.Decimal{
							Unscaled: 8070, // 0.8070 USD/EUR (slightly better rate)
							Exponent: -4,
						},
					},
					{
						ClientQuoteId: "eur-quote-id-3",
						MaxAmount: &common.Decimal{
							Unscaled: 1500000, // 15000.00 EUR
							Exponent: -2,
						},
						Rate: &common.Decimal{
							Unscaled: 8071, // 0.8071 USD/EUR (best rate)
							Exponent: -4,
						},
					},
				},
				Expiration: timestamppb.New(time.Now().Add(15 * time.Minute)),
				Timestamp:  timestamppb.Now(),
			},
		},
	}

	_, err = networkClient.UpdateQuote(context.Background(), connect.NewRequest(&req))
	if err != nil {
		fmt.Println(err)
	}

	// Example will fail as it tries to connect to a fake address using an unknown key.
	// Output:
	// unavailable: dial tcp 0.0.0.0:8080: connect: connection refused
}
