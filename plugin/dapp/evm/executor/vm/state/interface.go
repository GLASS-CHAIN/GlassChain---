// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package state

import (
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/model"
)

type EVMStateDB interface {
	CreateAccount(string, string, string, string)

	SubBalance(string, string, uint64)

	AddBalance(string, string, uint64)

	GetBalance(string) uint64

	GetNonce(string) uint64

	SetNonce(string, uint64)

	GetCodeHash(string) common.hash

	GetCode(string) []byte

	SetCode(string, []byte)

	GetCodeSize(string) int

	SetAbi(addr, abi string)

	GetAbi(addr string) string

	AddRefund(uint64)

	GetRefund() uint64

	GetState(string, common.Hash) common.Hash

	SetState(string, common.Hash, common.Hash)

	Suicide(string) bool

	HasSuicided(string) bool

	Exist(string) bool

	Empty(string) bool

	RevertToSnapshot(int)

	Snapshot() int

	TransferStateData(addr string)

	AddLog(*model.ContractLog)

	AddPreimage(common.Hash, []byte)


	CanTransfer(sender string, amount uint64) bool

	Transfer(sender, recipient string, amount uint64) bool

	GetBlockHeight() int64

	GetConfig() *types.Chain33Config
}
