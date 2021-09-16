// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"github.com/holiman/uint256"

	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
)

type Tracer interface {
	CaptureStart(from common.Address, to common.Address, call bool, input []byte, gas uint64, value uint64) error

	CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, rData []byte, contract *Contract, depth int, err error) error

	CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error

	CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error
}

type JSONLogger struct {
	encoder *json.Encoder
}

// Storage represents a contract's storage.
type Storage map[common.Hash]common.Hash

// Copy duplicates the current storage.
func (s Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range s {
		cpy[key] = value
	}
	return cpy
}

// LogConfig are the configuration options for structured logger the EVM
type LogConfig struct {
	DisableMemory     bool // disable memory capture
	DisableStack      bool // disable stack capture
	DisableStorage    bool // disable storage capture
	DisableReturnData bool // disable return data capture
	Debug             bool // print output during capture end
	Limit             int  // maximum length of output, but zero means unlimited
}

type StructLog struct {

	Pc uint64 `json:"pc"`

	Op OpCode `json:"op"`

	Gas uint64 `json:"gas"`

	GasCost uint64 `json:"gasCost"`

	Memory []string `json:"memory"`

	MemorySize int `json:"memSize"`

	Stack []*big.Int `json:"stack"`

	ReturnStack []uint32 `json:"returnStack"`

	ReturnData []byte `json:"returnData"`

	Storage map[common.Hash]common.Hash `json:"-"`

	Depth int `json:"depth"`

	RefundCounter uint64 `json:"refund"`

	Err error `json:"-"`
}

func NewJSONLogger(writer io.Writer) *JSONLogger {
	return &JSONLogger{json.NewEncoder(writer)}
}

func (logger *JSONLogger) CaptureStart(from common.Address, to common.Address, create bool, input []byte, gas uint64, value uint64) error {
	return nil
}

func (logger *JSONLogger) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, rData []byte, contract *Contract, depth int, err error) error {
	log := StructLog{
		Pc:         pc,
		Op:         op,
		Gas:        gas,
		GasCost:    cost,
		MemorySize: memory.Len(),
		Storage:    nil,
		Depth:      depth,
		Err:        err,
	}
	log.Memory = formatMemory(memory.Data())
	log.Stack = formatStack(stack.Data())
	log.ReturnData = rData
	return logger.encoder.Encode(log)
}

func formatStack(data []uint256.Int) (res []*big.Int) {
	for _, v := range data {
		res = append(res, v.ToBig())
	}
	return
}

func formatMemory(data []byte) (res []string) {
	for idx := 0; idx < len(data); idx += 32 {
		res = append(res, common.Bytes2HexTrim(data[idx:idx+32]))
	}
	return
}

func (logger *JSONLogger) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error {
	return nil
}

func (logger *JSONLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	type endLog struct {
		Output  string        `json:"output"`
		GasUsed int64         `json:"gasUsed"`
		Time    time.Duration `json:"time"`
		Err     string        `json:"error,omitempty"`
	}

	if err != nil {
		return logger.encoder.Encode(endLog{common.Bytes2Hex(output), int64(gasUsed), t, err.Error()})
	}
	return logger.encoder.Encode(endLog{common.Bytes2Hex(output), int64(gasUsed), t, ""})
}

type mdLogger struct {
	out io.Writer
	cfg *LogConfig
}

// NewMarkdownLogger creates a logger which outputs information in a format adapted
// for human readability, and is also a valid markdown table
func NewMarkdownLogger(cfg *LogConfig, writer io.Writer) *mdLogger {
	l := &mdLogger{writer, cfg}
	if l.cfg == nil {
		l.cfg = &LogConfig{}
	}
	return l
}

func (t *mdLogger) CaptureStart(from common.Address, to common.Address, create bool, input []byte, gas uint64, value uint64) error {
	if !create {
		fmt.Fprintf(t.out, "From: `%v`\nTo: `%v`\nData: `0x%x`\nGas: `%d`\nValue `%v` wei\n",
			from.String(), to.String(),
			input, gas, value)
	} else {
		fmt.Fprintf(t.out, "From: `%v`\nCreate at: `%v`\nData: `0x%x`\nGas: `%d`\nValue `%v` wei\n",
			from.String(), to.String(),
			input, gas, value)
	}

	fmt.Fprintf(t.out, `
|  Pc   |      Op     | Cost |   Stack   |   RStack  |  Refund |
|-------|-------------|------|-----------|-----------|---------|
`)
	return nil
}

func (t *mdLogger) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, rData []byte, contract *Contract, depth int, err error) error {
	fmt.Fprintf(t.out, "| %4d  | %10v  |  %3d |", pc, op, cost)

	if !t.cfg.DisableStack {
		// format stack
		var a []string
		for _, elem := range stack.data {
			a = append(a, fmt.Sprintf("%v", elem.String()))
		}
		b := fmt.Sprintf("[%v]", strings.Join(a, ","))
		fmt.Fprintf(t.out, "%10v |", b)
	}
	fmt.Fprintf(t.out, "%10v |", env.StateDB.GetRefund())
	fmt.Fprintln(t.out, "")
	if err != nil {
		fmt.Fprintf(t.out, "Error: %v\n", err)
	}
	return nil
}

func (t *mdLogger) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error {

	fmt.Fprintf(t.out, "\nError: at pc=%d, op=%v: %v\n", pc, op, err)

	return nil
}

func (t *mdLogger) CaptureEnd(output []byte, gasUsed uint64, tm time.Duration, err error) error {
	fmt.Fprintf(t.out, "\nOutput: `0x%x`\nConsumed gas: `%d`\nError: `%v`\n",
		output, gasUsed, err)
	return nil
}
