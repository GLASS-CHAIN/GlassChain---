// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/paracross/types"
)

// ： ， paracros titl  , ，

//NewParaAccount create new paracross account
// {}, ， ， 
// paracros 
// execName:  user.p.{guodun}.paracross
// symbol: coins.bty, token.{TEST}
//  mavl-{paracross}-coins.bty-{user-address} title{paracross}
// paracross  malv-coins-bty-exec-{Address(paracross)}:{Address(user.p.{guodun}.paracross)}
func NewParaAccount(cfg *types.Chain33Config, paraTitle, mainExecName, mainSymbol string, db db.KV) (*account.DB, error) {
	// ， title   "." 
	// paraExec := paraTitle + types.ParaX
	paraExec := pt.ParaX // ，
	//  ) tokne  symbol   coins  symbol ，  localExecName 
	paraSymbol := mainExecName + "." + mainSymbol
	return account.NewAccountDB(cfg, paraExec, paraSymbol, db)
}

//NewMainAccount create new Main account
// ， paracros 
// execName: paracross
// symbol: user.p.{guodun}.coins.{guodun}  user.p.{guodun}.token.{TEST}
//  mavl-paracross-user.p.{guodun}.coins.guodun-{user-address}
//   mavl-coins-{guodun}-exec-{Address(paracross)}:{Address(paracross)}
func NewMainAccount(cfg *types.Chain33Config, paraTitle, paraExecName, paraSymbol string, db db.KV) (*account.DB, error) {
	mainSymbol := paraTitle + paraExecName + "." + paraSymbol
	return account.NewAccountDB(cfg, pt.ParaX, mainSymbol, db)
}

func assetDepositBalance(acc *account.DB, addr string, amount int64) (*types.Receipt, error) {
	if !acc.CheckAmount(amount) {
		return nil, types.ErrAmount
	}
	acc1 := acc.LoadAccount(addr)
	copyacc := *acc1
	acc1.Balance += amount
	receiptBalance := &types.ReceiptAccountTransfer{
		Prev:    &copyacc,
		Current: acc1,
	}
	acc.SaveAccount(acc1)
	ty := int32(pt.TyLogParaAssetDeposit)
	log1 := &types.ReceiptLog{
		Ty:  ty,
		Log: types.Encode(receiptBalance),
	}
	kv := acc.GetKVSet(acc1)
	return &types.Receipt{
		Ty:   types.ExecOk,
		KV:   kv,
		Logs: []*types.ReceiptLog{log1},
	}, nil
}

func assetWithdrawBalance(acc *account.DB, addr string, amount int64) (*types.Receipt, error) {
	if !acc.CheckAmount(amount) {
		return nil, types.ErrAmount
	}
	acc1 := acc.LoadAccount(addr)
	if acc1.Balance-amount < 0 {
		return nil, types.ErrNoBalance
	}
	copyacc := *acc1
	acc1.Balance -= amount
	receiptBalance := &types.ReceiptAccountTransfer{
		Prev:    &copyacc,
		Current: acc1,
	}
	acc.SaveAccount(acc1)
	ty := int32(pt.TyLogParaAssetWithdraw)
	log1 := &types.ReceiptLog{
		Ty:  ty,
		Log: types.Encode(receiptBalance),
	}
	kv := acc.GetKVSet(acc1)
	return &types.Receipt{
		Ty:   types.ExecOk,
		KV:   kv,
		Logs: []*types.ReceiptLog{log1},
	}, nil
}

//                          trade add                                user address
// mavl-token-test-exec-1HPkPopVe3ERfvaAgedDtJQ792taZFEHCe:13DP8mVru5Rtu6CrjXQMvLsjvve3epRR1i
// mavl-conis-bty-exec-{para}1e:13DP8mVru5Rtu6CrjXQMvLsjvve3epRR1i

// 
//      mavl-  ` -  ` - 

// 
//      mavl-  ` -  ` -  exec -  ` ：                             10    - 5 |  5
//      mavl-  ` -  ` -  exec -  ` ：  paracros `                      |  5

// 
//     mavl-  ` -  ` - `
//
// titl hu

// 
//    `    paracross  ` :     user.p.guodun.paracross`
//     `         coins.bty
// mavl-  ` -  ` - 
//

// transfer / withdraw
//

// mavl -exec  -  symbol - addr

//  token TEST   -> trade
//                                                                         token-symbol{TEST}-addr{trade}:addr{user}
//  token TEST -> paracross ->  paracross： token.TEST ->  trade:   token.TEST@user.p.guodun.paracross
//           TEST     token-TEST-addr{paracross}:addr{user}
//                                    paracross-symbol{token.TEST}:addr{user}
//                                                                         paracross-symbol{token.TEST}-addr{trade}:addr{user}
//  token  TEST -> trade
//                                                                         token-symbol{TEST}-addr{trade}:addr{user}
//  TEST, 
//  token  TEST -> paracross
//                                                                        token-symbol{TEST}-addr{paracross}:addr{user}
//  ， ,   exec + symbol

// ，  ， symbol  symbol@host-title,   @bityuan
//  ".", host-title.exec.symbol, host-title , , 
//    titleFro    titleTo
// token TEST -> paracross     ->  paracross ->   trade
//                token-symbol{TEST}-addr{paracross}:addr{user}
//                                  paracross-symbol{token.TEST@tileFrom}-addr{user}
//                                                   paracross-symbol{token.TEST@tileFrom}-addr{trade}:addr{user}
//                                        ->   titleTo  paracross -> titleTo.trade
//                                              paracross-symbol{token.TEST@tileFrom}-addr{user}
//                                                                      paracross-symbol{token.TEST@tileFrom}-addr{trade}:addr{user}

/*
 ， trad ， 

 ？ ， 
  1. 
  1. bty
  1. 
  1. token  YCC
 

 trade ， 。
*/
