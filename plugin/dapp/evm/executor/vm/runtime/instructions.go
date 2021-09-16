// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"fmt"
	"math/big"

	"github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common/crypto"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/model"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/params"
	"github.com/holiman/uint256"
)

func opAdd(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	y.Add(&x, y)
	return nil, nil
}

func opSub(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	y.Sub(&x, y)
	return nil, nil
}

func opMul(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	y.Mul(&x, y)
	return nil, nil
}

func opDiv(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	y.Div(&x, y)
	return nil, nil
}

func opSdiv(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	y.SDiv(&x, y)
	return nil, nil
}

func opMod(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	y.Mod(&x, y)
	return nil, nil
}

func opSmod(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	y.SMod(&x, y)
	return nil, nil
}

func opExp(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	base, exponent := callContext.stack.pop(), callContext.stack.peek()
	exponent.Exp(&base, exponent)
	return nil, nil
}

func opSignExtend(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	back, num := callContext.stack.pop(), callContext.stack.peek()
	num.ExtendSign(num, &back)
	return nil, nil
}

func opNot(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x := callContext.stack.peek()
	x.Not(x)
	return nil, nil
}

func opLt(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	if x.Lt(y) {
		y.SetOne()
	} else {
		y.Clear()
	}
	return nil, nil
}

func opGt(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	if x.Gt(y) {
		y.SetOne()
	} else {
		y.Clear()
	}
	return nil, nil
}

func opSlt(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	if x.Slt(y) {
		y.SetOne()
	} else {
		y.Clear()
	}
	return nil, nil
}

func opSgt(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	if x.Sgt(y) {
		y.SetOne()
	} else {
		y.Clear()
	}
	return nil, nil
}

func opEq(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	if x.Eq(y) {
		y.SetOne()
	} else {
		y.Clear()
	}
	return nil, nil
}

func opIszero(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x := callContext.stack.peek()
	if x.IsZero() {
		x.SetOne()
	} else {
		x.Clear()
	}
	return nil, nil
}

func opAnd(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	y.And(&x, y)
	return nil, nil
}

func opOr(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	y.Or(&x, y)
	return nil, nil
}

func opXor(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y := callContext.stack.pop(), callContext.stack.peek()
	y.Xor(&x, y)
	return nil, nil
}

func opByte(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	th, val := callContext.stack.pop(), callContext.stack.peek()
	val.Byte(&th)
	return nil, nil
}

func opAddmod(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y, z := callContext.stack.pop(), callContext.stack.pop(), callContext.stack.peek()
	if z.IsZero() {
		z.Clear()
	} else {
		z.AddMod(&x, &y, z)
	}
	return nil, nil
}

func opMulmod(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x, y, z := callContext.stack.pop(), callContext.stack.pop(), callContext.stack.peek()
	z.MulMod(&x, &y, z)
	return nil, nil
}

func opSHL(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := callContext.stack.pop(), callContext.stack.peek()
	if shift.LtUint64(256) {
		value.Lsh(value, uint(shift.Uint64()))
	} else {
		value.Clear()
	}

	return nil, nil
}

func opSHR(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := callContext.stack.pop(), callContext.stack.peek()
	if shift.LtUint64(256) {
		value.Rsh(value, uint(shift.Uint64()))
	} else {
		value.Clear()
	}

	return nil, nil
}

func opSAR(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	shift, value := callContext.stack.pop(), callContext.stack.peek()
	if shift.GtUint64(256) {
		if value.Sign() >= 0 {
			value.Clear()
		} else {
			// Max negative shift: all bits set
			value.SetAllOne()
		}
		return nil, nil
	}
	n := uint(shift.Uint64())
	value.SRsh(value, n)

	return nil, nil
}

func opSha3(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	offset, size := callContext.stack.pop(), callContext.stack.peek()
	data := callContext.memory.GetPtr(int64(offset.Uint64()), int64(size.Uint64()))
	hash := crypto.Keccak256(data)

	if evm.VMConfig.EnablePreimageRecording {
		evm.StateDB.AddPreimage(common.BytesToHash(hash), data)
	}
	size.SetBytes(hash)
	return nil, nil
}

func opAddress(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.push(new(uint256.Int).SetBytes(callContext.contract.Address().Bytes()))
	return nil, nil
}

func opBalance(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {

	slot := callContext.stack.peek()
	address := common.Uint256ToAddress(slot)
	slot.SetUint64(evm.StateDB.GetBalance(address.String()))
	return nil, nil
}

func opOrigin(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.push(new(uint256.Int).SetBytes(evm.Origin.Bytes()))
	return nil, nil
}

func opCaller(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.push(new(uint256.Int).SetBytes(callContext.contract.Caller().Bytes()))
	return nil, nil
}

func opCallValue(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.push(new(uint256.Int).SetUint64(callContext.contract.value))
	return nil, nil
}

func opCallDataLoad(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	x := callContext.stack.peek()
	if offset, overflow := x.Uint64WithOverflow(); !overflow {
		data := common.GetData(callContext.contract.Input, offset, 32)
		x.SetBytes(data)
	} else {
		x.Clear()
	}
	return nil, nil
}

func opCallDataSize(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.push(new(uint256.Int).SetUint64(uint64(len(callContext.contract.Input))))
	return nil, nil
}

func opCallDataCopy(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	var (
		memOffset  = callContext.stack.pop()
		dataOffset = callContext.stack.pop()
		length     = callContext.stack.pop()
	)
	dataOffset64, overflow := dataOffset.Uint64WithOverflow()
	if overflow {
		dataOffset64 = 0xffffffffffffffff
	}
	// These values are checked for overflow during gas cost calculation
	memOffset64 := memOffset.Uint64()
	length64 := length.Uint64()
	callContext.memory.Set(memOffset64, length64, common.GetData(callContext.contract.Input, dataOffset64, length64))
	return nil, nil
}

func opReturnDataSize(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.push(new(uint256.Int).SetUint64(uint64(len(evm.Interpreter.returnData))))
	return nil, nil
}

func opReturnDataCopy(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	var (
		memOffset  = callContext.stack.pop()
		dataOffset = callContext.stack.pop()
		length     = callContext.stack.pop()
	)

	offset64, overflow := dataOffset.Uint64WithOverflow()
	if overflow {
		return nil, model.ErrReturnDataOutOfBounds
	}
	// we can reuse dataOffset now (aliasing it for clarity)
	var end = dataOffset
	end.Add(&dataOffset, &length)
	end64, overflow := end.Uint64WithOverflow()
	if overflow || uint64(len(evm.Interpreter.returnData)) < end64 {
		return nil, model.ErrReturnDataOutOfBounds
	}
	callContext.memory.Set(memOffset.Uint64(), length.Uint64(), evm.Interpreter.returnData[offset64:end64])
	return nil, nil
}

func opExtCodeSize(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	slot := callContext.stack.peek()
	address := common.Uint256ToAddress(slot)
	slot.SetUint64(uint64(evm.StateDB.GetCodeSize(address.String())))
	return nil, nil
}

func opCodeSize(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	l := new(uint256.Int)
	l.SetUint64(uint64(len(callContext.contract.Code)))
	callContext.stack.push(l)
	return nil, nil
}

func opCodeCopy(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	var (
		memOffset  = callContext.stack.pop()
		codeOffset = callContext.stack.pop()
		length     = callContext.stack.pop()
	)
	uint64CodeOffset, overflow := codeOffset.Uint64WithOverflow()
	if overflow {
		uint64CodeOffset = 0xffffffffffffffff
	}
	codeCopy := common.GetData(callContext.contract.Code, uint64CodeOffset, length.Uint64())
	callContext.memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	return nil, nil
}

func opExtCodeCopy(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	var (
		stack      = callContext.stack
		a          = stack.pop()
		memOffset  = stack.pop()
		codeOffset = stack.pop()
		length     = stack.pop()
	)
	uint64CodeOffset, overflow := codeOffset.Uint64WithOverflow()
	if overflow {
		uint64CodeOffset = 0xffffffffffffffff
	}
	addr := common.Uint256ToAddress(&a)
	codeCopy := common.GetData(evm.StateDB.GetCode(addr.String()), uint64CodeOffset, length.Uint64())
	callContext.memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)
	return nil, nil
}

func opExtCodeHash(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	slot := callContext.stack.peek()
	addr := common.Uint256ToAddress(slot)
	if evm.StateDB.Empty(addr.String()) {
		slot.Clear()
	} else {
		slot.SetBytes(evm.StateDB.GetCodeHash(addr.String()).Bytes())
	}
	return nil, nil
}

func opGasprice(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	v := uint256.NewInt().SetUint64(uint64(evm.GasPrice))
	callContext.stack.push(v)
	return nil, nil
}

func opBlockhash(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	num := callContext.stack.peek()
	num64, overflow := num.Uint64WithOverflow()
	if overflow {
		num.Clear()
		return nil, nil
	}
	var upper, lower uint64
	upper = evm.BlockNumber.Uint64()
	if upper < 257 {
		lower = 0
	} else {
		lower = upper - 256
	}
	if num64 >= lower && num64 < upper {
		num.SetBytes(evm.GetHash(num64).Bytes())
	} else {
		num.Clear()
	}
	return nil, nil
}

func opCoinbase(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	if evm.Coinbase == nil {
		callContext.stack.push(new(uint256.Int).SetBytes(callContext.contract.CallerAddress.Bytes()))
	} else {
		callContext.stack.push(new(uint256.Int).SetBytes(evm.Coinbase.Bytes()))
	}
	return nil, nil
}

func opTimestamp(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	v, _ := uint256.FromBig(evm.Time)
	callContext.stack.push(v)
	return nil, nil
}

func opNumber(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	v, _ := uint256.FromBig(evm.BlockNumber)
	callContext.stack.push(v)
	return nil, nil
}

func opDifficulty(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	v, _ := uint256.FromBig(evm.Difficulty)
	callContext.stack.push(v)
	return nil, nil
}

func opGasLimit(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.push(new(uint256.Int).SetUint64(evm.GasLimit))
	return nil, nil
}

func opPop(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.pop()
	return nil, nil
}

func opMload(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	v := callContext.stack.peek()
	offset := int64(v.Uint64())
	v.SetBytes(callContext.memory.GetPtr(offset, 32))
	return nil, nil
}

func opMstore(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	// pop value of the stack
	mStart, val := callContext.stack.pop(), callContext.stack.pop()
	callContext.memory.Set32(mStart.Uint64(), &val)
	return nil, nil
}

func opMstore8(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	off, val := callContext.stack.pop(), callContext.stack.pop()
	callContext.memory.store[off.Uint64()] = byte(val.Uint64())
	return nil, nil
}

func opSload(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	loc := callContext.stack.peek()
	hash := common.BytesToHash(loc.Bytes())
	val := evm.StateDB.GetState(callContext.contract.Address().String(), hash)
	loc.SetBytes(val.Bytes())

	return nil, nil
}

func opSstore(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	loc := callContext.stack.pop()
	val := callContext.stack.pop()

	evm.StateDB.SetState(callContext.contract.Address().String(),
		common.BytesToHash(loc.Bytes()), common.BytesToHash(val.Bytes()))
	return nil, nil
}

func opJump(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	pos := callContext.stack.pop()
	contract := callContext.contract
	if !contract.Jumpdests.Has(contract.CodeHash, contract.Code, &pos) {
		nop := contract.GetOp(pos.Uint64())
		return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
	}
	*pc = pos.Uint64()
	return nil, nil
}

func opJumpi(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	pos, cond := callContext.stack.pop(), callContext.stack.pop()
	if !cond.IsZero() {
		if !callContext.contract.validJumpdest(&pos) {
			return nil, model.ErrInvalidJump
		}
		*pc = pos.Uint64()
	} else {
		*pc++
	}
	return nil, nil
}

func opJumpdest(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	return nil, nil
}

//func opBeginSub(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
//	return nil, model.ErrInvalidSubroutineEntry
//}

//func opJumpSub(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
//	if callContext.rstack.Len() >= 1023 {
//		return nil, model.ErrReturnStackExceeded
//	}
//	pos := callContext.stack.pop()
//	if !pos.IsUint64() {
//		return nil, model.ErrInvalidJump
//	}
//	posU64 := pos.Uint64()
//	if !callContext.contract.validJumpSubdest(posU64) {
//		return nil, model.ErrInvalidJump
//	}
//	callContext.rstack.push(uint32(*pc))
//	*pc = posU64 + 1
//	return nil, nil
//}
//func opReturnSub(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
//	if callContext.rstack.Len() == 0 {
//		return nil, model.ErrInvalidRetsub
//	}
//	// Other than the check that the return stack is not empty, there is no
//	// need to validate the pc from 'returns', since we only ever push valid
//	//values onto it via jumpsub.
//	*pc = uint64(callContext.rstack.pop()) + 1
//	return nil, nil
//}

func opPc(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.push(new(uint256.Int).SetUint64(*pc))
	return nil, nil
}

func opMsize(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.push(new(uint256.Int).SetUint64(uint64(callContext.memory.Len())))
	return nil, nil
}

func opGas(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	callContext.stack.push(new(uint256.Int).SetUint64(callContext.contract.Gas))
	return nil, nil
}

func opCreate(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {

	var (
		value        = callContext.stack.pop()
		offset, size = callContext.stack.pop(), callContext.stack.pop()
		input        = callContext.memory.GetCopy(int64(offset.Uint64()), int64(size.Uint64()))
		gas          = callContext.contract.Gas
	)

	callContext.contract.UseGas(gas)
	stackvalue := size

	addr := common.NewContractAddress(evm.Origin, evm.TxHash)
	res, _, returnGas, suberr := evm.Create(callContext.contract, addr, input, gas, "innerContract", "", value.Uint64())

	if suberr != nil && suberr != model.ErrCodeStoreOutOfGas {
		log15.Error("evm contract opCreate instruction error,value", suberr, value)
		stackvalue.Clear()
	} else {
		stackvalue.SetBytes(addr.Bytes())
	}
	callContext.stack.push(&stackvalue)

	callContext.contract.Gas += returnGas

	if suberr == model.ErrExecutionReverted {
		log15.Error("evm contract opCreate instruction error,value", suberr, value)
		return res, nil
	}
	return nil, nil
}

//CREATE2
func opCreate2(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	var (
		endowment    = callContext.stack.pop()
		offset, size = callContext.stack.pop(), callContext.stack.pop()
		salt         = callContext.stack.pop()
		input        = callContext.memory.GetCopy(int64(offset.Uint64()), int64(size.Uint64()))
		gas          = callContext.contract.Gas
	)

	// Apply EIP150
	gas = gas / 2
	callContext.contract.UseGas(gas)
	// reuse size int for stackvalue
	stackvalue := size
	//TODO: use uint256.Int instead of converting with toBig()
	bigEndowment := big0
	if !endowment.IsZero() {
		bigEndowment = endowment.ToBig()
	}

	//res, addr, returnGas, suberr := evm.Create2(callContext.contract, input, gas,
	//	bigEndowment, &salt)

	log15.Info("opCreate2", "initHash", common.Bytes2Hex(crypto.Keccak256Hash(input).Bytes()))

	newContractAddr := crypto.CreateAddress2(callContext.contract.Address(), salt.Bytes32(), crypto.Keccak256Hash(input).Bytes())
	log15.Info("opCreate2", "newContractAddr", newContractAddr.String(), "newContractAddr byte", common.Bytes2Hex(newContractAddr.Bytes()))
	saltSlice := salt.Bytes32()
	saltStr := common.Bytes2Hex(saltSlice[:])
	res, _, returnGas, suberr := evm.Create(callContext.contract, newContractAddr, input, gas, saltStr, "", endowment.Uint64())
	// push item on the stack based on the returned error.

	log15.Info("opCreate2", "callContext.contract.Address()", callContext.contract.Address(),
		"salt", saltStr, "suberr", suberr)
	if suberr != nil && suberr != model.ErrCodeStoreOutOfGas {
		log15.Error("evm contract opCreate instruction error,endowment,salt", suberr, bigEndowment, salt)
		stackvalue.Clear()
	} else {
		stackvalue.SetBytes(newContractAddr.Bytes())
	}
	if suberr != nil {
		stackvalue.Clear()
	} else {
		stackvalue.SetBytes(newContractAddr.Bytes())
	}
	callContext.stack.push(&stackvalue)
	callContext.contract.Gas += returnGas

	log15.Error("opCreate2", "returnGas", returnGas, "callContext.contract.Gas", callContext.contract.Gas)

	if suberr == ErrExecutionReverted {
		log15.Error("opCreate2", "\n\n suberr == ErrExecutionReverted \n\n")
		return res, nil
	}
	return nil, nil
}


func opCall(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	stack := callContext.stack
	// pop gas. The actual gas in interpreter.evm.callGasTemp.
	// We can use this as a temporary value
	temp := stack.pop()
	gas := evm.callGasTemp
	// pop other call parameters.
	addr, value, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := common.Uint256ToAddress(&addr)
	// Get the arguments from the memory.
	args := callContext.memory.GetPtr(int64(inOffset.Uint64()), int64(inSize.Uint64()))
	log15.Info("evm contract opCall", "toAddr", toAddr.String(), "value:", value.Uint64(), "input len", len(args), "input", common.Bytes2Hex(args))
	if !value.IsZero() {
		gas += params.CallStipend
	}
	ret, _, returnGas, err := evm.Call(callContext.contract, toAddr, args, gas, value.Uint64())
	if err != nil {
		temp.Clear()
		log15.Error("evm contract opCall instruction error", "error", err)
	} else {
		temp.SetOne()
	}
	stack.push(&temp)
	if err == nil || err == model.ErrExecutionReverted {
		callContext.memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	callContext.contract.Gas += returnGas
	return ret, nil
}

func opCallCode(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	// pop gas. The actual gas is in interpreter.evm.callGasTemp.
	stack := callContext.stack
	// We use it as a temporary value
	temp := stack.pop()
	gas := evm.callGasTemp
	// pop other call parameters.
	addr, value, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := common.Uint256ToAddress(&addr)
	// Get arguments from the memory.
	args := callContext.memory.GetPtr(int64(inOffset.Uint64()), int64(inSize.Uint64()))

	if !value.IsZero() {
		gas += params.CallStipend
	}
	ret, returnGas, err := evm.CallCode(callContext.contract, toAddr, args, gas, value.Uint64())
	if err != nil {
		temp.Clear()
		log15.Error("evm contract opCallCode instruction error", err)
	} else {
		temp.SetOne()
	}
	stack.push(&temp)
	if err == nil || err == model.ErrExecutionReverted {
		callContext.memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	callContext.contract.Gas += returnGas
	return ret, nil
}

func opDelegateCall(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	stack := callContext.stack
	// pop gas. The actual gas is in interpreter.evm.callGasTemp.
	// We use it as a temporary value
	temp := stack.pop()
	gas := evm.callGasTemp
	// pop other call parameters.
	addr, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := common.Uint256ToAddress(&addr)
	// Get arguments from the memory.
	args := callContext.memory.GetPtr(int64(inOffset.Uint64()), int64(inSize.Uint64()))

	ret, returnGas, err := evm.DelegateCall(callContext.contract, toAddr, args, gas)
	if err != nil {
		temp.Clear()
		log15.Error("evm contract opDelegateCall instruction error", err)
	} else {
		temp.SetOne()
	}
	stack.push(&temp)
	if err == nil || err == model.ErrExecutionReverted {
		callContext.memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	callContext.contract.Gas += returnGas
	return ret, nil
}

func opStaticCall(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	// pop gas. The actual gas is in interpreter.evm.callGasTemp.
	stack := callContext.stack
	// We use it as a temporary value
	temp := stack.pop()
	gas := evm.callGasTemp
	// pop other call parameters.
	addr, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := common.Uint256ToAddress(&addr)
	// Get arguments from the memory.
	args := callContext.memory.GetPtr(int64(inOffset.Uint64()), int64(inSize.Uint64()))

	log15.Info("opStaticCall", "caller", callContext.contract.self.Address(),
		"toaddr", toAddr, "gas", gas)
	ret, returnGas, err := evm.StaticCall(callContext.contract, toAddr, args, gas)
	if err != nil {
		temp.Clear()
		log15.Error("evm contract opStaticCall instruction error", "err", err)
	} else {
		temp.SetOne()
	}
	stack.push(&temp)
	if err == nil || err == model.ErrExecutionReverted {
		callContext.memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	callContext.contract.Gas += returnGas
	return ret, nil
}

func opReturn(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	offset, size := callContext.stack.pop(), callContext.stack.pop()
	ret := callContext.memory.GetPtr(int64(offset.Uint64()), int64(size.Uint64()))
	return ret, nil
}

func opRevert(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	offset, size := callContext.stack.pop(), callContext.stack.pop()
	ret := callContext.memory.GetPtr(int64(offset.Uint64()), int64(size.Uint64()))
	log15.Info("opRevert", "info", string(ret))
	return ret, nil
}

func opStop(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	return nil, nil
}

func opSuicide(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	beneficiary := callContext.stack.pop()
	balance := evm.StateDB.GetBalance(callContext.contract.Address().String())

	evm.StateDB.AddBalance(common.Uint256ToAddress(&beneficiary).String(), (*callContext.contract.CodeAddr).String(), balance)
	evm.StateDB.Suicide(callContext.contract.Address().String())
	return nil, nil
}

func makeLog(size int) executionFunc {
	return func(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
		topics := make([]common.Hash, size)
		stack := callContext.stack
		mStart, mSize := stack.pop(), stack.pop()
		for i := 0; i < size; i++ {
			addr := stack.pop()
			topics[i] = common.Uint256ToHash(&addr)
		}

		d := callContext.memory.GetCopy(int64(mStart.Uint64()), int64(mSize.Uint64()))
		evm.StateDB.AddLog(&model.ContractLog{
			Address: callContext.contract.Address(),
			Topics:  topics,
			Data:    d,
			// This is a non-consensus field, but assigned here because
			// core/state doesn't know the current block number.
			BlockNumber: evm.BlockNumber.Uint64(),
		})
		log15.Info("makeLog End", "data", string(d), "data in hex", common.Bytes2Hex(d))
		return nil, nil
	}
}

// opPush1 is a specialized version of pushN
func opPush1(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	var (
		codeLen = uint64(len(callContext.contract.Code))
		integer = new(uint256.Int)
	)
	*pc++
	if *pc < codeLen {
		callContext.stack.push(integer.SetUint64(uint64(callContext.contract.Code[*pc])))
	} else {
		callContext.stack.push(integer.Clear())
	}
	return nil, nil
}

func makePush(size uint64, pushByteSize int) executionFunc {
	return func(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
		codeLen := len(callContext.contract.Code)

		startMin := codeLen
		if int(*pc+1) < startMin {
			startMin = int(*pc + 1)
		}

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			endMin = startMin + pushByteSize
		}

		integer := new(uint256.Int)
		callContext.stack.push(integer.SetBytes(common.RightPadBytes(
			callContext.contract.Code[startMin:endMin], pushByteSize)))

		*pc += size
		return nil, nil
	}
}

func makeDup(size int64) executionFunc {
	return func(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
		callContext.stack.dup(int(size))
		return nil, nil
	}
}

func makeSwap(size int64) executionFunc {
	// switch n + 1 otherwise n would be swapped with n
	size++
	return func(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
		callContext.stack.swap(int(size))
		return nil, nil
	}
}

func opSelfBalance(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	balance := uint256.NewInt().SetUint64(evm.StateDB.GetBalance(callContext.contract.Address().String()))
	callContext.stack.push(balance)
	return nil, nil
}

// opChainID implements CHAINID opcode
func opChainID(pc *uint64, evm *EVM, callContext *callCtx) ([]byte, error) {
	chainId, _ := uint256.FromBig(big.NewInt(int64(evm.cfg.GetChainID())))
	callContext.stack.push(chainId)
	return nil, nil
}
