module github.com/rosetta-dogecoin/rosetta-dogecoin

go 1.13

require (
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/coinbase/rosetta-bitcoin v0.0.9 // indirect
	github.com/coinbase/rosetta-sdk-go v0.6.10
	github.com/dgraph-io/badger/v2 v2.2007.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.16.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/tools v0.1.0 // indirect
	honnef.co/go/tools v0.1.3 // indirect
)

replace github.com/coinbase/rosetta-bitcoin/configuration => ./configuration
