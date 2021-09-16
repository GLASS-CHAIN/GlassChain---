// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"strings"

	"github.com/33cn/chain33/common/address"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
)

var (
	elog = log.New("module", "exectype.evm")

	actionName = map[string]int32{
		"EvmCreate": EvmCreateAction,
		"EvmCall":   EvmCallAction,
	}
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, ExecerEvm)
	types.RegFork(ExecutorName, InitFork)
	types.RegExec(ExecutorName, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.Chain33Config) {
	cfg.RegisterDappFork(ExecutorName, EVMEnable, 500000)

	cfg.RegisterDappFork(ExecutorName, ForkEVMState, 650000)

	cfg.RegisterDappFork(ExecutorName, ForkEVMKVHash, 1000000)

	cfg.RegisterDappFork(ExecutorName, ForkEVMABI, 1250000)

	cfg.RegisterDappFork(ExecutorName, ForkEVMFrozen, 1300000)

	cfg.RegisterDappFork(ExecutorName, ForkEVMYoloV1, 9500000)

	cfg.RegisterDappFork(ExecutorName, ForkEVMTxGroup, 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(ExecutorName, NewType(cfg))
}

type EvmType struct {
	types.ExecTypeBase
}

func NewType(cfg *types.Chain33Config) *EvmType {
	c := &EvmType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

func (evm *EvmType) GetName() string {
	return ExecutorName
}

func (evm *EvmType) GetPayload() types.Message {
	return &EVMContractAction{}
}

func (evm EvmType) ActionName(tx *types.Transaction) string {
	cfg := evm.GetConfig()
	if strings.EqualFold(tx.To, address.ExecAddress(cfg.ExecName(ExecutorName))) {
		return "createEvmContract"
	}
	return "callEvmContract"
}

func (evm *EvmType) GetTypeMap() map[string]int32 {
	return actionName
}

func (evm EvmType) GetRealToAddr(tx *types.Transaction) string {
	if string(tx.Execer) == ExecutorName {
		return tx.To
	}
	var action EVMContractAction
	err := types.Decode(tx.Payload, &action)
	if err != nil {
		return tx.To
	}
	return tx.To
}

func (evm EvmType) Amount(tx *types.Transaction) (int64, error) {
	return 0, nil
}

func (evm *EvmType) GetLogMap() map[int64]*types.LogInfo {
	return logInfo
}
