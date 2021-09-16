// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gas

import (
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/mm"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/model"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/params"
)

type (
	CalcGasFunc func(Table, *params.EVMParam, *params.GasParam, *mm.Stack, *mm.Memory, uint64) (uint64, error) // last parameter is the requested memory size as a uint64
)

type Table struct {
	ExtcodeSize uint64
	ExtcodeCopy uint64
	Balance uint64
	SLoad uint64
	Calls uint64
	Suicide uint64
	ExpByte uint64
}

var (
	TableHomestead = Table{
		ExtcodeSize: 20,
		ExtcodeCopy: 20,
		Balance:     20,
		SLoad:       50,
		Calls:       40,
		Suicide:     0,
		ExpByte:     10,
	}
)

func memoryGasCost(mem *mm.Memory, newMemSize uint64) (uint64, error) {
	if newMemSize == 0 {
		return 0, nil
	}

	if newMemSize > MaxNewMemSize {
		return 0, model.ErrGasUintOverflow
	}

	newMemSizeWords := common.ToWordSize(newMemSize)

	newMemSize = newMemSizeWords * 32

	if newMemSize > uint64(mem.Len()) {
		square := newMemSizeWords * newMemSizeWords
		linCoef := newMemSizeWords * params.MemoryGas
		quadCoef := square / params.QuadCoeffDiv
		newTotalFee := linCoef + quadCoef

		fee := newTotalFee - mem.LastGasCost
		mem.LastGasCost = newTotalFee

		return fee, nil
	}
	return 0, nil
}

func ConstGasFunc(gas uint64) CalcGasFunc {
	return func(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
		return gas, nil
	}
}

func CallDataCopy(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}

	var overflow bool
	if gas, overflow = common.SafeAdd(gas, GasFastestStep); overflow {
		return 0, model.ErrGasUintOverflow
	}

	words, overflow := stack.Back(2).Uint64WithOverflow()
	if overflow {
		return 0, model.ErrGasUintOverflow
	}

	if words, overflow = common.SafeMul(common.ToWordSize(words), params.CopyGas); overflow {
		return 0, model.ErrGasUintOverflow
	}

	if gas, overflow = common.SafeAdd(gas, words); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func ReturnDataCopy(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}

	var overflow bool
	if gas, overflow = common.SafeAdd(gas, GasFastestStep); overflow {
		return 0, model.ErrGasUintOverflow
	}

	words, overflow := stack.Back(2).Uint64WithOverflow()
	if overflow {
		return 0, model.ErrGasUintOverflow
	}

	if words, overflow = common.SafeMul(common.ToWordSize(words), params.CopyGas); overflow {
		return 0, model.ErrGasUintOverflow
	}

	if gas, overflow = common.SafeAdd(gas, words); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func SStore(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	var (
		y, x = stack.Back(1), stack.Back(0)
		val  = evm.StateDB.GetState(contractGas.Address.String(), common.Uint256ToHash(x))
	)

	if val == (common.Hash{}) && y.Sign() != 0 {

		// 0 => non 0
		return params.SstoreSetGas, nil
	} else if val != (common.Hash{}) && y.Sign() == 0 {

		// non 0 => 0
		evm.StateDB.AddRefund(params.SstoreRefundGas)
		return params.SstoreClearGas, nil
	} else {

		// non 0 => non 0
		return params.SstoreResetGas, nil
	}
}


func MakeGasLog(n uint64) CalcGasFunc {
	return func(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
		requestedSize, overflow := stack.Back(1).Uint64WithOverflow()
		if overflow {
			return 0, model.ErrGasUintOverflow
		}

		gas, err := memoryGasCost(mem, memorySize)
		if err != nil {
			return 0, err
		}

		if gas, overflow = common.SafeAdd(gas, params.LogGas); overflow {
			return 0, model.ErrGasUintOverflow
		}
		if gas, overflow = common.SafeAdd(gas, n*params.LogTopicGas); overflow {
			return 0, model.ErrGasUintOverflow
		}

		var memorySizeGas uint64
		if memorySizeGas, overflow = common.SafeMul(requestedSize, params.LogDataGas); overflow {
			return 0, model.ErrGasUintOverflow
		}
		if gas, overflow = common.SafeAdd(gas, memorySizeGas); overflow {
			return 0, model.ErrGasUintOverflow
		}
		return gas, nil
	}
}

func Sha3(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}

	if gas, overflow = common.SafeAdd(gas, params.Sha3Gas); overflow {
		return 0, model.ErrGasUintOverflow
	}

	wordGas, overflow := stack.Back(1).Uint64WithOverflow()
	if overflow {
		return 0, model.ErrGasUintOverflow
	}
	if wordGas, overflow = common.SafeMul(common.ToWordSize(wordGas), params.Sha3WordGas); overflow {
		return 0, model.ErrGasUintOverflow
	}
	if gas, overflow = common.SafeAdd(gas, wordGas); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func CodeCopy(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}

	var overflow bool
	if gas, overflow = common.SafeAdd(gas, GasFastestStep); overflow {
		return 0, model.ErrGasUintOverflow
	}

	wordGas, overflow := stack.Back(2).Uint64WithOverflow()
	if overflow {
		return 0, model.ErrGasUintOverflow
	}
	if wordGas, overflow = common.SafeMul(common.ToWordSize(wordGas), params.CopyGas); overflow {
		return 0, model.ErrGasUintOverflow
	}
	if gas, overflow = common.SafeAdd(gas, wordGas); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func ExtCodeCopy(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}

	var overflow bool
	if gas, overflow = common.SafeAdd(gas, gt.ExtcodeCopy); overflow {
		return 0, model.ErrGasUintOverflow
	}

	wordGas, overflow := stack.Back(3).Uint64WithOverflow()
	if overflow {
		return 0, model.ErrGasUintOverflow
	}

	if wordGas, overflow = common.SafeMul(common.ToWordSize(wordGas), params.CopyGas); overflow {
		return 0, model.ErrGasUintOverflow
	}

	if gas, overflow = common.SafeAdd(gas, wordGas); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func MLoad(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, model.ErrGasUintOverflow
	}
	if gas, overflow = common.SafeAdd(gas, GasFastestStep); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func MStore8(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, model.ErrGasUintOverflow
	}
	if gas, overflow = common.SafeAdd(gas, GasFastestStep); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func MStore(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, model.ErrGasUintOverflow
	}
	if gas, overflow = common.SafeAdd(gas, GasFastestStep); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func Create(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	if gas, overflow = common.SafeAdd(gas, params.CreateGas); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func Balance(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	return gt.Balance, nil
}

func ExtCodeSize(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	return gt.ExtcodeSize, nil
}

func SLoad(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	return gt.SLoad, nil
}

func Exp(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	expByteLen := uint64((stack.Data()[stack.Len()-2].BitLen() + 7) / 8)

	var (
		gas      = expByteLen * gt.ExpByte // no overflow check required. Max is 256 * ExpByte gas
		overflow bool
	)
	if gas, overflow = common.SafeAdd(gas, GasSlowStep); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func Call(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	var (
		gas            = gt.Calls
		transfersValue = stack.Back(2).Sign() != 0
		address        = common.Uint256ToAddress(stack.Back(1))
	)
	if !evm.StateDB.Exist(address.String()) {
		gas += params.CallNewAccountGas
	}
	if transfersValue {
		gas += params.CallValueTransferGas
	}
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	var overflow bool
	if gas, overflow = common.SafeAdd(gas, memoryGas); overflow {
		return 0, model.ErrGasUintOverflow
	}

	evm.CallGasTemp, err = callGas(gt, contractGas.Gas, gas, stack.Back(0))
	if err != nil {
		return 0, err
	}
	if gas, overflow = common.SafeAdd(gas, evm.CallGasTemp); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func CallCode(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	gas := gt.Calls
	if stack.Back(2).Sign() != 0 {
		gas += params.CallValueTransferGas
	}
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	var overflow bool
	if gas, overflow = common.SafeAdd(gas, memoryGas); overflow {
		return 0, model.ErrGasUintOverflow
	}

	evm.CallGasTemp, err = callGas(gt, contractGas.Gas, gas, stack.Back(0))
	if err != nil {
		return 0, err
	}
	if gas, overflow = common.SafeAdd(gas, evm.CallGasTemp); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func Return(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	return memoryGasCost(mem, memorySize)
}

func Revert(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	return memoryGasCost(mem, memorySize)
}

func Suicide(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	var gas uint64
	if !evm.StateDB.HasSuicided(contractGas.Address.String()) {
		evm.StateDB.AddRefund(params.SelfdestructRefundGas)
	}
	return gas, nil
}

func DelegateCall(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	var overflow bool
	if gas, overflow = common.SafeAdd(gas, gt.Calls); overflow {
		return 0, model.ErrGasUintOverflow
	}

	evm.CallGasTemp, err = callGas(gt, contractGas.Gas, gas, stack.Back(0))
	if err != nil {
		return 0, err
	}
	if gas, overflow = common.SafeAdd(gas, evm.CallGasTemp); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func StaticCall(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	var overflow bool
	if gas, overflow = common.SafeAdd(gas, gt.Calls); overflow {
		return 0, model.ErrGasUintOverflow
	}

	evm.CallGasTemp, err = callGas(gt, contractGas.Gas, gas, stack.Back(0))
	if err != nil {
		return 0, err
	}
	if gas, overflow = common.SafeAdd(gas, evm.CallGasTemp); overflow {
		return 0, model.ErrGasUintOverflow
	}
	return gas, nil
}

func Push(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	return GasFastestStep, nil
}

func Swap(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	return GasFastestStep, nil
}

func Dup(gt Table, evm *params.EVMParam, contractGas *params.GasParam, stack *mm.Stack, mem *mm.Memory, memorySize uint64) (uint64, error) {
	return GasFastestStep, nil
}
