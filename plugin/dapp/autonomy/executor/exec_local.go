// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	auty "github.com/33cn/plugin/plugin/dapp/autonomy/types"
)


func (a *Autonomy) ExecLocal_PropBoard(payload *auty.ProposalBoard, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalBoard(tx, receiptData)
}

func (a *Autonomy) ExecLocal_RvkPropBoard(payload *auty.RevokeProposalBoard, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalBoard(tx, receiptData)
}

func (a *Autonomy) ExecLocal_VotePropBoard(payload *auty.VoteProposalBoard, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalBoard(tx, receiptData)
}

func (a *Autonomy) ExecLocal_TmintPropBoard(payload *auty.TerminateProposalBoard, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalBoard(tx, receiptData)
}


func (a *Autonomy) ExecLocal_PropProject(payload *auty.ProposalProject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalProject(tx, receiptData)
}

func (a *Autonomy) ExecLocal_RvkPropProject(payload *auty.RevokeProposalProject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalProject(tx, receiptData)
}

func (a *Autonomy) ExecLocal_VotePropProject(payload *auty.VoteProposalProject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalProject(tx, receiptData)
}

func (a *Autonomy) ExecLocal_PubVotePropProject(payload *auty.PubVoteProposalProject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalProject(tx, receiptData)
}

func (a *Autonomy) ExecLocal_TmintPropProject(payload *auty.TerminateProposalProject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalProject(tx, receiptData)
}


func (a *Autonomy) ExecLocal_PropRule(payload *auty.ProposalRule, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalRule(tx, receiptData)
}

func (a *Autonomy) ExecLocal_RvkPropRule(payload *auty.RevokeProposalRule, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalRule(tx, receiptData)
}

func (a *Autonomy) ExecLocal_VotePropRule(payload *auty.VoteProposalRule, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalRule(tx, receiptData)
}

func (a *Autonomy) ExecLocal_TmintPropRule(payload *auty.TerminateProposalRule, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalRule(tx, receiptData)
}

func (a *Autonomy) ExecLocal_CommentProp(payload *auty.Comment, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalCommentProp(tx, receiptData)
}


func (a *Autonomy) ExecLocal_PropChange(payload *auty.ProposalChange, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalChange(tx, receiptData)
}

func (a *Autonomy) ExecLocal_RvkPropChange(payload *auty.RevokeProposalChange, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalChange(tx, receiptData)
}

func (a *Autonomy) ExecLocal_VotePropChange(payload *auty.VoteProposalChange, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalChange(tx, receiptData)
}

func (a *Autonomy) ExecLocal_TmintPropChange(payload *auty.TerminateProposalChange, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalChange(tx, receiptData)
}
