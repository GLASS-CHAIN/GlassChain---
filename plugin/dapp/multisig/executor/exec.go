// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	mty "github.com/33cn/plugin/plugin/dapp/multisig/types"
)

func (m *MultiSig) Exec_MultiSigAccCreate(payload *mty.MultiSigAccCreate, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigAccCreate(payload)
}

func (m *MultiSig) Exec_MultiSigOwnerOperate(payload *mty.MultiSigOwnerOperate, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigOwnerOperate(payload)
}

func (m *MultiSig) Exec_MultiSigAccOperate(payload *mty.MultiSigAccOperate, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigAccOperate(payload)
}

func (m *MultiSig) Exec_MultiSigConfirmTx(payload *mty.MultiSigConfirmTx, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigConfirmTx(payload)
}

func (m *MultiSig) Exec_MultiSigExecTransferTo(payload *mty.MultiSigExecTransferTo, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigExecTransferTo(payload)
}

func (m *MultiSig) Exec_MultiSigExecTransferFrom(payload *mty.MultiSigExecTransferFrom, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigExecTransferFrom(payload)
}
