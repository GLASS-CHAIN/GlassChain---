// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	mty "github.com/33cn/plugin/plugin/dapp/multisig/types"
)

func (m *MultiSig) ExecDelLocal_MultiSigAccCreate(payload *mty.MultiSigAccCreate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.execLocalMultiSigReceipt(receiptData, tx, false)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}

func (m *MultiSig) ExecDelLocal_MultiSigOwnerOperate(payload *mty.MultiSigOwnerOperate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.execLocalMultiSigReceipt(receiptData, tx, false)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}

func (m *MultiSig) ExecDelLocal_MultiSigAccOperate(payload *mty.MultiSigAccOperate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.execLocalMultiSigReceipt(receiptData, tx, false)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}

func (m *MultiSig) ExecDelLocal_MultiSigConfirmTx(payload *mty.MultiSigConfirmTx, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.execLocalMultiSigReceipt(receiptData, tx, false)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}

func (m *MultiSig) ExecDelLocal_MultiSigExecTransferTo(payload *mty.MultiSigExecTransferTo, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.saveMultiSigTransfer(tx, mty.IsSubmit, false)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}

func (m *MultiSig) ExecDelLocal_MultiSigExecTransferFrom(payload *mty.MultiSigExecTransferFrom, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.execLocalMultiSigReceipt(receiptData, tx, false)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}
