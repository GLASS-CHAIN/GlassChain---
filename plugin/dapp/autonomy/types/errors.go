// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "errors"

var (
	ErrVotePeriod = errors.New("ErrVotePeriod")
	ErrProposalStatus = errors.New("ErrProposalStatus")
	ErrRepeatVoteAddr = errors.New("ErrRepeatVoteAddr")
	ErrRevokeProposalPeriod = errors.New("ErrRevokeProposalPeriod")
	ErrRevokeProposalPower = errors.New("ErrRevokeProposalPower")
	ErrTerminatePeriod = errors.New("ErrTerminatePeriod")
	ErrNoActiveBoard = errors.New("ErrNoActiveBoard")
	ErrNoAutonomyExec = errors.New("ErrNoAutonomyExec")
	ErrNoPeriodAmount = errors.New("ErrNoPeriodAmount")
	ErrMinerAddr = errors.New("ErrMinerAddr")
	ErrBindAddr = errors.New("ErrBindAddr")
	ErrChangeBoardAddr = errors.New("ErrChangeBoardAddr")
	ErrBoardNumber = errors.New("ErrBoardNumber")
	ErrRepeatAddr = errors.New("ErrRepeatAddr")
	ErrNotEnoughFund = errors.New("ErrNotEnoughFund")
	ErrSetBlockHeight = errors.New("ErrSetBlockHeight")
)
