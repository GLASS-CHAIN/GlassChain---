// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
)


type ContractLog struct {
	Address common.Address

	TxHash common.Hash

	Index int

	Topics []common.Hash

	Data []byte
	BlockNumber uint64 `json:"blockNumber"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex"`
	// hash of the block in which the transaction was included
	BlockHash common.Hash `json:"blockHash"`
}

func (log *ContractLog) PrintLog() {
	log15.Debug("!Contract Log!", "Contract address", log.Address.String(), "TxHash", log.TxHash.Hex(), "Log Index", log.Index, "Log Topics", log.Topics, "Log Data", common.Bytes2Hex(log.Data))
}
