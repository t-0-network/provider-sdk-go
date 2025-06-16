package api

import (
	"github.com/t-0-network/provider-sdk-go/service/gen/proto/network/networkconnect"
)

var _ networkconnect.ProviderServiceHandler = (*ProviderService)(nil)

type ProviderService struct {
	networkClient networkconnect.NetworkServiceClient
}

func NewProviderService(
	networkClient networkconnect.NetworkServiceClient,
) *ProviderService {
	return &ProviderService{
		networkClient: networkClient,
	}
}
