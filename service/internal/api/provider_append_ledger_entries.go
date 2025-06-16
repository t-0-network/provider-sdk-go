package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/service/gen/proto/network"
)

func (s *ProviderService) AppendLedgerEntries(
	ctx context.Context,
	req *connect.Request[network.AppendLedgerEntriesRequest],
) (*connect.Response[network.AppendLedgerEntriesResponse], error) {
	return connect.NewResponse(&network.AppendLedgerEntriesResponse{}), nil
}
