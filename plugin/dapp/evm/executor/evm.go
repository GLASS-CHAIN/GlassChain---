// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"bytes"
	"math/big"
	"os"

	"reflect"

	"github.com/33cn/chain33/common/address"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/runtime"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/state"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
)

var (
	evmDebugInited = false
	EvmAddress = ""
	driverName = evmtypes.ExecutorName
)

func Init(name string, cfg *types.Chain33Config, sub []byte) {
	driverName = name
	drivers.Register(cfg, driverName, newEVMDriver, cfg.GetDappFork(driverName, evmtypes.EVMEnable))
	EvmAddress = address.ExecAddress(cfg.ExecName(name))

	state.InitForkData()
	InitExecType()
}

// InitExecType Init Exec Type
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&EVMExecutor{}))
}

func GetName() string {
	return newEVMDriver().GetName()
}

func newEVMDriver() drivers.Driver {
	evm := NewEVMExecutor()
	return evm
}

type EVMExecutor struct {
	drivers.DriverBase
	vmCfg    *runtime.Config
	mStateDB *state.MemoryStateDB
}

func NewEVMExecutor() *EVMExecutor {
	exec := &EVMExecutor{}

	exec.vmCfg = &runtime.Config{}
	//exec.vmCfg.Tracer = runtime.NewJSONLogger(os.Stdout)
	exec.vmCfg.Tracer = runtime.NewMarkdownLogger(
		&runtime.LogConfig{
			DisableMemory:     false,
			DisableStack:      false,
			DisableStorage:    false,
			DisableReturnData: false,
			Debug:             true,
			Limit:             0,
		},
		os.Stdout,
	)

	exec.SetChild(exec)
	exec.SetExecutorType(types.LoadExecutorType(driverName))
	return exec
}

func (evm *EVMExecutor) GetFuncMap() map[string]reflect.Method {
	ety := types.LoadExecutorType(driverName)
	return ety.GetExecFuncMap()
}

func (evm *EVMExecutor) GetDriverName() string {
	return evmtypes.ExecutorName
}

func (evm *EVMExecutor) ExecutorOrder() int64 {
	cfg := evm.GetAPI().GetConfig()
	if cfg.IsFork(evm.GetHeight(), "ForkLocalDBAccess") {
		return drivers.ExecLocalSameTime
	}
	return evm.DriverBase.ExecutorOrder()
}

func (evm *EVMExecutor) Allow(tx *types.Transaction, index int) error {
	err := evm.DriverBase.Allow(tx, index)
	if err == nil {
		return nil
	}

	cfg := evm.GetAPI().GetConfig()
	exec := cfg.GetParaExec(tx.Execer)
	if evm.AllowIsUserDot2(exec) {
		return nil
	}
	return types.ErrNotAllow
}

func (evm *EVMExecutor) IsFriend(myexec, writekey []byte, othertx *types.Transaction) bool {
	if othertx == nil {
		return false
	}
	cfg := evm.GetAPI().GetConfig()
	exec := cfg.GetParaExec(othertx.Execer)
	if exec == nil || len(bytes.TrimSpace(exec)) == 0 {
		return false
	}
	if bytes.HasPrefix(exec, evmtypes.UserPrefix) || bytes.Equal(exec, evmtypes.ExecerEvm) {
		if bytes.HasPrefix(writekey, []byte("mavl-evm-")) {
			return true
		}
	}
	return false
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (evm *EVMExecutor) CheckReceiptExecOk() bool {
	return true
}

func (evm *EVMExecutor) getNewAddr(txHash []byte) common.Address {
	cfg := evm.GetAPI().GetConfig()
	return common.NewAddress(cfg, txHash)
}

// createContractAddress creates an ethereum address given the bytes and the nonce
func (evm *EVMExecutor) createContractAddress(b common.Address, txHash []byte) common.Address {
	return common.NewContractAddress(b, txHash)
}


func (evm *EVMExecutor) CheckTx(tx *types.Transaction, index int) error {
	return nil
}

func (evm *EVMExecutor) GetActionName(tx *types.Transaction) string {
	cfg := evm.GetAPI().GetConfig()
	if bytes.Equal(tx.Execer, []byte(cfg.ExecName(evmtypes.ExecutorName))) {
		return cfg.ExecName(evmtypes.ExecutorName)
	}
	return tx.ActionName()
}

func (evm *EVMExecutor) GetMStateDB() *state.MemoryStateDB {
	return evm.mStateDB
}

func (evm *EVMExecutor) GetVMConfig() *runtime.Config {
	return evm.vmCfg
}

func (evm *EVMExecutor) NewEVMContext(msg *common.Message, txHash []byte) runtime.Context {
	return runtime.Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(evm.GetAPI()),
		Origin:      msg.From(),
		Coinbase:    nil,
		BlockNumber: new(big.Int).SetInt64(evm.GetHeight()),
		Time:        new(big.Int).SetInt64(evm.GetBlockTime()),
		Difficulty:  new(big.Int).SetUint64(evm.GetDifficulty()),
		GasLimit:    msg.GasLimit(),
		GasPrice:    msg.GasPrice(),
		TxHash:      txHash,
	}
}
