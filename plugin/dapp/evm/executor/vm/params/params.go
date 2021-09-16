// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package params

import (
	"math/big"

	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/state"
)


type GasParam struct {
	Gas uint64

	Address common.Address
}

type EVMParam struct {

	StateDB state.EVMStateDB

	CallGasTemp uint64

	BlockNumber *big.Int
}
