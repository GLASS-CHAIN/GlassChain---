// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"math/big"
	"sync/atomic"

	log "github.com/33cn/chain33/common/log/log15"

	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/gas"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/model"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/params"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/state"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
)

type (
	CanTransferFunc func(state.EVMStateDB, common.Address, uint64) bool

	TransferFunc func(state.EVMStateDB, common.Address, common.Address, uint64) bool


	GetHashFunc func(uint64) common.Hash
)

func run(evm *EVM, contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {
	if contract.CodeAddr != nil {
		precompiles := PrecompiledContractsBerlin
		if p := precompiles[contract.CodeAddr.ToHash160()]; p != nil {
			ret, contract.Gas, err = RunPrecompiledContract(p, input, contract.Gas)
			return
		}
	}
	ret, err = evm.Interpreter.Run(contract, input, readOnly)
	if err != nil {
		log.Error("error occurs while run evm contract", "error info", err)
	}

	return ret, err
}

type Context struct {

	CanTransfer CanTransferFunc

	Transfer TransferFunc

	GetHash GetHashFunc

	Origin common.Address
	GasPrice uint32

	Coinbase *common.address
	GasLimit uint64

	// TxHash
	TxHash []byte

	BlockNumber *big.Int

	Time *big.Int

	Difficulty *big.Int
}


type EVM struct {

	context

	StateDB state.EVMStateDB
	
	depth int

	VMConfig Config

	Interpreter *Interpreter


	abort int32


	callGasTemp uint64

	maxCodeSize int

	cfg *types.Chain33Config
}

func NewEVM(ctx Context, statedb state.EVMStateDB, vmConfig Config, cfg *types.Chain33Config) *EVM {
	evm := &EVM{
		Context:     ctx,
		StateDB:     statedb,
		VMConfig:    vmConfig,
		maxCodeSize: params.MaxCodeSize,
		cfg:         cfg,
	}

	evm.Interpreter = NewInterpreter(evm, vmConfig)
	return evm
}

func (evm *EVM) GasTable(num *big.Int) gas.Table {
	return gas.TableHomestead
}

func (evm *EVM) Cancel() {
	atomic.StoreInt32(&evm.abort, 1)
}

func (evm *EVM) SetMaxCodeSize(maxCodeSize int) {
	if maxCodeSize < 1 || maxCodeSize > params.MaxCodeSize {
		return
	}

	evm.maxCodeSize = maxCodeSize
}

func (evm *EVM) preCheck(caller ContractRef, value uint64) (pass bool, err error) {

	if evm.VMConfig.NoRecursion && evm.depth > 0 {
		return false, nil
	}

	if evm.depth > int(params.CallCreateDepth) {
		return false, model.ErrDepth
	}

	if value > 0 {
		if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
			return false, model.ErrInsufficientBalance
		}
	}

	return true, nil
}


func (evm *EVM) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value uint64) (ret []byte, snapshot int, leftOverGas uint64, err error) {
	pass, err := evm.preCheck(caller, value)
	if !pass {
		return nil, -1, gas, err
	}

	p, isPrecompile := evm.precompile(addr)
	if !evm.StateDB.Exist(addr.String()) {
		if !isPrecompile {
			if len(input) > 0 || value == 0 {
				if EVMDebugOn == evm.VMConfig.Debug && evm.depth == 0 {
					evm.VMConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)
					evm.VMConfig.Tracer.CaptureEnd(ret, 0, 0, nil)
				}
				return nil, -1, gas, model.ErrAddrNotExists
			}
		} else {

		}
	}

	if evm.StateDB.HasSuicided(addr.String()) {
		return nil, -1, gas, model.ErrDestruct
	}

	snapshot = evm.StateDB.Snapshot()
	to := AccountRef(addr)

	evm.Transfer(evm.StateDB, caller.Address(), to.Address(), value)
	log.Info("evm call", "caller address", caller.Address().String(), "contract address", to.Address().String(), "value", value)

	cfg := evm.StateDB.GetConfig()
	if cfg.IsDappFork(evm.BlockNumber.Int64(), "evm", evmtypes.ForkEVMState) {
		evm.StateDB.TransferStateData(addr.String())
	}

	if isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		code := evm.StateDB.GetCode(addr.String())
		if len(code) == 0 {
			ret, err = nil, nil // gas is unchanged
		} else {
			contract := NewContract(caller, to, value, gas)
			contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr.String()), evm.StateDB.GetCode(addr.String()))

			start := types.Now()

			if EVMDebugOn == evm.VMConfig.Debug && evm.depth == 0 {
				evm.VMConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)

				defer func() {
					evm.VMConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, types.Since(start), err)
				}()
			}
			ret, err = run(evm, contract, input, false)
			gas = contract.Gas
		}
	}

	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != model.ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, snapshot, gas, err
}

func (evm *EVM) CallCode(caller ContractRef, addr common.Address, input []byte, gas uint64, value uint64) (ret []byte, leftOverGas uint64, err error) {
	pass, err := evm.preCheck(caller, value)
	if !pass {
		return nil, gas, err
	}

	if evm.StateDB.HasSuicided(addr.String()) {
		return nil, gas, model.ErrDestruct
	}

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	// It is allowed to call precompiles, even via delegatecall
	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		contract := NewContract(caller, to, value, gas)
		contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr.String()), evm.StateDB.GetCode(addr.String()))
		ret, err = run(evm, contract, input, false)
		gas = contract.Gas

	}

	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

func (evm *EVM) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	pass, err := evm.preCheck(caller, 0)
	if !pass {
		return nil, gas, err
	}

	if evm.StateDB.HasSuicided(addr.String()) {
		return nil, gas, model.ErrDestruct
	}

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		contract := NewContract(caller, to, 0, gas).AsDelegate()
		contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr.String()), evm.StateDB.GetCode(addr.String()))

		ret, err = run(evm, contract, input, true)

	}

	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != model.ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

func (evm *EVM) StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	addrecrecover := common.BytesToAddress(common.RightPadBytes([]byte{1}, 20))
	log.Info("StaticCall", "input", common.Bytes2Hex(input),
		"addr slice", common.Bytes2Hex(addr.Bytes()),
		"addrecrecover", addrecrecover.String(),
		"addrecrecoverslice", common.Bytes2Hex(addrecrecover.Bytes()))

	log.Info("StaticCall contract info", "caller", caller.Address(), "gas", gas)

	pass, err := evm.preCheck(caller, 0)
	if !pass {
		return nil, gas, err
	}

	isPrecompile := false
	precompiles := PrecompiledContractsByzantium
	if !evm.StateDB.Exist(addr.String()) {

		if evm.cfg.IsDappFork(evm.StateDB.GetBlockHeight(), "evm", evmtypes.ForkEVMYoloV1) {
			precompiles = PrecompiledContractsIstanbul
		}
		if precompiles[addr.ToHash160()] == nil {
			if len(input) > 0 {

				return nil, gas, model.ErrAddrNotExists
			}
		} else {
			isPrecompile = true
			log.Info("StaticCall", "addr.Bytes()", common.Bytes2Hex(addr.Bytes()),
				"isPrecompile", isPrecompile)
		}
	}

	log.Info("StaticCall debug", "hhhh", 1)

	if evm.StateDB.HasSuicided(addr.String()) {
		return nil, gas, model.ErrDestruct
	}

	if !evm.Interpreter.readOnly {
		evm.Interpreter.readOnly = true
		defer func() { evm.Interpreter.readOnly = false }()
	}

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)

	contract := NewContract(caller, to, 0, gas)
	if isPrecompile {
		ret, gas, err = RunPrecompiledContract(precompiles[addr.ToHash160()], input, gas)
	} else {
		contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr.String()), evm.StateDB.GetCode(addr.String()))
		ret, err = run(evm, contract, input, false)
	}

	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)

		if err != model.ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

func (evm *EVM) Create(caller ContractRef, contractAddr common.Address, code []byte, gas uint64, execName, alias string, value uint64) (ret []byte, snapshot int, leftOverGas uint64, err error) {
	pass, err := evm.preCheck(caller, value)
	if !pass {
		return nil, -1, gas, err
	}

	evm.Transfer(evm.StateDB, caller.Address(), contractAddr, value)

	contract := NewContract(caller, AccountRef(contractAddr), value, gas)
	contract.SetCallCode(&contractAddr, common.ToHash(code), code)

	snapshot = evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(contractAddr.String(), contract.CallerAddress.String(), execName, alias)

	if EVMDebugOn == evm.VMConfig.Debug && evm.depth == 0 {
		evm.VMConfig.Tracer.CaptureStart(caller.Address(), contractAddr, true, code, gas, 0)
	}
	start := types.Now()

	ret, err = run(evm, contract, nil, false)

	maxCodeSizeExceeded := len(ret) > evm.maxCodeSize

	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			evm.StateDB.SetCode(contractAddr.String(), ret)
		} else {
			err = model.ErrCodeStoreOutOfGas
		}
	}


	if maxCodeSizeExceeded || (err != nil && err != model.ErrCodeStoreOutOfGas) {
		evm.StateDB.RevertToSnapshot(snapshot)

		if err != model.ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	if maxCodeSizeExceeded && err == nil {
		err = model.ErrMaxCodeSizeExceeded
	}

	if EVMDebugOn == evm.VMConfig.Debug && evm.depth == 0 {
		evm.VMConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, types.Since(start), err)
	}

	return ret, snapshot, contract.Gas, err
}

func (evm *EVM) precompile(addr common.Address) (PrecompiledContract, bool) {
	p, ok := PrecompiledContractsBerlin[addr.ToHash160()]
	return p, ok
}
