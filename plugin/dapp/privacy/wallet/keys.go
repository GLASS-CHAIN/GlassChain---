// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import "fmt"

const (
	// PrivacyDBVersion  KE 
	PrivacyDBVersion = "Privacy-DBVersion"
	// Privacy4Addr KE 
	// KE   	Privacy4Addr 
	// VALU  types.WalletAccountPrivacy， 
	Privacy4Addr = "Privacy-Addr"
	// AvailUTXOs UTX KE 
	// KE   	AvailUTXOs-tokenname-address-outtxhash-outindex outtxhas UTX  common.Byte2Hex( 
	// VALU  types.PrivacyDBStore UTX 
	AvailUTXOs = "Privacy-UTXO"
	// UTXOsSpentInTx  UTX   UTX 
	// KE  	UTXOsSpentInTx：costtxhash  costtxhas  common.Byte2Hex( 
	// VALU 	types.FTXOsSTXOsInOneTx
	UTXOsSpentInTx = "Privacy-UTXOsSpentInTx"
	// FrozenUTXOs  UTX  KE 
	// KE   	FrozenUTXOs:tokenname-address-costtxhash costtxhas UTX  common.Byte2Hex( 
	// VALU  UTXOsSpentInT KE 
	FrozenUTXOs = "Privacy-FUTXO4Tx"
	// PrivacySTXO  UTX UTXO UTX KE 
	// KE 	PrivacySTXO-tokenname-address-costtxhash costtxhas UTX  common.Byte2Hex( 
	// VALU  UTXOsSpentInT KE 
	PrivacySTXO = "Privacy-SUTXO"
	// STXOs4Tx UTX KE 
	// KE 	STXOs4Tx：costtxhash costtxhas UTX  common.Byte2Hex( 
	// VALU  UTXOsSpentInT KE 
	STXOs4Tx = "Privacy-SUTXO4Tx"
	// RevertSendtx UTX 
	// KE 	RevertSendtx:tokenname-address-costtxhash costtxhas UTX  common.Byte2Hex( 
	// VALU  UTXOsSpentInT KE 
	RevertSendtx = "Privacy-RevertSendtx"
	// RecvPrivacyTx KE 
	// KE 	RecvPrivacyTx:tokenname-address-heighstr heighst types.MaxTxsPerBloc index
	// VALU  PrivacyT KE 
	RecvPrivacyTx = "Privacy-RecvTX"
	// SendPrivacyTx KE 
	// KE 	SendPrivacyTx:tokenname-address-heighstr heighst types.MaxTxsPerBloc index
	// VALU  PrivacyT KE 
	SendPrivacyTx = "Privacy-SendTX"
	// PrivacyTX KE 
	// KE 	PrivacyTX:heighstr heighst types.MaxTxsPerBloc index
	// VALU 	types.WalletTxDetail
	PrivacyTX = "Privacy-TX"
	// ScanPrivacyInput  UTX 
	// KE 	ScanPrivacyInput-outtxhash-outindex outtxhas UTX  common.Byte2Hex( 
	// VALU 	types.UTXOGlobalIndex
	ScanPrivacyInput = "Privacy-ScaneInput"
	// ReScanUtxosFlag UTX 
	// KE 	ReScanUtxosFlag
	// VALU 	types.Int64 
	//		UtxoFlagNoScan  int32 = 0
	//		UtxoFlagScaning int32 = 1
	//		UtxoFlagScanEnd int32 = 2
	ReScanUtxosFlag = "Privacy-RescanFlag"
)

func calcPrivacyDBVersion() []byte {
	return []byte(PrivacyDBVersion)
}

// calcUTXOKey UTX   
//key and prefix for privacy
//types.PrivacyDBStor calcUTXOKe key，
//1 utx  ke value calcUTXOKey4TokenAdd ke k ；
//2 ，calcUTXOKey4TokenAdd k  calcPrivacyFUTXOKe ke k  key，
//  utx futxo；
//3  futx   ke stx ，
//4  del bloc  
// 4.a stx  stx ftx ，
// 4.b utx ftx  utx ftx  types.PrivacyDBStor 
// 4.c stx    utx   
func calcUTXOKey(txhash string, index int) []byte {
	return []byte(fmt.Sprintf("%s-%s-%d", AvailUTXOs, txhash, index))
}

func calcKey4UTXOsSpentInTx(key string) []byte {
	return []byte(fmt.Sprintf("%s:%s", UTXOsSpentInTx, key))
}

// calcPrivacyAddrKey 
func calcPrivacyAddrKey(addr string) []byte {
	return []byte(fmt.Sprintf("%s-%s", Privacy4Addr, addr))
}

//calcAddrKey add Accoun 
func calcAddrKey(addr string) []byte {
	return []byte(fmt.Sprintf("Addr:%s", addr))
}

// calcPrivacyUTXOPrefix4Addr UTX KE 
func calcPrivacyUTXOPrefix4Addr(assetExec, token, addr string) []byte {
	return []byte(fmt.Sprintf("%s-%s-%s-%s-", AvailUTXOs, assetExec, token, addr))
}

// calcFTXOsKeyPrefix UTX KE 
func calcFTXOsKeyPrefix(assetExec, token, addr string) []byte {
	return []byte(fmt.Sprintf("%s:%s-%s-%s-", FrozenUTXOs, assetExec, token, addr))
}

// calcSendPrivacyTxKey 
// add 
func calcSendPrivacyTxKey(assetExec, assetSymbol, addr, txHeightIndex string) []byte {
	return []byte(fmt.Sprintf("%s:%s-%s-%s-%s", SendPrivacyTx, assetExec, assetSymbol, addr, txHeightIndex))
}

// calcRecvPrivacyTxKey 
// add 
// ke calcTxKey(heightstr 
func calcRecvPrivacyTxKey(assetExec, tokenname, addr, key string) []byte {
	return []byte(fmt.Sprintf("%s:%s-%s-%s-%s", RecvPrivacyTx, assetExec, tokenname, addr, key))
}

// calcUTXOKey4TokenAddr UTX Ke 
func calcUTXOKey4TokenAddr(assetExec, token, addr, txhash string, index int) []byte {
	return []byte(fmt.Sprintf("%s-%s-%s-%s-%s-%d", AvailUTXOs, assetExec, token, addr, txhash, index))
}

// calcKey4FTXOsInTx  UTX 
func calcKey4FTXOsInTx(assetExec, token, addr, txhash string) []byte {
	return []byte(fmt.Sprintf("%s:%s-%s-%s-%s", FrozenUTXOs, assetExec, token, addr, txhash))
}

// calcRescanUtxosFlagKey UTX 
func calcRescanUtxosFlagKey(addr string) []byte {
	return []byte(fmt.Sprintf("%s-%s", ReScanUtxosFlag, addr))
}

func calcScanPrivacyInputUTXOKey(txhash string, index int) []byte {
	return []byte(fmt.Sprintf("%s-%s-%d", ScanPrivacyInput, txhash, index))
}

func calcKey4STXOsInTx(txhash string) []byte {
	return []byte(fmt.Sprintf("%s:%s", STXOs4Tx, txhash))
}

// calcSTXOTokenAddrTxKey UTXO
func calcSTXOTokenAddrTxKey(assetExec, token, addr, txhash string) []byte {
	return []byte(fmt.Sprintf("%s-%s-%s-%s-%s", PrivacySTXO, assetExec, token, addr, txhash))
}

func calcSTXOPrefix4Addr(assetExec, token, addr string) []byte {
	return []byte(fmt.Sprintf("%s-%s-%s-%s-", PrivacySTXO, assetExec, token, addr))
}

// calcRevertSendTxKey UTX UTX 
func calcRevertSendTxKey(assetExec, tokenname, addr, txhash string) []byte {
	return []byte(fmt.Sprintf("%s:%s-%s-%s-%s", RevertSendtx, assetExec, tokenname, addr, txhash))
}

/ height*100000+index T 
//key:Tx:height*100000+index
func calcTxKey(key string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrivacyTX, key))
}
