// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
)

var multisiglog = log15.New("module", "execs.multisig")

// OwnerAdd : 
var (
	OwnerAdd     uint64 = 1
	OwnerDel     uint64 = 2
	OwnerModify  uint64 = 3
	OwnerReplace uint64 = 4
	//AccWeightOp 
	AccWeightOp     = true
	AccDailyLimitOp = false
	//OwnerOperate  ，owne ，accoun 
	OwnerOperate    uint64 = 1
	AccountOperate  uint64 = 2
	TransferOperate uint64 = 3
	//IsSubmit ：
	IsSubmit  = true
	IsConfirm = false

	MultiSigX            = "multisig"
	OneDaySecond   int64 = 24 * 3600
	MinOwnersInit        = 2
	MinOwnersCount       = 1  / owner
	MaxOwnersCount       = 20 / 2 owner

	Multisiglog = log15.New("module", MultiSigX)
)

// MultiSig actionid
const (
	ActionMultiSigAccCreate        = 10000
	ActionMultiSigOwnerOperate     = 10001
	ActionMultiSigAccOperate       = 10002
	ActionMultiSigConfirmTx        = 10003
	ActionMultiSigExecTransferTo   = 10004
	ActionMultiSigExecTransferFrom = 10005
)

/ logid
const (
	TyLogMultiSigAccCreate = 10000 / 

	TyLogMultiSigOwnerAdd     = 10001 / ad owner：add weight
	TyLogMultiSigOwnerDel     = 10002 / de owner：add weight
	TyLogMultiSigOwnerModify  = 10003 / modif owner：preweigh currentweight
	TyLogMultiSigOwnerReplace = 10004 / ol owne  owne ：addr+weight

	TyLogMultiSigAccWeightModify     = 10005 / ：preReqWeigh curReqWeight
	TyLogMultiSigAccDailyLimitAdd    = 10006 / ad DailyLimit：Symbo DailyLimit
	TyLogMultiSigAccDailyLimitModify = 10007 / modif DailyLimit：preDailyLimi currentDailyLimit

	TyLogMultiSigConfirmTx       = 10008 / 
	TyLogMultiSigConfirmTxRevoke = 10009 / 

	TyLogDailyLimitUpdate = 10010 //DailyLimi ，DailyLimi Submi Confir 
	TyLogMultiSigTx       = 10011 / Submi 
	TyLogTxCountUpdate    = 10012 //txcoun Submi 

)

//AccAssetsResult cl  amoun 
type AccAssetsResult struct {
	Execer   string `json:"execer,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	Currency int32  `json:"currency,omitempty"`
	Balance  string `json:"balance,omitempty"`
	Frozen   string `json:"frozen,omitempty"`
	Receiver string `json:"receiver,omitempty"`
	Addr     string `json:"addr,omitempty"`
}

//DailyLimitResult cli
type DailyLimitResult struct {
	Symbol     string `json:"symbol,omitempty"`
	Execer     string `json:"execer,omitempty"`
	DailyLimit string `json:"dailyLimit,omitempty"`
	SpentToday string `json:"spent,omitempty"`
	LastDay    string `json:"lastday,omitempty"`
}

//MultiSigResult cli
type MultiSigResult struct {
	CreateAddr     string              `json:"createAddr,omitempty"`
	MultiSigAddr   string              `json:"multiSigAddr,omitempty"`
	Owners         []*Owner            `json:"owners,omitempty"`
	DailyLimits    []*DailyLimitResult `json:"dailyLimits,omitempty"`
	TxCount        uint64              `json:"txCount,omitempty"`
	RequiredWeight uint64              `json:"requiredWeight,omitempty"`
}

//UnSpentAssetsResult cli
type UnSpentAssetsResult struct {
	Symbol  string `json:"symbol,omitempty"`
	Execer  string `json:"execer,omitempty"`
	UnSpent string `json:"unspent,omitempty"`
}

//IsAssetsInvalid ，Symbol  ：BTY,coins.BTY。exec types.AllowUserExe 
func IsAssetsInvalid(exec, symbol string) error {

	//exe 
	allowExeName := types.AllowUserExec
	nameLen := len(allowExeName)
	execValid := false
	for i := 0; i < nameLen; i++ {
		if exec == string(allowExeName[i]) {
			execValid = true
			break
		}
	}
	if !execValid {
		multisiglog.Error("IsAssetsInvalid", "exec", exec)
		return ErrInvalidExec
	}
	//Symbo 
	return nil
}
