package client

import (
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/transaction"
)

// Contract defines the contract used to encapsulate the IOTA Client; separated to facilitate unit testing.
type Contract interface {
	// FindTransactionObjects is called to retrieve transactions for the specified query.
	FindTransactionObjects(query api.FindTransactionsQuery) (transaction.Transactions, error)
}
