// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	auty "github.com/33cn/plugin/plugin/dapp/autonomy/types"
)

func (a *Autonomy) Query_GetProposalBoard(in *types.ReqString) (types.Message, error) {
	return a.getProposalBoard(in)
}

func (a *Autonomy) Query_ListProposalBoard(in *auty.ReqQueryProposalBoard) (types.Message, error) {
	return a.listProposalBoard(in)
}

func (a *Autonomy) Query_GetActiveBoard(in *types.ReqString) (types.Message, error) {
	return a.getActiveBoard()
}

func (a *Autonomy) Query_GetProposalProject(in *types.ReqString) (types.Message, error) {
	return a.getProposalProject(in)
}

func (a *Autonomy) Query_ListProposalProject(in *auty.ReqQueryProposalProject) (types.Message, error) {
	return a.listProposalProject(in)
}

func (a *Autonomy) Query_GetProposalRule(in *types.ReqString) (types.Message, error) {
	return a.getProposalRule(in)
}

func (a *Autonomy) Query_ListProposalRule(in *auty.ReqQueryProposalRule) (types.Message, error) {
	return a.listProposalRule(in)
}

func (a *Autonomy) Query_GetActiveRule(in *types.ReqString) (types.Message, error) {
	return a.getActiveRule()
}

func (a *Autonomy) Query_ListProposalComment(in *auty.ReqQueryProposalComment) (types.Message, error) {
	return a.listProposalComment(in)
}

func (a *Autonomy) Query_GetProposalChange(in *types.ReqString) (types.Message, error) {
	return a.getProposalChange(in)
}

func (a *Autonomy) Query_ListProposalChange(in *auty.ReqQueryProposalChange) (types.Message, error) {
	return a.listProposalChange(in)
}
