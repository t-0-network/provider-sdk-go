package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/service/gen/proto/network"
)

func (s *ProviderService) UpdatePayment(
	ctx context.Context,
	req *connect.Request[network.UpdatePaymentRequest],
) (*connect.Response[network.UpdatePaymentResponse], error) {
	return connect.NewResponse(&network.UpdatePaymentResponse{}), nil
}
