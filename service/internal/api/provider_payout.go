package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/service/gen/proto/network"
)

func (s *ProviderService) PayOut(
	ctx context.Context,
	_ *connect.Request[network.PayoutRequest],
) (*connect.Response[network.PayoutResponse], error) {
	_, err := s.networkClient.UpdateQuote(ctx, connect.NewRequest(&network.UpdateQuoteRequest{}))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&network.PayoutResponse{}), nil
}
