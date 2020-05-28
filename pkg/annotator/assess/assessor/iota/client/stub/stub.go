/*******************************************************************************
 * Copyright 2020 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package stub

import (
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/transaction"
)

// client is a receiver that encapsulates dependencies.
type client struct {
	resultIndex          uint
	findTxObjectsResults []transaction.Transactions
	findTxObjectsErr     error
}

// New is a factory function that returns an initialized instance.
func New(findTxObjectsResults []transaction.Transactions, findTxObjectErr error) *client {
	return &client{
		findTxObjectsResults: findTxObjectsResults,
		findTxObjectsErr:     findTxObjectErr,
	}
}

// FindTransactionObjects is called to retrieve transactions for the specified query.
func (c *client) FindTransactionObjects(query api.FindTransactionsQuery) (transaction.Transactions, error) {
	if c.findTxObjectsErr != nil {
		return nil, c.findTxObjectsErr
	}
	result := c.findTxObjectsResults[c.resultIndex]
	c.resultIndex++
	return result, nil
}
