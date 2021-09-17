// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"strings"

	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/paracross/types"
	token "github.com/33cn/plugin/plugin/dapp/token/types"
	"github.com/pkg/errors"
)

//SymbolBty ...
const SymbolBty = "bty"

/*
 　=　assetExec + assetSymbol 

				exec              		symbol								tx.title=user.p.test1   tx.title=user.p.test2
1. ：
				coins					bty                     	  		transfer                 transfer
				paracross				user.p.test1.coins.fzm    			withdraw                 transfer

2. ：
				user.p.test1.coins		fzm              					transfer                 NAN
    			user.p.test1.paracross	coins.bty    						withdraw                 NAN
    			user.p.test1.paracross	paracross.user.p.test2.coins.cny	withdraw                 NAN

 user.p.test1.paracross.paracross.user.p.test2.coins.cn ：
user.p.test1.paracross paracros ，　paracross.user.p.test2.coins.cn paracros paracros user.p.test2.coins.cn 
*/
func getCrossAction(transfer *pt.CrossAssetTransfer, txExecer string) (int64, error) {
	paraTitle, ok := types.GetParaExecTitleName(txExecer)
	if !ok {
		return pt.ParacrossNoneTransfer, errors.Wrapf(types.ErrInvalidParam, "asset cross transfer execer:%s should be user.p.xx", txExecer)
	}
	/ 
	if types.IsParaExecName(transfer.AssetExec) && !strings.Contains(transfer.AssetExec, paraTitle) {
		return pt.ParacrossNoneTransfer, errors.Wrapf(types.ErrInvalidParam, "asset execer=%s not belong to title=%s", transfer.AssetExec, paraTitle)
	}

	//paracros withdra  
	if types.IsParaExecName(transfer.AssetExec) {
		if strings.Contains(transfer.AssetExec, pt.ParaX) {
			return pt.ParacrossMainAssetWithdraw, nil
		}
		return pt.ParacrossParaAssetTransfer, nil
	}

	if strings.Contains(transfer.AssetExec, pt.ParaX) && strings.Contains(transfer.AssetSymbol, paraTitle) {
		return pt.ParacrossParaAssetWithdraw, nil
	}
	return pt.ParacrossMainAssetTransfer, nil

}

/*
 
								      								type			realExec    realSymbol
coins+bty															mainTransfer	coins		bty
paracross+user.p.test1.coins.bty									paraWithdraw	coins		bty
user.p.test1.coins+bty												paraTransfer    coins		bty
user.p.test1.paracross+coins.bty									mainWithdraw	coins		bty
paracross+user.p.test1.coins.bty(->user.p.test2)					mainTransfer 	paracross   user.p.test1.coins.bty
user.p.test2.paracross+paracross.user.p.test1.coins.bty 			mainWithdraw	paracross   user.p.test1.coins.bty
 :
1. user.p.test1.coins+bt  coins accoun mavl-coins-bty-  mavl-user.p.test.coins-bty-
2. paracros  withdra  symbo 
　　a. 　mavl-paracross-coins.bty-exec-addr(user)
　　b. 　mavl-coins-bty-exec-addr{paracross}:addr{user}, coin 
*/
func amendTransferParam(transfer *pt.CrossAssetTransfer, act int64) (*pt.CrossAssetTransfer, error) {
	newTransfer := *transfer
	//exec=user.p.test1.coins -> exec=coins
	if types.IsParaExecName(transfer.AssetExec) {
		elements := strings.Split(transfer.AssetExec, ".")
		newTransfer.AssetExec = elements[len(elements)-1]
	}

	//paracross　exec's symbol should contain ".", non-paracross exec should not contain "."
	if newTransfer.AssetExec == pt.ParaX && !strings.Contains(newTransfer.AssetSymbol, ".") {
		return nil, errors.Wrapf(types.ErrInvalidParam, "paracross exec=%s, the symbol=%s should contain '.'", newTransfer.AssetExec, transfer.AssetSymbol)
	}

	if newTransfer.AssetExec != pt.ParaX && strings.Contains(newTransfer.AssetSymbol, ".") {
		return nil, errors.Wrapf(types.ErrInvalidParam, "non-paracross exec=%s, symbol=%s should not contain '.'", newTransfer.AssetExec, transfer.AssetSymbol)
	}

	if act == pt.ParacrossMainAssetWithdraw {
		e := strings.Split(transfer.AssetSymbol, ".")
		if len(e) <= 1 {
			return nil, errors.Wrapf(types.ErrInvalidParam, "main asset withdraw symbol=%s should be exec.symbol", transfer.AssetSymbol)
		}
		newTransfer.AssetExec = e[0]
		newTransfer.AssetSymbol = strings.Join(e[1:], ".")
		return &newTransfer, nil
	}

	/ user.p.{para}.coins.ccny prefi  coins.ccny
	if act == pt.ParacrossParaAssetWithdraw {
		e := strings.Split(transfer.AssetSymbol, ".")
		if len(e) <= 1 {
			return nil, errors.Wrapf(types.ErrInvalidParam, "para asset withdraw symbol=%s should be exec.symbol", transfer.AssetSymbol)
		}
		newTransfer.AssetSymbol = e[len(e)-1]
		newTransfer.AssetExec = e[len(e)-2]
		return &newTransfer, nil
	}
	return &newTransfer, nil
}

func (a *action) crossAssetTransfer(transfer *pt.CrossAssetTransfer, act int64, actTx *types.Transaction) (*types.Receipt, error) {
	newTransfer, err := amendTransferParam(transfer, act)
	if err != nil {
		return nil, err
	}
	clog.Info("paracross.crossAssetTransfer", "action", act, "newExec", newTransfer.AssetExec, "newSymbol", newTransfer.AssetSymbol,
		"ori.exec", transfer.AssetExec, "ori.symbol", transfer.AssetSymbol, "txHash", common.ToHex(actTx.Hash()))
	switch act {
	case pt.ParacrossMainAssetTransfer:
		return a.mainAssetTransfer(newTransfer, actTx)
	case pt.ParacrossMainAssetWithdraw:
		return a.mainAssetWithdraw(newTransfer, actTx)
	case pt.ParacrossParaAssetTransfer:
		return a.paraAssetTransfer(newTransfer, actTx)
	case pt.ParacrossParaAssetWithdraw:
		return a.paraAssetWithdraw(newTransfer, actTx)
	default:
		return nil, types.ErrNotSupport
	}
}

/ transfer, create　asset,  rollback
func (a *action) mainAssetTransfer(transfer *pt.CrossAssetTransfer, transferTx *types.Transaction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	isPara := cfg.IsPara()
	/ , 
	if !isPara {
		return a.execTransfer(transfer, transferTx)
	}
	return a.execCreateAsset(transfer, transferTx)
}

/ ， withdraw
func (a *action) mainAssetWithdraw(withdraw *pt.CrossAssetTransfer, withdrawTx *types.Transaction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	isPara := cfg.IsPara()
	/  ，a.t 
	if !isPara {
		return a.execWithdraw(withdraw, withdrawTx)
	}
	return a.execDestroyAsset(withdraw, withdrawTx)
}

/ ， create asset
func (a *action) paraAssetTransfer(transfer *pt.CrossAssetTransfer, transferTx *types.Transaction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	isPara := cfg.IsPara()
	/  
	if isPara {
		return a.execTransfer(transfer, transferTx)
	}
	/ 
	return a.execCreateAsset(transfer, transferTx)
}

/ ，  ，  
func (a *action) paraAssetWithdraw(withdraw *pt.CrossAssetTransfer, withdrawTx *types.Transaction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	isPara := cfg.IsPara()
	/  
	if isPara {
		return a.execWithdraw(withdraw, withdrawTx)
	}
	return a.execDestroyAsset(withdraw, withdrawTx)
}

func (a *action) execTransfer(transfer *pt.CrossAssetTransfer, transferTx *types.Transaction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	accDB, err := a.createAccount(cfg, a.db, transfer.AssetExec, transfer.AssetSymbol)
	if err != nil {
		return nil, errors.Wrapf(err, "execTransfer.createAccount,exec=%s,symbol=%s", transfer.AssetExec, transfer.AssetSymbol)
	}

	/ toAdd user.p.xx.paracros 
	execAddr := address.ExecAddress(pt.ParaX)
	toAddr := address.ExecAddress(string(transferTx.Execer))
	/ toAdd paracros 
	if cfg.IsPara() {
		execAddr = address.ExecAddress(string(transferTx.Execer))
		toAddr = address.ExecAddress(pt.ParaX)
	}

	clog.Debug("paracross.execTransfer", "execer", string(transferTx.Execer), "assetexec", transfer.AssetExec, "symbol", transfer.AssetSymbol,
		"txHash", common.ToHex(transferTx.Hash()))

	/ paracros  paracros    
	if transfer.AssetExec == pt.ParaX {
		r, err := accDB.Transfer(transferTx.From(), toAddr, transfer.Amount)
		if err != nil {
			return nil, errors.Wrapf(err, "assetTransfer,assetExec=%s,assetSym=%s", transfer.AssetExec, transfer.AssetSymbol)
		}
		return r, nil
	}

	fromAcc := accDB.LoadExecAccount(transferTx.From(), execAddr)
	if fromAcc.Balance < transfer.Amount {
		return nil, errors.Wrapf(types.ErrNoBalance, "execTransfer,acctBalance=%d,assetExec=%s,assetSym=%s", fromAcc.Balance, transfer.AssetExec, transfer.AssetSymbol)
	}
	r, err := accDB.ExecTransfer(transferTx.From(), toAddr, execAddr, transfer.Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "assetTransfer,assetExec=%s,assetSym=%s", transfer.AssetExec, transfer.AssetSymbol)
	}
	return r, nil
}

//withdra ，a.t ，　withdrawT withdra tx
func (a *action) execWithdraw(withdraw *pt.CrossAssetTransfer, withdrawTx *types.Transaction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	accDB, err := a.createAccount(cfg, a.db, withdraw.AssetExec, withdraw.AssetSymbol)
	if err != nil {
		return nil, errors.Wrapf(err, "execWithdraw.createAccount,exec=%s,symbol=%s", withdraw.AssetExec, withdraw.AssetSymbol)
	}
	execAddr := address.ExecAddress(pt.ParaX)
	fromAddr := address.ExecAddress(string(withdrawTx.Execer))
	if cfg.IsPara() {
		execAddr = address.ExecAddress(string(withdrawTx.Execer))
		fromAddr = address.ExecAddress(pt.ParaX)
	}

	clog.Debug("Paracross.execWithdraw", "amount", withdraw.Amount, "from", fromAddr,
		"assetExec", withdraw.AssetExec, "symbol", withdraw.AssetSymbol, "execAddr", execAddr, "txHash", common.ToHex(withdrawTx.Hash()))

	/ paracros  paracros    
	if withdraw.AssetExec == pt.ParaX {
		r, err := accDB.Transfer(fromAddr, withdraw.ToAddr, withdraw.Amount)
		if err != nil {
			return nil, errors.Wrapf(err, "assetWithdraw,assetExec=%s,assetSym=%s", withdraw.AssetExec, withdraw.AssetSymbol)
		}
		return r, nil
	}

	r, err := accDB.ExecTransfer(fromAddr, withdraw.ToAddr, execAddr, withdraw.Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "assetWithdraw,assetExec=%s,assetSym=%s", withdraw.AssetExec, withdraw.AssetSymbol)
	}
	return r, nil
}

/ Alic toke user.p.bb  mavl-paracross-token.symbol-Addr(Alice) Addr(user.p.bb.paracross 
/ toke mavl-paracross-user.p.aa.token.symbol-exec-Addr(Alice) user.p.bb  transfe paracros 
/ b  mavl-paracross-paracross.user.p.aa.token.symbol-exec-Addr(Alice) paracros paracross
func (a *action) createParaAccount(cross *pt.CrossAssetTransfer, crossTx *types.Transaction) (*account.DB, error) {
	cfg := a.api.GetConfig()
	paraTitle, err := getTitleFrom(crossTx.Execer)
	if err != nil {
		return nil, errors.Wrapf(err, "createParaAccount call getTitleFrom failed,exec=%s", string(crossTx.Execer))
	}

	assetExec := cross.AssetExec
	assetSymbol := cross.AssetSymbol
	if !cfg.IsPara() {
		assetExec = string(paraTitle) + assetExec
	}
	paraAcc, err := NewParaAccount(cfg, string(paraTitle), assetExec, assetSymbol, a.db)
	clog.Debug("createParaAccount", "assetExec", assetExec, "symbol", assetSymbol, "txHash", common.ToHex(crossTx.Hash()))
	if err != nil {
		return nil, errors.Wrapf(err, "createParaAccount,exec=%s,symbol=%s,title=%s", assetExec, assetSymbol, paraTitle)
	}
	return paraAcc, nil
}

func (a *action) execCreateAsset(transfer *pt.CrossAssetTransfer, transferTx *types.Transaction) (*types.Receipt, error) {
	paraAcc, err := a.createParaAccount(transfer, transferTx)
	if err != nil {
		return nil, errors.Wrapf(err, "createAsset")
	}
	clog.Debug("paracross.execCreateAsset", "assetExec", transfer.AssetExec, "symbol", transfer.AssetSymbol,
		"txHash", common.ToHex(transferTx.Hash()))

	r, err := assetDepositBalance(paraAcc, transfer.ToAddr, transfer.Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "createParaAsset,assetExec=%s,assetSym=%s", transfer.AssetExec, transfer.AssetSymbol)
	}
	return r, nil
}

func (a *action) execDestroyAsset(withdraw *pt.CrossAssetTransfer, withdrawTx *types.Transaction) (*types.Receipt, error) {
	paraAcc, err := a.createParaAccount(withdraw, withdrawTx)
	if err != nil {
		return nil, errors.Wrapf(err, "destroyAsset")
	}
	clog.Debug("paracross.execDestroyAsset", "assetExec", withdraw.AssetExec, "symbol", withdraw.AssetSymbol,
		"txHash", common.ToHex(withdrawTx.Hash()), "from", withdrawTx.From(), "amount", withdraw.Amount)
	r, err := assetWithdrawBalance(paraAcc, withdrawTx.From(), withdraw.Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "destroyAsset,assetExec=%s,assetSym=%s", withdraw.AssetExec, withdraw.AssetSymbol)
	}
	return r, nil
}

/  
func (a *action) assetTransfer(transfer *types.AssetsTransfer) (*types.Receipt, error) {
	tr := &pt.CrossAssetTransfer{
		AssetSymbol: transfer.Cointoken,
		Amount:      transfer.Amount,
		Note:        string(transfer.Note),
		ToAddr:      transfer.To,
	}
	adaptNullAssetExec(tr)
	return a.mainAssetTransfer(tr, a.tx)
}

/  
func (a *action) assetWithdraw(withdraw *types.AssetsWithdraw, withdrawTx *types.Transaction) (*types.Receipt, error) {
	tr := &pt.CrossAssetTransfer{
		AssetExec:   withdraw.ExecName,
		AssetSymbol: withdraw.Cointoken,
		Amount:      withdraw.Amount,
		Note:        string(withdraw.Note),
		ToAddr:      withdraw.To,
	}
	/ withdra  cointoke  token 
	if withdraw.Cointoken != "" {
		tr.AssetExec = token.TokenX
	}
	adaptNullAssetExec(tr)
	return a.mainAssetWithdraw(tr, withdrawTx)
}

func (a *action) assetTransferRollback(transfer *pt.CrossAssetTransfer, transferTx *types.Transaction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	isPara := cfg.IsPara()
	/ 
	if !isPara {
		accDB, err := a.createAccount(cfg, a.db, transfer.AssetExec, transfer.AssetSymbol)
		if err != nil {
			return nil, errors.Wrap(err, "assetTransferRollback.createAccount failed")
		}
		execAddr := address.ExecAddress(pt.ParaX)
		fromAcc := address.ExecAddress(string(transferTx.Execer))
		clog.Debug("paracross.AssetTransferRbk ", "exec", transfer.AssetExec, "sym", transfer.AssetSymbol,
			"transfer.txHash", common.ToHex(transferTx.Hash()), "curTx", common.ToHex(a.tx.Hash()))
		return accDB.ExecTransfer(fromAcc, transferTx.From(), execAddr, transfer.Amount)
	}
	return nil, nil
}

/ withdra   
func (a *action) paraAssetWithdrawRollback(wtw *pt.CrossAssetTransfer, withdrawTx *types.Transaction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	isPara := cfg.IsPara()
	/ 
	if !isPara {
		withdraw, err := amendTransferParam(wtw, pt.ParacrossParaAssetWithdraw)
		if err != nil {
			return nil, errors.Wrapf(err, "paraAssetWithdrawRollback amend param")
		}
		paraAcc, err := a.createParaAccount(withdraw, withdrawTx)
		if err != nil {
			return nil, errors.Wrapf(err, "createAsset")
		}
		clog.Debug("paracross.paraAssetWithdrawRollback", "exec", withdraw.AssetExec, "sym", withdraw.AssetSymbol,
			"transfer.txHash", common.ToHex(withdrawTx.Hash()), "curTx", common.ToHex(a.tx.Hash()))
		return assetDepositBalance(paraAcc, withdrawTx.From(), withdraw.Amount)
	}
	return nil, nil
}

func (a *action) createAccount(cfg *types.Chain33Config, db db.KV, exec, symbol string) (*account.DB, error) {
	var accDB *account.DB
	if symbol == "" {
		accDB = account.NewCoinsAccount(cfg)
		accDB.SetDB(db)
		return accDB, nil
	}
	if exec == "" {
		exec = token.TokenX
	}
	return account.NewAccountDB(cfg, exec, symbol, db)
}

/ assetTransfer,assetWithdra  AssetExec coins coins   。
/ AssetExec  tom 
/ coinsx coinsx.bt ，withdra coinsx
func adaptNullAssetExec(transfer *pt.CrossAssetTransfer) {
	if transfer.AssetSymbol == "" {
		transfer.AssetExec = "coins"
		transfer.AssetSymbol = SymbolBty
		return
	}
	if transfer.AssetExec == "" {
		transfer.AssetExec = token.TokenX
	}
}
