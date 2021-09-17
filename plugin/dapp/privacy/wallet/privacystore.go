// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	wcom "github.com/33cn/chain33/wallet/common"
	privacy "github.com/33cn/plugin/plugin/dapp/privacy/crypto"
	privacytypes "github.com/33cn/plugin/plugin/dapp/privacy/types"

	"github.com/golang/protobuf/proto"
)

const (
	// PRIVACYDBVERSION 
	PRIVACYDBVERSION int64 = 1
)

func newStore(db db.DB, cfg *types.Chain33Config) *privacyStore {
	return &privacyStore{Store: wcom.NewStore(db), Chain33Config: cfg}
}

// privacyStore 
type privacyStore struct {
	*types.Chain33Config
	*wcom.Store
}

func (store *privacyStore) getVersion() int64 {
	var version int64
	data, err := store.Get(calcPrivacyDBVersion())
	if err != nil || data == nil {
		bizlog.Debug("getVersion", "db.Get error", err)
		return 0
	}
	err = json.Unmarshal(data, &version)
	if err != nil {
		bizlog.Error("getVersion", "json.Unmarshal error", err)
		return 0
	}
	return version
}

func (store *privacyStore) setVersion() error {
	version := PRIVACYDBVERSION
	data, err := json.Marshal(&version)
	if err != nil || data == nil {
		bizlog.Error("setVersion", "json.Marshal error", err)
		return err
	}
	err = store.GetDB().SetSync(calcPrivacyDBVersion(), data)
	if err != nil {
		bizlog.Error("setVersion", "db.SetSync error", err)
	}
	return err
}

func (store *privacyStore) getAccountByPrefix(addr string) ([]*types.WalletAccountStore, error) {
	if len(addr) == 0 {
		bizlog.Error("getAccountByPrefix addr is nil")
		return nil, types.ErrInvalidParam
	}
	list := store.NewListHelper()
	accbytes := list.PrefixScan([]byte(addr))
	if len(accbytes) == 0 {
		bizlog.Error("getAccountByPrefix addr not exist")
		return nil, types.ErrAccountNotExist
	}
	WalletAccountStores := make([]*types.WalletAccountStore, len(accbytes))
	for index, accbyte := range accbytes {
		var walletaccount types.WalletAccountStore
		err := proto.Unmarshal(accbyte, &walletaccount)
		if err != nil {
			bizlog.Error("GetAccountByAddr", "proto.Unmarshal err:", err)
			return nil, types.ErrUnmarshal
		}
		WalletAccountStores[index] = &walletaccount
	}
	return WalletAccountStores, nil
}

func (store *privacyStore) getWalletAccountPrivacy(addr string) (*privacytypes.WalletAccountPrivacy, error) {
	if len(addr) == 0 {
		bizlog.Error("GetWalletAccountPrivacy addr is nil")
		return nil, types.ErrInvalidParam
	}

	privacyByte, err := store.Get(calcPrivacyAddrKey(addr))
	if err != nil {
		bizlog.Error("GetWalletAccountPrivacy", "db Get error ", err)
		return nil, err
	}
	if nil == privacyByte {
		return nil, privacytypes.ErrPrivacyNotEnabled
	}
	var accPrivacy privacytypes.WalletAccountPrivacy
	err = proto.Unmarshal(privacyByte, &accPrivacy)
	if err != nil {
		bizlog.Error("GetWalletAccountPrivacy", "proto.Unmarshal err:", err)
		return nil, types.ErrUnmarshal
	}
	return &accPrivacy, nil
}

func (store *privacyStore) getAccountByAddr(addr string) (*types.WalletAccountStore, error) {
	var account types.WalletAccountStore
	if len(addr) == 0 {
		bizlog.Error("GetAccountByAddr addr is nil")
		return nil, types.ErrInvalidParam
	}
	data, err := store.Get(calcAddrKey(addr))
	if data == nil || err != nil {
		if err != db.ErrNotFoundInDb {
			bizlog.Debug("GetAccountByAddr addr", "err", err)
		}
		return nil, types.ErrAddrNotExist
	}
	err = proto.Unmarshal(data, &account)
	if err != nil {
		bizlog.Error("GetAccountByAddr", "proto.Unmarshal err:", err)
		return nil, types.ErrUnmarshal
	}
	return &account, nil
}

func (store *privacyStore) setWalletAccountPrivacy(addr string, privacy *privacytypes.WalletAccountPrivacy) error {
	if len(addr) == 0 {
		bizlog.Error("SetWalletAccountPrivacy addr is nil")
		return types.ErrInvalidParam
	}
	if privacy == nil {
		bizlog.Error("SetWalletAccountPrivacy privacy is nil")
		return types.ErrInvalidParam
	}

	privacybyte := types.Encode(privacy)

	newbatch := store.NewBatch(true)
	newbatch.Set(calcPrivacyAddrKey(addr), privacybyte)
	newbatch.Write()

	return nil
}
func (store *privacyStore) listAvailableUTXOs(assetExec, token, addr string) ([]*privacytypes.PrivacyDBStore, error) {
	if 0 == len(addr) {
		bizlog.Error("listWalletPrivacyAccount addr is nil")
		return nil, types.ErrInvalidParam
	}

	list := store.NewListHelper()
	onetimeAccbytes := list.PrefixScan(calcPrivacyUTXOPrefix4Addr(assetExec, token, addr))
	if len(onetimeAccbytes) == 0 {
		bizlog.Error("listWalletPrivacyAccount ", "addr not exist", addr)
		return nil, nil
	}

	privacyDBStoreSlice := make([]*privacytypes.PrivacyDBStore, len(onetimeAccbytes))
	for index, acckeyByte := range onetimeAccbytes {
		var accPrivacy privacytypes.PrivacyDBStore
		accByte, err := store.Get(acckeyByte)
		if err != nil {
			bizlog.Error("listWalletPrivacyAccount", "db Get err:", err)
			return nil, err
		}
		err = proto.Unmarshal(accByte, &accPrivacy)
		if err != nil {
			bizlog.Error("listWalletPrivacyAccount", "proto.Unmarshal err:", err)
			return nil, types.ErrUnmarshal
		}
		privacyDBStoreSlice[index] = &accPrivacy
	}
	return privacyDBStoreSlice, nil
}

func (store *privacyStore) listFrozenUTXOs(assetExec, token, addr string) ([]*privacytypes.FTXOsSTXOsInOneTx, error) {
	if 0 == len(addr) {
		bizlog.Error("listFrozenUTXOs addr is nil")
		return nil, types.ErrInvalidParam
	}
	list := store.NewListHelper()
	values := list.List(calcFTXOsKeyPrefix(assetExec, token, addr), nil, 0, 0)
	if len(values) == 0 {
		bizlog.Error("listFrozenUTXOs ", "addr not exist", addr)
		return nil, nil
	}

	ftxoslice := make([]*privacytypes.FTXOsSTXOsInOneTx, 0)
	for _, acckeyByte := range values {
		var ftxotx privacytypes.FTXOsSTXOsInOneTx
		accByte, err := store.Get(acckeyByte)
		if err != nil {
			bizlog.Error("listFrozenUTXOs", "db Get err:", err)
			return nil, err
		}

		err = proto.Unmarshal(accByte, &ftxotx)
		if err != nil {
			bizlog.Error("listFrozenUTXOs", "proto.Unmarshal err:", err)
			return nil, types.ErrUnmarshal
		}
		ftxoslice = append(ftxoslice, &ftxotx)
	}
	return ftxoslice, nil
}

func (store *privacyStore) getWalletPrivacyTxDetails(param *privacytypes.ReqPrivacyTransactionList) (*types.WalletTxDetails, error) {
	if param == nil {
		bizlog.Error("getWalletPrivacyTxDetails param is nil")
		return nil, types.ErrInvalidParam
	}
	if param.SendRecvFlag != sendTx && param.SendRecvFlag != recvTx {
		bizlog.Error("procPrivacyTransactionList", "invalid sendrecvflag ", param.SendRecvFlag)
		return nil, types.ErrInvalidParam
	}

	list := store.NewListHelper()
	var txKeyBytes [][]byte
	if len(param.StartTxHeightIndex) == 0 {
		var keyPrefix []byte
		if param.SendRecvFlag == sendTx {
			keyPrefix = calcSendPrivacyTxKey(param.AssetExec, param.AssetSymbol, param.Address, "")
		} else {
			keyPrefix = calcRecvPrivacyTxKey(param.AssetExec, param.AssetSymbol, param.Address, "")
		}
		txKeyBytes = list.IteratorScanFromLast(keyPrefix, param.Count, db.ListDESC)

	} else {
		if param.SendRecvFlag == sendTx {
			txKeyBytes = list.IteratorScan([]byte(SendPrivacyTx), calcSendPrivacyTxKey(param.AssetExec, param.AssetSymbol, param.Address, param.StartTxHeightIndex), param.Count, param.Direction)
		} else {
			txKeyBytes = list.IteratorScan([]byte(RecvPrivacyTx), calcRecvPrivacyTxKey(param.AssetExec, param.AssetSymbol, param.Address, param.StartTxHeightIndex), param.Count, param.Direction)
		}
	}
	txDetails := &types.WalletTxDetails{}
	for _, keyByte := range txKeyBytes {
		value, err := store.Get(keyByte)
		if err != nil || value == nil {
			bizlog.Error("getWalletPrivacyTxDetails", "db Get error", err)
			continue
		}

		txDetail := &types.WalletTxDetail{}
		err = types.Decode(value, txDetail)
		if err != nil {
			bizlog.Error("getWalletPrivacyTxDetails", "proto.Unmarshal err:", err)
			return nil, types.ErrUnmarshal
		}
		txDetail.Txhash = txDetail.GetTx().Hash()
		if txDetail.GetTx().IsWithdraw(store.Chain33Config.GetCoinExec()) {
			//swap from and to
			txDetail.Fromaddr, txDetail.Tx.To = txDetail.Tx.To, txDetail.Fromaddr
		}
		txDetails.TxDetails = append(txDetails.TxDetails, txDetail)
	}
	return txDetails, nil
}

func (store *privacyStore) getPrivacyTokenUTXOs(assetExec, token, addr string) (*walletUTXOs, error) {
	list := store.NewListHelper()
	prefix := calcPrivacyUTXOPrefix4Addr(assetExec, token, addr)
	values := list.List(prefix, nil, 0, 0)
	wutxos := new(walletUTXOs)
	if len(values) == 0 {
		return wutxos, nil
	}
	for _, value := range values {
		if len(value) == 0 {
			continue
		}
		accByte, err := store.Get(value)
		if err != nil {
			return nil, types.ErrDataBaseDamage
		}
		privacyDBStore := new(privacytypes.PrivacyDBStore)
		err = types.Decode(accByte, privacyDBStore)
		if err != nil {
			bizlog.Error("getPrivacyTokenUTXOs", "decode PrivacyDBStore error. ", err)
			return nil, types.ErrDataBaseDamage
		}
		wutxo := &walletUTXO{
			height: privacyDBStore.Height,
			outinfo: &txOutputInfo{
				amount:           privacyDBStore.Amount,
				txPublicKeyR:     privacyDBStore.TxPublicKeyR,
				onetimePublicKey: privacyDBStore.OnetimePublicKey,
				utxoGlobalIndex: &privacytypes.UTXOGlobalIndex{
					Outindex: privacyDBStore.OutIndex,
					Txhash:   privacyDBStore.Txhash,
				},
			},
		}
		wutxos.utxos = append(wutxos.utxos, wutxo)
	}
	return wutxos, nil
}

//calcUTXOKey4TokenAddr---X--->calcUTXOKey  toke utx 
//calcKey4UTXOsSpentInTx------>types.FTXOsSTXOsInOneTx utx  ftxo has 
//calcKey4FTXOsInTx----------->calcKey4UTXOsSpentInTx utx 
/  utx ftxo t utxo FTX STXO
func (store *privacyStore) moveUTXO2FTXO(expire int64, assetExec, token, sender, txhash string, selectedUtxos []*txOutputInfo) {
	FTXOsInOneTx := &privacytypes.FTXOsSTXOsInOneTx{
		AssetExec: assetExec,
		Tokenname: token,
		Sender:    sender,
		Txhash:    txhash,
	}
	newbatch := store.NewBatch(true)
	for _, txOutputInfo := range selectedUtxos {
		key := calcUTXOKey4TokenAddr(assetExec, token, sender, hex.EncodeToString(txOutputInfo.utxoGlobalIndex.Txhash), int(txOutputInfo.utxoGlobalIndex.Outindex))
		newbatch.Delete(key)
		utxo := &privacytypes.UTXO{
			Amount: txOutputInfo.amount,
			UtxoBasic: &privacytypes.UTXOBasic{
				UtxoGlobalIndex: txOutputInfo.utxoGlobalIndex,
				OnetimePubkey:   txOutputInfo.onetimePublicKey,
			},
		}
		FTXOsInOneTx.Utxos = append(FTXOsInOneTx.Utxos, utxo)
	}
	FTXOsInOneTx.SetExpire(expire)
	/ UTXO
	key1 := calcKey4UTXOsSpentInTx(txhash)
	value1 := types.Encode(FTXOsInOneTx)
	newbatch.Set(key1, value1)

	/ ftx key utxo
	key2 := calcKey4FTXOsInTx(assetExec, token, sender, txhash)
	value2 := key1
	newbatch.Set(key2, value2)

	newbatch.Write()
}

func (store *privacyStore) getRescanUtxosFlag4Addr(req *privacytypes.ReqRescanUtxos) (*privacytypes.RepRescanUtxos, error) {
	var storeAddrs []string
	if len(req.Addrs) == 0 {
		WalletAccStores, err := store.getAccountByPrefix("Account")
		if err != nil || len(WalletAccStores) == 0 {
			bizlog.Info("getRescanUtxosFlag4Addr", "GetAccountByPrefix:err", err)
			return nil, types.ErrNotFound
		}
		for _, WalletAccStore := range WalletAccStores {
			storeAddrs = append(storeAddrs, WalletAccStore.Addr)
		}
	} else {
		storeAddrs = append(storeAddrs, req.Addrs...)
	}

	var repRescanUtxos privacytypes.RepRescanUtxos
	for _, addr := range storeAddrs {
		value, err := store.Get(calcRescanUtxosFlagKey(addr))
		if err != nil {
			bizlog.Error("getRescanUtxosFlag4Addr", "Failed to get calcRescanUtxosFlagKey(addr) for value", addr)
			continue
		}

		var data types.Int64
		err = types.Decode(value, &data)
		if nil != err {
			bizlog.Error("getRescanUtxosFlag4Addr", "Failed to decode types.Int64 for value", value)
			continue
		}
		result := &privacytypes.RepRescanResult{
			Addr: addr,
			Flag: int32(data.Data),
		}
		repRescanUtxos.RepRescanResults = append(repRescanUtxos.RepRescanResults, result)
	}

	if len(repRescanUtxos.RepRescanResults) == 0 {
		return nil, types.ErrNotFound
	}

	repRescanUtxos.Flag = req.Flag

	return &repRescanUtxos, nil
}

func (store *privacyStore) saveREscanUTXOsAddresses(addrs []string, scanFlag int32) {
	newbatch := store.NewBatch(true)
	for _, addr := range addrs {
		data := &types.Int64{
			Data: int64(scanFlag),
		}
		value := types.Encode(data)
		newbatch.Set(calcRescanUtxosFlagKey(addr), value)
	}
	newbatch.Write()
}

func (store *privacyStore) setScanPrivacyInputUTXO(count int32) []*privacytypes.UTXOGlobalIndex {
	prefix := []byte(ScanPrivacyInput)
	list := store.NewListHelper()
	values := list.List(prefix, nil, count, 0)
	var utxoGlobalIndexs []*privacytypes.UTXOGlobalIndex
	if len(values) != 0 {
		var utxoGlobalIndex privacytypes.UTXOGlobalIndex
		for _, value := range values {
			err := types.Decode(value, &utxoGlobalIndex)
			if err == nil {
				utxoGlobalIndexs = append(utxoGlobalIndexs, &utxoGlobalIndex)
			}
		}
	}
	return utxoGlobalIndexs
}

func (store *privacyStore) isUTXOExist(txhash string, outindex int) (*privacytypes.PrivacyDBStore, error) {
	value1, err := store.Get(calcUTXOKey(txhash, outindex))
	if err != nil {
		bizlog.Error("IsUTXOExist", "Get calcUTXOKey error:", err)
		return nil, err
	}
	var accPrivacy privacytypes.PrivacyDBStore
	err = proto.Unmarshal(value1, &accPrivacy)
	if err != nil {
		bizlog.Error("IsUTXOExist", "proto.Unmarshal err:", err)
		return nil, err
	}
	return &accPrivacy, nil
}

func (store *privacyStore) updateScanInputUTXOs(utxoGlobalIndexs []*privacytypes.UTXOGlobalIndex) {
	if len(utxoGlobalIndexs) <= 0 {
		return
	}
	newbatch := store.NewBatch(true)
	var utxos []*privacytypes.UTXO
	var owner, token, txhash, assetExec string
	for _, utxoGlobal := range utxoGlobalIndexs {
		accPrivacy, err := store.isUTXOExist(hex.EncodeToString(utxoGlobal.Txhash), int(utxoGlobal.Outindex))
		if err == nil && accPrivacy != nil {
			utxo := &privacytypes.UTXO{
				Amount: accPrivacy.Amount,
				UtxoBasic: &privacytypes.UTXOBasic{
					UtxoGlobalIndex: utxoGlobal,
					OnetimePubkey:   accPrivacy.OnetimePublicKey,
				},
			}
			utxos = append(utxos, utxo)
			owner = accPrivacy.Owner
			token = accPrivacy.Tokenname
			txhash = hex.EncodeToString(accPrivacy.Txhash)
		}
		key := calcScanPrivacyInputUTXOKey(hex.EncodeToString(utxoGlobal.Txhash), int(utxoGlobal.Outindex))
		newbatch.Delete(key)
	}
	if len(utxos) > 0 {
		store.moveUTXO2STXO(assetExec, owner, token, txhash, utxos, newbatch)
	}
	newbatch.Write()
}

func (store *privacyStore) moveUTXO2STXO(assetExec, owner, token, txhash string, utxos []*privacytypes.UTXO, newbatch db.Batch) {
	if len(utxos) == 0 {
		return
	}

	FTXOsInOneTx := &privacytypes.FTXOsSTXOsInOneTx{
		AssetExec: assetExec,
		Tokenname: token,
		Sender:    owner,
		Txhash:    txhash,
		Utxos:     utxos,
	}

	for _, utxo := range utxos {
		Txhash := utxo.UtxoBasic.UtxoGlobalIndex.Txhash
		Outindex := utxo.UtxoBasic.UtxoGlobalIndex.Outindex
		/ UTX 
		key := calcUTXOKey4TokenAddr(assetExec, token, owner, hex.EncodeToString(Txhash), int(Outindex))
		newbatch.Delete(key)
	}

	/ UTXO
	key1 := calcKey4UTXOsSpentInTx(txhash)
	value1 := types.Encode(FTXOsInOneTx)
	newbatch.Set(key1, value1)

	/ stx key
	key2 := calcKey4STXOsInTx(txhash)
	value2 := key1
	newbatch.Set(key2, value2)

	// stxo-token-addr-txhash key 
	key3 := calcSTXOTokenAddrTxKey(assetExec, FTXOsInOneTx.Tokenname, FTXOsInOneTx.Sender, FTXOsInOneTx.Txhash)
	value3 := key1
	newbatch.Set(key3, value3)

	bizlog.Info("moveUTXO2STXO", "tx hash", txhash)
}

func (store *privacyStore) selectPrivacyTransactionToWallet(txDetals *types.TransactionDetails, privacyInfo []addrAndprivacy) {
	newbatch := store.NewBatch(true)
	for _, txdetal := range txDetals.Txs {
		if !bytes.Equal([]byte(privacytypes.PrivacyX), txdetal.Tx.Execer) {
			continue
		}
		store.selectCurrentWalletPrivacyTx(txdetal, int32(txdetal.Index), privacyInfo, newbatch)
	}
	newbatch.Write()
}

func (store *privacyStore) selectCurrentWalletPrivacyTx(txDetal *types.TransactionDetail, index int32, privacyInfo []addrAndprivacy, newbatch db.Batch) {
	tx := txDetal.Tx
	amount, err := tx.Amount()
	if err != nil {
		bizlog.Error("selectCurrentWalletPrivacyTx failed to tx.Amount()")
		return
	}

	txExecRes := txDetal.Receipt.Ty
	height := txDetal.Height

	txhashInbytes := tx.Hash()
	txhash := hex.EncodeToString(txhashInbytes)
	var privateAction privacytypes.PrivacyAction
	if err := types.Decode(tx.GetPayload(), &privateAction); err != nil {
		bizlog.Error("selectCurrentWalletPrivacyTx failed to decode payload")
		return
	}
	bizlog.Info("selectCurrentWalletPrivacyTx", "tx hash", txhash)
	var RpubKey []byte
	var privacyOutput *privacytypes.PrivacyOutput
	var privacyInput *privacytypes.PrivacyInput
	var tokenname, assetExec string
	if privacytypes.ActionPublic2Privacy == privateAction.Ty {
		RpubKey = privateAction.GetPublic2Privacy().GetOutput().GetRpubKeytx()
		privacyOutput = privateAction.GetPublic2Privacy().GetOutput()
		tokenname = privateAction.GetPublic2Privacy().GetTokenname()
		assetExec = privateAction.GetPublic2Privacy().GetAssetExec()
	} else if privacytypes.ActionPrivacy2Privacy == privateAction.Ty {
		RpubKey = privateAction.GetPrivacy2Privacy().GetOutput().GetRpubKeytx()
		privacyOutput = privateAction.GetPrivacy2Privacy().GetOutput()
		tokenname = privateAction.GetPrivacy2Privacy().GetTokenname()
		privacyInput = privateAction.GetPrivacy2Privacy().GetInput()
		assetExec = privateAction.GetPrivacy2Privacy().GetAssetExec()
	} else if privacytypes.ActionPrivacy2Public == privateAction.Ty {
		RpubKey = privateAction.GetPrivacy2Public().GetOutput().GetRpubKeytx()
		privacyOutput = privateAction.GetPrivacy2Public().GetOutput()
		tokenname = privateAction.GetPrivacy2Public().GetTokenname()
		privacyInput = privateAction.GetPrivacy2Public().GetInput()
		assetExec = privateAction.GetPrivacy2Public().GetAssetExec()
	}

	if assetExec == "" {
		assetExec = store.Chain33Config.GetCoinExec()
	}

	/ output
	if nil != privacyOutput && len(privacyOutput.Keyoutput) > 0 {
		utxoProcessed := make([]bool, len(privacyOutput.Keyoutput))
		for _, info := range privacyInfo {
			bizlog.Debug("SelectCurrentWalletPrivacyTx", "individual privacyInfo's addr", *info.Addr)
			privacykeyParirs := info.PrivacyKeyPair
			bizlog.Debug("SelectCurrentWalletPrivacyTx", "individual ViewPubkey", hex.EncodeToString(privacykeyParirs.ViewPubkey.Bytes()),
				"individual SpendPubkey", hex.EncodeToString(privacykeyParirs.SpendPubkey.Bytes()))

			var utxos []*privacytypes.UTXO
			for indexoutput, output := range privacyOutput.Keyoutput {
				if utxoProcessed[indexoutput] {
					continue
				}
				priv, err := privacy.RecoverOnetimePriKey(RpubKey, privacykeyParirs.ViewPrivKey, privacykeyParirs.SpendPrivKey, int64(indexoutput))
				if err == nil {
					recoverPub := priv.PubKey().Bytes()[:]
					if bytes.Equal(recoverPub, output.Onetimepubkey) {
						/  
						/ ，
						//1   ，
						//2  change  
						utxoProcessed[indexoutput] = true
						bizlog.Debug("SelectCurrentWalletPrivacyTx got privacy tx belong to current wallet",
							"Address", *info.Addr, "tx with hash", txhash, "Amount", amount)
						/ UTX 
						if types.ExecOk == txExecRes {

							// UTX has  
							accPrivacy, err := store.isUTXOExist(hex.EncodeToString(txhashInbytes), indexoutput)
							if err == nil && accPrivacy != nil {
								continue
							}

							info2store := &privacytypes.PrivacyDBStore{
								AssetExec:        assetExec,
								Txhash:           txhashInbytes,
								Tokenname:        tokenname,
								Amount:           output.Amount,
								OutIndex:         int32(indexoutput),
								TxPublicKeyR:     RpubKey,
								OnetimePublicKey: output.Onetimepubkey,
								Owner:            *info.Addr,
								Height:           height,
								Txindex:          index,
								//Blockhash:        block.Block.Hash(),
							}

							utxoGlobalIndex := &privacytypes.UTXOGlobalIndex{
								Outindex: int32(indexoutput),
								Txhash:   txhashInbytes,
							}

							utxoCreated := &privacytypes.UTXO{
								Amount: output.Amount,
								UtxoBasic: &privacytypes.UTXOBasic{
									UtxoGlobalIndex: utxoGlobalIndex,
									OnetimePubkey:   output.Onetimepubkey,
								},
							}

							utxos = append(utxos, utxoCreated)
							store.setUTXO(info2store, txhash, newbatch)
						}
					}
				}
			}
		}
	}

	/ input
	if nil != privacyInput && len(privacyInput.Keyinput) > 0 {
		var utxoGlobalIndexs []*privacytypes.UTXOGlobalIndex
		for _, input := range privacyInput.Keyinput {
			utxoGlobalIndexs = append(utxoGlobalIndexs, input.UtxoGlobalIndex...)
		}

		if len(utxoGlobalIndexs) > 0 {
			store.storeScanPrivacyInputUTXO(utxoGlobalIndexs, newbatch)
		}
	}
}

// setUTXO UTX 
// addr UTX  UTX 
// txhash UTX  0x
// outindex UTX 
// dbStore UTX 
//UTXO---->moveUTXO2FTXO---->FTXO---->moveFTXO2STXO---->STXO
//1.calcUTXOKey------------>types.PrivacyDBStore k d  UTX 
//2.calcUTXOKey4TokenAddr-->calcUTXOKey kv toke utxo
func (store *privacyStore) setUTXO(utxoInfo *privacytypes.PrivacyDBStore, txHash string, newbatch db.Batch) error {

	privacyStorebyte := types.Encode(utxoInfo)
	outIndex := int(utxoInfo.OutIndex)
	utxoKey := calcUTXOKey(txHash, outIndex)
	bizlog.Debug("setUTXO", "addr", utxoInfo.Owner, "tx with hash", txHash, "amount:", utxoInfo.Amount/store.GetCoinPrecision())
	newbatch.Set(calcUTXOKey4TokenAddr(utxoInfo.AssetExec, utxoInfo.Tokenname, utxoInfo.Owner, txHash, outIndex), utxoKey)
	newbatch.Set(utxoKey, privacyStorebyte)
	return nil
}

func (store *privacyStore) storeScanPrivacyInputUTXO(utxoGlobalIndexs []*privacytypes.UTXOGlobalIndex, newbatch db.Batch) {
	for _, utxoGlobalIndex := range utxoGlobalIndexs {
		key1 := calcScanPrivacyInputUTXOKey(hex.EncodeToString(utxoGlobalIndex.Txhash), int(utxoGlobalIndex.Outindex))
		utxoIndex := &privacytypes.UTXOGlobalIndex{
			Txhash:   utxoGlobalIndex.Txhash,
			Outindex: utxoGlobalIndex.Outindex,
		}
		value1 := types.Encode(utxoIndex)
		newbatch.Set(key1, value1)
	}
}

func (store *privacyStore) listSpendUTXOs(assetExec, token, addr string) (*privacytypes.UTXOHaveTxHashs, error) {
	if 0 == len(addr) {
		bizlog.Error("listSpendUTXOs addr is nil")
		return nil, types.ErrInvalidParam
	}
	prefix := calcSTXOPrefix4Addr(assetExec, token, addr)
	list := store.NewListHelper()
	Key4FTXOsInTxs := list.PrefixScan(prefix)
	//if len(Key4FTXOsInTxs) == 0 {
	//	bizlog.Error("listSpendUTXOs ", "addr not exist", addr)
	//	return nil, types.ErrNotFound
	//}

	var utxoHaveTxHashs privacytypes.UTXOHaveTxHashs
	utxoHaveTxHashs.UtxoHaveTxHashs = make([]*privacytypes.UTXOHaveTxHash, 0)
	for _, Key4FTXOsInTx := range Key4FTXOsInTxs {
		value, err := store.Get(Key4FTXOsInTx)
		if err != nil {
			continue
		}
		var ftxosInOneTx privacytypes.FTXOsSTXOsInOneTx
		err = types.Decode(value, &ftxosInOneTx)
		if nil != err {
			bizlog.Error("listSpendUTXOs", "Failed to decode FTXOsSTXOsInOneTx for value", value)
			return nil, types.ErrInvalidParam
		}

		for _, ftxo := range ftxosInOneTx.Utxos {
			utxohash := hex.EncodeToString(ftxo.UtxoBasic.UtxoGlobalIndex.Txhash)
			value1, err := store.Get(calcUTXOKey(utxohash, int(ftxo.UtxoBasic.UtxoGlobalIndex.Outindex)))
			if err != nil {
				continue
			}
			var accPrivacy privacytypes.PrivacyDBStore
			err = proto.Unmarshal(value1, &accPrivacy)
			if err != nil {
				bizlog.Error("listWalletPrivacyAccount", "proto.Unmarshal err:", err)
				return nil, types.ErrUnmarshal
			}

			utxoBasic := &privacytypes.UTXOBasic{
				UtxoGlobalIndex: &privacytypes.UTXOGlobalIndex{
					Outindex: accPrivacy.OutIndex,
					Txhash:   accPrivacy.Txhash,
				},
				OnetimePubkey: accPrivacy.OnetimePublicKey,
			}

			var utxoHaveTxHash privacytypes.UTXOHaveTxHash
			utxoHaveTxHash.Amount = accPrivacy.Amount
			utxoHaveTxHash.TxHash = ftxosInOneTx.Txhash
			utxoHaveTxHash.UtxoBasic = utxoBasic

			utxoHaveTxHashs.UtxoHaveTxHashs = append(utxoHaveTxHashs.UtxoHaveTxHashs, &utxoHaveTxHash)
		}
	}
	return &utxoHaveTxHashs, nil
}

func (store *privacyStore) getWalletFtxoStxo(prefix string) ([]*privacytypes.FTXOsSTXOsInOneTx, []string, error) {
	list := store.NewListHelper()
	values := list.List([]byte(prefix), nil, 0, 0)
	var Ftxoes []*privacytypes.FTXOsSTXOsInOneTx
	var key []string
	for _, value := range values {
		value1, err := store.Get(value)
		if err != nil {
			continue
		}

		FTXOsInOneTx := &privacytypes.FTXOsSTXOsInOneTx{}
		err = types.Decode(value1, FTXOsInOneTx)
		if nil != err {
			bizlog.Error("DecodeString Error", "Error", err.Error())
			return nil, nil, err
		}

		Ftxoes = append(Ftxoes, FTXOsInOneTx)
		key = append(key, string(value))
	}
	return Ftxoes, key, nil
}

func (store *privacyStore) getFTXOlist() ([]*privacytypes.FTXOsSTXOsInOneTx, [][]byte) {
	curFTXOTxs, _, _ := store.getWalletFtxoStxo(FrozenUTXOs)
	revertFTXOTxs, _, _ := store.getWalletFtxoStxo(RevertSendtx)
	var keys [][]byte
	for _, ftxo := range curFTXOTxs {
		keys = append(keys, calcKey4FTXOsInTx(ftxo.AssetExec, ftxo.Tokenname, ftxo.Sender, ftxo.Txhash))
	}
	for _, ftxo := range revertFTXOTxs {
		keys = append(keys, calcRevertSendTxKey(ftxo.AssetExec, ftxo.Tokenname, ftxo.Sender, ftxo.Txhash))
	}
	curFTXOTxs = append(curFTXOTxs, revertFTXOTxs...)
	return curFTXOTxs, keys
}

//calcKey4FTXOsInTx-----x------>calcKey4UTXOsSpentInTx ，
//calcKey4STXOsInTx------------>calcKey4UTXOsSpentInTx
/ types.FTXOsSTXOsInOneT 
func (store *privacyStore) moveFTXO2STXO(key1 []byte, txhash string, newbatch db.Batch) error {
	/ UTXO
	value1, err := store.Get(key1)
	if err != nil {
		bizlog.Error("moveFTXO2STXO", "Get(key1) error ", err, "key1", string(key1))
		return err
	}
	if value1 == nil {
		bizlog.Error("moveFTXO2STXO", "Get nil value for txhash", txhash)
		return types.ErrNotFound
	}
	newbatch.Delete(key1)

	/ ftx key utxo
	key2 := calcKey4STXOsInTx(txhash)
	value2 := value1
	newbatch.Set(key2, value2)

	// stxo-token-addr-txhash key 
	key := value2
	value, err := store.Get(key)
	if err != nil {
		bizlog.Error("moveFTXO2STXO", "db Get(key) error ", err, "key", key)
	}
	var ftxosInOneTx privacytypes.FTXOsSTXOsInOneTx
	err = types.Decode(value, &ftxosInOneTx)
	if nil != err {
		bizlog.Error("moveFTXO2STXO", "Failed to decode FTXOsSTXOsInOneTx for value", value)
	}
	key3 := calcSTXOTokenAddrTxKey(ftxosInOneTx.AssetExec, ftxosInOneTx.Tokenname, ftxosInOneTx.Sender, ftxosInOneTx.Txhash)
	newbatch.Set(key3, value2)
	newbatch.Write()

	bizlog.Info("moveFTXO2STXO", "tx hash", txhash)
	return nil
}

/ FTX UTXO
// moveFTXO2UTXO  UTX UTX 
// UTX   UTX  UTX 
func (store *privacyStore) moveFTXO2UTXO(key1 []byte, newbatch db.Batch) {
	/ ftx key utxo
	value1, err := store.Get(key1)
	if err != nil {
		bizlog.Error("moveFTXO2UTXO", "db Get(key1) error ", err)
		return
	}
	if nil == value1 {
		bizlog.Error("moveFTXO2UTXO", "Get nil value for key", string(key1))
		return

	}
	newbatch.Delete(key1)

	key2 := value1
	value2, err := store.Get(key2)
	if err != nil {
		bizlog.Error("moveFTXO2UTXO", "db Get(key2) error ", err)
		return
	}
	if nil == value2 {
		bizlog.Error("moveFTXO2UTXO", "Get nil value for key", string(key2))
		return
	}
	newbatch.Delete(key2)

	var ftxosInOneTx privacytypes.FTXOsSTXOsInOneTx
	err = types.Decode(value2, &ftxosInOneTx)
	if nil != err {
		bizlog.Error("moveFTXO2UTXO", "Failed to decode FTXOsSTXOsInOneTx for value", value2)
		return
	}
	for _, ftxo := range ftxosInOneTx.Utxos {
		utxohash := hex.EncodeToString(ftxo.UtxoBasic.UtxoGlobalIndex.Txhash)
		outindex := int(ftxo.UtxoBasic.UtxoGlobalIndex.Outindex)
		key := calcUTXOKey4TokenAddr(ftxosInOneTx.AssetExec, ftxosInOneTx.Tokenname, ftxosInOneTx.Sender, utxohash, outindex)
		value := calcUTXOKey(utxohash, int(ftxo.UtxoBasic.UtxoGlobalIndex.Outindex))
		bizlog.Debug("moveFTXO2UTXO", "addr", ftxosInOneTx.Sender, "tx with hash", utxohash, "amount", ftxo.Amount/store.GetCoinPrecision())
		newbatch.Set(key, value)
	}
	bizlog.Debug("moveFTXO2UTXO", "addr", ftxosInOneTx.Sender, "tx with hash", ftxosInOneTx.Txhash)
}

// unsetUTXO   UTX  
// 1 UTX UTX 
// 2 UTX UTX 
// 3 UTX UTX UTX 
// 4 UTX 
// addr UTX 
// txhash UTX  0x
func (store *privacyStore) unsetUTXO(assetExec, token, addr, txhash string, outindex int, newbatch db.Batch) error {
	if 0 == len(addr) || 0 == len(txhash) || outindex < 0 || len(token) <= 0 {
		bizlog.Error("unsetUTXO", "InvalidParam addr", addr, "txhash", txhash, "outindex", outindex, "token", token)
		return types.ErrInvalidParam
	}
	// 1 UTX 
	ftxokey := calcUTXOKey(txhash, outindex)
	newbatch.Delete(ftxokey)
	// 2 UTX 
	ftxokey = calcKey4FTXOsInTx(assetExec, token, addr, txhash)
	newbatch.Delete(ftxokey)
	// 3 UTX 
	ftxokey = calcRevertSendTxKey(assetExec, token, addr, txhash)
	newbatch.Delete(ftxokey)
	// 4 UTX 
	utxokey := calcUTXOKey4TokenAddr(assetExec, token, addr, txhash, outindex)
	newbatch.Delete(utxokey)

	bizlog.Debug("PrivacyTrading unsetUTXO", "addr", addr, "tx with hash", txhash, "outindex", outindex)
	return nil
}

/   stx ftxo，
/  
func (store *privacyStore) moveSTXO2FTXO(tx *types.Transaction, txhash string, newbatch db.Batch) error {
	/ ftx key utxo
	key2 := calcKey4STXOsInTx(txhash)
	value2, err := store.Get(key2)
	if err != nil {
		bizlog.Error("moveSTXO2FTXO", "Get(key2) error ", err)
		return err
	}
	if value2 == nil {
		bizlog.Debug("moveSTXO2FTXO", "Get nil value for txhash", txhash)
		return types.ErrNotFound
	}
	newbatch.Delete(key2)

	key := value2
	value, err := store.Get(key)
	if err != nil {
		bizlog.Error("moveSTXO2FTXO", "db Get(key) error ", err)
	}

	var ftxosInOneTx privacytypes.FTXOsSTXOsInOneTx
	err = types.Decode(value, &ftxosInOneTx)
	if nil != err {
		bizlog.Error("moveSTXO2FTXO", "Failed to decode FTXOsSTXOsInOneTx for value", value)
	}

	/ stxo-token-addr-txhash key
	key3 := calcSTXOTokenAddrTxKey(ftxosInOneTx.AssetExec, ftxosInOneTx.Tokenname, ftxosInOneTx.Sender, ftxosInOneTx.Txhash)
	newbatch.Delete(key3)

	/ UTXO
	key1 := calcRevertSendTxKey(ftxosInOneTx.AssetExec, ftxosInOneTx.Tokenname, ftxosInOneTx.Sender, txhash)
	value1 := value2
	newbatch.Set(key1, value1)
	bizlog.Info("moveSTXO2FTXO", "txhash ", txhash)

	ftxosInOneTx.SetExpire(tx.GetExpire())
	value = types.Encode(&ftxosInOneTx)
	newbatch.Set(key, value)

	newbatch.Write()
	return nil
}
