package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/service/gen/proto/network"
)

func (s *ProviderService) CreatePayInDetails(
	ctx context.Context,
	req *connect.Request[network.CreatePayInDetailsRequest],
) (*connect.Response[network.CreatePayInDetailsResponse], error) {
	return connect.NewResponse(&network.CreatePayInDetailsResponse{}), nil
}
