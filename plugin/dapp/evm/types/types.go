// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"reflect"

	"github.com/33cn/chain33/types"
)

const (

	EvmCreateAction = 1

	EvmCallAction = 2

	TyLogContractData = 601

	TyLogContractState = 602

	TyLogCallContract = 603

	TyLogEVMStateChangeItem = 604

	TyLogEVMEventData = 605

	MaxGasLimit = (100000000 * 5)
)

const (
	EVMEnable = "Enable"

	ForkEVMState = "ForkEVMState"
	ForkEVMKVHash = "ForkEVMKVHash"
	ForkEVMABI = "ForkEVMABI"

	ForkEVMFrozen = "ForkEVMFrozen"

	ForkEVMYoloV1 = "ForkEVMYoloV1"

	ForkEVMTxGroup = "ForkEVMTxGroup"
)

var (

	EvmPrefix = "user.evm."

	ExecutorName = "evm"

	ExecerEvm = []byte(ExecutorName)

	UserPrefix = []byte(EvmPrefix)

	logInfo = map[int64]*types.LogInfo{
		TyLogCallContract:       {Ty: reflect.TypeOf(ReceiptEVMContract{}), Name: "LogCallContract"},
		TyLogContractData:       {Ty: reflect.TypeOf(EVMContractData{}), Name: "LogContractData"},
		TyLogContractState:      {Ty: reflect.TypeOf(EVMContractState{}), Name: "LogContractState"},
		TyLogEVMStateChangeItem: {Ty: reflect.TypeOf(EVMStateChangeItem{}), Name: "LogEVMStateChangeItem"},
		TyLogEVMEventData:       {Ty: reflect.TypeOf(types.EVMLog{}), Name: "LogEVMEventData"},
	}
)
