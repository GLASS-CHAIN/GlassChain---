// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	auty "github.com/33cn/plugin/plugin/dapp/autonomy/types"
)

func (a *Autonomy) Exec_PropBoard(payload *auty.ProposalBoard, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.propBoard(payload)
}

func (a *Autonomy) Exec_RvkPropBoard(payload *auty.RevokeProposalBoard, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.rvkPropBoard(payload)
}

func (a *Autonomy) Exec_VotePropBoard(payload *auty.VoteProposalBoard, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.votePropBoard(payload)
}

func (a *Autonomy) Exec_TmintPropBoard(payload *auty.TerminateProposalBoard, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.tmintPropBoard(payload)
}

func (a *Autonomy) Exec_PropProject(payload *auty.ProposalProject, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.propProject(payload)
}

func (a *Autonomy) Exec_RvkPropProject(payload *auty.RevokeProposalProject, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.rvkPropProject(payload)
}

func (a *Autonomy) Exec_VotePropProject(payload *auty.VoteProposalProject, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.votePropProject(payload)
}

func (a *Autonomy) Exec_PubVotePropProject(payload *auty.PubVoteProposalProject, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.pubVotePropProject(payload)
}

func (a *Autonomy) Exec_TmintPropProject(payload *auty.TerminateProposalProject, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.tmintPropProject(payload)
}


func (a *Autonomy) Exec_PropRule(payload *auty.ProposalRule, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.propRule(payload)
}

func (a *Autonomy) Exec_RvkPropRule(payload *auty.RevokeProposalRule, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.rvkPropRule(payload)
}

func (a *Autonomy) Exec_VotePropRule(payload *auty.VoteProposalRule, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.votePropRule(payload)
}

func (a *Autonomy) Exec_TmintPropRule(payload *auty.TerminateProposalRule, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.tmintPropRule(payload)
}

func (a *Autonomy) Exec_Transfer(payload *auty.TransferFund, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.transfer(payload)
}

func (a *Autonomy) Exec_CommentProp(payload *auty.Comment, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.commentProp(payload)
}


func (a *Autonomy) Exec_PropChange(payload *auty.ProposalChange, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.propChange(payload)
}

func (a *Autonomy) Exec_RvkPropChange(payload *auty.RevokeProposalChange, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.rvkPropChange(payload)
}

func (a *Autonomy) Exec_VotePropChange(payload *auty.VoteProposalChange, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.votePropChange(payload)
}

func (a *Autonomy) Exec_TmintPropChange(payload *auty.TerminateProposalChange, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(a, tx, int32(index))
	return action.tmintPropChange(payload)
}
