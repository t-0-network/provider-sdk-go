module github.com/t-0-network/provider-sdk-go-examples

go 1.24.4

require (
	connectrpc.com/connect v1.18.1
	github.com/google/uuid v1.6.0
	github.com/shopspring/decimal v1.4.0
	github.com/t-0-network/provider-sdk-go v0.0.0-20250625131743-be3754b69efc
	google.golang.org/protobuf v1.36.6
)

replace github.com/t-0-network/provider-sdk-go => ../../provider-sdk-go

require (
	github.com/btcsuite/btcd/btcec/v2 v2.3.5 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
)
