package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/service/gen/proto/network"
)

func (s *ProviderService) UpdateLimit(
	ctx context.Context,
	req *connect.Request[network.UpdateLimitRequest],
) (*connect.Response[network.UpdateLimitResponse], error) {
	return connect.NewResponse(&network.UpdateLimitResponse{}), nil
}
