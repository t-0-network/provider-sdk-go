package examples_test

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/t-0-network/provider-sdk-go/api/ivms101/v1/ivms"
	"github.com/t-0-network/provider-sdk-go/api/tzero/v1/common"
	"github.com/t-0-network/provider-sdk-go/api/tzero/v1/payment"
	"github.com/t-0-network/provider-sdk-go/api/tzero/v1/payment/paymentconnect"
	"github.com/t-0-network/provider-sdk-go/network"
)

// ExampleNewServiceClient demonstrates how to create a new network service client
// to interact with the T-0 Network.
func ExampleNewServiceClient() {
	// Replace with your actual private key in hex format.
	yourPrivateKey := network.PrivateKeyHexed("0x7795db2f4499c04d80062c1f1614ff1e427c148e47ed23e387d62829f437b5d8")

	networkClient, err := network.NewServiceClient(
		yourPrivateKey,
		paymentconnect.NewNetworkServiceClient,
	)
	if err != nil {
		log.Fatalf("Failed to create network service client: %v", err)
	}

	resp, err := networkClient.CreatePayment(context.Background(), connect.NewRequest(&payment.CreatePaymentRequest{
		PaymentClientId: uuid.NewString(),
		Amount: &payment.PaymentAmount{
			Amount: &payment.PaymentAmount_SettlementAmount{SettlementAmount: &common.Decimal{
				Unscaled: 2,
				Exponent: 0,
			}},
		},
		Currency: "PHP",
		PaymentDetails: &common.PaymentDetails{
			Details: &common.PaymentDetails_Pesonet_{
				Pesonet: &common.PaymentDetails_Pesonet{
					RecipientFinancialInstitution: "TestInsitution",
					RecipientIdentifier:           "123456",
					RecipientAccountName:          "TestAccount",
				},
			},
		},
		TravelRuleData: &payment.CreatePaymentRequest_TravelRuleData{
			Originator: []*ivms.Person{
				{Person: &ivms.Person_NaturalPerson{
					NaturalPerson: &ivms.NaturalPerson{
						Name: &ivms.NaturalPersonName{
							LocalNameIdentifiers: []*ivms.LocalNaturalPersonNameId{
								{
									PrimaryIdentifier:   "Test",
									SecondaryIdentifier: "Person",
									NameIdentifierType:  ivms.NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_BIRT,
								},
							},
						},
						GeographicAddresses: []*ivms.Address{{
							AddressType:    ivms.AddressTypeCode_ADDRESS_TYPE_CODE_HOME,
							StreetName:     "TestStreet",
							BuildingNumber: "54",
							PostCode:       "54642",
							TownName:       "TestTown",
							Country:        "TestCountry",
						}},
					},
				},
				},
			},
			Beneficiary: []*ivms.Person{
				{Person: &ivms.Person_NaturalPerson{
					NaturalPerson: &ivms.NaturalPerson{
						Name: &ivms.NaturalPersonName{
							LocalNameIdentifiers: []*ivms.LocalNaturalPersonNameId{
								{
									PrimaryIdentifier:   "Test",
									SecondaryIdentifier: "Person",
									NameIdentifierType:  ivms.NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_BIRT,
								},
							},
						},
						GeographicAddresses: []*ivms.Address{{
							AddressType:    ivms.AddressTypeCode_ADDRESS_TYPE_CODE_HOME,
							StreetName:     "TestStreet",
							BuildingNumber: "54",
							PostCode:       "54642",
							TownName:       "TestTown",
							Country:        "TestCountry",
						}},
					},
				},
				},
			},
		},
	}))

	if err != nil {
		log.Fatalf("Failed to create payment: %v", err)
	}
	fmt.Printf(" Response: %+v\n", resp.Msg.Result)

	// Example will fail as it tries to connect to a fake address using an unknown key.
	// Output:
	// unavailable: dial tcp 0.0.0.0:8080: connect: connection refused
}
