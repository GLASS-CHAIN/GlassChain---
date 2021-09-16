// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/holiman/uint256"
)

type ContractRef interface {
	Address() common.Address
}

type AccountRef common.Address

func (ar AccountRef) Address() common.Address { return (common.Address)(ar) }

type Contract struct {
	CallerAddress common.Address

	caller ContractRef

	self ContractRef

	Jumpdests Destinations

	// Locally cached result of JUMPDEST analysis
	analysis bitvec

	Code []bytecode

	CodeHash common.Hash

	CodeAddr *common.Address

	Input []byte

	Gas uint64

	value uint64

	DelegateCall bool
}

func NewContract(caller ContractRef, object ContractRef, value uint64, gas uint64) *Contract {

	c := &Contract{CallerAddress: caller.Address(), caller: caller, self: object}

	if parent, ok := caller.(*Contract); ok {
		c.Jumpdests = parent.Jumpdests
	} else {
		c.Jumpdests = make(Destinations)
	}

	c.Gas = gas
	c.value = value

	return c
}

func (c *Contract) validJumpdest(dest *uint256.Int) bool {
	udest, overflow := dest.Uint64WithOverflow()
	// PC cannot go beyond len(code) and certainly can't be bigger than 63bits.
	// Don't bother checking for JUMPDEST in that case.
	if overflow || udest >= uint64(len(c.Code)) {
		return false
	}
	// Only JUMPDESTs allowed for destinations
	if OpCode(c.Code[udest]) != JUMPDEST {
		return false
	}
	return c.isCode(udest)
}

func (c *Contract) validJumpSubdest(udest uint64) bool {
	// PC cannot go beyond len(code) and certainly can't be bigger than 63 bits.
	// Don't bother checking for BEGINSUB in that case.
	if int64(udest) < 0 || udest >= uint64(len(c.Code)) {
		return false
	}
	// Only BEGINSUBs allowed for destinations
	if OpCode(c.Code[udest]) != BEGINSUB {
		return false
	}
	return c.isCode(udest)
}

func (c *Contract) isCode(udest uint64) bool {
	// Do we have a contract hash already?
	if c.CodeHash != (common.Hash{}) {
		// Does parent context have the analysis?
		analysis, exist := c.Jumpdests[c.CodeHash]
		if !exist {
			// Do the analysis and save in parent context
			// We do not need to store it in c.analysis
			analysis = codeBitmap(c.Code)
			c.Jumpdests[c.CodeHash] = analysis
		}
		// Also stash it in current contract for faster access
		c.analysis = analysis
		return analysis.codeSegment(udest)
	}
	// We don't have the code hash, most likely a piece of initcode not already
	// in state trie. In that case, we do an analysis, and save it locally, so
	// we don't have to recalculate it for every JUMP instruction in the execution
	// However, we don't save it within the parent context
	if c.analysis == nil {
		c.analysis = codeBitmap(c.Code)
	}
	return c.analysis.codeSegment(udest)
}


func (c *Contract) AsDelegate() *Contract {
	c.DelegateCall = true

	parent := c.caller.(*Contract)

	c.CallerAddress = parent.CallerAddress

	c.value = parent.value
	return c
}

func (c *Contract) GetOp(n uint64) OpCode {
	return OpCode(c.GetByte(n))
}

func (c *Contract) GetByte(n uint64) byte {
	if n < uint64(len(c.Code)) {
		return c.Code[n]
	}

	return 0
}

func (c *Contract) Caller() common.Address {
	return c.CallerAddress
}

func (c *Contract) UseGas(gas uint64) (ok bool) {
	if c.Gas < gas {
		return false
	}
	c.Gas -= gas
	return true
}

func (c *Contract) Address() common.Address {
	return c.self.Address()
}

func (c *Contract) Value() uint64 {
	return c.value
}

func (c *Contract) SetCode(hash common.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
}

func (c *Contract) SetCallCode(addr *common.Address, hash common.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
	c.CodeAddr = addr
}
