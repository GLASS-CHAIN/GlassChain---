// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	dty "github.com/33cn/plugin/plugin/dapp/dposvote/types"
)


func (d *DPos) Exec_Regist(payload *dty.DposCandidatorRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.Regist(payload)
}


func (d *DPos) Exec_CancelRegist(payload *dty.DposCandidatorCancelRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.CancelRegist(payload)
}


func (d *DPos) Exec_ReRegist(payload *dty.DposCandidatorRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.ReRegist(payload)
}

func (d *DPos) Exec_Vote(payload *dty.DposVote, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.Vote(payload)
}

func (d *DPos) Exec_CancelVote(payload *dty.DposCancelVote, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.CancelVote(payload)
}

func (d *DPos) Exec_RegistVrfM(payload *dty.DposVrfMRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.RegistVrfM(payload)
}

func (d *DPos) Exec_RegistVrfRP(payload *dty.DposVrfRPRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.RegistVrfRP(payload)
}

func (d *DPos) Exec_RecordCB(payload *dty.DposCBInfo, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.RecordCB(payload)
}

func (d *DPos) Exec_RegistTopN(payload *dty.TopNCandidatorRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.RegistTopN(payload)
}
