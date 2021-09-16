// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"sync/atomic"

	"github.com/33cn/chain33/common/log/log15"

	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"

	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common/math"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/model"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/params"
)

type Config struct {
	Debug int32
	Tracer Tracer
	NoRecursion bool
	EnablePreimageRecording bool
	JumpTable [256]*operation
}

type Interpreter struct {
	evm *EVM
	cfg Config
	readOnly bool
	returnData []byte
}

const (
	EVMDebugOn  = int32(1)
	EVMDebugOff = int32(0)
)

func NewInterpreter(evm *EVM, cfg Config) *Interpreter {
	if cfg.JumpTable[STOP] == nil {
		cfg.JumpTable = berlinInstructionSet
		if evm.cfg.IsDappFork(evm.StateDB.GetBlockHeight(), "evm", evmtypes.ForkEVMYoloV1) {

			cfg.JumpTable = berlinInstructionSet
		}
	}

	return &Interpreter{
		evm: evm,
		cfg: cfg,
	}
}

func (in *Interpreter) enforceRestrictions(op OpCode, operation *operation, stack *Stack) error {
	if in.readOnly {
		if operation.writes || (op == CALL && stack.Back(2).BitLen() > 0) {
			return model.ErrWriteProtection
		}
	}
	return nil
}

func (in *Interpreter) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {

	in.evm.depth++
	defer func() { in.evm.depth-- }()

	// Make sure the readOnly is only set if we aren't in readOnly yet.
	// This makes also sure that the readOnly flag isn't removed for child calls.
	if readOnly && !in.readOnly {
		in.readOnly = true
		defer func() { in.readOnly = false }()
	}

	in.returnData = nil

	if len(contract.Code) == 0 {
		return nil, nil
	}

	var (
		op OpCode
		mem = NewMemory()
		stack = newstack()
		//returns     = mm.NewReturnStack() // local returns stack
		callContext = &callCtx{
			memory:   mem,
			stack:    stack,
			contract: contract,
		}

		pc = uint64(0)

		cost uint64
		pcCopy  uint64
		gasCopy uint64
		logged  bool
		res []byte
	)
	contract.Input = input

	defer func() {
		returnStack(stack)
	}()

	if EVMDebugOn == in.cfg.Debug {
		defer func() {
			if err != nil {
				if !logged {
					in.cfg.Tracer.CaptureState(in.evm, pcCopy, op, gasCopy, cost, mem, stack, in.returnData, contract, in.evm.depth, err)
				} else {
					in.cfg.Tracer.CaptureFault(in.evm, pcCopy, op, gasCopy, cost, mem, stack, contract, in.evm.depth, err)
				}
			}
		}()
	}
	steps := 0
	for {
		steps++
		if steps%1000 == 0 && atomic.LoadInt32(&in.evm.abort) != 0 {
			break
		}
		if EVMDebugOn == in.cfg.Debug {
			logged, pcCopy, gasCopy = false, pc, contract.Gas
		}

		op = contract.GetOp(pc)
		operation := in.cfg.JumpTable[op]
		if operation == nil {
			log15.Error("can't found operation:%s", op)
			return nil, &ErrInvalidOpCode{opcode: op}
		}
		// Validate stack
		if sLen := stack.len(); sLen < operation.minStack {
			return nil, &ErrStackUnderflow{stackLen: sLen, required: operation.minStack}
		} else if sLen > operation.maxStack {
			return nil, &ErrStackOverflow{stackLen: sLen, limit: operation.maxStack}
		}

		if err := in.enforceRestrictions(op, operation, stack); err != nil {
			return nil, err
		}

		// Static portion of gas
		cost = operation.constantGas // For tracing
		if !contract.UseGas(operation.constantGas) {
			log15.Error("Run:outOfGas", "op=", op.String(), "contract addr=", contract.self.Address().String(),
				"CallerAddress=", contract.CallerAddress.String(),
				"caller=", contract.caller.Address().String())
			return nil, ErrOutOfGas
		}

		var memorySize uint64
		// Memory check needs to be done prior to evaluating the dynamic gas portion,
		// to detect calculation overflows
		if operation.memorySize != nil {
			memSize, overflow := operation.memorySize(stack)
			if overflow {
				return nil, ErrGasUintOverflow
			}
			// memory is expanded in words of 32 bytes. Gas
			// is also calculated in words.
			if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
				return nil, ErrGasUintOverflow
			}
		}
		// Dynamic portion of gas
		// consume the gas and return an error if not enough gas is available.
		// cost is explicitly set so that the capture state defer method can get the proper cost

		if operation.dynamicGas != nil {
			var dynamicCost uint64
			dynamicCost, err = operation.dynamicGas(in.evm, contract, stack, mem, memorySize)
			cost += dynamicCost // total cost, for debug tracing
			if err != nil || !contract.UseGas(dynamicCost) {
				log15.Error("Run:outOfGas", "op=", op.String(), "contract addr=", contract.self.Address().String(),
					"CallerAddress=", contract.CallerAddress.String(),
					"caller=", contract.caller.Address().String())
				return nil, ErrOutOfGas
			}
		}
		if memorySize > 0 {
			mem.Resize(memorySize)
		}

		if EVMDebugOn == in.cfg.Debug {
			in.cfg.Tracer.CaptureState(in.evm, pc, op, gasCopy, cost, mem, stack, in.returnData, contract, in.evm.depth, err)
			logged = true
		}

		res, err = operation.execute(&pc, in.evm, callContext)
		if operation.returns {
			in.returnData = common.CopyBytes(res)
		}

		switch {
		case err != nil:
			return nil, err
		case operation.reverts:
			return res, model.ErrExecutionReverted
		case operation.halts:
			return res, nil
		case !operation.jumps:
			pc++
		}
	}
	return nil, nil
}

// CanRun tells if the contract, passed as an argument, can be
// run by the current interpreter.
func (in *Interpreter) CanRun(code []byte) bool {
	return true
}

func buildGasParam(contract *Contract) *params.GasParam {
	return &params.GasParam{Gas: contract.Gas, Address: contract.Address()}
}

func buildEVMParam(evm *EVM) *params.EVMParam {
	return &params.EVMParam{
		StateDB:     evm.StateDB,
		CallGasTemp: evm.callGasTemp,
		BlockNumber: evm.BlockNumber,
	}
}

func fillEVM(param *params.EVMParam, evm *EVM) {
	evm.callGasTemp = param.CallGasTemp
}

// callCtx contains the things that are per-call, such as stack and memory,
// but not transients like pc and gas
type callCtx struct {
	memory *Memory
	stack  *Stack
	//rstack   *ReturnStack
	contract *Contract
}
