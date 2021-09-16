// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"math/big"

	"github.com/holiman/uint256"

	"github.com/33cn/chain33/common"
)

const (
	HashLength = 32

	Hash160Length = 20

	AddressLength = 20
)

type Hash common.Hash

func (h Hash) Str() string { return string(h[:]) }

func (h Hash) Bytes() []byte { return h[:] }

func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

func (h Hash) Hex() string { return Bytes2Hex(h[:]) }

func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

func BigToHash(b *big.Int) Hash {
	return Hash(common.BytesToHash(b.Bytes()))
}

func Uint256ToHash(u *uint256.Int) Hash {
	return Hash(common.BytesToHash(u.Bytes()))
}

func BytesToHash(b []byte) Hash {
	return Hash(common.BytesToHash(b))
}

func ToHash(data []byte) Hash {
	return BytesToHash(common.Sha256(data))
}
