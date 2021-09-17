// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"
)

const (
	//MultiSigPrefix statedb
	MultiSigPrefix   = "mavl-multisig-"
	MultiSigTxPrefix = "mavl-multisig-tx-"

	//MultiSigLocalPrefix localdb multisig account count
	MultiSigLocalPrefix = "LODB-multisig-"
	MultiSigAccCount    = "acccount"
	MultiSigAcc         = "account"
	MultiSigAllAcc      = "allacc"
	MultiSigTx          = "tx"
	MultiSigRecvAssets  = "assets"
	MultiSigAccCreate   = "create"
)

func calcMultiSigAccountKey(multiSigAccAddr string) (key []byte) {
	return []byte(fmt.Sprintf(MultiSigPrefix+"%s", multiSigAccAddr))
}

func calcMultiSigAccTxKey(multiSigAccAddr string, txid uint64) (key []byte) {
	txstr := fmt.Sprintf("%018d", txid)
	return []byte(fmt.Sprintf(MultiSigTxPrefix+"%s-%s", multiSigAccAddr, txstr))
}


func calcMultiSigAccCountKey() []byte {
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s", MultiSigAccCount))
}
func calcMultiSigAllAcc(accindex int64) (key []byte) {
	accstr := fmt.Sprintf("%018d", accindex)
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s-%s", MultiSigAllAcc, accstr))
}

func calcMultiSigAcc(addr string) (key []byte) {
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s-%s", MultiSigAcc, addr))
}

func calcMultiSigAccCreateAddr(createAddr string) (key []byte) {
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s:-%s", MultiSigAccCreate, createAddr))
}

func calcMultiSigAccTx(addr string, txid uint64) (key []byte) {
	accstr := fmt.Sprintf("%018d", txid)

	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s-%s-%s", MultiSigTx, addr, accstr))
}

func calcAddrRecvAmountKey(addr, execname, symbol string) []byte {
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s-%s-%s-%s", MultiSigRecvAssets, addr, execname, symbol))
}

func calcAddrRecvAmountPrefix(addr string) []byte {
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s-%s-", MultiSigRecvAssets, addr))
}
