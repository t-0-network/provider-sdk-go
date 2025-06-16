package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/service/gen/proto/network"
)

func (s *ProviderService) PayOut(
	ctx context.Context,
	req *connect.Request[network.PayoutRequest],
) (*connect.Response[network.PayoutResponse], error) {
	return connect.NewResponse(&network.PayoutResponse{}), nil
}
