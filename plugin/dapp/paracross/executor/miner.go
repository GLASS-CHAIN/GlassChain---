// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/paracross/executor/minerrewards"
	pt "github.com/33cn/plugin/plugin/dapp/paracross/types"
	"github.com/pkg/errors"
)

/ miner t  t  ，preHas blockchai 
//note: Mine Height= ， 
/ bu  coi  coi  coin 
func (a *action) Miner(miner *pt.ParacrossMinerAction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	/ coin
	if miner.AddIssueCoins > 0 {
		return a.addIssueCoins(miner.AddIssueCoins)
	}

	if miner.Status.Title != cfg.GetTitle() || miner.Status.MainBlockHash == nil {
		return nil, pt.ErrParaMinerExecErr
	}

	var logs []*types.ReceiptLog
	var receipt = &pt.ReceiptParacrossMiner{}

	log := &types.ReceiptLog{}
	log.Ty = pt.TyLogParacrossMiner
	receipt.Status = miner.Status

	log.Log = types.Encode(receipt)
	logs = append(logs, log)

	minerReceipt := &types.Receipt{Ty: types.ExecOk, KV: nil, Logs: logs}

	on, err := a.isSelfConsensOn(miner)
	if err != nil {
		return nil, err
	}
	/ 
	if on {
		r, err := a.issueCoins(miner)
		if err != nil {
			return nil, err
		}

		minerReceipt = mergeReceipt(minerReceipt, r)
	}

	return minerReceipt, nil
}

// Non   manager  paracros  
func (a *action) addIssueCoins(amount int64) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	if !isSuperManager(cfg, a.fromaddr) {
		return nil, errors.Wrapf(types.ErrNotAllow, "addr=%s,is not super manager", a.fromaddr)
	}

	issueReceipt, err := a.coinsAccount.ExecIssueCoins(a.execaddr, amount)
	if err != nil {
		clog.Error("paracross miner issue err", "execAddr", a.execaddr, "amount", amount)
		return nil, errors.Wrap(err, "issueCoins")
	}
	return issueReceipt, nil

}

func (a *action) isSelfConsensOn(miner *pt.ParacrossMinerAction) (bool, error) {
	cfg := a.api.GetConfig()
	//ForkParaInitMinerHeigh   paracros    
	if cfg.IsDappFork(a.height, pt.ParaX, pt.ForkParaFullMinerHeight) {
		return true, nil
	}

	isSelfConsensOn := miner.IsSelfConsensus

	/    a.height
	/ 100  10 80~9 2  2 10   8 100
	/ miner.Status.Heigh   10 

	if cfg.IsDappFork(a.height, pt.ParaX, pt.ForkParaSelfConsStages) {
		var err error
		isSelfConsensOn, err = isSelfConsOn(a.db, miner.Status.Height)
		if err != nil && errors.Cause(err) != pt.ErrKeyNotExist {
			clog.Error("paracross miner getConsensus ", "height", miner.Status.Height, "err", err)
			return false, err
		}
	}
	return isSelfConsensOn, nil
}

func (a *action) issueCoins(miner *pt.ParacrossMinerAction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()

	mode := cfg.MGStr("mver.consensus.paracross.minerMode", a.height)
	if _, ok := minerrewards.MinerRewards[mode]; !ok {
		panic("getTotalReward not be set depend on consensus.paracross.minerMode")
	}

	coinReward, coinFundReward, _ := minerrewards.MinerRewards[mode].GetConfigReward(cfg, a.height)
	totalReward := coinReward + coinFundReward
	if totalReward > 0 {
		issueReceipt, err := a.coinsAccount.ExecIssueCoins(a.execaddr, totalReward)
		if err != nil {
			clog.Error("paracross miner issue err", "height", miner.Status.Height,
				"execAddr", a.execaddr, "amount", totalReward)
			return nil, err
		}
		return issueReceipt, nil
	}
	return nil, nil
}
